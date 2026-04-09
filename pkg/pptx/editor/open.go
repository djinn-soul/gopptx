package editor

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

const (
	legacyPPTHeaderLength = 8
)

// OpenPresentationEditor opens a PPTX package for in-place slide editing.
func OpenPresentationEditor(filePath string) (*PresentationEditor, error) {
	ps, err := openPartStore(filePath)
	if err != nil {
		return nil, err
	}
	editor, err := NewPresentationEditorFromParts(ps)
	if err != nil {
		_ = ps.Close()
		return nil, err
	}
	return editor, nil
}

// OpenPresentationEditorFromBytes opens a PPTX package from a byte slice for in-place slide editing.
func OpenPresentationEditorFromBytes(data []byte) (*PresentationEditor, error) {
	ps, err := OpenPartStoreFromBytes(data)
	if err != nil {
		return nil, err
	}
	editor, err := NewPresentationEditorFromParts(ps)
	if err != nil {
		_ = ps.Close()
		return nil, err
	}
	return editor, nil
}

// OpenPresentationEditorFromReader opens a PPTX package from an io.Reader for in-place slide editing.
// The entire reader is read into memory before processing.
func OpenPresentationEditorFromReader(r io.Reader) (*PresentationEditor, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read presentation: %w", err)
	}
	return OpenPresentationEditorFromBytes(data)
}

// OpenPartStoreFromBytes opens a part store from a byte slice.
func OpenPartStoreFromBytes(data []byte) (*PartStore, error) {
	if len(data) >= 8 && bytes.Equal(data[:8], []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}) {
		return nil, errors.New(
			"legacy proprietary .ppt (OLE2) files are not supported. please use the interop package to convert to .pptx first",
		)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("not a valid ZIP archive: %w", err)
	}
	// Note: nil os.File is fine for in-memory or byte-backed PartStore as long as zip.Reader is used for lazy reads.
	return newPartStoreFromZip(nil, zr), nil
}

// NewPresentationEditorFromParts creates a new presentation editor from an existing part store.
func NewPresentationEditorFromParts(ps *PartStore) (*PresentationEditor, error) {
	if !ps.Has(common.ContentTypesPath) {
		return nil, fmt.Errorf("missing required package part %q", common.ContentTypesPath)
	}
	presentationXMLBytes, ok := ps.Get(common.PresentationXMLPath)
	if !ok {
		return nil, fmt.Errorf("missing required package part %q", common.PresentationXMLPath)
	}
	presentationRelsBytes, ok := ps.Get(common.PresentationRelPath)
	if !ok {
		return nil, fmt.Errorf("missing required package part %q", common.PresentationRelPath)
	}

	rels, err := parseRelationshipsXML(presentationRelsBytes)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", common.PresentationRelPath, err)
	}
	slideIDRefs, err := editorslide.ParsePresentationSlideIDs(presentationXMLBytes)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", common.PresentationXMLPath, err)
	}
	slideRefs, nonSlideRels, err := resolveSlideReferences(slideIDRefs, rels, ps)
	if err != nil {
		return nil, err
	}

	editor := &PresentationEditor{
		parts:           ps,
		slides:          slideRefs,
		nonSlideRels:    nonSlideRels,
		presentationXML: string(presentationXMLBytes),
		embeddedFontLst: extractEmbeddedFontLst(presentationXMLBytes),
		imagePathCache:  make(map[string]imagePathCacheEntry),
	}
	slideSize, err := parsePresentationSlideSize(presentationXMLBytes)
	if err != nil {
		return nil, fmt.Errorf("parse %s slide size: %w", common.PresentationXMLPath, err)
	}
	coreData, _ := ps.Get(common.CorePropsPath)
	coreProps, _ := parseCoreProperties(coreData)
	editor.metadata = common.Metadata{}
	editor.metadata.Title = coreProps.Title
	editor.metadata.SlideCount = len(slideRefs)
	editor.metadata.SlideSize = slideSize
	editor.metadata.CoreProperties = coreProps

	if vbaData, hasVBA := ps.Get(vba.PackagePath); hasVBA {
		editor.metadata.VBA = vba.FromData(vbaData)
	}

	editor.nextSlideID = editorslide.NextSlideID(slideRefs)
	editor.nextRelIDNum = common.NextRelationshipNumber(rels)
	editor.nextSlideNum = editorslide.NextSlidePartNumber(slideRefs)

	partKeys := ps.Keys()
	editor.mediaInventory, editor.nextMediaNum = editorslide.ParseMediaInventory(ps, partKeys)
	if sectionData, sectionOK := ps.Get("ppt/sectionList.xml"); sectionOK {
		sections, err := parseSectionListXML(sectionData)
		if err != nil {
			return nil, fmt.Errorf("parse ppt/sectionList.xml: %w", err)
		}
		editor.sections = sections
	}

	editor.chartEmbeddings, editor.nextChartNum, editor.nextExcelNum = editorslide.ParseChartInventory(ps, partKeys)
	editor.notesInventory, editor.nextNotesNum = editorslide.ParseNotesInventory(ps, partKeys)
	editor.nextDiagramNum = editorslide.ParseDiagramInventory(partKeys)

	editor.metadata.CustomXML = editorslide.ParseCustomXMLInventory(ps, partKeys)

	editor.populateSlideTitlesConcurrently()
	return editor, nil
}

