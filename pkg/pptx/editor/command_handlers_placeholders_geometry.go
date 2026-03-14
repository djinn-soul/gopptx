package editor

import "regexp"

var (
	prstGeomPattern = regexp.MustCompile(`(?s)<a:prstGeom\b.*?</a:prstGeom>`)
	custGeomPattern = regexp.MustCompile(`(?s)<a:custGeom\b.*?</a:custGeom>`)
)

func extractPlaceholderGeometryXML(shapeXML []byte) string {
	if match := custGeomPattern.Find(shapeXML); len(match) > 0 {
		return string(match)
	}
	if match := prstGeomPattern.Find(shapeXML); len(match) > 0 {
		return string(match)
	}
	return ""
}
