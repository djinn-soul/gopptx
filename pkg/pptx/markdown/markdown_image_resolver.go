package markdown

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func (p *markdownASTParser) resolveMarkdownImage(image markdownImage) (shapes.Image, error) {
	dest := strings.TrimSpace(image.Dest)
	if dest == "" {
		return shapes.Image{}, errors.New("markdown image destination cannot be empty")
	}

	x, y, cx, cy := p.nextEmbeddedImageFrame()
	alt := strings.TrimSpace(image.Alt)
	if alt == "" {
		alt = "markdown image"
	}

	if isHTTPImageURL(dest) {
		return shapes.NewImageFromURL(dest, x, y, cx, cy).WithAltText(alt), nil
	}
	if strings.HasPrefix(strings.ToLower(dest), "data:") {
		data, format, err := decodeMarkdownDataURIImage(dest)
		if err != nil {
			return shapes.Image{}, err
		}
		return shapes.NewImageFromBytes(data, format, x, y, cx, cy).WithAltText(alt), nil
	}

	resolvedPath := dest
	if p.options.BaseDir != "" && !filepath.IsAbs(dest) {
		resolvedPath = filepath.Join(p.options.BaseDir, dest)
	}
	resolvedPath = filepath.Clean(resolvedPath)
	return shapes.NewImage(resolvedPath, x, y, cx, cy).WithAltText(alt), nil
}

func isHTTPImageURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil {
		return false
	}
	scheme := strings.ToLower(parsed.Scheme)
	return scheme == "http" || scheme == "https"
}

func decodeMarkdownDataURIImage(dataURI string) ([]byte, string, error) {
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 {
		return nil, "", errors.New("invalid data URI image payload")
	}

	meta := strings.TrimSpace(parts[0])
	payload := strings.TrimSpace(parts[1])
	if !strings.HasPrefix(strings.ToLower(meta), "data:") {
		return nil, "", errors.New("invalid data URI image payload")
	}
	if !strings.Contains(strings.ToLower(meta), ";base64") {
		return nil, "", errors.New("data URI image payload must be base64 encoded")
	}

	mimeType := strings.TrimPrefix(strings.SplitN(meta, ";", 2)[0], "data:")
	format, ok := dataURIImageMimeToFormat(strings.ToLower(mimeType))
	if !ok {
		return nil, "", fmt.Errorf("unsupported data URI image mime type %q", mimeType)
	}

	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, "", fmt.Errorf("invalid data URI base64 payload: %w", err)
	}
	if len(data) == 0 {
		return nil, "", errors.New("data URI image payload cannot be empty")
	}
	return data, format, nil
}

func dataURIImageMimeToFormat(mimeType string) (string, bool) {
	switch mimeType {
	case "image/png":
		return "png", true
	case "image/jpeg":
		return "jpg", true
	case "image/jpg":
		return "jpg", true
	case "image/gif":
		return "gif", true
	case "image/bmp":
		return "bmp", true
	case "image/tiff":
		return "tiff", true
	default:
		// Keep alignment with supported media extensions.
		ext := media.NormalizeExtension(mimeType)
		switch ext {
		case "png", "jpg", "jpeg", "gif", "bmp", "tiff":
			return ext, true
		default:
			return "", false
		}
	}
}
