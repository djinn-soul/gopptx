package text

import (
	"errors"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Paragraph represents a single paragraph of text with runs and styling.
type Paragraph struct {
	Runs  []Run
	Style ParagraphStyle
}

// NewParagraph creates a new paragraph with default style.
func NewParagraph() Paragraph {
	return Paragraph{
		Runs:  make([]Run, 0),
		Style: DefaultParagraphStyle(),
	}
}

// AddRun appends a text run to the paragraph.
func (p Paragraph) AddRun(run Run) Paragraph {
	p.Runs = append(p.Runs, run)
	return p
}

// WithStyle sets the paragraph style.
func (p Paragraph) WithStyle(style ParagraphStyle) Paragraph {
	p.Style = style
	return p
}

// ParagraphStyle describes paragraph-level formatting for one bullet line.
type ParagraphStyle struct {
	Align          string
	SpaceBeforePt  int
	SpaceAfterPt   int
	LineSpacingPct int
	LineSpacingPts int
	BulletStyle    string
	BulletChar     string
	BulletColor    string
	BulletSize     int
	TabStops       []styling.Length // EMU
	Level          int
	LeftIndent     styling.Length // EMU
	RightIndent    styling.Length // EMU
	HangingIndent  styling.Length // EMU
	RTL            bool
}

// DefaultParagraphStyle returns the standard paragraph styling.
func DefaultParagraphStyle() ParagraphStyle {
	return ParagraphStyle{
		BulletStyle: BulletStyleBullet,
	}
}

// NewParagraphStyle creates one paragraph style with default settings.
func NewParagraphStyle() ParagraphStyle {
	return NormalizeParagraphStyle(ParagraphStyle{})
}

// WithAlign sets paragraph alignment.
func (p ParagraphStyle) WithAlign(align string) ParagraphStyle {
	p.Align = NormalizeTextAlign(align)
	return p
}

// WithAlignLeft sets left alignment.
func (p ParagraphStyle) WithAlignLeft() ParagraphStyle {
	p.Align = TextAlignLeft
	return p
}

// WithAlignCenter sets center alignment.
func (p ParagraphStyle) WithAlignCenter() ParagraphStyle {
	p.Align = TextAlignCenter
	return p
}

// WithAlignRight sets right alignment.
func (p ParagraphStyle) WithAlignRight() ParagraphStyle {
	p.Align = TextAlignRight
	return p
}

// WithAlignJustify sets justified alignment.
func (p ParagraphStyle) WithAlignJustify() ParagraphStyle {
	p.Align = TextAlignJustify
	return p
}

// WithNumbered sets the bullet style to numbered.
func (p ParagraphStyle) WithNumbered() ParagraphStyle {
	p.BulletStyle = BulletStyleNumber
	return p
}

// WithBulletStyle sets the bullet style by name.
func (p ParagraphStyle) WithBulletStyle(style string) ParagraphStyle {
	p.BulletStyle = NormalizeBulletStyle(style)
	return p
}

// WithLetteredLower sets lowercase lettered list style.
func (p ParagraphStyle) WithLetteredLower() ParagraphStyle {
	p.BulletStyle = BulletStyleLetterLower
	return p
}

// WithRomanUpper sets uppercase roman list style.
func (p ParagraphStyle) WithRomanUpper() ParagraphStyle {
	p.BulletStyle = BulletStyleRomanUpper
	return p
}

// WithCustomBullet sets custom single-character bullet style.
func (p ParagraphStyle) WithCustomBullet(char string) ParagraphStyle {
	p.BulletStyle = BulletStyleCustom
	p.BulletChar = strings.TrimSpace(char)
	return p
}

// WithNoBullet sets the bullet style to none.
func (p ParagraphStyle) WithNoBullet() ParagraphStyle {
	p.BulletStyle = BulletStyleNone
	return p
}

// WithLevel sets paragraph bullet nesting level (0..8).
func (p ParagraphStyle) WithLevel(level int) ParagraphStyle {
	p.Level = level
	return p
}

// WithBulletChar sets custom bullet character.
func (p ParagraphStyle) WithBulletChar(char string) ParagraphStyle {
	p.BulletChar = strings.TrimSpace(char)
	return p
}

// WithBulletColor sets hex color for bullet.
func (p ParagraphStyle) WithBulletColor(color string) ParagraphStyle {
	p.BulletColor = common.NormalizeHexColor(color)
	return p
}

// WithBulletSize sets bullet size as percentage of text size.
func (p ParagraphStyle) WithBulletSize(size int) ParagraphStyle {
	p.BulletSize = size
	return p
}

// WithSpaceBeforePt sets space before paragraph in points.
func (p ParagraphStyle) WithSpaceBeforePt(pt int) ParagraphStyle {
	p.SpaceBeforePt = pt
	return p
}

// WithSpaceAfterPt sets space after paragraph in points.
func (p ParagraphStyle) WithSpaceAfterPt(pt int) ParagraphStyle {
	p.SpaceAfterPt = pt
	return p
}

// WithLineSpacingPct sets line spacing as percentage (e.g. 100).
func (p ParagraphStyle) WithLineSpacingPct(pct int) ParagraphStyle {
	p.LineSpacingPct = pct
	return p
}

// WithLineSpacingPts sets line spacing as points (e.g. 18).
func (p ParagraphStyle) WithLineSpacingPts(pt int) ParagraphStyle {
	p.LineSpacingPts = pt
	return p
}

// WithTabStops sets paragraph tab stops in EMU.
func (p ParagraphStyle) WithTabStops(stops ...styling.Length) ParagraphStyle {
	if len(stops) == 0 {
		p.TabStops = nil
		return p
	}
	p.TabStops = append(make([]styling.Length, 0, len(stops)), stops...)
	return p
}

// WithLeftIndent sets the left margin for the paragraph in EMUs.
func (p ParagraphStyle) WithLeftIndent(emu styling.Length) ParagraphStyle {
	p.LeftIndent = emu
	return p
}

// WithRightIndent sets the right margin for the paragraph in EMUs.
func (p ParagraphStyle) WithRightIndent(emu styling.Length) ParagraphStyle {
	p.RightIndent = emu
	return p
}

// WithRTL sets the Right-To-Left flag for this paragraph.
func (p ParagraphStyle) WithRTL(rtl bool) ParagraphStyle {
	p.RTL = rtl
	return p
}

// WithHangingIndent sets the hanging indent (indent of the first line) in EMUs.
// Note: Usually negative to shift the first line to the left of the rest.
func (p ParagraphStyle) WithHangingIndent(emu styling.Length) ParagraphStyle {
	p.HangingIndent = emu
	return p
}

// Validate checks for invalid text paragraph style properties.
func (p ParagraphStyle) Validate() error {
	if p.BulletColor != "" && !common.IsHexColor(p.BulletColor) {
		return errors.New("bullet color must be hex")
	}
	if p.SpaceBeforePt < 0 {
		return errors.New("space-before must be >= 0")
	}
	if p.SpaceAfterPt < 0 {
		return errors.New("space-after must be >= 0")
	}
	if p.LineSpacingPct < 0 {
		return errors.New("line-spacing must be >= 0")
	}
	if p.LineSpacingPts < 0 {
		return errors.New("line-spacing points must be >= 0")
	}
	if p.LineSpacingPct > 0 && p.LineSpacingPts > 0 {
		return errors.New("line-spacing percent and points are mutually exclusive")
	}
	for _, tabStop := range p.TabStops {
		if tabStop.Emu() < 0 {
			return errors.New("tab-stops must be >= 0")
		}
	}
	switch p.Align {
	case "", TextAlignLeft, TextAlignCenter, TextAlignRight, TextAlignJustify:
		// Valid
	default:
		return errors.New("align must be one of l|ctr|r|just")
	}
	return nil
}
