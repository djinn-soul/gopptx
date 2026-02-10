package charts

import "strings"

func NormalizeHexColor(color string) string {
	clean := strings.TrimPrefix(strings.TrimSpace(color), "#")
	return strings.ToUpper(clean)
}

func IsHexColor(color string) bool {
	return hexColorPattern.MatchString(NormalizeHexColor(color))
}

func IsLegendPosition(position string) bool {
	switch strings.ToLower(strings.TrimSpace(position)) {
	case LegendPositionRight, LegendPositionLeft, LegendPositionTop, LegendPositionBottom:
		return true
	default:
		return false
	}
}

func IsValueAxisCrossBetween(mode string) bool {
	switch strings.TrimSpace(mode) {
	case ValueAxisCrossBetweenBetween, ValueAxisCrossBetweenMidCategory:
		return true
	default:
		return false
	}
}

func CopyStringSlice(s []string) []string {
	if s == nil {
		return nil
	}
	res := make([]string, len(s))
	copy(res, s)
	return res
}

func CopyFloat64Slice(s []float64) []float64 {
	if s == nil {
		return nil
	}
	res := make([]float64, len(s))
	copy(res, s)
	return res
}

func CopyFloat64Pointer(p *float64) *float64 {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}
