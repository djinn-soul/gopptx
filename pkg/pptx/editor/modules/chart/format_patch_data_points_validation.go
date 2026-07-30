package chart

import (
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

func validateChartDataPoints(points []common.ChartDataPoint, clearSeries []int) error {
	for _, point := range points {
		if err := validateChartDataPoint(point); err != nil {
			return err
		}
	}
	for _, index := range clearSeries {
		if index < 0 {
			return fmt.Errorf("clear_data_point_series index must not be negative, got %d", index)
		}
	}
	return nil
}

func validateChartDataPoint(point common.ChartDataPoint) error {
	if point.SeriesIndex < 0 {
		return fmt.Errorf("data point series_index must not be negative, got %d", point.SeriesIndex)
	}
	if point.PointIndex < 0 {
		return fmt.Errorf("data point point_index must not be negative, got %d", point.PointIndex)
	}
	for _, pair := range []struct {
		name  string
		value *string
	}{{"fill_color", point.FillColor}, {"line_color", point.LineColor}} {
		if pair.value == nil {
			continue
		}
		if _, err := editorshape.NormalizeHexColor(*pair.value); err != nil {
			return fmt.Errorf("data point %s: %w", pair.name, err)
		}
	}
	if point.LineWidthEMU != nil && *point.LineWidthEMU < 0 {
		return fmt.Errorf("data point line_width_emu must not be negative, got %d", *point.LineWidthEMU)
	}
	if point.Explosion != nil && (*point.Explosion < 0 || *point.Explosion > explosionMax) {
		return fmt.Errorf("data point explosion must be between 0 and %d, got %d", explosionMax, *point.Explosion)
	}
	return validateChartDataPointMarker(point)
}

func validateChartDataPointMarker(point common.ChartDataPoint) error {
	for _, pair := range []struct {
		name  string
		value *string
	}{{"marker_fill_color", point.MarkerFillColor}, {"marker_line_color", point.MarkerLineColor}} {
		if pair.value == nil {
			continue
		}
		if _, err := editorshape.NormalizeHexColor(*pair.value); err != nil {
			return fmt.Errorf("data point %s: %w", pair.name, err)
		}
	}
	if point.MarkerSymbol != nil && !markerSymbols[*point.MarkerSymbol] {
		return fmt.Errorf("data point marker_symbol %q is not a CT_MarkerStyle value", *point.MarkerSymbol)
	}
	if point.MarkerSize != nil && (*point.MarkerSize < markerSizeMin || *point.MarkerSize > markerSizeMax) {
		return fmt.Errorf(
			"data point marker_size must be between %d and %d, got %d",
			markerSizeMin, markerSizeMax, *point.MarkerSize,
		)
	}
	return nil
}
