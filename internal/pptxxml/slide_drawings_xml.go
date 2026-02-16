package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

// ShapeFillSpec describes solid fill properties for a custom shape.
type ShapeFillSpec struct {
	Color        string
	Transparency *float64
}

// ShapeGradientStopSpec describes one gradient stop for a custom shape.
type ShapeGradientStopSpec struct {
	PositionPct  int
	Color        string
	Transparency *float64
}

// ShapeGradientFillSpec describes gradient fill properties for a custom shape.
type ShapeGradientFillSpec struct {
	Type     string
	Stops    []ShapeGradientStopSpec
	AngleDeg *int
}

// ShapeLineSpec describes line properties for a custom shape or connector.
type ShapeLineSpec struct {
	Color string
	Width int64
	Dash  string
	Cap   string
	Join  string
}

// ShapeSpec describes one custom shape rendered as p:sp.
type ShapeSpec struct {
	Type         string
	X            int64
	Y            int64
	CX           int64
	CY           int64
	Fill         *ShapeFillSpec
	GradientFill *ShapeGradientFillSpec
	Line         *ShapeLineSpec
	Text         string
	RotationDeg  *int
	Hyperlink    *HyperlinkSpec // Legacy: mapped to ClickAction
	ClickAction  *HyperlinkSpec
	HoverAction  *HyperlinkSpec
	AltText      string
	IsDecorative bool
	TextFrame    *TextFrameSpec
}

// TextFrameSpec describes the text layout within a shape.
type TextFrameSpec struct {
	MarginLeft   int64
	MarginRight  int64
	MarginTop    int64
	MarginBottom int64
	Anchor       string
	Wrap         string
	AutoFit      string
}

// ConnectorSpec describes one custom connector rendered as p:cxnSp.
type ConnectorSpec struct {
	Type            string
	StartX, StartY  int64
	EndX, EndY      int64
	Line            ShapeLineSpec
	StartArrow      string
	StartArrowWidth string
	StartArrowLen   string
	EndArrow        string
	EndArrowWidth   string
	EndArrowLen     string
	StartShapeIndex int
	StartSiteIndex  *int
	EndShapeIndex   int
	EndSiteIndex    *int
	Label           string
	AltText         string
	IsDecorative    bool
	Adjustments     []ConnectorAdjustmentSpec
}

// ConnectorAdjustmentSpec describes one connector adjustment entry (a:gd).
type ConnectorAdjustmentSpec struct {
	Name    string
	Formula string
}

