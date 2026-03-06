package editor

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodmedia "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/media"
)

var embeddedImageRelPattern = regexp.MustCompile(`r:embed="([^"]+)"`)

const (
	mimePNG    = "image/png"
	mimeJPEG   = "image/jpeg"
	mimeGIF    = "image/gif"
	mimeBMP    = "image/bmp"
	mimeTIFF   = "image/tiff"
	formatPNG  = "png"
	formatJPEG = "jpeg"
	formatGIF  = "gif"
	formatBMP  = "bmp"
	formatTIFF = "tiff"
)

type imagePathCacheEntry struct {
	data     []byte
	format   string
	partPath string
}

// AddImageFromBase64 adds an image from raw base64 bytes or a data URI payload.
func (e *PresentationEditor) AddImageFromBase64(
	slideIndex int,
	base64Data string,
	format string,
	x, y, w, h float64,
	opts *common.ShapeUpdate,
) (int, error) {
	data, detectedFormat, err := editormodmedia.DecodeBase64ImagePayload(base64Data)
	if err != nil {
		return 0, err
	}
	resolvedFormat := strings.TrimSpace(format)
	if resolvedFormat == "" {
		resolvedFormat = detectedFormat
	}
	if resolvedFormat == "" {
		return 0, errors.New("image format is required for base64 image payload")
	}
	return e.AddImageFromBytes(slideIndex, data, resolvedFormat, x, y, w, h, opts)
}

// AddImageFromURL downloads an image and embeds it into the specified slide.
func (e *PresentationEditor) AddImageFromURL(
	slideIndex int,
	sourceURL string,
	x, y, w, h float64,
	opts *common.ShapeUpdate,
) (int, error) {
	data, format, err := editormodmedia.FetchImageFromURL(sourceURL)
	if err != nil {
		return 0, err
	}
	return e.AddImageFromBytes(slideIndex, data, format, x, y, w, h, opts)
}

func buildImageMetadata(data []byte, cfg image.Config, format string) *common.ImageMetadata {
	return &common.ImageMetadata{
		Width:       cfg.Width,
		Height:      cfg.Height,
		Format:      strings.ToLower(strings.TrimSpace(format)),
		ContentType: imageContentType(data, format),
		Hash:        imageSHA256Hex(data),
	}
}

func imageSHA256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func imageContentType(data []byte, format string) string {
	contentType := strings.TrimSpace(http.DetectContentType(data))
	if contentType != "" && contentType != "application/octet-stream" {
		return contentType
	}

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "jpg", formatJPEG:
		return mimeJPEG
	case formatPNG:
		return mimePNG
	case formatGIF:
		return mimeGIF
	case formatBMP:
		return mimeBMP
	case "tif", formatTIFF:
		return mimeTIFF
	default:
		return contentType
	}
}

func (e *PresentationEditor) registerImageFromPath(imagePath, formatHint string) (string, error) {
	cleanPath := filepath.Clean(imagePath)
	e.imagePathMu.RLock()
	entry, ok := e.imagePathCache[cleanPath]
	e.imagePathMu.RUnlock()
	if ok && entry.partPath != "" {
		return entry.partPath, nil
	}

	data, format, err := e.loadImageFromPath(cleanPath, formatHint)
	if err != nil {
		return "", err
	}

	partPath, err := e.RegisterImage(data, format)
	if err != nil {
		return "", err
	}

	e.imagePathMu.Lock()
	e.imagePathCache[cleanPath] = imagePathCacheEntry{data: nil, format: format, partPath: partPath}
	e.imagePathMu.Unlock()
	return partPath, nil
}

func (e *PresentationEditor) loadImageFromPath(cleanPath, formatHint string) ([]byte, string, error) {
	e.imagePathMu.RLock()
	entry, ok := e.imagePathCache[cleanPath]
	e.imagePathMu.RUnlock()
	if ok && len(entry.data) > 0 {
		format := normalizeImageFormatHint(formatHint)
		if format == "" {
			format = entry.format
		}
		return entry.data, format, nil
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, "", err
	}

	format := normalizeImageFormatHint(formatHint)
	if format == "" {
		format = normalizeImageFormatHint(filepath.Ext(cleanPath))
	}
	if format == "" {
		return nil, "", fmt.Errorf("image path %q has no detectable format", cleanPath)
	}

	e.imagePathMu.Lock()
	entry = e.imagePathCache[cleanPath]
	if len(entry.data) == 0 {
		entry.data = data
	}
	if entry.format == "" {
		entry.format = format
	}
	e.imagePathCache[cleanPath] = entry
	e.imagePathMu.Unlock()
	return entry.data, entry.format, nil
}

func normalizeImageFormatHint(format string) string {
	trimmed := strings.TrimSpace(strings.ToLower(format))
	return strings.TrimPrefix(trimmed, ".")
}

// ListSlideImages returns image relationships for a slide in relationship order.
func (e *PresentationEditor) ListSlideImages(slideIndex int) ([]common.SlideImageRef, error) {
	if e == nil {
		return nil, errors.New("editor cannot be nil")
	}
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, fmt.Errorf("slide index %d out of range [0,%d)", slideIndex, len(e.slides))
	}
	rels, err := e.slideRelationships(e.slides[slideIndex].Part)
	if err != nil {
		return nil, err
	}

	out := make([]common.SlideImageRef, 0)
	for _, rel := range rels {
		if rel.Type != common.RelTypeImage {
			continue
		}
		out = append(out, common.SlideImageRef{Index: len(out), RelID: rel.ID, Target: rel.Target})
	}
	return out, nil
}

// SwapImageByIndex replaces one slide image relationship target by image index.
func (e *PresentationEditor) SwapImageByIndex(slideIndex, imageIndex int, data []byte, format string) error {
	images, err := e.ListSlideImages(slideIndex)
	if err != nil {
		return err
	}
	if imageIndex < 0 || imageIndex >= len(images) {
		return fmt.Errorf("image index %d out of range (found %d images)", imageIndex, len(images))
	}
	return e.swapImageByRelID(slideIndex, images[imageIndex].RelID, data, format)
}

// SwapImageByRelID replaces one slide image relationship target by relationship ID.
func (e *PresentationEditor) SwapImageByRelID(slideIndex int, relID string, data []byte, format string) error {
	if strings.TrimSpace(relID) == "" {
		return errors.New("relationship id cannot be empty")
	}
	return e.swapImageByRelID(slideIndex, relID, data, format)
}

func (e *PresentationEditor) swapImageByRelID(slideIndex int, relID string, data []byte, format string) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range [0,%d)", slideIndex, len(e.slides))
	}
	slidePart := e.slides[slideIndex].Part
	relsPath := common.SlideRelsPartName(slidePart)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return fmt.Errorf("missing slide relationships part %q", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return fmt.Errorf("parse %s: %w", relsPath, err)
	}

	updated := false
	for i := range rels {
		if rels[i].ID == relID && rels[i].Type == common.RelTypeImage {
			updated = true
			break
		}
	}
	if !updated {
		return fmt.Errorf("image relationship %q not found on slide %d", relID, slideIndex)
	}
	newPart, err := e.RegisterImage(data, format)
	if err != nil {
		return err
	}
	for i := range rels {
		if rels[i].ID == relID && rels[i].Type == common.RelTypeImage {
			rels[i].Target = common.MakeRelativePath(slidePart, newPart)
			break
		}
	}

	rendered := renderRelationshipsXML(rels)
	e.parts.Set(relsPath, []byte(rendered))
	return nil
}
