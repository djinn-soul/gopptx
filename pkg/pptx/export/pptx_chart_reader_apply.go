package export

import (
	"strconv"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	parsedChartKindBar = "bar"
	parsedChartKindPie = "pie"
	seriesIndex0       = 0
	seriesIndex1       = 1
	seriesIndex2       = 2
	seriesIndex3       = 3
)

func applyParsedCharts(slide *elements.SlideContent, chartList []parsedChart) {
	for _, pc := range chartList {
		applyParsedChart(slide, pc)
	}
}

// setIfNonEmpty sets *dest to value only when value is non-empty.
func setIfNonEmpty(dest *string, value string) {
	if value != "" {
		*dest = value
	}
}

// chartApplyCtx bundles the common fields derived from a parsedChart for chart construction.
type chartApplyCtx struct {
	slide          *elements.SlideContent
	pc             parsedChart
	px, py, pw, ph styling.Length
	title          string
	seriesColor    string
	cats           []string
	vals           []float64
	axisMin        *float64 // nil → auto
	axisMax        *float64 // nil → auto
}

func newChartApplyCtx(slide *elements.SlideContent, pc parsedChart) chartApplyCtx {
	ctx := chartApplyCtx{
		slide: slide,
		pc:    pc,
		px:    styling.Emu(pc.X),
		py:    styling.Emu(pc.Y),
		pw:    styling.Emu(pc.CX),
		ph:    styling.Emu(pc.CY),
		title: nonEmpty(pc.Title, "Chart"),
		cats:  categoriesForChart(pc),
		vals:  valuesForChart(pc),
	}
	if len(pc.Series) > 0 {
		ctx.seriesColor = pc.Series[0].Color
	}
	ctx.axisMin = pc.AxisMinValue
	ctx.axisMax = pc.AxisMaxValue
	return ctx
}

func applyParsedChart(slide *elements.SlideContent, pc parsedChart) {
	ctx := newChartApplyCtx(slide, pc)
	switch {
	case pc.Kind == parsedChartKindBar || pc.Kind == "barHorizontal" ||
		pc.Kind == "barStacked" || pc.Kind == "barStacked100":
		applyBarLikeChart(ctx)
	case pc.Kind == "line" || pc.Kind == "lineMarkers" || pc.Kind == "lineStacked":
		applyLineLikeChart(ctx)
	case pc.Kind == "area" || pc.Kind == "areaStacked" || pc.Kind == "areaStacked100":
		applyAreaLikeChart(ctx)
	default:
		applyOtherChart(ctx)
	}
}

func applyBarLikeChart(ctx chartApplyCtx) {
	slide, pc, cats, vals, title, px, py, pw, ph, color :=
		ctx.slide, ctx.pc, ctx.cats, ctx.vals, ctx.title, ctx.px, ctx.py, ctx.pw, ctx.ph, ctx.seriesColor
	switch pc.Kind {
	case parsedChartKindBar:
		c := charts.NewBarChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		setIfNonEmpty(&c.BarColor, color)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.Chart = &c
	case "barHorizontal":
		c := charts.NewBarHorizontalChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		setIfNonEmpty(&c.BarColor, color)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.BarHorizontal = &c
	case "barStacked":
		c := charts.NewBarStackedChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		setIfNonEmpty(&c.BarColor, color)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.BarStacked = &c
	case "barStacked100":
		c := charts.NewBarStacked100Chart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		setIfNonEmpty(&c.BarColor, color)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.BarStacked100 = &c
	}
}

func applyLineLikeChart(ctx chartApplyCtx) {
	slide, pc, cats, vals, title, px, py, pw, ph, color :=
		ctx.slide, ctx.pc, ctx.cats, ctx.vals, ctx.title, ctx.px, ctx.py, ctx.pw, ctx.ph, ctx.seriesColor
	switch pc.Kind {
	case "line":
		c := charts.NewLineChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		setIfNonEmpty(&c.LineColor, color)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.Line = &c
	case "lineMarkers":
		c := charts.NewLineMarkersChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		setIfNonEmpty(&c.LineColor, color)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.LineMarkers = &c
	case "lineStacked":
		c := charts.NewLineStackedChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		setIfNonEmpty(&c.LineColor, color)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.LineStacked = &c
	}
}

func applyAreaLikeChart(ctx chartApplyCtx) {
	slide, pc, cats, vals, title, px, py, pw, ph, color :=
		ctx.slide, ctx.pc, ctx.cats, ctx.vals, ctx.title, ctx.px, ctx.py, ctx.pw, ctx.ph, ctx.seriesColor
	switch pc.Kind {
	case "area":
		c := charts.NewAreaChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		setIfNonEmpty(&c.AreaColor, color)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.Area = &c
	case "areaStacked":
		c := charts.NewAreaStackedChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		setIfNonEmpty(&c.AreaColor, color)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.AreaStacked = &c
	case "areaStacked100":
		c := charts.NewAreaStacked100Chart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		setIfNonEmpty(&c.AreaColor, color)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.AreaStacked100 = &c
	}
}

