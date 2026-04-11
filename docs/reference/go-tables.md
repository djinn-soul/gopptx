# Go Tables Reference

This page documents the table helpers exposed by `pkg/pptx`.

Primary source files:

- `pkg/pptx/table_compat.go`
- `pkg/pptx/tables/table.go`
- `pkg/pptx/tables/table_cell.go`

## Constructors

### `NewTable(columnWidths []Length) Table`

Create a table with explicit column widths.

### `NewTableCell(text string) TableCell`

Create a table cell with text content.

## Table fluent methods

- `WithAltText(text string) Table` — accessibility alt text
- `WithDecorative(enabled bool) Table` — mark as decorative (screen readers skip it)
- `AddColumn(width Length) Table` — append a column
- `AddRow(cells []string) Table` — append a row of plain-text cells
- `AddStyledRow(cells []TableCell) Table` — append a row of styled cells
- `WithData(data [][]string) Table` — set all rows from a 2-D string slice
- `WithStyledData(data [][]TableCell) Table` — set all rows from a 2-D cell slice
- `WithRowHeights(heights []Length) Table` — set per-row heights
- `Position(x, y Length) Table` — set position on the slide
- `Size(cx, cy Length) Table` — set total table size

## TableCell fluent methods

### Text and font

- `WithBold(enabled bool) TableCell`
- `WithSizePt(pt float64) TableCell`
- `WithFontName(name string) TableCell`

### Color and alignment

- `WithBackgroundColor(color string) TableCell`
- `WithAlign(align string) TableCell`
- `WithAlignLeft() TableCell`
- `WithAlignCenter() TableCell`
- `WithAlignRight() TableCell`
- `WithAlignJustify() TableCell`
- `WithVAlign(vAlign string) TableCell`
- `WithVAlignTop() TableCell`
- `WithVAlignMiddle() TableCell`
- `WithVAlignBottom() TableCell`

### Span

- `WithRowSpan(span int) TableCell`
- `WithColSpan(span int) TableCell`

### Borders

- `WithBorder(widthPt float64, color string) TableCell` — uniform border on all sides
- `WithBorderStyle(widthPt float64, color string, dash string) TableCell` — uniform with dash style
- `WithLeftBorder(widthPt float64, color string) TableCell`
- `WithLeftBorderStyle(widthPt float64, color string, dash string) TableCell`
- `WithRightBorder(widthPt float64, color string) TableCell`
- `WithRightBorderStyle(widthPt float64, color string, dash string) TableCell`
- `WithTopBorder(widthPt float64, color string) TableCell`
- `WithTopBorderStyle(widthPt float64, color string, dash string) TableCell`
- `WithBottomBorder(widthPt float64, color string) TableCell`
- `WithBottomBorderStyle(widthPt float64, color string, dash string) TableCell`

## Table constants

- `TableAlignLeft`
- `TableAlignCenter`
- `TableAlignRight`
- `TableAlignJustify`
- `TableVAlignTop`
- `TableVAlignMiddle`
- `TableVAlignBottom`
- `TableBorderDashSolid`
- `TableBorderDashDash`
- `TableBorderDashDot`
- `TableBorderDashDashDot`
- `TableBorderDashLongDash`

## Typical usage

```go
table := pptx.NewTable([]pptx.Length{pptx.Inches(2.6), pptx.Inches(2.6), pptx.Inches(2.6)})
table = table.WithStyledData([][]pptx.TableCell{
    {
        pptx.NewTableCell("Workstream"),
        pptx.NewTableCell("Owner"),
        pptx.NewTableCell("Status"),
    },
})
```

See also:

- [Go API Reference](go-api.md)
