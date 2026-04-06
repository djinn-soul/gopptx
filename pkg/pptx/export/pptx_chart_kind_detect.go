package export

import (
	"encoding/xml"
	"strings"
)

const (
	chartGroupingStacked = "stacked"
	stockOHLCSeriesMin   = 4
	chartKindNone        = "none"
	chartKindArea        = "area"
	chartKindAreaStacked = "areaStacked"
	chartKindArea100     = "areaStacked100"
	chartKindLine        = "line"
	chartKindLineMarkers = "lineMarkers"
	chartKindLineStacked = "lineStacked"
	chartKindBar         = "bar"
	chartKindBarStacked  = "barStacked"
	chartKindBar100      = "barStacked100"
	chartKindBarHoriz    = "barHorizontal"
	chartKindPie         = "pie"
	chartKindDoughnut    = "doughnut"
	chartKindBubble      = "bubble"
	chartKindScatter     = "scatter"
	chartKindRadar       = "radar"
	chartKindRadarFilled = "radarFilled"
	chartKindStockHLC    = "stockHLC"
	chartKindStockOHLC   = "stockOHLC"
	chartKindCombo       = "combo"
	// xmlBarDirHoriz is the OOXML barDir attribute value for horizontal bars.
	// Distinct from chartKindBar, which is the gopptx chart kind for column (vertical) charts.
	xmlBarDirHoriz = chartKindBar
)

// chartXMLInfo holds structural facts extracted from a chart XML part.
type chartXMLInfo struct {
	chartTypes   map[string]bool
	barDir       string
	grouping     string
	radarStyle   string
	scatterStyle string
	hasMarker    bool
	seriesCount  int
}

// scanChartXMLInfo walks XML tokens to collect chart type and key settings.
func scanChartXMLInfo(rawXML string) chartXMLInfo {
	info := chartXMLInfo{chartTypes: make(map[string]bool)}
	dec := xml.NewDecoder(strings.NewReader(rawXML))
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		switch start.Name.Local {
		case "stockChart", "barChart", "lineChart", "doughnutChart", "pieChart",
			"bubbleChart", "scatterChart", "radarChart", "areaChart":
			info.chartTypes[start.Name.Local] = true
		case "barDir":
			info.barDir = xmlAttrVal(start)
		case "grouping":
			if info.grouping == "" {
				info.grouping = xmlAttrVal(start)
			}
		case "radarStyle":
			info.radarStyle = xmlAttrVal(start)
		case "scatterStyle":
			info.scatterStyle = xmlAttrVal(start)
		case "symbol":
			// Only treat markers as visible when the symbol is not "none".
			// Line charts emit <c:marker><c:symbol val="none"/> to explicitly
			// suppress markers; lineMarkers charts use val="circle" (or similar).
			if val := xmlAttrVal(start); val != "" && val != chartKindNone {
				info.hasMarker = true
			}
		case "ser":
			info.seriesCount++
		}
	}
	return info
}

func xmlAttrVal(start xml.StartElement) string {
	for _, attr := range start.Attr {
		if attr.Name.Local == "val" {
			return attr.Value
		}
	}
	return ""
}

func detectChartKind(rawXML string) string {
	info := scanChartXMLInfo(rawXML)
	switch {
	case info.chartTypes["stockChart"]:
		if info.seriesCount >= stockOHLCSeriesMin {
			return chartKindStockOHLC
		}
		return chartKindStockHLC
	case info.chartTypes["barChart"] && info.chartTypes["lineChart"]:
		return chartKindCombo
	case info.chartTypes["doughnutChart"]:
		return chartKindDoughnut
	case info.chartTypes["pieChart"]:
		return chartKindPie
	case info.chartTypes["bubbleChart"]:
		return chartKindBubble
	case info.chartTypes["scatterChart"]:
		return chartKindScatter
	case info.chartTypes["radarChart"]:
		if info.radarStyle == "filled" {
			return chartKindRadarFilled
		}
		return chartKindRadar
	case info.chartTypes["areaChart"]:
		if info.grouping == "percentStacked" {
			return chartKindArea100
		}
		if info.grouping == chartGroupingStacked {
			return chartKindAreaStacked
		}
		return chartKindArea
	case info.chartTypes["lineChart"]:
		if info.grouping == chartGroupingStacked {
			return chartKindLineStacked
		}
		if info.hasMarker {
			return chartKindLineMarkers
		}
		return chartKindLine
	case info.chartTypes["barChart"]:
		if info.grouping == "percentStacked" {
			return chartKindBar100
		}
		if info.grouping == chartGroupingStacked {
			return chartKindBarStacked
		}
		if info.barDir == xmlBarDirHoriz {
			return chartKindBarHoriz
		}
		return chartKindBar
	default:
		return chartKindBar
	}
}
