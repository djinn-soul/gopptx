package chart

import (
	"errors"
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func DetectChartKind(chartXML []byte) Kind {
	s := string(chartXML)
	switch {
	case strings.Contains(s, "<c:bubbleChart"):
		return KindBubble
	case strings.Contains(s, "<c:scatterChart"):
		return KindScatter
	default:
		return KindCategory
	}
}

func ValidateChartUpdatePayload(kind Kind, req common.ChartDataUpdate) error {
	if len(req.Series) == 0 {
		return errors.New("chart update requires at least one series")
	}

	var err error
	switch kind {
	case KindCategory:
		err = validateCategoryChartUpdatePayload(req)
	case KindScatter:
		err = validateScatterChartUpdatePayload(req.Series)
	case KindBubble:
		err = validateBubbleChartUpdatePayload(req.Series)
	default:
		err = errors.New("unsupported chart type")
	}
	return err
}

func ResolveChartSelector(
	refs []common.SlideChartRef,
	selector common.ChartSelector,
	slideIndex int,
) (common.SlideChartRef, error) {
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
