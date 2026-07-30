package chart

import "strings"

// normalizeDataLabelPosition maps raw position tokens into valid OOXML ST_DLblPos values.
// For bar charts (<c:barChart>), invalid positions like "r" or "l" are mapped to valid
// bar chart positions ("outEnd" / "inBase") to prevent Office 365 web file corruption (Issue #1134).
func normalizeDataLabelPosition(rawPos string, isBarChart bool) string {
	pos := strings.TrimSpace(rawPos)
	posLower := strings.ToLower(pos)

	switch posLower {
	case "right", "r":
		if isBarChart {
			return dataLabelPositionOutsideEnd
		}
		return "r"
	case "left", "l":
		if isBarChart {
			return dataLabelPositionInsideBase
		}
		return "l"
	case "top", "t":
		if isBarChart {
			return dataLabelPositionOutsideEnd
		}
		return "t"
	case "bottom", "b":
		if isBarChart {
			return dataLabelPositionInsideBase
		}
		return "b"
	case dataLabelPositionCenterInput, dataLabelPositionCenter:
		return dataLabelPositionCenter
	case "outside_end", "outsideend", "outend":
		return dataLabelPositionOutsideEnd
	case "inside_end", "insideend", "inend":
		return dataLabelPositionInsideEnd
	case "inside_base", "insidebase", "inbase":
		return dataLabelPositionInsideBase
	case "best_fit", "bestfit":
		return dataLabelPositionBestFit
	default:
		return pos
	}
}
