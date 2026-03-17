package table

import "bytes"

func buildTableCellInfo(rowIndex, colIndex int, cell CellXML) map[string]any {
	rowSpan := normalizeSpan(cell.RowSpan)
	colSpan := normalizeSpan(cell.GridSpan)
	vMerge := TruthyAttr(cell.VMerge)
	hMerge := TruthyAttr(cell.HMerge)
	info := map[string]any{
		"row":             rowIndex,
		"col":             colIndex,
		"row_span":        rowSpan,
		"col_span":        colSpan,
		"v_merge":         vMerge,
		"h_merge":         hMerge,
		"is_merge_origin": rowSpan > 1 || colSpan > 1,
		"is_spanned":      vMerge || hMerge,
		"text":            tableCellText(cell),
	}
	applyTableCellLayoutInfo(info, cell.TcPr)
	applyTableCellBorderInfo(info, cell.TcPr)
	return info
}

func applyTableCellLayoutInfo(info map[string]any, props cellPropertiesXML) {
	if align := normalizeTableCellVAlign(props.Anchor); align != "" {
		info["v_align"] = align
	}
	if props.MarL != nil {
		info["margin_left"] = *props.MarL
	}
	if props.MarR != nil {
		info["margin_right"] = *props.MarR
	}
	if props.MarT != nil {
		info["margin_top"] = *props.MarT
	}
	if props.MarB != nil {
		info["margin_bottom"] = *props.MarB
	}
}

func normalizeTableCellVAlign(anchor string) string {
	switch anchor {
	case "t":
		return "top"
	case "b":
		return "bottom"
	case "ctr":
		return "middle"
	default:
		return ""
	}
}

func applyTableCellBorderInfo(info map[string]any, props cellPropertiesXML) {
	if border := lineToBorderInfo(props.LnL); border != nil {
		info["border_left"] = border
	}
	if border := lineToBorderInfo(props.LnR); border != nil {
		info["border_right"] = border
	}
	if border := lineToBorderInfo(props.LnT); border != nil {
		info["border_top"] = border
	}
	if border := lineToBorderInfo(props.LnB); border != nil {
		info["border_bottom"] = border
	}
}

func lineToBorderInfo(line *linePropertiesXML) map[string]any {
	if line == nil || line.NoFill != nil {
		return nil
	}
	border := map[string]any{}
	if line.Width > 0 {
		border["width"] = line.Width
	}
	if line.PrstDash != nil && line.PrstDash.Val != "" {
		border["dash"] = line.PrstDash.Val
	}
	if color := lineColorToken(line); color != "" {
		border["color"] = color
	}
	if len(border) == 0 {
		return nil
	}
	return border
}

func lineColorToken(line *linePropertiesXML) string {
	if line == nil || line.SolidFill == nil {
		return ""
	}
	if line.SolidFill.SrgbClr != nil && line.SolidFill.SrgbClr.Val != "" {
		return line.SolidFill.SrgbClr.Val
	}
	scheme := line.SolidFill.SchemeClr
	if scheme == nil || scheme.Val == "" {
		return ""
	}
	token := "scheme:" + scheme.Val
	if scheme.LumMod != nil && scheme.LumMod.Val != "" {
		token += "|lumMod=" + scheme.LumMod.Val
	}
	if scheme.LumOff != nil && scheme.LumOff.Val != "" {
		token += "|lumOff=" + scheme.LumOff.Val
	}
	if scheme.Tint != nil && scheme.Tint.Val != "" {
		token += "|tint=" + scheme.Tint.Val
	}
	if scheme.Shade != nil && scheme.Shade.Val != "" {
		token += "|shade=" + scheme.Shade.Val
	}
	return token
}

func normalizeSpan(span int) int {
	if span <= 0 {
		return 1
	}
	return span
}

func tableCellText(cell CellXML) string {
	var textBuf bytes.Buffer
	for i, p := range cell.TxBody.Paragraphs {
		if i > 0 {
			textBuf.WriteString("\n")
		}
		for _, r := range p.Runs {
			textBuf.WriteString(r.Text)
		}
	}
	return textBuf.String()
}
