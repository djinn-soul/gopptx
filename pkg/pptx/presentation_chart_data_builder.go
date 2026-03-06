package pptx

import (
	"errors"
	"fmt"
)

// ChartDataBuilder converts a chart-data builder model into an update payload
// for existing chart objects.
type ChartDataBuilder interface {
	chartDataUpdate() (ChartDataUpdate, error)
}

// CategoryChartData models categorical chart replacement/update input.
type CategoryChartData struct {
	Categories           []string
	MultiLevelCategories [][]string
	Series               []CategorySeries
}

// CategorySeries models one categorical series.
type CategorySeries struct {
	Name   string
	Values []float64
}

// NewCategoryChartData creates a categorical chart-data builder.
func NewCategoryChartData(categories []string) *CategoryChartData {
	return &CategoryChartData{Categories: append([]string(nil), categories...)}
}

// AddSeries appends one categorical series to this builder.
func (d *CategoryChartData) AddSeries(name string, values []float64) *CategoryChartData {
	if d == nil {
		return nil
	}
	d.Series = append(d.Series, CategorySeries{
		Name:   name,
		Values: append([]float64(nil), values...),
	})
	return d
}

// AddCategoryLevel appends one category level for multi-level category charts.
// Every level must have the same leaf count.
func (d *CategoryChartData) AddCategoryLevel(categories []string) *CategoryChartData {
	if d == nil {
		return nil
	}
	d.MultiLevelCategories = append(
		d.MultiLevelCategories,
		append([]string(nil), categories...),
	)
	return d
}

func (d *CategoryChartData) chartDataUpdate() (ChartDataUpdate, error) {
	if d == nil {
		return ChartDataUpdate{}, errors.New("category chart data builder is nil")
	}
	if len(d.Series) == 0 {
		return ChartDataUpdate{}, errors.New("category chart data requires at least one series")
	}
	if len(d.MultiLevelCategories) > 0 {
		leafCount := len(d.MultiLevelCategories[0])
		if leafCount == 0 {
			return ChartDataUpdate{}, errors.New("multi-level categories require at least one leaf value")
		}
		for i := 1; i < len(d.MultiLevelCategories); i++ {
			if len(d.MultiLevelCategories[i]) != leafCount {
				return ChartDataUpdate{}, fmt.Errorf("multi-level category level %d length mismatch", i)
			}
		}
	}
	series := make([]ChartSeriesData, 0, len(d.Series))
	for _, s := range d.Series {
		if len(d.MultiLevelCategories) > 0 && len(s.Values) != len(d.MultiLevelCategories[0]) {
			return ChartDataUpdate{}, errors.New("series values length must match multi-level category leaf count")
		}
		sCopy := s.Name
		series = append(series, ChartSeriesData{
			Name:   &sCopy,
			Values: append([]float64(nil), s.Values...),
		})
	}
	return ChartDataUpdate{
		Categories: append([]string(nil), d.Categories...),
		MultiLevelCategories: func() [][]string {
			if len(d.MultiLevelCategories) == 0 {
				return nil
			}
			out := make([][]string, 0, len(d.MultiLevelCategories))
			for _, lvl := range d.MultiLevelCategories {
				out = append(out, append([]string(nil), lvl...))
			}
			return out
		}(),
		Series: series,
	}, nil
}

// XyChartData models scatter/XY chart replacement/update input.
type XyChartData struct {
	Series []XySeries
}

// XySeries models one scatter/XY series.
type XySeries struct {
	Name    string
	XValues []float64
	YValues []float64
}

// NewXyChartData creates an XY chart-data builder.
func NewXyChartData() *XyChartData {
	return &XyChartData{}
}

// AddSeries appends one XY series to this builder.
func (d *XyChartData) AddSeries(name string, xValues, yValues []float64) *XyChartData {
	if d == nil {
		return nil
	}
	d.Series = append(d.Series, XySeries{
		Name:    name,
		XValues: append([]float64(nil), xValues...),
		YValues: append([]float64(nil), yValues...),
	})
	return d
}

