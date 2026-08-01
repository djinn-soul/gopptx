package common

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

// Image fit modes for placing a picture inside a fixed box such as a
// placeholder. The OOXML default is stretch, which is what a bare
// <a:stretch><a:fillRect/></a:stretch> does: the image is distorted to the box.
const (
	// ImageFitStretch fills the box exactly, ignoring the image aspect ratio.
	ImageFitStretch = "stretch"
	// ImageFitContain shrinks the image inside the box, preserving aspect ratio
	// and centering the letterboxed result. The box is not filled.
	ImageFitContain = "contain"
	// ImageFitCover fills the box, preserving aspect ratio by cropping the
	// overflowing axis evenly on both sides.
	ImageFitCover = "cover"
)

// cropPercentDenominator is the OOXML fixed-point unit for a:srcRect insets:
// 100000 means 100%.
const cropPercentDenominator = 100000

// ImageFit is the geometry to render for one fitted image: the box to draw in
// and, for cover, the crop insets that hide the overflow.
type ImageFit struct {
	X    int64
	Y    int64
	CX   int64
	CY   int64
	Crop *pptxxml.ImageCropRef
}

// NormalizeImageFit maps a caller-supplied fit mode onto a supported constant.
// An empty mode means stretch, which is the historical behavior.
func NormalizeImageFit(fit string) (string, error) {
	switch fit {
	case "", ImageFitStretch:
		return ImageFitStretch, nil
	case ImageFitContain, ImageFitCover:
		return fit, nil
	default:
		return "", fmt.Errorf(
			"image fit must be one of %q, %q or %q, got %q",
			ImageFitStretch, ImageFitContain, ImageFitCover, fit,
		)
	}
}

// FitImageToBox places an image of imgW x imgH pixels inside the box at
// (boxX, boxY) sized boxCX x boxCY EMU, according to fit.
//
// Stretch returns the box unchanged, which is what every caller got before fit
// modes existed. Contain and Cover both need real pixel dimensions; a zero or
// negative dimension is an error rather than a silent fall back to stretch,
// because silently distorting the image is the bug being fixed.
func FitImageToBox(fit string, boxX, boxY, boxCX, boxCY int64, imgW, imgH int) (ImageFit, error) {
	mode, err := NormalizeImageFit(fit)
	if err != nil {
		return ImageFit{}, err
	}

	box := ImageFit{X: boxX, Y: boxY, CX: boxCX, CY: boxCY}
	if mode == ImageFitStretch {
		return box, nil
	}
	if boxCX <= 0 || boxCY <= 0 {
		return ImageFit{}, fmt.Errorf(
			"image fit %q needs a sized box, got cx=%d cy=%d", mode, boxCX, boxCY,
		)
	}
	if imgW <= 0 || imgH <= 0 {
		return ImageFit{}, fmt.Errorf(
			"image fit %q needs the image dimensions, got %dx%d", mode, imgW, imgH,
		)
	}

	if mode == ImageFitContain {
		return containImage(box, int64(imgW), int64(imgH)), nil
	}
	return coverImage(box, int64(imgW), int64(imgH)), nil
}

// containImage shrinks the drawn extent to the larger of the two axis fits and
// centers it, leaving the box partly empty on the other axis.
func containImage(box ImageFit, imgW, imgH int64) ImageFit {
	// Compare imgW/imgH against boxCX/boxCY without floating point:
	// the image is relatively wider when imgW*boxCY > imgH*boxCX.
	fitted := box
	if imgW*box.CY > imgH*box.CX {
		// Width-bound: keep the box width, derive the height.
		fitted.CY = divRound(box.CX*imgH, imgW)
	} else {
		// Height-bound: keep the box height, derive the width.
		fitted.CX = divRound(box.CY*imgW, imgH)
	}
	fitted.X = box.X + (box.CX-fitted.CX)/2
	fitted.Y = box.Y + (box.CY-fitted.CY)/2
	return fitted
}

// coverImage keeps the box and crops the overflowing axis evenly, which is what
// PowerPoint's own "Fill" does to a picture placeholder.
func coverImage(box ImageFit, imgW, imgH int64) ImageFit {
	filled := box
	switch {
	case imgW*box.CY > imgH*box.CX:
		// Image is relatively wider than the box: crop left and right.
		// Visible fraction of the width is (boxCX/boxCY) / (imgW/imgH).
		visible := divRound(box.CX*imgH*cropPercentDenominator, box.CY*imgW)
		inset := (cropPercentDenominator - visible) / 2
		if inset > 0 {
			filled.Crop = &pptxxml.ImageCropRef{Left: inset, Right: inset}
		}
	case imgW*box.CY < imgH*box.CX:
		// Image is relatively taller than the box: crop top and bottom.
		visible := divRound(box.CY*imgW*cropPercentDenominator, box.CX*imgH)
		inset := (cropPercentDenominator - visible) / 2
		if inset > 0 {
			filled.Crop = &pptxxml.ImageCropRef{Top: inset, Bottom: inset}
		}
	}
	// Equal aspect ratios need no crop at all.
	return filled
}

// divRound divides two positive int64 values, rounding to nearest.
func divRound(numerator, denominator int64) int64 {
	if denominator == 0 {
		return 0
	}
	return (numerator + denominator/2) / denominator
}
