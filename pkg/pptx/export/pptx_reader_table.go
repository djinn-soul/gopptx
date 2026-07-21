package export

import (
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

//nolint:gocognit
func extractTableContent(
	ed *editor.PresentationEditor,
	slideIndex int,
	shape editorcommon.Shape,
) *tables.Table {
	if shape.ID <= 0 {
		return nil
	}
	info, err := ed.GetTable(slideIndex, shape.ID)
	if err != nil {
		return nil
	}
	tableAny, ok := info["table"]
	if !ok {
		return nil
	}
	tableMeta, ok := tableAny.(map[string]any)
	if !ok {
		return nil
	}
	rowsMeta := toMapSlice(tableMeta["rows"])
	if len(rowsMeta) == 0 {
		return nil
	}
	firstRow, _ := tableMeta["first_row"].(bool)
	bandRow, _ := tableMeta["band_row"].(bool)
	colWidths := toLengthSlice(tableMeta["column_widths"])
	rowHeights := toLengthSlice(tableMeta["row_heights"])
	table := tables.NewTable(colWidths)
	table = table.Position(styling.Emu(int64(shape.X)), styling.Emu(int64(shape.Y)))
	table = table.Size(styling.Emu(int64(shape.W)), styling.Emu(int64(shape.H)))
	if shape.AltText != "" {
		table = table.WithAltText(shape.AltText)
	}
	if shape.IsDecorative {
		table = table.WithDecorative(true)
	}
	if len(rowHeights) > 0 {
		table = table.WithRowHeights(rowHeights)
	}
	for rowIndex, row := range rowsMeta {
		cellsMeta := toMapSlice(row["cells"])
		if len(cellsMeta) == 0 {
			continue
		}
		values := make([]string, 0, len(cellsMeta))
		styled := make([]tables.TableCell, 0, len(cellsMeta))
		for _, cell := range cellsMeta {
			value, _ := cell["text"].(string)
			values = append(values, value)
			styledCell := tables.NewTableCell(value)
			applyTableCellLayout(&styledCell, cell)
			applyTableCellBorders(&styledCell, cell)
			if firstRow && rowIndex == 0 {
				styledCell = styledCell.WithBold(true).WithBackgroundColor("4472C4")
				styledCell.Color = "FFFFFF"
			} else if bandRow && rowIndex%2 == 1 {
				styledCell = styledCell.WithBackgroundColor("D0D8E8")
			}
			styled = append(styled, styledCell)
		}
		if firstRow || bandRow {
			table = table.AddStyledRow(styled)
			continue
		}
		table = table.AddRow(values)
	}
	return &table
}

func applyTableCellLayout(cell *tables.TableCell, cellMeta map[string]any) {
	if rawAlign, _ := cellMeta["v_align"].(string); rawAlign != "" {
		switch strings.ToLower(strings.TrimSpace(rawAlign)) {
		case "top":
			*cell = cell.WithVAlignTop()
		case anchorBottom:
			*cell = cell.WithVAlignBottom()
		case "middle", "center", "ctr":
			*cell = cell.WithVAlignMiddle()
		}
	}
	applyMargin := func(
		margin any,
		setter func(float64) tables.TableCell,
	) {
		emu, ok := parseNumericInt64(margin)
		if !ok || emu <= 0 {
			return
		}
		pts := math.Max(emuToPt(emu), 0)
		*cell = setter(pts)
	}
	applyMargin(cellMeta["margin_left"], cell.WithMarginLeftPt)
	applyMargin(cellMeta["margin_right"], cell.WithMarginRightPt)
	applyMargin(cellMeta["margin_top"], cell.WithMarginTopPt)
	applyMargin(cellMeta["margin_bottom"], cell.WithMarginBottomPt)
}

func applyTableCellBorders(cell *tables.TableCell, cellMeta map[string]any) {
	applyBorder := func(
		key string,
		setter func(float64, string) tables.TableCell,
	) {
		borderMeta, ok := cellMeta[key].(map[string]any)
		if !ok {
			return
		}
		widthEmu, ok := parseNumericInt64(borderMeta["width"])
		if !ok || widthEmu <= 0 {
			return
		}
		color := "B4B4B4"
		if c, ok := borderMeta["color"].(string); ok && strings.TrimSpace(c) != "" {
			color = c
		}
		*cell = setter(math.Max(emuToPt(widthEmu), minTableBorderWidthPt), color)
	}
	applyBorder("border_left", cell.WithLeftBorder)
	applyBorder("border_right", cell.WithRightBorder)
	applyBorder("border_top", cell.WithTopBorder)
	applyBorder("border_bottom", cell.WithBottomBorder)
}

func parseNumericInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}

func toMapSlice(value any) []map[string]any {
	switch items := value.(type) {
	case []map[string]any:
		return items
	case []any:
		out := make([]map[string]any, 0, len(items))
		for _, item := range items {
			m, ok := item.(map[string]any)
			if ok {
				out = append(out, m)
			}
		}
		return out
	default:
		return nil
	}
}

func toLengthSlice(value any) []styling.Length {
	switch values := value.(type) {
	case []int64:
		out := make([]styling.Length, 0, len(values))
		for _, v := range values {
			out = append(out, styling.Emu(v))
		}
		return out
	case []any:
		out := make([]styling.Length, 0, len(values))
		for _, raw := range values {
			switch v := raw.(type) {
			case int64:
				out = append(out, styling.Emu(v))
			case int:
				out = append(out, styling.Emu(int64(v)))
			case float64:
				out = append(out, styling.Emu(int64(v)))
			}
		}
		return out
	default:
		return nil
	}
}
