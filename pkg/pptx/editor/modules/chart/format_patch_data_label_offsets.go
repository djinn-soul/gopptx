package chart

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	reSerBlocks       = regexp.MustCompile(`(?s)<c:ser>.*?</c:ser>`)
	reSeriesDLblsOpen = regexp.MustCompile(`<c:dLbls\b[^>]*>`)
)

// patchDataLabelOffsets writes a manual layout on individual data labels.
//
// Doughnut and pie charts reject most c:dLblPos values — setting one makes
// PowerPoint report the file as unreadable — so moving a single label means
// giving its <c:dLbl> a <c:layout><c:manualLayout> with an offset, which is what
// PowerPoint itself records when a label is dragged (upstream issue #1025).
func patchDataLabelOffsets(xml string, offsets []common.DataLabelOffset) string {
	if len(offsets) == 0 {
		return xml
	}

	bySeries := map[int][]common.DataLabelOffset{}
	for _, offset := range offsets {
		bySeries[offset.SeriesIndex] = append(bySeries[offset.SeriesIndex], offset)
	}

	// Labels are usually enabled on the plot-level <c:dLbls>, outside <c:ser>,
	// so that block supplies the flags when the series carries none.
	plotFlags := dataLabelFlagsFrom(seriesDataLabelsBlock(stripSeriesBlocks(xml)))

	seriesIndex := -1
	return reSerBlocks.ReplaceAllStringFunc(xml, func(ser string) string {
		seriesIndex++
		seriesOffsets, ok := bySeries[seriesIndex]
		if !ok {
			return ser
		}
		return applySeriesDataLabelOffsets(ser, seriesOffsets, plotFlags)
	})
}

// stripSeriesBlocks removes every <c:ser> so plot-level elements can be read
// without a series' own copies shadowing them.
func stripSeriesBlocks(xml string) string {
	return reSerBlocks.ReplaceAllString(xml, "")
}

func applySeriesDataLabelOffsets(
	ser string,
	offsets []common.DataLabelOffset,
	plotFlags string,
) string {
	// Highest point index first: each <c:dLbl> is inserted at the start of the
	// dLbls block, so descending order leaves them ascending in the output.
	sorted := append([]common.DataLabelOffset(nil), offsets...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].PointIndex > sorted[j].PointIndex })

	for _, offset := range sorted {
		ser = applyDataLabelOffset(ser, offset, plotFlags)
	}
	return ser
}

func applyDataLabelOffset(ser string, offset common.DataLabelOffset, plotFlags string) string {
	if offset.X == nil && offset.Y == nil {
		return ser
	}

	existing, start, end := findDataLabelBlock(ser, offset.PointIndex)
	x, y := 0.0, 0.0
	if existing != "" {
		x, y = parseManualLayoutXY(reChartTitleLayout.FindString(existing))
	}
	if offset.X != nil {
		x = *offset.X
	}
	if offset.Y != nil {
		y = *offset.Y
	}

	if existing != "" {
		patched := patchExistingDataLabelLayout(existing, x, y)
		return ser[:start] + patched + ser[end:]
	}
	flags := seriesDataLabelFlags(ser, existing, plotFlags)
	label := buildDataLabelBlock(offset.PointIndex, x, y, flags)
	return insertDataLabelBlock(ser, label)
}

func patchExistingDataLabelLayout(existing string, x, y float64) string {
	layout := buildManualDataLabelLayout(x, y)
	if current := reChartTitleLayout.FindString(existing); current != "" {
		return strings.Replace(existing, current, layout, 1)
	}

	insertAt := strings.Index(existing, "<c:idx")
	if insertAt >= 0 {
		if idxEnd := strings.Index(existing[insertAt:], "/>"); idxEnd >= 0 {
			insertAt += idxEnd + len("/>")
		} else if idxEnd := strings.Index(existing[insertAt:], "</c:idx>"); idxEnd >= 0 {
			insertAt += idxEnd + len("</c:idx>")
		} else {
			insertAt = -1
		}
	}
	if deleteAt := strings.Index(existing, "<c:delete"); deleteAt >= 0 {
		if deleteEnd := strings.Index(existing[deleteAt:], "/>"); deleteEnd >= 0 {
			insertAt = deleteAt + deleteEnd + len("/>")
		}
	}
	if insertAt < 0 {
		return existing
	}
	return existing[:insertAt] + layout + existing[insertAt:]
}

// dataLabelFlagNames returns the CT_DLbl display flags, in schema order.
func dataLabelFlagNames() []string {
	return []string{
		flagShowLegendKey, flagShowValue, flagShowCategory,
		flagShowSeriesName, flagShowPercent, flagShowBubbleSize,
	}
}

