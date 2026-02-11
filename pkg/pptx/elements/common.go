package elements

import (
	"regexp"
	"strings"
)

var hexColorPattern = regexp.MustCompile(`^[0-9A-F]{6}$`)

// NormalizeHexColor sanitizes hex color strings.
func NormalizeHexColor(color string) string {
	clean := strings.TrimPrefix(strings.TrimSpace(color), "#")
	return strings.ToUpper(clean)
}

// IsHexColor checks if a string is a valid 6-digit RGB hex color.
func IsHexColor(color string) bool {
	return hexColorPattern.MatchString(NormalizeHexColor(color))
}
