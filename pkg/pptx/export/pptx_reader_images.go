package export

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// attachSlideImages appends all embedded images for the given slide index.
func attachSlideImages(sc *elements.SlideContent, slideImages [][]SlideImage, idx int) {
	if idx >= len(slideImages) {
		return
	}
	for _, img := range slideImages[idx] {
		sc.Images = append(sc.Images, shapes.Image{
			Data:     img.Bytes,
			Format:   img.Format,
			X:        styling.Emu(img.X),
			Y:        styling.Emu(img.Y),
			CX:       styling.Emu(img.CX),
			CY:       styling.Emu(img.CY),
			Rotation: img.Rotation,
			Crop: shapes.ImageCrop{
				Left:   img.CropLeft,
				Right:  img.CropRight,
				Top:    img.CropTop,
				Bottom: img.CropBottom,
			},
			FlipH:        img.FlipH,
			FlipV:        img.FlipV,
			Shadow:       img.Shadow,
			Reflection:   img.Reflection,
			AltText:      img.AltText,
			IsDecorative: img.IsDecorative,

			InnerShadow:  img.InnerShadow,
			Glow:         img.Glow,
			SoftEdges:    img.SoftEdges,
			GlowSpec:     readerImageGlow(img),
			BlurSpec:     readerImageBlur(img),
			SoftEdgeSpec: readerImageSoftEdge(img),
		})
	}
}

// readerImageGlow keeps the glow radius read back from the package, so a
// round trip does not fall back to the default radius.
func readerImageGlow(img SlideImage) *shapes.ShapeGlow {
	if !img.Glow || img.GlowRadiusEmu <= 0 {
		return nil
	}
	return &shapes.ShapeGlow{RadiusEmu: img.GlowRadiusEmu}
}

func readerImageBlur(img SlideImage) *shapes.ShapeBlur {
	if !img.Blur {
		return nil
	}
	return &shapes.ShapeBlur{RadiusEmu: img.BlurRadiusEmu}
}

func readerImageSoftEdge(img SlideImage) *shapes.ShapeSoftEdge {
	if !img.SoftEdges || img.SoftEdgeRadiusEmu <= 0 {
		return nil
	}
	return &shapes.ShapeSoftEdge{RadiusEmu: img.SoftEdgeRadiusEmu}
}
