package pptxxml

import "fmt"

// ChartFrame describes chart placement in slide XML.
type ChartFrame struct {
	RelID string
	X     int64
	Y     int64
	CX    int64
	CY    int64

	// Accessibility
	AltText      string
	IsDecorative bool
}

func chartFrameShape(chart *ChartFrame, shapeID int) string {
	return fmt.Sprintf(`
<p:graphicFrame>
<p:nvGraphicFramePr>
<p:cNvPr id="%d" name="Chart %d"%s/>
<p:cNvGraphicFramePr/>
<p:nvPr/>
</p:nvGraphicFramePr>
<p:xfrm>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</p:xfrm>
<a:graphic>
<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/chart">
<c:chart xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" r:id="%s"/>
</a:graphicData>
</a:graphic>
</p:graphicFrame>`,
		shapeID,
		shapeID,
		makeCNvPrAttrs(chart.AltText, chart.IsDecorative),
		chart.X,
		chart.Y,
		chart.CX,
		chart.CY,
		Escape(chart.RelID),
	)
}
