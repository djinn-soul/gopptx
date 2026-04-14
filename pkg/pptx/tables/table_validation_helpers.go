package tables

import "strings"

func normalizeTableAlign(align string) string {
	return strings.ToLower(strings.TrimSpace(align))
}

func normalizeTableVAlign(vAlign string) string {
	return strings.ToLower(strings.TrimSpace(vAlign))
}

func isTableAlign(align string) bool {
	switch normalizeTableAlign(align) {
	case TableAlignLeft, TableAlignCenter, TableAlignRight, TableAlignJustify:
		return true
	default:
		return false
	}
}

func isTableVAlign(vAlign string) bool {
	switch normalizeTableVAlign(vAlign) {
	case TableVAlignTop, TableVAlignMiddle, TableVAlignBottom:
		return true
	default:
		return false
	}
}

// NormalizeTableBorderDash sanitizes table border dash styles.
func NormalizeTableBorderDash(dash string) string {
	switch strings.ToLower(strings.TrimSpace(dash)) {
	case "", "solid":
		return TableBorderDashSolid
	case "dash":
		return TableBorderDashDash
	case "dot":
		return TableBorderDashDot
	case "dashdot", "dash-dot", "dash_dot":
		return TableBorderDashDashDot
	case "lgdash", "lg-dash", "longdash", "long-dash", "long_dash":
		return TableBorderDashLongDash
	default:
		return strings.TrimSpace(dash)
	}
}

func isTableBorderDash(dash string) bool {
	switch NormalizeTableBorderDash(dash) {
	case TableBorderDashSolid, TableBorderDashDash, TableBorderDashDot, TableBorderDashDashDot, TableBorderDashLongDash:
		return true
	default:
		return false
	}
}
