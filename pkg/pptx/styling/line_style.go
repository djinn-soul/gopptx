package styling

// Shared constants for line dashing and styles (parity with shapes subpackage).
const (
	LineDashSolid   = "solid"
	LineDashDash    = "dash"
	LineDashDot     = "dot"
	LineDashDashDot = "dashDot"
	// LineDashDashDotDot is the dash-dot-dot preset. ST_PresetLineDashVal has no
	// bare "dashDotDot": the dot-dot forms exist only as "lgDashDotDot" and
	// "sysDashDotDot".
	LineDashDashDotDot     = "lgDashDotDot"
	LineDashLongDash       = "lgDash"
	LineDashLongDashDot    = "lgDashDot"
	LineDashLongDashDotDot = "lgDashDotDot"
	LineDashSystemDash     = "sysDash"
	LineDashSystemDot      = "sysDot"
	LineDashSystemDashDot  = "sysDashDot"
	// LineDashSystemDashDotDot is the other dot-dot member of the enumeration.
	LineDashSystemDashDotDot = "sysDashDotDot"
)
