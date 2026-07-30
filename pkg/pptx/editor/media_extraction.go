package editor

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Media kinds reported by ListSlideMedia.
const (
	MediaKindImage = "image"
	MediaKindAudio = "audio"
	MediaKindVideo = "video"
)

// ListSlideMedia reports every media relationship on a slide, not only images.
// A video or audio asset embedded in a deck was previously reachable only by
// digging through relationships by hand (upstream #1049).
//
// PowerPoint declares a movie twice: once as a video (or audio) relationship
// and once as a Microsoft media relationship pointing at the same part. Both
// are reported so a caller sees exactly what the package contains, and the part
// path lets a caller deduplicate.
func (e *PresentationEditor) ListSlideMedia(slideIndex int) ([]common.SlideMediaRef, error) {
	if e == nil || e.parts == nil {
		return nil, errors.New("editor cannot be nil")
	}
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, fmt.Errorf("slide index %d out of range [0,%d)", slideIndex, len(e.slides))
	}

	slidePart := e.slides[slideIndex].Part
	rels, err := e.slideRelationships(slidePart)
	if err != nil {
		return nil, err
	}

	out := make([]common.SlideMediaRef, 0)
	for _, rel := range rels {
		kind, ok := mediaKindForRelType(rel.Type)
		if !ok {
			continue
		}
		out = append(out, e.buildMediaRef(len(out), slidePart, kind, rel))
	}
	return out, nil
}

func (e *PresentationEditor) buildMediaRef(
	index int,
	slidePart, kind string,
	rel common.EditorRelationship,
) common.SlideMediaRef {
	ref := common.SlideMediaRef{
		Index:  index,
		RelID:  rel.ID,
		Kind:   kind,
		Target: rel.Target,
	}
	if strings.EqualFold(rel.TargetMode, "External") {
		// A linked movie lives outside the package; there are no bytes to read.
		ref.External = true
		return ref
	}
	ref.PartPath = common.CanonicalPartPath(path.Join(path.Dir(slidePart), rel.Target))
	if data, ok := e.parts.Get(ref.PartPath); ok {
		ref.SizeBytes = len(data)
	}
	ref.ContentType = e.declaredContentType(ref.PartPath)
	return ref
}

// ExtractMedia returns the bytes of a media part, so a caller can write an
// embedded movie or sound out as its own file.
func (e *PresentationEditor) ExtractMedia(partPath string) ([]byte, error) {
	if e == nil || e.parts == nil {
		return nil, errors.New("editor cannot be nil")
	}
	canonical := common.CanonicalPartPath(partPath)
	data, ok := e.parts.Get(canonical)
	if !ok {
		return nil, fmt.Errorf("media part %s not found", canonical)
	}
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

func mediaKindForRelType(relType string) (string, bool) {
	switch relType {
	case common.RelTypeImage:
		return MediaKindImage, true
	case common.RelTypeAudio:
		return MediaKindAudio, true
	case common.RelTypeVideo, common.RelTypeMedia:
		return MediaKindVideo, true
	default:
		return "", false
	}
}

// declaredContentType reports what [Content_Types].xml says a part is, falling
// back to the extension mapping when the package declares nothing, so the
// answer matches the file that would be written out.
func (e *PresentationEditor) declaredContentType(partPath string) string {
	data, ok := e.parts.Get(common.ContentTypesPath)
	if ok {
		if declared := contentTypeFromDeclarations(string(data), partPath); declared != "" {
			return declared
		}
	}
	return mediaContentTypeForExtension(path.Ext(partPath))
}

func contentTypeFromDeclarations(contentTypesXML, partPath string) string {
	if override := overrideContentTypePattern(partPath).FindStringSubmatch(contentTypesXML); len(override) > 1 {
		return override[1]
	}
	ext := strings.TrimPrefix(strings.ToLower(path.Ext(partPath)), ".")
	if ext == "" {
		return ""
	}
	if match := defaultContentTypePattern(ext).FindStringSubmatch(contentTypesXML); len(match) > 1 {
		return match[1]
	}
	return ""
}

func overrideContentTypePattern(partPath string) *regexp.Regexp {
	return regexp.MustCompile(
		`(?i)<Override[^>]*PartName="/` + regexp.QuoteMeta(partPath) +
			`"[^>]*ContentType="([^"]+)"`,
	)
}

func defaultContentTypePattern(ext string) *regexp.Regexp {
	return regexp.MustCompile(
		`(?i)<Default[^>]*Extension="` + regexp.QuoteMeta(ext) +
			`"[^>]*ContentType="([^"]+)"`,
	)
}

// mediaContentTypeForExtension answers for packages that declare neither an
// Override nor a Default for the part, reusing the MediaFormat table so the
// extraction path and the insertion path agree on every container.
func mediaContentTypeForExtension(ext string) string {
	normalized := strings.TrimPrefix(strings.ToLower(ext), ".")
	if mime := MediaFormat(normalized).MIMEType(); mime != "" {
		return mime
	}
	return contentTypeForImageFormat(normalized)
}