// seriesDataLabelFlags returns the display flags a per-point label must repeat.
//
// A <c:dLbl> does not inherit the surrounding flags: whatever it omits is off.
// They are taken from the label being replaced, else the series <c:dLbls>, else
// the plot-level one, and only fall back to "value only" when there is none.
func seriesDataLabelFlags(ser string, existing string, plotFlags string) string {
	for _, source := range []string{existing, seriesDataLabelsBlock(ser)} {
		if flags := dataLabelFlagsFrom(source); flags != "" {
			return completeDataLabelFlags(flags)
		}
	}
	if plotFlags != "" {
		return completeDataLabelFlags(plotFlags)
	}
	return `<c:showLegendKey val="0"/><c:showVal val="1"/><c:showCatName val="0"/>` +
		`<c:showSerName val="0"/><c:showPercent val="0"/><c:showBubbleSize val="0"/>`
}

// completeDataLabelFlags fills in the flags the source omitted, defaulting them
// to off. PowerPoint ignores a <c:dLbl> that carries only some of the six
// display flags — label, number format and font all revert to the parent's —
// so an inherited partial set has to be completed before it is written back.
func completeDataLabelFlags(flags string) string {
	values := dataLabelFlagValues(flags)
	var b strings.Builder
	for _, name := range dataLabelFlagNames() {
		value, ok := values[name]
		if !ok {
			value = "0"
		}
		b.WriteString(`<c:` + name + ` val="` + value + `"/>`)
	}
	return b.String()
}

// dataLabelFlagsFrom copies the display flags out of a dLbls or dLbl block,
// preserving schema order.
func dataLabelFlagsFrom(source string) string {
	if source == "" {
		return ""
	}
	values := dataLabelFlagValues(source)
	var b strings.Builder
	for _, flag := range dataLabelFlagNames() {
		if value, ok := values[flag]; ok {
			b.WriteString(`<c:` + flag + ` val="` + value + `"/>`)
		}
	}
	return b.String()
}

// seriesDataLabelsBlock returns the series <c:dLbls> content without its
// per-point <c:dLbl> children, so their flags are not mistaken for the
// series-wide ones.
func seriesDataLabelsBlock(ser string) string {
	block := reDataLabelsBlock.FindString(ser)
	if block == "" {
		return ""
	}
	for {
		start := strings.Index(block, "<c:dLbl>")
		if start < 0 {
			return block
		}
		closeRel := strings.Index(block[start:], "</c:dLbl>")
		if closeRel < 0 {
			return block
		}
		block = block[:start] + block[start+closeRel+len("</c:dLbl>"):]
	}
}

// buildDataLabelBlock renders one <c:dLbl>. CT_DLbl orders idx, then the
// delete/layout group, then the display flags.
func buildDataLabelBlock(pointIndex int, x, y float64, flags string) string {
	return `<c:dLbl><c:idx val="` + strconv.Itoa(pointIndex) + `"/>` +
		buildManualDataLabelLayout(x, y) +
		flags +
		`</c:dLbl>`
}

func buildManualDataLabelLayout(x, y float64) string {
	return `<c:layout><c:manualLayout>` +
		`<c:x val="` + formatChartFraction(x) + `"/>` +
		`<c:y val="` + formatChartFraction(y) + `"/>` +
		`</c:manualLayout></c:layout>`
}

// findDataLabelBlock locates the <c:dLbl> for pointIndex, returning it with its
// bounds in ser. The block is empty when the point has no explicit label.
func findDataLabelBlock(ser string, pointIndex int) (string, int, int) {
	needle := `<c:idx val="` + strconv.Itoa(pointIndex) + `"/>`
	cursor := 0
	for {
		rel := strings.Index(ser[cursor:], "<c:dLbl>")
		if rel < 0 {
			return "", 0, 0
		}
		start := cursor + rel
		closeRel := strings.Index(ser[start:], "</c:dLbl>")
		if closeRel < 0 {
			return "", 0, 0
		}
		end := start + closeRel + len("</c:dLbl>")
		block := ser[start:end]
		if strings.Contains(block, needle) {
			return block, start, end
		}
		cursor = end
	}
}

// insertDataLabelBlock puts a <c:dLbl> at the front of the series dLbls block,
// creating that block when the series has none. CT_DLbls requires the
// per-point labels before the series-wide flags.
func insertDataLabelBlock(ser string, label string) string {
	if loc := reSeriesDLblsOpen.FindStringIndex(ser); loc != nil {
		if strings.HasSuffix(ser[loc[0]:loc[1]], "/>") {
			// Self-closed <c:dLbls/>: expand it so the label has a home.
			return ser[:loc[0]] + `<c:dLbls>` + label + `</c:dLbls>` + ser[loc[1]:]
		}
		return ser[:loc[1]] + label + ser[loc[1]:]
	}

	// No dLbls yet. CT_Ser puts dLbls after tx/spPr and before the data refs.
	for _, anchor := range []string{"<c:cat>", "<c:xVal>", "<c:val>"} {
		if idx := strings.Index(ser, anchor); idx >= 0 {
			return ser[:idx] + `<c:dLbls>` + label + `</c:dLbls>` + ser[idx:]
		}
	}
	return ser
}
