package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	shapeDescrAttrEmpty = ` descr=""`
	strokeDashSolid     = "solid"
	arrowTypeNone       = "none"
	emusPerDegree       = 60000
	shadowScaleBase     = 100000
	transparencyBase    = 100000
	defaultMargin       = 457200
	customShapeGrowCap  = 2048
	midpointDivisor     = 2
	normAutoFitToken    = "normAutoFit"
)

func customShapeXML(shape ShapeSpec, shapeID int) string {
	cNvPrContent := customShapeNonVisualProperties(shape)
	xfrmXML := customShapeTransform(shape)
	fillXML := customShapeFill(shape)

	lineXML := ""
	if shape.RichLine != nil {
		lineXML = richShapeLineXML(*shape.RichLine)
	} else if shape.Line != nil {
		lineXML = shapeLineXML(*shape.Line)
	}

	descrAttr := customShapeAltText(shape)

	name := shape.Name
	if name == "" {
		name = fmt.Sprintf("Shape %d", shapeID)
	}

	var b strings.Builder
	b.Grow(customShapeGrowCap)
	b.WriteString(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="`)
	b.WriteString(strconv.Itoa(shapeID))
	b.WriteString(`" name="`)
	b.WriteString(Escape(name))
	b.WriteString(`"`)
	b.WriteString(descrAttr)
	b.WriteString(cNvPrContent)
	b.WriteString(`
<p:cNvSpPr/>
<p:nvPr/>
</p:nvSpPr>
<p:spPr>
`)
	b.WriteString(xfrmXML)
	b.WriteString(`
<a:prstGeom prst="`)
	b.WriteString(Escape(shape.Type))
	b.WriteString(`">`)
	b.WriteString(shapeAdjustments(shape))
	b.WriteString(`</a:prstGeom>`)
	b.WriteString(fillXML)
	b.WriteString(lineXML)
	b.WriteString(shapeEffectsXML(shape.Effects, shape.RichShadow))
	b.WriteString(`
</p:spPr>`)

	b.WriteString(customShapeTextBody(shape))
	b.WriteString(`
</p:sp>`)
	return b.String()
}

func shapeEffectsXML(effects *ShapeEffectsSpec, richShadow *RichShapeShadowSpec) string {
	// Check if we have rich shadow
	hasRichShadow := richShadow != nil && richShadow.Type != ""

	// Check if we have legacy effects
	hasLegacyEffects := effects != nil && (effects.Shadow || effects.Glow || effects.SoftEdges || effects.Reflection)

	if !hasRichShadow && !hasLegacyEffects {
		return ""
	}

	var b strings.Builder
	b.WriteString("<a:effectLst>")

	// Rich shadow takes precedence over legacy shadow
	if hasRichShadow {
		b.WriteString(richShapeShadowXML(*richShadow))
	} else if effects != nil && effects.Shadow {
		b.WriteString(`<a:outerShdw blurRad="40000" dist="20000" dir="5400000" rotWithShape="0">`)
		b.WriteString(`<a:srgbClr val="000000"><a:alpha val="40000"/></a:srgbClr>`)
		b.WriteString(`</a:outerShdw>`)
	}

	if effects != nil {
		if effects.Glow {
			b.WriteString(`<a:glow rad="6350">`)
			b.WriteString(`<a:srgbClr val="4472C4"><a:alpha val="35000"/></a:srgbClr>`)
			b.WriteString(`</a:glow>`)
		}
		if effects.SoftEdges {
			b.WriteString(`<a:softEdge rad="38100"/>`)
		}
		if effects.Reflection {
			b.WriteString(`<a:ref blurRad="6350" stA="50000" endA="300" endPos="35000" dist="0"`)
			b.WriteString(` dir="5400000" sy="-100000" algn="bl" rotWithShape="0"/>`)
		}
	}

	b.WriteString("</a:effectLst>")
	return b.String()
}

func shapeAdjustments(shape ShapeSpec) string {
	if len(shape.Adjustments) == 0 {
		return "<a:avLst/>"
	}
	var av strings.Builder
	av.WriteString("<a:avLst>")
	for _, adj := range shape.Adjustments {
		av.WriteString(`<a:gd name="` + Escape(adj.Name) + `" fmla="` + Escape(adj.Formula) + `"/>`)
	}
	av.WriteString("</a:avLst>")
	return av.String()
}

func customShapeNonVisualProperties(shape ShapeSpec) string {
	hyperlinkXML := ""
	if shape.ClickAction != nil {
		hyperlinkXML = HyperlinkXML(*shape.ClickAction, "a:hlinkClick")
	} else if shape.Hyperlink != nil {
		hyperlinkXML = HyperlinkXML(*shape.Hyperlink, "a:hlinkClick")
	}

	if shape.HoverAction != nil {
		hyperlinkXML += HyperlinkXML(*shape.HoverAction, "a:hlinkHover")
	}

	if hyperlinkXML != "" {
		return ">" + hyperlinkXML + "</p:cNvPr>"
	}
	return "/>"
}

func customShapeTransform(shape ShapeSpec) string {
	rotationAttr := ""
	if shape.RotationDeg != nil {
		rotationAttr = ` rot="` + strconv.Itoa(*shape.RotationDeg*emusPerDegree) + `"`
	}
	return `
<a:xfrm` + rotationAttr + `>
<a:off x="` + strconv.FormatInt(shape.X, 10) + `" y="` + strconv.FormatInt(shape.Y, 10) + `"/>
<a:ext cx="` + strconv.FormatInt(shape.CX, 10) + `" cy="` + strconv.FormatInt(shape.CY, 10) + `"/>
</a:xfrm>`
}

func customShapeFill(shape ShapeSpec) string {
	// Rich fill takes precedence over legacy fill
	if shape.RichFill != nil {
		return richShapeFillXML(*shape.RichFill)
	}
	switch {
	case shape.GradientFill != nil:
		return shapeGradientFillXML(*shape.GradientFill)
	case shape.Fill != nil:
		return shapeSolidFillXML(*shape.Fill)
	default:
		return `
<a:noFill/>`
	}
}

func customShapeAltText(shape ShapeSpec) string {
	if shape.IsDecorative {
		return shapeDescrAttrEmpty
	}
	if shape.AltText != "" {
		escaped := Escape(shape.AltText)
		return ` descr="` + escaped + `" title="` + escaped + `"`
	}
	return shapeDescrAttrEmpty
}

func customShapeTextBody(shape ShapeSpec) string {
	if strings.TrimSpace(shape.Text) == "" {
		return `
<p:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p/>
</p:txBody>`
	}

	autoFitXML := `<a:spAutoFit/>`
	var bodyPrChildren string
	bodyPrAttr := ` wrap="square" rtlCol="0" anchor="ctr" lIns="` + strconv.Itoa(
		defaultMargin,
	) + `" tIns="` + strconv.Itoa(
		defaultMargin,
	) + `" rIns="` + strconv.Itoa(
		defaultMargin,
	) + `" bIns="` + strconv.Itoa(
		defaultMargin,
	) + `"`

	if shape.TextFrame != nil {
		bodyPrAttr = ` wrap="` + Escape(
			shape.TextFrame.Wrap,
		) + `" rtlCol="0" anchor="` + Escape(
			shape.TextFrame.Anchor,
		) + `" lIns="` + strconv.FormatInt(
			shape.TextFrame.MarginLeft,
			10,
		) + `" tIns="` + strconv.FormatInt(
			shape.TextFrame.MarginTop,
			10,
		) + `" rIns="` + strconv.FormatInt(
			shape.TextFrame.MarginRight,
			10,
		) + `" bIns="` + strconv.FormatInt(
			shape.TextFrame.MarginBottom,
			10,
		) + `"`
		if shape.TextFrame.Rotation != nil {
			bodyPrAttr += ` rot="` + strconv.FormatInt(*shape.TextFrame.Rotation, 10) + `"`
		}
		switch shape.TextFrame.AutoFit {
		case "spAutoFit":
			autoFitXML = `<a:spAutoFit/>`
		case normAutoFitToken:
			autoFitXML = `<a:normAutofit/>`
		case "none":
			autoFitXML = `<a:noAutofit/>`
		default:
			autoFitXML = ""
		}
	}
	if shape.TextFrame != nil && shape.TextFrame.AutoFit == normAutoFitToken {
		bodyPrChildren = `<a:prstTxWarp prst="textNoShape"><a:avLst/></a:prstTxWarp>` + "\n" + autoFitXML
	} else {
		bodyPrChildren = autoFitXML
	}

	hyperlinkXML := ""
	if shape.ClickAction != nil {
		hyperlinkXML = HyperlinkXML(*shape.ClickAction, "a:hlinkClick")
	} else if shape.Hyperlink != nil {
		hyperlinkXML = HyperlinkXML(*shape.Hyperlink, "a:hlinkClick")
	}

	return `
<p:txBody>
<a:bodyPr` + bodyPrAttr + `>
` + bodyPrChildren + `
</a:bodyPr>
<a:lstStyle/>
<a:p>
<a:pPr algn="l"/>
<a:r>
<a:rPr lang="en-US" sz="` + shapeTextSizeXML(shape) + `" b="0" i="0" u="none" dirty="0">` + shapeTextRunPropertiesXML(shape) + hyperlinkXML + `</a:rPr>
<a:t>` + Escape(shape.Text) + `</a:t>
</a:r>
</a:p>
</p:txBody>`
}