func applyOtherChart(ctx chartApplyCtx) {
	slide, pc, cats, vals, title, px, py, pw, ph :=
		ctx.slide, ctx.pc, ctx.cats, ctx.vals, ctx.title, ctx.px, ctx.py, ctx.pw, ctx.ph
	switch pc.Kind {
	case "scatter":
		xs, ys := scatterForChart(pc)
		c := charts.NewScatterChart(xs, ys).WithTitle(title).Position(px, py).Size(pw, ph)
		setIfNonEmpty(&c.LineColor, ctx.seriesColor)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.ScatterStyle = pc.ScatterStyle
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.Scatter = &c
	case parsedChartKindPie:
		c := charts.NewPieChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		slide.Pie = &c
	case "doughnut":
		c := charts.NewDoughnutChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		slide.Doughnut = &c
	case "bubble":
		xs, ys, sizes := bubbleForChart(pc)
		c := charts.NewBubbleChart(xs, ys, sizes).
			WithTitle(title).Position(px.Emu(), py.Emu()).Size(pw.Emu(), ph.Emu())
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.Bubble = &c
	case "radar":
		c := charts.NewRadarChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		slide.Radar = &c
	case "radarFilled":
		c := charts.NewRadarFilledChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		slide.RadarFilled = &c
	case "stockHLC":
		high, low, closeVals := stockTriplet(pc)
		c := charts.NewStockHLCChart(cats, high, low, closeVals).WithTitle(title).Position(px, py).Size(pw, ph)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.StockHLC = &c
	case "stockOHLC":
		openVals, high, low, closeVals := stockQuad(pc)
		c := charts.NewStockOHLCChart(cats, openVals, high, low, closeVals).
			WithTitle(title).Position(px, py).Size(pw, ph)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.StockOHLC = &c
	case "combo":
		barSeries, lineSeries := comboSeries(pc)
		c := charts.NewComboChart(cats, barSeries, lineSeries).WithTitle(title).Position(px, py).Size(pw, ph)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		c.MinValue, c.MaxValue = ctx.axisMin, ctx.axisMax
		slide.Combo = &c
	}
}

func categoriesForChart(c parsedChart) []string {
	if len(c.Series) > 0 && len(c.Series[0].Categories) > 0 {
		return c.Series[0].Categories
	}
	vals := valuesForChart(c)
	out := make([]string, len(vals))
	for i := range vals {
		out[i] = "Item " + strconv.Itoa(i+1)
	}
	return out
}

func valuesForChart(c parsedChart) []float64 {
	if len(c.Series) > 0 && len(c.Series[0].Values) > 0 {
		return c.Series[0].Values
	}
	return []float64{1, 2, 3}
}

func scatterForChart(c parsedChart) ([]float64, []float64) {
	if len(c.Series) > 0 && len(c.Series[0].XValues) > 0 && len(c.Series[0].YValues) > 0 {
		return c.Series[0].XValues, c.Series[0].YValues
	}
	return []float64{1, 2, 3}, []float64{1, 2, 3}
}

func bubbleForChart(c parsedChart) ([]float64, []float64, []float64) {
	if len(c.Series) > 0 {
		s := c.Series[0]
		if len(s.XValues) > 0 && len(s.YValues) > 0 && len(s.Sizes) > 0 {
			return s.XValues, s.YValues, s.Sizes
		}
	}
	return []float64{1, 2, 3}, []float64{1, 2, 3}, []float64{12, 18, 24}
}

func stockTriplet(c parsedChart) ([]float64, []float64, []float64) {
	high := pickSeriesValues(c, seriesIndex0)
	low := pickSeriesValues(c, seriesIndex1)
	closeVals := pickSeriesValues(c, seriesIndex2)
	return high, low, closeVals
}

func stockQuad(c parsedChart) ([]float64, []float64, []float64, []float64) {
	openVals := pickSeriesValues(c, seriesIndex0)
	high := pickSeriesValues(c, seriesIndex1)
	low := pickSeriesValues(c, seriesIndex2)
	closeVals := pickSeriesValues(c, seriesIndex3)
	return openVals, high, low, closeVals
}

func comboSeries(c parsedChart) ([]charts.Series, []charts.Series) {
	// First series → bar, second → line; use names from the parsed data when available.
	barName := "Bar"
	lineNameStr := "Line"
	if len(c.Series) > 0 && c.Series[0].Name != "" {
		barName = c.Series[0].Name
	}
	if len(c.Series) > 1 && c.Series[1].Name != "" {
		lineNameStr = c.Series[1].Name
	}
	barVals := pickSeriesValues(c, seriesIndex0)
	lineVals := pickSeriesValues(c, seriesIndex1)
	return []charts.Series{{Name: barName, Values: barVals}}, []charts.Series{{Name: lineNameStr, Values: lineVals}}
}

func pickSeriesValues(c parsedChart, idx int) []float64 {
	if idx >= 0 && idx < len(c.Series) && len(c.Series[idx].Values) > 0 {
		return c.Series[idx].Values
	}
	return valuesForChart(c)
}

func nonEmpty(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
