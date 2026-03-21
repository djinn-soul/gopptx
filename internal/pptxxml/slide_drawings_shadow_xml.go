package pptxxml

import "fmt"

func shadowDirEMU(angle float64) int {
	return int(angle * emusPerDegree)
}

func shadowAlphaValue(transparency float64) int {
	return alphaFromNormalizedTransparency(transparency)
}

func shadowBlurDistDirAttrs(shadow RichShapeShadowSpec) string {
	return fmt.Sprintf(
		`blurRad="%d" dist="%d" dir="%d"`,
		shadow.BlurRadius,
		shadow.Distance,
		shadowDirEMU(shadow.Angle),
	)
}

func shadowDistDirAttrs(shadow RichShapeShadowSpec) string {
	return fmt.Sprintf(
		`dist="%d" dir="%d"`,
		shadow.Distance,
		shadowDirEMU(shadow.Angle),
	)
}

func richInnerShadowXML(shadow RichShapeShadowSpec) string {
	attrs := shadowBlurDistDirAttrs(shadow)
	alphaVal := shadowAlphaValue(shadow.Transparency)

	return fmt.Sprintf(`<a:innerShdw %s><a:srgbClr val="%s"><a:alpha val="%d"/></a:srgbClr></a:innerShdw>`,
		attrs, Escape(shadow.Color), alphaVal)
}

func richPerspectiveShadowXML(shadow RichShapeShadowSpec) string {
	attrs := shadowDistDirAttrs(shadow)

	if shadow.SkewX != 0 || shadow.SkewY != 0 {
		attrs += fmt.Sprintf(` sx="%d" sy="%d"`,
			int(shadow.SkewX*shadowScaleBase), int(shadow.SkewY*shadowScaleBase))
	}

	if shadow.ScaleX != 1.0 || shadow.ScaleY != 1.0 {
		attrs += fmt.Sprintf(` kx="%d" ky="%d"`,
			int(shadow.ScaleX*shadowScaleBase), int(shadow.ScaleY*shadowScaleBase))
	}

	if shadow.Alignment != "" {
		attrs += fmt.Sprintf(` algn="%s"`, Escape(shadow.Alignment))
	}

	alphaVal := shadowAlphaValue(shadow.Transparency)

	return fmt.Sprintf(
		`<a:prstShdw prst="shdw1" %s><a:srgbClr val="%s"><a:alpha val="%d"/></a:srgbClr></a:prstShdw>`,
		attrs, Escape(shadow.Color), alphaVal,
	)
}
