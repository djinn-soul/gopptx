package pptxxml

import (
	"fmt"
	"strings"
)

func tableRowHeightsForRender(table *TableSpec, rowCount int) []int64 {
	if rowCount <= 0 {
		return nil
	}
	if len(table.RowHeights) == rowCount {
		out := make([]int64, rowCount)
		for i, height := range table.RowHeights {
			if height <= 0 {
				out[i] = 1
				continue
			}
			out[i] = height
		}
		return out
	}

	defaultHeight := table.CY / int64(rowCount)
	if defaultHeight <= 0 {
		defaultHeight = 1
	}
	out := make([]int64, rowCount)
	for i := range out {
		out[i] = defaultHeight
	}
	return out
}

func tableCellBodyPrXML(cell TableCellSpec) string {
	if cell.WrapText == nil {
		return "<a:bodyPr/>"
	}
	if *cell.WrapText {
		return `<a:bodyPr wrap="square"/>`
	}
	return `<a:bodyPr wrap="none"/>`
}

func hasTableCellMargins(cell TableCellSpec) bool {
	return cell.MarginLeft != nil || cell.MarginRight != nil || cell.MarginTop != nil || cell.MarginBottom != nil
}

func appendTableCellMarginAttrs(builder *strings.Builder, cell TableCellSpec) {
	appendMargin(builder, "marL", cell.MarginLeft)
	appendMargin(builder, "marR", cell.MarginRight)
	appendMargin(builder, "marT", cell.MarginTop)
	appendMargin(builder, "marB", cell.MarginBottom)
}

func appendMargin(builder *strings.Builder, attr string, value *int64) {
	if value == nil {
		return
	}
	builder.WriteString(` `)
	builder.WriteString(attr)
	builder.WriteString(`="`)
	_, _ = fmt.Fprintf(builder, "%d", *value)
	builder.WriteString(`"`)
}
