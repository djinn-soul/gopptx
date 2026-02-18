package editorcommon

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

// Shared constants for editor logic.
const (
	PresentationRelPath = "ppt/_rels/presentation.xml.rels"
	PresentationXMLPath = "ppt/presentation.xml"
	ContentTypesPath    = "[Content_Types].xml"
	CorePropsPath       = "docProps/core.xml"

	// DCNamespace and related metadata XML namespaces.
	DCNamespace       = "http://purl.org/dc/elements/1.1/"
	DCTermsNamespace  = "http://purl.org/dc/terms/"
	DCMITypeNamespace = "http://purl.org/dc/dcmitype/"
	CPNamespace       = "http://schemas.openxmlformats.org/package/2006/metadata/core-properties"
	XSINamespace      = "http://www.w3.org/2001/XMLSchema-instance"

	RelTypeSlide       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide"
	RelTypeSlideMaster = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster"
	RelTypeSlideLayout = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout"
	RelTypeNotesSlide  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide"
	RelTypeNotesMaster = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesMaster"
	RelTypeHyperlink   = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink"
	RelTypeImage       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"
	RelTypeChart       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart"
	RelTypeAudio       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/audio"
	RelTypeVideo       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/video"
	RelTypeTheme       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme"
	RelTypeSectionList = "http://schemas.microsoft.com/office/2007/relationships/sectionList"
	RelTypePackage     = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/package"

	RelationshipsXMLNS = "http://schemas.openxmlformats.org/package/2006/relationships"
	ContentTypesXMLNS  = "http://schemas.openxmlformats.org/package/2006/content-types"
	SlideContentType   = "application/vnd.openxmlformats-officedocument.presentationml.slide+xml"
)

// EditorRelationship describes an OOXML relationship entry.
type EditorRelationship struct {
	ID         string
	Type       string
	Target     string
	TargetMode string
}

// EditorSlideRef describes internal slide tracking data within the editor.
type EditorSlideRef struct {
	SlideID int64
	RelID   string
	Target  string
	Part    string
	Title   string
}

// XMLEscape provides basic XML attribute escaping.
func XMLEscape(value string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(value))
	return b.String()
}

// CanonicalPartPath cleans a path to use forward slashes and removes leading slash.
func CanonicalPartPath(target string) string {
	clean := strings.TrimSpace(strings.ReplaceAll(target, "\\", "/"))
	if after, ok := strings.CutPrefix(clean, "/"); ok {
		return after
	}
	return clean
}

// SlideRelsPartName returns the relative path to a slide's relationships part.
func SlideRelsPartName(slidePart string) string {
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

// RelsPathFor returns the relationships part path for any given part path.
func RelsPathFor(partPath string) string {
	clean := strings.TrimSpace(strings.ReplaceAll(partPath, "\\", "/"))
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

// ParseRelationshipNumber extracts the numeric part of an rId string.
func ParseRelationshipNumber(id string) (int, bool) {
	if !strings.HasPrefix(id, "rId") {
		return 0, false
	}
	var num int
	_, err := fmt.Sscanf(id, "rId%d", &num)
	if err != nil {
		return 0, false
	}
	return num, true
}

// ParseSlidePartNumber extracts the numeric part of a slide part name (e.g., slide1.xml).
func ParseSlidePartNumber(partPath string) (int, bool) {
	lastSlash := strings.LastIndex(partPath, "/")
	base := partPath
	if lastSlash >= 0 {
		base = partPath[lastSlash+1:]
	}
	if !strings.HasPrefix(base, "slide") || !strings.HasSuffix(base, ".xml") {
		return 0, false
	}
	var num int
	_, err := fmt.Sscanf(base, "slide%d.xml", &num)
	if err != nil {
		return 0, false
	}
	return num, true
}

// Shape represents a simplified view of a slide shape for editing.
type Shape struct {
	ID   int
	Name string
	Type string
	Text string
	X, Y int
	W, H int
}

// ShapeSearchQuery filters shapes for editor-wide search.
type ShapeSearchQuery struct {
	NameContains  string
	TypeEquals    string
	TextContains  string
	CaseSensitive bool
}

// ShapeSearchResult identifies one matched shape and its slide index.
type ShapeSearchResult struct {
	SlideIndex int
	Shape      Shape
}

// ChartSelector identifies a slide chart by index and/or relationship ID.
type ChartSelector struct {
	Index *int   `json:"index,omitempty"`
	RelID string `json:"rel_id,omitempty"`
}

// ChartSeriesData carries one chart series worth of input data.
type ChartSeriesData struct {
	Name       *string   `json:"name,omitempty"`
	Categories []string  `json:"categories,omitempty"`
	Values     []float64 `json:"values,omitempty"`
	XValues    []float64 `json:"x_values,omitempty"`
	YValues    []float64 `json:"y_values,omitempty"`
	Sizes      []float64 `json:"sizes,omitempty"`
}

// ChartDataUpdate is the complete chart update payload.
type ChartDataUpdate struct {
	Categories []string          `json:"categories,omitempty"`
	Series     []ChartSeriesData `json:"series,omitempty"`
}

// SlideChartRef describes a chart relationship discovered on a slide.
type SlideChartRef struct {
	Index     int
	RelID     string
	ChartPart string
}

// SlideLayoutInfo describes one available slide layout part.
type SlideLayoutInfo struct {
	Part       string
	Name       string
	MasterPart string
}

// SlideMasterCloneResult summarizes an in-package layout/master clone operation.
type SlideMasterCloneResult struct {
	MasterPart string
	ThemePart  string
	LayoutMap  map[string]string
}

// ShapeProps defines optional properties when creating a shape.
type ShapeProps struct {
	Name string `json:"name,omitempty"`
	// Add other properties as needed for Phase 1
}

// ShapeUpdate defines fields that can be updated on a shape.
// Pointers are used to indicate which fields are being updated (non-nil).
type ShapeUpdate struct {
	Text *string `json:"text,omitempty"`
	X    *int    `json:"x,omitempty"`
	Y    *int    `json:"y,omitempty"`
	W    *int    `json:"w,omitempty"`
	H    *int    `json:"h,omitempty"`
}

// SlideImageRef describes one image relationship on a slide.
type SlideImageRef struct {
	Index  int
	RelID  string
	Target string
}
