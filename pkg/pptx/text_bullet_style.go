package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

const (
	BulletStyleBullet      = elements.BulletStyleBullet
	BulletStyleNumber      = elements.BulletStyleNumber
	BulletStyleLetterLower = elements.BulletStyleLetterLower
	BulletStyleLetterUpper = elements.BulletStyleLetterUpper
	BulletStyleRomanLower  = elements.BulletStyleRomanLower
	BulletStyleRomanUpper  = elements.BulletStyleRomanUpper
	BulletStyleCustom      = elements.BulletStyleCustom
	BulletStyleNone        = elements.BulletStyleNone
)

func DefaultTextParagraphStyle() TextParagraphStyle {
	return elements.NormalizeTextParagraphStyle(TextParagraphStyle{})
}

func normalizeBulletStyle(style string) string {
	return elements.NormalizeBulletStyle(style)
}
