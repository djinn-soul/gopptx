package mermaid

import (
	"strconv"
	"strings"
)

func parseQuadrant(code string) *QuadrantDiagram {
	lines := ParseLines(code)
	quadrant := &QuadrantDiagram{}

	for _, line := range lines {
		consumeQuadrantLine(quadrant, strings.TrimSpace(line))
	}

	return quadrant
}

func consumeQuadrantLine(quadrant *QuadrantDiagram, trimmed string) {
	if trimmed == "" {
		return
	}

	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "quadrantchart") {
		return
	}
	if value, ok := parseQuadrantTitle(trimmed, lower); ok {
		quadrant.Title = value
		return
	}
	if value, ok := parseQuadrantAxis(trimmed, lower, "x-axis"); ok {
		quadrant.XAxis = value
		return
	}
	if value, ok := parseQuadrantAxis(trimmed, lower, "y-axis"); ok {
		quadrant.YAxis = value
		return
	}
	if idx, text, ok := parseQuadrantLabel(trimmed, lower); ok {
		quadrant.Quadrants[idx] = text
		return
	}
	if point, ok := parseQuadrantPoint(trimmed); ok {
		quadrant.Points = append(quadrant.Points, point)
	}
}

func parseQuadrantTitle(trimmed string, lower string) (string, bool) {
	if !strings.HasPrefix(lower, "title") {
		return "", false
	}
	return strings.TrimSpace(trimmed[5:]), true
}

func parseQuadrantAxis(trimmed string, lower string, key string) (string, bool) {
	if !strings.HasPrefix(lower, key) {
		return "", false
	}
	return strings.TrimSpace(trimmed[len(key):]), true
}

func parseQuadrantLabel(trimmed string, lower string) (int, string, bool) {
	if !strings.HasPrefix(lower, "quadrant-") || len(trimmed) < 10 {
		return 0, "", false
	}
	idx, err := strconv.Atoi(trimmed[9:10])
	if err != nil || idx < 1 || idx > 4 {
		return 0, "", false
	}
	return idx - 1, strings.TrimSpace(trimmed[10:]), true
}

func parseQuadrantPoint(trimmed string) (QuadrantPoint, bool) {
	if !strings.Contains(trimmed, ":") || !strings.Contains(trimmed, "[") || !strings.Contains(trimmed, "]") {
		return QuadrantPoint{}, false
	}
	labelPart, coordPart, ok := strings.Cut(trimmed, ":")
	if !ok {
		return QuadrantPoint{}, false
	}
	x, y, ok := parseQuadrantCoords(coordPart)
	if !ok {
		return QuadrantPoint{}, false
	}
	return QuadrantPoint{
		Label: strings.TrimSpace(labelPart),
		X:     x,
		Y:     y,
	}, true
}

func parseQuadrantCoords(coordPart string) (float64, float64, bool) {
	coords := strings.Trim(strings.TrimSpace(coordPart), "[]")
	xPart, yPart, ok := strings.Cut(coords, ",")
	if !ok {
		return 0, 0, false
	}
	x, errX := strconv.ParseFloat(strings.TrimSpace(xPart), 64)
	y, errY := strconv.ParseFloat(strings.TrimSpace(yPart), 64)
	if errX != nil || errY != nil {
		return 0, 0, false
	}
	return x, y, true
}
