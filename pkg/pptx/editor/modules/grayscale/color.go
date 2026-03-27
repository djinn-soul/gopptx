package grayscale

import (
	"fmt"
	"strconv"
	"strings"
)

const rgbHexLen = 6

// HexColor converts an RGB hex string into its grayscale equivalent.
func HexColor(raw string) (string, error) {
	clean := strings.TrimPrefix(strings.TrimSpace(raw), "#")
	if len(clean) != rgbHexLen {
		return "", fmt.Errorf("expected %d-digit RGB hex, got %q", rgbHexLen, raw)
	}
	r, err := strconv.ParseUint(clean[0:2], 16, 8)
	if err != nil {
		return "", fmt.Errorf("parse red channel: %w", err)
	}
	g, err := strconv.ParseUint(clean[2:4], 16, 8)
	if err != nil {
		return "", fmt.Errorf("parse green channel: %w", err)
	}
	b, err := strconv.ParseUint(clean[4:6], 16, 8)
	if err != nil {
		return "", fmt.Errorf("parse blue channel: %w", err)
	}
	luma := int((299*int(r) + 587*int(g) + 114*int(b) + 500) / 1000)
	return fmt.Sprintf("%02X%02X%02X", luma, luma, luma), nil
}
