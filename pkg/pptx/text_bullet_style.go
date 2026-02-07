package pptx

import "strings"

const (
	// BulletStyleBullet renders a standard bullet character.
	BulletStyleBullet = "bullet"
	// BulletStyleNumber renders arabic numbering (1., 2., 3.).
	BulletStyleNumber = "number"
	// BulletStyleLetterLower renders lowercase lettering (a., b., c.).
	BulletStyleLetterLower = "letter_lower"
	// BulletStyleLetterUpper renders uppercase lettering (A., B., C.).
	BulletStyleLetterUpper = "letter_upper"
	// BulletStyleRomanLower renders lowercase roman numerals (i., ii., iii.).
	BulletStyleRomanLower = "roman_lower"
	// BulletStyleRomanUpper renders uppercase roman numerals (I., II., III.).
	BulletStyleRomanUpper = "roman_upper"
	// BulletStyleCustom renders one caller-provided bullet character.
	BulletStyleCustom = "custom"
	// BulletStyleNone renders no bullet marker.
	BulletStyleNone = "none"
)

const maxBulletLevel = 8

func defaultTextParagraphStyle() TextParagraphStyle {
	return TextParagraphStyle{
		BulletStyle: BulletStyleBullet,
	}
}

// WithBulletStyle sets the bullet style for this paragraph.
func (p TextParagraphStyle) WithBulletStyle(style string) TextParagraphStyle {
	p.BulletStyle = normalizeBulletStyle(style)
	return p
}

// WithBullet sets standard bullet style.
func (p TextParagraphStyle) WithBullet() TextParagraphStyle {
	return p.WithBulletStyle(BulletStyleBullet)
}

// WithNumbered sets numbered list style.
func (p TextParagraphStyle) WithNumbered() TextParagraphStyle {
	return p.WithBulletStyle(BulletStyleNumber)
}

// WithLetteredLower sets lowercase lettered list style.
func (p TextParagraphStyle) WithLetteredLower() TextParagraphStyle {
	return p.WithBulletStyle(BulletStyleLetterLower)
}

// WithLetteredUpper sets uppercase lettered list style.
func (p TextParagraphStyle) WithLetteredUpper() TextParagraphStyle {
	return p.WithBulletStyle(BulletStyleLetterUpper)
}

// WithRomanLower sets lowercase roman list style.
func (p TextParagraphStyle) WithRomanLower() TextParagraphStyle {
	return p.WithBulletStyle(BulletStyleRomanLower)
}

// WithRomanUpper sets uppercase roman list style.
func (p TextParagraphStyle) WithRomanUpper() TextParagraphStyle {
	return p.WithBulletStyle(BulletStyleRomanUpper)
}

// WithCustomBullet sets custom single-character bullet style.
func (p TextParagraphStyle) WithCustomBullet(char string) TextParagraphStyle {
	p.BulletStyle = BulletStyleCustom
	p.BulletChar = strings.TrimSpace(char)
	return p
}

// WithNoBullet disables bullet marker rendering for this paragraph.
func (p TextParagraphStyle) WithNoBullet() TextParagraphStyle {
	return p.WithBulletStyle(BulletStyleNone)
}

// WithLevel sets paragraph bullet nesting level (0..8).
func (p TextParagraphStyle) WithLevel(level int) TextParagraphStyle {
	p.Level = level
	return p
}

// WithBulletStyle sets the default bullet style for subsequent AddBullet calls.
func (s SlideContent) WithBulletStyle(style string) SlideContent {
	s.DefaultBulletStyle = defaultTextParagraphStyle().WithBulletStyle(style)
	return s
}

// AddStyledBullet appends a bullet with explicit paragraph style.
func (s SlideContent) AddStyledBullet(text string, style TextParagraphStyle) SlideContent {
	return s.AddBulletWithStyle(text, style)
}

// AddNumbered appends one numbered bullet.
func (s SlideContent) AddNumbered(text string) SlideContent {
	return s.AddBulletWithStyle(text, defaultTextParagraphStyle().WithNumbered())
}

// AddLettered appends one lower-letter bullet.
func (s SlideContent) AddLettered(text string) SlideContent {
	return s.AddBulletWithStyle(text, defaultTextParagraphStyle().WithLetteredLower())
}

// AddRomanLower appends one lowercase-roman bullet.
func (s SlideContent) AddRomanLower(text string) SlideContent {
	return s.AddBulletWithStyle(text, defaultTextParagraphStyle().WithRomanLower())
}

// AddRomanUpper appends one uppercase-roman bullet.
func (s SlideContent) AddRomanUpper(text string) SlideContent {
	return s.AddBulletWithStyle(text, defaultTextParagraphStyle().WithRomanUpper())
}

// AddCustomBullet appends one custom-character bullet.
func (s SlideContent) AddCustomBullet(text string, char string) SlideContent {
	return s.AddBulletWithStyle(text, defaultTextParagraphStyle().WithCustomBullet(char))
}

// AddSubBullet appends one indented bullet at level 1.
func (s SlideContent) AddSubBullet(text string) SlideContent {
	style := s.DefaultBulletStyle
	style.Level++
	return s.AddBulletWithStyle(text, style)
}

func normalizeBulletStyle(style string) string {
	normalized := strings.ToLower(strings.TrimSpace(style))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")

	switch normalized {
	case "":
		return ""
	case BulletStyleBullet:
		return BulletStyleBullet
	case BulletStyleNumber, "numbered":
		return BulletStyleNumber
	case BulletStyleLetterLower, "lettered", "letter", "letterlower", "alphalower":
		return BulletStyleLetterLower
	case BulletStyleLetterUpper, "letterupper", "alphaupper":
		return BulletStyleLetterUpper
	case BulletStyleRomanLower, "romanlower":
		return BulletStyleRomanLower
	case BulletStyleRomanUpper, "roman", "romanupper":
		return BulletStyleRomanUpper
	case BulletStyleCustom:
		return BulletStyleCustom
	case BulletStyleNone:
		return BulletStyleNone
	default:
		return normalized
	}
}
