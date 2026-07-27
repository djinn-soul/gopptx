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

// tableCellBorderXML renders one a:lnL/R/T/B. CT_LineProperties fixes both the
// attribute set and the child order: fill, then prstDash, then the join.
func tableCellBorderXML(side string, border TableCellBorderSpec) string {
	dash := tableCellBorderDash(border.Dash)

	var b strings.Builder
	b.WriteString(`<a:` + side + ` w="` + strconv.FormatInt(border.Width, 10) + `"`)
	if lineCap := tableCellBorderCap(border.Cap); lineCap != "" {
		b.WriteString(` cap="` + lineCap + `"`)
	}
	if compound := tableCellBorderCompound(border.Compound); compound != "" {
		b.WriteString(` cmpd="` + compound + `"`)
	}
	if border.Inset {
		b.WriteString(` algn="in"`)
	}
	b.WriteString(`>`)
	b.WriteString(`<a:solidFill><a:srgbClr val="` + Escape(border.Color) + `"/></a:solidFill>`)
	b.WriteString(`<a:prstDash val="` + Escape(dash) + `"/>`)
	b.WriteString(tableCellBorderJoinXML(border))
	b.WriteString(`</a:` + side + `>`)
	return b.String()
}

func tableCellBorderCap(lineCap string) string {
	switch strings.ToLower(strings.TrimSpace(lineCap)) {
	case string(LineCapStyleRound), string(LineJoinStyleRound):
		return string(LineCapStyleRound)
	case string(LineCapStyleSquare), lineCapSquareLong:
		return string(LineCapStyleSquare)
	default:
		return "" // flat is the default; no attribute needed
	}
}

// Compound line types (ST_CompoundLine).
const (
	compoundDouble    = "dbl"
	compoundThickThin = "thickThin"
	compoundThinThick = "thinThick"
	compoundTriple    = "tri"
	lineCapSquareLong = "square"
)

func tableCellBorderCompound(compound string) string {
	switch strings.ToLower(strings.TrimSpace(compound)) {
	case compoundDouble, "double":
		return compoundDouble
	case "thickthin", "thick-thin", "thick_thin":
		return compoundThickThin
	case "thinthick", "thin-thick", "thin_thick":
		return compoundThinThick
	case compoundTriple, "triple":
		return compoundTriple
	default:
		return "" // sng is the default; no attribute needed
	}
}

const defaultMiterLimitPct = 800000.0

func tableCellBorderJoinXML(border TableCellBorderSpec) string {
	switch strings.ToLower(strings.TrimSpace(border.Join)) {
	case string(LineJoinStyleRound), string(LineCapStyleRound):
		return joinRoundXML
	case string(LineJoinStyleBevel):
		return joinBevelXML
	case string(LineJoinStyleMiter):
		limit := border.MiterLimitPct
		if limit <= 0 {
			limit = defaultMiterLimitPct
		}
		return `<a:miter lim="` + strconv.FormatInt(int64(limit), 10) + `"/>`
	default:
		return ""
	}
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
