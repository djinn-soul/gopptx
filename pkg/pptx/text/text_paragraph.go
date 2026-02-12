package text

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// TextParagraph represents a single paragraph of text with runs and styling.
type TextParagraph struct {
	Runs  []TextRun
	Style TextParagraphStyle
}

// NewTextParagraph creates a new paragraph with default style.
func NewTextParagraph() TextParagraph {
	return TextParagraph{
		Runs:  make([]TextRun, 0),
		Style: DefaultTextParagraphStyle(),
	}
}

// AddRun appends a text run to the paragraph.
func (p TextParagraph) AddRun(run TextRun) TextParagraph {
	p.Runs = append(p.Runs, run)
	return p
}

// WithStyle sets the paragraph style.
func (p TextParagraph) WithStyle(style TextParagraphStyle) TextParagraph {
	p.Style = style
	return p
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
	LeftIndent     styling.Length // EMU
	RightIndent    styling.Length // EMU
	HangingIndent  styling.Length // EMU
}

// DefaultTextParagraphStyle returns the standard paragraph styling.
func DefaultTextParagraphStyle() TextParagraphStyle {
	return TextParagraphStyle{
		BulletStyle: BulletStyleBullet,
	}
}

// NewTextParagraphStyle creates one paragraph style with default settings.
func NewTextParagraphStyle() TextParagraphStyle {
	return NormalizeTextParagraphStyle(TextParagraphStyle{})
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
	p.BulletColor = common.NormalizeHexColor(color)
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

// WithLeftIndent sets the left margin for the paragraph in EMUs.
func (p TextParagraphStyle) WithLeftIndent(emu styling.Length) TextParagraphStyle {
	p.LeftIndent = emu
	return p
}

// WithRightIndent sets the right margin for the paragraph in EMUs.
func (p TextParagraphStyle) WithRightIndent(emu styling.Length) TextParagraphStyle {
	p.RightIndent = emu
	return p
}

// WithHangingIndent sets the hanging indent (indent of the first line) in EMUs.
// Note: Usually negative to shift the first line to the left of the rest.
func (p TextParagraphStyle) WithHangingIndent(emu styling.Length) TextParagraphStyle {
	p.HangingIndent = emu
	return p
}

// Validate checks for invalid text paragraph style properties.
func (p TextParagraphStyle) Validate() error {
	if p.BulletColor != "" && !common.IsHexColor(p.BulletColor) {
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
