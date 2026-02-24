package shapes

import (
	"encoding/base64"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// ImageCrop defines cropping details for an image.
type ImageCrop struct {
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

// Image describes one image placement.
type Image struct {
	Path      string
	SourceURL string
	Data      []byte
	Format    string
	X         styling.Length
	Y         styling.Length
	CX        styling.Length
	CY        styling.Length

	Rotation     float64
	Crop         ImageCrop
	FlipH        bool
	FlipV        bool
	Shadow       bool
	Reflection   bool
	AltText      string
	IsDecorative bool
	Placeholder  *Placeholder
}

// NewImage creates an image placement.
func NewImage(path string, x, y, cx, cy styling.Length) Image {
	return Image{Path: path, X: x, Y: y, CX: cx, CY: cy}
}

// NewImageFromBytes creates an image placement from raw bytes.
func NewImageFromBytes(data []byte, format string, x, y, cx, cy styling.Length) Image {
	return Image{Data: data, Format: format, X: x, Y: y, CX: cx, CY: cy}
}

// NewImageFromBase64 creates an image placement from a base64 string.
func NewImageFromBase64(b64 string, format string, x, y, cx, cy styling.Length) (Image, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return Image{}, fmt.Errorf("invalid base64 image data: %w", err)
	}
	return NewImageFromBytes(data, format, x, y, cx, cy), nil
}

// NewImageFromURL creates an image placement from a URL.
func NewImageFromURL(url string, x, y, cx, cy styling.Length) Image {
	return Image{SourceURL: url, X: x, Y: y, CX: cx, CY: cy}
}

// WithShadow adds an outer shadow effect to the image.
func (img Image) WithShadow(enabled bool) Image {
	img.Shadow = enabled
	return img
}

// WithReflection adds a reflection effect to the image.
func (img Image) WithReflection(enabled bool) Image {
	img.Reflection = enabled
	return img
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

// WithAltText sets the alternative text for accessibility.
func (img Image) WithAltText(text string) Image {
	img.AltText = text
	return img
}

// WithDecorative marks the image as decorative (ignored by screen readers).
func (img Image) WithDecorative(enabled bool) Image {
	img.IsDecorative = enabled
	return img
}

// Validate checks the image for common constraints.
func (img Image) Validate(slideIndex, imageIndex int) error {
	if !img.IsDecorative && len(img.AltText) > common.MaxAltTextLength {
		return fmt.Errorf(
			"slide %d image %d alt text exceeds %d characters",
			slideIndex,
			imageIndex,
			common.MaxAltTextLength,
		)
	}

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
