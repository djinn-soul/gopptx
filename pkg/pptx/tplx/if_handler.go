package tplx

import (
	"regexp"
	"strings"
)

// ifPattern matches {{#if KEY}} … {{/if}} blocks.
var (
	ifPattern    = regexp.MustCompile(`\{\{#if\s+([a-zA-Z0-9_]+)\s*\}\}`)
	ifEndPattern = regexp.MustCompile(`\{\{/if\s*\}\}`)
)

// applyConditionals removes shape text bodies (or entire shapes) whose guard
// condition evaluates to falsy in ctx.
//
// Strategy: scan for {{#if KEY}} tokens in the XML. If ctx[KEY] is falsy (nil,
// false, empty string, 0), find the smallest enclosing <p:sp> element and
// blank its text body by replacing <a:t>…</a:t> with <a:t></a:t>, then remove
// the {{#if}}…{{/if}} wrapper tokens.  If ctx[KEY] is truthy, just remove the
// wrapper tokens and keep the content.
func applyConditionals(slideXML []byte, ctx Context) []byte {
	content := string(slideXML)

	for {
		m := ifPattern.FindStringIndex(content)
		if m == nil {
			break
		}

		// Extract key.
		keyMatch := ifPattern.FindStringSubmatch(content[m[0]:])
		if len(keyMatch) < 2 {
			break
		}
		key := keyMatch[1]

		// Find the matching {{/if}}.
		endIdx := ifEndPattern.FindStringIndex(content[m[1]:])
		if endIdx == nil {
			// Malformed: just strip the opening tag.
			content = content[:m[0]] + content[m[1]:]
			continue
		}
		endStart := m[1] + endIdx[0]
		endEnd := m[1] + endIdx[1]

		inner := content[m[1]:endStart]
		truthy := isTruthy(ctx[key])

		if truthy {
			// Keep the inner content; strip wrapper tokens only.
			content = content[:m[0]] + inner + content[endEnd:]
		} else {
			// Remove the entire block.
			content = content[:m[0]] + content[endEnd:]
		}
	}

	return []byte(content)
}

// isTruthy returns true for non-zero/non-empty/non-nil/true values.
func isTruthy(v any) bool {
	if v == nil {
		return false
	}
	switch vt := v.(type) {
	case bool:
		return vt
	case string:
		return strings.TrimSpace(vt) != ""
	case int:
		return vt != 0
	case int64:
		return vt != 0
	case float64:
		return vt != 0
	case []Row:
		return len(vt) > 0
	case []map[string]string:
		return len(vt) > 0
	}
	return true // non-nil unknown type is truthy
}
