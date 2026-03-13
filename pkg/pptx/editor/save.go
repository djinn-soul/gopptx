package editor

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation/protection"
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

	var zipBuf bytes.Buffer
	zw := zip.NewWriter(&zipBuf)

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

	if err := zw.Close(); err != nil {
		return fmt.Errorf("finalize zip stream: %w", err)
	}

	output := zipBuf.Bytes()
	if password := strings.TrimSpace(e.metadata.Protection.EncryptPassword); password != "" {
		encrypted, err := protection.EncryptAgilePackage(output, password)
		if err != nil {
			return fmt.Errorf("encrypt presentation package: %w", err)
		}
		output = encrypted
	}

	if err := os.WriteFile(filePath, output, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", filePath, err)
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

	presentationRelsXML, err := editorslide.RenderPresentationRelsXML(e.nonSlideRels, e.slides, hasSections, hasVBA)
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
	contentTypesXML, err := editorslide.RewriteContentTypes(
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
		out["ppt/sectionList.xml"] = []byte(editorslide.BuildSectionListXML(e.sections))
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

func (e *PresentationEditor) renderPresentationXMLWithSections() (string, error) {
	presentationXML, err := editorslide.RewritePresentationSlideList([]byte(e.presentationXML), e.slides)
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

	presentationXML, err = editorslide.RewritePresentationNotesMasterList(
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

	presentationXML, err = editorslide.RewritePresentationSections([]byte(presentationXML), e.sections)
	if err != nil {
		return "", fmt.Errorf("rewrite sections: %w", err)
	}

	presentationXML, err = editorslide.RewritePresentationEmbeddedFonts([]byte(presentationXML), e.embeddedFontLst)
	if err != nil {
		return "", fmt.Errorf("rewrite embedded fonts: %w", err)
	}

	return presentationXML, nil
}
