package editor

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	editormodcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/common"
)

// applyPlaceholderImageFit adjusts a placeholder picture so it keeps its aspect
// ratio inside the placeholder box.
//
// It runs after the target shape has been resolved because a placeholder that
// inherits its geometry from the layout has no bounds of its own until then,
// and a fit without a box is meaningless.
func (e *PresentationEditor) applyPlaceholderImageFit(
	slideIndex int,
	payload map[string]any,
	v *PayloadValidator,
	phSpec *pptxxml.PlaceholderOverrideSpec,
	shapeXML []byte,
) error {
	if phSpec.Image == nil {
		return nil
	}

	fit, err := editormodcommon.NormalizeImageFit(v.OptionalString(payload, "image_fit"))
	if err != nil {
		return NewBridgeError(ErrCodeInvalidValue, err.Error())
	}
	if fit == editormodcommon.ImageFitStretch {
		return nil
	}

	box, ok := placeholderImageBox(phSpec.Image, shapeXML)
	if !ok {
		return NewBridgeError(ErrCodeInvalidValue, fmt.Sprintf(
			"image_fit %q needs a sized placeholder: this placeholder inherits its "+
				"geometry from the layout, so pass explicit 'bounds'", fit,
		))
	}

	width, height, err := e.imagePixelSize(slideIndex, phSpec.Image.RelID)
	if err != nil {
		return err
	}

	fitted, err := editormodcommon.FitImageToBox(fit, box.X, box.Y, box.CX, box.CY, width, height)
	if err != nil {
		return NewBridgeError(ErrCodeInvalidValue, err.Error())
	}

	phSpec.Image.X = fitted.X
	phSpec.Image.Y = fitted.Y
	phSpec.Image.CX = fitted.CX
	phSpec.Image.CY = fitted.CY
	phSpec.Image.Crop = fitted.Crop
	return nil
}

// placeholderImageBox prefers the bounds the caller passed and falls back to
// the bounds already on the shape.
func placeholderImageBox(img *pptxxml.ImageRef, shapeXML []byte) (placeholderBounds, bool) {
	if img.CX > 0 && img.CY > 0 {
		return placeholderBounds{X: img.X, Y: img.Y, CX: img.CX, CY: img.CY}, true
	}
	return extractPlaceholderBounds(shapeXML)
}

// imagePixelSize reads the pixel dimensions of an already-registered image part.
func (e *PresentationEditor) imagePixelSize(slideIndex int, relID string) (int, int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, 0, NewBridgeError(ErrCodeInvalidValue, "slide index out of range")
	}
	data, partPath, err := e.imagePartData(e.slides[slideIndex].Part, relID)
	if err != nil {
		return 0, 0, err
	}
	config, _ := decodeImageConfig(data, partPath)
	if config.Width <= 0 || config.Height <= 0 {
		return 0, 0, NewBridgeError(ErrCodeInvalidValue, fmt.Sprintf(
			"could not read the dimensions of %s, so the image cannot be fitted", partPath,
		))
	}
	return config.Width, config.Height, nil
}
