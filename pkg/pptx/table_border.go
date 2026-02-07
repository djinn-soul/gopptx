package pptx

import "math"

const tableBorderPtToEMU = 12700.0

func tableBorderWidthEMU(widthPt float64) int64 {
	if widthPt <= 0 {
		return 0
	}
	width := int64(math.Round(widthPt * tableBorderPtToEMU))
	if width < 1 {
		return 1
	}
	return width
}
