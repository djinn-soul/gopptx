package editor

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type parsedSlideIDRef struct {
	SlideID int64
	RelID   string
}

const (
	defaultRelsCapacity     = 8
	defaultSlideIDsCapacity = 8
)

// OpenPresentationEditor opens a PPTX package for in-place slide editing.
func OpenPresentationEditor(filePath string) (*PresentationEditor, error) {
	ps, err := openPartStore(filePath)
	if err != nil {
		return nil, err
	}
	editor, err := newPresentationEditorFromParts(ps)
	if err != nil {
		_ = ps.Close()
		return nil, err
	}
	return editor, nil
}

func newPresentationEditorFromParts(ps *PartStore) (*PresentationEditor, error) {
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
	slideIDRefs, err := parsePresentationSlideIDs(presentationXMLBytes)
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
		imagePathCache:  make(map[string]imagePathCacheEntry),
	}
	slideSize, err := parsePresentationSlideSize(presentationXMLBytes)
	if err != nil {
		return nil, fmt.Errorf("parse %s slide size: %w", common.PresentationXMLPath, err)
	}
	coreData, _ := ps.Get(common.CorePropsPath)
	coreProps, _ := parseCoreProperties(coreData)
	editor.metadata = common.PresentationMetadata{
		Title:          coreProps.Title,
		SlideCount:     len(slideRefs),
		SlideSize:      slideSize,
		CoreProperties: coreProps,
	}
	editor.nextSlideID = nextSlideID(slideRefs)
	editor.nextRelIDNum = nextRelationshipNumber(rels)
	editor.nextSlideNum = nextSlidePartNumber(slideRefs)

	partKeys := ps.Keys()
	editor.mediaInventory, editor.nextMediaNum = parseMediaInventory(ps, partKeys)
	if sectionData, sectionOK := ps.Get("ppt/sectionList.xml"); sectionOK {
		sections, _ := parseSectionListXML(sectionData)
		editor.sections = sections
	}

	editor.chartEmbeddings, editor.nextChartNum, editor.nextExcelNum = parseChartInventory(ps, partKeys)
	editor.notesInventory, editor.nextNotesNum = parseNotesInventory(ps, partKeys)

	editor.populateSlideTitlesConcurrently()
	return editor, nil
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

	zr, err := zip.NewReader(file, meta.Size())
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("invalid PPTX zip archive: %w", err)
	}

	return newPartStoreFromZip(file, zr), nil
}

func resolveSlideReferences(
	slideIDs []parsedSlideIDRef,
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
		target := normalizePresentationTarget(rel.Target)
		partName := common.CanonicalPartPath(path.Join("ppt", target))
		if !ps.Has(partName) {
			return nil, nil, fmt.Errorf("slide part %q not found", partName)
		}
		if err := editorEnsureSlideRelsExistPS(ps, partName); err != nil {
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
	out := make([]common.EditorRelationship, 0, defaultRelsCapacity)

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
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
			return nil, errors.New("relationship with missing Id/Type/Target")
		}
		out = append(out, rel)
	}
	return out, nil
}

func parsePresentationSlideIDs(content []byte) ([]parsedSlideIDRef, error) {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	out := make([]parsedSlideIDRef, 0, defaultSlideIDsCapacity)

	for {
		start, ok, err := nextSlideIDStartElement(decoder)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		if !isPresentationSlideIDElement(start) {
			continue
		}

		ref, err := parseSlideIDRef(start)
		if err != nil {
			return nil, err
		}
		if ref.SlideID == 0 || ref.RelID == "" {
			return nil, errors.New("slide id entry missing id or r:id")
		}
		out = append(out, ref)
	}
	return out, nil
}

func nextSlideIDStartElement(decoder *xml.Decoder) (xml.StartElement, bool, error) {
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return xml.StartElement{}, false, nil
			}
			return xml.StartElement{}, false, err
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		return start, true, nil
	}
}

func isPresentationSlideIDElement(start xml.StartElement) bool {
	if start.Name.Local != "sldId" {
		return false
	}
	if start.Name.Space == "" {
		return true
	}
	return start.Name.Space == "http://schemas.openxmlformats.org/presentationml/2006/main"
}

