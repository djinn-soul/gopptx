package editor

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

const commentAuthorsRelType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/commentAuthors"

// Save writes the edited presentation back to a PPTX file.
//
//nolint:gocognit // Save flow intentionally sequences materialize/validate/write/cleanup steps with explicit guards.
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

		var (
			w         io.Writer
			createErr error
		)
		if editorslide.SaveZipMethod(name) == zip.Store {
			header := &zip.FileHeader{Name: name, Method: zip.Store}
			w, createErr = zw.CreateHeader(header)
		} else {
			w, createErr = zw.Create(name)
		}
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

	e.authorCacheMu.RLock()
	authorCache := e.authorCache
	e.authorCacheMu.RUnlock()
	authorXML, hasAuthors, err := editorslide.SerializeCommentAuthorsIfPopulated(authorCache, e.GetAuthors)
	if err != nil {
		return nil, err
	}
	if hasAuthors {
		e.parts.Set("ppt/commentAuthors.xml", authorXML)
	}

	customXMLPropsPaths, err := e.processCustomXMLParts(out)
	if err != nil {
		return nil, err
	}
	e.filterRootCustomXMLRelationships(out)

	// Check for commentAuthors existence and relationship injection
	hasCommentAuthors := e.parts.Has("ppt/commentAuthors.xml")
	if hasCommentAuthors {
		e.nonSlideRels, e.nextRelIDNum = editorslide.EnsureCommentAuthorsRelationship(
			e.nonSlideRels,
			e.nextRelIDNum,
			commentAuthorsRelType,
			"commentAuthors.xml",
		)
	}

	presentationXML, err := e.renderPresentationXMLWithSections()
	if err != nil {
		return nil, err
	}
	out[common.PresentationXMLPath] = []byte(presentationXML)

	hasSections := len(e.sections) > 0
	hasNotesMaster := e.parts.Has("ppt/notesMasters/notesMaster1.xml")
	vbaProject, hasVBA := editorslide.VbaProjectFromMetadata(e.metadata.VBA)

	presentationRelsXML, err := renderPresentationRelsXML(e.nonSlideRels, e.slides, hasSections, hasVBA)
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

	hasHandoutMaster := e.parts.Has("ppt/handoutMasters/handoutMaster1.xml")
	if err := e.writeContentTypesPart(
		out,
		hasSections,
		hasNotesMaster,
		hasCommentAuthors,
		hasVBA,
		hasHandoutMaster,
		customXMLPropsPaths,
	); err != nil {
		return nil, err
	}

	e.writeOptionalPresentationParts(
		out,
		hasSections,
		hasNotesMaster,
		hasHandoutMaster,
		hasVBA,
		vbaProject,
	)
	return out, nil
}

func (e *PresentationEditor) writeContentTypesPart(
	out map[string][]byte,
	hasSections bool,
	hasNotesMaster bool,
	hasCommentAuthors bool,
	hasVBA bool,
	hasHandoutMaster bool,
	customXMLPropsPaths []string,
) error {
	mediaPaths := editorslide.MapValues(e.mediaInventory)
	filteredChartPaths := editorslide.FilterXMLPartPaths(e.parts.KeysWithPrefix("ppt/charts/chart"))
	notesPaths := editorslide.MapValues(e.notesInventory)

	themePaths := e.parts.KeysWithPrefix("ppt/theme/theme")
	layoutPaths := e.parts.KeysWithPrefix("ppt/slideLayouts/slideLayout")
	masterPaths := e.parts.KeysWithPrefix("ppt/slideMasters/slideMaster")

	filteredCommentPaths := editorslide.FilterXMLPartPaths(e.parts.KeysWithPrefix("ppt/comments/comment"))

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
		hasVBA,
		hasHandoutMaster,
		customXMLPropsPaths,
	)
	if err != nil {
		return err
	}
	out[common.ContentTypesPath] = []byte(contentTypesXML)
	return nil
}

func (e *PresentationEditor) writeOptionalPresentationParts(
	out map[string][]byte,
	hasSections bool,
	hasNotesMaster bool,
	hasHandoutMaster bool,
	hasVBA bool,
	vbaProject *vba.VBAProject,
) {
	if hasSections {
		out["ppt/sectionList.xml"] = []byte(buildSectionListXML(e.sections))
	}

	if hasNotesMaster {
		// Ensure notes master rels are also persisted if they were injected
		if masterRels, ok := e.parts.Get("ppt/notesMasters/_rels/notesMaster1.xml.rels"); ok {
			out["ppt/notesMasters/_rels/notesMaster1.xml.rels"] = masterRels
		}
	}
	if hasHandoutMaster {
		if masterXML, ok := e.parts.Get("ppt/handoutMasters/handoutMaster1.xml"); ok {
			out["ppt/handoutMasters/handoutMaster1.xml"] = masterXML
		}
		if masterRels, ok := e.parts.Get("ppt/handoutMasters/_rels/handoutMaster1.xml.rels"); ok {
			out["ppt/handoutMasters/_rels/handoutMaster1.xml.rels"] = masterRels
		}
	}

	if hasVBA {
		out["ppt/vbaProject.bin"] = vbaProject.Data
	}
}

