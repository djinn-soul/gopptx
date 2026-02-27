package pptx

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation"
	"github.com/djinn-soul/gopptx/pkg/pptx/validation/logical"
	"github.com/djinn-soul/gopptx/pkg/pptx/validation/structural"
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
	// Section defines a presentation section.
	Section = presentation.Section
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

// Validate checks a PPTX byte slice for structural issues.
func Validate(pptxData []byte) ([]structural.Issue, error) {
	// Try full editor open first for deep validation (relationships, etc.)
	ed, err := editor.OpenPresentationEditorFromBytes(pptxData)
	if err == nil {
		defer func() { _ = ed.Close() }()
		return ed.Validate(), nil
	}

	// If full open fails, it might be due to invalid XML or missing parts.
	// Fallback to structural validation on raw parts.
	ps, psErr := editor.OpenPartStoreFromBytes(pptxData)
	if psErr != nil {
		return nil, psErr
	}
	defer func() { _ = ps.Close() }()

	v := structural.NewValidator(ps)
	v.AddChecker(&logical.Checker{})
	return v.Validate(), nil
}

// Repair attempts to fix structural issues in a PPTX byte slice and returns the fixed bytes.
func Repair(pptxData []byte) ([]byte, structural.RepairResult, error) {
	ps, err := editor.OpenPartStoreFromBytes(pptxData)
	if err != nil {
		return nil, structural.RepairResult{}, err
	}
	defer func() {
		if ps != nil {
			_ = ps.Close()
		}
	}()

	v := structural.NewValidator(ps)
	v.AddChecker(&logical.Checker{})
	issues := v.Validate()
	if len(issues) == 0 {
		return pptxData, structural.RepairResult{}, nil
	}

	r := structural.NewRepairer(ps)
	result := r.Repair(issues)

	// Save the repaired parts
	// We need a temporary editor to save, or we can use a simpler save mechanism.
	// Since we already have the parts, we can try to re-open the editor from these parts.
	// If it still fails, it's a critical error.

	// Create a new editor from the repaired part store
	ed, err := editor.NewPresentationEditorFromParts(ps)
	if err != nil {
		// If it still fails to open as a full editor, we might have repaired some things but not enough.
		// For now, we'll try a raw save if possible, or return the error.
		return nil, result, fmt.Errorf("repair produced unusable package: %w", err)
	}
	// Transfer ownership of ps to ed; ed.Close() will close ps.
	// Set ps to nil to prevent double-close from the deferred close above.
	ps = nil
	defer func() { _ = ed.Close() }()

	// Use a secure temporary file with proper cleanup
	tmpFile, err := os.CreateTemp("", "repair-*.pptx")
	if err != nil {
		return nil, result, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	if err := ed.Save(tmpPath); err != nil {
		return nil, result, err
	}
	repairedData, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, result, err
	}

	return repairedData, result, nil
}
