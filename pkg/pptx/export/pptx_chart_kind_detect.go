package export

import (
	"encoding/xml"
	"strings"
)

const (
	chartGroupingStacked = "stacked"
	stockOHLCSeriesMin   = 4
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
			info.barDir = xmlAttrVal(start, "val")
		case "grouping":
			if info.grouping == "" {
				info.grouping = xmlAttrVal(start, "val")
			}
		case "radarStyle":
			info.radarStyle = xmlAttrVal(start, "val")
		case "scatterStyle":
			info.scatterStyle = xmlAttrVal(start, "val")
		case "symbol":
			// Only treat markers as visible when the symbol is not "none".
			// Line charts emit <c:marker><c:symbol val="none"/> to explicitly
			// suppress markers; lineMarkers charts use val="circle" (or similar).
			if val := xmlAttrVal(start, "val"); val != "" && val != "none" {
				info.hasMarker = true
			}
		case "ser":
			info.seriesCount++
		}
	}
	return info
}

func xmlAttrVal(start xml.StartElement, name string) string {
	for _, attr := range start.Attr {
		if attr.Name.Local == name {
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
			return "stockOHLC"
		}
		return "stockHLC"
	case info.chartTypes["barChart"] && info.chartTypes["lineChart"]:
		return "combo"
	case info.chartTypes["doughnutChart"]:
		return "doughnut"
	case info.chartTypes["pieChart"]:
		return parsedChartKindPie
	case info.chartTypes["bubbleChart"]:
		return "bubble"
	case info.chartTypes["scatterChart"]:
		return "scatter"
	case info.chartTypes["radarChart"]:
		if info.radarStyle == "filled" {
			return "radarFilled"
		}
		return "radar"
	case info.chartTypes["areaChart"]:
		if info.grouping == "percentStacked" {
			return "areaStacked100"
		}
		if info.grouping == chartGroupingStacked {
			return "areaStacked"
		}
		return "area"
	case info.chartTypes["lineChart"]:
		if info.grouping == chartGroupingStacked {
			return "lineStacked"
		}
		if info.hasMarker {
			return "lineMarkers"
		}
		return "line"
	case info.chartTypes["barChart"]:
		if info.grouping == "percentStacked" {
			return "barStacked100"
		}
		if info.grouping == chartGroupingStacked {
			return "barStacked"
		}
		if info.barDir == parsedChartKindBar {
			return "barHorizontal"
		}
		return parsedChartKindBar
	default:
		return parsedChartKindBar
	}
}
