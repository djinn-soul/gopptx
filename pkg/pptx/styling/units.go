package styling

import "math"

const (
	emuPerInch = 914400.0
	emuPerCM   = 360000.0
	emuPerPT   = 12700.0

	// MaxEMU is the maximum EMU value allowed by OOXML (int32 max).
	MaxEMU = 2147483647 // 0x7FFFFFFF
)

// Inches converts inches to EMU units with overflow protection.
func Inches(value float64) int64 {
	return clampToEMU(value * emuPerInch)
}

// InchesToEMU is an alias for Inches.
func InchesToEMU(value float64) int64 {
	return Inches(value)
}

// Centimeters converts centimeters to EMU units with overflow protection.
func Centimeters(value float64) int64 {
	return clampToEMU(value * emuPerCM)
}

// CMToEMU is an alias for Centimeters.
func CMToEMU(value float64) int64 {
	return Centimeters(value)
}

// Points converts points to EMU units with overflow protection.
func Points(value float64) int64 {
	return clampToEMU(value * emuPerPT)
}

// PointsToEMU is an alias for Points.
func PointsToEMU(value float64) int64 {
	return Points(value)
}

// FontSize converts points to OOXML size units (hundredths of a point).
func FontSize(pt float64) int {
	val := math.Round(pt * 100)
	if val > math.MaxInt32 {
		return math.MaxInt32
	}
	if val < math.MinInt32 {
		return math.MinInt32
	}
	return int(val)
}

func clampToEMU(val float64) int64 {
	if val > MaxEMU {
		return MaxEMU
	}
	if val < -MaxEMU {
		return -MaxEMU
	}
	return int64(math.Round(val))
}
