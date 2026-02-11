package common

import (
	"regexp"
	"strings"
)

var hexColorPattern = regexp.MustCompile(`^[0-9A-F]{6}$`)

// NormalizeHexColor sanitizes hex color strings and expands 3-digit codes.
func NormalizeHexColor(color string) string {
	clean := strings.TrimPrefix(strings.TrimSpace(color), "#")
	if len(clean) == 3 {
		return strings.ToUpper(string([]byte{clean[0], clean[0], clean[1], clean[1], clean[2], clean[2]}))
	}
	return strings.ToUpper(clean)
}

// IsHexColor checks if a string is a valid 6-digit RGB hex color.
func IsHexColor(color string) bool {
	return hexColorPattern.MatchString(NormalizeHexColor(color))
}
