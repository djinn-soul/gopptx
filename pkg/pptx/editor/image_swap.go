package editor

import (
	"errors"
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// ListSlideImages returns image relationships for a slide in relationship order.
func (e *PresentationEditor) ListSlideImages(slideIndex int) ([]common.SlideImageRef, error) {
	if e == nil {
		return nil, errors.New("editor cannot be nil")
	}
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, fmt.Errorf("slide index %d out of range [0,%d)", slideIndex, len(e.slides))
	}
	rels, err := e.slideRelationships(e.slides[slideIndex].Part)
	if err != nil {
		return nil, err
	}

	out := make([]common.SlideImageRef, 0)
	for _, rel := range rels {
		if rel.Type != common.RelTypeImage {
			continue
		}
		out = append(out, common.SlideImageRef{
			Index:  len(out),
			RelID:  rel.ID,
			Target: rel.Target,
		})
	}
	return out, nil
}

// SwapImageByIndex replaces one slide image relationship target by image index.
func (e *PresentationEditor) SwapImageByIndex(slideIndex, imageIndex int, data []byte, format string) error {
	images, err := e.ListSlideImages(slideIndex)
	if err != nil {
		return err
	}
	if imageIndex < 0 || imageIndex >= len(images) {
		return fmt.Errorf("image index %d out of range (found %d images)", imageIndex, len(images))
	}
	return e.swapImageByRelID(slideIndex, images[imageIndex].RelID, data, format)
}

// SwapImageByRelID replaces one slide image relationship target by relationship ID.
func (e *PresentationEditor) SwapImageByRelID(slideIndex int, relID string, data []byte, format string) error {
	if strings.TrimSpace(relID) == "" {
		return errors.New("relationship id cannot be empty")
	}
	return e.swapImageByRelID(slideIndex, relID, data, format)
}

func (e *PresentationEditor) swapImageByRelID(slideIndex int, relID string, data []byte, format string) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range [0,%d)", slideIndex, len(e.slides))
	}
	slidePart := e.slides[slideIndex].Part
	relsPath := common.SlideRelsPartName(slidePart)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return fmt.Errorf("missing slide relationships part %q", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return fmt.Errorf("parse %s: %w", relsPath, err)
	}

	updated := false
	for i := range rels {
		if rels[i].ID == relID && rels[i].Type == common.RelTypeImage {
			updated = true
			break
		}
	}
	if !updated {
		return fmt.Errorf("image relationship %q not found on slide %d", relID, slideIndex)
	}
	newPart, err := e.RegisterImage(data, format)
	if err != nil {
		return err
	}
	for i := range rels {
		if rels[i].ID == relID && rels[i].Type == common.RelTypeImage {
			rels[i].Target = common.MakeRelativePath(slidePart, newPart)
			break
		}
	}

	rendered := renderRelationshipsXML(rels)
	e.parts.Set(relsPath, []byte(rendered))
	return nil
}
