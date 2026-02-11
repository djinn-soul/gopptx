package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

type (
	// Image describes one image placement.
	Image = shapes.Image
	// ImageCrop defines cropping details for an image.
	ImageCrop = shapes.ImageCrop
)

// NewImage creates a new image descriptor with a local file path.
var NewImage = media.NewImage

// NewImageFromBytes creates a new image descriptor with raw data.
var NewImageFromBytes = media.NewImageFromBytes

// NewImageFromBase64 creates a new image descriptor with base64 encoded data.
var NewImageFromBase64 = media.NewImageFromBase64

// NewImageFromURL creates a new image descriptor with a remote URL.
var NewImageFromURL = media.NewImageFromURL
