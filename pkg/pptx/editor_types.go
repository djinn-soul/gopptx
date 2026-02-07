package pptx

import (
	"fmt"
	"sort"
	"strings"
)

const (
	presentationRelPath = "ppt/_rels/presentation.xml.rels"
	presentationXMLPath = "ppt/presentation.xml"
	contentTypesPath    = "[Content_Types].xml"
	corePropsPath       = "docProps/core.xml"

	relTypeSlide       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide"
	relTypeSlideLayout = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout"
	relTypeNotesSlide  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide"
)

// PresentationMetadata describes summary information parsed from an existing PPTX package.
type PresentationMetadata struct {
	Title      string
	SlideCount int
}

// SlideMetadata describes one slide entry inside an editable presentation.
type SlideMetadata struct {
	Index          int
	SlideID        int64
	RelationshipID string
	PartName       string
	Title          string
}

type editorSlideRef struct {
	SlideID int64
	RelID   string
	Target  string
	Part    string
	Title   string
}

type editorRelationship struct {
	ID         string
	Type       string
	Target     string
	TargetMode string
}

// PresentationEditor provides read/modify/save operations for existing PPTX files.
//
// The current editor supports slide operations for content that does not require
// external embedded assets during add/update operations (for example: images/charts/media).
// Unsupported cases fail fast with explicit errors.
type PresentationEditor struct {
	parts map[string][]byte

	slides       []editorSlideRef
	nextSlideID  int64
	nextRelIDNum int
	nextSlideNum int

	metadata        PresentationMetadata
	nonSlideRels    []editorRelationship
	presentationXML string
}

// Metadata returns presentation-level metadata parsed from the package.
func (e *PresentationEditor) Metadata() PresentationMetadata {
	return e.metadata
}

// SlideCount returns the number of slides currently tracked by the editor.
func (e *PresentationEditor) SlideCount() int {
	if e == nil {
		return 0
	}
	return len(e.slides)
}

// Slides returns ordered slide metadata snapshots (0-based indexes).
func (e *PresentationEditor) Slides() []SlideMetadata {
	if e == nil || len(e.slides) == 0 {
		return nil
	}
	out := make([]SlideMetadata, 0, len(e.slides))
	for idx, slide := range e.slides {
		out = append(out, SlideMetadata{
			Index:          idx,
			SlideID:        slide.SlideID,
			RelationshipID: slide.RelID,
			PartName:       slide.Part,
			Title:          slide.Title,
		})
	}
	return out
}

func cloneParts(parts map[string][]byte) map[string][]byte {
	out := make(map[string][]byte, len(parts))
	for path, content := range parts {
		clone := make([]byte, len(content))
		copy(clone, content)
		out[path] = clone
	}
	return out
}

func requirePart(parts map[string][]byte, path string) ([]byte, error) {
	content, ok := parts[path]
	if !ok {
		return nil, fmt.Errorf("missing required package part %q", path)
	}
	return content, nil
}

func partNamesSorted(parts map[string][]byte) []string {
	names := make([]string, 0, len(parts))
	for path := range parts {
		names = append(names, path)
	}
	sort.Strings(names)
	return names
}

func canonicalPartPath(target string) string {
	clean := strings.TrimSpace(strings.ReplaceAll(target, "\\", "/"))
	if strings.HasPrefix(clean, "/") {
		return strings.TrimPrefix(clean, "/")
	}
	return clean
}

func slideRelsPartName(slidePart string) string {
	clean := strings.TrimSpace(strings.ReplaceAll(slidePart, "\\", "/"))
	if clean == "" {
		return ""
	}
	lastSlash := strings.LastIndex(clean, "/")
	if lastSlash < 0 {
		return "_rels/" + clean + ".rels"
	}
	dir := clean[:lastSlash]
	file := clean[lastSlash+1:]
	return dir + "/_rels/" + file + ".rels"
}
