package editor

import (
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/fonts"
)

// defaultFitShapeFontSizePt is the size the text is measured at when neither the
// caller nor the shape's runs name one. It matches PowerPoint's default body
// size for a new text box.
const defaultFitShapeFontSizePt = 18.0

// FitShapeToTextRequest asks for a shape to be grown so its text fits.
type FitShapeToTextRequest struct {
	// FontPath is a TrueType file used to measure the text. Without it the text
	// cannot be measured and the shape cannot be sized.
	FontPath string
	// FontSizePt is the size the text is measured at. Zero means the largest
	// size found on the shape's runs, or defaultFitShapeFontSizePt.
	FontSizePt float64
	// MinHeightEmu and MaxHeightEmu bound the result. Zero means unbounded.
	MinHeightEmu int
	MaxHeightEmu int
	// Shrink lets the shape get shorter than it is. By default the shape only
	// grows, which is what "resize shape to fit text" does in PowerPoint.
	Shrink bool
}

// FitShapeToTextResult reports the geometry the fit produced.
type FitShapeToTextResult struct {
	HeightEmu int  `json:"height_emu"`
	LineCount int  `json:"line_count"`
	Resized   bool `json:"resized"`
}

// FitShapeToText grows a shape until its text fits, the opposite of
// FitShapeText: the font size is left alone and the box moves.
//
// This is the "autosize the shape to the given text" that upstream issue #970
// asks for, where fit_text only ever shrank the text. The frame is switched to
// spAutoFit so PowerPoint keeps the shape sized to the text after an edit.
func (e *PresentationEditor) FitShapeToText(
	slideIndex, shapeID int,
	request FitShapeToTextRequest,
) (FitShapeToTextResult, error) {
	shape, err := e.getShapeForTextOps(slideIndex, shapeID)
	if err != nil {
		return FitShapeToTextResult{}, err
	}
	if request.FontPath == "" {
		return FitShapeToTextResult{}, errors.New(
			"fitting a shape to its text needs a font file to measure with",
		)
	}

	insetX, insetY := shapeTextInsetsEmu(shape.TextFrame)
	widthPt := (float64(shape.W) - insetX) / emuPerPoint
	if widthPt <= 0 {
		return FitShapeToTextResult{}, fmt.Errorf(
			"shape %d is too narrow to hold any text (%d EMU wide)", shapeID, shape.W,
		)
	}

	measurer, err := fonts.NewTextMeasurer(request.FontPath)
	if err != nil {
		return FitShapeToTextResult{}, err
	}
	textHeightPt, lineCount, err := measurer.WrappedHeightPt(
		shape.Text, widthPt, fitShapeFontSizePt(request, shape),
	)
	if err != nil {
		return FitShapeToTextResult{}, err
	}

	height := int(textHeightPt*emuPerPoint) + int(insetY)
	height = boundFittedHeight(height, shape.H, request)

	// The frame is switched to spAutoFit even when the height did not move: that
	// is what keeps the shape sized to its text through later edits, and it is
	// the contract of this call. Returning early here left the shape without it.
	wordWrap := true
	autoFitType := "shape"
	update := common.ShapeUpdate{
		TextFrame: &common.TextFrame{WordWrap: &wordWrap, AutoFitType: &autoFitType},
	}
	resized := height != shape.H
	if resized {
		update.H = &height
	}
	if updateErr := e.UpdateShape(slideIndex, shapeID, update); updateErr != nil {
		return FitShapeToTextResult{}, updateErr
	}
	return FitShapeToTextResult{HeightEmu: height, LineCount: lineCount, Resized: resized}, nil
}

// fitShapeFontSizePt picks the size the text is measured at: the caller's, else
// the largest size any run carries, else the default body size.
func fitShapeFontSizePt(request FitShapeToTextRequest, shape parsedShape) float64 {
	if request.FontSizePt > 0 {
		return request.FontSizePt
	}
	largest := 0.0
	for _, run := range shape.Runs {
		if run.SizePt != nil && *run.SizePt > largest {
			largest = *run.SizePt
		}
	}
	if largest > 0 {
		return largest
	}
	return defaultFitShapeFontSizePt
}

// boundFittedHeight applies the caller's bounds and the grow-only default.
func boundFittedHeight(height, currentHeight int, request FitShapeToTextRequest) int {
	if request.MinHeightEmu > 0 && height < request.MinHeightEmu {
		height = request.MinHeightEmu
	}
	if request.MaxHeightEmu > 0 && height > request.MaxHeightEmu {
		height = request.MaxHeightEmu
	}
	if !request.Shrink && height < currentHeight {
		return currentHeight
	}
	return height
}
