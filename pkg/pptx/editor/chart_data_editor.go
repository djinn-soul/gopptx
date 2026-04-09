package editor

import (
	"bytes"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodchart "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/chart"
)

func (e *PresentationEditor) ListSlideCharts(slideIndex int) ([]common.SlideChartRef, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, fmt.Errorf("slide index %d out of range", slideIndex)
	}

	slideRef := e.slides[slideIndex]
	slideXML, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return nil, fmt.Errorf("slide part %s not found", slideRef.Part)
	}
	relsPath := common.RelsPathFor(slideRef.Part)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return nil, fmt.Errorf("rels part %s not found", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return nil, fmt.Errorf("parse slide rels: %w", err)
	}

	relToChart := make(map[string]string)
	for _, rel := range rels {
		if rel.Type != common.RelTypeChart {
			continue
		}
		relToChart[rel.ID] = common.ResolveRelationshipTarget(slideRef.Part, rel.Target)
	}

	rIDs := extractChartRelIDs(slideXML)
	out := make([]common.SlideChartRef, 0, len(rIDs))
	for i, relID := range rIDs {
		chartPart, found := relToChart[relID]
		if !found {
			continue
		}
		out = append(out, common.SlideChartRef{
			Index:     i,
			RelID:     relID,
			ChartPart: chartPart,
		})
	}
	return out, nil
}

func (e *PresentationEditor) UpdateChartData(
	slideIndex int,
	selector common.ChartSelector,
	req common.ChartDataUpdate,
) error {
	refs, err := e.ListSlideCharts(slideIndex)
	if err != nil {
		return err
	}
	chartRef, err := editormodchart.ResolveChartSelector(refs, selector, slideIndex)
	if err != nil {
		return err
	}
	return e.applyChartDataUpdateByRef(chartRef, req)
}

func (e *PresentationEditor) UpdateChartDataBatch(
	slideIndex int,
	updates []common.ChartDataBatchItem,
) error {
	if len(updates) == 0 {
		return nil
	}

	refs, err := e.ListSlideCharts(slideIndex)
	if err != nil {
		return err
	}
	refsByIndex := make(map[int]common.SlideChartRef, len(refs))
	refsByRelID := make(map[string]common.SlideChartRef, len(refs))
	for _, ref := range refs {
		refsByIndex[ref.Index] = ref
		if ref.RelID != "" {
			refsByRelID[ref.RelID] = ref
		}
	}

	for _, item := range updates {
		chartRef, resolveErr := resolveChartRefFast(
			item.ChartSelector,
			refsByIndex,
			refsByRelID,
			slideIndex,
		)
		if resolveErr != nil {
			return resolveErr
		}
		if applyErr := e.applyChartDataUpdateByRef(chartRef, item.Data); applyErr != nil {
			return applyErr
		}
	}
	return nil
}

func resolveChartRefFast(
	selector common.ChartSelector,
	refsByIndex map[int]common.SlideChartRef,
	refsByRelID map[string]common.SlideChartRef,
	slideIndex int,
) (common.SlideChartRef, error) {
	if selector.Index == nil && selector.RelID == "" {
		return common.SlideChartRef{}, fmt.Errorf(
			"slide %d: chart selector requires index or rel_id",
			slideIndex,
		)
	}
	if selector.RelID != "" {
		ref, ok := refsByRelID[selector.RelID]
		if !ok {
			return common.SlideChartRef{}, fmt.Errorf(
				"slide %d: chart with rel_id %q not found",
				slideIndex,
				selector.RelID,
			)
		}
		if selector.Index != nil && ref.Index != *selector.Index {
			return common.SlideChartRef{}, fmt.Errorf(
				"slide %d: selector mismatch index=%d rel_id=%s",
				slideIndex,
				*selector.Index,
				selector.RelID,
			)
		}
		return ref, nil
	}
	ref, ok := refsByIndex[*selector.Index]
	if !ok {
		return common.SlideChartRef{}, fmt.Errorf(
			"slide %d: chart with index %d not found",
			slideIndex,
			*selector.Index,
		)
	}
	return ref, nil
}

