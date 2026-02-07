package pptxxml

import (
	"fmt"
	"strings"
)

// ShapeFillSpec describes solid fill properties for a custom shape.
type ShapeFillSpec struct {
	Color           string
	TransparencyPct *int
}

// ShapeLineSpec describes line properties for a custom shape or connector.
type ShapeLineSpec struct {
	Color string
	Width int64
	Dash  string
}

// ShapeSpec describes one custom shape rendered as p:sp.
type ShapeSpec struct {
	Type        string
	X           int64
	Y           int64
	CX          int64
	CY          int64
	Fill        *ShapeFillSpec
	Line        *ShapeLineSpec
	Text        string
	RotationDeg *int
}

// ConnectorSpec describes one custom connector rendered as p:cxnSp.
type ConnectorSpec struct {
	Type            string
	StartX          int64
	StartY          int64
	EndX            int64
	EndY            int64
	Line            ShapeLineSpec
	StartArrow      string
	EndArrow        string
	ArrowSize       string
	StartShapeIndex int
	StartSiteIndex  *int
	EndShapeIndex   int
	EndSiteIndex    *int
	Label           string
}

func customShapeXML(shape ShapeSpec, shapeID int) string {
	rotationAttr := ""
	if shape.RotationDeg != nil {
		rotationAttr = fmt.Sprintf(` rot="%d"`, *shape.RotationDeg*60000)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="%d" name="Shape %d"/>
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
		rotationAttr,
		shape.X,
		shape.Y,
		shape.CX,
		shape.CY,
		Escape(shape.Type),
	))

	if shape.Fill != nil {
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
		b.WriteString(`
<p:txBody>
<a:bodyPr wrap="square" rtlCol="0" anchor="ctr"/>
<a:lstStyle/>
<a:p>
<a:pPr algn="ctr"/>
<a:r>
<a:rPr lang="en-US" sz="1800" b="0" i="0" u="none" dirty="0"/>
<a:t>` + Escape(shape.Text) + `</a:t>
</a:r>
</a:p>
</p:txBody>`)
	}
	b.WriteString(`
</p:sp>`)
	return b.String()
}

func shapeSolidFillXML(fill ShapeFillSpec) string {
	alpha := ""
	if fill.TransparencyPct != nil {
		alpha = fmt.Sprintf(`<a:alpha val="%d"/>`, alphaFromTransparencyPct(*fill.TransparencyPct))
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
	return `
<a:ln w="` + fmt.Sprintf("%d", line.Width) + `">
<a:solidFill><a:srgbClr val="` + Escape(line.Color) + `"/></a:solidFill>
` + dash + `
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

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
<p:cxnSp>
<p:nvCxnSpPr>
<p:cNvPr id="%d" name="Connector %d"/>
<p:cNvCxnSpPr>`, shapeID, shapeID))
	if startShapeID > 0 && connector.StartSiteIndex != nil {
		b.WriteString(fmt.Sprintf(`
<a:stCxn id="%d" idx="%d"/>`, startShapeID, *connector.StartSiteIndex))
	}
	if endShapeID > 0 && connector.EndSiteIndex != nil {
		b.WriteString(fmt.Sprintf(`
<a:endCxn id="%d" idx="%d"/>`, endShapeID, *connector.EndSiteIndex))
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
<a:prstGeom prst="%s"><a:avLst/></a:prstGeom>
<a:ln w="%d">
<a:solidFill><a:srgbClr val="%s"/></a:solidFill>`,
		flipH,
		flipV,
		x,
		y,
		cx,
		cy,
		Escape(connector.Type),
		connector.Line.Width,
		Escape(connector.Line.Color),
	))
	if strings.TrimSpace(connector.Line.Dash) != "" && connector.Line.Dash != "solid" {
		b.WriteString(`
<a:prstDash val="` + Escape(connector.Line.Dash) + `"/>`)
	}
	if connector.StartArrow != "none" {
		b.WriteString(`
<a:headEnd type="` + Escape(connector.StartArrow) + `" w="` + Escape(connector.ArrowSize) + `" len="` + Escape(connector.ArrowSize) + `"/>`)
	}
	if connector.EndArrow != "none" {
		b.WriteString(`
<a:tailEnd type="` + Escape(connector.EndArrow) + `" w="` + Escape(connector.ArrowSize) + `" len="` + Escape(connector.ArrowSize) + `"/>`)
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

func alphaFromTransparencyPct(transparencyPct int) int {
	return (100 - transparencyPct) * 1000
}

func connectorBounds(connector ConnectorSpec) (x int64, y int64, cx int64, cy int64) {
	x = connector.StartX
	if connector.EndX < x {
		x = connector.EndX
	}
	y = connector.StartY
	if connector.EndY < y {
		y = connector.EndY
	}
	cx = connector.EndX - connector.StartX
	if cx < 0 {
		cx = -cx
	}
	cy = connector.EndY - connector.StartY
	if cy < 0 {
		cy = -cy
	}
	return x, y, cx, cy
}
