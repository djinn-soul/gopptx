package editor

import (
	"archive/zip"
	"fmt"
	"os"
	"sort"
	"strings"

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

	// Iterate over ALL unique part names from both existing state and updates
	allNamesSet := make(map[string]struct{})
	for k := range e.parts {
		allNamesSet[k] = struct{}{}
	}
	for k := range updatedParts {
		allNamesSet[k] = struct{}{}
	}

	allNames := make([]string, 0, len(allNamesSet))
	for k := range allNamesSet {
		allNames = append(allNames, k)
	}
	sort.Strings(allNames)

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
	hasNotesMaster := false
	if _, ok := e.parts["ppt/notesMasters/notesMaster1.xml"]; ok {
		hasNotesMaster = true
	}
	notesMasterRelID := ""
	if hasNotesMaster {
		for _, rel := range e.nonSlideRels {
			if rel.Type == common.RelTypeNotesMaster {
				notesMasterRelID = rel.ID
				break
			}
		}
		if strings.TrimSpace(notesMasterRelID) == "" {
			return nil, fmt.Errorf("notes master part exists but presentation relationship is missing")
		}
	}
	presentationXML, err = rewritePresentationNotesMasterList([]byte(presentationXML), notesMasterRelID, hasNotesMaster)
	if err != nil {
		return nil, err
	}
	out[common.PresentationXMLPath] = []byte(presentationXML)

	hasSections := len(e.sections) > 0
	presentationRelsXML, err := renderPresentationRelsXML(e.nonSlideRels, e.slides, hasSections)
	if err != nil {
		return nil, err
	}
	out[common.PresentationRelPath] = []byte(presentationRelsXML)

	mediaPaths := make([]string, 0, len(e.mediaInventory))
	for _, p := range e.mediaInventory {
		mediaPaths = append(mediaPaths, p)
	}

	chartPaths := make([]string, 0)
	for p := range e.parts {
		if strings.HasPrefix(p, "ppt/charts/chart") && strings.HasSuffix(p, ".xml") {
			chartPaths = append(chartPaths, p)
		}
	}

	notesPaths := make([]string, 0)
	for _, p := range e.notesInventory {
		notesPaths = append(notesPaths, p)
	}

	themePaths := make([]string, 0)
	for p := range e.parts {
		if strings.HasPrefix(p, "ppt/theme/theme") {
			themePaths = append(themePaths, p)
		}
	}

	contentTypesXML, err := rewriteContentTypes(e.parts[common.ContentTypesPath], e.slides, mediaPaths, hasSections, chartPaths, notesPaths, themePaths, hasNotesMaster)
	if err != nil {
		return nil, err
	}
	out[common.ContentTypesPath] = []byte(contentTypesXML)

	if hasSections {
		out["ppt/sectionList.xml"] = []byte(buildSectionListXML(e.sections))
	}

	if hasNotesMaster {
		// Ensure notes master rels are also persisted if they were injected
		if masterRels, ok := e.parts["ppt/notesMasters/_rels/notesMaster1.xml.rels"]; ok {
			out["ppt/notesMasters/_rels/notesMaster1.xml.rels"] = masterRels
		}
	}

	return out, nil
}
