package elements

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// SlideBackgroundType defines the filling method for a slide background.
type SlideBackgroundType string

const (
	SlideBackgroundSolid    SlideBackgroundType = "solid"
	SlideBackgroundGradient SlideBackgroundType = "gradient"
	SlideBackgroundPicture  SlideBackgroundType = "picture"
)

// SlideBackground defines how a slide's background is rendered.
type SlideBackground struct {
	Type         SlideBackgroundType
	SolidFill    *shapes.ShapeFill
	GradientFill *shapes.ShapeGradientFill
	PictureFill  *shapes.Image
}

// NewSolidBackground creates a solid color background.
func NewSolidBackground(color string) SlideBackground {
	fill := shapes.NewShapeFill(color)
	return SlideBackground{
		Type:      SlideBackgroundSolid,
		SolidFill: &fill,
	}
}

// NewGradientBackground creates a background with a gradient fill.
func NewGradientBackground(gradient shapes.ShapeGradientFill) SlideBackground {
	return SlideBackground{
		Type:         SlideBackgroundGradient,
		GradientFill: &gradient,
	}
}

// NewPictureBackground creates a background using an image.
func NewPictureBackground(img shapes.Image) SlideBackground {
	return SlideBackground{
		Type:        SlideBackgroundPicture,
		PictureFill: &img,
	}
}

// Validate checks if the background configuration is valid.
func (b SlideBackground) Validate() error {
	switch b.Type {
	case SlideBackgroundSolid:
		if b.SolidFill == nil {
			return fmt.Errorf("solid background must have fill properties")
		}
		return b.SolidFill.Validate()
	case SlideBackgroundGradient:
		if b.GradientFill == nil {
			return fmt.Errorf("gradient background must have fill properties")
		}
		return b.GradientFill.Validate()
	case SlideBackgroundPicture:
		if b.PictureFill == nil {
			return fmt.Errorf("picture background must have an image")
		}
		// Minimal validation for path/data
		if b.PictureFill.Path == "" && len(b.PictureFill.Data) == 0 {
			return fmt.Errorf("picture background image must have a path or data")
		}
	case "":
		return nil
	default:
		return fmt.Errorf("invalid background type %q", b.Type)
	}
	return nil
}

// NormalizeSlideBackgroundType ensures the type string is canonical.
func NormalizeSlideBackgroundType(t string) SlideBackgroundType {
	lower := strings.ToLower(strings.TrimSpace(t))
	switch lower {
	case "solid":
		return SlideBackgroundSolid
	case "gradient":
		return SlideBackgroundGradient
	case "picture", "image":
		return SlideBackgroundPicture
	default:
		return SlideBackgroundSolid
	}
}
