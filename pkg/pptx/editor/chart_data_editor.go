package editor

import (
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

	chartXML, ok := e.parts.Get(chartRef.ChartPart)
	if !ok {
		return fmt.Errorf("chart part %s not found", chartRef.ChartPart)
	}

	kind := editormodchart.DetectChartKind(chartXML)
	if validateErr := editormodchart.ValidateChartUpdatePayload(kind, req); validateErr != nil {
		return validateErr
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
