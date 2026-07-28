package chart

import (
	"errors"
	"strconv"
	"strings"
)

func patchChartTitle(
	xml string,
	show *bool,
	title *string,
	overlay *bool,
	titleX, titleY *float64,
) (string, error) {
	match := reChartTitleBlock.FindString(xml)
	if show != nil && !*show {
		xml = strings.Replace(xml, match, "", 1)
		return setAutoTitleDeleted(xml, true), nil
	}
	if match == "" {
		if show != nil || title != nil || overlay != nil || titleX != nil || titleY != nil {
			inserted := insertChartTitleBlock(xml, title, overlay)
			return patchChartTitleLayoutIn(inserted, titleX, titleY), nil
		}
		return "", errors.New("chart title block not found")
	}
	block := patchExistingChartTitleBlock(match, title, overlay)
	block = patchChartTitleLayout(block, titleX, titleY)
	xml = strings.Replace(xml, match, block, 1)
	return setAutoTitleDeleted(xml, false), nil
}

func patchChartTitleLayoutIn(xml string, titleX, titleY *float64) string {
	if titleX == nil && titleY == nil {
		return xml
	}
	match := reChartTitleBlock.FindString(xml)
	if match == "" {
		return xml
	}
	return strings.Replace(xml, match, patchChartTitleLayout(match, titleX, titleY), 1)
}

// patchChartTitleLayout writes the manual title position. CT_Title orders its
// children tx, layout, overlay, spPr, txPr, so the layout is spliced after
// </c:tx> rather than appended.
func patchChartTitleLayout(block string, titleX, titleY *float64) string {
	if titleX == nil && titleY == nil {
		return block
	}

	existing := reChartTitleLayout.FindString(block)
	x, y := 0.0, 0.0
	if existing != "" {
		x, y = parseManualLayoutXY(existing)
	}
	if titleX != nil {
		x = *titleX
	}
	if titleY != nil {
		y = *titleY
	}

	layout := `<c:layout><c:manualLayout>` +
		`<c:xMode val="edge"/><c:yMode val="edge"/>` +
		`<c:x val="` + formatChartFraction(x) + `"/>` +
		`<c:y val="` + formatChartFraction(y) + `"/>` +
		`</c:manualLayout></c:layout>`

	if existing != "" {
		return strings.Replace(block, existing, layout, 1)
	}
	if strings.Contains(block, "</c:tx>") {
		return strings.Replace(block, "</c:tx>", "</c:tx>"+layout, 1)
	}
	return strings.Replace(block, "<c:title>", "<c:title>"+layout, 1)
}

func parseManualLayoutXY(layout string) (float64, float64) {
	var x, y float64
	if m := reLayoutX.FindStringSubmatch(layout); len(m) > 1 {
		x, _ = strconv.ParseFloat(m[1], 64)
	}
	if m := reLayoutY.FindStringSubmatch(layout); len(m) > 1 {
		y, _ = strconv.ParseFloat(m[1], 64)
	}
	return x, y
}

func formatChartFraction(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

func insertChartTitleBlock(xml string, title *string, overlay *bool) string {
	titleText := "Chart"
	if title != nil {
		titleText = *title
	}
	overlayVal := "0"
	if overlay != nil {
		overlayVal = boolToOneZero(*overlay)
	}
	titleBlock := `<c:title><c:tx><c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr lang="en-US"/><a:t>` +
		xmlEscape(titleText) +
		`</a:t></a:r></a:p></c:rich></c:tx><c:overlay val="` + overlayVal + `"/></c:title>`
	if strings.Contains(xml, "<c:autoTitleDeleted") {
		xml = strings.Replace(xml, "<c:autoTitleDeleted", titleBlock+"<c:autoTitleDeleted", 1)
	} else {
		xml = strings.Replace(xml, "<c:plotArea>", titleBlock+"<c:plotArea>", 1)
	}
	return setAutoTitleDeleted(xml, false)
}

func patchExistingChartTitleBlock(match string, title *string, overlay *bool) string {
	block := match
	if title != nil {
		block = patchChartTitleText(block, *title)
	}
	if overlay != nil {
		block = patchChartTitleOverlay(block, *overlay)
	}
	return block
}

func patchChartTitleText(block string, title string) string {
	escaped := `<a:t>` + xmlEscape(title) + `</a:t>`
	if reTitleText.MatchString(block) {
		return reTitleText.ReplaceAllLiteralString(block, escaped)
	}
	return strings.Replace(
		block,
		"</c:tx>",
		`<c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr lang="en-US"/>`+escaped+`</a:r></a:p></c:rich></c:tx>`,
		1,
	)
}

func patchChartTitleOverlay(block string, overlay bool) string {
	overlayNode := `<c:overlay val="` + boolToOneZero(overlay) + `"/>`
	if reOverlay.MatchString(block) {
		return reOverlay.ReplaceAllString(block, overlayNode)
	}
	return strings.Replace(block, "</c:title>", overlayNode+"</c:title>", 1)
}

func setAutoTitleDeleted(xml string, deleted bool) string {
	node := `<c:autoTitleDeleted val="` + boolToOneZero(deleted) + `"/>`
	if reAutoTitleDelete.MatchString(xml) {
		return reAutoTitleDelete.ReplaceAllString(xml, node)
	}
	return strings.Replace(xml, "<c:plotArea>", node+"<c:plotArea>", 1)
}
