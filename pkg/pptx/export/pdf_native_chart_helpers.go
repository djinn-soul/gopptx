//nolint:mnd // Chart helper math uses tuned visual constants for native PDF fidelity.
package export

import (
	"math"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

// legendEntry describes one entry in a chart legend (series name + color).
type legendEntry struct {
	Name    string
	R, G, B uint8
}

// chartSeriesOpts holds optional rendering hints for chart renderers.
type chartSeriesOpts struct {
	color          string   // hex color override; empty = use renderer default
	minValue       *float64 // axis min override
	maxValue       *float64 // axis max override
	showLegend     bool
	legendPosition string // "r","l","t","b"
	seriesName     string
	showDataLabels bool
	dataLabels     charts.DataLabelSettings // c:dLbls content and placement
	// categories and valueTotal are filled in by the renderer just before it
	// draws, so a data label can name its category or state its share of the
	// series without every drawing helper having to carry them.
	categories []string
	valueTotal float64
	// plot is the solved plot rect: data labels wrap against it and are kept
	// inside it. hasPlot says whether the renderer has supplied it yet.
	plot               chartRect
	hasPlot            bool
	labelWrapWidth     float64
	showCatName        bool // explicitly show category names in data labels (Pie/Doughnut)
	catAxisTitle       string
	valAxisTitle       string
	scatterStyle       string // "marker" | "lineMarker" | "smoothMarker"
	smooth             bool   // draw line chart with Catmull-Rom smooth curves
	showMajorGridlines bool   // draw horizontal value-axis gridlines
	showCatGridlines   bool   // draw vertical category-axis gridlines (horizontal charts)
	titleOverlay       bool   // title overlaps plot area; don't reserve top padding
	valueFormat        string // Excel-style number format ("General", "0%", "$#,##0", …)
	bubbleScale        int    // bubble size scale percent (1–300; 0 = use renderer default)
}

// categoryLabel returns the label for index i, falling back to "Q<i+1>".
func categoryLabel(categories []string, i int) string {
	if i < len(categories) && categories[i] != "" {
		return categories[i]
	}
	return "Q" + strconv.Itoa(i+1)
}

// formatTickValue formats a numeric axis-tick value using the given Excel-style format.
// Supported: "General"/empty → rounded integer; contains "%" → append "%";
// starts with "$" → prepend "$". Anything else falls back to rounded integer.
func formatTickValue(v float64, format string) string {
	if format == "" || format == chartGeneralNumberFormat {
		return strconv.Itoa(int(math.Round(v)))
	}
	if strings.Contains(format, "%") {
		return strconv.Itoa(int(math.Round(v))) + "%"
	}
	if strings.HasPrefix(format, "$") {
		return "$" + strconv.FormatFloat(v, 'f', 0, 64)
	}
	return strconv.Itoa(int(math.Round(v)))
}

func pieColor(i int) (uint8, uint8, uint8) {
	palette := [][3]uint8{
		{79, 129, 189},
		{192, 80, 77},
		{155, 187, 89},
		{128, 100, 162},
		{75, 172, 198},
		{247, 150, 70},
	}
	c := palette[i%len(palette)]
	return c[0], c[1], c[2]
}

func drawWedge(pdf *gopdf.GoPdf, cx, cy, radius, start, end float64, r, g, b uint8) {
	pdf.SetFillColor(r, g, b)
	pts := []gopdf.Point{{X: cx, Y: cy}}
	steps := max(8, int(math.Ceil((end-start)/(math.Pi/18))))
	for i := 0; i <= steps; i++ {
		t := start + (end-start)*float64(i)/float64(steps)
		pts = append(pts, gopdf.Point{X: cx + radius*math.Cos(t), Y: cy + radius*math.Sin(t)})
	}
	pdf.Polygon(pts, "F")
}

func maxFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	v := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > v {
			v = values[i]
		}
	}
	return v
}

func sumFloat(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total
}

func minMax(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	minV, maxV := values[0], values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < minV {
			minV = values[i]
		}
		if values[i] > maxV {
			maxV = values[i]
		}
	}
	return minV, maxV
}

// Fallback padding used when no labels are supplied to measure against.
const (
	chartFallbackLeftPadPt   = 36.0
	chartFallbackRightPadPt  = 12.0
	chartFallbackTopPadPt    = 24.0
	chartFallbackBottomPadPt = 26.0
	chartOverlayTopPadPt     = 4.0
	chartAxisTickPt          = 4.0

	// chartMinLeftPadPt is the smallest gutter PowerPoint leaves between the
	// chart edge and the plot area, regardless of how narrow the axis labels are.
	chartMinLeftPadPt = 24.0
)

func widestChartLabel(pdf *gopdf.GoPdf, labels []string) float64 {
	widest := 0.0
	for _, label := range labels {
		if w := chartLabelWidth(pdf, label, chartLabelFontSize); w > widest {
			widest = w
		}
	}
	return widest
}

// chartLegendReservePt is how much width a side legend takes: its marker, the
// widest entry name, and the clearance PowerPoint keeps either side of it.
func chartLegendReservePt(pdf *gopdf.GoPdf, names []string) float64 {
	widest := widestChartLabel(pdf, names)
	if widest <= 0 {
		return chartLegendFallbackWidthPt
	}
	reserve := chartLegendMarkerWPt + chartLegendMarkerGapPt + widest + chartLegendEdgeGapPt
	return math.Min(math.Max(reserve, chartLegendMinWidthPt), chartLegendMaxWidthPt)
}

const (
	// chartLegendFallbackWidthPt is used when the entry names are not known at
	// the point the plot area is sized.
	chartLegendFallbackWidthPt = 70.0
	chartLegendMarkerGapPt     = 4.0
	chartLegendEdgeGapPt       = 12.0
	chartLegendMinWidthPt      = 40.0
	chartLegendMaxWidthPt      = 140.0
)

// chartRectWithLegendMargin shrinks the chart rect to leave room for the
// legend.
//
// The side reserve is measured from the entry names rather than fixed at 110pt:
// PowerPoint gives a legend only the width its labels need, and reserving more
// than that squeezed the plot area and shifted every gridline, bar and label
// left of where PowerPoint draws them.
func chartRectWithLegendMargin(pdf *gopdf.GoPdf, r chartRect, pos string, names []string) chartRect {
	legendW := chartLegendReservePt(pdf, names)
	const legendH = 36
	switch pos {
	case "l":
		return chartRect{r.x + legendW, r.y, r.w - legendW, r.h}
	case "t":
		return chartRect{r.x, r.y + legendH, r.w, r.h - legendH}
	case "b":
		return chartRect{r.x, r.y, r.w, r.h - legendH}
	default: // "r"
		return chartRect{r.x, r.y, r.w - legendW, r.h}
	}
}
