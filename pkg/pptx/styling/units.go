package styling

import "math"

// Length represents a distance in English Metric Units (EMU).
type Length int64

const (
	emuPerInch float64 = 914400
	emuPerCM   float64 = 360000
	emuPerPT   float64 = 12700

	// MaxEMU is the maximum EMU value allowed by OOXML (int32 max).
	MaxEMU Length = 2147483647 // 0x7FFFFFFF
)

// Inches converts inches to Length (EMU) with overflow protection.
func Inches(value float64) Length {
	return clampToLength(value * emuPerInch)
}

// InchesToEMU is an alias for Inches.
func InchesToEMU(value float64) Length {
	return Inches(value)
}

// Centimeters converts centimeters to Length (EMU) with overflow protection.
func Centimeters(value float64) Length {
	return clampToLength(value * emuPerCM)
}

// CMToEMU is an alias for Centimeters.
func CMToEMU(value float64) Length {
	return Centimeters(value)
}

// Points converts points to Length (EMU) with overflow protection.
func Points(value float64) Length {
	return clampToLength(value * emuPerPT)
}

// PointsToEMU is an alias for Points.
func PointsToEMU(value float64) Length {
	return Points(value)
}

// Emu returns value as Length.
func Emu(value int64) Length {
	return Length(value)
}

// Inches returns the length in inches.
func (l Length) Inches() float64 {
	return float64(l) / emuPerInch
}

// Cm returns the length in centimeters.
func (l Length) Cm() float64 {
	return float64(l) / emuPerCM
}

// Pt returns the length in points.
func (l Length) Pt() float64 {
	return float64(l) / emuPerPT
}

// Emu returns the length in EMU units.
func (l Length) Emu() int64 {
	return int64(l)
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

func clampToLength(val float64) Length {
	if math.IsNaN(val) {
		return 0
	}
	if math.IsInf(val, 1) {
		return MaxEMU
	}
	if math.IsInf(val, -1) {
		return -MaxEMU
	}
	if val > float64(MaxEMU) {
		return MaxEMU
	}
	if val < float64(-MaxEMU) {
		return -MaxEMU
	}
	return Length(math.Round(val))
}
