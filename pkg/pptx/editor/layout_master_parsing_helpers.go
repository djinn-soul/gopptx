package editor

import (
	"regexp"
	"strconv"
)

// Attribute and element patterns used when parsing master/layout XML. They are
// fixed literals, so compiling them once at init beats any runtime cache.
var (
	phIdxAttrPattern  = regexp.MustCompile(`idx="([^"]+)"`)
	phTypeAttrPattern = regexp.MustCompile(`type="([^"]+)"`)
	phNameAttrPattern = regexp.MustCompile(`name="([^"]+)"`)
)

func parsePlaceholderAttrIndex(match string) int {
	raw := parsePlaceholderAttrString(match, phIdxAttrPattern)
	if raw == "" {
		return 0
	}
	idx, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return idx
}

func parsePlaceholderAttrString(match string, pattern *regexp.Regexp) string {
	ms := pattern.FindStringSubmatch(match)
	if len(ms) <= 1 {
		return ""
	}
	return ms[1]
}
