package slide

import (
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type NewGUIDFn func() (string, error)

func BuildSectionSlideIDs(slides []common.EditorSlideRef, slideIndices []int) ([]int64, error) {
	ids := make([]int64, 0, len(slideIndices))
	for _, idx := range slideIndices {
		if idx < 0 || idx >= len(slides) {
			return nil, fmt.Errorf("slide index %d out of range", idx)
		}
		ids = append(ids, slides[idx].SlideID)
	}
	return ids, nil
}

func AddSectionData(current []SectionData, name string, slideIDs []int64, newGUID NewGUIDFn) ([]SectionData, error) {
	if name == "" {
		return nil, errors.New("section name cannot be empty")
	}
	guid, err := newGUID()
	if err != nil {
		return nil, fmt.Errorf("generate section GUID: %w", err)
	}
	next := make([]SectionData, 0, len(current)+1)
	next = append(next, current...)
	next = append(next, SectionData{
		Name:     name,
		GUID:     guid,
		SlideIDs: slideIDs,
	})
	return next, nil
}

func RemoveSectionData(current []SectionData, name string) ([]SectionData, error) {
	next := make([]SectionData, 0, len(current))
	found := false
	for _, s := range current {
		if s.Name == name {
			found = true
			continue
		}
		next = append(next, s)
	}
	if !found {
		return nil, fmt.Errorf("section %q not found", name)
	}
	return next, nil
}

func RenameSectionData(current []SectionData, oldName, newName string) ([]SectionData, error) {
	if newName == "" {
		return nil, errors.New("new section name cannot be empty")
	}
	next := make([]SectionData, 0, len(current))
	found := false
	for _, s := range current {
		if s.Name == oldName {
			s.Name = newName
			found = true
		}
		next = append(next, s)
	}
	if !found {
		return nil, fmt.Errorf("section %q not found", oldName)
	}
	return next, nil
}
