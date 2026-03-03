package charts

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestCharts_Creation(t *testing.T) {
	cats := []string{"A", "B"}
	vals := []float64{10, 20}

	t.Run("Bar", func(t *testing.T) {
		c := NewBarChart(cats, vals)
		if len(c.Categories) != 2 || c.Values[0] != 10 {
			t.Error("NewBarChart failed")
		}
	})

	t.Run("Line", func(t *testing.T) {
		c := NewLineChart(cats, vals)
		if len(c.Categories) != 2 {
			t.Error("NewLineChart failed")
		}
	})

	t.Run("Pie", func(t *testing.T) {
		c := NewPieChart(cats, vals)
		if len(c.Categories) != 2 {
			t.Error("NewPieChart failed")
		}
	})

	t.Run("Area", func(t *testing.T) {
		c := NewAreaChart(cats, vals)
		if len(c.Categories) != 2 {
			t.Error("NewAreaChart failed")
		}
	})

	t.Run("Doughnut", func(t *testing.T) {
		c := NewDoughnutChart(cats, vals)
		if len(c.Categories) != 2 {
			t.Error("NewDoughnutChart failed")
		}
	})

	t.Run("Radar", func(t *testing.T) {
		c := NewRadarChart(cats, vals)
		if len(c.Categories) != 2 {
			t.Error("NewRadarChart failed")
		}
	})

	t.Run("Scatter", func(t *testing.T) {
		c := NewScatterChart([]float64{1}, []float64{2})
		if len(c.XValues) != 1 {
			t.Error("NewScatterChart failed")
		}
	})
}

func TestCharts_BarOptions(t *testing.T) {
	minVal, maxVal := 0.0, 100.0
	c := NewBarChart(nil, nil).
		WithTitle("Title").
		WithBarColor("FF0000").
		WithAltText("Alt").
		WithDecorative(true).
		Position(styling.Inches(1), styling.Inches(1)).
		Size(styling.Inches(5), styling.Inches(3))

	if c.Title != "Title" || c.BarColor != "FF0000" || c.AltText != "Alt" || !c.IsDecorative {
		t.Error("Bar options failed")
	}

	// Test variants
	b2 := BarHorizontalChart{BarChart: c}.WithSeriesName("S")
	if b2.SeriesName != "S" {
		t.Error("WithSeriesName failed")
	}

	b3 := BarStackedChart{BarChart: c}.WithLegend(true).WithLegendPosition("b")
	if !b3.ShowLegend || b3.LegendPosition != "b" {
		t.Error("WithLegend failed")
	}

	b4 := BarStacked100Chart{BarChart: c}.WithTitleOverlay(true).WithLegendOverlay(true)
	if !b4.TitleOverlay || !b4.LegendOverlay {
		t.Error("Overlay failed")
	}

	b5 := BarHorizontalChart{BarChart: c}.WithDataLabels(true).WithAxisTitles("X", "Y")
	if !b5.ShowDataLabels || b5.CategoryAxisTitle != "X" {
		t.Error("DataLabels/Axis failed")
	}

	b6 := BarStackedChart{BarChart: c}.WithMajorGridlines(true).WithValueFormat("0.00")
	if !b6.ShowMajorGridlines || b6.ValueFormat != "0.00" {
		t.Error("Grid/Format failed")
	}

	b7 := BarStacked100Chart{BarChart: c}.WithValueAxisCrossBetween("midCat").WithValueRange(minVal, maxVal)
	if b7.ValueAxisCrossBetween != "midCat" || *b7.MinValue != 0 {
		t.Error("Cross/Range failed")
	}
}

