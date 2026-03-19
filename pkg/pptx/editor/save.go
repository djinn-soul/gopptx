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
const rawZipCopyBufferSize = 32 * 1024

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

	// Saving in-place to the source path requires closing the source file handle first.
	// Keep the old materialize behavior for this path to preserve editor usability.
	if e.parts.IsBackedByPath(filePath) {
		if err := e.parts.Materialize(); err != nil {
			return fmt.Errorf("failed to materialize lazy PPTX parts from source archive: %w", err)
		}
	}

	updatedParts, err := e.collectUpdatedParts(vbaProject, hasVBA)
	if err != nil {
		return fmt.Errorf("prepare updated parts: %w", err)
	}

	keys := e.parts.Keys()

	var zipBuf bytes.Buffer
	zw := zip.NewWriter(&zipBuf)
	rawZipCopyBuffer := make([]byte, rawZipCopyBufferSize)

	// Merge existing sorted keys with updatedParts keys.
	// keys is already sorted; updatedParts typically contains only a handful of entries
	// (presentationXML, presentationRels, coreProps, contentTypes, …), nearly all of
	// which already exist in keys.  Binary-search for each updatedParts key avoids
	// building a full N-entry set map just to deduplicate a tiny delta.
	var extraKeys []string
	for k := range updatedParts {
		if i := sort.SearchStrings(keys, k); i >= len(keys) || keys[i] != k {
			extraKeys = append(extraKeys, k)
		}
	}
	var allNames []string
	if len(extraKeys) == 0 {
		// Common case: no new part names — reuse the sorted slice directly (zero allocs).
		allNames = keys
	} else {
		sort.Strings(extraKeys)
		allNames = make([]string, 0, len(keys)+len(extraKeys))
		ki, ei := 0, 0
		for ki < len(keys) && ei < len(extraKeys) {
			if keys[ki] <= extraKeys[ei] {
				allNames = append(allNames, keys[ki])
				ki++
			} else {
				allNames = append(allNames, extraKeys[ei])
				ei++
			}
		}
		allNames = append(allNames, keys[ki:]...)
		allNames = append(allNames, extraKeys[ei:]...)
	}

	for _, name := range allNames {
		if updated, updatedOK := updatedParts[name]; updatedOK {
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
			if _, writeErr := w.Write(updated); writeErr != nil {
				return fmt.Errorf("write zip entry %q: %w", name, writeErr)
			}
			continue
		}

		// Fast path: unchanged part from original archive can be copied raw.
		if sourceEntry, sourceOK := e.parts.SourceZipEntry(name); sourceOK {
			if err := copyZipEntryRaw(zw, sourceEntry, rawZipCopyBuffer); err != nil {
				return fmt.Errorf("copy source zip entry %q: %w", name, err)
			}
			continue
		}

		// Fallback for in-memory or already materialized parts.
		content, partOK := e.parts.Get(name)
		if !partOK {
			return fmt.Errorf("failed to retrieve part %q during save", name)
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

func copyZipEntryRaw(zw *zip.Writer, sourceEntry *zip.File, copyBuffer []byte) error {
	header := sourceEntry.FileHeader
	writer, err := zw.CreateRaw(&header)
	if err != nil {
		return err
	}
	reader, err := sourceEntry.OpenRaw()
	if err != nil {
		return err
	}
	if _, err := io.CopyBuffer(writer, reader, copyBuffer); err != nil {
		return err
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
	manifestParts, err := e.buildManifestParts(
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
