package presentation

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormedia "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// A clip reaches the package as three things: the media part itself, the pair
// of relationships PowerPoint expects on the slide (the typed video/audio one
// and the p14 media one), and a p:pic in the shape tree. The markup is the
// editor's, which is what PowerPoint round-trips.

// mediaSlideParts is the result of attaching a slide's media: the rewritten
// slide XML and relationships, plus the binary parts to write.
type mediaSlideParts struct {
	SlideXML string
	RelsXML  string
	Parts    map[string][]byte
}

// slideMedia returns the clips on a slide that carry something to write.
func slideMedia(slide elements.SlideContent) []shapes.Media {
	if len(slide.Media) == 0 {
		return nil
	}
	out := make([]shapes.Media, 0, len(slide.Media))
	for _, clip := range slide.Media {
		if clip.IsOnline() || len(clip.Data) > 0 || clip.Path != "" {
			out = append(out, clip)
		}
	}
	return out
}

// mediaPartStartIndex is the 1-based media part number of the first clip on
// slide slideIdx. Parts are numbered across the package in slide order so the
// content-type pass and the render pass agree.
func mediaPartStartIndex(slides []elements.SlideContent, slideIdx int) int {
	index := 1
	for i := 0; i < slideIdx && i < len(slides); i++ {
		index += len(slideMedia(slides[i]))
	}
	return index
}

// mediaExtensions lists the distinct file extensions the package's clips use,
// each of which needs a Default entry in [Content_Types].xml.
func mediaExtensions(slides []elements.SlideContent) []string {
	seen := map[string]bool{}
	var out []string
	for _, slide := range slides {
		for _, clip := range slideMedia(slide) {
			if clip.IsOnline() {
				continue
			}
			ext := clip.FormatOrExtension()
			if seen[ext] {
				continue
			}
			seen[ext] = true
			out = append(out, ext)
		}
		for _, clip := range slideMedia(slide) {
			if len(clip.PosterData) == 0 {
				continue
			}
			ext := posterExtension(clip)
			if seen[ext] {
				continue
			}
			seen[ext] = true
			out = append(out, ext)
		}
	}
	return out
}

// mediaMIMETypes maps the extensions a deck's clips use to the content types
// PowerPoint expects. An unlisted extension falls back to a generic type,
// which PowerPoint accepts for playback it can decode.
//
//nolint:gochecknoglobals // read-only lookup table
var mediaMIMETypes = map[string]string{
	"mp4":                  "video/mp4",
	"m4v":                  "video/mp4",
	"mov":                  "video/quicktime",
	"avi":                  "video/x-msvideo",
	"wmv":                  "video/x-ms-wmv",
	"webm":                 "video/webm",
	"mp3":                  "audio/mpeg",
	"m4a":                  "audio/mp4",
	"wav":                  "audio/wav",
	"wma":                  "audio/x-ms-wma",
	"aac":                  "audio/aac",
	"ogg":                  "audio/ogg",
	posterDefaultExtension: "image/png",
	"jpg":                  "image/jpeg",
	"jpeg":                 "image/jpeg",
	"gif":                  "image/gif",
}

// addMediaContentTypes registers a Default entry per media extension used.
func addMediaContentTypes(contentTypes string, slides []elements.SlideContent) string {
	for _, ext := range mediaExtensions(slides) {
		mime, ok := mediaMIMETypes[ext]
		if !ok {
			mime = "application/octet-stream"
		}
		contentTypes = pptxxml.WithContentTypeDefault(contentTypes, ext, mime)
	}
	return contentTypes
}

func posterExtension(clip shapes.Media) string {
	if format := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(clip.PosterFormat)), "."); format != "" {
		return format
	}
	return posterDefaultExtension
}

// posterDefaultExtension is the format a poster with no stated one is assumed
// to be.
const posterDefaultExtension = "png"

// attachMedia writes each clip's part, its relationships and its picture into
// the slide.
func attachMedia(
	slideXML, relsXML string,
	clips []shapes.Media,
	startPartIndex int,
	slideIndex, slideCount int,
) (mediaSlideParts, error) {
	result := mediaSlideParts{SlideXML: slideXML, RelsXML: relsXML, Parts: map[string][]byte{}}
	if len(clips) == 0 {
		return result, nil
	}

	nextRID := nextRelationshipID(relsXML)
	shapeID := nextShapeID(slideXML)

	var markup strings.Builder
	timings := make([]mediaTiming, 0, len(clips))
	for i, clip := range clips {
		partIndex := startPartIndex + i
		clipXML, err := attachOneClip(&result, clip, partIndex, &nextRID, shapeID)
		if err != nil {
			return result, err
		}
		markup.WriteString(clipXML)
		timings = append(timings, mediaTiming{clip: clip, shapeID: shapeID})
		shapeID++
	}

	result.SlideXML = insertBeforeShapeTreeClose(result.SlideXML, markup.String())
	if err := applyMediaTimings(&result, timings, slideIndex, slideCount); err != nil {
		return result, err
	}
	return result, nil
}

