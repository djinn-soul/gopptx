package pptx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"

	"github.com/vegito/goppt/internal/pptxxml"
)

// Create builds a valid PPTX with generated slide titles.
func Create(title string, slideCount int) ([]byte, error) {
	if title == "" {
		return nil, fmt.Errorf("presentation title cannot be empty")
	}
	if slideCount < 1 {
		return nil, fmt.Errorf("slide count must be at least 1")
	}

	slides := make([]SlideContent, 0, slideCount)
	for i := 1; i <= slideCount; i++ {
		slideTitle := title
		if i > 1 {
			slideTitle = fmt.Sprintf("Slide %d", i)
		}
		slides = append(slides, NewSlide(slideTitle))
	}

	return CreateWithSlides(title, slides)
}

// CreateWithSlides builds a PPTX from caller-provided slide content.
func CreateWithSlides(title string, slides []SlideContent) ([]byte, error) {
	if title == "" {
		return nil, fmt.Errorf("presentation title cannot be empty")
	}
	if len(slides) == 0 {
		return nil, fmt.Errorf("at least one slide is required")
	}
	for i, slide := range slides {
		if err := validateSlide(slide, i+1); err != nil {
			return nil, err
		}
	}

	buf := bytes.NewBuffer(nil)
	zw := zip.NewWriter(buf)
	count := len(slides)

	if err := writePackageFiles(zw, title, slides, count); err != nil {
		_ = zw.Close()
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WriteFile is a convenience helper that writes the generated PPTX to disk.
func WriteFile(path string, title string, slides []SlideContent) error {
	data, err := CreateWithSlides(title, slides)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func writePackageFiles(zw *zip.Writer, title string, slides []SlideContent, slideCount int) error {
	files := []struct {
		name    string
		content string
	}{
		{"[Content_Types].xml", pptxxml.ContentTypes(slideCount)},
		{"_rels/.rels", pptxxml.RootRelationships()},
		{"ppt/_rels/presentation.xml.rels", pptxxml.PresentationRelationships(slideCount)},
		{"ppt/presentation.xml", pptxxml.Presentation(title, slideCount)},
		{"ppt/slideLayouts/slideLayout1.xml", pptxxml.SlideLayout()},
		{"ppt/slideLayouts/_rels/slideLayout1.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideMasters/slideMaster1.xml", pptxxml.SlideMaster()},
		{"ppt/slideMasters/_rels/slideMaster1.xml.rels", pptxxml.SlideMasterRelationships()},
		{"ppt/theme/theme1.xml", pptxxml.Theme()},
		{"docProps/core.xml", pptxxml.CoreProperties(title)},
		{"docProps/app.xml", pptxxml.AppProperties(slideCount)},
	}

	for _, item := range files {
		if err := writeFile(zw, item.name, item.content); err != nil {
			return err
		}
	}

	for i, slide := range slides {
		slideNumber := i + 1
		slideXML := pptxxml.SlideWithContent(slide.Title, slide.Bullets)
		slidePath := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
		if err := writeFile(zw, slidePath, slideXML); err != nil {
			return err
		}

		relsPath := fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", slideNumber)
		if err := writeFile(zw, relsPath, pptxxml.SlideRelationships()); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(zw *zip.Writer, path string, content string) error {
	w, err := zw.Create(path)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(content))
	return err
}
