package slide

import (
	"encoding/xml"
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const commentsPartType = "application/vnd.openxmlformats-officedocument.presentationml.comments+xml"

// contentTypesEscaper escapes XML attribute values for content types output.
// Package-level to avoid per-call allocation (strings.NewReplacer allocates).
//
//nolint:gochecknoglobals // read-only package-level replacer, never mutated
var contentTypesEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
)

// ContentTypesBase is an opaque pre-parsed [Content_Types].xml representation.
// Callers obtain it via ParseContentTypesBase and pass it to RewriteContentTypesFromBase
// to skip xml.Unmarshal on repeated saves when the bytes have not changed.
type ContentTypesBase interface {
	// unexported method keeps the interface opaque; only this package can implement it.
	contentTypesBase() contentTypesDocument
}

type contentTypesBaseImpl struct {
	doc contentTypesDocument
}

func (c contentTypesBaseImpl) contentTypesBase() contentTypesDocument { return c.doc }

// ParseContentTypesBase parses [Content_Types].xml and returns an opaque value
// suitable for caching and passing to RewriteContentTypesFromBase.
func ParseContentTypesBase(current []byte) (ContentTypesBase, error) {
	doc, err := parseContentTypesDocument(current)
	if err != nil {
		return nil, err
	}
	return contentTypesBaseImpl{doc}, nil
}

// RewriteContentTypesFromBase is identical to RewriteContentTypes but accepts a
// pre-parsed base (from ParseContentTypesBase) to avoid re-parsing on each call.
func RewriteContentTypesFromBase(
	base ContentTypesBase,
	slides []common.EditorSlideRef,
	mediaPaths []string,
	hasSections bool,
	chartPaths []string,
	notesPaths []string,
	themePaths []string,
	layoutPaths []string,
	masterPaths []string,
	hasNotesMaster bool,
	hasCommentAuthors bool,
	commentPaths []string,
	hasVBA bool,
	hasHandoutMaster bool,
	customXMLPropsPaths []string,
) (string, error) {
	doc := base.contentTypesBase()
	// Copy Defaults slice before passing to ensureContentTypeDefaults, which may
	// append to it.  This prevents the cached base from growing across saves.
	if len(mediaPaths) > 0 || hasVBA {
		cp := make([]contentTypeDefault, len(doc.Defaults))
		copy(cp, doc.Defaults)
		doc.Defaults = cp
	}
	return rewriteContentTypesFromDoc(doc, slides, mediaPaths, hasSections, chartPaths,
		notesPaths, themePaths, layoutPaths, masterPaths, hasNotesMaster, hasCommentAuthors,
		commentPaths, hasVBA, hasHandoutMaster, customXMLPropsPaths)
}

func RewriteContentTypes(
	current []byte,
	slides []common.EditorSlideRef,
	mediaPaths []string,
	hasSections bool,
	chartPaths []string,
	notesPaths []string,
	themePaths []string,
	layoutPaths []string,
	masterPaths []string,
	hasNotesMaster bool,
	hasCommentAuthors bool,
	commentPaths []string,
	hasVBA bool,
	hasHandoutMaster bool,
	customXMLPropsPaths []string,
) (string, error) {
	doc, err := parseContentTypesDocument(current)
	if err != nil {
		return "", err
	}
	return rewriteContentTypesFromDoc(doc, slides, mediaPaths, hasSections, chartPaths,
		notesPaths, themePaths, layoutPaths, masterPaths, hasNotesMaster, hasCommentAuthors,
		commentPaths, hasVBA, hasHandoutMaster, customXMLPropsPaths)
}

func rewriteContentTypesFromDoc(
	doc contentTypesDocument,
	slides []common.EditorSlideRef,
	mediaPaths []string,
	hasSections bool,
	chartPaths []string,
	notesPaths []string,
	themePaths []string,
	layoutPaths []string,
	masterPaths []string,
	hasNotesMaster bool,
	hasCommentAuthors bool,
	commentPaths []string,
	hasVBA bool,
	hasHandoutMaster bool,
	customXMLPropsPaths []string,
) (string, error) {
	ensureContentTypeDefaults(&doc, mediaPaths, hasVBA)

	overrides := filterDynamicOverrides(doc.Overrides, len(slides))
	overrides = appendSlideOverrides(overrides, slides)
	overrides = appendOptionalContentTypeOverride(overrides, hasSections, "/ppt/sectionList.xml",
		"application/vnd.microsoft.powerpoint.sectionList+xml")
	overrides = appendPathOverrides(overrides, chartPaths,
		"application/vnd.openxmlformats-officedocument.drawingml.chart+xml")
	overrides = appendPathOverrides(overrides, notesPaths,
		"application/vnd.openxmlformats-officedocument.presentationml.notesSlide+xml")
	overrides = appendPathOverrides(overrides, themePaths,
		"application/vnd.openxmlformats-officedocument.theme+xml")
	overrides = appendPathOverrides(overrides, layoutPaths,
		"application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml")
	overrides = appendPathOverrides(overrides, masterPaths,
		"application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasNotesMaster, "/ppt/notesMasters/notesMaster1.xml",
		"application/vnd.openxmlformats-officedocument.presentationml.notesMaster+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasHandoutMaster, "/ppt/handoutMasters/handoutMaster1.xml",
		"application/vnd.openxmlformats-officedocument.presentationml.handoutMaster+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasCommentAuthors, "/ppt/commentAuthors.xml",
		"application/vnd.openxmlformats-officedocument.presentationml.commentAuthors+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasVBA, "/ppt/vbaProject.bin",
		"application/vnd.ms-office.vbaProject")
	overrides = appendPathOverrides(overrides, commentPaths, commentsPartType)
	overrides = appendPathOverrides(overrides, customXMLPropsPaths,
		"application/vnd.openxmlformats-officedocument.customXmlProperties+xml")
	overrides = dedupeContentTypeOverrides(overrides)

	sort.Slice(overrides, func(i, j int) bool { return overrides[i].PartName < overrides[j].PartName })
	doc.Overrides = overrides
	return renderContentTypesDocument(doc)
}

