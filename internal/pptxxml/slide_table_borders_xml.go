package pptxxml

import (
	"strconv"
	"strings"
)

func tableCellBordersForRender(cell TableCellSpec) tableCellBorderSet {
	borders := tableCellBorderSet{
		Left:   cloneTableCellBorderSpec(cell.BorderLeft),
		Right:  cloneTableCellBorderSpec(cell.BorderRight),
		Top:    cloneTableCellBorderSpec(cell.BorderTop),
		Bottom: cloneTableCellBorderSpec(cell.BorderBottom),
	}
	if borders.Left == nil && borders.Right == nil && borders.Top == nil && borders.Bottom == nil {
		if cell.BorderWidth > 0 && strings.TrimSpace(cell.BorderColor) != "" {
			legacy := &TableCellBorderSpec{Width: cell.BorderWidth, Color: cell.BorderColor, Dash: strokeDashSolid}
			borders.Left = cloneTableCellBorderSpec(legacy)
			borders.Right = cloneTableCellBorderSpec(legacy)
			borders.Top = cloneTableCellBorderSpec(legacy)
			borders.Bottom = cloneTableCellBorderSpec(legacy)
		}
	}
	return borders
}

func cloneTableCellBorderSpec(border *TableCellBorderSpec) *TableCellBorderSpec {
	if border == nil {
		return nil
	}
	clone := *border
	return &clone
}

func tableCellBorderXML(side string, border TableCellBorderSpec) string {
	dash := tableCellBorderDash(border.Dash)
	return `<a:` + side +
		` w="` + strconv.FormatInt(border.Width, 10) +
		`"><a:solidFill><a:srgbClr val="` + Escape(border.Color) +
		`"/></a:solidFill><a:prstDash val="` + Escape(dash) + `"/></a:` + side + `>`
}

func tableCellBorderDash(dash string) string {
	switch strings.ToLower(strings.TrimSpace(dash)) {
	case "", strokeDashSolid:
		return strokeDashSolid
	case strokeDashDash:
		return strokeDashDash
	case "dot":
		return "dot"
	case "dashdot", "dash-dot", "dash_dot":
		return "dashDot"
	case "lgdash", "lg-dash", "longdash", "long-dash", "long_dash":
		return "lgDash"
	default:
		return strings.TrimSpace(dash)
	}
}
