package elements

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

type (
	// TextRun describes a single piece of text with uniform styling.
	TextRun = text.TextRun

	// TextParagraphStyle describes paragraph-level formatting for one bullet line.
	TextParagraphStyle = text.TextParagraphStyle
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

func NewTextRun(txt string) TextRun {
	return text.NewTextRun(txt)
}

func DefaultTextParagraphStyle() TextParagraphStyle {
	return text.DefaultTextParagraphStyle()
}

func NormalizeTextAlign(align string) string {
	return text.NormalizeTextAlign(align)
}

func NormalizeBulletStyle(style string) string {
	return text.NormalizeBulletStyle(style)
}

func NormalizeTextParagraphStyle(style TextParagraphStyle) TextParagraphStyle {
	return text.NormalizeTextParagraphStyle(style)
}

func NormalizeTextRuns(runs []TextRun) []TextRun {
	return text.NormalizeTextRuns(runs)
}

func RunsToPlainText(runs []TextRun) string {
	return text.RunsToPlainText(runs)
}
