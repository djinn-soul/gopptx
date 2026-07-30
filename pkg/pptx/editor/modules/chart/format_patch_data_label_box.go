package chart

import (
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

// The data-label box, upstream #662 and #716: the background fill and the
// outline of a label live in the c:spPr of a c:dLbls or a single c:dLbl, which
// is separate from the c:txPr carrying its font.

func validateChartDataLabelBox(req common.ChartFormatUpdate) error {
	if req.DataLabelFillColor != nil {
		if _, err := editorshape.NormalizeHexColor(*req.DataLabelFillColor); err != nil {
			return fmt.Errorf("data_label_fill_color: %w", err)
		}
	}
	return validateChartLineFormat("data_label_border", req.DataLabelBorder)
}

func validateDataLabelPointBox(point common.ChartDataLabelPoint) error {
	if point.FillColor != nil {
		if _, err := editorshape.NormalizeHexColor(*point.FillColor); err != nil {
			return fmt.Errorf("data label point fill_color: %w", err)
		}
	}
	return validateChartLineFormat("data label point border", point.Border)
}

// patchChartDataLabelBox writes the fill and outline of every plot-wide c:dLbls.
// A chart with no labels yet gets the default block first, since the styling
// has nowhere else to live.
func patchChartDataLabelBox(xml string, req common.ChartFormatUpdate) string {
	fill := chartFillNode(req.DataLabelFillColor, req.DataLabelNoFill)
	line := buildChartLine(req.DataLabelBorder)
	if fill == "" && line == "" {
		return xml
	}
	if !reDataLabelsBlock.MatchString(xml) {
		xml = insertDefaultDataLabels(xml)
	}
	return reDataLabelsBlock.ReplaceAllStringFunc(xml, func(block string) string {
		return writeDataLabelShapeProperties(block, fill, line)
	})
}

// writeDataLabelShapeProperties merges a fill and a line into the c:spPr of a
// label block. CT_DLbls orders numFmt, spPr, txPr, dLblPos and the display
// flags, so a new c:spPr is spliced in after the number format.
func writeDataLabelShapeProperties(block string, fill string, line string) string {
	current := reDataLabelShapePr.FindString(block)
	if current == "" || strings.HasSuffix(current, "/>") {
		node := "<c:spPr>" + fill + line + "</c:spPr>"
		if current != "" {
			return strings.Replace(block, current, node, 1)
		}
		return strings.Replace(block, dataLabelShapePropsAnchor(block), dataLabelShapePropsAnchor(block)+node, 1)
	}
	updated := current
	if fill != "" {
		updated = setShapePropertiesFill(updated, fill)
	}
	updated = setShapePropertiesLine(updated, line)
	return strings.Replace(block, current, updated, 1)
}

// dataLabelShapePropsAnchor returns the text a new c:spPr follows: the number
// format when the block has one, else the opening tag.
func dataLabelShapePropsAnchor(block string) string {
	if numFmt := reDataLabelNumFmt.FindString(block); numFmt != "" {
		return numFmt
	}
	return "<c:dLbls>"
}

// buildDataLabelShapeProperties renders the c:spPr of one c:dLbl, keeping what
// the label already carried when the request does not mention fill or border.
func buildDataLabelShapeProperties(existing string, point common.ChartDataLabelPoint) string {
	current := reDataLabelShapePr.FindString(existing)
	fill := chartFillNode(point.FillColor, point.NoFill)
	line := buildChartLine(point.Border)
	if fill == "" && line == "" {
		return current
	}
	if current == "" || strings.HasSuffix(current, "/>") {
		return "<c:spPr>" + fill + line + "</c:spPr>"
	}
	updated := current
	if fill != "" {
		updated = setShapePropertiesFill(updated, fill)
	}
	return setShapePropertiesLine(updated, line)
}
