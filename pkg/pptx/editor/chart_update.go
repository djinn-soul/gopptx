package editor

import (
	"errors"
	"fmt"
	"path"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodchart "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/chart"
)

type chartKind int

const (
	chartKindCategory chartKind = iota
	chartKindScatter
	chartKindBubble
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
		relToChart[rel.ID] = resolveRelationshipTarget(slideRef.Part, rel.Target)
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
	chartRef, err := e.resolveChartSelector(slideIndex, selector)
	if err != nil {
		return err
	}

	chartXML, ok := e.parts.Get(chartRef.ChartPart)
	if !ok {
		return fmt.Errorf("chart part %s not found", chartRef.ChartPart)
	}

	kind := detectChartKind(chartXML)
	if validateErr := validateChartUpdatePayload(kind, req); validateErr != nil {
		return validateErr
	}

	workbook, err := editormodchart.GenerateExcelForChartUpdate(chartKindToModuleKind(kind), req)
	if err != nil {
		return fmt.Errorf("generate excel: %w", err)
	}

	excelPartPath, err := e.resolveChartEmbeddingPath(chartRef.ChartPart)
	if err != nil {
		return err
	}
	e.parts.Set(excelPartPath, workbook)
	e.chartEmbeddings[chartRef.ChartPart] = excelPartPath

	patchedChartXML, err := patchChartDataCache(chartXML, kind, req)
	if err != nil {
		return err
	}
	e.parts.Set(chartRef.ChartPart, patchedChartXML)
	return nil
}

func (e *PresentationEditor) resolveChartSelector(
	slideIndex int,
	selector common.ChartSelector,
) (common.SlideChartRef, error) {
	refs, err := e.ListSlideCharts(slideIndex)
	if err != nil {
		return common.SlideChartRef{}, err
	}
	if len(refs) == 0 {
		return common.SlideChartRef{}, fmt.Errorf("slide %d has no charts", slideIndex)
	}

	var idxMatch *common.SlideChartRef
	if selector.Index != nil {
		if *selector.Index < 0 || *selector.Index >= len(refs) {
			return common.SlideChartRef{}, fmt.Errorf(
				"chart index %d out of range (found %d charts)",
				*selector.Index,
				len(refs),
			)
		}
		ref := refs[*selector.Index]
		idxMatch = &ref
	}

	relID := strings.TrimSpace(selector.RelID)
	var relMatch *common.SlideChartRef
	if relID != "" {
		for i := range refs {
			if refs[i].RelID == relID {
				ref := refs[i]
				relMatch = &ref
				break
			}
		}
		if relMatch == nil {
			return common.SlideChartRef{}, fmt.Errorf("chart rel_id %q not found on slide %d", relID, slideIndex)
		}
	}

	if idxMatch != nil && relMatch != nil {
		if idxMatch.RelID != relMatch.RelID {
			return common.SlideChartRef{}, fmt.Errorf(
				"chart selector mismatch: index=%d rel_id=%q",
				*selector.Index,
				selector.RelID,
			)
		}
		return *idxMatch, nil
	}
	if idxMatch != nil {
		return *idxMatch, nil
	}
	if relMatch != nil {
		return *relMatch, nil
	}
	return common.SlideChartRef{}, errors.New("chart_selector must include index and/or rel_id")
}

func detectChartKind(chartXML []byte) chartKind {
	s := string(chartXML)
	switch {
	case strings.Contains(s, "<c:bubbleChart"):
		return chartKindBubble
	case strings.Contains(s, "<c:scatterChart"):
		return chartKindScatter
	default:
		return chartKindCategory
	}
}

func chartKindToModuleKind(kind chartKind) editormodchart.Kind {
	switch kind {
	case chartKindScatter:
		return editormodchart.KindScatter
	case chartKindBubble:
		return editormodchart.KindBubble
	default:
		return editormodchart.KindCategory
	}
}

func validateChartUpdatePayload(kind chartKind, req common.ChartDataUpdate) error {
	if len(req.Series) == 0 {
		return errors.New("chart update requires at least one series")
	}

	var err error
	switch kind {
	case chartKindCategory:
		err = validateCategoryChartUpdatePayload(req)
	case chartKindScatter:
		err = validateScatterChartUpdatePayload(req.Series)
	case chartKindBubble:
		err = validateBubbleChartUpdatePayload(req.Series)
	default:
		err = errors.New("unsupported chart type")
	}
	return err
}

func validateCategoryChartUpdatePayload(req common.ChartDataUpdate) error {
	baseCats := len(req.Categories)
	if baseCats == 0 {
		baseCats = firstSeriesCategoryCount(req.Series)
	}

	for i, s := range req.Series {
		if err := validateCategorySeries(i, s, baseCats); err != nil {
			return err
		}
	}
	return nil
}

func firstSeriesCategoryCount(series []common.ChartSeriesData) int {
	if len(series) == 0 {
		return 0
	}
	return len(series[0].Categories)
}

func validateCategorySeries(index int, series common.ChartSeriesData, baseCats int) error {
	if len(series.Values) == 0 {
		return fmt.Errorf("series %d requires values", index)
	}

	catLen := baseCats
	if len(series.Categories) > 0 {
		catLen = len(series.Categories)
	}
	if catLen == 0 {
		return errors.New("category chart requires categories")
	}
	if len(series.Values) == catLen {
		return nil
	}

	return fmt.Errorf(
		"series %d values length (%d) must equal category length (%d)",
		index,
		len(series.Values),
		catLen,
	)
}

func validateScatterChartUpdatePayload(series []common.ChartSeriesData) error {
	for i, s := range series {
		if len(s.XValues) == 0 || len(s.YValues) == 0 {
			return fmt.Errorf("scatter series %d requires x_values and y_values", i)
		}
		if len(s.XValues) != len(s.YValues) {
			return fmt.Errorf("scatter series %d x/y length mismatch", i)
		}
	}
	return nil
}

func validateBubbleChartUpdatePayload(series []common.ChartSeriesData) error {
	for i, s := range series {
		if len(s.XValues) == 0 || len(s.YValues) == 0 || len(s.Sizes) == 0 {
			return fmt.Errorf("bubble series %d requires x_values, y_values, and sizes", i)
		}
		if len(s.XValues) != len(s.YValues) || len(s.XValues) != len(s.Sizes) {
			return fmt.Errorf("bubble series %d x/y/size length mismatch", i)
		}
	}
	return nil
}

func (e *PresentationEditor) resolveChartEmbeddingPath(chartPart string) (string, error) {
	if excelPath, ok := e.chartEmbeddings[chartPart]; ok {
		return excelPath, nil
	}

	relsPath := common.RelsPathFor(chartPart)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return "", fmt.Errorf("chart rels part not found: %s", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return "", fmt.Errorf("parse chart rels: %w", err)
	}
	for _, rel := range rels {
		if rel.Type != common.RelTypePackage {
			continue
		}
		excelPath := common.CanonicalPartPath(path.Join(path.Dir(chartPart), rel.Target))
		if !e.parts.Has(excelPath) {
			return "", fmt.Errorf("chart embedding part missing: %s", excelPath)
		}
		e.chartEmbeddings[chartPart] = excelPath
		return excelPath, nil
	}
	return "", fmt.Errorf("associated excel part not found for chart %s", chartPart)
}
