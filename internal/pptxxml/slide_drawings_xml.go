package pptxxml

import (
	"fmt"
	"strings"
)

const (
	shapeDescrAttrEmpty = ` descr=""`
	strokeDashSolid     = "solid"
	arrowTypeNone       = "none"
	emusPerDegree       = 60000
	transparencyBase    = 100000
	defaultMargin       = 457200
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
	cNvPrContent := customShapeNonVisualProperties(shape)
	xfrmXML := customShapeTransform(shape)
	fillXML := customShapeFill(shape)

	lineXML := ""
	if shape.Line != nil {
		lineXML = shapeLineXML(*shape.Line)
	}

	descrAttr := customShapeAltText(shape)

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="%d" name="Shape %d"%s%s
<p:cNvSpPr/>
<p:nvPr/>
</p:nvSpPr>
<p:spPr>
%s
<a:prstGeom prst="%s"><a:avLst/></a:prstGeom>%s%s
</p:spPr>`, shapeID, shapeID, descrAttr, cNvPrContent, xfrmXML, Escape(shape.Type), fillXML, lineXML))

	b.WriteString(customShapeTextBody(shape))
	b.WriteString(`
</p:sp>`)
	return b.String()
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
		rotationAttr = fmt.Sprintf(` rot="%d"`, *shape.RotationDeg*emusPerDegree)
	}
	return fmt.Sprintf(`
<a:xfrm%s>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</a:xfrm>`, rotationAttr, shape.X, shape.Y, shape.CX, shape.CY)
}

func customShapeFill(shape ShapeSpec) string {
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
		return fmt.Sprintf(` descr="%s"`, Escape(shape.AltText))
	}
	return ""
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
	bodyPrAttr := fmt.Sprintf(
		` wrap="square" rtlCol="0" anchor="ctr" lIns="%d" tIns="%d" rIns="%d" bIns="%d"`,
		defaultMargin, defaultMargin, defaultMargin, defaultMargin,
	)

	if shape.TextFrame != nil {
		bodyPrAttr = fmt.Sprintf(` wrap="%s" rtlCol="0" anchor="%s" lIns="%d" tIns="%d" rIns="%d" bIns="%d"`,
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

	hyperlinkXML := ""
	if shape.ClickAction != nil {
		hyperlinkXML = HyperlinkXML(*shape.ClickAction, "a:hlinkClick")
	} else if shape.Hyperlink != nil {
		hyperlinkXML = HyperlinkXML(*shape.Hyperlink, "a:hlinkClick")
	}

	return fmt.Sprintf(`
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
</p:txBody>`, bodyPrAttr, autoFitXML, shapeTextSizeXML(shape), shapeTextRunPropertiesXML(shape), hyperlinkXML, Escape(shape.Text))
}

func connectorXML(connector ConnectorSpec, shapeID int, startShapeID int, endShapeID int) string {
	xfrm := connectorTransform(connector)
	avLst := connectorAdjustments(connector)
	descrAttr := connectorAltText(connector)
	connections := connectorCxn(startShapeID, endShapeID, connector.StartSiteIndex, connector.EndSiteIndex)

	capAttr := ""
	if strings.TrimSpace(connector.Line.Cap) != "" {
		capAttr = ` cap="` + Escape(connector.Line.Cap) + `"`
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
<p:cxnSp>
<p:nvCxnSpPr>
<p:cNvPr id="%d" name="Connector %d"%s/>
<p:cNvCxnSpPr>%s</p:cNvCxnSpPr>
<p:nvPr/>
</p:nvCxnSpPr>
<p:spPr>
%s
<a:prstGeom prst="%s">%s</a:prstGeom>
<a:ln w="%d"%s>
<a:solidFill><a:srgbClr val="%s"/></a:solidFill>`,
		shapeID, shapeID, descrAttr, connections, xfrm, Escape(connector.Type), avLst,
		connector.Line.Width, capAttr, Escape(connector.Line.Color),
	))

	if strings.TrimSpace(connector.Line.Dash) != "" && connector.Line.Dash != strokeDashSolid {
		b.WriteString(`
<a:prstDash val="` + Escape(connector.Line.Dash) + `"/>`)
	}
	b.WriteString(connectorArrows(connector))
	b.WriteString(connectorLineJoin(connector.Line.Join))

	b.WriteString(`
</a:ln>
</p:spPr>`)
	b.WriteString(`
</p:cxnSp>`)
	return b.String()
}

func connectorTransform(connector ConnectorSpec) string {
	x, y, cx, cy := connectorBounds(connector)
	flipH := ""
	if connector.EndX < connector.StartX {
		flipH = ` flipH="1"`
	}
	flipV := ""
	if connector.EndY < connector.StartY {
		flipV = ` flipV="1"`
	}
	return fmt.Sprintf(`
<a:xfrm%s%s>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</a:xfrm>`, flipH, flipV, x, y, cx, cy)
}

func connectorCxn(startID, endID int, startIdx, endIdx *int) string {
	res := ""
	if startID > 0 && startIdx != nil {
		res += fmt.Sprintf(`
<a:stCxn id="%d" idx="%d"/>`, startID, *startIdx)
	}
	if endID > 0 && endIdx != nil {
		res += fmt.Sprintf(`
<a:endCxn id="%d" idx="%d"/>`, endID, *endIdx)
	}
	return res
}

func connectorAdjustments(connector ConnectorSpec) string {
	if len(connector.Adjustments) == 0 {
		return "<a:avLst/>"
	}
	var av strings.Builder
	av.WriteString("<a:avLst>")
	for _, adj := range connector.Adjustments {
		av.WriteString(`<a:gd name="` + Escape(adj.Name) + `" fmla="` + Escape(adj.Formula) + `"/>`)
	}
	av.WriteString("</a:avLst>")
	return av.String()
}

func connectorAltText(connector ConnectorSpec) string {
	if connector.IsDecorative {
		return shapeDescrAttrEmpty
	}
	if connector.AltText != "" {
		return fmt.Sprintf(` descr="%s"`, Escape(connector.AltText))
	}
	return ""
}

func connectorArrows(connector ConnectorSpec) string {
	res := ""
	if connector.StartArrow != arrowTypeNone {
		res += `
<a:headEnd type="` + Escape(connector.StartArrow) + `" w="` + Escape(connector.StartArrowWidth) + `" len="` + Escape(connector.StartArrowLen) + `"/>`
	}
	if connector.EndArrow != arrowTypeNone {
		res += `
<a:tailEnd type="` + Escape(connector.EndArrow) + `" w="` + Escape(connector.EndArrowWidth) + `" len="` + Escape(connector.EndArrowLen) + `"/>`
	}
	return res
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
	return fmt.Sprintf(`
<a:ln w="%d"%s>
<a:solidFill><a:srgbClr val="%s"/></a:solidFill>
%s
%s
</a:ln>`, line.Width, lineCapAttr, Escape(line.Color), dash, join)
}

func connectorLineJoin(join string) string {
	switch strings.TrimSpace(join) {
	case "bevel":
		return `
<a:bevel/>`
	case "miter":
		return `
<a:miter/>`
	case "round":
		return `
<a:round/>`
	default:
		return ""
	}
}

func alphaFromNormalizedTransparency(transparency float64) int {
	return int((1.0 - transparency) * transparencyBase)
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
