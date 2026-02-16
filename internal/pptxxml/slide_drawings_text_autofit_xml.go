package pptxxml

import (
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	shapeTextMinPt = 10
	shapeTextMaxPt = 36
)

func shapeTextSizeXML(shape ShapeSpec) string {
	return strconv.Itoa(autoFitShapeTextSizePt(shape) * 100)
}

func autoFitShapeTextSizePt(shape ShapeSpec) int {
	chars := utf8.RuneCountInString(strings.TrimSpace(shape.Text))
	if chars <= 0 {
		return 18
	}

	dimensionPts := float64(minInt64(shape.CX, shape.CY)) / 12700
	// More conservative sizing for shapes (Star, Heart, etc.)
	// Most shapes have internal margins or narrow areas.
	sizeByBounds := int(math.Round(dimensionPts * 0.28))
	sizeByChars := int(math.Round(42 - 0.50*float64(chars)))

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
