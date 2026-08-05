package export

import (
	"bytes"
	"fmt"
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

const (
	imageShadowOffsetPt      = 4.0
	imageReflectionGapPt     = 3.0
	imageReflectionScale     = 0.35
	imageReflectionMaxAlpha  = 0.28
	imageShadowDefaultAlpha  = 0.25
	imageShadowDefaultMode   = "Normal"
	imageReflectionBlendMode = "Normal"

	// A srcRect keeping less than this fraction of an edge is treated as
	// keeping nothing, which also guards the division by the visible fraction.
	minVisibleCropFraction = 0.001
)

func renderPDFImageWithEffects(pdf *gopdf.GoPdf, img shapes.Image) error {
	if len(img.Data) == 0 {
		return nil
	}
	// gopdf does not support EMF/WMF (vector formats) — skip them.
	if img.Format == formatEMF || img.Format == "wmf" {
		return nil
	}

	x := emuToPt(img.X.Emu())
	y := emuToPt(img.Y.Emu())
	w := emuToPt(img.CX.Emu())
	h := emuToPt(img.CY.Emu())
	if w <= 0 || h <= 0 {
		return nil
	}

	imgHolder, err := gopdf.ImageHolderByReader(bytes.NewReader(img.Data))
	if err != nil {
		return fmt.Errorf("create image holder: %w", err)
	}

	if img.Shadow {
		renderPDFImageShadow(pdf, x, y, w, h)
	}
	if err := drawPDFImage(pdf, imgHolder, x, y, w, h, img, nil); err != nil {
		return err
	}
	if img.Reflection {
		renderPDFImageReflection(pdf, imgHolder, x, y, w, h, img)
	}
	return nil
}

// drawPDFImage paints one picture into the box at (x, y, w, h).
//
// Cropped and uncropped pictures take the same path: the crop decides where the
// whole picture is laid down (an uncropped one simply fills the box), a clip
// keeps it inside the box, and rotation turns the page about the box centre.
// gopdf's own DegreeAngle would turn the picture about its own centre, which is
// the box centre only while nothing is cropped — so rotation is done here for
// both, and the two cases cannot drift apart.
func drawPDFImage(
	pdf *gopdf.GoPdf,
	holder gopdf.ImageHolder,
	x, y, w, h float64,
	img shapes.Image,
	transparency *gopdf.Transparency,
) error {
	full, ok := croppedImagePlacement(img.Crop, x, y, w, h)
	// A srcRect that keeps nothing (or is malformed) leaves no picture to draw.
	if !ok {
		return nil
	}

	rotated := math.Abs(img.Rotation) > nearZeroEpsilon
	if rotated {
		pdf.Rotate(img.Rotation, x+w/2, y+h/2)
	}

	// Clipping is only needed when the picture is larger than its box, which is
	// exactly when something is cropped away.
	clipped := hasImageCrop(img.Crop)
	if clipped {
		pdf.SaveGraphicsState()
		pdf.ClipPolygon([]gopdf.Point{
			{X: x, Y: y},
			{X: x + w, Y: y},
			{X: x + w, Y: y + h},
			{X: x, Y: y + h},
		})
	}
	err := pdf.ImageByHolderWithOptions(holder, gopdf.ImageOptions{
		X:              full.X,
		Y:              full.Y,
		Rect:           &gopdf.Rect{W: full.W, H: full.H},
		HorizontalFlip: img.FlipH,
		VerticalFlip:   img.FlipV,
		Transparency:   transparency,
	})
	if clipped {
		pdf.RestoreGraphicsState()
	}
	if rotated {
		pdf.RotateReset()
	}
	if err != nil {
		return fmt.Errorf("draw image: %w", err)
	}
	return nil
}

// hasImageCrop reports whether a:srcRect asks for anything to be cut away.
// Negative insets are legal in OOXML and pad the picture instead of cropping it,
// so they count too.
func hasImageCrop(crop shapes.ImageCrop) bool {
	return crop.Left != 0 || crop.Right != 0 || crop.Top != 0 || crop.Bottom != 0
}

// imageRect is a placement rectangle in points.
type imageRect struct {
	X, Y, W, H float64
}

// croppedImagePlacement is where the *whole* picture goes so that the part
// a:srcRect keeps lands exactly on the shape box at (x, y, w, h). It reports
// false when the crop leaves nothing visible.
func croppedImagePlacement(crop shapes.ImageCrop, x, y, w, h float64) (imageRect, bool) {
	visibleW := 1 - crop.Left - crop.Right
	visibleH := 1 - crop.Top - crop.Bottom
	if visibleW <= minVisibleCropFraction || visibleH <= minVisibleCropFraction {
		return imageRect{}, false
	}
	fullW := w / visibleW
	fullH := h / visibleH
	return imageRect{
		X: x - crop.Left*fullW,
		Y: y - crop.Top*fullH,
		W: fullW,
		H: fullH,
	}, true
}

func renderPDFImageShadow(pdf *gopdf.GoPdf, x, y, w, h float64) {
	shadowAlpha, err := gopdf.NewTransparency(imageShadowDefaultAlpha, imageShadowDefaultMode)
	if err == nil {
		_ = pdf.SetTransparency(shadowAlpha)
	}
	pdf.SetFillColor(0, 0, 0)
	pdf.RectFromUpperLeftWithStyle(
		x+imageShadowOffsetPt,
		y+imageShadowOffsetPt,
		w,
		h,
		"F",
	)
	if err == nil {
		pdf.ClearTransparency()
	}
}

func renderPDFImageReflection(
	pdf *gopdf.GoPdf,
	holder gopdf.ImageHolder,
	x, y, w, h float64,
	img shapes.Image,
) {
	refH := h * imageReflectionScale
	if refH <= 1 {
		return
	}
	alpha, err := gopdf.NewTransparency(imageReflectionMaxAlpha, imageReflectionBlendMode)
	if err != nil {
		return
	}
	refImg := img
	refImg.FlipV = !img.FlipV
	_ = drawPDFImage(
		pdf,
		holder,
		x,
		y+h+imageReflectionGapPt,
		w,
		refH,
		refImg,
		&alpha,
	)
}
