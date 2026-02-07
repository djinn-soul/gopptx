package pptx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// Save writes the edited package to disk.
func (e *PresentationEditor) Save(filePath string) error {
	data, err := e.Bytes()
	if err != nil {
		return err
	}
	parent := filepath.Dir(filePath)
	if parent != "." && parent != "" {
		if err := os.MkdirAll(parent, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(filePath, data, 0o600)
}

// Bytes serializes the edited package into a PPTX byte slice.
func (e *PresentationEditor) Bytes() ([]byte, error) {
	if e == nil {
		return nil, fmt.Errorf("editor cannot be nil")
	}
	if err := e.rebuildControlParts(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, partName := range partNamesSorted(e.parts) {
		w, err := zw.Create(partName)
		if err != nil {
			_ = zw.Close()
			return nil, err
		}
		if _, err := w.Write(e.parts[partName]); err != nil {
			_ = zw.Close()
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *PresentationEditor) rebuildControlParts() error {
	if e == nil {
		return fmt.Errorf("editor cannot be nil")
	}

	currentPresentationXML := e.parts[presentationXMLPath]
	if len(currentPresentationXML) == 0 {
		currentPresentationXML = []byte(e.presentationXML)
	}
	presentationXML, err := rewritePresentationSlideList(currentPresentationXML, e.slides)
	if err != nil {
		return err
	}
	e.parts[presentationXMLPath] = []byte(presentationXML)
	e.presentationXML = presentationXML

	relsXML, err := renderPresentationRelsXML(e.nonSlideRels, e.slides)
	if err != nil {
		return err
	}
	e.parts[presentationRelPath] = []byte(relsXML)

	contentTypesXML, err := rewriteContentTypesSlideOverrides(e.parts[contentTypesPath], e.slides)
	if err != nil {
		return err
	}
	e.parts[contentTypesPath] = []byte(contentTypesXML)
	return nil
}
