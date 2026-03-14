package editor

import (
	"regexp"
	"strconv"
)

func parsePlaceholderAttrIndex(match string) int {
	raw := parsePlaceholderAttrString(match, `idx="([^"]+)"`)
	if raw == "" {
		return 0
	}
	idx, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return idx
}

func parsePlaceholderAttrString(match, pattern string) string {
	re := commonCompile(pattern)
	if re == nil {
		return ""
	}
	ms := re.FindStringSubmatch(match)
	if len(ms) <= 1 {
		return ""
	}
	return ms[1]
}

func commonCompile(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}
