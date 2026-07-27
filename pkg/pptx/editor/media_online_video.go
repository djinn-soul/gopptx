package editor

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodmedia "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/media"
)

// AddOnlineVideo places a shape that links to a video hosted elsewhere, the
// Insert > Media > Video > Online Video case from upstream issue #1034.
//
// Unlike AddVideo the media is not embedded: the slide keeps an external
// relationship to videoURL and shows posterData (or a default poster) as the
// click target.
func (e *PresentationEditor) AddOnlineVideo(
	slideIndex int,
	videoURL string,
	posterData []byte,
	posterFormat string,
	altText string,
	x, y, w, h float64,
) (int, error) {
	if err := editormodmedia.ValidateMediaSlideIndex(slideIndex, len(e.slides)); err != nil {
		return 0, err
	}
	normalizedURL, err := normalizeOnlineVideoURL(videoURL)
	if err != nil {
		return 0, err
	}

	posterRelID, err := e.registerOnlineVideoPoster(slideIndex, posterData, posterFormat)
	if err != nil {
		return 0, err
	}

	videoRelID, err := e.addExternalSlideRel(slideIndex, normalizedURL, common.RelTypeVideo)
	if err != nil {
		return 0, err
	}

	slideRef := e.slides[slideIndex]
	content, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return 0, errors.New("read slide part: not found")
	}

	newID := e.maxObjectID(slideRef.Part, content) + 1
	e.reserveObjectIDs(slideRef.Part, newID)

	videoXML := editormodmedia.BuildOnlineVideoShapeXML(
		newID, videoRelID, posterRelID, altText, x, y, w, h,
	)
	updatedContent, err := editormodmedia.AppendShapeXMLToSlide(content, videoXML)
	if err != nil {
		return 0, err
	}
	e.parts.Set(slideRef.Part, updatedContent)
	return newID, nil
}

func (e *PresentationEditor) registerOnlineVideoPoster(
	slideIndex int,
	posterData []byte,
	posterFormat string,
) (string, error) {
	data := posterData
	format := strings.TrimSpace(posterFormat)
	if len(data) == 0 {
		data = editormodmedia.DefaultVideoPosterPNG()
		format = formatPNG
	}
	if format == "" {
		format = formatPNG
	}

	posterPart, err := e.RegisterImage(data, format)
	if err != nil {
		return "", fmt.Errorf("register poster frame: %w", err)
	}
	relID, err := e.getOrCreateSlideRel(slideIndex, posterPart)
	if err != nil {
		return "", fmt.Errorf("create poster rel: %w", err)
	}
	return relID, nil
}

// addExternalSlideRel adds a relationship whose target is a URL rather than a
// package part, so it carries TargetMode="External".
func (e *PresentationEditor) addExternalSlideRel(
	slideIndex int,
	target string,
	relType string,
) (string, error) {
	slideRef := e.slides[slideIndex]
	rels, err := e.slideRelationships(slideRef.Part)
	if err != nil {
		return "", err
	}
	for _, rel := range rels {
		if rel.Type == relType && rel.Target == target && rel.TargetMode == targetModeExternal {
			return rel.ID, nil
		}
	}

	relID := fmt.Sprintf("rId%d", common.NextRelationshipNumber(rels))
	rels = append(rels, common.EditorRelationship{
		ID:         relID,
		Type:       relType,
		Target:     target,
		TargetMode: targetModeExternal,
	})
	e.parts.Set(common.SlideRelsPartName(slideRef.Part), []byte(renderRelationshipsXML(rels)))
	return relID, nil
}

// normalizeOnlineVideoURL keeps obviously unusable targets out of the package.
// Only http(s) is accepted: PowerPoint will not play anything else, and a
// file:// or javascript: target in a relationship is a hazard rather than a
// feature.
func normalizeOnlineVideoURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("video URL is required")
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("invalid video URL: %w", err)
	}
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https":
	default:
		return "", fmt.Errorf("video URL scheme %q is not supported, use http or https", parsed.Scheme)
	}
	if parsed.Host == "" {
		return "", errors.New("video URL has no host")
	}
	return trimmed, nil
}
