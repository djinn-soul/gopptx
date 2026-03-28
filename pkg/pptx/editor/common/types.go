package editorcommon

import (
	"strconv"
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

	RelTypeSlide          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide"
	RelTypeSlideMaster    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster"
	RelTypeSlideLayout    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout"
	RelTypeNotesSlide     = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide"
	RelTypeNotesMaster    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesMaster"
	RelTypeHyperlink      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink"
	RelTypeImage          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"
	RelTypeChart          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart"
	RelTypeAudio          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/audio"
	RelTypeVideo          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/video"
	RelTypeMedia          = "http://schemas.microsoft.com/office/2007/relationships/media"
	RelTypeTheme          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme"
	RelTypeSectionList    = "http://schemas.microsoft.com/office/2007/relationships/sectionList"
	RelTypePackage        = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/package"
	RelTypeTableStyles    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/tableStyles"
	RelTypeCustomXML      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXml"
	RelTypeCustomXMLProps = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXmlProps"

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
	Hidden  bool
}

// xmlEscaper is a package-level replacer so XMLEscape never allocates.
//
//nolint:gochecknoglobals // read-only package-level replacer, never mutated
var xmlEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&apos;",
)

// XMLEscape provides basic XML attribute escaping.
func XMLEscape(value string) string {
	return xmlEscaper.Replace(value)
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
// Uses strconv.Atoi (zero-alloc) instead of fmt.Sscanf to avoid scanner allocation.
func ParseRelationshipNumber(id string) (int, bool) {
	if !strings.HasPrefix(id, "rId") {
		return 0, false
	}
	num, err := strconv.Atoi(id[3:])
	if err != nil {
		return 0, false
	}
	return num, true
}

// NextRelationshipNumber returns the next available relationship number
// (starting from 1) that is not present in the existing relationships list.
func NextRelationshipNumber(rels []EditorRelationship) int {
	nextNum := 1
	for _, r := range rels {
		if num, ok := ParseRelationshipNumber(r.ID); ok && num >= nextNum {
			nextNum = num + 1
		}
	}
	return nextNum
}

// ParseSlidePartNumber extracts the numeric part of a slide part name (e.g., slide1.xml).
// Uses strconv.Atoi (zero-alloc) instead of fmt.Sscanf to avoid scanner allocation.
func ParseSlidePartNumber(partPath string) (int, bool) {
	lastSlash := strings.LastIndex(partPath, "/")
	base := partPath
	if lastSlash >= 0 {
		base = partPath[lastSlash+1:]
	}
	const prefix = "slide"
	const suffix = ".xml"
	if !strings.HasPrefix(base, prefix) || !strings.HasSuffix(base, suffix) {
		return 0, false
	}
	numStr := base[len(prefix) : len(base)-len(suffix)]
	num, err := strconv.Atoi(numStr)
	if err != nil || num <= 0 {
		return 0, false
	}
	return num, true
}
