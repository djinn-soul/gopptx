package command

import common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"

type OptionalInt64FieldFn func(map[string]any, string) (int64, bool)
type ParseStringSliceFieldFn func(map[string]any, string) ([]string, bool)
type ParseFloat64SliceFieldFn func(map[string]any, string) ([]float64, bool)

type SectionAddRequest struct {
	Name         string
	SlideIndices []int
}

func ParseSectionAddRequest(
	payload map[string]any,
	parseStringField ParseStringFieldFn,
	parseIntSliceField ParseIntSliceFieldFn,
) (SectionAddRequest, bool) {
	name, ok := parseStringField(payload, "name")
	if !ok {
		return SectionAddRequest{}, false
	}
	slideIndices, ok := parseIntSliceField(payload, "slide_indices")
	if !ok {
		return SectionAddRequest{}, false
	}
	return SectionAddRequest{
		Name:         name,
		SlideIndices: slideIndices,
	}, true
}

type SectionNameRequest struct {
	Name string
}

func ParseSectionNameRequest(
	payload map[string]any,
	parseStringField ParseStringFieldFn,
) (SectionNameRequest, bool) {
	name, ok := parseStringField(payload, "name")
	if !ok {
		return SectionNameRequest{}, false
	}
	return SectionNameRequest{Name: name}, true
}

type SectionRenameRequest struct {
	OldName string
	NewName string
}

func ParseSectionRenameRequest(
	payload map[string]any,
	parseStringField ParseStringFieldFn,
) (SectionRenameRequest, bool) {
	oldName, ok := parseStringField(payload, "old_name")
	if !ok {
		return SectionRenameRequest{}, false
	}
	newName, ok := parseStringField(payload, "new_name")
	if !ok {
		return SectionRenameRequest{}, false
	}
	return SectionRenameRequest{
		OldName: oldName,
		NewName: newName,
	}, true
}

type SlideSizeRequest struct {
	Width  int64
	Height int64
}

func ParseSlideSizeRequest(
	payload map[string]any,
	parseInt64Field ParseInt64FieldFn,
) (SlideSizeRequest, bool) {
	width, ok := parseInt64Field(payload, "width")
	if !ok {
		return SlideSizeRequest{}, false
	}
	height, ok := parseInt64Field(payload, "height")
	if !ok {
		return SlideSizeRequest{}, false
	}
	return SlideSizeRequest{
		Width:  width,
		Height: height,
	}, true
}

type SlideTitleRequest struct {
	SlideIndex int
	Title      string
}

func ParseSlideTitleRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
	parseStringField ParseStringFieldFn,
) (SlideTitleRequest, bool) {
	slideIndex, ok := parseSlideIndex(payload)
	if !ok {
		return SlideTitleRequest{}, false
	}
	title, ok := parseStringField(payload, "title")
	if !ok {
		return SlideTitleRequest{}, false
	}
	return SlideTitleRequest{
		SlideIndex: slideIndex,
		Title:      title,
	}, true
}

type MergeFromFileRequest struct {
	Path string
}

func ParseMergeFromFileRequest(
	payload map[string]any,
	parseStringField ParseStringFieldFn,
) (MergeFromFileRequest, bool) {
	path, ok := parseStringField(payload, "path")
	if !ok {
		return MergeFromFileRequest{}, false
	}
	return MergeFromFileRequest{Path: path}, true
}

func ParseCorePropertiesRequest(
	payload map[string]any,
	optionalString OptionalStringFieldFn,
) common.CoreProperties {
	return common.CoreProperties{
		Title:          optionalString(payload, "title"),
		Subject:        optionalString(payload, "subject"),
		Creator:        optionalString(payload, "creator"),
		Keywords:       optionalString(payload, "keywords"),
		Description:    optionalString(payload, "description"),
		LastModifiedBy: optionalString(payload, "lastModifiedBy"),
		Revision:       optionalString(payload, "revision"),
		Created:        optionalString(payload, "created"),
		Modified:       optionalString(payload, "modified"),
		Category:       optionalString(payload, "category"),
		ContentStatus:  optionalString(payload, "contentStatus"),
		Identifier:     optionalString(payload, "identifier"),
		Language:       optionalString(payload, "language"),
		LastPrinted:    optionalString(payload, "lastPrinted"),
		Version:        optionalString(payload, "version"),
	}
}

type UpdateSlideRequest struct {
	SlideIndex int
	Title      string
	Layout     string
	Bullets    []string
}

func ParseUpdateSlideRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
	optionalString OptionalStringFieldFn,
	optionalStringSlice OptionalStringSliceFieldFn,
) (UpdateSlideRequest, bool) {
	slideIndex, ok := parseSlideIndex(payload)
	if !ok {
		return UpdateSlideRequest{}, false
	}
	bullets, _ := optionalStringSlice(payload, "bullets")
	return UpdateSlideRequest{
		SlideIndex: slideIndex,
		Title:      optionalString(payload, "title"),
		Layout:     optionalString(payload, "layout"),
		Bullets:    bullets,
	}, true
}

type ChartSeriesRequest struct {
	Name   string
	Values []float64
}

type AddChartRequest struct {
	SlideIndex int
	ChartType  string
	Title      string
	Categories []string
	Values     []float64
	BarSeries  []ChartSeriesRequest
	LineSeries []ChartSeriesRequest
	X          int64
	Y          int64
	W          int64
	H          int64
}

func ParseAddChartRequest(
	payload map[string]any,
	parseSlideIndex ParseSlideIndexFn,
	parseStringField ParseStringFieldFn,
	optionalString OptionalStringFieldFn,
	parseStringSliceField ParseStringSliceFieldFn,
	parseFloat64SliceField ParseFloat64SliceFieldFn,
	optionalInt64 OptionalInt64FieldFn,
) (AddChartRequest, bool) {
	slideIndex, ok := parseSlideIndex(payload)
	if !ok {
		return AddChartRequest{}, false
	}
	chartType, ok := parseStringField(payload, "chart_type")
	if !ok {
		return AddChartRequest{}, false
	}
	categories, ok := parseStringSliceField(payload, "categories")
	if !ok {
		return AddChartRequest{}, false
	}
	values, ok := parseFloat64SliceField(payload, "values")
	if !ok {
		return AddChartRequest{}, false
	}
	x, _ := optionalInt64(payload, "x")
	y, _ := optionalInt64(payload, "y")
	w, _ := optionalInt64(payload, "w")
	h, _ := optionalInt64(payload, "h")
	barSeries := parseChartSeriesList(payload, "bar_series")
	lineSeries := parseChartSeriesList(payload, "line_series")

	return AddChartRequest{
		SlideIndex: slideIndex,
		ChartType:  chartType,
		Title:      optionalString(payload, "title"),
		Categories: categories,
		Values:     values,
		BarSeries:  barSeries,
		LineSeries: lineSeries,
		X:          x,
		Y:          y,
		W:          w,
		H:          h,
	}, true
}

func parseChartSeriesList(payload map[string]any, key string) []ChartSeriesRequest {
	raw, ok := payload[key]
	if !ok {
		return nil
	}
	items, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]ChartSeriesRequest, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		name, _ := m["name"].(string)
		var vals []float64
		if rawVals, ok := m["values"].([]any); ok {
			for _, v := range rawVals {
				switch n := v.(type) {
				case float64:
					vals = append(vals, n)
				case int:
					vals = append(vals, float64(n))
				}
			}
		}
		out = append(out, ChartSeriesRequest{Name: name, Values: vals})
	}
	return out
}
