package editor

import (
	"crypto/sha256"
	"encoding/hex"
	"image"
	"net/http"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var embeddedImageRelPattern = regexp.MustCompile(`r:embed="([^"]+)"`)

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
