package compress

import (
	"regexp"
	"strings"
)

// Compiled once, immutable.
var (
	relationshipPattern   = regexp.MustCompile(`(?s)<Relationship\b[^>]*/>`)
	overridePattern       = regexp.MustCompile(`(?s)<Override\b[^>]*/>`)
	targetAttrPattern     = regexp.MustCompile(`Target="([^"]*)"`)
	partNameAttrPattern   = regexp.MustCompile(`PartName="([^"]*)"`)
	interElementWhitespce = regexp.MustCompile(`>[ \t\r\n]*\n[ \t\r\n]*<`)
)

// relationshipTargets returns every Target attribute in a `.rels` part.
func relationshipTargets(data []byte) []string {
	matches := targetAttrPattern.FindAllSubmatch(data, -1)
	targets := make([]string, 0, len(matches))
	for _, m := range matches {
		targets = append(targets, string(m[1]))
	}
	return targets
}

// dropRelationships removes Relationship entries whose resolved target is in
// removed. Relationship IDs of surviving entries are left untouched, because
// slide XML references them by ID.
func dropRelationships(data []byte, owner string, removed map[string]bool) []byte {
	return relationshipPattern.ReplaceAllFunc(data, func(entry []byte) []byte {
		m := targetAttrPattern.FindSubmatch(entry)
		if m == nil {
			return entry
		}
		target := string(m[1])
		if strings.Contains(target, "://") {
			return entry
		}
		if removed[resolveTarget(owner, target)] {
			return nil
		}
		return entry
	})
}

// dropContentTypeOverrides removes Override entries for parts that are gone.
func dropContentTypeOverrides(data []byte, removed map[string]bool) []byte {
	return overridePattern.ReplaceAllFunc(data, func(entry []byte) []byte {
		m := partNameAttrPattern.FindSubmatch(entry)
		if m == nil {
			return entry
		}
		if removed[normalizeName(string(m[1]))] {
			return nil
		}
		return entry
	})
}

// minifyXML strips the indentation between elements. Only whitespace runs that
// contain a newline are removed, so single spaces inside `<a:t>` text runs
// survive; a text run that itself spans lines keeps its content because the
// replacement leaves the element delimiters in place.
func minifyXML(data []byte) []byte {
	return interElementWhitespce.ReplaceAll(data, []byte("><"))
}

func isXMLPart(name string) bool {
	n := strings.ToLower(normalizeName(name))
	return strings.HasSuffix(n, ".xml") || strings.HasSuffix(n, ".rels")
}

func isImagePart(name string) bool {
	n := strings.ToLower(normalizeName(name))
	if !strings.HasPrefix(n, "ppt/media/") {
		return false
	}
	return strings.HasSuffix(n, ".jpg") ||
		strings.HasSuffix(n, ".jpeg") ||
		strings.HasSuffix(n, ".png")
}
