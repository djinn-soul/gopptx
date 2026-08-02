package fonts

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/signintech/gopdf/fontmaker/core"
)

// Bounds on the font sizes the fitter will consider, in points.
const (
	// MinFitFontSizePt is the smallest size fit_text will shrink to. Below this
	// the text is unreadable, so reporting the floor beats returning something
	// nobody can read.
	MinFitFontSizePt = 6
	// DefaultMaxFitFontSizePt is the starting size when the caller names none.
	DefaultMaxFitFontSizePt = 18
)

// pointsPerEmHeightRatio approximates the line height of a font as a multiple
// of its point size. PowerPoint's own single line spacing is about 1.2 em for
// the fonts it ships.
const pointsPerEmHeightRatio = 1.2

// ErrNoFontMetrics reports that text cannot be measured because no usable font
// file was supplied. Callers that only want the autofit flags should not treat
// this as fatal.
var ErrNoFontMetrics = errors.New("no font metrics available to measure text")

// TextMeasurer measures strings against one font's real glyph advances.
type TextMeasurer struct {
	// unitsPerEm scales glyph advances to em units.
	unitsPerEm float64
	// advances maps a rune to its advance width in font units.
	advances map[rune]float64
	// fallback is the advance used for a rune the font has no glyph for.
	fallback float64
}

// NewTextMeasurer reads glyph advances out of a TrueType font file.
func NewTextMeasurer(fontPath string) (*TextMeasurer, error) {
	if strings.TrimSpace(fontPath) == "" {
		return nil, ErrNoFontMetrics
	}
	parser := core.TTFParser{}
	if err := parser.Parse(fontPath); err != nil {
		return nil, fmt.Errorf("parse font %s: %w", fontPath, err)
	}

	widths := parser.Widths()
	unitsPerEm := float64(parser.UnitsPerEm())
	if unitsPerEm <= 0 || len(widths) == 0 {
		return nil, fmt.Errorf("font %s has no usable metrics: %w", fontPath, ErrNoFontMetrics)
	}

	advances := make(map[rune]float64, len(parser.Chars()))
	for code, glyphIndex := range parser.Chars() {
		if int(glyphIndex) >= len(widths) {
			continue
		}
		advances[rune(code)] = float64(widths[glyphIndex])
	}

	// An unmapped rune is measured as the font's own space, or half an em when
	// even that is missing, rather than silently counting as zero width.
	fallback := unitsPerEm / 2 //nolint:mnd // half an em is the conventional default advance
	if space, ok := advances[' ']; ok && space > 0 {
		fallback = space
	}
	return &TextMeasurer{unitsPerEm: unitsPerEm, advances: advances, fallback: fallback}, nil
}

// WidthPt returns the rendered width of text at the given point size.
func (m *TextMeasurer) WidthPt(text string, fontSizePt float64) float64 {
	total := 0.0
	for _, r := range text {
		advance, ok := m.advances[r]
		if !ok {
			advance = m.fallback
		}
		total += advance
	}
	return total / m.unitsPerEm * fontSizePt
}

// FitRequest describes the box the text has to fit inside.
type FitRequest struct {
	Text string
	// WidthPt and HeightPt are the usable text area, i.e. the shape extent
	// minus its insets.
	WidthPt  float64
	HeightPt float64
	// MaxSizePt is the largest size to try; zero means DefaultMaxFitFontSizePt.
	MaxSizePt float64
	// MinSizePt is the smallest size to try; zero means MinFitFontSizePt.
	MinSizePt float64
}

// FitResult reports the chosen size and whether the text actually fit.
type FitResult struct {
	// FontSizePt is the largest size that fits, or MinSizePt when nothing does.
	FontSizePt float64
	// Fits is false when even MinSizePt overflows the box.
	Fits bool
	// LineCount is how many lines the text wraps into at FontSizePt.
	LineCount int
}

// FitText picks the largest whole point size at which text wraps inside the box.
//
// Unlike the python-pptx implementation this issue is about, a word longer than
// the line is not an error: it is placed on a line of its own and allowed to
// overflow, which is what PowerPoint does. That is the crash in #1026, where an
// unbreakable run like "580.9m↓-11.3m592.2m" made the line breaker return None.
func (m *TextMeasurer) FitText(request FitRequest) (FitResult, error) {
	if m == nil {
		return FitResult{}, ErrNoFontMetrics
	}
	if request.WidthPt <= 0 || request.HeightPt <= 0 {
		return FitResult{}, errors.New("fit area must have positive width and height")
	}

	maxSize := request.MaxSizePt
	if maxSize <= 0 {
		maxSize = DefaultMaxFitFontSizePt
	}
	minSize := request.MinSizePt
	if minSize <= 0 {
		minSize = MinFitFontSizePt
	}
	if minSize > maxSize {
		return FitResult{}, fmt.Errorf("min size %v exceeds max size %v", minSize, maxSize)
	}

	for size := maxSize; size >= minSize; size-- {
		lines := m.WrapLines(request.Text, request.WidthPt, size)
		if float64(len(lines))*size*pointsPerEmHeightRatio > request.HeightPt {
			continue
		}
		// Height alone is not enough: a word with nothing to break on stays on
		// one line and can still be wider than the box, which is exactly the
		// text these issues are about.
		if m.widestLinePt(lines, size) > request.WidthPt {
			continue
		}
		return FitResult{FontSizePt: size, Fits: true, LineCount: len(lines)}, nil
	}

	lines := m.WrapLines(request.Text, request.WidthPt, minSize)
	return FitResult{FontSizePt: minSize, Fits: false, LineCount: len(lines)}, nil
}

// WrappedHeightPt reports how tall text is once wrapped to widthPt at
// fontSizePt, along with the line count. It is the measurement behind growing a
// shape to its text instead of shrinking the text to the shape.
func (m *TextMeasurer) WrappedHeightPt(
	text string,
	widthPt, fontSizePt float64,
) (heightPt float64, lineCount int, err error) {
	if m == nil {
		return 0, 0, ErrNoFontMetrics
	}
	if widthPt <= 0 {
		return 0, 0, errors.New("fit area must have positive width")
	}
	if fontSizePt <= 0 {
		return 0, 0, errors.New("font size must be positive")
	}
	lines := m.WrapLines(text, widthPt, fontSizePt)
	return float64(len(lines)) * fontSizePt * pointsPerEmHeightRatio, len(lines), nil
}

// widestLinePt returns the width of the longest wrapped line.
func (m *TextMeasurer) widestLinePt(lines []string, fontSizePt float64) float64 {
	widest := 0.0
	for _, line := range lines {
		if width := m.WidthPt(line, fontSizePt); width > widest {
			widest = width
		}
	}
	return widest
}

// WrapLines breaks text into lines that fit the given width at the given size.
// A single word wider than the line gets its own line rather than being dropped.
func (m *TextMeasurer) WrapLines(text string, widthPt, fontSizePt float64) []string {
	var lines []string
	for _, paragraph := range strings.Split(text, "\n") {
		lines = append(lines, m.wrapParagraph(paragraph, widthPt, fontSizePt)...)
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func (m *TextMeasurer) wrapParagraph(paragraph string, widthPt, fontSizePt float64) []string {
	words := strings.FieldsFunc(paragraph, unicode.IsSpace)
	if len(words) == 0 {
		return []string{""}
	}

	lines := make([]string, 0, len(words))
	current := words[0]
	for _, word := range words[1:] {
		candidate := current + " " + word
		if m.WidthPt(candidate, fontSizePt) <= widthPt {
			current = candidate
			continue
		}
		lines = append(lines, current)
		current = word
	}
	return append(lines, current)
}
