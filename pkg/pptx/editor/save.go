package editor

import (
	"archive/zip"
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Save writes the edited presentation back to a PPTX file.
func (e *PresentationEditor) Save(filePath string) error {
	if e == nil {
		return fmt.Errorf("nil editor")
	}

	updatedParts, err := e.collectUpdatedParts()
	if err != nil {
		return fmt.Errorf("prepare updated parts: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create %s: %w", filePath, err)
	}
	defer func() { _ = file.Close() }()

	zw := zip.NewWriter(file)
	defer func() { _ = zw.Close() }()

	allNames := partNamesSorted(e.parts)
	for _, name := range allNames {
		content := e.parts[name]
		if updated, ok := updatedParts[name]; ok {
			content = updated
		}

		w, err := zw.Create(name)
		if err != nil {
			return fmt.Errorf("create zip entry %q: %w", name, err)
		}
		if _, err := w.Write(content); err != nil {
			return fmt.Errorf("write zip entry %q: %w", name, err)
		}
	}

	return nil
}

func (e *PresentationEditor) collectUpdatedParts() (map[string][]byte, error) {
	out := make(map[string][]byte)

	presentationXML, err := rewritePresentationSlideList([]byte(e.presentationXML), e.slides)
	if err != nil {
		return nil, err
	}
	out[common.PresentationXMLPath] = []byte(presentationXML)

	presentationRelsXML, err := renderPresentationRelsXML(e.nonSlideRels, e.slides)
	if err != nil {
		return nil, err
	}
	out[common.PresentationRelPath] = []byte(presentationRelsXML)

	contentTypesXML, err := rewriteContentTypesSlideOverrides(e.parts[common.ContentTypesPath], e.slides)
	if err != nil {
		return nil, err
	}
	out[common.ContentTypesPath] = []byte(contentTypesXML)

	return out, nil
}
