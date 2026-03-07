package chart

import (
	"fmt"
	"regexp"
	"strings"
)

func validateAxisTickLabelPosition(field string, value *string) error {
	if value == nil {
		return nil
	}
	switch strings.TrimSpace(*value) {
	case "nextTo", "high", "low", "none":
		return nil
	default:
		return fmt.Errorf("%s must be one of nextTo,high,low,none", field)
	}
}

func patchAxisTickLabelPosition(xml string, axisTag string, value *string) string {
	if value == nil {
		return xml
	}
	startTag := "<c:" + axisTag + ">"
	endTag := "</c:" + axisTag + ">"
	reTick := regexp.MustCompile(`<c:tickLblPos val="[^"]*"/>`)
	node := `<c:tickLblPos val="` + strings.TrimSpace(*value) + `"/>`

	start := strings.Index(xml, startTag)
	for start >= 0 {
		endRel := strings.Index(xml[start:], endTag)
		if endRel < 0 {
			break
		}
		end := start + endRel + len(endTag)
		block := xml[start:end]
		nextBlock, ok := patchAxisBlock(block, reTick, node)
		if !ok {
			start = nextAxisStart(xml[end:], startTag, end)
			continue
		}
		block = nextBlock
		xml = xml[:start] + block + xml[end:]
		nextStart := strings.Index(xml[start+len(block):], startTag)
		if nextStart < 0 {
			break
		}
		start = start + len(block) + nextStart
	}
	return xml
}

func patchAxisBlock(block string, reTick *regexp.Regexp, node string) (string, bool) {
	if reTick.MatchString(block) {
		return reTick.ReplaceAllString(block, node), true
	}
	insertAt := strings.Index(block, "<c:crosses")
	if insertAt < 0 {
		insertAt = strings.Index(block, "</c:")
	}
	if insertAt < 0 {
		return "", false
	}
	return block[:insertAt] + node + block[insertAt:], true
}

func nextAxisStart(segment string, startTag string, base int) int {
	start := strings.Index(segment, startTag)
	if start < 0 {
		return -1
	}
	return base + start
}
