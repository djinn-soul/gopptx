package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

type (
	// Image describes one image placement.
	Image = shapes.Image
	// ImageCrop defines cropping details for an image.
	ImageCrop = shapes.ImageCrop
)

// NewImage creates a new image descriptor with a local file path.
func NewImage(path string, x, y, cx, cy styling.Length) shapes.Image {
	return media.NewImage(path, x, y, cx, cy)
}

// NewImageFromBytes creates a new image descriptor with raw data.
func NewImageFromBytes(data []byte, format string, x, y, cx, cy styling.Length) shapes.Image {
	return media.NewImageFromBytes(data, format, x, y, cx, cy)
}

// NewImageFromBase64 creates a new image descriptor with base64 encoded data.
func NewImageFromBase64(b64 string, format string, x, y, cx, cy styling.Length) (shapes.Image, error) {
	return media.NewImageFromBase64(b64, format, x, y, cx, cy)
}

// NewImageFromURL creates a new image descriptor with a remote URL.
func NewImageFromURL(url string, x, y, cx, cy styling.Length) shapes.Image {
	return media.NewImageFromURL(url, x, y, cx, cy)
}
