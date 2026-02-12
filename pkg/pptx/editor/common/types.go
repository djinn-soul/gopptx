package common

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

	RelTypeSlide       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide"
	RelTypeSlideLayout = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout"
	RelTypeNotesSlide  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide"
	RelTypeHyperlink   = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink"
	RelTypeImage       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"
	RelTypeChart       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart"
	RelTypeSectionList = "http://schemas.microsoft.com/office/2007/relationships/sectionList"

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
	if strings.HasPrefix(clean, "/") {
		return strings.TrimPrefix(clean, "/")
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
