package tplx

import (
	"regexp"
	"strings"
)

const (
	minSubmatchLen      = 2
	intToStrBufferSize  = 20
	numericBaseTen      = 10
	textElementOpenTag  = "<a:t>"
	textElementCloseTag = "</a:t>"
)

// scalarPattern matches {{ key }} with optional whitespace.
var scalarPattern = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.\-]+)\s*\}\}`)

// dotPattern matches {{.field}} for loop-item access.
var dotPattern = regexp.MustCompile(`\{\{\s*\.([a-zA-Z0-9_.\-]+)\s*\}\}`)

// interpolateText replaces all {{key}} scalars in text using ctx.
// Unresolved keys are left as-is. Values are XML-escaped for safe insertion into <a:t> elements.
func interpolateText(text string, ctx Context, item Row) string {
	escape := xmlEscape

	// First resolve item-scoped {{.field}} tokens.
	if item != nil {
		text = dotPattern.ReplaceAllStringFunc(text, func(m string) string {
			matches := dotPattern.FindStringSubmatch(m)
			if len(matches) < minSubmatchLen {
				return m
			}
			key := matches[1]
			if val, ok := item[key]; ok {
				return escape(val)
			}
			return m
		})
	}

	// Then resolve {{key}} scalars from the main context.
	text = scalarPattern.ReplaceAllStringFunc(text, func(m string) string {
		matches := scalarPattern.FindStringSubmatch(m)
		if len(matches) < minSubmatchLen {
			return m
		}
		key := matches[1]
		val, ok := ctx[key]
		if !ok {
			return m // leave untouched in lenient mode
		}
		switch v := val.(type) {
		case string:
			return escape(v)
		case fmtStringer:
			return escape(v.String())
		default:
			return escape(strings.TrimSpace(fmtSprint(v)))
		}
	})
	return text
}

func xmlEscape(value string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	).Replace(value)
}

// fmtStringer is a copy-free equivalent of [fmt.Stringer] to avoid importing fmt.
type fmtStringer interface{ String() string }

// fmtSprint converts any value to string without importing fmt directly.
func fmtSprint(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(fmtStringer); ok {
		return s.String()
	}
	// Fall back to type assertion for common primitives.
	switch vt := v.(type) {
	case bool:
		if vt {
			return "true"
		}
		return "false"
	case int:
		return intToStr(int64(vt))
	case int64:
		return intToStr(vt)
	case float64:
		return floatToStr(vt)
	}
	return ""
}

func intToStr(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := make([]byte, intToStrBufferSize)
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

func floatToStr(f float64) string {
	// Simple integer check.
	if f == float64(int64(f)) {
		return intToStr(int64(f))
	}
	return strconvFmtFloat(f)
}

func strconvFmtFloat(f float64) string {
	// Import-free basic float format.
	// Use strconv indirectly via a helper to avoid adding import at top.
	return strconvAppendFloat(f)
}

// interpolateXMLPart replaces {{key}} tokens in all text runs within an XML part.
// It first runs the run-merger pass, then does string substitution.
func interpolateXMLPart(xmlBytes []byte, ctx Context, item Row, strict bool) []byte {
	// Step 1: merge fragmented runs.
	merged := mergeAdjacentRuns(xmlBytes)
	// Step 2: simple text substitution inside <a:t> elements.
	result := replaceInTextElements(merged, ctx, item, strict)
	return result
}

// replaceInTextElements walks the raw XML bytes finding <a:t>…</a:t> text and replaces tokens.
// We do this with simple string scanning (not full re-parse) so we preserve all other bytes exactly.
func replaceInTextElements(src []byte, ctx Context, item Row, _ bool) []byte {
	result := make([]byte, 0, len(src))
	s := string(src)
	pos := 0
	for {
		start := strings.Index(s[pos:], textElementOpenTag)
		if start < 0 {
			result = append(result, s[pos:]...)
			break
		}
		start += pos
		end := strings.Index(s[start+len(textElementOpenTag):], textElementCloseTag)
		if end < 0 {
			result = append(result, s[pos:]...)
			break
		}
		end = start + len(textElementOpenTag) + end

		// Append everything up to (and including) the opening <a:t>.
		result = append(result, s[pos:start+len(textElementOpenTag)]...)

		// The text content.
		text := s[start+len(textElementOpenTag) : end]
		replaced := interpolateText(text, ctx, item)
		result = append(result, replaced...)

		// Append </a:t>.
		result = append(result, textElementCloseTag...)
		pos = end + len(textElementCloseTag)
	}
	return result
}