func parseContentTypesDocument(current []byte) (contentTypesDocument, error) {
	if len(current) == 0 {
		return contentTypesDocument{}, errors.New("missing content types content")
	}

	var doc contentTypesDocument
	if err := xml.Unmarshal(current, &doc); err != nil {
		return contentTypesDocument{}, fmt.Errorf("parse content types: %w", err)
	}
	if strings.TrimSpace(doc.XMLNS) == "" {
		doc.XMLNS = common.ContentTypesXMLNS
	}
	return doc, nil
}

func ensureContentTypeDefaults(doc *contentTypesDocument, mediaPaths []string, hasVBA bool) {
	exts := make(map[string]struct{}, len(doc.Defaults))
	for _, d := range doc.Defaults {
		exts[strings.ToLower(d.Extension)] = struct{}{}
	}
	for _, mediaPath := range mediaPaths {
		ext := strings.TrimPrefix(strings.ToLower(path.Ext(strings.TrimSpace(mediaPath))), ".")
		if _, ok := exts[ext]; ok {
			continue
		}
		contentType := contentTypeForExtension(ext)
		if contentType == "" {
			continue
		}
		doc.Defaults = append(doc.Defaults, contentTypeDefault{
			Extension:   ext,
			ContentType: contentType,
		})
		exts[ext] = struct{}{}
	}

	if hasVBA {
		if _, ok := exts["bin"]; !ok {
			doc.Defaults = append(doc.Defaults, contentTypeDefault{
				Extension:   "bin",
				ContentType: "application/vnd.ms-office.vbaProject",
			})
			exts["bin"] = struct{}{}
		}
	}
}

func contentTypeForExtension(ext string) string {
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")
	switch ext {
	case "png":
		return "image/png"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "gif":
		return "image/gif"
	case "bmp":
		return "image/bmp"
	case "tif", "tiff":
		return "image/tiff"
	case "wav":
		return "audio/wav"
	case "mp3":
		return "audio/mpeg"
	case "m4a":
		return "audio/mp4"
	case "wma":
		return "audio/x-ms-wma"
	case "ogg":
		return "audio/ogg"
	case "flac":
		return "audio/flac"
	case "aac":
		return "audio/aac"
	case "mp4":
		return "video/mp4"
	case "webm":
		return "video/webm"
	case "avi":
		return "video/x-msvideo"
	case "wmv":
		return "video/x-ms-wmv"
	case "mov":
		return "video/quicktime"
	case "mkv":
		return "video/x-matroska"
	case "m4v":
		return "video/x-m4v"
	case "bin":
		return "application/vnd.openxmlformats-officedocument.oleObject"
	case "xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	default:
		return ""
	}
}

func filterDynamicOverrides(existing []contentTypeOverride, slideCapacity int) []contentTypeOverride {
	filtered := make([]contentTypeOverride, 0, len(existing)+slideCapacity+1)
	for _, override := range existing {
		part := common.CanonicalPartPath(override.PartName)
		if shouldSkipOverridePart(part) {
			continue
		}
		filtered = append(filtered, override)
	}
	return filtered
}

func shouldSkipOverridePart(part string) bool {
	if isSlidePartOverride(part) {
		return true
	}
	return part == "ppt/sectionList.xml" ||
		part == "ppt/commentAuthors.xml" ||
		strings.HasPrefix(part, "ppt/charts/chart") ||
		strings.HasPrefix(part, "ppt/notesSlides/notesSlide") ||
		strings.HasPrefix(part, "ppt/notesMasters/notesMaster") ||
		strings.HasPrefix(part, "ppt/theme/theme") ||
		strings.HasPrefix(part, "ppt/slideLayouts/slideLayout") ||
		strings.HasPrefix(part, "ppt/slideMasters/slideMaster") ||
		strings.HasPrefix(part, "ppt/comments/comment")
}

