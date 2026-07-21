package shape

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const colorHexLength = 6
const (
	arrowTypeNone     = "none"
	arrowTypeTriangle = "triangle"
	arrowTypeStealth  = "stealth"
	arrowTypeDiamond  = "diamond"
	arrowTypeOval     = "oval"
	arrowTypeArrow    = "arrow"
)

func GetStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ParseIntAttr(value *string) int {
	if value == nil {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(*value))
	if err != nil {
		return 0
	}
	return n
}

func ParseXMLBoolAttr(value *string) bool {
	if value == nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(*value)) {
	case "1", boolTrueLiteral, "on", "yes":
		return true
	default:
		return false
	}
}

// XMLEscape replaces XML-sensitive characters with entity references.
//
// This uses the shared package-level [strings.Replacer] rather than
// [xml.EscapeText]: the latter allocated a fresh buffer and a []byte copy per
// call (~6.7x slower, 20 allocs vs 2). The two differ in that xml.EscapeText
// also escapes \n, \r and \t; the sole caller escapes preset geometry tokens
// (e.g. "pct5"), which never contain whitespace, so the behavior is equivalent
// here.
func XMLEscape(value string) string {
	return xmlEscapeReplacer.Replace(value)
}

//nolint:gochecknoglobals // Replacer is stateless and reused to avoid per-call allocation.
var xmlEscapeReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&apos;",
)

func NormalizeHexColor(raw string) (string, error) {
	color := strings.TrimSpace(strings.TrimPrefix(raw, "#"))
	if len(color) != colorHexLength {
		return "", fmt.Errorf("expected 6 hex digits, got %q", raw)
	}
	for _, ch := range color {
		if (ch < '0' || ch > '9') && (ch < 'a' || ch > 'f') && (ch < 'A' || ch > 'F') {
			return "", fmt.Errorf("expected hex color, got %q", raw)
		}
	}
	return strings.ToUpper(color), nil
}

func NormalizeLineDashStyle(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", errors.New("must be non-empty")
	}
	switch s {
	case lineDashSolid,
		lineDashDash,
		lineDashDashDot,
		lineDashLgDash,
		lineDashLgDashDot,
		lineDashLgDashDotDot,
		lineDashSysDot,
		lineDashSysDash,
		lineDashSysDashDot,
		lineDashSysDashDotDot:
		return s, nil
	}
	key := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(s, "-", "_"), " ", "_"))
	key = strings.ReplaceAll(key, "__", "_")
	aliases := map[string]string{
		"dash_dot":          lineDashDashDot,
		"dashdot":           lineDashDashDot,
		"dash_dot_dot":      lineDashLgDashDotDot,
		"dashdotdot":        lineDashLgDashDotDot,
		"long_dash":         lineDashLgDash,
		"longdash":          lineDashLgDash,
		"long_dash_dot":     lineDashLgDashDot,
		"longdashdot":       lineDashLgDashDot,
		"long_dash_dot_dot": lineDashLgDashDotDot,
		"longdashdotdot":    lineDashLgDashDotDot,
		"round_dot":         lineDashSysDot,
		"rounddot":          lineDashSysDot,
		"square_dot":        lineDashSysDash,
		"squaredot":         lineDashSysDash,
		"sys_dash":          lineDashSysDash,
		"sysdash":           lineDashSysDash,
		"sys_dot":           lineDashSysDot,
		"sysdot":            lineDashSysDot,
		"sys_dash_dot":      lineDashSysDashDot,
		"sysdashdot":        lineDashSysDashDot,
		"sys_dash_dot_dot":  lineDashSysDashDotDot,
		"sysdashdotdot":     lineDashSysDashDotDot,
	}
	if normalized, ok := aliases[key]; ok {
		return normalized, nil
	}
	return "", fmt.Errorf("unsupported value %q", raw)
}

func NormalizeArrowType(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", errors.New("must be non-empty")
	}
	switch s {
	case arrowTypeNone, arrowTypeTriangle, arrowTypeStealth, arrowTypeDiamond, arrowTypeOval, arrowTypeArrow:
		return s, nil
	}
	key := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(s, "-", "_"), " ", "_"))
	key = strings.ReplaceAll(key, "__", "_")
	switch key {
	case "n":
		return arrowTypeNone, nil
	case "t", arrowTypeTriangle:
		return arrowTypeTriangle, nil
	case "s", "stealth":
		return arrowTypeStealth, nil
	case "d", "diamond":
		return arrowTypeDiamond, nil
	case "o", "oval":
		return arrowTypeOval, nil
	case "a", "open", arrowTypeArrow:
		return arrowTypeArrow, nil
	default:
		return "", fmt.Errorf("unsupported value %q", raw)
	}
}

func NormalizeArrowSize(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", errors.New("must be non-empty")
	}
	switch s {
	case "sm", "med", "lg":
		return s, nil
	}
	key := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(s, "-", "_"), " ", "_"))
	switch key {
	case "s", "small":
		return "sm", nil
	case "m", "medium":
		return "med", nil
	case "l", "large":
		return "lg", nil
	default:
		return "", fmt.Errorf("unsupported value %q", raw)
	}
}

// boolTrueLiteral is the spelled-out true accepted by the bool attribute parsers.
const boolTrueLiteral = "true"

// Canonical OOXML line dash style tokens. NormalizeLineDashStyle accepts these
// directly and also maps a set of spelled-out aliases onto them.
const (
	lineDashSolid         = "solid"
	lineDashDash          = "dash"
	lineDashDashDot       = "dashDot"
	lineDashLgDash        = "lgDash"
	lineDashLgDashDot     = "lgDashDot"
	lineDashLgDashDotDot  = "lgDashDotDot"
	lineDashSysDot        = "sysDot"
	lineDashSysDash       = "sysDash"
	lineDashSysDashDot    = "sysDashDot"
	lineDashSysDashDotDot = "sysDashDotDot"
)
