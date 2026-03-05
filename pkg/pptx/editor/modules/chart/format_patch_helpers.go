package chart

import "bytes"

func boolToOneZero(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func xmlEscape(value string) string {
	var b bytes.Buffer
	b.Grow(len(value))
	for _, r := range value {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&apos;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
