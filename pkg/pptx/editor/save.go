package editor

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

const commentAuthorsRelType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/commentAuthors"

var editorCustomXMLNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9._-]*$`)

//nolint:gochecknoglobals // Reused escaper table avoids repeated allocations while serializing custom XML values.
var editorCustomXMLEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&apos;",
)

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
		if saveZipMethod(name) == zip.Store {
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

//nolint:gocognit,funlen // Part collection intentionally aggregates many conditional package part rewrites.
func (e *PresentationEditor) collectUpdatedParts() (map[string][]byte, error) {
	out := make(map[string][]byte)

	if err := e.serializeAuthorCacheIfPopulated(); err != nil {
		return nil, err
	}

	// 1. Process Custom XML and inject relationships Early.
	var customXMLPropsPaths []string

	for i, cXML := range e.metadata.CustomXML {
		var itemStr string
		var err error
		if cXML.RootElement != "" {
			itemStr, err = generateCustomXMLItem(cXML)
			if err != nil {
				return nil, fmt.Errorf("custom XML part %d: %w", i+1, err)
			}
		} else {
			itemStr = cXML.Content
		}
		itemPath := fmt.Sprintf("customXml/item%d.xml", i+1)
		out[itemPath] = []byte(itemStr)

		schemaRefs := "<ds:schemaRefs></ds:schemaRefs>"
		if cXML.Namespace != "" {
			schemaRefs = fmt.Sprintf(
				`<ds:schemaRefs><ds:schemaRef ds:uri="%s"/></ds:schemaRefs>`,
				escapeEditorCustomXML(cXML.Namespace),
			)
		}

		itemID := cXML.ItemID
		if itemID == "" {
			guid, err := common.NewGUID()
			if err != nil {
				return nil, err
			}
			itemID = guid
		}

		propsContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ds:datastoreItem ds:itemID="%s" xmlns:ds="http://schemas.openxmlformats.org/officeDocument/2006/customXml">
%s
</ds:datastoreItem>`, itemID, schemaRefs)
		propsPath := fmt.Sprintf("customXml/itemProps%d.xml", i+1)
		out[propsPath] = []byte(propsContent)
		customXMLPropsPaths = append(customXMLPropsPaths, propsPath)

		// Injection into presentation.xml.rels
		itemTarget := "../" + itemPath
		foundItemRel := false
		for _, r := range e.nonSlideRels {
			if r.Type == common.RelTypeCustomXML && r.Target == itemTarget {
				foundItemRel = true
				break
			}
		}
		if !foundItemRel {
			e.nonSlideRels = append(e.nonSlideRels, common.EditorRelationship{
				ID:     fmt.Sprintf("rId%d", e.nextRelIDNum),
				Type:   common.RelTypeCustomXML,
				Target: itemTarget,
			})
			e.nextRelIDNum++
		}

		// Create itemN.xml.rels
		itemRelContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXmlProps" Target="itemProps%d.xml"/>
</Relationships>`, i+1)
		out[fmt.Sprintf("customXml/_rels/item%d.xml.rels", i+1)] = []byte(itemRelContent)
	}

	// 2. Filter root package relationships (_rels/.rels) to remove any misplaced CustomXML rels.
	packageRelsData, ok := e.parts.Get("_rels/.rels")
	//nolint:nestif // Relationship parsing + filtering keeps error/condition handling local to this mutation block.
	if ok {
		rels, err := parseRelationshipsXML(packageRelsData)
		if err == nil {
			filtered := make([]common.EditorRelationship, 0, len(rels))
			changed := false
			for _, r := range rels {
				// CustomXML rels belong in presentation.xml.rels in PowerPoint, not at the root.
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
	vbaProject, hasVBA := vbaProjectFromMetadata(e.metadata.VBA)

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

	mediaPaths := mapValues(e.mediaInventory)
	filteredChartPaths := filterXMLPartPaths(e.parts.KeysWithPrefix("ppt/charts/chart"))
	notesPaths := mapValues(e.notesInventory)

	themePaths := e.parts.KeysWithPrefix("ppt/theme/theme")
	layoutPaths := e.parts.KeysWithPrefix("ppt/slideLayouts/slideLayout")
	masterPaths := e.parts.KeysWithPrefix("ppt/slideMasters/slideMaster")

	filteredCommentPaths := filterXMLPartPaths(e.parts.KeysWithPrefix("ppt/comments/comment"))

	contentTypesData, _ := e.parts.Get(common.ContentTypesPath)
	hasHandoutMaster := e.parts.Has("ppt/handoutMasters/handoutMaster1.xml")
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

	return out, nil
}

func generateCustomXMLItem(part common.CustomXMLPart) (string, error) {
	if !editorCustomXMLNamePattern.MatchString(part.RootElement) {
		return "", fmt.Errorf("invalid root element name %q", part.RootElement)
	}

	nsAttr := ""
	if part.Namespace != "" {
		nsAttr = fmt.Sprintf(` xmlns="%s"`, escapeEditorCustomXML(part.Namespace))
	}

	inner := part.Content
	if inner == "" && len(part.Properties) > 0 {
		var propsSb strings.Builder
		for j, kv := range part.Properties {
			if !editorCustomXMLNamePattern.MatchString(kv.Key) {
				return "", fmt.Errorf("invalid property element name %q", kv.Key)
			}
			if j > 0 {
				propsSb.WriteString("\n  ")
			}
			propsSb.WriteString(fmt.Sprintf("<%s>%s</%s>", kv.Key, escapeEditorCustomXML(kv.Value), kv.Key))
		}
		inner = propsSb.String()
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<%s%s>
  %s
</%s>`, part.RootElement, nsAttr, inner, part.RootElement), nil
}

func escapeEditorCustomXML(value string) string {
	return editorCustomXMLEscaper.Replace(value)
}

func vbaProjectFromMetadata(raw any) (*vba.VBAProject, bool) {
	project, ok := raw.(*vba.VBAProject)
	if !ok || project == nil {
		return nil, false
	}
	return project, project.IsMacroEnabled()
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

	presentationXML, err = rewritePresentationEmbeddedFonts([]byte(presentationXML), e.embeddedFontLst)
	if err != nil {
		return "", fmt.Errorf("rewrite embedded fonts: %w", err)
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

func saveZipMethod(path string) uint16 {
	lowerPath := strings.ToLower(path)
	if strings.HasPrefix(lowerPath, "ppt/notes") {
		return zip.Store
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tif", ".tiff", ".mp3", ".m4a", ".wav", ".mp4", ".avi":
		return zip.Store
	}
	return zip.Deflate
}
