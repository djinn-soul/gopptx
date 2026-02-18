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
	Name         string
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

	name := shape.Name
	if name == "" {
		name = fmt.Sprintf("Shape %d", shapeID)
	}

	var b strings.Builder
	b.Grow(2048)
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
	b.WriteString(`"><a:avLst/></a:prstGeom>`)
	b.WriteString(fillXML)
	b.WriteString(lineXML)
	b.WriteString(`
</p:spPr>`)

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
		rotationAttr = ` rot="` + strconv.Itoa(*shape.RotationDeg*emusPerDegree) + `"`
	}
	return `
<a:xfrm` + rotationAttr + `>
<a:off x="` + strconv.FormatInt(shape.X, 10) + `" y="` + strconv.FormatInt(shape.Y, 10) + `"/>
<a:ext cx="` + strconv.FormatInt(shape.CX, 10) + `" cy="` + strconv.FormatInt(shape.CY, 10) + `"/>
</a:xfrm>`
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
	bodyPrChildren := autoFitXML
	bodyPrAttr := ` wrap="square" rtlCol="0" anchor="ctr" lIns="` + strconv.Itoa(defaultMargin) + `" tIns="` + strconv.Itoa(defaultMargin) + `" rIns="` + strconv.Itoa(defaultMargin) + `" bIns="` + strconv.Itoa(defaultMargin) + `"`

	if shape.TextFrame != nil {
		bodyPrAttr = ` wrap="` + Escape(shape.TextFrame.Wrap) + `" rtlCol="0" anchor="` + Escape(shape.TextFrame.Anchor) + `" lIns="` + strconv.FormatInt(shape.TextFrame.MarginLeft, 10) + `" tIns="` + strconv.FormatInt(shape.TextFrame.MarginTop, 10) + `" rIns="` + strconv.FormatInt(shape.TextFrame.MarginRight, 10) + `" bIns="` + strconv.FormatInt(shape.TextFrame.MarginBottom, 10) + `"`
		switch shape.TextFrame.AutoFit {
		case "spAutoFit":
			autoFitXML = `<a:spAutoFit/>`
		case "normAutoFit":
			autoFitXML = `<a:normAutofit/>`
		default:
			autoFitXML = ""
		}
	}
	if shape.TextFrame != nil && shape.TextFrame.AutoFit == "normAutoFit" {
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
	b.WriteString(`
<p:cxnSp>
<p:nvCxnSpPr>
<p:cNvPr id="`)
	b.WriteString(strconv.Itoa(shapeID))
	b.WriteString(`" name="Connector `)
	b.WriteString(strconv.Itoa(shapeID))
	b.WriteString(`"`)
	b.WriteString(descrAttr)
	b.WriteString(`/>
<p:cNvCxnSpPr>`)
	b.WriteString(connections)
	b.WriteString(`</p:cNvCxnSpPr>
<p:nvPr/>
</p:nvCxnSpPr>
<p:spPr>
`)
	b.WriteString(xfrm)
	b.WriteString(`
<a:prstGeom prst="`)
	b.WriteString(Escape(connector.Type))
	b.WriteString(`">`)
	b.WriteString(avLst)
	b.WriteString(`</a:prstGeom>
<a:ln w="`)
	b.WriteString(strconv.FormatInt(connector.Line.Width, 10))
	b.WriteString(`"`)
	b.WriteString(capAttr)
	b.WriteString(`>
<a:solidFill><a:srgbClr val="`)
	b.WriteString(Escape(connector.Line.Color))
	b.WriteString(`"/></a:solidFill>`)

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

func connectorLabelShape(connector ConnectorSpec, shapeID int) string {
	label := strings.TrimSpace(connector.Label)
	if label == "" {
		return ""
	}
	x := (connector.StartX + connector.EndX) / 2
	y := (connector.StartY + connector.EndY) / 2
	const (
		labelWidth  = 914400
		labelHeight = 228600
	)
	return `
<p:sp>
<p:nvSpPr>
<p:cNvPr id="` + strconv.Itoa(shapeID) + `" name="Connector Label ` + strconv.Itoa(shapeID) + `"/>
<p:cNvSpPr txBox="1"/>
<p:nvPr/>
</p:nvSpPr>
<p:spPr>
<a:xfrm>
<a:off x="` + strconv.FormatInt(x-labelWidth/2, 10) + `" y="` + strconv.FormatInt(y-labelHeight/2, 10) + `"/>
<a:ext cx="` + strconv.Itoa(labelWidth) + `" cy="` + strconv.Itoa(labelHeight) + `"/>
</a:xfrm>
<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
<a:noFill/>
<a:ln><a:noFill/></a:ln>
</p:spPr>
<p:txBody>
<a:bodyPr wrap="square" rtlCol="0" anchor="ctr"/>
<a:lstStyle/>
<a:p>
<a:pPr algn="ctr"/>
<a:r>
<a:rPr lang="en-US" sz="1000" b="0" i="0" u="none" dirty="0"/>
<a:t>` + Escape(label) + `</a:t>
</a:r>
</a:p>
</p:txBody>
</p:sp>`
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
	return `
<a:xfrm` + flipH + flipV + `>
<a:off x="` + strconv.FormatInt(x, 10) + `" y="` + strconv.FormatInt(y, 10) + `"/>
<a:ext cx="` + strconv.FormatInt(cx, 10) + `" cy="` + strconv.FormatInt(cy, 10) + `"/>
</a:xfrm>`
}

func connectorCxn(startID, endID int, startIdx, endIdx *int) string {
	res := ""
	if startID > 0 && startIdx != nil {
		res += `
<a:stCxn id="` + strconv.Itoa(startID) + `" idx="` + strconv.Itoa(*startIdx) + `"/>`
	}
	if endID > 0 && endIdx != nil {
		res += `
<a:endCxn id="` + strconv.Itoa(endID) + `" idx="` + strconv.Itoa(*endIdx) + `"/>`
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
		escaped := Escape(connector.AltText)
		return ` descr="` + escaped + `" title="` + escaped + `"`
	}
	return shapeDescrAttrEmpty
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
