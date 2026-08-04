package editor

import (
	"errors"
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
	tablemod "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/table"
)

// Picture fill modes. Stretch scales one copy of the image to the shape;
// tile repeats it at its natural size.
const (
	PictureFillStretch = "stretch"
	PictureFillTile    = "tile"
)

// pictureFillCropScale converts a crop fraction to the per-hundred-thousand
// units <a:srcRect> uses.
const pictureFillCropScale = 100000

// PictureFillRequest describes a picture (<a:blipFill>) fill to write onto a
// shape or a table cell. Exactly one of ImagePath and ImageData must be set.
type PictureFillRequest struct {
	ImagePath string
	ImageData string // base64
	// Mode is PictureFillStretch (default) or PictureFillTile.
	Mode string
	// Crop trims the source image before it is drawn, as fractions in [0,1)
	// taken off each edge. Nil leaves the image uncropped.
	Crop *PictureFillCropRequest
	// Row and Column select a table cell inside the shape. Both are nil for a
	// plain shape fill.
	Row    *int
	Column *int
}

// PictureFillCropRequest holds the four <a:srcRect> insets as fractions.
type PictureFillCropRequest struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

// SetShapePictureFill writes an <a:blipFill> onto a shape, or onto one cell of
// a table shape when Row and Column are set. The image is registered as a media
// part with a slide relationship first, which is why this cannot go through the
// generic fill model: common.PictureFill is only a read view.
func (e *PresentationEditor) SetShapePictureFill(
	slideIndex, shapeID int,
	request PictureFillRequest,
) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	if err := validatePictureFillRequest(request); err != nil {
		return err
	}

	partPath := e.slides[slideIndex].Part
	relID, err := addBgImageRelationship(e, partPath, SlideBackground{
		Type:      "image",
		ImagePath: request.ImagePath,
		ImageData: request.ImageData,
	})
	if err != nil {
		return err
	}
	fillXML := renderPictureFillXML(relID, request)

	if request.Row != nil || request.Column != nil {
		return e.setTableCellPictureFill(slideIndex, shapeID, request, fillXML)
	}
	return e.setShapeSpPrPictureFill(partPath, shapeID, fillXML)
}

func validatePictureFillRequest(request PictureFillRequest) error {
	hasPath := strings.TrimSpace(request.ImagePath) != ""
	hasData := strings.TrimSpace(request.ImageData) != ""
	if hasPath == hasData {
		return errors.New("picture fill requires exactly one of image_path or image_data")
	}
	switch request.Mode {
	case "", PictureFillStretch, PictureFillTile:
	default:
		return fmt.Errorf("unknown picture fill mode %q", request.Mode)
	}
	if (request.Row == nil) != (request.Column == nil) {
		return errors.New("picture fill on a table cell needs both row and column")
	}
	if request.Crop != nil {
		for name, value := range map[string]float64{
			"left":   request.Crop.Left,
			"top":    request.Crop.Top,
			"right":  request.Crop.Right,
			"bottom": request.Crop.Bottom,
		} {
			if value < 0 || value >= 1 {
				return fmt.Errorf("picture fill crop %s must be in [0,1), got %v", name, value)
			}
		}
	}
	return nil
}

// renderPictureFillXML builds the <a:blipFill> element. CT_BlipFillProperties
// orders blip, srcRect, then the tile/stretch choice.
func renderPictureFillXML(relID string, request PictureFillRequest) string {
	var b strings.Builder
	b.WriteString(`<a:blipFill rotWithShape="1"><a:blip r:embed="`)
	b.WriteString(common.XMLEscape(relID))
	b.WriteString(`"/>`)
	if crop := request.Crop; crop != nil {
		fmt.Fprintf(
			&b,
			`<a:srcRect l="%d" t="%d" r="%d" b="%d"/>`,
			int(crop.Left*pictureFillCropScale),
			int(crop.Top*pictureFillCropScale),
			int(crop.Right*pictureFillCropScale),
			int(crop.Bottom*pictureFillCropScale),
		)
	}
	if request.Mode == PictureFillTile {
		b.WriteString(`<a:tile tx="0" ty="0" sx="100000" sy="100000" flip="none" algn="tl"/>`)
	} else {
		b.WriteString(`<a:stretch><a:fillRect/></a:stretch>`)
	}
	b.WriteString(`</a:blipFill>`)
	return b.String()
}

// setShapeSpPrPictureFill swaps the shape's existing fill for the blipFill,
// leaving its line and effects alone.
func (e *PresentationEditor) setShapeSpPrPictureFill(partPath string, shapeID int, fillXML string) error {
	content, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}
	parsed, err := parseSlideShapes(content)
	if err != nil {
		return fmt.Errorf("parse shapes: %w", err)
	}

	found := false
	updated := replaceShapeNodes(content, parsed, func(_ int, shape *parsedShape) ([]byte, bool) {
		if shape.ID != shapeID {
			return nil, false
		}
		found = true
		return editorshape.ReplaceStyleInSpPr(content[shape.Start:shape.End], fillXML, true, false, false), true
	})
	if !found {
		return fmt.Errorf("shape id %d not found in part %s", shapeID, partPath)
	}
	e.parts.Set(partPath, updated)
	return nil
}

// setTableCellPictureFill writes the blipFill into one cell's <a:tcPr>.
func (e *PresentationEditor) setTableCellPictureFill(
	slideIndex, shapeID int,
	request PictureFillRequest,
	fillXML string,
) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	rowIdx, colIdx := *request.Row, *request.Column
	parsed, err := tablemod.ParseTable(frame)
	if err != nil {
		return err
	}
	rows, cols := tablemod.Dimensions(parsed)
	if rowIdx < 0 || rowIdx >= rows || colIdx < 0 || colIdx >= cols {
		return fmt.Errorf("table cell [%d,%d] out of range", rowIdx, colIdx)
	}

	newFrame, err := tablemod.MutateTableRows(frame, rowIdx, rowIdx,
		func(_ int, rowContent []byte) ([]byte, error) {
			return tablemod.MutateTableCells(rowContent, colIdx, colIdx,
				func(_ int, cellContent []byte) ([]byte, error) {
					return tablemod.SetCellFillXML(cellContent, fillXML)
				})
		})
	if err != nil {
		return err
	}

	out := make([]byte, 0, len(slideContent)-(frameEnd-frameStart)+len(newFrame))
	out = append(out, slideContent[:frameStart]...)
	out = append(out, newFrame...)
	out = append(out, slideContent[frameEnd:]...)
	e.parts.Set(partPath, out)
	return nil
}
