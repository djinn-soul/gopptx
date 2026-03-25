package pptxxml

import (
	"strconv"
	"strings"
)

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
</p:spPr>
</p:cxnSp>`)
	return b.String()
}

func connectorLabelShape(connector ConnectorSpec, shapeID int) string {
	label := strings.TrimSpace(connector.Label)
	if label == "" {
		return ""
	}
	x := (connector.StartX + connector.EndX) / midpointDivisor
	y := (connector.StartY + connector.EndY) / midpointDivisor
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
