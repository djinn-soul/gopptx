package chart

import (
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// patchDataLabelFlags rewrites the display flags of one c:dLbls as a complete
// set in schema order. CT_DLbls puts them after numFmt, spPr, txPr and dLblPos
// and before c:separator, and a partial set makes PowerPoint fall back to its
// own defaults, which is how a per-point label loses its number format.
func patchDataLabelFlags(block string, req common.ChartFormatUpdate) string {
	overrides := map[string]*bool{
		flagShowLegendKey:  req.DataLabelShowLegendKey,
		flagShowValue:      req.DataLabelShowValue,
		flagShowCategory:   req.DataLabelShowCategory,
		flagShowSeriesName: req.DataLabelShowSeriesName,
		flagShowPercent:    req.DataLabelShowPercent,
		flagShowBubbleSize: req.DataLabelShowBubbleSize,
	}

	values := dataLabelFlagValues(block)
	present := len(values) > 0
	for name, override := range overrides {
		if override == nil {
			continue
		}
		values[name] = boolToOneZero(*override)
		present = true
	}
	if !present {
		return block
	}

	block = reDataLabelFlag.ReplaceAllLiteralString(block, "")
	var flags strings.Builder
	for _, name := range dataLabelFlagNames() {
		value, ok := values[name]
		if !ok {
			value = "0"
		}
		flags.WriteString(`<c:` + name + ` val="` + value + `"/>`)
	}

	if index := strings.Index(block, "<c:separator>"); index >= 0 {
		return block[:index] + flags.String() + block[index:]
	}
	return strings.Replace(block, "</c:dLbls>", flags.String()+"</c:dLbls>", 1)
}

func insertDefaultDataLabels(xml string) string {
	start, end := firstChartBlockBounds(xml)
	if start < 0 || end <= start {
		return xml
	}
	chartBlock := xml[start:end]
	insertAt := strings.Index(chartBlock, "<c:axId")
	if insertAt < 0 {
		insertAt = strings.LastIndex(chartBlock, chartElementClosePrefix)
		if insertAt < 0 {
			return xml
		}
	}
	labels := `<c:dLbls><c:showVal val="1"/></c:dLbls>`
	patched := chartBlock[:insertAt] + labels + chartBlock[insertAt:]
	return xml[:start] + patched + xml[end:]
}
