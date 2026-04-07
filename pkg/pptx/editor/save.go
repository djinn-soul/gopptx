package editor

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation/protection"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

const commentAuthorsRelType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/commentAuthors"
const rawZipCopyBufferSize = 32 * 1024

// Save writes the edited presentation back to a PPTX file.
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

	allNames := mergedPartNames(e.parts.Keys(), updatedParts)
	if password := strings.TrimSpace(e.metadata.Protection.EncryptPassword); password != "" {
		output, err := e.buildZipStream(allNames, updatedParts)
		if err != nil {
			return err
		}
		encrypted, err := protection.EncryptAgilePackage(output, password)
		if err != nil {
			return fmt.Errorf("encrypt presentation package: %w", err)
		}
		if err := os.WriteFile(filePath, encrypted, 0o600); err != nil {
			return fmt.Errorf("write %s: %w", filePath, err)
		}
		return nil
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open %s: %w", filePath, err)
	}
	defer func() { _ = file.Close() }()

	if err := e.buildZipToWriter(file, allNames, updatedParts); err != nil {
		return err
	}
	return nil
}

// SaveToBytes serializes the presentation to a byte slice without writing to disk.
func (e *PresentationEditor) SaveToBytes() ([]byte, error) {
	if e == nil {
		return nil, errors.New("nil editor")
	}
	vbaProject, hasVBA := editorslide.VbaProjectFromMetadata(e.metadata.VBA)

	updatedParts, err := e.collectUpdatedParts(vbaProject, hasVBA)
	if err != nil {
		return nil, fmt.Errorf("prepare updated parts: %w", err)
	}

	allNames := mergedPartNames(e.parts.Keys(), updatedParts)

	output, err := e.buildZipStream(allNames, updatedParts)
	if err != nil {
		return nil, err
	}

	if password := strings.TrimSpace(e.metadata.Protection.EncryptPassword); password != "" {
		encrypted, encErr := protection.EncryptAgilePackage(output, password)
		if encErr != nil {
			return nil, fmt.Errorf("encrypt presentation package: %w", encErr)
		}
		output = encrypted
	}
	return output, nil
}

// SaveToWriter serializes the presentation and writes it to the provided io.Writer.
func (e *PresentationEditor) SaveToWriter(w io.Writer) error {
	if e == nil {
		return errors.New("nil editor")
	}
	vbaProject, hasVBA := editorslide.VbaProjectFromMetadata(e.metadata.VBA)

	updatedParts, err := e.collectUpdatedParts(vbaProject, hasVBA)
	if err != nil {
		return fmt.Errorf("prepare updated parts: %w", err)
	}
	allNames := mergedPartNames(e.parts.Keys(), updatedParts)

	password := strings.TrimSpace(e.metadata.Protection.EncryptPassword)
	if password == "" {
		return e.buildZipToWriter(w, allNames, updatedParts)
	}

	output, err := e.buildZipStream(allNames, updatedParts)
	if err != nil {
		return err
	}
	encrypted, err := protection.EncryptAgilePackage(output, password)
	if err != nil {
		return fmt.Errorf("encrypt presentation package: %w", err)
	}
	_, err = w.Write(encrypted)
	return err
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
