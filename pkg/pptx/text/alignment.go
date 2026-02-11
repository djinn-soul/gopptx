package text

import "strings"

const (
	// TextAlignLeft aligns paragraph text to the left.
	TextAlignLeft = "l"
	// TextAlignCenter aligns paragraph text to the center.
	TextAlignCenter = "ctr"
	// TextAlignRight aligns paragraph text to the right.
	TextAlignRight = "r"
	// TextAlignJustify aligns paragraph text with justification.
	TextAlignJustify = "just"
)

// NormalizeTextAlign sanitizes alignment strings.
func NormalizeTextAlign(align string) string {
	return strings.ToLower(strings.TrimSpace(align))
}
