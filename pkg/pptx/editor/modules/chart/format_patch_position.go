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
			return "outEnd"
		}
		return "r"
	case "left", "l":
		if isBarChart {
			return "inBase"
		}
		return "l"
	case "top", "t":
		if isBarChart {
			return "outEnd"
		}
		return "t"
	case "bottom", "b":
		if isBarChart {
			return "inBase"
		}
		return "b"
	case "center", "ctr":
		return "ctr"
	case "outside_end", "outsideend", "outend":
		return "outEnd"
	case "inside_end", "insideend", "inend":
		return "inEnd"
	case "inside_base", "insidebase", "inbase":
		return "inBase"
	case "best_fit", "bestfit":
		return "bestFit"
	default:
		return pos
	}
}
