package pptxxml

import (
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
	var b strings.Builder
	b.Grow(96)
	b.WriteString(`<a:solidFill><a:srgbClr val="`)
	b.WriteString(Escape(fill.Color))
	b.WriteString(`">`)
	if fill.Transparency > 0 {
		alphaVal := int((1.0 - fill.Transparency) * transparencyBase)
		b.WriteString(`<a:alpha val="`)
		b.WriteString(strconv.Itoa(alphaVal))
		b.WriteString(`"/>`)
	}
	b.WriteString(`</a:srgbClr></a:solidFill>`)
	return b.String()
}

func richPatternFillXML(fill PatternFillSpec) string {
	var b strings.Builder
	b.Grow(160)
	b.WriteString(`<a:pattFill prst="`)
	b.WriteString(Escape(fill.Pattern))
	b.WriteString(`"><a:fgClr><a:srgbClr val="`)
	b.WriteString(Escape(fill.FgColor))
	b.WriteString(`"/></a:fgClr><a:bgClr><a:srgbClr val="`)
	b.WriteString(Escape(fill.BgColor))
	b.WriteString(`"/></a:bgClr></a:pattFill>`)
	return b.String()
}

func richShapeLineXML(line RichShapeLineSpec) string {
	var b strings.Builder
	b.Grow(160)
	b.WriteString(`<a:ln w="`)
	b.WriteString(strconv.FormatInt(line.Width, 10))
	b.WriteString(`"`)
	if line.CapStyle != "" {
		b.WriteString(` cap="`)
		b.WriteString(string(line.CapStyle))
		b.WriteString(`"`)
	}
	b.WriteString(`><a:solidFill><a:srgbClr val="`)
	b.WriteString(Escape(line.Color))
	b.WriteString(`">`)
	if line.Transparency > 0 {
		alphaVal := int((1.0 - line.Transparency) * transparencyBase)
		b.WriteString(`<a:alpha val="`)
		b.WriteString(strconv.Itoa(alphaVal))
		b.WriteString(`"/>`)
	}
	b.WriteString(`</a:srgbClr></a:solidFill>`)
	if line.DashStyle != "" && line.DashStyle != LineDashStyleSolid {
		b.WriteString(`<a:prstDash val="`)
		b.WriteString(string(line.DashStyle))
		b.WriteString(`"/>`)
	}
	switch line.JoinStyle {
	case LineJoinStyleBevel:
		b.WriteString(`<a:bevel/>`)
	case LineJoinStyleMiter:
		b.WriteString(`<a:miter/>`)
	case LineJoinStyleRound:
		b.WriteString(`<a:round/>`)
	}
	b.WriteString(`</a:ln>`)
	return b.String()
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
	var b strings.Builder
	b.Grow(180)
	b.WriteString(`<a:outerShdw `)
	writeShadowBlurDistDirAttrs(&b, shadow)
	if shadow.Alignment != "" {
		b.WriteString(` algn="`)
		b.WriteString(Escape(shadow.Alignment))
		b.WriteString(`"`)
	}
	if !shadow.RotateWithShape {
		b.WriteString(` rotWithShape="0"`)
	}
	b.WriteString(`><a:srgbClr val="`)
	b.WriteString(Escape(shadow.Color))
	b.WriteString(`"><a:alpha val="`)
	b.WriteString(strconv.Itoa(shadowAlphaValue(shadow.Transparency)))
	b.WriteString(`"/></a:srgbClr></a:outerShdw>`)
	return b.String()
}
