package main

import (
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

type chartTypeInfo struct {
	name    string
	builder func() pptx.SlideContent
}

func titleOnlySlide(title string) pptx.SlideContent {
	s := pptx.NewSlide(title)
	s.Layout = pptx.SlideLayoutTitleOnly
	return s
}

func allChartTypes(data chartDemoData) []chartTypeInfo {
	types := append(coreChartTypes(data), lineAndScatterChartTypes(data)...)
	types = append(types, stockAndComboChartTypes(data)...)
	return types
}

func coreChartTypes(data chartDemoData) []chartTypeInfo {
	return []chartTypeInfo{
		{name: "Area", builder: areaChartSlide(data)},
		{name: "Area Stacked", builder: areaStackedChartSlide(data)},
		{name: "Area Stacked 100%", builder: areaStacked100ChartSlide(data)},
		{name: "Column / Bar", builder: barChartSlide(data)},
		{name: "Bar Horizontal", builder: barHorizontalChartSlide(data)},
		{name: "Bar Stacked", builder: barStackedChartSlide(data)},
		{name: "Bar Stacked 100%", builder: barStacked100ChartSlide(data)},
	}
}

func lineAndScatterChartTypes(data chartDemoData) []chartTypeInfo {
	return []chartTypeInfo{
		{name: "Bubble", builder: bubbleChartSlide(data)},
		{name: "Doughnut", builder: doughnutChartSlide(data)},
		{name: "Line", builder: lineChartSlide(data)},
		{name: "Line with Markers", builder: lineMarkersChartSlide(data)},
		{name: "Line Stacked", builder: lineStackedChartSlide(data)},
		{name: "Pie", builder: pieChartSlide(data)},
		{name: "Radar", builder: radarChartSlide(data, false)},
		{name: "Radar Filled", builder: radarChartSlide(data, true)},
		{name: "Scatter", builder: scatterChartSlide(data)},
	}
}

func stockAndComboChartTypes(data chartDemoData) []chartTypeInfo {
	return []chartTypeInfo{
		{name: "Stock HLC", builder: stockHLCChartSlide(data)},
		{name: "Stock OHLC", builder: stockOHLCChartSlide(data)},
	}
}

func areaChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Area Chart").WithAreaChart(
			pptx.NewAreaChart(data.cats, data.vals).WithTitle("Area Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func areaStackedChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Area Stacked Chart").WithAreaStackedChart(
			pptx.NewAreaStackedChart(data.cats, data.vals).WithTitle("Area Stacked Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func areaStacked100ChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Area Stacked 100% Chart").WithAreaStacked100Chart(
			pptx.NewAreaStacked100Chart(data.cats, data.vals).WithTitle("Area Stacked 100% Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func barChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Column / Bar Chart").WithBarChart(
			pptx.NewBarChart(data.cats, data.vals).WithTitle("Column / Bar Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func barHorizontalChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Bar Horizontal Chart").WithBarHorizontalChart(
			pptx.NewBarHorizontalChart(data.cats, data.vals).WithTitle("Bar Horizontal Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func barStackedChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Bar Stacked Chart").WithBarStackedChart(
			pptx.NewBarStackedChart(data.cats, data.vals).WithTitle("Bar Stacked Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func barStacked100ChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Bar Stacked 100% Chart").WithBarStacked100Chart(
			pptx.NewBarStacked100Chart(data.cats, data.vals).WithTitle("Bar Stacked 100% Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func bubbleChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Bubble Chart").WithBubbleChart(
			pptx.NewBubbleChart(data.xVals, data.yVals, data.sizes).WithTitle("Bubble Demo").
				Position(data.chartX.Emu(), data.chartY.Emu()).Size(data.chartW.Emu(), data.chartH.Emu()),
		)
	}
}

func doughnutChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Doughnut Chart").WithDoughnutChart(
			pptx.NewDoughnutChart(data.cats, data.vals).WithTitle("Doughnut Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func lineChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Line Chart").WithLineChart(
			pptx.NewLineChart(data.cats, data.vals).WithTitle("Line Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func lineMarkersChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Line with Markers Chart").WithLineMarkersChart(
			pptx.NewLineMarkersChart(data.cats, data.vals).WithTitle("Line with Markers Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func lineStackedChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Line Stacked Chart").WithLineStackedChart(
			pptx.NewLineStackedChart(data.cats, data.vals).WithTitle("Line Stacked Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func pieChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Pie Chart").WithPieChart(
			pptx.NewPieChart(data.cats, data.vals).WithTitle("Pie Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func radarChartSlide(data chartDemoData, filled bool) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		if filled {
			return titleOnlySlide("Radar Filled Chart").WithRadarFilledChart(
				pptx.NewRadarFilledChart(data.cats, data.vals).WithTitle("Radar Filled Demo").
					Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
			)
		}
		return titleOnlySlide("Radar Chart").WithRadarChart(
			pptx.NewRadarChart(data.cats, data.vals).WithTitle("Radar Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func scatterChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Scatter Chart").WithScatterChart(
			pptx.NewScatterChart(data.xVals, data.yVals).WithTitle("Scatter Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func stockHLCChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Stock HLC Chart").WithStockHLCChart(
			pptx.NewStockHLCChart(data.hlcCats, data.high, data.low, data.closeV).WithTitle("Stock HLC Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func stockOHLCChartSlide(data chartDemoData) func() pptx.SlideContent {
	return func() pptx.SlideContent {
		return titleOnlySlide("Stock OHLC Chart").WithStockOHLCChart(
			pptx.NewStockOHLCChart(
				data.hlcCats, data.openV, data.high, data.low, data.closeV,
			).WithTitle("Stock OHLC Demo").
				Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
		)
	}
}

func buildSlides(data chartDemoData) []pptx.SlideContent {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Chart Types Demo").
			AddBullet("This deck demonstrates chart creation from Go API.").
			AddBullet("It outputs both PPTX and PDF artifacts."),
	}
	for _, ct := range allChartTypes(data) {
		slides = append(slides, ct.builder())
	}
	slides = append(slides, comboSlide(data))
	slides = append(slides, referenceSlides()...)
	return slides
}

func comboSlide(data chartDemoData) pptx.SlideContent {
	return titleOnlySlide("Combo Chart").WithComboChart(
		pptx.NewComboChart(
			data.cats,
			[]pptx.Series{{Name: "Revenue", Values: []float64{180, 220, 210, 260}}},
			[]pptx.Series{{Name: "Growth %", Values: []float64{8, 11, 10, 14}}},
		).WithTitle("Combo Demo").Position(data.chartX, data.chartY).Size(data.chartW, data.chartH),
	)
}

func referenceSlides() []pptx.SlideContent {
	goChartConstructors := []struct{ constructor, kind string }{
		{"pptx.NewBarChart", "bar"},
		{"pptx.NewBarChart (BAR alias)", "bar"},
		{"pptx.NewBarHorizontalChart", "barHorizontal"},
		{"pptx.NewBarStackedChart", "barStacked"},
		{"pptx.NewBarStacked100Chart", "barStacked100"},
		{"pptx.NewLineChart", "line"},
		{"pptx.NewLineMarkersChart", "lineMarkers"},
		{"pptx.NewLineStackedChart", "lineStacked"},
		{"pptx.NewScatterChart", "scatter"},
		{"pptx.NewAreaChart", "area"},
		{"pptx.NewAreaStackedChart", "areaStacked"},
		{"pptx.NewAreaStacked100Chart", "areaStacked100"},
		{"pptx.NewPieChart", "pie"},
		{"pptx.NewDoughnutChart", "doughnut"},
		{"pptx.NewBubbleChart", "bubble"},
		{"pptx.NewRadarChart", "radar"},
		{"pptx.NewRadarFilledChart", "radarFilled"},
		{"pptx.NewStockHLCChart", "stockHLC"},
		{"pptx.NewStockOHLCChart", "stockOHLC"},
		{"pptx.NewComboChart", "combo"},
	}
	lines := make([]string, 0, len(goChartConstructors))
	for _, c := range goChartConstructors {
		lines = append(lines, fmt.Sprintf("%s = %q", c.constructor, c.kind))
	}
	half := len(lines) / 2

	uniqueKinds := make(map[string]struct{}, len(goChartConstructors))
	for _, c := range goChartConstructors {
		uniqueKinds[c.kind] = struct{}{}
	}
	s1 := pptx.NewSlide("Go Chart Surface (1/2)").
		AddBullet(fmt.Sprintf("Named constructors: %d", len(goChartConstructors))).
		AddBullet(fmt.Sprintf("Unique chart kinds: %d", len(uniqueKinds)))
	for _, l := range lines[:half] {
		s1 = s1.AddBullet(l)
	}

	s2 := pptx.NewSlide("Go Chart Surface (2/2)")
	for _, l := range lines[half:] {
		s2 = s2.AddBullet(l)
	}

	return []pptx.SlideContent{s1, s2}
}
