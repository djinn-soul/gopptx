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

func dedupeContentTypeOverrides(overrides []contentTypeOverride) []contentTypeOverride {
	if len(overrides) == 0 {
		return overrides
	}
	seen := make(map[string]contentTypeOverride, len(overrides))
	order := make([]string, 0, len(overrides))
	for _, override := range overrides {
		key := "/" + common.CanonicalPartPath(strings.TrimPrefix(strings.TrimSpace(override.PartName), "/"))
		if _, exists := seen[key]; !exists {
			order = append(order, key)
		}
		seen[key] = contentTypeOverride{
			PartName:    key,
			ContentType: strings.TrimSpace(override.ContentType),
		}
	}
	deduped := make([]contentTypeOverride, 0, len(order))
	for _, key := range order {
		deduped = append(deduped, seen[key])
	}
	return deduped
}

func renderContentTypesDocument(doc contentTypesDocument) (string, error) {
	rendered, err := xml.MarshalIndent(doc, "", "")
	if err != nil {
		return "", err
	}
	return xml.Header + strings.TrimSpace(string(rendered)), nil
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
