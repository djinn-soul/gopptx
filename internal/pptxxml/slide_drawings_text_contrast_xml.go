package pptxxml

import (
	"encoding/hex"
	"math"
	"strings"
)

const (
	shapeTextColorLight = "FFFFFF"
	shapeTextColorDark  = "000000"
	// Color algorithms use standard coefficients.
	sRGBGamma      = 2.4
	sRGBOffset     = 0.055
	sRGBDivisor    = 1.055
	sRGBThreshold  = 0.04045
	sRFBSlope      = 12.92
	lumaR          = 0.2126
	lumaG          = 0.7152
	lumaB          = 0.0722
	contrastOffset = 0.05
	maxByteValue   = 255
	hexColorLen    = 6
)

func shapeTextRunPropertiesXML(shape ShapeSpec) string {
	return `<a:solidFill><a:srgbClr val="` + Escape(shapeTextColorForShape(shape)) + `"/></a:solidFill>`
}

func shapeTextColorForShape(shape ShapeSpec) string {
	r, g, b, ok := shapeBackgroundColor(shape)
	if !ok {
		return shapeTextColorDark
	}
	if contrastRatio(r, g, b, maxByteValue, maxByteValue, maxByteValue) >= contrastRatio(r, g, b, 0, 0, 0) {
		return shapeTextColorLight
	}
	return shapeTextColorDark
}

func shapeBackgroundColor(shape ShapeSpec) (int, int, int, bool) {
	if shape.Fill != nil {
		return colorWithTransparency(shape.Fill.Color, shape.Fill.Transparency)
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
		r, g, b, ok := colorWithTransparency(stop.Color, stop.Transparency)
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

func colorWithTransparency(color string, transparency *float64) (int, int, int, bool) {
	r, g, b, ok := parseHexColor(color)
	if !ok {
		return 0, 0, 0, false
	}
	alpha := 1.0
	if transparency != nil {
		percent := *transparency
		if percent < 0 {
			percent = 0
		}
		if percent > 1 {
			percent = 1
		}
		alpha = 1.0 - percent
	}
	return blendWithWhite(r, alpha), blendWithWhite(g, alpha), blendWithWhite(b, alpha), true
}

func blendWithWhite(channel int, alpha float64) int {
	value := alpha*float64(channel) + (1-alpha)*maxByteValue
	return int(math.Round(value))
}

func parseHexColor(color string) (int, int, int, bool) {
	clean := strings.TrimPrefix(strings.TrimSpace(color), "#")
	if len(clean) != hexColorLen {
		return 0, 0, 0, false
	}
	rgb, err := hex.DecodeString(clean)
	if err != nil {
		return 0, 0, 0, false
	}
	return int(rgb[0]), int(rgb[1]), int(rgb[2]), true
}

func contrastRatio(r1 int, g1 int, b1 int, r2 int, g2 int, b2 int) float64 {
	l1 := relativeLuminance(r1, g1, b1)
	l2 := relativeLuminance(r2, g2, b2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + contrastOffset) / (l2 + contrastOffset)
}

func relativeLuminance(r int, g int, b int) float64 {
	rLinear := srgbToLinear(float64(r) / maxByteValue)
	gLinear := srgbToLinear(float64(g) / maxByteValue)
	bLinear := srgbToLinear(float64(b) / maxByteValue)
	return lumaR*rLinear + lumaG*gLinear + lumaB*bLinear
}

func srgbToLinear(value float64) float64 {
	if value <= sRGBThreshold {
		return value / sRFBSlope
	}
	return math.Pow((value+sRGBOffset)/sRGBDivisor, sRGBGamma)
}
