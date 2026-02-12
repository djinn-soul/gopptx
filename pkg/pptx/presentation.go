package pptx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation"
)

type (
	// PresentationMetadata defines non-content properties of a PPTX.
	PresentationMetadata = presentation.PresentationMetadata
	// SlideSize defines presentation dimensions in EMUs.
	SlideSize = presentation.SlideSize
	// CustomXMLPart represents an embedded custom XML document.
	CustomXMLPart = common.CustomXMLPart
)

// Default slide sizes.
var (
	SlideSize4x3  = presentation.SlideSize4x3
	SlideSize16x9 = presentation.SlideSize16x9
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

	return CreateWithMetadata(PresentationMetadata{PresentationMetadata: common.PresentationMetadata{Title: title}}, slides)
}

// CreateWithSlides builds a PPTX from caller-provided slide content.
func CreateWithSlides(title string, slides []SlideContent) ([]byte, error) {
	return CreateWithMetadata(PresentationMetadata{PresentationMetadata: common.PresentationMetadata{Title: title}}, slides)
}

// CreateWithMetadata builds a PPTX from metadata and caller-provided slide content.
func CreateWithMetadata(meta PresentationMetadata, slides []SlideContent) ([]byte, error) {
	if meta.Title == "" {
		return nil, fmt.Errorf("presentation title cannot be empty")
	}
	if len(slides) == 0 {
		return nil, fmt.Errorf("at least one slide is required")
	}
	for i, slide := range slides {
		if err := slide.Validate(i + 1); err != nil {
			return nil, err
		}
	}

	if meta.SlideSize.Width == 0 || meta.SlideSize.Height == 0 {
		meta.SlideSize = SlideSize4x3
	}

	buf := bytes.NewBuffer(nil)
	zw := zip.NewWriter(buf)
	count := len(slides)

	if err := presentation.WritePackageFiles(zw, meta, slides, count); err != nil {
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
	data, err := CreateWithMetadata(PresentationMetadata{PresentationMetadata: common.PresentationMetadata{Title: title, SlideSize: SlideSize4x3}}, slides)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
