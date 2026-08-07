package elements

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// CodeBlock is a positioned, sized block of source code. Code could previously
// only reach a slide through the Markdown importer, which appends it as bullets
// in the body placeholder — so there was no way to put a code listing at chosen
// coordinates, or to have one anywhere but the body.
type CodeBlock struct {
	// Code is the source text. Tabs are kept as written.
	Code string
	// Language selects the lexer; empty or "text" renders unhighlighted.
	Language string

	X  styling.Length
	Y  styling.Length
	CX styling.Length
	CY styling.Length

	// FontSizePt overrides the default monospace size.
	FontSizePt int
	// Font overrides the default monospace typeface.
	Font string
	// BackgroundColor is the hex fill behind the code. Empty uses the dark
	// background the highlighter's palette is designed against.
	BackgroundColor string
	// ShowLanguageLabel prints the language above the listing.
	ShowLanguageLabel bool
	// AltText describes the listing for screen readers.
	AltText string
}

// Code block defaults, sized for a listing that stays readable on a slide.
const (
	DefaultCodeFontSizePt = 12
	DefaultCodeFont       = "Consolas"
	// DefaultCodeBackground is the Solarized Dark background the highlighter's
	// token colours are chosen against.
	DefaultCodeBackground = "002B36"
)

// NewCodeBlock creates a code block at the given position and size.
func NewCodeBlock(code, language string, x, y, cx, cy styling.Length) CodeBlock {
	return CodeBlock{Code: code, Language: language, X: x, Y: y, CX: cx, CY: cy}
}

// Position moves the block.
func (c CodeBlock) Position(x, y styling.Length) CodeBlock {
	c.X, c.Y = x, y
	return c
}

// Size resizes the block.
func (c CodeBlock) Size(cx, cy styling.Length) CodeBlock {
	c.CX, c.CY = cx, cy
	return c
}

// WithFontSize pins the monospace size in points.
func (c CodeBlock) WithFontSize(pt int) CodeBlock {
	c.FontSizePt = pt
	return c
}

// WithFont sets the monospace typeface.
func (c CodeBlock) WithFont(font string) CodeBlock {
	c.Font = font
	return c
}

// WithBackground sets the fill behind the code.
func (c CodeBlock) WithBackground(hexColor string) CodeBlock {
	c.BackgroundColor = hexColor
	return c
}

// WithLanguageLabel prints the language above the listing.
func (c CodeBlock) WithLanguageLabel(show bool) CodeBlock {
	c.ShowLanguageLabel = show
	return c
}

// WithAltText describes the listing for screen readers.
func (c CodeBlock) WithAltText(text string) CodeBlock {
	c.AltText = text
	return c
}

// FontSizeOrDefault is the size to render at.
func (c CodeBlock) FontSizeOrDefault() int {
	if c.FontSizePt > 0 {
		return c.FontSizePt
	}
	return DefaultCodeFontSizePt
}

// FontOrDefault is the typeface to render with.
func (c CodeBlock) FontOrDefault() string {
	if c.Font != "" {
		return c.Font
	}
	return DefaultCodeFont
}

// BackgroundOrDefault is the fill to draw behind the code.
func (c CodeBlock) BackgroundOrDefault() string {
	if c.BackgroundColor != "" {
		return c.BackgroundColor
	}
	return DefaultCodeBackground
}

// AddCodeBlock appends a code block to the slide.
func (s SlideContent) AddCodeBlock(block CodeBlock) SlideContent {
	s.CodeBlocks = append(s.CodeBlocks, block)
	return s
}
