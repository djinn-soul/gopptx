package editor

import (
	"regexp"
	"strconv"
)

var (
	prstGeomPattern = regexp.MustCompile(`(?s)<a:prstGeom\b.*?</a:prstGeom>`)
	custGeomPattern = regexp.MustCompile(`(?s)<a:custGeom\b.*?</a:custGeom>`)
	offPattern      = regexp.MustCompile(`<a:off\s+x="(-?\d+)"\s+y="(-?\d+)"\s*/>`)
	extPattern      = regexp.MustCompile(`<a:ext\s+cx="(\d+)"\s+cy="(\d+)"\s*/>`)
)

// placeholderBounds is a shape's own a:off/a:ext pair, in EMU.
type placeholderBounds struct {
	X  int64
	Y  int64
	CX int64
	CY int64
}

// extractPlaceholderBounds reads the a:off/a:ext pair from a shape's own
// a:xfrm. A placeholder that inherits its geometry from the layout has neither,
// in which case ok is false and the caller has no box to fit an image into.
func extractPlaceholderBounds(shapeXML []byte) (placeholderBounds, bool) {
	off := offPattern.FindSubmatch(shapeXML)
	ext := extPattern.FindSubmatch(shapeXML)
	if off == nil || ext == nil {
		return placeholderBounds{}, false
	}
	raws := [][]byte{off[1], off[2], ext[1], ext[2]}
	values := make([]int64, 0, len(raws))
	for _, raw := range raws {
		value, err := strconv.ParseInt(string(raw), 10, 64)
		if err != nil {
			return placeholderBounds{}, false
		}
		values = append(values, value)
	}
	return placeholderBounds{X: values[0], Y: values[1], CX: values[2], CY: values[3]}, true
}

func extractPlaceholderGeometryXML(shapeXML []byte) string {
	if match := custGeomPattern.Find(shapeXML); len(match) > 0 {
		return string(match)
	}
	if match := prstGeomPattern.Find(shapeXML); len(match) > 0 {
		return string(match)
	}
	return ""
}
