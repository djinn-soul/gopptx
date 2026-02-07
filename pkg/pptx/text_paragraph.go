package pptx

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

// TextParagraphStyle describes paragraph-level formatting for one bullet line.
type TextParagraphStyle struct {
	Align          string
	SpaceBeforePt  int
	SpaceAfterPt   int
	LineSpacingPct int
	BulletStyle    string
	BulletChar     string
	Level          int
}

// NewTextParagraphStyle creates one paragraph style with default settings.
func NewTextParagraphStyle() TextParagraphStyle {
	return defaultTextParagraphStyle()
}

// WithAlign sets paragraph alignment.
func (p TextParagraphStyle) WithAlign(align string) TextParagraphStyle {
	p.Align = normalizeTextAlign(align)
	return p
}

// WithAlignLeft sets paragraph alignment to left.
func (p TextParagraphStyle) WithAlignLeft() TextParagraphStyle {
	return p.WithAlign(TextAlignLeft)
}

// WithAlignCenter sets paragraph alignment to center.
func (p TextParagraphStyle) WithAlignCenter() TextParagraphStyle {
	return p.WithAlign(TextAlignCenter)
}

// WithAlignRight sets paragraph alignment to right.
func (p TextParagraphStyle) WithAlignRight() TextParagraphStyle {
	return p.WithAlign(TextAlignRight)
}

// WithAlignJustify sets paragraph alignment to justify.
func (p TextParagraphStyle) WithAlignJustify() TextParagraphStyle {
	return p.WithAlign(TextAlignJustify)
}

// WithSpaceBeforePt sets space before paragraph in points.
func (p TextParagraphStyle) WithSpaceBeforePt(points int) TextParagraphStyle {
	p.SpaceBeforePt = points
	return p
}

// WithSpaceAfterPt sets space after paragraph in points.
func (p TextParagraphStyle) WithSpaceAfterPt(points int) TextParagraphStyle {
	p.SpaceAfterPt = points
	return p
}

// WithLineSpacingPct sets line spacing percentage (100 = single spacing).
func (p TextParagraphStyle) WithLineSpacingPct(percent int) TextParagraphStyle {
	p.LineSpacingPct = percent
	return p
}

// AddBulletWithStyle appends a plain bullet plus paragraph style.
func (s SlideContent) AddBulletWithStyle(text string, style TextParagraphStyle) SlideContent {
	s = s.AddBullet(text)
	s.BulletStyles[len(s.BulletStyles)-1] = normalizeTextParagraphStyle(style)
	return s
}

// AddBulletRunsWithStyle appends rich bullet runs plus paragraph style.
func (s SlideContent) AddBulletRunsWithStyle(runs []TextRun, style TextParagraphStyle) SlideContent {
	s = s.AddBulletRuns(runs)
	s.BulletStyles[len(s.BulletStyles)-1] = normalizeTextParagraphStyle(style)
	return s
}

func normalizeTextAlign(align string) string {
	return strings.ToLower(strings.TrimSpace(align))
}

func normalizeTextParagraphStyle(style TextParagraphStyle) TextParagraphStyle {
	normalizedBulletStyle := normalizeBulletStyle(style.BulletStyle)
	if normalizedBulletStyle == "" {
		normalizedBulletStyle = BulletStyleBullet
	}
	return TextParagraphStyle{
		Align:          normalizeTextAlign(style.Align),
		SpaceBeforePt:  style.SpaceBeforePt,
		SpaceAfterPt:   style.SpaceAfterPt,
		LineSpacingPct: style.LineSpacingPct,
		BulletStyle:    normalizedBulletStyle,
		BulletChar:     strings.TrimSpace(style.BulletChar),
		Level:          style.Level,
	}
}
