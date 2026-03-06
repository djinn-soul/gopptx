package pptx

import (
	"errors"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// ChartSelector identifies a slide chart by index and/or relationship ID.
type ChartSelector = common.ChartSelector

// ChartSeriesData carries one chart series worth of input data.
type ChartSeriesData = common.ChartSeriesData

// ChartDataUpdate is the complete chart update payload.
type ChartDataUpdate = common.ChartDataUpdate

// ChartFormatUpdate is a partial formatting patch for an existing chart.
type ChartFormatUpdate = common.ChartFormatUpdate

// SlideChartRef describes a chart relationship discovered on a slide.
type SlideChartRef = common.SlideChartRef

// ListSlideCharts returns chart references discovered on the given slide.
func (p *Presentation) ListSlideCharts(slideIndex int) ([]SlideChartRef, error) {
	if p == nil {
		return nil, errors.New("presentation is nil")
	}
	if p.editor == nil {
		return nil, errors.New("presentation editor is not initialized")
	}
	return p.editor.ListSlideCharts(slideIndex)
}

// UpdateChartData updates chart workbook/cache data for a selected chart on a
// slide.
func (p *Presentation) UpdateChartData(slideIndex int, selector ChartSelector, data ChartDataUpdate) error {
	if p == nil {
		return errors.New("presentation is nil")
	}
	if p.editor == nil {
		return errors.New("presentation editor is not initialized")
	}
	return p.editor.UpdateChartData(slideIndex, selector, data)
}

// UpdateChartDataByIndex updates chart data for a chart selected by 0-based
// chart index on a slide.
func (p *Presentation) UpdateChartDataByIndex(slideIndex int, chartIndex int, data ChartDataUpdate) error {
	idx := chartIndex
	return p.UpdateChartData(slideIndex, ChartSelector{Index: &idx}, data)
}

// UpdateChartDataByRelID updates chart data for a chart selected by
// relationship ID on a slide.
func (p *Presentation) UpdateChartDataByRelID(slideIndex int, relID string, data ChartDataUpdate) error {
	return p.UpdateChartData(slideIndex, ChartSelector{RelID: relID}, data)
}

// UpdateChartFormatting applies a partial formatting patch to a selected chart.
func (p *Presentation) UpdateChartFormatting(
	slideIndex int,
	selector ChartSelector,
	format ChartFormatUpdate,
) error {
	if p == nil {
		return errors.New("presentation is nil")
	}
	if p.editor == nil {
		return errors.New("presentation editor is not initialized")
	}
	return p.editor.UpdateChartFormatting(slideIndex, selector, format)
}

// UpdateChartFormattingByIndex applies chart formatting by chart index.
func (p *Presentation) UpdateChartFormattingByIndex(
	slideIndex int,
	chartIndex int,
	format ChartFormatUpdate,
) error {
	idx := chartIndex
	return p.UpdateChartFormatting(slideIndex, ChartSelector{Index: &idx}, format)
}

// UpdateChartFormattingByRelID applies chart formatting by relationship ID.
func (p *Presentation) UpdateChartFormattingByRelID(
	slideIndex int,
	relID string,
	format ChartFormatUpdate,
) error {
	return p.UpdateChartFormatting(slideIndex, ChartSelector{RelID: relID}, format)
}

// ReplaceChartData updates a chart identified by its slide-local index with a
// single category/value series.
func (p *Presentation) ReplaceChartData(
	slideIndex int,
	chartIndex int,
	categories []string,
	values []float64,
) error {
	if p == nil {
		return errors.New("presentation is nil")
	}
	if p.editor == nil {
		return errors.New("presentation editor is not initialized")
	}
	return p.editor.ReplaceChartData(slideIndex, chartIndex, categories, values)
}

// ReplaceChartDataByRelID updates a category/value chart selected by
// relationship ID.
func (p *Presentation) ReplaceChartDataByRelID(
	slideIndex int,
	relID string,
	categories []string,
	values []float64,
) error {
	return p.UpdateChartDataByRelID(slideIndex, relID, ChartDataUpdate{
		Categories: categories,
		Series: []ChartSeriesData{
			{Values: values},
		},
	})
}
