package chart

import "strings"

func boolToOneZero(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func xmlEscape(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(text)
}
