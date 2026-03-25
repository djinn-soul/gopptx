package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

func shapeSolidFillXML(fill ShapeFillSpec) string {
	alphaXML := ""
	if fill.Transparency != nil {
		alphaXML = `<a:alpha val="` + strconv.Itoa(alphaFromNormalizedTransparency(*fill.Transparency)) + `"/>`
	}
	return `
<a:solidFill>
<a:srgbClr val="` + Escape(fill.Color) + `">` + alphaXML + `</a:srgbClr>
</a:solidFill>`
}

func shapeLineXML(line ShapeLineSpec) string {
	dash := ""
	if strings.TrimSpace(line.Dash) != "" && line.Dash != strokeDashSolid {
		dash = `<a:prstDash val="` + Escape(line.Dash) + `"/>`
	}
	lineCapAttr := ""
	if strings.TrimSpace(line.Cap) != "" {
		lineCapAttr = ` cap="` + Escape(line.Cap) + `"`
	}
	join := ""
	switch strings.TrimSpace(line.Join) {
	case "bevel":
		join = `<a:bevel/>`
	case "miter":
		join = `<a:miter/>`
	case "round":
		join = `<a:round/>`
	}
	return `
<a:ln w="` + strconv.FormatInt(line.Width, 10) + `"` + lineCapAttr + `>
<a:solidFill><a:srgbClr val="` + Escape(line.Color) + `"/></a:solidFill>
` + dash + `
` + join + `
</a:ln>`
}

func alphaFromNormalizedTransparency(transparency float64) int {
	return int((1.0 - transparency) * transparencyBase)
}

func richShapeFillXML(fill RichShapeFillSpec) string {
	switch fill.Type {
	case FillTypeSolid:
		if fill.Solid != nil {
			return richSolidFillXML(*fill.Solid)
		}
	case FillTypeGradient:
		if fill.Gradient != nil {
			return shapeGradientFillXML(*fill.Gradient)
		}
	case FillTypePattern:
		if fill.Pattern != nil {
			return richPatternFillXML(*fill.Pattern)
		}
	case FillTypeNoFill:
		return `<a:noFill/>`
	}
	return `<a:noFill/>`
}

func richSolidFillXML(fill SolidFillSpec) string {
	alphaXML := ""
	if fill.Transparency > 0 {
		alphaVal := int((1.0 - fill.Transparency) * transparencyBase)
		alphaXML = fmt.Sprintf(`<a:alpha val="%d"/>`, alphaVal)
	}
	return fmt.Sprintf(`<a:solidFill><a:srgbClr val="%s">%s</a:srgbClr></a:solidFill>`,
		Escape(fill.Color), alphaXML)
}

func richPatternFillXML(fill PatternFillSpec) string {
	return fmt.Sprintf(
		`<a:pattFill prst="%s">`+
			`<a:fgClr><a:srgbClr val="%s"/></a:fgClr>`+
			`<a:bgClr><a:srgbClr val="%s"/></a:bgClr>`+
			`</a:pattFill>`,
		Escape(fill.Pattern),
		Escape(fill.FgColor),
		Escape(fill.BgColor),
	)
}

func richShapeLineXML(line RichShapeLineSpec) string {
	attrs := fmt.Sprintf(`w="%d"`, line.Width)
	if line.CapStyle != "" {
		attrs += fmt.Sprintf(` cap="%s"`, string(line.CapStyle))
	}

	dashXML := ""
	if line.DashStyle != "" && line.DashStyle != LineDashStyleSolid {
		dashXML = fmt.Sprintf(`<a:prstDash val="%s"/>`, string(line.DashStyle))
	}

	joinXML := ""
	switch line.JoinStyle {
	case LineJoinStyleBevel:
		joinXML = `<a:bevel/>`
	case LineJoinStyleMiter:
		joinXML = `<a:miter/>`
	case LineJoinStyleRound:
		joinXML = `<a:round/>`
	}

	alphaXML := ""
	if line.Transparency > 0 {
		alphaVal := int((1.0 - line.Transparency) * transparencyBase)
		alphaXML = fmt.Sprintf(`<a:alpha val="%d"/>`, alphaVal)
	}
	return fmt.Sprintf(`<a:ln %s><a:solidFill><a:srgbClr val="%s">%s</a:srgbClr></a:solidFill>%s%s</a:ln>`,
		attrs, Escape(line.Color), alphaXML, dashXML, joinXML)
}

func richShapeShadowXML(shadow RichShapeShadowSpec) string {
	if shadow.Type == "" {
		return ""
	}
	switch shadow.Type {
	case ShadowTypeOuter:
		return richOuterShadowXML(shadow)
	case ShadowTypeInner:
		return richInnerShadowXML(shadow)
	case ShadowTypePerspective:
		return richPerspectiveShadowXML(shadow)
	}
	return ""
}

func richOuterShadowXML(shadow RichShapeShadowSpec) string {
	attrs := shadowBlurDistDirAttrs(shadow)
	if shadow.Alignment != "" {
		attrs += fmt.Sprintf(` algn="%s"`, Escape(shadow.Alignment))
	}
	if !shadow.RotateWithShape {
		attrs += ` rotWithShape="0"`
	}
	alphaVal := shadowAlphaValue(shadow.Transparency)
	return fmt.Sprintf(`<a:outerShdw %s><a:srgbClr val="%s"><a:alpha val="%d"/></a:srgbClr></a:outerShdw>`,
		attrs, Escape(shadow.Color), alphaVal)
}
