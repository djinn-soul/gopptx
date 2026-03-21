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
