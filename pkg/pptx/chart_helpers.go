package pptx

import "strings"

func normalizeHexColor(color string) string {
	clean := strings.TrimPrefix(strings.TrimSpace(color), "#")
	return strings.ToUpper(clean)
}

func isHexColor(color string) bool {
	return hexColorPattern.MatchString(normalizeHexColor(color))
}

func isLegendPosition(position string) bool {
	switch strings.ToLower(strings.TrimSpace(position)) {
	case LegendPositionRight, LegendPositionLeft, LegendPositionTop, LegendPositionBottom:
		return true
	default:
		return false
	}
}
