package pptxxml

import (
	"fmt"
	"strings"
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
	if stop.TransparencyPct != nil {
		alpha = fmt.Sprintf(`<a:alpha val="%d"/>`, alphaFromTransparencyPct(*stop.TransparencyPct))
	}
	return fmt.Sprintf(`
<a:gs pos="%d">
<a:srgbClr val="%s">%s</a:srgbClr>
</a:gs>`, stop.PositionPct*1000, Escape(stop.Color), alpha)
}

func shapeGradientPathXML(fill ShapeGradientFillSpec) string {
	switch normalizeShapeGradientType(fill.Type) {
	case "radial":
		return `
<a:path path="circle"/>`
	case "rectangular":
		return `
<a:path path="rect"/>`
	case "path":
		return `
<a:path path="shape"/>`
	default:
		angle := 0
		if fill.AngleDeg != nil {
			angle = *fill.AngleDeg
		}
		return fmt.Sprintf(`
<a:lin ang="%d" scaled="1"/>`, angle*60000)
	}
}

func normalizeShapeGradientType(fillType string) string {
	switch strings.ToLower(strings.TrimSpace(fillType)) {
	case "radial":
		return "radial"
	case "rectangular":
		return "rectangular"
	case "path":
		return "path"
	default:
		return "linear"
	}
}
