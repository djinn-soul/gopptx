package editor

import (
	"crypto/sha256"
	"encoding/hex"
	"image"
	"net/http"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Metafile magic numbers. http.DetectContentType does not know EMF or WMF, so
// they have to be sniffed explicitly or both report application/octet-stream.
const (
	metafileMagicLen   = 4
	emfRecordHeader    = "\x01\x00\x00\x00"
	emfSignature       = " EMF" // at offset 40 of EMR_HEADER
	emfSignatureOffset = 40
	emfSignatureEnd    = emfSignatureOffset + len(emfSignature)
	wmfPlaceableMagic  = "\xd7\xcd\xc6\x9a"
	wmfDiskHeader      = "\x01\x00\x09\x00"
	wmfMemoryHeader    = "\x02\x00\x09\x00"
)

func buildImageMetadata(data []byte, cfg image.Config, format string) *common.ImageMetadata {
	return &common.ImageMetadata{
		Width:       cfg.Width,
		Height:      cfg.Height,
		Format:      strings.ToLower(strings.TrimSpace(format)),
		ContentType: imageContentType(data, format),
		Hash:        imageSHA256Hex(data),
		Data:        data,
	}
}

func imageSHA256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func imageContentType(data []byte, format string) string {
	if metafile := metafileContentType(data); metafile != "" {
		return metafile
	}

	contentType := strings.TrimSpace(http.DetectContentType(data))
	if contentType != "" && contentType != "application/octet-stream" {
		return contentType
	}

	if mime := contentTypeForImageFormat(format); mime != "" {
		return mime
	}
	return contentType
}

func contentTypeForImageFormat(format string) string {
	switch normalizeImageFormatHint(format) {
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
	case formatEMF:
		return mimeEMF
	case formatWMF:
		return mimeWMF
	case "wdp", "hdp":
		return mimeMSPhoto
	case "svg", "image/svg+xml":
		return "image/svg+xml"
	default:
		return ""
	}
}

// metafileContentType returns the EMF or WMF content type when data carries the
// corresponding magic number, or "" when it is neither.
func metafileContentType(data []byte) string {
	if len(data) < metafileMagicLen {
		return ""
	}
	switch string(data[:metafileMagicLen]) {
	case wmfPlaceableMagic, wmfDiskHeader, wmfMemoryHeader:
		return mimeWMF
	case emfRecordHeader:
		if len(data) >= emfSignatureEnd && string(data[emfSignatureOffset:emfSignatureEnd]) == emfSignature {
			return mimeEMF
		}
	}
	return ""
}
