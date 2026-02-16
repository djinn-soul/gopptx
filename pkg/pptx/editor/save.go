package editor

import (
	"archive/zip"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
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

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create %s: %w", filePath, err)
	}
	defer func() { _ = file.Close() }()

	zw := zip.NewWriter(file)
	defer func() { _ = zw.Close() }()

	// Iterate over ALL unique part names from both existing state and updates
	allNamesSet := make(map[string]struct{})
	for _, k := range e.parts.Keys() {
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

	// Serialize authors if cache is populated
	e.authorCacheMu.RLock()
	cachePopulated := e.authorCache != nil
	e.authorCacheMu.RUnlock()

	if cachePopulated {
		// Convert map to slice
		authors, _ := e.GetAuthors() // Acquires lock internally
		// Sort by ID
		sort.Slice(authors, func(i, j int) bool {
			return authors[i].ID < authors[j].ID
		})

		xmlContent := pptxxml.CommentAuthorsXML(authors)
		e.parts.Set("ppt/commentAuthors.xml", []byte(xmlContent))
	}

	// Check for commentAuthors existence and relationship injection
	hasCommentAuthors := e.parts.Has("ppt/commentAuthors.xml")
	if hasCommentAuthors {
		found := false
		for _, rel := range e.nonSlideRels {
			if rel.Type == commentAuthorsRelType {
				found = true
				break
			}
		}
		if !found {
			// Add rel
			relID := fmt.Sprintf("rId%d", e.nextRelIDNum)
			e.nextRelIDNum++
			e.nonSlideRels = append(e.nonSlideRels, common.EditorRelationship{
				ID:     relID,
				Type:   commentAuthorsRelType,
				Target: "commentAuthors.xml",
			})
		}
	}

	presentationXML, err := rewritePresentationSlideList([]byte(e.presentationXML), e.slides)
	if err != nil {
		return nil, err
	}
	hasNotesMaster := e.parts.Has("ppt/notesMasters/notesMaster1.xml")
	notesMasterRelID := ""
	if hasNotesMaster {
		for _, rel := range e.nonSlideRels {
			if rel.Type == common.RelTypeNotesMaster {
				notesMasterRelID = rel.ID
				break
			}
		}
		if strings.TrimSpace(notesMasterRelID) == "" {
			return nil, errors.New("notes master part exists but presentation relationship is missing")
		}
	}
	presentationXML, err = rewritePresentationNotesMasterList([]byte(presentationXML), notesMasterRelID, hasNotesMaster)
	if err != nil {
		return nil, err
	}
	out[common.PresentationXMLPath] = []byte(presentationXML)

	// Inject Sections into presentation.xml extension list (Required for PPT 2010+)
	if len(e.sections) > 0 {
		pXML, rewriteErr := rewritePresentationSections([]byte(presentationXML), e.sections)
		if rewriteErr != nil {
			return nil, fmt.Errorf("rewrite sections: %w", rewriteErr)
		}
		out[common.PresentationXMLPath] = []byte(pXML)
	}

	hasSections := len(e.sections) > 0
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

	mediaPaths := make([]string, 0, len(e.mediaInventory))
	for _, p := range e.mediaInventory {
		mediaPaths = append(mediaPaths, p)
	}

	chartPaths := e.parts.KeysWithPrefix("ppt/charts/chart")
	filteredChartPaths := make([]string, 0, len(chartPaths))
	for _, p := range chartPaths {
		if strings.HasSuffix(p, ".xml") {
			filteredChartPaths = append(filteredChartPaths, p)
		}
	}

	notesPaths := make([]string, 0)
	for _, p := range e.notesInventory {
		notesPaths = append(notesPaths, p)
	}

	themePaths := e.parts.KeysWithPrefix("ppt/theme/theme")
	layoutPaths := e.parts.KeysWithPrefix("ppt/slideLayouts/slideLayout")
	masterPaths := e.parts.KeysWithPrefix("ppt/slideMasters/slideMaster")

	commentPaths := e.parts.KeysWithPrefix("ppt/comments/comment")
	// Filter just in case
	filteredCommentPaths := make([]string, 0, len(commentPaths))
	for _, p := range commentPaths {
		if strings.HasSuffix(p, ".xml") {
			filteredCommentPaths = append(filteredCommentPaths, p)
		}
	}

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
