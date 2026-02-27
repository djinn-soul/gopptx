package tplx

import "strconv"

// strconvAppendFloat formats a float64 to string using strconv.
func strconvAppendFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
