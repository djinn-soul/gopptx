package gopptx

import "encoding/xml"

const (
	NS_P = "http://schemas.openxmlformats.org/presentationml/2006/main"
)

// Presentation represents the main presentation XML component.
type Presentation struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/presentationml/2006/main presentation"`
	Slides  []*Slide `xml:"sldIdLst>sldId"`
}

// AddSlide adds a new blank slide to the presentation.
func (p *Presentation) AddSlide() *Slide {
	slide := &Slide{}
	p.Slides = append(p.Slides, slide)
	return slide
}

// Slide represents a single slide XML component.
type Slide struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/presentationml/2006/main sld"`
}
