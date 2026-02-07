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

func tableMarginEMU(marginPt *float64) *int64 {
	if marginPt == nil {
		return nil
	}
	value := tableBorderWidthEMU(*marginPt)
	if *marginPt == 0 {
		value = 0
	}
	return &value
}
