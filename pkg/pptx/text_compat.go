package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

type (
	// Run is an alias for text.Run.
	Run = text.Run
	// ParagraphStyle is an alias for text.ParagraphStyle.
	ParagraphStyle = text.ParagraphStyle
	// Paragraph is an alias for text.Paragraph.
	Paragraph = text.Paragraph

	TextRun            = Run
	TextParagraphStyle = ParagraphStyle
	TextParagraph      = Paragraph
)

const (
	TextAlignLeft    = text.TextAlignLeft
	TextAlignCenter  = text.TextAlignCenter
	TextAlignRight   = text.TextAlignRight
	TextAlignJustify = text.TextAlignJustify

	BulletStyleBullet      = text.BulletStyleBullet
	BulletStyleNumber      = text.BulletStyleNumber
	BulletStyleLetterLower = text.BulletStyleLetterLower
	BulletStyleLetterUpper = text.BulletStyleLetterUpper
	BulletStyleRomanLower  = text.BulletStyleRomanLower
	BulletStyleRomanUpper  = text.BulletStyleRomanUpper
	BulletStyleCustom      = text.BulletStyleCustom
	BulletStyleNone        = text.BulletStyleNone

	TextSizeTitle    = text.TextSizeTitle
	TextSizeSubtitle = text.TextSizeSubtitle
	TextSizeHeading  = text.TextSizeHeading
	TextSizeBody     = text.TextSizeBody
	TextSizeSmall    = text.TextSizeSmall
	TextSizeCaption  = text.TextSizeCaption
	TextSizeCode     = text.TextSizeCode
	TextSizeLarge    = text.TextSizeLarge
	TextSizeXLarge   = text.TextSizeXLarge
)

// NewRun creates a new text run.
func NewRun(txt string) Run {
	return text.NewRun(txt)
}

// NewTextRun is an alias for NewRun.
func NewTextRun(txt string) TextRun {
	return NewRun(txt)
}

// NewParagraph creates a new text paragraph.
func NewParagraph() Paragraph {
	return text.NewParagraph()
}

// NewTextParagraph is an alias for NewParagraph.
func NewTextParagraph() TextParagraph {
	return NewParagraph()
}

// NewParagraphStyle creates a new paragraph style.
func NewParagraphStyle() ParagraphStyle {
	return text.NewParagraphStyle()
}

// NewTextParagraphStyle is an alias for NewParagraphStyle.
func NewTextParagraphStyle() TextParagraphStyle {
	return NewParagraphStyle()
}

// DefaultParagraphStyle returns the standard paragraph styling.
func DefaultParagraphStyle() ParagraphStyle {
	return text.DefaultParagraphStyle()
}

// DefaultTextParagraphStyle is an alias for DefaultParagraphStyle.
func DefaultTextParagraphStyle() TextParagraphStyle {
	return DefaultParagraphStyle()
}

func NormalizeTextAlign(align string) string {
	return text.NormalizeTextAlign(align)
}

// NormalizeRuns combines adjacent runs with identical styling.
func NormalizeRuns(runs []Run) []Run {
	return text.NormalizeRuns(runs)
}

// NormalizeTextRuns is an alias for NormalizeRuns.
func NormalizeTextRuns(runs []TextRun) []Run {
	return NormalizeRuns(runs)
}
