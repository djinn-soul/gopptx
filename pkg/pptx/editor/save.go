package editor

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation/protection"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

const commentAuthorsRelType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/commentAuthors"
const manifestBuildWorkers = 4

// Save writes the edited presentation back to a PPTX file.
//
//nolint:gocognit,funlen // Save flow intentionally sequences materialize/validate/write/cleanup steps with explicit guards.
func (e *PresentationEditor) Save(filePath string) error {
	if e == nil {
		return errors.New("nil editor")
	}
	vbaProject, hasVBA := editorslide.VbaProjectFromMetadata(e.metadata.VBA)
	if hasVBA {
		ext := strings.ToLower(strings.TrimSpace(filepath.Ext(filePath)))
		if ext != ".pptm" {
			return fmt.Errorf("macro-enabled presentations must be saved with .pptm extension, got %q", ext)
		}
	}

	// Materialize all lazy parts into memory and release the source file handle.
	if err := e.parts.Materialize(); err != nil {
		return fmt.Errorf("failed to materialize lazy PPTX parts from source archive: %w", err)
	}

	updatedParts, err := e.collectUpdatedParts(vbaProject, hasVBA)
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

func (e *PresentationEditor) collectUpdatedParts(vbaProject *vba.VBAProject, hasVBA bool) (map[string][]byte, error) {
	out := make(map[string][]byte)

	if err := e.verifyMediaInventoryChecksumsParallel(); err != nil {
		return nil, fmt.Errorf("verify media inventory: %w", err)
	}

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

	hasSections := len(e.sections) > 0
	hasNotesMaster := e.parts.Has("ppt/notesMasters/notesMaster1.xml")
	hasHandoutMaster := e.parts.Has("ppt/handoutMasters/handoutMaster1.xml")
	manifestParts, err := e.buildManifestPartsParallel(
		hasSections,
		hasNotesMaster,
		hasCommentAuthors,
		hasVBA,
		hasHandoutMaster,
		customXMLPropsPaths,
	)
	if err != nil {
		return nil, err
	}
	out[common.PresentationXMLPath] = manifestParts.presentationXML
	out[common.PresentationRelPath] = manifestParts.presentationRelsXML
	out[common.CorePropsPath] = manifestParts.corePropsXML
	out[common.ContentTypesPath] = manifestParts.contentTypesXML

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