func (e *PresentationEditor) processCustomXMLParts(out map[string][]byte) ([]string, error) {
	customXMLPropsPaths := make([]string, 0, len(e.metadata.CustomXML))
	for i, cXML := range e.metadata.CustomXML {
		itemPath, propsPath, err := e.writeCustomXMLPart(out, i, cXML)
		if err != nil {
			return nil, err
		}
		customXMLPropsPaths = append(customXMLPropsPaths, propsPath)
		e.ensureCustomXMLRelationship(itemPath)
		e.writeCustomXMLPropsRelationship(out, i+1)
	}
	return customXMLPropsPaths, nil
}

func (e *PresentationEditor) writeCustomXMLPart(
	out map[string][]byte,
	index int,
	part common.CustomXMLPart,
) (string, string, error) {
	itemContent, err := buildCustomXMLItemContent(part, index+1)
	if err != nil {
		return "", "", err
	}
	itemPath := fmt.Sprintf("customXml/item%d.xml", index+1)
	propsPath := fmt.Sprintf("customXml/itemProps%d.xml", index+1)
	out[itemPath] = []byte(itemContent)
	propsContent, err := buildCustomXMLPropsContent(part)
	if err != nil {
		return "", "", err
	}
	out[propsPath] = []byte(propsContent)
	return itemPath, propsPath, nil
}

func buildCustomXMLItemContent(part common.CustomXMLPart, index int) (string, error) {
	if part.RootElement == "" {
		return part.Content, nil
	}
	itemContent, err := editorslide.GenerateCustomXMLItem(part)
	if err != nil {
		return "", fmt.Errorf("custom XML part %d: %w", index, err)
	}
	return itemContent, nil
}

func buildCustomXMLPropsContent(part common.CustomXMLPart) (string, error) {
	schemaRefs := "<ds:schemaRefs></ds:schemaRefs>"
	if part.Namespace != "" {
		schemaRefs = fmt.Sprintf(
			`<ds:schemaRefs><ds:schemaRef ds:uri="%s"/></ds:schemaRefs>`,
			editorslide.EscapeCustomXML(part.Namespace),
		)
	}
	itemID := part.ItemID
	if itemID == "" {
		guid, err := common.NewGUID()
		if err != nil {
			return "", err
		}
		itemID = guid
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ds:datastoreItem ds:itemID="%s" xmlns:ds="http://schemas.openxmlformats.org/officeDocument/2006/customXml">
%s
</ds:datastoreItem>`, itemID, schemaRefs), nil
}

func (e *PresentationEditor) ensureCustomXMLRelationship(itemPath string) {
	itemTarget := "../" + itemPath
	for _, r := range e.nonSlideRels {
		if r.Type == common.RelTypeCustomXML && r.Target == itemTarget {
			return
		}
	}
	e.nonSlideRels = append(e.nonSlideRels, common.EditorRelationship{
		ID:     fmt.Sprintf("rId%d", e.nextRelIDNum),
		Type:   common.RelTypeCustomXML,
		Target: itemTarget,
	})
	e.nextRelIDNum++
}

func (e *PresentationEditor) writeCustomXMLPropsRelationship(out map[string][]byte, index int) {
	itemRelContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXmlProps" Target="itemProps%d.xml"/>
</Relationships>`, index)
	out[fmt.Sprintf("customXml/_rels/item%d.xml.rels", index)] = []byte(itemRelContent)
}

func (e *PresentationEditor) filterRootCustomXMLRelationships(out map[string][]byte) {
	packageRelsData, ok := e.parts.Get("_rels/.rels")
	if !ok {
		return
	}
	rels, err := parseRelationshipsXML(packageRelsData)
	if err != nil {
		return
	}
	filtered := make([]common.EditorRelationship, 0, len(rels))
	changed := false
	for _, r := range rels {
		if r.Type == common.RelTypeCustomXML || r.Type == common.RelTypeCustomXMLProps {
			changed = true
			continue
		}
		filtered = append(filtered, r)
	}
	if changed {
		out["_rels/.rels"] = []byte(renderRelationshipsXML(filtered))
	}
}

func (e *PresentationEditor) renderPresentationXMLWithSections() (string, error) {
	presentationXML, err := rewritePresentationSlideList([]byte(e.presentationXML), e.slides)
	if err != nil {
		return "", err
	}

	hasNotesMaster := e.parts.Has("ppt/notesMasters/notesMaster1.xml")
	notesMasterRelID, err := editorslide.ResolveNotesMasterRelID(
		e.nonSlideRels,
		hasNotesMaster,
		common.RelTypeNotesMaster,
	)
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

	presentationXML, err = rewritePresentationEmbeddedFonts([]byte(presentationXML), e.embeddedFontLst)
	if err != nil {
		return "", fmt.Errorf("rewrite embedded fonts: %w", err)
	}

	return presentationXML, nil
}