// newPresentationEditorFromParts is kept as a package-local compatibility shim
// for existing internal tests and call sites.
func newPresentationEditorFromParts(ps *PartStore) (*PresentationEditor, error) {
	return NewPresentationEditorFromParts(ps)
}

func openPartStore(filePath string) (*PartStore, error) {
	meta, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	if !meta.Mode().IsRegular() {
		return nil, fmt.Errorf("path is not a regular file: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	// Check for legacy .ppt OLE2 magic number
	header := make([]byte, legacyPPTHeaderLength)
	if n, err := file.Read(header); err == nil && n == legacyPPTHeaderLength {
		if bytes.Equal(header, []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}) {
			_ = file.Close()
			return nil, errors.New(
				"legacy proprietary .ppt (OLE2) files are not supported. please use the interop package to convert to .pptx first",
			)
		}
	}
	// Reset file pointer after reading header
	if _, err := file.Seek(0, 0); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("failed to seek file: %w", err)
	}

	zr, err := zip.NewReader(file, meta.Size())
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("not a valid ZIP archive: %w", err)
	}

	return newPartStoreFromZip(file, zr), nil
}

func resolveSlideReferences(
	slideIDs []editorslide.ParsedSlideIDRef,
	rels []common.EditorRelationship,
	ps *PartStore,
) ([]common.EditorSlideRef, []common.EditorRelationship, error) {
	relByID := make(map[string]common.EditorRelationship, len(rels))
	nonSlide := make([]common.EditorRelationship, 0, len(rels))
	for _, rel := range rels {
		relByID[rel.ID] = rel
		if rel.Type != common.RelTypeSlide {
			nonSlide = append(nonSlide, rel)
		}
	}

	if len(slideIDs) == 0 {
		return nil, nonSlide, nil
	}

	out := make([]common.EditorSlideRef, 0, len(slideIDs))
	for _, item := range slideIDs {
		rel, ok := relByID[item.RelID]
		if !ok {
			return nil, nil, fmt.Errorf("presentation.xml references missing relationship %q", item.RelID)
		}
		if rel.Type != common.RelTypeSlide {
			return nil, nil, fmt.Errorf("relationship %q is not a slide relationship", item.RelID)
		}
		target := editorslide.NormalizePresentationTarget(rel.Target)
		partName := common.CanonicalPartPath(path.Join("ppt", target))
		if !ps.Has(partName) {
			return nil, nil, fmt.Errorf("slide part %q not found", partName)
		}
		slideXML, _ := ps.Get(partName)
		hiddenFromSlideXML, err := editorslide.ParseSlideHidden(slideXML)
		if err != nil {
			return nil, nil, fmt.Errorf("parse %s hidden flag: %w", partName, err)
		}
		if err := editorslide.EnsureSlideRelsExist(ps.Has, partName); err != nil {
			return nil, nil, err
		}
		out = append(out, common.EditorSlideRef{
			SlideID: item.SlideID,
			RelID:   rel.ID,
			Target:  target,
			Part:    partName,
			// Keep legacy p:sldId show="0" support for older files while
			// preferring the schema-valid slide root show="0" marker.
			Hidden: item.Hidden || hiddenFromSlideXML,
		})
	}
	return out, nonSlide, nil
}

func extractEmbeddedFontLst(xml []byte) string {
	return editorslide.ExtractEmbeddedFontList(xml)
}

func parseRelationshipsXML(content []byte) ([]common.EditorRelationship, error) {
	return editorslide.ParseRelationshipsXML(content)
}

type xmlSectionList struct {
	Sections []xmlSection `xml:"section"`
}

type xmlSection struct {
	Name     string             `xml:"name,attr"`
	GUID     string             `xml:"id,attr"`
	SlideIDs []xmlSectionSLDRef `xml:"sldIdLst>sldId"`
}

type xmlSectionSLDRef struct {
	ID int64 `xml:"id,attr"`
}

func parseSectionListXML(data []byte) ([]Section, error) {
	var list xmlSectionList
	if err := xml.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	out := make([]Section, 0, len(list.Sections))
	for _, s := range list.Sections {
		ids := make([]int64, 0, len(s.SlideIDs))
		for _, item := range s.SlideIDs {
			ids = append(ids, item.ID)
		}
		out = append(out, Section{
			Name:     s.Name,
			GUID:     s.GUID,
			SlideIDs: ids,
		})
	}
	return out, nil
}
