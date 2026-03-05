package editor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type imagePathCacheEntry struct {
	data     []byte
	format   string
	partPath string
}

func (e *PresentationEditor) registerImageFromPath(imagePath, formatHint string) (string, error) {
	cleanPath := filepath.Clean(imagePath)
	e.imagePathMu.RLock()
	entry, ok := e.imagePathCache[cleanPath]
	e.imagePathMu.RUnlock()
	if ok && entry.partPath != "" {
		return entry.partPath, nil
	}

	data, format, err := e.loadImageFromPath(cleanPath, formatHint)
	if err != nil {
		return "", err
	}

	partPath, err := e.RegisterImage(data, format)
	if err != nil {
		return "", err
	}

	e.imagePathMu.Lock()
	e.imagePathCache[cleanPath] = imagePathCacheEntry{
		data:     nil,
		format:   format,
		partPath: partPath,
	}
	e.imagePathMu.Unlock()
	return partPath, nil
}

func (e *PresentationEditor) loadImageFromPath(cleanPath, formatHint string) ([]byte, string, error) {
	e.imagePathMu.RLock()
	entry, ok := e.imagePathCache[cleanPath]
	e.imagePathMu.RUnlock()
	if ok && len(entry.data) > 0 {
		format := normalizeImageFormatHint(formatHint)
		if format == "" {
			format = entry.format
		}
		return entry.data, format, nil
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, "", err
	}

	format := normalizeImageFormatHint(formatHint)
	if format == "" {
		format = normalizeImageFormatHint(filepath.Ext(cleanPath))
	}
	if format == "" {
		return nil, "", fmt.Errorf("image path %q has no detectable format", cleanPath)
	}

	e.imagePathMu.Lock()
	entry = e.imagePathCache[cleanPath]
	if len(entry.data) == 0 {
		entry.data = data
	}
	if entry.format == "" {
		entry.format = format
	}
	e.imagePathCache[cleanPath] = entry
	e.imagePathMu.Unlock()
	return entry.data, entry.format, nil
}

func normalizeImageFormatHint(format string) string {
	trimmed := strings.TrimSpace(strings.ToLower(format))
	return strings.TrimPrefix(trimmed, ".")
}

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
