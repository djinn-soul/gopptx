package pptxxml

import (
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	shapeTextMinPt      = 10
	shapeTextMaxPt      = 36
	shapeTextDefaultPt  = 18
	shapeSizingBase     = 42
	shapeSizingSlope    = 0.50
	shapeSizingBoundMod = 0.28
	emuPerPoint         = 12700
)

func shapeTextSizeXML(shape ShapeSpec) string {
	return strconv.Itoa(autoFitShapeTextSizePt(shape) * ptFactor)
}

func autoFitShapeTextSizePt(shape ShapeSpec) int {
	chars := utf8.RuneCountInString(strings.TrimSpace(shape.Text))
	if chars <= 0 {
		return shapeTextDefaultPt
	}

	dimensionPts := float64(minInt64(shape.CX, shape.CY)) / float64(emuPerPoint)
	// More conservative sizing for shapes (Star, Heart, etc.)
	// Most shapes have internal margins or narrow areas.
	sizeByBounds := int(math.Round(dimensionPts * shapeSizingBoundMod))
	sizeByChars := int(math.Round(shapeSizingBase - shapeSizingSlope*float64(chars)))

	sizePt := min(sizeByBounds, sizeByChars)
	return clampInt(sizePt, shapeTextMinPt, shapeTextMaxPt)
}

func minInt64(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func clampInt(value int, minValue int, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}
