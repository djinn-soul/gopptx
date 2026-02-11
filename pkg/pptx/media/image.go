package media

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// NewImage creates a new image descriptor with a local file path.
func NewImage(path string, x, y, cx, cy int64) shapes.Image {
	return shapes.NewImage(path, x, y, cx, cy)
}

// NewImageFromBytes creates a new image descriptor with raw data.
func NewImageFromBytes(data []byte, format string, x, y, cx, cy int64) shapes.Image {
	return shapes.NewImageFromBytes(data, format, x, y, cx, cy)
}

// NewImageFromBase64 creates a new image descriptor with base64 encoded data.
func NewImageFromBase64(b64 string, format string, x, y, cx, cy int64) (shapes.Image, error) {
	return shapes.NewImageFromBase64(b64, format, x, y, cx, cy)
}

// NewImageFromURL creates a new image descriptor with a remote URL.
func NewImageFromURL(url string, x, y, cx, cy int64) shapes.Image {
	return shapes.NewImageFromURL(url, x, y, cx, cy)
}
