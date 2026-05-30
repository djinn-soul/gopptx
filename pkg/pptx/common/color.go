package common

import (
	"regexp"
	"strings"
)

var hexColorPattern = regexp.MustCompile(`^[0-9A-F]{6}$`)

// NormalizeHexColor sanitizes hex color strings and expands 3-digit codes.
// Fast-path: already-normalized 6-char uppercase hex returns without allocation.
func NormalizeHexColor(color string) string {
	const hex6DigitLen = 6
	if len(color) == hex6DigitLen && isUpperHex6(color) {
		return color
	}
	clean := strings.TrimPrefix(strings.TrimSpace(color), "#")
	const hex3DigitLen = 3
	if len(clean) == hex3DigitLen {
		return strings.ToUpper(string([]byte{clean[0], clean[0], clean[1], clean[1], clean[2], clean[2]}))
	}
	return strings.ToUpper(clean)
}

func isUpperHex6(s string) bool {
	for i := range 6 {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}

// IsHexColor checks if a string is a valid 6-digit RGB hex color.
func IsHexColor(color string) bool {
	return hexColorPattern.MatchString(NormalizeHexColor(color))
}