func (d *XyChartData) chartDataUpdate() (ChartDataUpdate, error) {
	if d == nil {
		return ChartDataUpdate{}, errors.New("xy chart data builder is nil")
	}
	if len(d.Series) == 0 {
		return ChartDataUpdate{}, errors.New("xy chart data requires at least one series")
	}
	series := make([]ChartSeriesData, 0, len(d.Series))
	for i, s := range d.Series {
		if len(s.XValues) != len(s.YValues) {
			return ChartDataUpdate{}, fmt.Errorf("xy series %d has mismatched x/y lengths", i)
		}
		sCopy := s.Name
		series = append(series, ChartSeriesData{
			Name:    &sCopy,
			XValues: append([]float64(nil), s.XValues...),
			YValues: append([]float64(nil), s.YValues...),
		})
	}
	return ChartDataUpdate{Series: series}, nil
}

// BubbleChartData models bubble chart replacement/update input.
type BubbleChartData struct {
	Series []BubbleSeries
}

// BubbleSeries models one bubble chart series.
type BubbleSeries struct {
	Name    string
	XValues []float64
	YValues []float64
	Sizes   []float64
}

// NewBubbleChartData creates a bubble chart-data builder.
func NewBubbleChartData() *BubbleChartData {
	return &BubbleChartData{}
}

// AddSeries appends one bubble series to this builder.
func (d *BubbleChartData) AddSeries(
	name string,
	xValues []float64,
	yValues []float64,
	sizes []float64,
) *BubbleChartData {
	if d == nil {
		return nil
	}
	d.Series = append(d.Series, BubbleSeries{
		Name:    name,
		XValues: append([]float64(nil), xValues...),
		YValues: append([]float64(nil), yValues...),
		Sizes:   append([]float64(nil), sizes...),
	})
	return d
}

func (d *BubbleChartData) chartDataUpdate() (ChartDataUpdate, error) {
	if d == nil {
		return ChartDataUpdate{}, errors.New("bubble chart data builder is nil")
	}
	if len(d.Series) == 0 {
		return ChartDataUpdate{}, errors.New("bubble chart data requires at least one series")
	}
	series := make([]ChartSeriesData, 0, len(d.Series))
	for i, s := range d.Series {
		if len(s.XValues) != len(s.YValues) || len(s.XValues) != len(s.Sizes) {
			return ChartDataUpdate{}, fmt.Errorf("bubble series %d has mismatched x/y/size lengths", i)
		}
		sCopy := s.Name
		series = append(series, ChartSeriesData{
			Name:    &sCopy,
			XValues: append([]float64(nil), s.XValues...),
			YValues: append([]float64(nil), s.YValues...),
			Sizes:   append([]float64(nil), s.Sizes...),
		})
	}
	return ChartDataUpdate{Series: series}, nil
}

// UpdateChartDataByIndexFromBuilder updates chart data by chart index using a
// chart-data builder.
func (p *Presentation) UpdateChartDataByIndexFromBuilder(
	slideIndex int,
	chartIndex int,
	builder ChartDataBuilder,
) error {
	if builder == nil {
		return errors.New("chart data builder is nil")
	}
	update, err := builder.chartDataUpdate()
	if err != nil {
		return err
	}
	return p.UpdateChartDataByIndex(slideIndex, chartIndex, update)
}

// UpdateChartDataByRelIDFromBuilder updates chart data by relationship ID using
// a chart-data builder.
func (p *Presentation) UpdateChartDataByRelIDFromBuilder(
	slideIndex int,
	relID string,
	builder ChartDataBuilder,
) error {
	if builder == nil {
		return errors.New("chart data builder is nil")
	}
	update, err := builder.chartDataUpdate()
	if err != nil {
		return err
	}
	return p.UpdateChartDataByRelID(slideIndex, relID, update)
}