func parseSlideIDRef(start xml.StartElement) (parsedSlideIDRef, error) {
	ref := parsedSlideIDRef{}
	for _, attr := range start.Attr {
		if attr.Name.Local != "id" {
			continue
		}
		if attr.Name.Space == "" {
			slideID, parseErr := strconv.ParseInt(strings.TrimSpace(attr.Value), 10, 64)
			if parseErr != nil {
				return parsedSlideIDRef{}, fmt.Errorf("invalid slide id %q", attr.Value)
			}
			ref.SlideID = slideID
			continue
		}
		ref.RelID = strings.TrimSpace(attr.Value)
	}
	return ref, nil
}

func parseMediaInventory(ps *PartStore, partKeys []string) (map[string]string, int) {
	inventory := make(map[string]string)
	maxNum := 0
	for _, partPath := range partKeys {
		if !strings.HasPrefix(partPath, "ppt/media/image") {
			continue
		}
		data, ok := ps.Get(partPath)
		if !ok {
			continue
		}
		hash := sha256.Sum256(data)
		inventory[hex.EncodeToString(hash[:])] = partPath

		num, ok := parseImagePartNumber(partPath)
		if ok && num > maxNum {
			maxNum = num
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

func parseChartInventory(ps *PartStore, partKeys []string) (map[string]string, int, int) {
	inventory := make(map[string]string)
	maxChart := 0
	maxExcel := 0

	for _, p := range partKeys {
		if !isChartPartPath(p) {
			continue
		}
		num, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(p, "ppt/charts/chart"), ".xml"))
		if num > maxChart {
			maxChart = num
		}

		excelPath, nextExcel := findChartEmbedding(ps, p, maxExcel)
		if excelPath == "" {
			continue
		}
		inventory[p] = excelPath
		maxExcel = nextExcel
	}
	return inventory, maxChart + 1, maxExcel + 1
}

func isChartPartPath(partPath string) bool {
	return strings.HasPrefix(partPath, "ppt/charts/chart") && strings.HasSuffix(partPath, ".xml")
}

func findChartEmbedding(ps *PartStore, chartPart string, currentMaxExcel int) (string, int) {
	relsPath := "ppt/charts/_rels/" + path.Base(chartPart) + ".rels"
	relsData, ok := ps.Get(relsPath)
	if !ok {
		return "", currentMaxExcel
	}

	rels, _ := parseRelationshipsXML(relsData)
	maxExcel := currentMaxExcel
	for _, rel := range rels {
		if rel.Type != common.RelTypePackage {
			continue
		}
		excelPath := common.CanonicalPartPath(path.Join("ppt/charts", rel.Target))
		maxExcel = maxExcelNumber(maxExcel, excelPath)
		return excelPath, maxExcel
	}
	return "", maxExcel
}

func maxExcelNumber(current int, excelPath string) int {
	base := path.Base(excelPath)
	after, ok := strings.CutPrefix(base, "Microsoft_Excel_Worksheet")
	if !ok {
		return current
	}

	enum, _ := strconv.Atoi(strings.TrimSuffix(after, ".xlsx"))
	if enum > current {
		return enum
	}
	return current
}

func parseNotesInventory(ps *PartStore, partKeys []string) (map[string]string, int) {
	inventory := make(map[string]string)
	maxNotes := 0

	for _, p := range partKeys {
		if !strings.HasPrefix(p, "ppt/slides/_rels/slide") {
			continue
		}
		if !strings.HasSuffix(p, ".xml.rels") {
			continue
		}
		slidePart := "ppt/slides/" + strings.TrimSuffix(path.Base(p), ".rels")
		relsData, ok := ps.Get(p)
		if !ok {
			continue
		}
		rels, _ := parseRelationshipsXML(relsData)
		for _, r := range rels {
			if r.Type == "http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide" {
				notesPath := common.CanonicalPartPath(path.Join("ppt/slides", r.Target))
				inventory[slidePart] = notesPath

				num, _ := strconv.Atoi(
					strings.TrimSuffix(strings.TrimPrefix(path.Base(notesPath), "notesSlide"), ".xml"),
				)
				if num > maxNotes {
					maxNotes = num
				}
			}
		}
	}
	return inventory, maxNotes + 1
}
