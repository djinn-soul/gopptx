package editor

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

const (
	minImageCropFraction = 0.0
	maxImageCropFraction = 1.0
	minImageRotationDeg  = -360.0
	maxImageRotationDeg  = 360.0
)

var (
	pictureBlipFillPattern = regexp.MustCompile(`(?s)<p:blipFill\b[^>]*>(.*?)</p:blipFill>`)
	pictureSrcRectPattern  = regexp.MustCompile(`(?s)<a:srcRect\b[^>]*/>`)
	pictureXfrmPattern     = regexp.MustCompile(`<a:xfrm\b([^>]*)>`)
)

func hasPictureUpdateFields(updates common.ShapeUpdate) bool {
	return updates.Crop != nil || updates.Rotation != nil || updates.FlipH != nil || updates.FlipV != nil
}

func validatePictureUpdateFields(updates common.ShapeUpdate) error {
	if updates.Crop != nil {
		crop := updates.Crop
		fields := []struct {
			name  string
			value float64
		}{
			{name: "crop.left", value: crop.Left},
			{name: "crop.right", value: crop.Right},
			{name: "crop.top", value: crop.Top},
			{name: "crop.bottom", value: crop.Bottom},
		}
		for _, field := range fields {
			if field.value < minImageCropFraction || field.value > maxImageCropFraction {
				return fmt.Errorf(
					"%s must be between %.1f and %.1f",
					field.name,
					minImageCropFraction,
					maxImageCropFraction,
				)
			}
		}
	}

	if updates.Rotation != nil {
		rotation := *updates.Rotation
		if rotation < minImageRotationDeg || rotation > maxImageRotationDeg {
			return fmt.Errorf(
				"rotation must be between %.0f and %.0f degrees",
				minImageRotationDeg,
				maxImageRotationDeg,
			)
		}
	}
	return nil
}

func applyPictureShapeUpdates(xmlData []byte, updates common.ShapeUpdate) ([]byte, error) {
	if err := validatePictureUpdateFields(updates); err != nil {
		return nil, err
	}
	updated := xmlData
	var err error

	if updates.Crop != nil {
		updated, err = replacePictureCrop(updated, updates.Crop)
		if err != nil {
			return nil, err
		}
	}
	if updates.Rotation != nil || updates.FlipH != nil || updates.FlipV != nil {
		updated, err = replacePictureTransformAttrs(updated, updates)
		if err != nil {
			return nil, err
		}
	}
	return updated, nil
}

func replacePictureCrop(xmlData []byte, crop *common.ImageCrop) ([]byte, error) {
	match := pictureBlipFillPattern.FindSubmatchIndex(xmlData)
	if match == nil {
		return nil, errors.New("picture shape missing blipFill")
	}

	inner := string(xmlData[match[2]:match[3]])
	inner = pictureSrcRectPattern.ReplaceAllString(inner, "")
	inner = strings.TrimSpace(inner)
	srcRect := buildImageCropXML(&common.ShapeUpdate{Crop: crop})
	inner = strings.TrimSpace(srcRect + inner)

	replacement := "<p:blipFill>" + inner + "</p:blipFill>"
	out := append([]byte{}, xmlData[:match[0]]...)
	out = append(out, replacement...)
	out = append(out, xmlData[match[1]:]...)
	return out, nil
}

func replacePictureTransformAttrs(xmlData []byte, updates common.ShapeUpdate) ([]byte, error) {
	match := pictureXfrmPattern.FindSubmatchIndex(xmlData)
	if match == nil {
		return nil, errors.New("picture shape missing xfrm")
	}

	attrs := string(xmlData[match[2]:match[3]])
	if updates.Rotation != nil {
		attrs = setOrRemoveXMLAttr(attrs, "rot", strconv.Itoa(int(*updates.Rotation*imageRotationScale)), true)
	}
	if updates.FlipH != nil {
		attrs = setOrRemoveXMLAttr(attrs, "flipH", "1", *updates.FlipH)
	}
	if updates.FlipV != nil {
		attrs = setOrRemoveXMLAttr(attrs, "flipV", "1", *updates.FlipV)
	}

	replacement := "<a:xfrm" + attrs + ">"
	out := append([]byte{}, xmlData[:match[0]]...)
	out = append(out, replacement...)
	out = append(out, xmlData[match[1]:]...)
	return out, nil
}

func setOrRemoveXMLAttr(attrs, name, value string, enabled bool) string {
	attrRe := regexp.MustCompile(`\s+` + regexp.QuoteMeta(name) + `="[^"]*"`)
	updated := attrRe.ReplaceAllString(attrs, "")
	if enabled {
		return updated + ` ` + name + `="` + value + `"`
	}
	return updated
}

