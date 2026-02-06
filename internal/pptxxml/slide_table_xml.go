package pptxxml

import (
	"fmt"
	"strings"
)

// TableSpec describes one table in a slide.
type TableSpec struct {
	X            int64
	Y            int64
	CX           int64
	CY           int64
	ColumnWidths []int64
	Rows         [][]string
}

func tableShape(table *TableSpec, shapeID int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
<p:graphicFrame>
<p:nvGraphicFramePr>
<p:cNvPr id="%d" name="Table 1"/>
<p:cNvGraphicFramePr><a:graphicFrameLocks noGrp="1"/></p:cNvGraphicFramePr>
<p:nvPr/>
</p:nvGraphicFramePr>
<p:xfrm>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</p:xfrm>
<a:graphic>
<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/table">
<a:tbl>
<a:tblPr firstRow="1" bandRow="1">
<a:tableStyleId>{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}</a:tableStyleId>
</a:tblPr>
<a:tblGrid>`,
		shapeID,
		table.X,
		table.Y,
		table.CX,
		table.CY,
	))

	for _, width := range table.ColumnWidths {
		b.WriteString(fmt.Sprintf(`
<a:gridCol w="%d"/>`, width))
	}

	b.WriteString(`
</a:tblGrid>`)

	rowHeight := table.CY / int64(len(table.Rows))
	if rowHeight <= 0 {
		rowHeight = 1
	}
	for _, row := range table.Rows {
		b.WriteString(fmt.Sprintf(`
<a:tr h="%d">`, rowHeight))
		for _, cell := range row {
			b.WriteString(fmt.Sprintf(`
<a:tc>
<a:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p>
<a:r>
<a:rPr lang="en-US" dirty="0"/>
<a:t>%s</a:t>
</a:r>
</a:p>
</a:txBody>
<a:tcPr/>
</a:tc>`, Escape(cell)))
		}
		b.WriteString(`
</a:tr>`)
	}

	b.WriteString(`
</a:tbl>
</a:graphicData>
</a:graphic>
</p:graphicFrame>`)
	return b.String()
}
