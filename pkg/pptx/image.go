package pptx

import (
	"encoding/base64"
	"fmt"
)

// ImageCrop defines cropping details for an image.
// Values are in percentage (0.0 to 1.0) or specific units depending on how we render.
// OOXML uses percentage (e.g. 100000 = 100%). We'll use 0-1 float for API and convert.
type ImageCrop struct {
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

// Image describes one image placement on a slide.
type Image struct {
	Path        string
	SourceURL   string
	Data        []byte
	Format      string // e.g. "png", "jpg" - required if Data is used
	X           int64
	Y           int64
	CX          int64
	CY          int64
	Rotation    float64
	Crop        ImageCrop
	FlipH       bool
	FlipV       bool
	Placeholder *Placeholder
}

// NewImage creates an image placement from a file path.
func NewImage(path string, x int64, y int64, cx int64, cy int64) Image {
	return Image{
		Path: path,
		X:    x,
		Y:    y,
		CX:   cx,
		CY:   cy,
	}
}

// NewImageFromBytes creates an image placement from raw bytes.
func NewImageFromBytes(data []byte, format string, x int64, y int64, cx int64, cy int64) Image {
	return Image{
		Data:   data,
		Format: format,
		X:      x,
		Y:      y,
		CX:     cx,
		CY:     cy,
	}
}

// NewImageFromBase64 creates an image placement from a base64 string.
func NewImageFromBase64(b64 string, format string, x int64, y int64, cx int64, cy int64) (Image, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return Image{}, fmt.Errorf("invalid base64 image data: %w", err)
	}
	return NewImageFromBytes(data, format, x, y, cx, cy), nil
}

// NewImageFromURL creates an image placement from a URL.
func NewImageFromURL(url string, x int64, y int64, cx int64, cy int64) Image {
	return Image{
		SourceURL: url,
		X:         x,
		Y:         y,
		CX:        cx,
		CY:        cy,
	}
}

// WithRotation adds rotation (degrees) to the image.
func (img Image) WithRotation(degrees float64) Image {
	img.Rotation = degrees
	return img
}

// WithCrop adds cropping to the image.
func (img Image) WithCrop(left, right, top, bottom float64) Image {
	img.Crop = ImageCrop{
		Left:   left,
		Right:  right,
		Top:    top,
		Bottom: bottom,
	}
	return img
}

// WithFlip adds horizontal/vertical flip.
func (img Image) WithFlip(horizontal, vertical bool) Image {
	img.FlipH = horizontal
	img.FlipV = vertical
	return img
}

// Validate checks the image for common constraints and required fields.
func (img Image) Validate(slideIndex int, imageIndex int) error {
	if img.Path == "" && len(img.Data) == 0 && img.SourceURL == "" {
		return fmt.Errorf("slide %d image %d has no source (Path, Data, or SourceURL)", slideIndex, imageIndex)
	}
	if len(img.Data) > 0 && img.Format == "" {
		return fmt.Errorf("slide %d image %d has Data but no Format", slideIndex, imageIndex)
	}
	if img.X < 0 || img.Y < 0 {
		return fmt.Errorf("slide %d image %d position cannot be negative", slideIndex, imageIndex)
	}
	if img.CX <= 0 || img.CY <= 0 {
		return fmt.Errorf("slide %d image %d size must be > 0", slideIndex, imageIndex)
	}
	return nil
}
