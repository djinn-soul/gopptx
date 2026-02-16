package elements

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

type (
	// Run describes a single piece of text with uniform styling.
	Run = text.Run

	// ParagraphStyle describes paragraph-level formatting for one bullet line.
	ParagraphStyle = text.ParagraphStyle

	// Paragraph represents a single paragraph of text with runs and styling.
	Paragraph = text.Paragraph
)

const (
	MaxBulletLevel = text.MaxBulletLevel

	BulletStyleBullet      = text.BulletStyleBullet
	BulletStyleNumber      = text.BulletStyleNumber
	BulletStyleLetterLower = text.BulletStyleLetterLower
	BulletStyleLetterUpper = text.BulletStyleLetterUpper
	BulletStyleRomanLower  = text.BulletStyleRomanLower
	BulletStyleRomanUpper  = text.BulletStyleRomanUpper
	BulletStyleCustom      = text.BulletStyleCustom
	BulletStyleNone        = text.BulletStyleNone

	TextAlignLeft    = text.TextAlignLeft
	TextAlignCenter  = text.TextAlignCenter
	TextAlignRight   = text.TextAlignRight
	TextAlignJustify = text.TextAlignJustify
)

// NewRun creates a new text run.
func NewRun(t string) Run {
	return text.NewRun(t)
}

// NewParagraph creates a new text paragraph with default style.
func NewParagraph() Paragraph {
	return text.NewParagraph()
}

// NewParagraphStyle creates a new paragraph style with default settings.
func NewParagraphStyle() ParagraphStyle {
	return text.NewParagraphStyle()
}

// DefaultParagraphStyle returns the standard paragraph styling.
func DefaultParagraphStyle() ParagraphStyle {
	return text.DefaultParagraphStyle()
}

// NormalizeTextAlign ensures the alignment string is one of the predefined constants.
func NormalizeTextAlign(align string) string {
	return text.NormalizeTextAlign(align)
}

// NormalizeBulletStyle ensures the bullet style string is one of the predefined constants.
func NormalizeBulletStyle(style string) string {
	return text.NormalizeBulletStyle(style)
}

// NormalizeParagraphStyle ensures all fields are within expected bounds.
func NormalizeParagraphStyle(style ParagraphStyle) ParagraphStyle {
	return text.NormalizeParagraphStyle(style)
}

func NormalizeRuns(runs []Run) []Run {
	return text.NormalizeRuns(runs)
}

func RunsToPlainText(runs []Run) string {
	return text.RunsToPlainText(runs)
}