func TestCharts_LineOptions(t *testing.T) {
	c := NewLineChart(nil, nil).
		WithLineColor("00FF00")
	if c.LineColor != "00FF00" {
		t.Error("WithLineColor failed")
	}

	// Test Line variants
	l2 := LineMarkersChart{LineChart: c}.WithSeriesName("S").WithLegend(true).WithLegendPosition("l").
		WithTitleOverlay(true).WithLegendOverlay(true).WithDataLabels(true).WithAxisTitles("X", "Y").
		WithMajorGridlines(true).WithValueFormat("0.00").WithValueAxisCrossBetween("midCat").
		WithValueRange(0, 100)
	if l2.SeriesName != "S" || !l2.ShowLegend || l2.CategoryAxisTitle != "X" {
		t.Error("LineMarkers failed")
	}

	l3 := LineStackedChart{LineChart: c}.WithSeriesName("S").WithLegend(true).WithLegendPosition("r").
		WithTitleOverlay(true).WithLegendOverlay(true).WithDataLabels(true).WithAxisTitles("X", "Y").
		WithMajorGridlines(true).WithValueFormat("0.00").WithValueAxisCrossBetween("midCat").
		WithValueRange(0, 100)
	if l3.SeriesName != "S" || l3.ValueFormat != "0.00" {
		t.Error("LineStacked failed")
	}
}

func TestCharts_AreaOptions(t *testing.T) {
	c := NewAreaChart(nil, nil).WithAltText("A").WithDecorative(true).Position(0, 0).Size(1, 1).
		WithAreaColor("FF00FF").WithSeriesName("S").WithLegend(true).WithLegendPosition("t").
		WithTitleOverlay(true).WithLegendOverlay(true).WithDataLabels(true).WithAxisTitles("X", "Y").
		WithMajorGridlines(true).WithValueFormat("0.00").WithValueAxisCrossBetween("midCat").
		WithValueRange(0, 100)
	if c.AreaColor != "FF00FF" {
		t.Error("Area base failed")
	}

	a2 := AreaStackedChart{AreaChart: c}.WithSeriesName("S2").WithAreaColor("0000FF")
	if a2.AreaColor != "0000FF" {
		t.Error("WithAreaColor failed")
	}

	a3 := AreaStacked100Chart{AreaChart: c}.WithSeriesName("S3").WithAreaColor("00FF00")
	if a3.AreaColor != "00FF00" {
		t.Error("WithAreaColor 100 failed")
	}
}

func TestCharts_BubbleOptions(t *testing.T) {
	c := NewBubbleChart(nil, nil, nil).WithAltText("B").WithDecorative(true).Position(0, 0).Size(1, 1).
		WithLineColor("000000").WithLegend(true).WithLegendPosition("b").WithTitleOverlay(true).
		WithLegendOverlay(true).WithDataLabels(true).WithAxisTitles("X", "Y").WithMajorGridlines(true).
		WithValueFormat("0.0").WithValueAxisCrossBetween("between").WithValueRange(0, 10)
	if c.AltText != "B" {
		t.Error("Bubble failed")
	}
}

func TestCharts_PieDoughnutOptions(t *testing.T) {
	p := NewPieChart(nil, nil).WithAltText("P").WithDecorative(true).Position(0, 0).Size(1, 1).
		WithTitle("PT").WithLegend(true).WithLegendPosition("l").WithTitleOverlay(true).
		WithLegendOverlay(true).WithDataLabels(true)
	if p.AltText != "P" {
		t.Error("Pie failed")
	}

	d := NewDoughnutChart(nil, nil).WithAltText("D").WithDecorative(true).Position(0, 0).Size(1, 1).
		WithTitle("DT").WithLegend(true).WithLegendPosition("r").WithTitleOverlay(true).
		WithLegendOverlay(true).WithDataLabels(true)
	if d.AltText != "D" {
		t.Error("Doughnut failed")
	}
}

func TestCharts_RadarOptions(t *testing.T) {
	r := NewRadarChart(nil, nil).WithAltText("R").WithDecorative(true).Position(0, 0).Size(1, 1).
		WithLineColor("111111").WithSeriesName("S").WithLegend(true).WithLegendPosition("b").
		WithTitleOverlay(true).WithLegendOverlay(true).WithDataLabels(true).WithAxisTitles("X", "Y").
		WithMajorGridlines(true).WithValueFormat("0").WithValueAxisCrossBetween("midCat").
		WithValueRange(0, 5)
	if r.AltText != "R" {
		t.Error("Radar failed")
	}

	rf := NewRadarFilledChart(nil, nil).WithAltText("RF").WithDecorative(true).
		WithLineColor("222222").WithSeriesName("SF").WithLegend(true).WithLegendPosition("t").
		WithTitleOverlay(true).WithLegendOverlay(true).WithDataLabels(true).WithAxisTitles("X", "Y").
		WithMajorGridlines(true).WithValueFormat("0").WithValueAxisCrossBetween("midCat").
		WithValueRange(0, 5).Position(0, 0).Size(1, 1)
	if rf.AltText != "RF" {
		t.Error("RadarFilled failed")
	}
}