// mediaTiming pairs a clip with the shape id its picture was given, which is
// how the timing tree refers to it.
type mediaTiming struct {
	clip    shapes.Media
	shapeID int
}

// applyMediaTimings writes the p:timing entries that carry autoplay, looping,
// muting, volume and play-across-slides. Without them a clip sits on the slide
// but never plays on its own.
func applyMediaTimings(result *mediaSlideParts, timings []mediaTiming, slideIndex, slideCount int) error {
	for _, timing := range timings {
		if timing.clip.IsOnline() {
			// An online video has no embedded part for the timing node to name.
			continue
		}
		kind := "video"
		if timing.clip.Kind == shapes.MediaAudio {
			kind = "audio"
		}
		updated, err := editormedia.ApplyMediaTiming(
			[]byte(result.SlideXML),
			kind,
			timing.shapeID,
			editormedia.MediaTimingOptions{
				AutoPlay:         timing.clip.Playback.AutoPlay,
				LoopPlayback:     timing.clip.Playback.Loop,
				Muted:            timing.clip.Playback.Muted,
				Volume:           timing.clip.Playback.Volume,
				ShowWhenStopped:  !timing.clip.Playback.HideWhenStopped,
				PlayAcrossSlides: timing.clip.Playback.PlayAcrossSlides,
				SlideIndex:       slideIndex,
				SlideCount:       slideCount,
			},
		)
		if err != nil {
			return fmt.Errorf("apply media timing: %w", err)
		}
		result.SlideXML = string(updated)
	}
	return nil
}

func attachOneClip(
	result *mediaSlideParts,
	clip shapes.Media,
	partIndex int,
	nextRID *int,
	shapeID int,
) (string, error) {
	posterRelID := attachPoster(result, clip, partIndex, nextRID)
	x, y := float64(clip.X), float64(clip.Y)
	cx, cy := clip.SizeOrDefault()

	if clip.IsOnline() {
		linkRelID := "rId" + strconv.Itoa(*nextRID)
		*nextRID++
		result.RelsXML = pptxxml.WithExternalRelationship(
			result.RelsXML, linkRelID, editorcommon.RelTypeVideo, clip.SourceURL,
		)
		return editormedia.BuildOnlineVideoShapeXML(
			shapeID, linkRelID, posterRelID, clip.AltText, x, y, float64(cx), float64(cy),
		), nil
	}

	data, err := clipBytes(clip)
	if err != nil {
		return "", err
	}
	partPath := fmt.Sprintf("ppt/media/media%d.%s", partIndex, clip.FormatOrExtension())
	result.Parts[partPath] = data

	// PowerPoint wants both relationships: the typed one names the clip, the
	// p14 one embeds it.
	typedRelID := "rId" + strconv.Itoa(*nextRID)
	*nextRID++
	mediaRelID := "rId" + strconv.Itoa(*nextRID)
	*nextRID++

	typedRelType := editorcommon.RelTypeVideo
	if clip.Kind == shapes.MediaAudio {
		typedRelType = editorcommon.RelTypeAudio
	}
	target := "../media/" + partPath[len("ppt/media/"):]
	result.RelsXML = pptxxml.WithRelationship(result.RelsXML, typedRelID, typedRelType, target)
	result.RelsXML = pptxxml.WithRelationship(result.RelsXML, mediaRelID, editorcommon.RelTypeMedia, target)

	if clip.Kind == shapes.MediaAudio {
		return editormedia.BuildAudioShapeXML(
			shapeID, typedRelID, mediaRelID, posterRelID, clip.AltText, x, y, float64(cx), float64(cy),
		), nil
	}
	return editormedia.BuildVideoShapeXML(
		shapeID, typedRelID, mediaRelID, posterRelID, clip.AltText, x, y, float64(cx), float64(cy),
	), nil
}

// attachPoster writes the still image a clip shows before playback, returning
// its relationship id or "" when the clip has none.
func attachPoster(result *mediaSlideParts, clip shapes.Media, partIndex int, nextRID *int) string {
	if len(clip.PosterData) == 0 {
		return ""
	}
	posterPath := fmt.Sprintf("ppt/media/mediaPoster%d.%s", partIndex, posterExtension(clip))
	result.Parts[posterPath] = clip.PosterData

	relID := "rId" + strconv.Itoa(*nextRID)
	*nextRID++
	result.RelsXML = pptxxml.WithRelationship(
		result.RelsXML,
		relID,
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/image",
		"../media/"+posterPath[len("ppt/media/"):],
	)
	return relID
}

func clipBytes(clip shapes.Media) ([]byte, error) {
	if len(clip.Data) > 0 {
		return clip.Data, nil
	}
	data, err := os.ReadFile(clip.Path)
	if err != nil {
		return nil, fmt.Errorf("read media %s: %w", clip.Path, err)
	}
	return data, nil
}
