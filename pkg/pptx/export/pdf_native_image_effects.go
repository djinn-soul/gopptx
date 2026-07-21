package export

import (
	"bytes"
	"fmt"

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

func drawPDFImage(
	pdf *gopdf.GoPdf,
	holder gopdf.ImageHolder,
	x, y, w, h float64,
	img shapes.Image,
	transparency *gopdf.Transparency,
) error {
	opts := gopdf.ImageOptions{
		X:              x,
		Y:              y,
		Rect:           &gopdf.Rect{W: w, H: h},
		DegreeAngle:    img.Rotation,
		HorizontalFlip: img.FlipH,
		VerticalFlip:   img.FlipV,
		Transparency:   transparency,
	}
	if err := pdf.ImageByHolderWithOptions(holder, opts); err != nil {
		return fmt.Errorf("draw image: %w", err)
	}
	return nil
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