func TestCharts_ScatterOptions(t *testing.T) {
	s := NewScatterChart(nil, nil).WithAltText("S").WithDecorative(true).Position(0, 0).Size(1, 1).
		WithTitle("ST").WithLegend(true).WithLegendPosition("l").WithTitleOverlay(true).
		WithLegendOverlay(true).WithDataLabels(true).WithAxisTitles("X", "Y").WithMajorGridlines(true).
		WithValueFormat("0").WithValueAxisCrossBetween("midCat").WithValueRange(0, 100)
	if s.AltText != "S" {
		t.Error("Scatter failed")
	}
}

func TestCharts_StockComboOptions(t *testing.T) {
	s1 := NewStockHLCChart(nil, nil, nil, nil).WithAltText("S1").WithDecorative(true).
		WithTitle("S1T").WithLegend(true).WithLegendPosition("b").WithTitleOverlay(true).
		WithLegendOverlay(true).WithDataLabels(true).WithAxisTitles("X", "Y").WithMajorGridlines(true).
		WithValueFormat("0").WithValueAxisCrossBetween("midCat").WithValueRange(0, 100).
		Position(0, 0).Size(1, 1)
	if s1.AltText != "S1" {
		t.Error("HLC failed")
	}

	s2 := NewStockOHLCChart(nil, nil, nil, nil, nil).WithAltText("S2").WithDecorative(true).
		WithTitle("S2T").WithLegend(true).WithLegendPosition("t").WithTitleOverlay(true).
		WithLegendOverlay(true).WithDataLabels(true).WithAxisTitles("X", "Y").WithMajorGridlines(true).
		WithValueFormat("0").WithValueAxisCrossBetween("midCat").WithValueRange(0, 100).
		Position(0, 0).Size(1, 1)
	if s2.AltText != "S2" {
		t.Error("OHLC failed")
	}

	c := NewComboChart(nil, nil, nil).WithAltText("C").WithDecorative(true).
		WithTitle("CT").WithLegend(true).WithLegendPosition("r").WithTitleOverlay(true).
		WithLegendOverlay(true).WithDataLabels(true).WithAxisTitles("X", "Y").WithMajorGridlines(true).
		WithValueFormat("0").WithValueAxisCrossBetween("midCat").WithValueRange(0, 100).
		Position(0, 0).Size(1, 1)
	if c.AltText != "C" {
		t.Error("Combo failed")
	}
}

func TestCharts_Getters(t *testing.T) {
	cats := []string{"A"}
	vals := []float64{1}

	if len(NewBarChart(cats, vals).GetCategories()) != 1 {
		t.Error("Bar Getter failed")
	}
	if len(NewLineChart(cats, vals).GetValues()) != 1 {
		t.Error("Line Getter failed")
	}
	if len(NewPieChart(cats, vals).GetCategories()) != 1 {
		t.Error("Pie Getter failed")
	}
	if len(NewDoughnutChart(cats, vals).GetValues()) != 1 {
		t.Error("Doughnut Getter failed")
	}
	if len(NewAreaChart(cats, vals).GetCategories()) != 1 {
		t.Error("Area Getter failed")
	}
	if len(NewRadarChart(cats, vals).GetValues()) != 1 {
		t.Error("Radar Getter failed")
	}

	sc := NewScatterChart([]float64{1}, []float64{2})
	if len(sc.GetCategories()) != 1 || len(sc.GetValues()) != 1 {
		t.Error("Scatter Getter failed")
	}
}

func TestCharts_Validate_Extended(t *testing.T) {
	tests := []struct {
		name    string
		chart   BarChart
		wantErr bool
	}{
		{"Valid", NewBarChart([]string{"A"}, []float64{1}), false},
		{"Mismatch", NewBarChart([]string{"A"}, []float64{1, 2}), true},
		{"Empty", NewBarChart(nil, nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.chart.Validate(1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