// GetShapeTextState returns a read snapshot of text, runs, and paragraph/frame settings for one shape.
func (e *PresentationEditor) GetShapeTextState(slideIndex, shapeID int) (common.ShapeTextState, error) {
	shape, err := e.getShapeForTextOps(slideIndex, shapeID)
	if err != nil {
		return common.ShapeTextState{}, err
	}

	state := common.ShapeTextState{
		Text:      shape.Text,
		Runs:      editorshape.CopyTextRuns(shape.Runs),
		TextFrame: shape.TextFrame,
		Paragraph: shape.Paragraph,
	}
	return state, nil
}

// GetSlideTextStates returns text-related state for every parsed shape on a slide.
func (e *PresentationEditor) GetSlideTextStates(slideIndex int) ([]common.SlideShapeTextState, error) {
	shapes, err := e.getShapesForTextOps(slideIndex)
	if err != nil {
		return nil, err
	}

	states := make([]common.SlideShapeTextState, 0, len(shapes))
	for _, shape := range shapes {
		states = append(states, common.SlideShapeTextState{
			ShapeID:   shape.ID,
			Text:      shape.Text,
			Runs:      editorshape.CopyTextRuns(shape.Runs),
			TextFrame: shape.TextFrame,
			Paragraph: shape.Paragraph,
		})
	}
	return states, nil
}

// GetShapeRuns returns a copy of the shape's text runs.
func (e *PresentationEditor) GetShapeRuns(slideIndex, shapeID int) ([]common.TextRun, error) {
	state, err := e.GetShapeTextState(slideIndex, shapeID)
	if err != nil {
		return nil, err
	}
	return state.Runs, nil
}

// SetShapeRuns replaces all runs on a shape.
func (e *PresentationEditor) SetShapeRuns(slideIndex, shapeID int, runs []common.TextRun) error {
	runsCopy := editorshape.CopyTextRuns(runs)
	updates := common.ShapeUpdate{Runs: &runsCopy}
	return e.UpdateShape(slideIndex, shapeID, updates)
}

// UpdateRunText updates the text of one run by index.
func (e *PresentationEditor) UpdateRunText(slideIndex, shapeID, runIndex int, text string) error {
	runs, err := e.GetShapeRuns(slideIndex, shapeID)
	if err != nil {
		return err
	}
	updatedRuns, err := editorshape.UpdateRunText(runs, runIndex, text)
	if err != nil {
		return err
	}
	return e.SetShapeRuns(slideIndex, shapeID, updatedRuns)
}

// AppendShapeRun appends one run to the shape's existing runs.
func (e *PresentationEditor) AppendShapeRun(slideIndex, shapeID int, run common.TextRun) error {
	runs, err := e.GetShapeRuns(slideIndex, shapeID)
	if err != nil {
		return err
	}
	runs = editorshape.AppendRun(runs, run)
	return e.SetShapeRuns(slideIndex, shapeID, runs)
}

// RemoveShapeRun removes one run by index from a shape's run list.
func (e *PresentationEditor) RemoveShapeRun(slideIndex, shapeID, runIndex int) error {
	runs, err := e.GetShapeRuns(slideIndex, shapeID)
	if err != nil {
		return err
	}
	if runIndex < 0 || runIndex >= len(runs) {
		return fmt.Errorf("run index %d out of range", runIndex)
	}
	updated := append([]common.TextRun{}, runs[:runIndex]...)
	updated = append(updated, runs[runIndex+1:]...)
	if len(updated) == 0 {
		empty := ""
		updates := common.ShapeUpdate{Text: &empty, Runs: &updated}
		return e.UpdateShape(slideIndex, shapeID, updates)
	}
	return e.SetShapeRuns(slideIndex, shapeID, updated)
}

// RemoveShapeParagraph removes one paragraph by index.
// The current text model normalizes one paragraph per shape, so only index 0 is valid.
func (e *PresentationEditor) RemoveShapeParagraph(slideIndex, shapeID, paragraphIndex int) error {
	if paragraphIndex != 0 {
		return fmt.Errorf("paragraph index %d out of range", paragraphIndex)
	}
	empty := ""
	emptyRuns := []common.TextRun{}
	updates := common.ShapeUpdate{Text: &empty, Runs: &emptyRuns}
	return e.UpdateShape(slideIndex, shapeID, updates)
}

func (e *PresentationEditor) getShapeForTextOps(slideIndex, shapeID int) (parsedShape, error) {
	shapes, err := e.getShapesForTextOps(slideIndex)
	if err != nil {
		return parsedShape{}, err
	}

	for _, shape := range shapes {
		if shape.ID == shapeID || (shape.PhType == placeholderTypeTitle && shapeID == 0) {
			return shape, nil
		}
	}

	return parsedShape{}, fmt.Errorf("shape id %d not found on slide %d", shapeID, slideIndex)
}

func (e *PresentationEditor) getShapesForTextOps(slideIndex int) ([]parsedShape, error) {
	if e == nil {
		return nil, errors.New("editor cannot be nil")
	}
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapes, err := parseSlideShapes(content)
	if err != nil {
		return nil, fmt.Errorf("parse shapes: %w", err)
	}
	return shapes, nil
}