func customShapeXML(shape ShapeSpec, shapeID int) string {
	rotationAttr := ""
	if shape.RotationDeg != nil {
		rotationAttr = fmt.Sprintf(` rot="%d"`, *shape.RotationDeg*60000)
	}

	hyperlinkXML := ""
	if shape.ClickAction != nil {
		hyperlinkXML = HyperlinkXML(*shape.ClickAction, "a:hlinkClick")
	} else if shape.Hyperlink != nil {
		hyperlinkXML = HyperlinkXML(*shape.Hyperlink, "a:hlinkClick")
	}

	if shape.HoverAction != nil {
		hyperlinkXML += HyperlinkXML(*shape.HoverAction, "a:hlinkHover")
	}

	// If hyperlink is present, p:cNvPr must have a child element, so we use a different format
	cNvPrContent := "/>"
	if hyperlinkXML != "" {
		cNvPrContent = ">" + hyperlinkXML + "</p:cNvPr>"
	}

	descrAttr := ""
	if shape.IsDecorative {
		descrAttr = ` descr=""`
	} else if shape.AltText != "" {
		descrAttr = fmt.Sprintf(` descr="%s"`, Escape(shape.AltText))
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="%d" name="Shape %d"%s%s
<p:cNvSpPr/>
<p:nvPr/>
</p:nvSpPr>
<p:spPr>
<a:xfrm%s>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</a:xfrm>
<a:prstGeom prst="%s"><a:avLst/></a:prstGeom>`,
		shapeID,
		shapeID,
		descrAttr, // We append descr to name/id attrs
		cNvPrContent,
		rotationAttr,
		shape.X,
		shape.Y,
		shape.CX,
		shape.CY,
		Escape(shape.Type),
	))

	if shape.GradientFill != nil {
		b.WriteString(shapeGradientFillXML(*shape.GradientFill))
	} else if shape.Fill != nil {
		b.WriteString(shapeSolidFillXML(*shape.Fill))
	} else {
		b.WriteString(`
<a:noFill/>`)
	}
	if shape.Line != nil {
		b.WriteString(shapeLineXML(*shape.Line))
	}
	b.WriteString(`
</p:spPr>`)

	if strings.TrimSpace(shape.Text) == "" {
		b.WriteString(`
<p:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p/>
</p:txBody>`)
	} else {
		autoFitXML := `<a:spAutoFit/>`
		bodyPrAttr := ` wrap="square" rtlCol="0" anchor="ctr" marL="45720" marT="45720" marR="45720" marB="45720"`

		if shape.TextFrame != nil {
			bodyPrAttr = fmt.Sprintf(` wrap="%s" rtlCol="0" anchor="%s" marL="%d" marT="%d" marR="%d" marB="%d"`,
				Escape(shape.TextFrame.Wrap),
				Escape(shape.TextFrame.Anchor),
				shape.TextFrame.MarginLeft,
				shape.TextFrame.MarginTop,
				shape.TextFrame.MarginRight,
				shape.TextFrame.MarginBottom,
			)
			switch shape.TextFrame.AutoFit {
			case "spAutoFit":
				autoFitXML = `<a:spAutoFit/>`
			case "normAutoFit":
				autoFitXML = `<a:normAutoFit/>`
			default:
				autoFitXML = ""
			}
		}

		b.WriteString(fmt.Sprintf(`
<p:txBody>
<a:bodyPr%s>
%s
</a:bodyPr>
<a:lstStyle/>
<a:p>
<a:pPr algn="l"/>
<a:r>
<a:rPr lang="en-US" sz="%s" b="0" i="0" u="none" dirty="0">%s%s</a:rPr>
<a:t>%s</a:t>
</a:r>
</a:p>
</p:txBody>`, bodyPrAttr, autoFitXML, shapeTextSizeXML(shape), shapeTextRunPropertiesXML(shape), hyperlinkXML, Escape(shape.Text)))
	}
	b.WriteString(`
</p:sp>`)
	return b.String()
}

func shapeSolidFillXML(fill ShapeFillSpec) string {
	alpha := ""
	if fill.Transparency != nil {
		alpha = fmt.Sprintf(`<a:alpha val="%d"/>`, alphaFromNormalizedTransparency(*fill.Transparency))
	}
	return `
<a:solidFill>
<a:srgbClr val="` + Escape(fill.Color) + `">` + alpha + `</a:srgbClr>
</a:solidFill>`
}

func shapeLineXML(line ShapeLineSpec) string {
	dash := ""
	if strings.TrimSpace(line.Dash) != "" && line.Dash != "solid" {
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

func connectorXML(connector ConnectorSpec, shapeID int, startShapeID int, endShapeID int) string {
	x, y, cx, cy := connectorBounds(connector)
	flipH := ""
	if connector.EndX < connector.StartX {
		flipH = ` flipH="1"`
	}
	flipV := ""
	if connector.EndY < connector.StartY {
		flipV = ` flipV="1"`
	}

	descrAttr := ""
	if connector.IsDecorative {
		descrAttr = ` descr=""`
	} else if connector.AltText != "" {
		descrAttr = fmt.Sprintf(` descr="%s"`, Escape(connector.AltText))
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
<p:cxnSp>
<p:nvCxnSpPr>
<p:cNvPr id="%d" name="Connector %d"%s/>
<p:cNvCxnSpPr>`, shapeID, shapeID, descrAttr))
	if startShapeID > 0 && connector.StartSiteIndex != nil {
		b.WriteString(fmt.Sprintf(`
<a:stCxn id="%d" idx="%d"/>`, startShapeID, *connector.StartSiteIndex))
	}
	if endShapeID > 0 && connector.EndSiteIndex != nil {
		b.WriteString(fmt.Sprintf(`
<a:endCxn id="%d" idx="%d"/>`, endShapeID, *connector.EndSiteIndex))
	}
	avLst := "<a:avLst/>"
	if len(connector.Adjustments) > 0 {
		var av strings.Builder
		av.WriteString("<a:avLst>")
		for _, adj := range connector.Adjustments {
			av.WriteString(`<a:gd name="` + Escape(adj.Name) + `" fmla="` + Escape(adj.Formula) + `"/>`)
		}
		av.WriteString("</a:avLst>")
		avLst = av.String()
	}

	capAttr := ""
	if strings.TrimSpace(connector.Line.Cap) != "" {
		capAttr = ` cap="` + Escape(connector.Line.Cap) + `"`
	}

	b.WriteString(fmt.Sprintf(`
</p:cNvCxnSpPr>
<p:nvPr/>
</p:nvCxnSpPr>
<p:spPr>
<a:xfrm%s%s>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</a:xfrm>
<a:prstGeom prst="%s">%s</a:prstGeom>
<a:ln w="%d"%s>
<a:solidFill><a:srgbClr val="%s"/></a:solidFill>`,
		flipH,
		flipV,
		x,
		y,
		cx,
		cy,
		Escape(connector.Type),
		avLst,
		connector.Line.Width,
		capAttr,
		Escape(connector.Line.Color),
	))
	if strings.TrimSpace(connector.Line.Dash) != "" && connector.Line.Dash != "solid" {
		b.WriteString(`
<a:prstDash val="` + Escape(connector.Line.Dash) + `"/>`)
	}
	if connector.StartArrow != "none" {
		b.WriteString(`
<a:headEnd type="` + Escape(connector.StartArrow) + `" w="` + Escape(connector.StartArrowWidth) + `" len="` + Escape(connector.StartArrowLen) + `"/>`)
	}
	if connector.EndArrow != "none" {
		b.WriteString(`
<a:tailEnd type="` + Escape(connector.EndArrow) + `" w="` + Escape(connector.EndArrowWidth) + `" len="` + Escape(connector.EndArrowLen) + `"/>`)
	}
	switch strings.TrimSpace(connector.Line.Join) {
	case "bevel":
		b.WriteString(`
<a:bevel/>`)
	case "miter":
		b.WriteString(`
<a:miter/>`)
	case "round":
		b.WriteString(`
<a:round/>`)
	}
	b.WriteString(`
</a:ln>
</p:spPr>`)
	if strings.TrimSpace(connector.Label) != "" {
		b.WriteString(`
<p:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p>
<a:r>
<a:rPr lang="en-US" sz="1000" b="0" i="0" u="none" dirty="0"/>
<a:t>` + Escape(connector.Label) + `</a:t>
</a:r>
</a:p>
</p:txBody>`)
	}
	b.WriteString(`
</p:cxnSp>`)
	return b.String()
}

func alphaFromNormalizedTransparency(transparency float64) int {
	return int((1.0 - transparency) * 100000)
}

func connectorBounds(connector ConnectorSpec) (int64, int64, int64, int64) {
	x := min(connector.EndX, connector.StartX)
	y := min(connector.EndY, connector.StartY)
	cx := connector.EndX - connector.StartX
	if cx < 0 {
		cx = -cx
	}
	cy := connector.EndY - connector.StartY
	if cy < 0 {
		cy = -cy
	}
	return x, y, cx, cy
}
