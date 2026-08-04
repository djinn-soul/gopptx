package enums

//go:generate go run ../../../cmd/gen_chart_types shape_chart.go ../../../internal/pptxxml/chart_spec.go ../../../python/gopptx/presentation/charts/chart_types.py

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
	XLChartTypeThreeDPie      XLChartType = pptxxml.ChartKindThreeDPie
	XLChartTypeThreeDColumn   XLChartType = pptxxml.ChartKindThreeDColumn
	XLChartTypeThreeDBar      XLChartType = pptxxml.ChartKindThreeDBar
	XLChartTypeThreeDLine     XLChartType = pptxxml.ChartKindThreeDLine
	XLChartTypeThreeDArea     XLChartType = pptxxml.ChartKindThreeDArea
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

// xlChartTypeByName maps every accepted spelling to its XLChartType. A table
// keeps the alias list additive: a new chart kind adds rows instead of growing
// a switch past the statement limit.
//
//nolint:gochecknoglobals // immutable lookup table
var xlChartTypeByName = map[string]XLChartType{
	"bar":           XLChartTypeBar,
	"barhorizontal": XLChartTypeBarHorizontal, "barh": XLChartTypeBarHorizontal,
	"barstacked":    XLChartTypeBarStacked,
	"barstacked100": XLChartTypeBarStacked100,
	"line":          XLChartTypeLine,
	"linemarkers":   XLChartTypeLineMarkers,
	"linestacked":   XLChartTypeLineStacked,
	"scatter":       XLChartTypeScatter, "xy": XLChartTypeScatter,
	"area":           XLChartTypeArea,
	"areastacked":    XLChartTypeAreaStacked,
	"areastacked100": XLChartTypeAreaStacked100,
	"pie":            XLChartTypePie,
	"pie3d":          XLChartTypeThreeDPie, "threedpie": XLChartTypeThreeDPie,
	"three_d_pie": XLChartTypeThreeDPie, "three-d-pie": XLChartTypeThreeDPie,
	"column3d": XLChartTypeThreeDColumn, "threedcolumn": XLChartTypeThreeDColumn,
	"three_d_column": XLChartTypeThreeDColumn, "three-d-column": XLChartTypeThreeDColumn,
	"bar3d": XLChartTypeThreeDBar, "threedbar": XLChartTypeThreeDBar,
	"three_d_bar": XLChartTypeThreeDBar, "three-d-bar": XLChartTypeThreeDBar,
	"line3d": XLChartTypeThreeDLine, "threedline": XLChartTypeThreeDLine,
	"three_d_line": XLChartTypeThreeDLine, "three-d-line": XLChartTypeThreeDLine,
	"area3d": XLChartTypeThreeDArea, "threedarea": XLChartTypeThreeDArea,
	"three_d_area": XLChartTypeThreeDArea, "three-d-area": XLChartTypeThreeDArea,
	chartTypeNameDoughnut: XLChartTypeDoughnut, "donut": XLChartTypeDoughnut,
	"bubble":      XLChartTypeBubble,
	"radar":       XLChartTypeRadar,
	"radarfilled": XLChartTypeRadarFilled,
	"stockhlc":    XLChartTypeStockHLC,
	"stockohlc":   XLChartTypeStockOHLC,
	"combo":       XLChartTypeCombo,
}

// ParseXLChartType resolves any accepted spelling of a chart type.
func ParseXLChartType(value string) (XLChartType, error) {
	if chartType, ok := xlChartTypeByName[normalizeKey(value)]; ok {
		return chartType, nil
	}
	return "", fmt.Errorf("invalid XL_CHART_TYPE value %q", value)
}

// chartTypeNameDoughnut is the canonical doughnut chart type token.
const chartTypeNameDoughnut = "doughnut"
