package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

type (
	TextRun            = text.TextRun
	TextParagraphStyle = text.TextParagraphStyle
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

func NewTextRun(txt string) TextRun {
	return text.NewTextRun(txt)
}

func NewTextParagraphStyle() TextParagraphStyle {
	return text.NewTextParagraphStyle()
}

func DefaultTextParagraphStyle() TextParagraphStyle {
	return text.DefaultTextParagraphStyle()
}

func NormalizeTextAlign(align string) string {
	return text.NormalizeTextAlign(align)
}

func NormalizeTextRuns(runs []TextRun) []TextRun {
	return text.NormalizeTextRuns(runs)
}
