package gopptx

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/internal/opc"
)

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

// Save writes the presentation to a .pptx file.
func (p *Presentation) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := opc.NewWriter(f)
	defer w.Close()

	// 1. Marshal Presentation
	presXML, err := xml.Marshal(p)
	if err != nil {
		return err
	}
	if err := w.AddFile("ppt/presentation.xml", presXML); err != nil {
		return err
	}

	// 2. Marshal Slides
	for i, slide := range p.Slides {
		slideXML, err := xml.Marshal(slide)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("ppt/slides/slide%d.xml", i+1)
		if err := w.AddFile(filename, slideXML); err != nil {
			return err
		}
	}

	return nil
}

// Slide represents a single slide XML component.
type Slide struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/presentationml/2006/main sld"`
}
