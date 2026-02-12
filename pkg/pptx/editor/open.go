package editor

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type parsedSlideIDRef struct {
	SlideID int64
	RelID   string
}

// OpenPresentationEditor opens a PPTX package for in-place slide editing.
func OpenPresentationEditor(path string) (*PresentationEditor, error) {
	parts, err := loadPackageParts(path)
	if err != nil {
		return nil, err
	}
	return newPresentationEditorFromParts(parts)
}

func newPresentationEditorFromParts(parts map[string][]byte) (*PresentationEditor, error) {
	if _, err := requirePart(parts, common.ContentTypesPath); err != nil {
		return nil, err
	}
	presentationXMLBytes, err := requirePart(parts, common.PresentationXMLPath)
	if err != nil {
		return nil, err
	}
	presentationRelsBytes, err := requirePart(parts, common.PresentationRelPath)
	if err != nil {
		return nil, err
	}

	rels, err := parseRelationshipsXML(presentationRelsBytes)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", common.PresentationRelPath, err)
	}
	slideIDRefs, err := parsePresentationSlideIDs(presentationXMLBytes)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", common.PresentationXMLPath, err)
	}
	slideRefs, nonSlideRels, err := resolveSlideReferences(slideIDRefs, rels, parts)
	if err != nil {
		return nil, err
	}

	editor := &PresentationEditor{
		parts:           cloneParts(parts),
		slides:          slideRefs,
		nonSlideRels:    nonSlideRels,
		presentationXML: string(presentationXMLBytes),
	}
	slideSize, err := parsePresentationSlideSize(presentationXMLBytes)
	if err != nil {
		return nil, fmt.Errorf("parse %s slide size: %w", common.PresentationXMLPath, err)
	}
	editor.metadata = common.PresentationMetadata{
		Title:      extractCoreTitle(parts[common.CorePropsPath]),
		SlideCount: len(slideRefs),
		SlideSize:  slideSize,
	}
	editor.nextSlideID = nextSlideID(slideRefs)
	editor.nextRelIDNum = nextRelationshipNumber(rels)
	editor.nextSlideNum = nextSlidePartNumber(slideRefs)

	editor.mediaInventory, editor.nextMediaNum = parseMediaInventory(parts)
	if sectionData, ok := parts["ppt/sectionList.xml"]; ok {
		sections, _ := parseSectionListXML(sectionData)
		editor.sections = sections
	}

	editor.populateSlideTitlesConcurrently()
	return editor, nil
}

func loadPackageParts(filePath string) (map[string][]byte, error) {
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
	defer func() { _ = file.Close() }()

	zr, err := zip.NewReader(file, meta.Size())
	if err != nil {
		return nil, fmt.Errorf("invalid PPTX zip archive: %w", err)
	}

	parts := make(map[string][]byte, len(zr.File))
	for _, entry := range zr.File {
		if entry.FileInfo().IsDir() {
			continue
		}
		reader, err := entry.Open()
		if err != nil {
			return nil, fmt.Errorf("open zip entry %q: %w", entry.Name, err)
		}
		data, err := io.ReadAll(reader)
		_ = reader.Close()
		if err != nil {
			return nil, fmt.Errorf("read zip entry %q: %w", entry.Name, err)
		}
		parts[common.CanonicalPartPath(entry.Name)] = data
	}
	return parts, nil
}

func resolveSlideReferences(
	slideIDs []parsedSlideIDRef,
	rels []common.EditorRelationship,
	parts map[string][]byte,
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
		target := normalizePresentationTarget(rel.Target)
		partName := common.CanonicalPartPath(path.Join("ppt", target))
		if _, ok := parts[partName]; !ok {
			return nil, nil, fmt.Errorf("slide part %q not found", partName)
		}
		if err := editorEnsureSlideRelsExist(parts, partName); err != nil {
			return nil, nil, err
		}
		out = append(out, common.EditorSlideRef{
			SlideID: item.SlideID,
			RelID:   rel.ID,
			Target:  target,
			Part:    partName,
		})
	}
	return out, nonSlide, nil
}

func normalizePresentationTarget(target string) string {
	clean := strings.TrimSpace(strings.ReplaceAll(target, "\\", "/"))
	clean = strings.TrimPrefix(clean, "/")
	clean = strings.TrimPrefix(clean, "ppt/")
	return path.Clean(clean)
}

func parseRelationshipsXML(content []byte) ([]common.EditorRelationship, error) {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	out := make([]common.EditorRelationship, 0, 8)

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "Relationship" {
			continue
		}

		rel := common.EditorRelationship{}
		for _, attr := range start.Attr {
			switch attr.Name.Local {
			case "Id":
				rel.ID = strings.TrimSpace(attr.Value)
			case "Type":
				rel.Type = strings.TrimSpace(attr.Value)
			case "Target":
				rel.Target = strings.TrimSpace(attr.Value)
			case "TargetMode":
				rel.TargetMode = strings.TrimSpace(attr.Value)
			}
		}
		if rel.ID == "" || rel.Type == "" || rel.Target == "" {
			return nil, fmt.Errorf("relationship with missing Id/Type/Target")
		}
		out = append(out, rel)
	}
	return out, nil
}

func parsePresentationSlideIDs(content []byte) ([]parsedSlideIDRef, error) {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	out := make([]parsedSlideIDRef, 0, 8)

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "sldId" {
			continue
		}

		ref := parsedSlideIDRef{}
		for _, attr := range start.Attr {
			if attr.Name.Local != "id" {
				continue
			}
			if attr.Name.Space == "" {
				slideID, err := strconv.ParseInt(strings.TrimSpace(attr.Value), 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid slide id %q", attr.Value)
				}
				ref.SlideID = slideID
				continue
			}
			ref.RelID = strings.TrimSpace(attr.Value)
		}
		if ref.SlideID == 0 || ref.RelID == "" {
			return nil, fmt.Errorf("slide id entry missing id or r:id")
		}
		out = append(out, ref)
	}
	return out, nil
}

func parseMediaInventory(parts map[string][]byte) (map[string]string, int) {
	inventory := make(map[string]string)
	maxNum := 0
	for partPath, data := range parts {
		if strings.HasPrefix(partPath, "ppt/media/image") {
			hash := sha1.Sum(data)
			inventory[hex.EncodeToString(hash[:])] = partPath

			num, ok := parseImagePartNumber(partPath)
			if ok && num > maxNum {
				maxNum = num
			}
		}
	}
	return inventory, maxNum + 1
}

func parseImagePartNumber(partPath string) (int, bool) {
	base := path.Base(partPath)
	if !strings.HasPrefix(base, "image") {
		return 0, false
	}
	ext := path.Ext(base)
	name := strings.TrimSuffix(base, ext)
	numStr := strings.TrimPrefix(name, "image")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, false
	}
	return num, true
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

func parseSectionListXML(data []byte) ([]EditorSection, error) {
	var list xmlSectionList
	if err := xml.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	out := make([]EditorSection, 0, len(list.Sections))
	for _, s := range list.Sections {
		ids := make([]int64, 0, len(s.SlideIDs))
		for _, item := range s.SlideIDs {
			ids = append(ids, item.ID)
		}
		out = append(out, EditorSection{
			Name:     s.Name,
			GUID:     s.GUID,
			SlideIDs: ids,
		})
	}
	return out, nil
}
