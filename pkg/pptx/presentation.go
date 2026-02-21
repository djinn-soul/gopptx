package pptx

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation"
)

type (
	// Metadata defines non-content properties of a PPTX.
	Metadata = presentation.Metadata
	// MetadataFields defines the basic descriptive fields for a PPTX.
	MetadataFields = common.Metadata
	// SlideSize defines presentation dimensions in EMUs.
	SlideSize = presentation.SlideSize
	// CustomXMLPart represents an embedded custom XML document.
	CustomXMLPart = common.CustomXMLPart
)

// SlideSize4x3 returns the standard 4:3 slide size.
func SlideSize4x3() SlideSize {
	return presentation.GetSlideSize4x3()
}

// SlideSize16x9 returns the standard 16:9 widescreen slide size.
func SlideSize16x9() SlideSize {
	return presentation.GetSlideSize16x9()
}

// Create builds a valid PPTX with generated slide titles.
func Create(title string, slideCount int) ([]byte, error) {
	if title == "" {
		return nil, errors.New("presentation title cannot be empty")
	}
	if slideCount < 1 {
		return nil, errors.New("slide count must be at least 1")
	}

	slides := make([]SlideContent, 0, slideCount)
	for i := 1; i <= slideCount; i++ {
		slideTitle := title
		if i > 1 {
			slideTitle = fmt.Sprintf("Slide %d", i)
		}
		slides = append(slides, NewSlide(slideTitle))
	}

	return CreateWithMetadata(
		Metadata{Metadata: common.Metadata{Title: title}},
		slides,
	)
}

// CreateWithSlides builds a PPTX from caller-provided slide content.
func CreateWithSlides(title string, slides []SlideContent) ([]byte, error) {
	return CreateWithMetadata(
		Metadata{Metadata: common.Metadata{Title: title}},
		slides,
	)
}

// CreateWithMetadata builds a PPTX from metadata and caller-provided slide content.
func CreateWithMetadata(meta Metadata, slides []SlideContent) ([]byte, error) {
	if meta.Title == "" {
		return nil, errors.New("presentation title cannot be empty")
	}
	if len(slides) == 0 {
		return nil, errors.New("at least one slide is required")
	}
	for i, slide := range slides {
		if err := slide.Validate(i + 1); err != nil {
			return nil, err
		}
	}
	if meta.NotesMaster != nil {
		if err := meta.NotesMaster.Validate(); err != nil {
			return nil, err
		}
	}

	if meta.SlideSize.Width == 0 || meta.SlideSize.Height == 0 {
		meta.SlideSize = SlideSize4x3()
	}

	buf := bytes.NewBuffer(nil)
	zw := zip.NewWriter(buf)
	count := len(slides)

	if err := presentation.WritePresentationPackage(zw, meta, slides, count); err != nil {
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
	data, err := CreateWithMetadata(
		Metadata{Metadata: common.Metadata{Title: title, SlideSize: SlideSize4x3()}},
		slides,
	)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