func (e *PresentationEditor) applyChartDataUpdateByRef(
	chartRef common.SlideChartRef,
	req common.ChartDataUpdate,
) error {
	chartXML, ok := e.parts.Get(chartRef.ChartPart)
	if !ok {
		return fmt.Errorf("chart part %s not found", chartRef.ChartPart)
	}

	kind := editormodchart.DetectChartKind(chartXML)
	if validateErr := editormodchart.ValidateChartUpdatePayload(kind, req); validateErr != nil {
		return validateErr
	}
	if err := e.ensureChartEmbeddingRel(chartRef.ChartPart); err != nil {
		return err
	}

	workbook, err := editormodchart.GenerateExcelForChartUpdate(kind, req)
	if err != nil {
		return fmt.Errorf("generate excel: %w", err)
	}

	// Use registerExcelEmbedding to handle deduplication and avoid side effects
	// if the Excel part is shared between multiple charts.
	excelPartPath, err := e.registerExcelEmbedding(workbook)
	if err != nil {
		return err
	}

	// Ensure the chart points to the correct Excel part.
	if err := e.updateChartEmbeddingRel(chartRef.ChartPart, excelPartPath); err != nil {
		return err
	}

	e.chartEmbeddings[chartRef.ChartPart] = excelPartPath

	patchedChartXML, err := editormodchart.PatchChartDataCache(chartXML, kind, req)
	if err != nil {
		return err
	}
	e.parts.Set(chartRef.ChartPart, patchedChartXML)
	return nil
}

func (e *PresentationEditor) ensureChartEmbeddingRel(chartPart string) error {
	relsPath := common.RelsPathFor(chartPart)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return fmt.Errorf("chart rels part not found: %s", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return fmt.Errorf("parse chart rels: %w", err)
	}
	for _, rel := range rels {
		if rel.Type == common.RelTypePackage {
			return nil
		}
	}
	return fmt.Errorf("chart embedding relationship not found: %s", relsPath)
}

// UpdateChartFormatting applies a partial formatting patch to an existing chart.
func (e *PresentationEditor) UpdateChartFormatting(
	slideIndex int,
	selector common.ChartSelector,
	req common.ChartFormatUpdate,
) error {
	refs, err := e.ListSlideCharts(slideIndex)
	if err != nil {
		return err
	}
	chartRef, err := editormodchart.ResolveChartSelector(refs, selector, slideIndex)
	if err != nil {
		return err
	}
	if err := editormodchart.ValidateChartFormatUpdate(req); err != nil {
		return err
	}

	chartXML, ok := e.parts.Get(chartRef.ChartPart)
	if !ok {
		return fmt.Errorf("chart part %s not found", chartRef.ChartPart)
	}
	patched, err := editormodchart.PatchChartFormatting(chartXML, req)
	if err != nil {
		return err
	}
	e.parts.Set(chartRef.ChartPart, patched)
	return nil
}

// GetChartState returns a read-only chart snapshot suitable for object-model traversal.
func (e *PresentationEditor) GetChartState(
	slideIndex int,
	selector common.ChartSelector,
) (common.ChartState, error) {
	refs, err := e.ListSlideCharts(slideIndex)
	if err != nil {
		return common.ChartState{}, err
	}
	chartRef, err := editormodchart.ResolveChartSelector(refs, selector, slideIndex)
	if err != nil {
		return common.ChartState{}, err
	}

	chartXML, ok := e.parts.Get(chartRef.ChartPart)
	if !ok {
		return common.ChartState{}, fmt.Errorf("chart part %s not found", chartRef.ChartPart)
	}
	return editormodchart.ExtractChartState(chartXML), nil
}

func (e *PresentationEditor) updateChartEmbeddingRel(chartPart, excelPath string) error {
	relsPath := common.RelsPathFor(chartPart)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return fmt.Errorf("chart rels part not found: %s", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return fmt.Errorf("parse chart rels: %w", err)
	}

	changed := false
	for i, rel := range rels {
		if rel.Type != common.RelTypePackage {
			continue
		}
		newRelTarget := common.MakeRelativePath(chartPart, excelPath)
		if rels[i].Target != newRelTarget {
			rels[i].Target = newRelTarget
			changed = true
		}
	}

	if changed {
		newRelsData := renderRelationshipsXML(rels)
		e.parts.Set(relsPath, []byte(newRelsData))
	}
	return nil
}

func (e *PresentationEditor) writeRelationships(
	path string,
	rels []common.EditorRelationship,
) error {
	e.parts.Set(path, []byte(renderRelationshipsXML(rels)))
	return nil
}

func (e *PresentationEditor) addContentTypeOverride(partName, contentType string) {
	ctPath := "[Content_Types].xml"
	data, ok := e.parts.Get(ctPath)
	if !ok {
		return
	}

	partNameRooted := "/" + partName
	if bytes.Contains(data, []byte(`PartName="`+partNameRooted+`"`)) {
		return
	}

	override := fmt.Sprintf(
		`<Override PartName="%s" ContentType="%s"/>`,
		partNameRooted,
		contentType,
	)
	replaced := bytes.Replace(data, []byte("</Types>"), []byte(override+"</Types>"), 1)
	e.parts.Set(ctPath, replaced)
}
