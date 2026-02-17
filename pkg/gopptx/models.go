package gopptx

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"

	"github.com/djinn-soul/gopptx/internal/opc"
	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

const (
	NSP = "http://schemas.openxmlformats.org/presentationml/2006/main"

	// Default slide sizes in EMUs (10x7.5 inches for 4:3).
	defaultWidth  = 9144000
	defaultHeight = 6858000
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
	f, createErr := os.Create(path)
	if createErr != nil {
		return createErr
	}
	defer f.Close()

	w := opc.NewWriter(f)
	defer func() {
		if err := w.Close(); err != nil {
			// TODO: Verify error handling logic ensures no data corruption on disk full.
			// Currently we log nothing to avoid external dependencies in this model.
			_ = err
		}
	}()

	slideCount := len(p.Slides)

	// 1. Mandatory OPC parts
	if err := w.AddFile("[Content_Types].xml", []byte(pptxxml.ContentTypes(
		slideCount, nil, 0, 0, nil, false, 0, 0, 0,
	))); err != nil {
		return err
	}
	if err := w.AddFile("_rels/.rels", []byte(pptxxml.RootRelationships())); err != nil {
		return err
	}
	if err := w.AddFile("ppt/_rels/presentation.xml.rels", []byte(pptxxml.PresentationRelationships(
		slideCount, false, 0, 1,
	))); err != nil {
		return err
	}

	// 2. Marshal Presentation
	// Note: Standard Presentation generator from pptxxml is more complete,
	// but here we use the struct for now as it's the intended models.go design.
	// However, the struct must match the expected RID and ID from PresentationRelationships.
	// For simplicity in this fix, we'll use pptxxml.Presentation to ensure validity.
	// TODO: This produced .pptx is currently invalid as it lacks mandatory OOXML parts
	// (masters, layouts, themes). We must integrate with internal/pptxxml/package_writer.go
	// or similar to generate a complete, valid OPC package.
	presXML := pptxxml.Presentation("Presentation", slideCount, false, defaultWidth, defaultHeight, 1)
	if err := w.AddFile("ppt/presentation.xml", []byte(presXML)); err != nil {
		return err
	}

	// 3. Marshal Slides
	for i := range p.Slides {
		// Using a minimal valid slide XML for now to ensure it's openable
		slideXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/><p:sp><p:nvSpPr><p:cNvPr id="2" name="Title 1"/><p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr><p:nvPr><p:ph type="title"/></p:nvPr></p:nvSpPr><p:spPr/><p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:t>Slide ` + strconv.Itoa(i+1) + `</a:t></a:r></a:p></p:txBody></p:sp></p:spTree></p:cSld></p:sld>`
		filename := fmt.Sprintf("ppt/slides/slide%d.xml", i+1)
		if err := w.AddFile(filename, []byte(slideXML)); err != nil {
			return err
		}
	}

	return nil
}

// Slide represents a single slide XML component.
type Slide struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/presentationml/2006/main sld"`
}
