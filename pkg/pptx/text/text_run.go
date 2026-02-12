package text

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

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
	AllCaps       bool
	SmallCaps     bool
	Hyperlink     *action.Hyperlink // Click behavior
	HoverAction   *action.Hyperlink // Hover behavior
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
	r.Color = common.NormalizeHexColor(color)
	return r
}

// WithHighlight sets highlight color.
func (r TextRun) WithHighlight(color string) TextRun {
	r.Highlight = common.NormalizeHexColor(color)
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

// WithAllCaps sets the text to be all uppercase.
func (r TextRun) WithAllCaps(allCaps bool) TextRun {
	r.AllCaps = allCaps
	if allCaps {
		r.SmallCaps = false
	}
	return r
}

// WithSmallCaps sets the text to be small caps.
func (r TextRun) WithSmallCaps(smallCaps bool) TextRun {
	r.SmallCaps = smallCaps
	if smallCaps {
		r.AllCaps = false
	}
	return r
}

// WithHyperlink sets a click action for the run.
func (r TextRun) WithHyperlink(link action.Hyperlink) TextRun {
	r.Hyperlink = &link
	return r
}

// WithHoverAction sets a hover action for the run.
func (r TextRun) WithHoverAction(link action.Hyperlink) TextRun {
	r.HoverAction = &link
	return r
}

// Validate checks for invalid text run properties.
func (r TextRun) Validate() error {
	if r.SizePt < 0 {
		return fmt.Errorf("size must be >= 0")
	}
	if r.Color != "" && !common.IsHexColor(r.Color) {
		return fmt.Errorf("color must be 6-digit RGB hex")
	}
	if r.Highlight != "" && !common.IsHexColor(r.Highlight) {
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
