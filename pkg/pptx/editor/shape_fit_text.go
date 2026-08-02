package editor

import (
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/fonts"
)

// emuPerPoint converts EMUs to points (1 pt = 12700 EMU).
const emuPerPoint = 12700.0

// Default text insets in EMUs, matching OOXML's implicit a:bodyPr values:
// 0.1" left/right and 0.05" top/bottom.
const (
	defaultTextInsetLREmu = 91440.0
	defaultTextInsetTBEmu = 45720.0
)

// FitTextRequest asks for a shape's text to be shrunk to fit its box.
type FitTextRequest struct {
	// FontPath is a TrueType file used to measure the text. Without it the text
	// cannot be measured and only the autofit flags can be set.
	FontPath string
	// MaxSizePt is the largest size to try; zero means the fitter's default.
	MaxSizePt float64
	// MinSizePt is the smallest size to try; zero means the fitter's floor.
	MinSizePt float64
	// WordWrap turns wrapping on, which is what makes shrinking meaningful.
	WordWrap bool
}

// FitTextResult reports what the fit produced.
type FitTextResult struct {
	FontSizePt float64 `json:"font_size_pt"`
	Fits       bool    `json:"fits"`
	LineCount  int     `json:"line_count"`
	Measured   bool    `json:"measured"`
}

// FitShapeText shrinks a shape's text to the largest size that fits its box.
//
// With a font file the text is measured against real glyph advances and every
// run is rewritten at the chosen size. Without one, only the autofit flags are
// set and Measured is false, so callers can tell a real fit from a hint to
// PowerPoint.
//
// Text that cannot be broken -- a long unspaced number, say -- is placed on its
// own line and allowed to overflow rather than raising, which is the crash the
// upstream issue reports.
func (e *PresentationEditor) FitShapeText(
	slideIndex, shapeID int,
	request FitTextRequest,
) (FitTextResult, error) {
	shape, err := e.getShapeForTextOps(slideIndex, shapeID)
	if err != nil {
		return FitTextResult{}, err
	}

	wordWrap := request.WordWrap
	autoFitType := "normal"
	frame := common.TextFrame{WordWrap: &wordWrap, AutoFitType: &autoFitType}

	if request.FontPath == "" {
		if updateErr := e.UpdateShape(slideIndex, shapeID, common.ShapeUpdate{
			TextFrame: &frame,
		}); updateErr != nil {
			return FitTextResult{}, updateErr
		}
		return FitTextResult{Measured: false}, nil
	}

	measurer, err := fonts.NewTextMeasurer(request.FontPath)
	if err != nil {
		return FitTextResult{}, err
	}

	widthPt := (float64(shape.W) - 2*defaultTextInsetLREmu) / emuPerPoint
	heightPt := (float64(shape.H) - 2*defaultTextInsetTBEmu) / emuPerPoint
	if widthPt <= 0 || heightPt <= 0 {
		return FitTextResult{}, fmt.Errorf(
			"shape %d is too small to hold any text (%dx%d EMU)", shapeID, shape.W, shape.H,
		)
	}

	fit, err := measurer.FitText(fonts.FitRequest{
		Text:      shape.Text,
		WidthPt:   widthPt,
		HeightPt:  heightPt,
		MaxSizePt: request.MaxSizePt,
		MinSizePt: request.MinSizePt,
	})
	if err != nil {
		return FitTextResult{}, err
	}

	runs, err := e.fitTextRuns(slideIndex, shapeID, fit.FontSizePt)
	if err != nil {
		return FitTextResult{}, err
	}
	if updateErr := e.UpdateShape(slideIndex, shapeID, common.ShapeUpdate{
		TextFrame: &frame,
		Runs:      &runs,
	}); updateErr != nil {
		return FitTextResult{}, updateErr
	}

	return FitTextResult{
		FontSizePt: fit.FontSizePt,
		Fits:       fit.Fits,
		LineCount:  fit.LineCount,
		Measured:   true,
	}, nil
}

// fitTextRuns re-reads the shape's runs and stamps the fitted size on each,
// leaving every other run property alone.
func (e *PresentationEditor) fitTextRuns(
	slideIndex, shapeID int,
	sizePt float64,
) ([]common.TextRun, error) {
	shape, err := e.getShapeForTextOps(slideIndex, shapeID)
	if err != nil {
		return nil, err
	}
	if len(shape.Runs) == 0 {
		return nil, errors.New("shape has no text runs to fit")
	}
	runs := make([]common.TextRun, len(shape.Runs))
	copy(runs, shape.Runs)
	for i := range runs {
		size := sizePt
		runs[i].SizePt = &size
	}
	return runs, nil
}
