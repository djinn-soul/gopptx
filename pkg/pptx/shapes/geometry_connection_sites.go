package shapes

import "strings"

const (
	connectionSiteTopIndex = iota
	connectionSiteRightIndex
	connectionSiteBottomIndex
	connectionSiteLeftIndex
	connectionSiteTopLeftIndex
	connectionSiteTopRightIndex
	connectionSiteBottomRightIndex
	connectionSiteBottomLeftIndex
	connectionSiteCenterIndex
)

const (
	// ConnectionSiteTop anchors on top-center.
	ConnectionSiteTop = "top"
	// ConnectionSiteRight anchors on right-center.
	ConnectionSiteRight = "right"
	// ConnectionSiteBottom anchors on bottom-center.
	ConnectionSiteBottom = "bottom"
	// ConnectionSiteLeft anchors on left-center.
	ConnectionSiteLeft = "left"
	// ConnectionSiteTopLeft anchors on top-left.
	ConnectionSiteTopLeft = "topLeft"
	// ConnectionSiteTopRight anchors on top-right.
	ConnectionSiteTopRight = "topRight"
	// ConnectionSiteBottomRight anchors on bottom-right.
	ConnectionSiteBottomRight = "bottomRight"
	// ConnectionSiteBottomLeft anchors on bottom-left.
	ConnectionSiteBottomLeft = "bottomLeft"
	// ConnectionSiteCenter anchors on center.
	ConnectionSiteCenter = "center"
)

func NormalizeConnectionSite(site string) string {
	t := strings.ToLower(strings.TrimSpace(site))
	switch t {
	case strings.ToLower(ConnectionSiteTop), "t":
		return ConnectionSiteTop
	case strings.ToLower(ConnectionSiteRight), "r":
		return ConnectionSiteRight
	case strings.ToLower(ConnectionSiteBottom), "b":
		return ConnectionSiteBottom
	case strings.ToLower(ConnectionSiteLeft), "l":
		return ConnectionSiteLeft
	case "topleft", "top-left", "top_left", "tl":
		return ConnectionSiteTopLeft
	case "topright", "top-right", "top_right", "tr":
		return ConnectionSiteTopRight
	case "bottomright", "bottom-right", "bottom_right", "br":
		return ConnectionSiteBottomRight
	case "bottomleft", "bottom-left", "bottom_left", "bl":
		return ConnectionSiteBottomLeft
	case strings.ToLower(ConnectionSiteCenter), string(TextAnchorMiddle), "c":
		return ConnectionSiteCenter
	default:
		return strings.TrimSpace(site)
	}
}

func ConnectionSiteIndex(site string) (int, bool) {
	switch NormalizeConnectionSite(site) {
	case ConnectionSiteTop:
		return connectionSiteTopIndex, true
	case ConnectionSiteRight:
		return connectionSiteRightIndex, true
	case ConnectionSiteBottom:
		return connectionSiteBottomIndex, true
	case ConnectionSiteLeft:
		return connectionSiteLeftIndex, true
	case ConnectionSiteTopLeft:
		return connectionSiteTopLeftIndex, true
	case ConnectionSiteTopRight:
		return connectionSiteTopRightIndex, true
	case ConnectionSiteBottomRight:
		return connectionSiteBottomRightIndex, true
	case ConnectionSiteBottomLeft:
		return connectionSiteBottomLeftIndex, true
	case ConnectionSiteCenter:
		return connectionSiteCenterIndex, true
	default:
		return 0, false
	}
}
