package elements

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// SlideMaster defines content that appears on all slides.
type SlideMaster struct {
	Background   *SlideBackground
	Shapes       []shapes.ShapeDefinition
	Images       []shapes.Image
	FooterText   string
	ColorMapping *ColorMapping
}

// ColorMapping defines how theme colors map to functional roles on slides.
type ColorMapping struct {
	BG1 string // e.g., "lt1", "dk1", "accent1"
	TX1 string // e.g., "dk1", "lt1"
}

// NewMaster creates a blank slide master.
func NewMaster() *SlideMaster {
	return &SlideMaster{}
}

// WithBackground sets the background for the slide master.
func (m *SlideMaster) WithBackground(bg SlideBackground) *SlideMaster {
	m.Background = &bg
	return m
}

// AddShape adds a shape (e.g., a logo) to the slide master.
func (m *SlideMaster) AddShape(sd shapes.ShapeDefinition) *SlideMaster {
	m.Shapes = append(m.Shapes, sd)
	return m
}

// AddImage adds an image to the slide master.
func (m *SlideMaster) AddImage(img shapes.Image) *SlideMaster {
	m.Images = append(m.Images, img)
	return m
}

// WithFooter sets the footer text for the slide master.
func (m *SlideMaster) WithFooter(text string) *SlideMaster {
	m.FooterText = text
	return m
}

// WithColorMapping sets the color mapping for the slide master.
func (m *SlideMaster) WithColorMapping(bg1, tx1 string) *SlideMaster {
	m.ColorMapping = &ColorMapping{BG1: bg1, TX1: tx1}
	return m
}
