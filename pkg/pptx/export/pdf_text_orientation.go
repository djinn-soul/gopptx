package export

import (
	"math"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func beginShapeTextOrientation(
	pdf *gopdf.GoPdf,
	frame *shapes.TextFrame,
	boxX, boxY, boxW, boxH, shapeX, shapeY, shapeW, shapeH float64,
) (float64, float64, float64, float64, func()) {
	if frame == nil {
		return boxX, boxY, boxW, boxH, func() {}
	}
	angle, rotate := shapeTextRotationAngle(frame)
	if !rotate {
		return boxX, boxY, boxW, boxH, func() {}
	}
	centerX := shapeX + shapeW/2
	centerY := shapeY + shapeH/2
	pdf.Rotate(angle, centerX, centerY)
	if isVerticalShapeText(frame.Orientation) {
		boxW, boxH = boxH, boxW
		boxX = centerX - boxW/2
		boxY = centerY - boxH/2
	}
	return boxX, boxY, boxW, boxH, pdf.RotateReset
}

func shapeTextRotationAngle(frame *shapes.TextFrame) (float64, bool) {
	if frame == nil {
		return 0, false
	}
	angle := 0.0
	if frame.RotationDeg != nil {
		angle += *frame.RotationDeg
	}
	switch strings.ToLower(strings.TrimSpace(frame.Orientation)) {
	case "vert", "wordartvert", "eavert", "mongolianvert":
		angle += 90
	case "vert270", "wordartvertrtl":
		angle -= 90
	}
	if math.Abs(angle) < 0.01 {
		return 0, false
	}
	return angle, true
}

func isVerticalShapeText(orientation string) bool {
	switch strings.ToLower(strings.TrimSpace(orientation)) {
	case "vert", "vert270", "wordartvert", "wordartvertrtl", "eavert", "mongolianvert":
		return true
	default:
		return false
	}
}
