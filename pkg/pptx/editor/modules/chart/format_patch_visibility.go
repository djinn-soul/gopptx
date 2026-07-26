package chart

import "strings"

func patchPlotVisibleOnly(xml string, value *bool) string {
	if value == nil {
		return xml
	}
	node := `<c:plotVisOnly val="` + boolToOneZero(*value) + `"/>`
	if rePlotVisOnly.MatchString(xml) {
		return rePlotVisOnly.ReplaceAllString(xml, node)
	}
	return strings.Replace(xml, "</c:chart>", node+"</c:chart>", 1)
}
