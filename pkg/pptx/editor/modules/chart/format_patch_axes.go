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

func validateAxisCrosses(field string, value *string) error {
	if value == nil {
		return nil
	}
	switch strings.TrimSpace(*value) {
	case "autoZero", "max", "min":
		return nil
	default:
		return fmt.Errorf("%s must be one of autoZero,max,min", field)
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

func patchAxisMajorGridlines(xml string, axisTag string, value *bool) string {
	reMajor := regexp.MustCompile(`(?s)<c:majorGridlines(?:\s*/>|>.*?</c:majorGridlines>)`)
	node := `<c:majorGridlines/>`
	return patchAxisGridlines(xml, axisTag, value, reMajor, node)
}

func patchAxisMinorGridlines(xml string, axisTag string, value *bool) string {
	reMinor := regexp.MustCompile(`(?s)<c:minorGridlines(?:\s*/>|>.*?</c:minorGridlines>)`)
	node := `<c:minorGridlines/>`
	return patchAxisGridlines(xml, axisTag, value, reMinor, node)
}

func patchAxisGridlines(
	xml string,
	axisTag string,
	value *bool,
	re *regexp.Regexp,
	node string,
) string {
	if value == nil {
		return xml
	}
	startTag := "<c:" + axisTag + ">"
	endTag := "</c:" + axisTag + ">"

	start := strings.Index(xml, startTag)
	for start >= 0 {
		endRel := strings.Index(xml[start:], endTag)
		if endRel < 0 {
			break
		}
		end := start + endRel + len(endTag)
		block := xml[start:end]
		nextBlock, ok := patchAxisGridBlock(block, *value, re, node)
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

func patchAxisGridBlock(
	block string,
	enable bool,
	reGrid *regexp.Regexp,
	node string,
) (string, bool) {
	if !enable {
		return reGrid.ReplaceAllString(block, ""), true
	}
	if reGrid.MatchString(block) {
		return block, true
	}
	insertAt := axisNodeInsertIndex(block)
	if insertAt < 0 {
		return "", false
	}
	return block[:insertAt] + node + block[insertAt:], true
}

func axisNodeInsertIndex(block string) int {
	insertAt := strings.Index(block, "<c:tickLblPos")
	if insertAt >= 0 {
		return insertAt
	}
	insertAt = strings.Index(block, "<c:crosses")
	if insertAt >= 0 {
		return insertAt
	}
	return strings.Index(block, "</c:")
}

func patchAxisCrosses(xml string, axisTag string, value *string) string {
	if value == nil {
		return xml
	}
	startTag := "<c:" + axisTag + ">"
	endTag := "</c:" + axisTag + ">"
	reCrosses := regexp.MustCompile(`<c:crosses val="[^"]*"/>`)
	node := `<c:crosses val="` + strings.TrimSpace(*value) + `"/>`

	start := strings.Index(xml, startTag)
	for start >= 0 {
		endRel := strings.Index(xml[start:], endTag)
		if endRel < 0 {
			break
		}
		end := start + endRel + len(endTag)
		block := xml[start:end]
		if reCrosses.MatchString(block) {
			block = reCrosses.ReplaceAllString(block, node)
		} else {
			insertAt := axisNodeInsertIndex(block)
			if insertAt < 0 {
				start = nextAxisStart(xml[end:], startTag, end)
				continue
			}
			block = block[:insertAt] + node + block[insertAt:]
		}
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
