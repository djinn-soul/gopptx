package pptxxml

import (
	"strconv"
	"strings"
)

const (
	richInnerShadowGrowCap       = 160
	richPerspectiveShadowGrowCap = 200
)

func shadowDirEMU(angle float64) int {
	return int(angle * emusPerDegree)
}

func shadowAlphaValue(transparency float64) int {
	return alphaFromNormalizedTransparency(transparency)
}

func writeShadowBlurDistDirAttrs(b *strings.Builder, shadow RichShapeShadowSpec) {
	b.WriteString(`blurRad="`)
	b.WriteString(strconv.Itoa(shadow.BlurRadius))
	b.WriteString(`" dist="`)
	b.WriteString(strconv.Itoa(shadow.Distance))
	b.WriteString(`" dir="`)
	b.WriteString(strconv.Itoa(shadowDirEMU(shadow.Angle)))
	b.WriteString(`"`)
}

func writeShadowDistDirAttrs(b *strings.Builder, shadow RichShapeShadowSpec) {
	b.WriteString(`dist="`)
	b.WriteString(strconv.Itoa(shadow.Distance))
	b.WriteString(`" dir="`)
	b.WriteString(strconv.Itoa(shadowDirEMU(shadow.Angle)))
	b.WriteString(`"`)
}

func richInnerShadowXML(shadow RichShapeShadowSpec) string {
	var b strings.Builder
	b.Grow(richInnerShadowGrowCap)
	b.WriteString(`<a:innerShdw `)
	writeShadowBlurDistDirAttrs(&b, shadow)
	b.WriteString(`><a:srgbClr val="`)
	b.WriteString(Escape(shadow.Color))
	b.WriteString(`"><a:alpha val="`)
	b.WriteString(strconv.Itoa(shadowAlphaValue(shadow.Transparency)))
	b.WriteString(`"/></a:srgbClr></a:innerShdw>`)
	return b.String()
}

func richPerspectiveShadowXML(shadow RichShapeShadowSpec) string {
	var b strings.Builder
	b.Grow(richPerspectiveShadowGrowCap)
	b.WriteString(`<a:prstShdw prst="shdw1" `)
	writeShadowDistDirAttrs(&b, shadow)
	if shadow.SkewX != 0 || shadow.SkewY != 0 {
		b.WriteString(` sx="`)
		b.WriteString(strconv.Itoa(int(shadow.SkewX * shadowScaleBase)))
		b.WriteString(`" sy="`)
		b.WriteString(strconv.Itoa(int(shadow.SkewY * shadowScaleBase)))
		b.WriteString(`"`)
	}
	if shadow.ScaleX != 1.0 || shadow.ScaleY != 1.0 {
		b.WriteString(` kx="`)
		b.WriteString(strconv.Itoa(int(shadow.ScaleX * shadowScaleBase)))
		b.WriteString(`" ky="`)
		b.WriteString(strconv.Itoa(int(shadow.ScaleY * shadowScaleBase)))
		b.WriteString(`"`)
	}
	if shadow.Alignment != "" {
		b.WriteString(` algn="`)
		b.WriteString(Escape(shadow.Alignment))
		b.WriteString(`"`)
	}
	b.WriteString(`><a:srgbClr val="`)
	b.WriteString(Escape(shadow.Color))
	b.WriteString(`"><a:alpha val="`)
	b.WriteString(strconv.Itoa(shadowAlphaValue(shadow.Transparency)))
	b.WriteString(`"/></a:srgbClr></a:prstShdw>`)
	return b.String()
}
