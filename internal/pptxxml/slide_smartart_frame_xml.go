package pptxxml

import "fmt"

// smartArtFrameShape renders the p:graphicFrame XML that embeds a SmartArt
// diagram on a slide. It uses dgm:relIds to reference the 4 DiagramML
// relationship parts (data, layout, quickStyle, colors).
func smartArtFrameShape(frame *SmartArtFrame, shapeID int) string {
	return fmt.Sprintf(`
<p:graphicFrame>
<p:nvGraphicFramePr>
<p:cNvPr id="%d" name="Diagram %d"/>
<p:cNvGraphicFramePr/>
<p:nvPr/>
</p:nvGraphicFramePr>
<p:xfrm>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</p:xfrm>
<a:graphic>
<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/diagram">
<dgm:relIds xmlns:dgm="http://schemas.openxmlformats.org/drawingml/2006/diagram" `+
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" `+
		`r:dm="%s" r:lo="%s" r:qs="%s" r:cs="%s"/>
</a:graphicData>
</a:graphic>
</p:graphicFrame>`,
		shapeID,
		shapeID,
		frame.X,
		frame.Y,
		frame.CX,
		frame.CY,
		Escape(frame.DataRelID),
		Escape(frame.LayoutRelID),
		Escape(frame.StyleRelID),
		Escape(frame.ColorRelID),
	)
}
