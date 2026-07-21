package editor

import (
	"regexp"
	"strconv"
	"sync"
)

// commonPatternCache memoizes compiled attribute patterns. The patterns come
// from a small fixed set of literals, so the cache is bounded; compiling on
// every call cost ~13x the cached lookup.
//
//nolint:gochecknoglobals // Shared compile cache; see above.
var commonPatternCache sync.Map // pattern string -> *regexp.Regexp

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
	if cached, ok := commonPatternCache.Load(pattern); ok {
		re, _ := cached.(*regexp.Regexp)
		return re
	}
	re := regexp.MustCompile(pattern)
	commonPatternCache.Store(pattern, re)
	return re
}
