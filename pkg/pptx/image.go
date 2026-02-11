package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type (
	// Image describes one image placement.
	Image = elements.Image
	// ImageCrop defines cropping details for an image.
	ImageCrop = elements.ImageCrop
)

func NewImage(path string, x, y, cx, cy int64) Image {
	return elements.NewImage(path, x, y, cx, cy)
}

func NewImageFromBytes(data []byte, format string, x, y, cx, cy int64) Image {
	return elements.NewImageFromBytes(data, format, x, y, cx, cy)
}

func NewImageFromBase64(b64 string, format string, x, y, cx, cy int64) (Image, error) {
	return elements.NewImageFromBase64(b64, format, x, y, cx, cy)
}

func NewImageFromURL(url string, x, y, cx, cy int64) Image {
	return elements.NewImageFromURL(url, x, y, cx, cy)
}
