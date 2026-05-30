package editor

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

// RegisterImage adds an image to the presentation or reuses an existing one based on its SHA-256 hash.
func (e *PresentationEditor) RegisterImage(data []byte, format string) (string, error) {
	if e == nil {
		return "", errors.New("editor cannot be nil")
	}
	if len(data) == 0 {
		return "", errors.New("image data cannot be empty")
	}

	hash := sha256.Sum256(data)
	hexHash := hex.EncodeToString(hash[:])

	e.mediaMu.Lock()
	defer e.mediaMu.Unlock()

	if part, ok := e.mediaInventory[hexHash]; ok {
		return part, nil
	}

	partPath := "ppt/media/image" + strconv.Itoa(e.nextMediaNum) + "." + format
	e.nextMediaNum++

	e.parts.Set(partPath, data)
	e.mediaInventory[hexHash] = partPath
	e.mediaInventoryDirty = true
	return partPath, nil
}

// AddSection creates a new grouped section for slides.
func (e *PresentationEditor) AddSection(name string, slideIndices []int) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	ids, err := editorslide.BuildSectionSlideIDs(e.slides, slideIndices)
	if err != nil {
		return err
	}
	return e.applySectionMutation(func(current []Section) ([]Section, error) {
		return editorslide.AddSectionData(current, name, ids, common.NewGUID)
	})
}

// RemoveSection removes a section by name.
func (e *PresentationEditor) RemoveSection(name string) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	return e.applySectionMutation(func(current []Section) ([]Section, error) {
		return editorslide.RemoveSectionData(current, name)
	})
}

// RenameSection renames a section.
func (e *PresentationEditor) RenameSection(oldName, newName string) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	return e.applySectionMutation(func(current []Section) ([]Section, error) {
		return editorslide.RenameSectionData(current, oldName, newName)
	})
}

func (e *PresentationEditor) applySectionMutation(
	mutate func([]Section) ([]Section, error),
) error {
	next, err := mutate(e.sections)
	if err != nil {
		return err
	}
	e.sections = next
	return nil
}

// Sections returns the current section list.
func (e *PresentationEditor) Sections() []Section {
	if e == nil {
		return nil
	}
	return e.sections
}
