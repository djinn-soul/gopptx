package editor

import (
	"archive/zip"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const commentAuthorsRelType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/commentAuthors"

// Save writes the edited presentation back to a PPTX file.
func (e *PresentationEditor) Save(filePath string) error {
	if e == nil {
		return errors.New("nil editor")
	}

	// Materialize all lazy parts into memory and release the source file handle.
	if err := e.parts.Materialize(); err != nil {
		return fmt.Errorf("failed to materialize lazy PPTX parts from source archive: %w", err)
	}

	updatedParts, err := e.collectUpdatedParts()
	if err != nil {
		return fmt.Errorf("prepare updated parts: %w", err)
	}

	// Pre-verification: Ensure all required parts are available before we start writing to disk.
	// This prevents truncated/corrupt files if a lazy-read fails mid-save.
	keys := e.parts.Keys()
	for _, name := range keys {
		if _, updatedOK := updatedParts[name]; !updatedOK {
			if _, ok := e.parts.Get(name); !ok {
				return fmt.Errorf("critical failure: part %q missing from store during pre-save scan", name)
			}
		}
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
	for _, k := range keys {
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
		var content []byte
		if updated, updatedOK := updatedParts[name]; updatedOK {
			content = updated
		} else {
			data, partOK := e.parts.Get(name)
			if !partOK {
				return fmt.Errorf("failed to retrieve part %q during save", name)
			}
			content = data
		}

		w, createErr := zw.Create(name)
		if createErr != nil {
			return fmt.Errorf("create zip entry %q: %w", name, createErr)
		}
		if _, writeErr := w.Write(content); writeErr != nil {
			return fmt.Errorf("write zip entry %q: %w", name, writeErr)
		}
	}

	return nil
}

func (e *PresentationEditor) collectUpdatedParts() (map[string][]byte, error) {
	out := make(map[string][]byte)

	if err := e.serializeAuthorCacheIfPopulated(); err != nil {
		return nil, err
	}

	// Check for commentAuthors existence and relationship injection
	hasCommentAuthors := e.parts.Has("ppt/commentAuthors.xml")
	if hasCommentAuthors {
		e.ensureCommentAuthorsRelationship()
	}

	presentationXML, err := e.renderPresentationXMLWithSections()
	if err != nil {
		return nil, err
	}
	out[common.PresentationXMLPath] = []byte(presentationXML)

	hasSections := len(e.sections) > 0
	hasNotesMaster := e.parts.Has("ppt/notesMasters/notesMaster1.xml")
	presentationRelsXML, err := renderPresentationRelsXML(e.nonSlideRels, e.slides, hasSections)
	if err != nil {
		return nil, err
	}
	out[common.PresentationRelPath] = []byte(presentationRelsXML)

	// Persist Core Properties
	corePropsXML, err := renderCoreProperties(e.metadata.CoreProperties)
	if err != nil {
		return nil, fmt.Errorf("render core properties: %w", err)
	}
	out[common.CorePropsPath] = corePropsXML

	mediaPaths := mapValues(e.mediaInventory)
	filteredChartPaths := filterXMLPartPaths(e.parts.KeysWithPrefix("ppt/charts/chart"))
	notesPaths := mapValues(e.notesInventory)

	themePaths := e.parts.KeysWithPrefix("ppt/theme/theme")
	layoutPaths := e.parts.KeysWithPrefix("ppt/slideLayouts/slideLayout")
	masterPaths := e.parts.KeysWithPrefix("ppt/slideMasters/slideMaster")

	filteredCommentPaths := filterXMLPartPaths(e.parts.KeysWithPrefix("ppt/comments/comment"))

	contentTypesData, _ := e.parts.Get(common.ContentTypesPath)
	contentTypesXML, err := rewriteContentTypes(
		contentTypesData,
		e.slides,
		mediaPaths,
		hasSections,
		filteredChartPaths,
		notesPaths,
		themePaths,
		layoutPaths,
		masterPaths,
		hasNotesMaster,
		hasCommentAuthors,
		filteredCommentPaths,
	)
	if err != nil {
		return nil, err
	}
	out[common.ContentTypesPath] = []byte(contentTypesXML)

	if hasSections {
		out["ppt/sectionList.xml"] = []byte(buildSectionListXML(e.sections))
	}

	if hasNotesMaster {
		// Ensure notes master rels are also persisted if they were injected
		if masterRels, ok := e.parts.Get("ppt/notesMasters/_rels/notesMaster1.xml.rels"); ok {
			out["ppt/notesMasters/_rels/notesMaster1.xml.rels"] = masterRels
		}
	}

	return out, nil
}

func (e *PresentationEditor) serializeAuthorCacheIfPopulated() error {
	e.authorCacheMu.RLock()
	cachePopulated := e.authorCache != nil
	e.authorCacheMu.RUnlock()
	if !cachePopulated {
		return nil
	}

	authors, err := e.GetAuthors()
	if err != nil {
		return fmt.Errorf("get authors: %w", err)
	}
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].ID < authors[j].ID
	})

	xmlContent := pptxxml.CommentAuthorsXML(authors)
	e.parts.Set("ppt/commentAuthors.xml", []byte(xmlContent))
	return nil
}

func (e *PresentationEditor) ensureCommentAuthorsRelationship() {
	for _, rel := range e.nonSlideRels {
		if rel.Type == commentAuthorsRelType {
			return
		}
	}

	relID := fmt.Sprintf("rId%d", e.nextRelIDNum)
	e.nextRelIDNum++
	e.nonSlideRels = append(e.nonSlideRels, common.EditorRelationship{
		ID:     relID,
		Type:   commentAuthorsRelType,
		Target: "commentAuthors.xml",
	})
}

func (e *PresentationEditor) renderPresentationXMLWithSections() (string, error) {
	presentationXML, err := rewritePresentationSlideList([]byte(e.presentationXML), e.slides)
	if err != nil {
		return "", err
	}

	hasNotesMaster := e.parts.Has("ppt/notesMasters/notesMaster1.xml")
	notesMasterRelID, err := e.resolveNotesMasterRelID(hasNotesMaster)
	if err != nil {
		return "", err
	}

	presentationXML, err = rewritePresentationNotesMasterList(
		[]byte(presentationXML),
		notesMasterRelID,
		hasNotesMaster,
	)
	if err != nil {
		return "", err
	}

	if len(e.sections) == 0 {
		return presentationXML, nil
	}

	presentationXML, err = rewritePresentationSections([]byte(presentationXML), e.sections)
	if err != nil {
		return "", fmt.Errorf("rewrite sections: %w", err)
	}
	return presentationXML, nil
}

func (e *PresentationEditor) resolveNotesMasterRelID(hasNotesMaster bool) (string, error) {
	if !hasNotesMaster {
		return "", nil
	}

	for _, rel := range e.nonSlideRels {
		if rel.Type == common.RelTypeNotesMaster {
			return rel.ID, nil
		}
	}
	return "", errors.New("notes master part exists but presentation relationship is missing")
}

func mapValues(m map[string]string) []string {
	values := make([]string, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

func filterXMLPartPaths(paths []string) []string {
	filtered := make([]string, 0, len(paths))
	for _, p := range paths {
		if strings.HasSuffix(p, ".xml") {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
