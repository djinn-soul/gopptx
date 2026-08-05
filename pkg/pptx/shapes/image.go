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
	// ZOrder is the picture's position in the slide shape tree. Lower paints
	// first (further back). Zero for images not read from a PPTX.
	ZOrder    int
	Path      string
	SourceURL string
	Data      []byte
	Format    string
	RelID     string
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

	// InnerShadow draws the shadow inside the picture edges instead of outside.
	InnerShadow bool
	// Glow enables a glow halo with PowerPoint's default radius and color.
	Glow bool
	// SoftEdges feathers the picture border with the default radius.
	SoftEdges bool
	// GlowSpec overrides the default glow radius and color.
	GlowSpec *ShapeGlow
	// BlurSpec blurs the picture; nil means no blur.
	BlurSpec *ShapeBlur
	// SoftEdgeSpec overrides the default soft-edge radius.
	SoftEdgeSpec *ShapeSoftEdge
	// ReflectionSpec overrides the default reflection blur and distance.
	ReflectionSpec *ShapeReflection
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

// WithInnerShadow adds an inner shadow effect to the image.
func (img Image) WithInnerShadow(enabled bool) Image {
	img.InnerShadow = enabled
	return img
}

// WithGlow adds a glow effect with PowerPoint's default radius and color.
func (img Image) WithGlow(enabled bool) Image {
	img.Glow = enabled
	return img
}

// WithSoftEdges feathers the image border with the default radius.
func (img Image) WithSoftEdges(enabled bool) Image {
	img.SoftEdges = enabled
	return img
}

// WithGlowSpec sets detailed glow settings and enables the glow flag.
func (img Image) WithGlowSpec(glow *ShapeGlow) Image {
	img.Glow = glow != nil
	img.GlowSpec = glow
	return img
}

// WithBlurSpec blurs the image by the given radius. A nil spec removes the blur.
func (img Image) WithBlurSpec(blur *ShapeBlur) Image {
	img.BlurSpec = blur
	return img
}

// WithSoftEdgeSpec sets detailed soft-edge settings and enables the flag.
func (img Image) WithSoftEdgeSpec(softEdge *ShapeSoftEdge) Image {
	img.SoftEdges = softEdge != nil
	img.SoftEdgeSpec = softEdge
	return img
}

// WithReflectionSpec sets detailed reflection settings and enables the flag.
func (img Image) WithReflectionSpec(reflection *ShapeReflection) Image {
	img.Reflection = reflection != nil
	img.ReflectionSpec = reflection
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

	if img.Path == "" && len(img.Data) == 0 && img.SourceURL == "" && img.RelID == "" {
		return fmt.Errorf("slide %d image %d has no source (Path, Data, SourceURL, or RelID)", slideIndex, imageIndex)
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
