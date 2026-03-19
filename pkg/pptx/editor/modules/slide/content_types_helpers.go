package slide

import (
	"encoding/xml"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

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

func filterDynamicOverrides(existing []contentTypeOverride, extraCapacity int) []contentTypeOverride {
	filtered := make([]contentTypeOverride, 0, len(existing)+extraCapacity)
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
	b.Grow(contentTypesBaseGrowCap + (len(doc.Defaults)+len(doc.Overrides))*contentTypesItemGrowCap)
	b.WriteString(`<Types xmlns="`)
	writeEscapedXMLAttr(&b, doc.XMLNS)
	b.WriteString(`">`)
	for _, d := range doc.Defaults {
		b.WriteString("\n<Default Extension=\"")
		writeEscapedXMLAttr(&b, d.Extension)
		b.WriteString("\" ContentType=\"")
		writeEscapedXMLAttr(&b, d.ContentType)
		b.WriteString("\"/>")
	}
	for _, o := range doc.Overrides {
		b.WriteString("\n<Override PartName=\"")
		writeEscapedXMLAttr(&b, o.PartName)
		b.WriteString("\" ContentType=\"")
		writeEscapedXMLAttr(&b, o.ContentType)
		b.WriteString("\"/>")
	}
	b.WriteString("\n</Types>")
	return xml.Header + b.String(), nil
}

func writeEscapedXMLAttr(b *strings.Builder, value string) {
	start := 0
	for i := range len(value) {
		var replacement string
		switch value[i] {
		case '&':
			replacement = "&amp;"
		case '<':
			replacement = "&lt;"
		case '>':
			replacement = "&gt;"
		case '"':
			replacement = "&quot;"
		default:
			continue
		}
		if start < i {
			b.WriteString(value[start:i])
		}
		b.WriteString(replacement)
		start = i + 1
	}
	if start == 0 {
		b.WriteString(value)
		return
	}
	if start < len(value) {
		b.WriteString(value[start:])
	}
}

func isSlidePartOverride(partName string) bool {
	clean := common.CanonicalPartPath(strings.TrimPrefix(strings.TrimSpace(partName), "/"))
	return strings.HasPrefix(clean, "ppt/slides/slide") && strings.HasSuffix(clean, ".xml")
}