func appendSlideOverrides(
	overrides []contentTypeOverride,
	slides []common.EditorSlideRef,
) []contentTypeOverride {
	for _, slide := range slides {
		overrides = append(overrides, contentTypeOverride{
			PartName:    "/" + common.CanonicalPartPath(slide.Part),
			ContentType: common.SlideContentType,
		})
	}
	return overrides
}

func appendPathOverrides(
	overrides []contentTypeOverride,
	paths []string,
	contentType string,
) []contentTypeOverride {
	for _, p := range paths {
		overrides = append(overrides, contentTypeOverride{
			PartName:    "/" + common.CanonicalPartPath(p),
			ContentType: contentType,
		})
	}
	return overrides
}

func appendOptionalContentTypeOverride(
	overrides []contentTypeOverride,
	include bool,
	partName,
	contentType string,
) []contentTypeOverride {
	if !include {
		return overrides
	}
	return append(overrides, contentTypeOverride{
		PartName:    partName,
		ContentType: contentType,
	})
}

// normalizeOverridePartName returns the canonical "/ppt/..." form.
// Fast path: if partName is already in that form (starts with "/" and the
// rest is already canonical), return partName directly without allocating.
func normalizeOverridePartName(partName string) string {
	trimmed := strings.TrimSpace(partName)
	canonical := common.CanonicalPartPath(strings.TrimPrefix(trimmed, "/"))
	// Check if trimmed == "/" + canonical without constructing the concatenation.
	if len(trimmed) == len(canonical)+1 && trimmed[0] == '/' && trimmed[1:] == canonical {
		return trimmed // already normalized — zero alloc
	}
	return "/" + canonical
}

func dedupeContentTypeOverrides(overrides []contentTypeOverride) []contentTypeOverride {
	if len(overrides) == 0 {
		return overrides
	}
	// Fast path: scan for duplicates without building the full output.
	// In the normal save flow the list is already unique; skip rebuilding it.
	seen := make(map[string]struct{}, len(overrides))
	hasDup := false
	for _, o := range overrides {
		key := normalizeOverridePartName(o.PartName)
		if _, dup := seen[key]; dup {
			hasDup = true
			break
		}
		seen[key] = struct{}{}
	}
	if !hasDup {
		return overrides // common case: no duplicates
	}

	// Slow path: rebuild with first-occurrence semantics.
	clear(seen)
	order := make([]string, 0, len(overrides))
	overrideByKey := make(map[string]contentTypeOverride, len(overrides))
	for _, o := range overrides {
		key := normalizeOverridePartName(o.PartName)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			order = append(order, key)
		}
		overrideByKey[key] = contentTypeOverride{
			PartName:    key,
			ContentType: strings.TrimSpace(o.ContentType),
		}
	}
	deduped := make([]contentTypeOverride, 0, len(order))
	for _, key := range order {
		deduped = append(deduped, overrideByKey[key])
	}
	return deduped
}

func renderContentTypesDocument(doc contentTypesDocument) (string, error) {
	// Hand-rolled renderer replaces xml.MarshalIndent to eliminate reflection overhead.
	// Each entry is ~100 bytes; pre-grow avoids builder reallocations.
	var b strings.Builder
	b.Grow(80 + (len(doc.Defaults)+len(doc.Overrides))*100)
	b.WriteString(`<Types xmlns="`)
	b.WriteString(contentTypesEscaper.Replace(doc.XMLNS))
	b.WriteString(`">`)
	for _, d := range doc.Defaults {
		b.WriteString("\n<Default Extension=\"")
		b.WriteString(contentTypesEscaper.Replace(d.Extension))
		b.WriteString("\" ContentType=\"")
		b.WriteString(contentTypesEscaper.Replace(d.ContentType))
		b.WriteString("\"/>")
	}
	for _, o := range doc.Overrides {
		b.WriteString("\n<Override PartName=\"")
		b.WriteString(contentTypesEscaper.Replace(o.PartName))
		b.WriteString("\" ContentType=\"")
		b.WriteString(contentTypesEscaper.Replace(o.ContentType))
		b.WriteString("\"/>")
	}
	b.WriteString("\n</Types>")
	return xml.Header + b.String(), nil
}

type contentTypesDocument struct {
	XMLName   xml.Name              `xml:"Types"`
	XMLNS     string                `xml:"xmlns,attr,omitempty"`
	Defaults  []contentTypeDefault  `xml:"Default"`
	Overrides []contentTypeOverride `xml:"Override"`
}

type contentTypeDefault struct {
	Extension   string `xml:"Extension,attr"`
	ContentType string `xml:"ContentType,attr"`
}

type contentTypeOverride struct {
	PartName    string `xml:"PartName,attr"`
	ContentType string `xml:"ContentType,attr"`
}

func isSlidePartOverride(partName string) bool {
	clean := common.CanonicalPartPath(strings.TrimPrefix(strings.TrimSpace(partName), "/"))
	return strings.HasPrefix(clean, "ppt/slides/slide") && strings.HasSuffix(clean, ".xml")
}
