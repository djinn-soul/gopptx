package tables

import "math"

const tableBorderPtToEMU = 12700.0

// TableBorderWidthEMU converts points to EMU.
func TableBorderWidthEMU(widthPt float64) int64 {
	if widthPt <= 0 {
		return 0
	}
	width := int64(math.Round(widthPt * tableBorderPtToEMU))
	if width < 1 {
		return 1
	}
	return width
}

// TableMarginEMU converts points to EMU pointer.
func TableMarginEMU(marginPt *float64) *int64 {
	if marginPt == nil {
		return nil
	}
	value := TableBorderWidthEMU(*marginPt)
	if *marginPt == 0 {
		value = 0
	}
	return &value
}
