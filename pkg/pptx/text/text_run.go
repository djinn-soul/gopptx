package text

import (
	"errors"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

const (
	UnderlineStyleNone   = "none"
	UnderlineStyleSingle = "sng"
	UnderlineStyleDouble = "dbl"
	UnderlineStyleDotted = "dotted"

	StrikethroughStyleNone   = "none"
	StrikethroughStyleSingle = "sngStrike"
	StrikethroughStyleDouble = "dblStrike"
)

// Run describes a single piece of text with uniform styling.
type Run struct {
	Text           string
	Bold           bool
	Italic         bool
	Underline      string // "none", "sng", "dbl", "dotted", etc.
	Strikethrough  string // "none", "sng", "dbl"
	Subscript      bool
	Superscript    bool
	Color          string
	Highlight      string
	Font           string
	SizePt         int
	Code           bool
	AllCaps        bool
	SmallCaps      bool
	OutlineColor   string            // Character stroke/outline color hex
	OutlineWidthPt float64           // Character stroke/outline width in points (default 1pt when OutlineColor is set)
	Hyperlink      *action.Hyperlink // Click behavior
	HoverAction    *action.Hyperlink // Hover behavior
}

// NewRun creates a simple text run.
func NewRun(text string) Run {
	return Run{Text: text}
}

// WithBold sets bold property.
func (r Run) WithBold(bold bool) Run {
	r.Bold = bold
	return r
}

// WithItalic sets italic property.
func (r Run) WithItalic(italic bool) Run {
	r.Italic = italic
	return r
}

// WithUnderline sets underline property (boolean compatibility).
func (r Run) WithUnderline(underline bool) Run {
	if underline {
		r.Underline = UnderlineStyleSingle
	} else {
		r.Underline = UnderlineStyleNone
	}
	return r
}

// WithUnderlineStyle sets a specific underline style.
func (r Run) WithUnderlineStyle(style string) Run {
	r.Underline = style
	return r
}

// WithStrikethrough sets strikethrough property (boolean compatibility).
func (r Run) WithStrikethrough(strikethrough bool) Run {
	if strikethrough {
		r.Strikethrough = StrikethroughStyleSingle
	} else {
		r.Strikethrough = StrikethroughStyleNone
	}
	return r
}

// WithStrikethroughStyle sets a specific strikethrough style.
func (r Run) WithStrikethroughStyle(style string) Run {
	r.Strikethrough = style
	return r
}

// WithSubscript sets subscript property.
func (r Run) WithSubscript(subscript bool) Run {
	r.Subscript = subscript
	if subscript {
		r.Superscript = false
	}
	return r
}

// WithSuperscript sets superscript property.
func (r Run) WithSuperscript(superscript bool) Run {
	r.Superscript = superscript
	if superscript {
		r.Subscript = false
	}
	return r
}

// WithColor sets hex color.
func (r Run) WithColor(color string) Run {
	r.Color = common.NormalizeHexColor(color)
	return r
}

// WithHighlight sets highlight color.
func (r Run) WithHighlight(color string) Run {
	r.Highlight = common.NormalizeHexColor(color)
	return r
}

// WithFont sets font name.
func (r Run) WithFont(font string) Run {
	r.Font = strings.TrimSpace(font)
	return r
}

// WithSizePt sets font size in points.
func (r Run) WithSizePt(size int) Run {
	r.SizePt = size
	return r
}

// WithCode sets code format (monospaced).
func (r Run) WithCode(code bool) Run {
	r.Code = code
	return r
}

// WithAllCaps sets the text to be all uppercase.
func (r Run) WithAllCaps(allCaps bool) Run {
	r.AllCaps = allCaps
	if allCaps {
		r.SmallCaps = false
	}
	return r
}

// WithSmallCaps sets the text to be small caps.
func (r Run) WithSmallCaps(smallCaps bool) Run {
	r.SmallCaps = smallCaps
	if smallCaps {
		r.AllCaps = false
	}
	return r
}

// WithOutline sets a character stroke/outline using the given hex color.
// An optional width in points can be provided; defaults to 1pt if omitted.
func (r Run) WithOutline(color string, widthPt ...float64) Run {
	r.OutlineColor = common.NormalizeHexColor(color)
	if len(widthPt) > 0 && widthPt[0] > 0 {
		r.OutlineWidthPt = widthPt[0]
	}
	return r
}

// WithHyperlink sets a click action for the run.
func (r Run) WithHyperlink(link action.Hyperlink) Run {
	r.Hyperlink = &link
	return r
}

// WithHoverAction sets a hover action for the run.
func (r Run) WithHoverAction(link action.Hyperlink) Run {
	r.HoverAction = &link
	return r
}

// Validate checks for invalid text run properties.
func (r Run) Validate() error {
	if r.SizePt < 0 {
		return errors.New("size must be >= 0")
	}
	if r.Color != "" && !common.IsHexColor(r.Color) {
		return errors.New("color must be 6-digit RGB hex")
	}
	if r.Highlight != "" && !common.IsHexColor(r.Highlight) {
		return errors.New("highlight must be 6-digit RGB hex")
	}
	if r.Subscript && r.Superscript {
		return errors.New("cannot be both subscript and superscript")
	}
	if r.Hyperlink != nil {
		if err := r.Hyperlink.Validate(); err != nil {
			return err
		}
	}
	return nil
}
