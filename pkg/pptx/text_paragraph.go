package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

const (
	TextAlignLeft    = elements.TextAlignLeft
	TextAlignCenter  = elements.TextAlignCenter
	TextAlignRight   = elements.TextAlignRight
	TextAlignJustify = elements.TextAlignJustify
)

type (
	// TextParagraphStyle describes paragraph-level formatting for one bullet line.
	TextParagraphStyle = elements.TextParagraphStyle
)

// NewTextParagraphStyle creates one paragraph style with default settings.
func NewTextParagraphStyle() TextParagraphStyle {
	return elements.NormalizeTextParagraphStyle(TextParagraphStyle{})
}

func NormalizeTextAlign(align string) string {
	return elements.NormalizeTextAlign(align)
}

func normalizeTextAlign(align string) string {
	return NormalizeTextAlign(align)
}
