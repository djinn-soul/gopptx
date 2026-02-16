package pptxxml

import (
	"fmt"
	"strings"
)

const (
	shapeGradientTypeRadial      = "radial"
	shapeGradientTypeRectangular = "rectangular"
	shapeGradientTypePath        = "path"
	gradientPosFactor            = 1000
)

func shapeGradientFillXML(fill ShapeGradientFillSpec) string {
	var b strings.Builder
	b.WriteString(`
<a:gradFill rotWithShape="1">
<a:gsLst>`)
	for _, stop := range fill.Stops {
		b.WriteString(shapeGradientStopXML(stop))
	}
	b.WriteString(`
</a:gsLst>`)
	b.WriteString(shapeGradientPathXML(fill))
	b.WriteString(`
</a:gradFill>`)
	return b.String()
}

func shapeGradientStopXML(stop ShapeGradientStopSpec) string {
	alpha := ""
	if stop.Transparency != nil {
		alpha = fmt.Sprintf(`<a:alpha val="%d"/>`, alphaFromNormalizedTransparency(*stop.Transparency))
	}
	return fmt.Sprintf(`
<a:gs pos="%d">
<a:srgbClr val="%s">%s</a:srgbClr>
</a:gs>`, stop.PositionPct*gradientPosFactor, Escape(stop.Color), alpha)
}

func shapeGradientPathXML(fill ShapeGradientFillSpec) string {
	switch normalizeShapeGradientType(fill.Type) {
	case shapeGradientTypeRadial:
		return `
<a:path path="circle"/>`
	case shapeGradientTypeRectangular:
		return `
<a:path path="rect"/>`
	case shapeGradientTypePath:
		return `
<a:path path="shape"/>`
	default:
		angle := 0
		if fill.AngleDeg != nil {
			angle = *fill.AngleDeg
		}
		return fmt.Sprintf(`
<a:lin ang="%d" scaled="1"/>`, angle*emusPerDegree)
	}
}

func normalizeShapeGradientType(fillType string) string {
	switch strings.ToLower(strings.TrimSpace(fillType)) {
	case shapeGradientTypeRadial:
		return shapeGradientTypeRadial
	case shapeGradientTypeRectangular:
		return shapeGradientTypeRectangular
	case shapeGradientTypePath:
		return shapeGradientTypePath
	default:
		return "linear"
	}
}
