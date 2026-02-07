package pptx

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateSlideTextParagraphStyles(s SlideContent, slideIndex int) error {
	if len(s.BulletStyles) == 0 {
		return nil
	}
	if len(s.BulletStyles) != len(s.Bullets) {
		return fmt.Errorf(
			"slide %d bullet styles count %d must match bullet count %d",
			slideIndex,
			len(s.BulletStyles),
			len(s.Bullets),
		)
	}
	for bulletIndex, style := range s.BulletStyles {
		if align := normalizeTextAlign(style.Align); align != "" && !isTextAlign(align) {
			return fmt.Errorf("slide %d bullet %d align must be one of l|ctr|r|just", slideIndex, bulletIndex+1)
		}
		if style.SpaceBeforePt < 0 {
			return fmt.Errorf("slide %d bullet %d space-before must be >= 0", slideIndex, bulletIndex+1)
		}
		if style.SpaceAfterPt < 0 {
			return fmt.Errorf("slide %d bullet %d space-after must be >= 0", slideIndex, bulletIndex+1)
		}
		if style.LineSpacingPct < 0 {
			return fmt.Errorf("slide %d bullet %d line-spacing must be >= 0", slideIndex, bulletIndex+1)
		}
		if style.Level < 0 || style.Level > maxBulletLevel {
			return fmt.Errorf("slide %d bullet %d level must be between 0 and %d", slideIndex, bulletIndex+1, maxBulletLevel)
		}
		bulletStyle := normalizeBulletStyle(style.BulletStyle)
		if bulletStyle == "" {
			bulletStyle = BulletStyleBullet
		}
		if !isBulletStyle(bulletStyle) {
			return fmt.Errorf("slide %d bullet %d bullet style must be one of bullet|number|letter_lower|letter_upper|roman_lower|roman_upper|custom|none", slideIndex, bulletIndex+1)
		}
		bulletChar := strings.TrimSpace(style.BulletChar)
		if bulletStyle == BulletStyleCustom {
			if utf8.RuneCountInString(bulletChar) != 1 {
				return fmt.Errorf("slide %d bullet %d custom bullet must be a single character", slideIndex, bulletIndex+1)
			}
		} else if bulletChar != "" {
			return fmt.Errorf("slide %d bullet %d bullet char is only allowed with custom bullet style", slideIndex, bulletIndex+1)
		}
	}
	return nil
}

func isTextAlign(align string) bool {
	switch normalizeTextAlign(align) {
	case TextAlignLeft, TextAlignCenter, TextAlignRight, TextAlignJustify:
		return true
	default:
		return false
	}
}

func isBulletStyle(style string) bool {
	switch normalizeBulletStyle(style) {
	case BulletStyleBullet,
		BulletStyleNumber,
		BulletStyleLetterLower,
		BulletStyleLetterUpper,
		BulletStyleRomanLower,
		BulletStyleRomanUpper,
		BulletStyleCustom,
		BulletStyleNone:
		return true
	default:
		return false
	}
}
