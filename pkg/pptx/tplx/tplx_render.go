package tplx

import "fmt"

func renderSlideParts(parts map[string][]byte, slideParts []string, ctx Context, strict bool) error {
	for _, name := range slideParts {
		data, ok := parts[name]
		if !ok {
			continue
		}
		updated, err := renderSlidePart(data, ctx, strict, name)
		if err != nil {
			return err
		}
		parts[name] = updated
	}
	return nil
}

func renderSlidePart(data []byte, ctx Context, strict bool, name string) ([]byte, error) {
	data = applyConditionals(data, ctx)
	data = expandTableRows(data, ctx)
	data = interpolateXMLPart(data, ctx, nil, strict)
	if err := validateStrictTokens(data, strict, name); err != nil {
		return nil, err
	}
	return data, nil
}

func renderNonSlideParts(parts map[string][]byte, ctx Context, strict bool) error {
	for name, data := range parts {
		if !isNonSlideTextPart(name) {
			continue
		}
		data = interpolateXMLPart(data, ctx, nil, strict)
		if err := validateStrictTokens(data, strict, name); err != nil {
			return err
		}
		parts[name] = data
	}
	return nil
}

func validateStrictTokens(data []byte, strict bool, partName string) error {
	if !strict {
		return nil
	}
	if tok := firstTemplateToken(data); tok != "" {
		return fmt.Errorf("tplx: unresolved token %q in %s", tok, partName)
	}
	return nil
}

func firstTemplateToken(data []byte) string {
	match := tokenPattern.Find(data)
	if len(match) == 0 {
		return ""
	}
	return string(match)
}

func protectEscapedTemplateDelimiters(parts map[string][]byte, slideParts []string) {
	for _, name := range slideParts {
		if data, ok := parts[name]; ok {
			parts[name] = protectEscapedTokenText(data)
		}
	}
	for name, data := range parts {
		if isNonSlideTextPart(name) {
			parts[name] = protectEscapedTokenText(data)
		}
	}
}

func restoreEscapedTemplateDelimiters(parts map[string][]byte, slideParts []string) {
	for _, name := range slideParts {
		if data, ok := parts[name]; ok {
			parts[name] = restoreEscapedTokenText(data)
		}
	}
	for name, data := range parts {
		if isNonSlideTextPart(name) {
			parts[name] = restoreEscapedTokenText(data)
		}
	}
}

func protectEscapedTokenText(data []byte) []byte {
	out := bytesReplaceAll(data, []byte(`\{{`), escapedTokenOpenBytes())
	out = bytesReplaceAll(out, []byte(`\}}`), escapedTokenCloseBytes())
	return out
}

func restoreEscapedTokenText(data []byte) []byte {
	out := bytesReplaceAll(data, escapedTokenOpenBytes(), []byte("{{"))
	out = bytesReplaceAll(out, escapedTokenCloseBytes(), []byte("}}"))
	return out
}

func escapedTokenOpenBytes() []byte {
	return []byte("@@TPLX_ESC_OPEN@@")
}

func escapedTokenCloseBytes() []byte {
	return []byte("@@TPLX_ESC_CLOSE@@")
}
