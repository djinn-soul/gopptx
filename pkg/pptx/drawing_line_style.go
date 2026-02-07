package pptx

import "strings"

const (
	// LineDashSolid emits a solid line.
	LineDashSolid = "solid"
	// LineDashDash emits a dashed line.
	LineDashDash = "dash"
	// LineDashDot emits a dotted line.
	LineDashDot = "dot"
	// LineDashDashDot emits a dash-dot line.
	LineDashDashDot = "dashDot"
	// LineDashDashDotDot emits a dash-dot-dot line.
	LineDashDashDotDot = "lgDashDotDot"
	// LineDashLongDash emits a long-dash line.
	LineDashLongDash = "lgDash"
	// LineDashLongDashDot emits a long-dash-dot line.
	LineDashLongDashDot = "lgDashDot"
)

func normalizeDrawingLineDash(dash string) string {
	switch strings.ToLower(strings.TrimSpace(dash)) {
	case "", "solid":
		return LineDashSolid
	case "dash":
		return LineDashDash
	case "dot":
		return LineDashDot
	case "dashdot", "dash-dot", "dash_dot":
		return LineDashDashDot
	case "dashdotdot", "dash-dot-dot", "dash_dot_dot", "lgdashdotdot", "lg-dash-dot-dot":
		return LineDashDashDotDot
	case "lgdash", "lg-dash", "longdash", "long-dash", "long_dash":
		return LineDashLongDash
	case "lgdashdot", "lg-dash-dot", "longdashdot", "long-dash-dot", "long_dash_dot":
		return LineDashLongDashDot
	default:
		return strings.TrimSpace(dash)
	}
}

func isDrawingLineDash(dash string) bool {
	switch normalizeDrawingLineDash(dash) {
	case LineDashSolid, LineDashDash, LineDashDot, LineDashDashDot, LineDashDashDotDot, LineDashLongDash, LineDashLongDashDot:
		return true
	default:
		return false
	}
}
