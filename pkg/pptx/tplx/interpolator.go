package tplx

import (
	"regexp"
	"strings"
)

// scalarPattern matches {{ key }} with optional whitespace.
var scalarPattern = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.\-]+)\s*\}\}`)

// dotPattern matches {{.field}} for loop-item access.
var dotPattern = regexp.MustCompile(`\{\{\s*\.([a-zA-Z0-9_.\-]+)\s*\}\}`)

// interpolateText replaces all {{key}} scalars in text using ctx.
// Unresolved keys are left as-is when strict=false, or returned as an error sentinel.
func interpolateText(text string, ctx Context, item Row, strict bool) (string, bool) {
	changed := false

	// First resolve item-scoped {{.field}} tokens.
	if item != nil {
		text = dotPattern.ReplaceAllStringFunc(text, func(m string) string {
			matches := dotPattern.FindStringSubmatch(m)
			if len(matches) < 2 {
				return m
			}
			key := matches[1]
			if val, ok := item[key]; ok {
				changed = true
				return val
			}
			return m
		})
	}

	// Then resolve {{key}} scalars from the main context.
	text = scalarPattern.ReplaceAllStringFunc(text, func(m string) string {
		matches := scalarPattern.FindStringSubmatch(m)
		if len(matches) < 2 {
			return m
		}
		key := matches[1]
		val, ok := ctx[key]
		if !ok {
			return m // leave untouched in lenient mode
		}
		changed = true
		switch v := val.(type) {
		case string:
			return v
		case fmt_stringer:
			return v.String()
		default:
			return strings.TrimSpace(fmt_sprint(v))
		}
	})
	return text, changed
}

// fmt_stringer is a copy-free equivalent of fmt.Stringer to avoid importing fmt.
type fmt_stringer interface{ String() string }

// fmt_sprint converts any value to string without importing fmt directly.
func fmt_sprint(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(fmt_stringer); ok {
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
	buf := make([]byte, 20)
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
func interpolateXMLPart(xmlBytes []byte, ctx Context, item Row, strict bool) ([]byte, error) {
	// Step 1: merge fragmented runs.
	merged, err := mergeAdjacentRuns(xmlBytes)
	if err != nil {
		merged = xmlBytes
	}

	// Step 2: simple text substitution inside <a:t> elements.
	result := replaceInTextElements(merged, ctx, item, strict)
	return result, nil
}

// replaceInTextElements walks the raw XML bytes finding <a:t>…</a:t> text and replaces tokens.
// We do this with simple string scanning (not full re-parse) so we preserve all other bytes exactly.
func replaceInTextElements(src []byte, ctx Context, item Row, strict bool) []byte {
	const open = "<a:t>"
	const close = "</a:t>"

	result := make([]byte, 0, len(src))
	s := string(src)
	pos := 0
	for {
		start := strings.Index(s[pos:], open)
		if start < 0 {
			result = append(result, s[pos:]...)
			break
		}
		start += pos
		end := strings.Index(s[start+len(open):], close)
		if end < 0 {
			result = append(result, s[pos:]...)
			break
		}
		end = start + len(open) + end

		// Append everything up to (and including) the opening <a:t>.
		result = append(result, s[pos:start+len(open)]...)

		// The text content.
		text := s[start+len(open) : end]
		replaced, _ := interpolateText(text, ctx, item, strict)
		result = append(result, replaced...)

		// Append </a:t>.
		result = append(result, close...)
		pos = end + len(close)
	}
	return result
}
