package elements

import (
	"fmt"
	"strings"
)

const (
	// MaxBulletLevel defines the maximum nested bullet depth supported by PowerPoint.
	MaxBulletLevel = 9

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

	// TextAlignLeft aligns paragraph text to the left.
	TextAlignLeft = "l"
	// TextAlignCenter aligns paragraph text to the center.
	TextAlignCenter = "ctr"
	// TextAlignRight aligns paragraph text to the right.
	TextAlignRight = "r"
	// TextAlignJustify aligns paragraph text with justification.
	TextAlignJustify = "just"
)

// DefaultTextParagraphStyle returns the standard paragraph styling.
func DefaultTextParagraphStyle() TextParagraphStyle {
	return TextParagraphStyle{
		BulletStyle: BulletStyleBullet,
	}
}

// TextRun describes a single piece of text with uniform styling.
type TextRun struct {
	Text          string
	Bold          bool
	Italic        bool
	Underline     bool
	Strikethrough bool
	Subscript     bool
	Superscript   bool
	Color         string
	Highlight     string
	Font          string
	SizePt        int
	Code          bool
	Hyperlink     *Hyperlink
}

// Validate checks for invalid text run properties.
func (r TextRun) Validate() error {
	if r.SizePt < 0 {
		return fmt.Errorf("size must be >= 0")
	}
	if r.Color != "" && !IsHexColor(r.Color) {
		return fmt.Errorf("color must be 6-digit RGB hex")
	}
	if r.Highlight != "" && !IsHexColor(r.Highlight) {
		return fmt.Errorf("highlight must be 6-digit RGB hex")
	}
	if r.Subscript && r.Superscript {
		return fmt.Errorf("cannot be both subscript and superscript")
	}
	if r.Hyperlink != nil {
		if err := r.Hyperlink.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// NewTextRun creates a simple text run.
func NewTextRun(text string) TextRun {
	return TextRun{Text: text}
}

// WithBold sets bold property.
func (r TextRun) WithBold(bold bool) TextRun {
	r.Bold = bold
	return r
}

// WithItalic sets italic property.
func (r TextRun) WithItalic(italic bool) TextRun {
	r.Italic = italic
	return r
}

// WithUnderline sets underline property.
func (r TextRun) WithUnderline(underline bool) TextRun {
	r.Underline = underline
	return r
}

// WithStrikethrough sets strikethrough property.
func (r TextRun) WithStrikethrough(strikethrough bool) TextRun {
	r.Strikethrough = strikethrough
	return r
}

// WithSubscript sets subscript property.
func (r TextRun) WithSubscript(subscript bool) TextRun {
	r.Subscript = subscript
	if subscript {
		r.Superscript = false
	}
	return r
}

// WithSuperscript sets superscript property.
func (r TextRun) WithSuperscript(superscript bool) TextRun {
	r.Superscript = superscript
	if superscript {
		r.Subscript = false
	}
	return r
}

// WithColor sets hex color.
func (r TextRun) WithColor(color string) TextRun {
	r.Color = NormalizeHexColor(color)
	return r
}

// WithHighlight sets highlight color.
func (r TextRun) WithHighlight(color string) TextRun {
	r.Highlight = NormalizeHexColor(color)
	return r
}

// WithFont sets font name.
func (r TextRun) WithFont(font string) TextRun {
	r.Font = strings.TrimSpace(font)
	return r
}

// WithSizePt sets font size in points.
func (r TextRun) WithSizePt(size int) TextRun {
	r.SizePt = size
	return r
}

// WithCode sets code format (monospaced).
func (r TextRun) WithCode(code bool) TextRun {
	r.Code = code
	return r
}

// WithHyperlink sets a hyperlink for the run.
func (r TextRun) WithHyperlink(link Hyperlink) TextRun {
	r.Hyperlink = &link
	return r
}

// TextParagraphStyle describes paragraph-level formatting for one bullet line.
type TextParagraphStyle struct {
	Align          string
	SpaceBeforePt  int
	SpaceAfterPt   int
	LineSpacingPct int
	BulletStyle    string
	BulletChar     string
	BulletColor    string
	BulletSize     int
	Level          int
}

// Validate checks for invalid text paragraph style properties.
func (p TextParagraphStyle) Validate() error {
	if p.BulletColor != "" && !IsHexColor(p.BulletColor) {
		return fmt.Errorf("bullet color must be hex")
	}
	if p.SpaceBeforePt < 0 {
		return fmt.Errorf("space-before must be >= 0")
	}
	if p.SpaceAfterPt < 0 {
		return fmt.Errorf("space-after must be >= 0")
	}
	if p.LineSpacingPct < 0 {
		return fmt.Errorf("line-spacing must be >= 0")
	}
	switch p.Align {
	case "", TextAlignLeft, TextAlignCenter, TextAlignRight, TextAlignJustify:
		// Valid
	default:
		return fmt.Errorf("align must be one of l|ctr|r|just")
	}
	return nil
}

// WithAlign sets paragraph alignment.
func (p TextParagraphStyle) WithAlign(align string) TextParagraphStyle {
	p.Align = NormalizeTextAlign(align)
	return p
}

// WithAlignLeft sets left alignment.
func (p TextParagraphStyle) WithAlignLeft() TextParagraphStyle {
	p.Align = TextAlignLeft
	return p
}

// WithAlignCenter sets center alignment.
func (p TextParagraphStyle) WithAlignCenter() TextParagraphStyle {
	p.Align = TextAlignCenter
	return p
}

// WithAlignRight sets right alignment.
func (p TextParagraphStyle) WithAlignRight() TextParagraphStyle {
	p.Align = TextAlignRight
	return p
}

// WithAlignJustify sets justified alignment.
func (p TextParagraphStyle) WithAlignJustify() TextParagraphStyle {
	p.Align = TextAlignJustify
	return p
}

// WithNumbered sets the bullet style to numbered.
func (p TextParagraphStyle) WithNumbered() TextParagraphStyle {
	p.BulletStyle = BulletStyleNumber
	return p
}

// WithBulletStyle sets the bullet style by name.
func (p TextParagraphStyle) WithBulletStyle(style string) TextParagraphStyle {
	p.BulletStyle = NormalizeBulletStyle(style)
	return p
}

// WithLetteredLower sets lowercase lettered list style.
func (p TextParagraphStyle) WithLetteredLower() TextParagraphStyle {
	p.BulletStyle = BulletStyleLetterLower
	return p
}

// WithRomanUpper sets uppercase roman list style.
func (p TextParagraphStyle) WithRomanUpper() TextParagraphStyle {
	p.BulletStyle = BulletStyleRomanUpper
	return p
}

// WithCustomBullet sets custom single-character bullet style.
func (p TextParagraphStyle) WithCustomBullet(char string) TextParagraphStyle {
	p.BulletStyle = BulletStyleCustom
	p.BulletChar = strings.TrimSpace(char)
	return p
}

// WithNoBullet sets the bullet style to none.
func (p TextParagraphStyle) WithNoBullet() TextParagraphStyle {
	p.BulletStyle = BulletStyleNone
	return p
}

// WithLevel sets paragraph bullet nesting level (0..8).
func (p TextParagraphStyle) WithLevel(level int) TextParagraphStyle {
	p.Level = level
	return p
}

// WithBulletChar sets custom bullet character.
func (p TextParagraphStyle) WithBulletChar(char string) TextParagraphStyle {
	p.BulletChar = strings.TrimSpace(char)
	return p
}

// WithBulletColor sets hex color for bullet.
func (p TextParagraphStyle) WithBulletColor(color string) TextParagraphStyle {
	p.BulletColor = NormalizeHexColor(color)
	return p
}

// WithBulletSize sets bullet size as percentage of text size.
func (p TextParagraphStyle) WithBulletSize(size int) TextParagraphStyle {
	p.BulletSize = size
	return p
}

// WithSpaceBeforePt sets space before paragraph in points.
func (p TextParagraphStyle) WithSpaceBeforePt(pt int) TextParagraphStyle {
	p.SpaceBeforePt = pt
	return p
}

// WithSpaceAfterPt sets space after paragraph in points.
func (p TextParagraphStyle) WithSpaceAfterPt(pt int) TextParagraphStyle {
	p.SpaceAfterPt = pt
	return p
}

// WithLineSpacingPct sets line spacing as percentage (e.g. 100).
func (p TextParagraphStyle) WithLineSpacingPct(pct int) TextParagraphStyle {
	p.LineSpacingPct = pct
	return p
}

// NormalizeTextAlign sanitizes alignment strings.
func NormalizeTextAlign(align string) string {
	return strings.ToLower(strings.TrimSpace(align))
}

// NormalizeBulletStyle sanitizes bullet style strings.
func NormalizeBulletStyle(style string) string {
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

// NormalizeTextParagraphStyle ensures all fields are within expected bounds.
func NormalizeTextParagraphStyle(style TextParagraphStyle) TextParagraphStyle {
	normalizedBulletStyle := NormalizeBulletStyle(style.BulletStyle)
	if normalizedBulletStyle == "" {
		normalizedBulletStyle = BulletStyleBullet
	}
	return TextParagraphStyle{
		Align:          NormalizeTextAlign(style.Align),
		SpaceBeforePt:  style.SpaceBeforePt,
		SpaceAfterPt:   style.SpaceAfterPt,
		LineSpacingPct: style.LineSpacingPct,
		BulletStyle:    normalizedBulletStyle,
		BulletChar:     strings.TrimSpace(style.BulletChar),
		Level:          style.Level,
	}
}

// NormalizeTextRuns removes empty runs and merges adjacent runs with identical styling.
func NormalizeTextRuns(runs []TextRun) []TextRun {
	if len(runs) == 0 {
		return nil
	}
	result := make([]TextRun, 0, len(runs))
	for _, run := range runs {
		if run.Text == "" {
			continue
		}
		if len(result) > 0 {
			last := &result[len(result)-1]
			if last.Bold == run.Bold && last.Italic == run.Italic && last.Code == run.Code &&
				last.Color == run.Color && last.SizePt == run.SizePt && last.Underline == run.Underline &&
				last.Strikethrough == run.Strikethrough && last.Subscript == run.Subscript &&
				last.Superscript == run.Superscript && last.Highlight == run.Highlight &&
				last.Font == run.Font && last.Hyperlink == run.Hyperlink {
				last.Text += run.Text
				continue
			}
		}
		result = append(result, run)
	}
	return result
}
