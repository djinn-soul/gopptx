package shape

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const colorHexLength = 6
const (
	arrowTypeNone    = "none"
	arrowTypeStealth = "stealth"
	arrowTypeDiamond = "diamond"
	arrowTypeOval    = "oval"
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
	case "1", "true", "on", "yes":
		return true
	default:
		return false
	}
}

func XMLEscape(value string) string {
	var buf bytes.Buffer
	if err := xml.EscapeText(&buf, []byte(value)); err != nil {
		return value
	}
	return buf.String()
}

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
	case "solid", "dash", "dashDot", "lgDash", "lgDashDot", "lgDashDotDot", "sysDot", "sysDash",
		"sysDashDot", "sysDashDotDot":
		return s, nil
	}
	key := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(s, "-", "_"), " ", "_"))
	key = strings.ReplaceAll(key, "__", "_")
	aliases := map[string]string{
		"dash_dot":          "dashDot",
		"dashdot":           "dashDot",
		"dash_dot_dot":      "lgDashDotDot",
		"dashdotdot":        "lgDashDotDot",
		"long_dash":         "lgDash",
		"longdash":          "lgDash",
		"long_dash_dot":     "lgDashDot",
		"longdashdot":       "lgDashDot",
		"long_dash_dot_dot": "lgDashDotDot",
		"longdashdotdot":    "lgDashDotDot",
		"round_dot":         "sysDot",
		"rounddot":          "sysDot",
		"square_dot":        "sysDash",
		"squaredot":         "sysDash",
		"sys_dash":          "sysDash",
		"sysdash":           "sysDash",
		"sys_dot":           "sysDot",
		"sysdot":            "sysDot",
		"sys_dash_dot":      "sysDashDot",
		"sysdashdot":        "sysDashDot",
		"sys_dash_dot_dot":  "sysDashDotDot",
		"sysdashdotdot":     "sysDashDotDot",
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
	case arrowTypeNone, "triangle", arrowTypeStealth, arrowTypeDiamond, arrowTypeOval, "arrow":
		return s, nil
	}
	key := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(s, "-", "_"), " ", "_"))
	key = strings.ReplaceAll(key, "__", "_")
	switch key {
	case "n":
		return arrowTypeNone, nil
	case "t", "triangle":
		return "triangle", nil
	case "s", "stealth":
		return arrowTypeStealth, nil
	case "d", "diamond":
		return arrowTypeDiamond, nil
	case "o", "oval":
		return arrowTypeOval, nil
	case "a", "open", "arrow":
		return "arrow", nil
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
