package pptxxml

import (
	"math"
	"strconv"
	"strings"
)

const (
	shapeTextColorLight = "FFFFFF"
	shapeTextColorDark  = "000000"
)

func shapeTextRunPropertiesXML(shape ShapeSpec) string {
	return `<a:solidFill><a:srgbClr val="` + Escape(shapeTextColorForShape(shape)) + `"/></a:solidFill>`
}

func shapeTextColorForShape(shape ShapeSpec) string {
	r, g, b, ok := shapeBackgroundColor(shape)
	if !ok {
		return shapeTextColorDark
	}
	if contrastRatio(r, g, b, 255, 255, 255) >= contrastRatio(r, g, b, 0, 0, 0) {
		return shapeTextColorLight
	}
	return shapeTextColorDark
}

func shapeBackgroundColor(shape ShapeSpec) (int, int, int, bool) {
	if shape.Fill != nil {
		return colorWithTransparency(shape.Fill.Color, shape.Fill.TransparencyPct)
	}
	if shape.GradientFill != nil {
		return gradientAverageColor(*shape.GradientFill)
	}
	return 0, 0, 0, false
}

func gradientAverageColor(fill ShapeGradientFillSpec) (int, int, int, bool) {
	if len(fill.Stops) == 0 {
		return 0, 0, 0, false
	}

	var sumR float64
	var sumG float64
	var sumB float64
	var count float64
	for _, stop := range fill.Stops {
		r, g, b, ok := colorWithTransparency(stop.Color, stop.TransparencyPct)
		if !ok {
			continue
		}
		sumR += float64(r)
		sumG += float64(g)
		sumB += float64(b)
		count++
	}
	if count == 0 {
		return 0, 0, 0, false
	}
	return int(math.Round(sumR / count)), int(math.Round(sumG / count)), int(math.Round(sumB / count)), true
}

func colorWithTransparency(color string, transparencyPct *int) (int, int, int, bool) {
	r, g, b, ok := parseHexColor(color)
	if !ok {
		return 0, 0, 0, false
	}
	alpha := 1.0
	if transparencyPct != nil {
		percent := *transparencyPct
		if percent < 0 {
			percent = 0
		}
		if percent > 100 {
			percent = 100
		}
		alpha = float64(100-percent) / 100
	}
	return blendWithWhite(r, alpha), blendWithWhite(g, alpha), blendWithWhite(b, alpha), true
}

func blendWithWhite(channel int, alpha float64) int {
	value := alpha*float64(channel) + (1-alpha)*255
	return int(math.Round(value))
}

func parseHexColor(color string) (int, int, int, bool) {
	clean := strings.TrimPrefix(strings.TrimSpace(color), "#")
	if len(clean) != 6 {
		return 0, 0, 0, false
	}
	value, err := strconv.ParseUint(clean, 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	return int((value >> 16) & 0xFF), int((value >> 8) & 0xFF), int(value & 0xFF), true
}

func contrastRatio(r1 int, g1 int, b1 int, r2 int, g2 int, b2 int) float64 {
	l1 := relativeLuminance(r1, g1, b1)
	l2 := relativeLuminance(r2, g2, b2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

func relativeLuminance(r int, g int, b int) float64 {
	rLinear := srgbToLinear(float64(r) / 255)
	gLinear := srgbToLinear(float64(g) / 255)
	bLinear := srgbToLinear(float64(b) / 255)
	return 0.2126*rLinear + 0.7152*gLinear + 0.0722*bLinear
}

func srgbToLinear(value float64) float64 {
	if value <= 0.04045 {
		return value / 12.92
	}
	return math.Pow((value+0.055)/1.055, 2.4)
}
