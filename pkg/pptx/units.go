package pptx

import "math"

const (
	emuPerInch = 914400.0
	emuPerCM   = 360000.0
	emuPerPT   = 12700.0
)

// Inches converts inches to EMU units.
func Inches(value float64) int64 {
	return int64(math.Round(value * emuPerInch))
}

// Centimeters converts centimeters to EMU units.
func Centimeters(value float64) int64 {
	return int64(math.Round(value * emuPerCM))
}

// Points converts points to EMU units.
func Points(value float64) int64 {
	return int64(math.Round(value * emuPerPT))
}
