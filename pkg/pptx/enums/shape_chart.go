package enums

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

type MSOShape string

func (s MSOShape) XMLValue() string {
	return string(s)
}

func ParseMSOShape(value string) (MSOShape, error) {
	key := normalizeKey(value)
	switch key {
	case "roundedrectangle":
		value = "roundRect"
	case "rectangle":
		value = "rect"
	case "oval", "circle":
		value = "ellipse"
	}
	normalized := shapes.NormalizeShapeType(value)
	if !shapes.IsShapeType(normalized) {
		return "", fmt.Errorf("invalid MSO_SHAPE value %q", value)
	}
	return MSOShape(normalized), nil
}

type XLChartType string

const (
	XLChartTypeBar            XLChartType = pptxxml.ChartKindBar
	XLChartTypeBarHorizontal  XLChartType = pptxxml.ChartKindBarHorizontal
	XLChartTypeBarStacked     XLChartType = pptxxml.ChartKindBarStacked
	XLChartTypeBarStacked100  XLChartType = pptxxml.ChartKindBarStacked100
	XLChartTypeLine           XLChartType = pptxxml.ChartKindLine
	XLChartTypeLineMarkers    XLChartType = pptxxml.ChartKindLineMarkers
	XLChartTypeLineStacked    XLChartType = pptxxml.ChartKindLineStacked
	XLChartTypeScatter        XLChartType = pptxxml.ChartKindScatter
	XLChartTypeArea           XLChartType = pptxxml.ChartKindArea
	XLChartTypeAreaStacked    XLChartType = pptxxml.ChartKindAreaStacked
	XLChartTypeAreaStacked100 XLChartType = pptxxml.ChartKindAreaStacked100
	XLChartTypePie            XLChartType = pptxxml.ChartKindPie
	XLChartTypeDoughnut       XLChartType = pptxxml.ChartKindDoughnut
	XLChartTypeBubble         XLChartType = pptxxml.ChartKindBubble
	XLChartTypeRadar          XLChartType = pptxxml.ChartKindRadar
	XLChartTypeRadarFilled    XLChartType = pptxxml.ChartKindRadarFilled
	XLChartTypeStockHLC       XLChartType = pptxxml.ChartKindStockHLC
	XLChartTypeStockOHLC      XLChartType = pptxxml.ChartKindStockOHLC
	XLChartTypeCombo          XLChartType = pptxxml.ChartKindCombo
)

func (t XLChartType) XMLValue() string {
	return string(t)
}

func ParseXLChartType(value string) (XLChartType, error) {
	switch normalizeKey(value) {
	case "bar":
		return XLChartTypeBar, nil
	case "barhorizontal", "barh":
		return XLChartTypeBarHorizontal, nil
	case "barstacked":
		return XLChartTypeBarStacked, nil
	case "barstacked100":
		return XLChartTypeBarStacked100, nil
	case "line":
		return XLChartTypeLine, nil
	case "linemarkers":
		return XLChartTypeLineMarkers, nil
	case "linestacked":
		return XLChartTypeLineStacked, nil
	case "scatter", "xy":
		return XLChartTypeScatter, nil
	case "area":
		return XLChartTypeArea, nil
	case "areastacked":
		return XLChartTypeAreaStacked, nil
	case "areastacked100":
		return XLChartTypeAreaStacked100, nil
	case "pie":
		return XLChartTypePie, nil
	case "doughnut", "donut":
		return XLChartTypeDoughnut, nil
	case "bubble":
		return XLChartTypeBubble, nil
	case "radar":
		return XLChartTypeRadar, nil
	case "radarfilled":
		return XLChartTypeRadarFilled, nil
	case "stockhlc":
		return XLChartTypeStockHLC, nil
	case "stockohlc":
		return XLChartTypeStockOHLC, nil
	case "combo":
		return XLChartTypeCombo, nil
	default:
		return "", fmt.Errorf("invalid XL_CHART_TYPE value %q", value)
	}
}
