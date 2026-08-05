// Package printsettings models the print configuration PowerPoint keeps in
// `ppt/presProps.xml` as the `<p:prnPr>` element: what to print (slides,
// handouts, notes, outline), the color mode, and the page options in the
// Print dialog.
//
// The handout header/footer/date/page-number toggles live on the handout
// master instead, see package handout.
package printsettings

import "strings"

// PrintWhat selects which view PowerPoint prints.
type PrintWhat int

const (
	// PrintSlides prints one full slide per page (the default).
	PrintSlides PrintWhat = iota
	// PrintHandouts1 prints handouts with 1 slide per page.
	PrintHandouts1
	// PrintHandouts2 prints handouts with 2 slides per page.
	PrintHandouts2
	// PrintHandouts3 prints handouts with 3 slides per page.
	PrintHandouts3
	// PrintHandouts4 prints handouts with 4 slides per page.
	PrintHandouts4
	// PrintHandouts6 prints handouts with 6 slides per page.
	PrintHandouts6
	// PrintHandouts9 prints handouts with 9 slides per page.
	PrintHandouts9
	// PrintNotes prints the notes pages.
	PrintNotes
	// PrintOutline prints the text-only outline view.
	PrintOutline
)

// The slide counts PowerPoint offers for handout pages. They are the values
// carried by PrintHandoutsN, not arbitrary numbers.
const (
	slidesPerHandout1 = 1
	slidesPerHandout2 = 2
	slidesPerHandout3 = 3
	slidesPerHandout4 = 4
	slidesPerHandout6 = 6
	slidesPerHandout9 = 9
)

// XMLValue returns the ST_PrintWhat token for the prnWhat attribute.
func (p PrintWhat) XMLValue() string {
	switch p {
	case PrintHandouts1:
		return "handouts1"
	case PrintHandouts2:
		return "handouts2"
	case PrintHandouts3:
		return "handouts3"
	case PrintHandouts4:
		return "handouts4"
	case PrintHandouts6:
		return "handouts6"
	case PrintHandouts9:
		return "handouts9"
	case PrintNotes:
		return "notes"
	case PrintOutline:
		return "outline"
	case PrintSlides:
		return "slides"
	default:
		return "slides"
	}
}

// SlidesPerPage returns how many slides one printed page holds, or 0 when the
// layout is not a handout layout.
func (p PrintWhat) SlidesPerPage() int {
	switch p {
	case PrintHandouts1:
		return slidesPerHandout1
	case PrintHandouts2:
		return slidesPerHandout2
	case PrintHandouts3:
		return slidesPerHandout3
	case PrintHandouts4:
		return slidesPerHandout4
	case PrintHandouts6:
		return slidesPerHandout6
	case PrintHandouts9:
		return slidesPerHandout9
	case PrintSlides, PrintNotes, PrintOutline:
		return 0
	default:
		return 0
	}
}

// HandoutsPerPage returns the PrintWhat that prints n slides per handout page.
// Values outside {1,2,3,4,6,9} fall back to PrintHandouts1.
func HandoutsPerPage(n int) PrintWhat {
	switch n {
	case slidesPerHandout1:
		return PrintHandouts1
	case slidesPerHandout2:
		return PrintHandouts2
	case slidesPerHandout3:
		return PrintHandouts3
	case slidesPerHandout4:
		return PrintHandouts4
	case slidesPerHandout6:
		return PrintHandouts6
	case slidesPerHandout9:
		return PrintHandouts9
	default:
		return PrintHandouts1
	}
}

// ColorMode selects the ink PowerPoint prints with.
type ColorMode int

const (
	// ColorModeColor prints in full color (the default).
	ColorModeColor ColorMode = iota
	// ColorModeGrayscale prints in grayscale.
	ColorModeGrayscale
	// ColorModeBlackAndWhite prints in pure black and white.
	ColorModeBlackAndWhite
)

// XMLValue returns the ST_PrintColorMode token for the clrMode attribute.
func (c ColorMode) XMLValue() string {
	switch c {
	case ColorModeGrayscale:
		return "gray"
	case ColorModeBlackAndWhite:
		return "bw"
	case ColorModeColor:
		return "clr"
	default:
		return "clr"
	}
}

// Settings is the print configuration written to `<p:prnPr>`.
//
//nolint:govet // Preserve public field order for source compatibility with positional literals.
type Settings struct {
	What              PrintWhat
	ColorMode         ColorMode
	FrameSlides       bool
	ScaleToFitPaper   bool
	PrintHiddenSlides bool
}

// New returns Settings with PowerPoint's defaults: full slides, color, no
// frame, no scaling, hidden slides skipped.
func New() *Settings {
	return &Settings{
		What:      PrintSlides,
		ColorMode: ColorModeColor,
	}
}

// WithPrintWhat selects the printed view.
func (s *Settings) WithPrintWhat(what PrintWhat) *Settings {
	s.What = what
	return s
}

// WithHandouts prints handouts with n slides per page.
func (s *Settings) WithHandouts(slidesPerPage int) *Settings {
	s.What = HandoutsPerPage(slidesPerPage)
	return s
}

// WithColorMode selects color, grayscale, or black and white.
func (s *Settings) WithColorMode(mode ColorMode) *Settings {
	s.ColorMode = mode
	return s
}

// WithFrameSlides draws a thin border around each printed slide.
func (s *Settings) WithFrameSlides(enabled bool) *Settings {
	s.FrameSlides = enabled
	return s
}

// WithScaleToFitPaper scales each slide to the printer's paper size.
func (s *Settings) WithScaleToFitPaper(enabled bool) *Settings {
	s.ScaleToFitPaper = enabled
	return s
}

// WithHiddenSlides includes hidden slides in the printed output.
func (s *Settings) WithHiddenSlides(enabled bool) *Settings {
	s.PrintHiddenSlides = enabled
	return s
}

// IsDefault reports whether the settings match PowerPoint's defaults, in which
// case the `<p:prnPr>` element can be omitted entirely.
func (s *Settings) IsDefault() bool {
	if s == nil {
		return true
	}
	return s.What == PrintSlides &&
		s.ColorMode == ColorModeColor &&
		!s.FrameSlides &&
		!s.ScaleToFitPaper &&
		!s.PrintHiddenSlides
}

func boolAttr(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

// PrnPrXML renders the `<p:prnPr/>` element. It returns an empty string when
// the settings are the PowerPoint defaults.
func (s *Settings) PrnPrXML() string {
	if s.IsDefault() {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<p:prnPr prnWhat="`)
	b.WriteString(s.What.XMLValue())
	b.WriteString(`" clrMode="`)
	b.WriteString(s.ColorMode.XMLValue())
	b.WriteString(`" hiddenSlides="`)
	b.WriteString(boolAttr(s.PrintHiddenSlides))
	b.WriteString(`" scaleToFitPaper="`)
	b.WriteString(boolAttr(s.ScaleToFitPaper))
	b.WriteString(`" frameSlides="`)
	b.WriteString(boolAttr(s.FrameSlides))
	b.WriteString(`"/>`)
	return b.String()
}
