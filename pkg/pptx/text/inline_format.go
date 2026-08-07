package text

import "strings"

// Inline markdown formatting — **bold**, *italic*, `code` — used to be
// understood only by the Markdown importer, so a caller building slides
// directly got the asterisks and backticks as literal text and had to hand-build
// runs to get formatting. The parser is the model's business, so it lives with
// the model.

const (
	boldMarkerLength    = 2
	inlineRunsCapacity  = 4
	italicMarkerLength  = 1
	codeMarkerLength    = 1
	inlineMarkerLiteral = "`*"
)

// ParseInlineRuns splits text on the inline markdown markers, returning the
// runs and whether any formatting was found. Unmatched markers stay literal.
func ParseInlineRuns(input string) ([]Run, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil, false
	}

	runs := make([]Run, 0, inlineRunsCapacity)
	styled := false
	for i := 0; i < len(trimmed); {
		if run, next, ok := parseStyledRun(trimmed, i); ok {
			runs = append(runs, run)
			styled = true
			i = next
			continue
		}

		next := nextInlineMarkerOffset(trimmed[i:])
		switch {
		case next < 0:
			runs = append(runs, Run{Text: trimmed[i:]})
			i = len(trimmed)
		case next == 0:
			// An unmatched marker: keep it as text and move past it.
			runs = append(runs, Run{Text: trimmed[i : i+1]})
			i++
		default:
			runs = append(runs, Run{Text: trimmed[i : i+next]})
			i += next
		}
	}

	return NormalizeRuns(runs), styled
}

// HasInlineFormatting reports whether the text carries markers ParseInlineRuns
// would act on.
func HasInlineFormatting(input string) bool {
	_, styled := ParseInlineRuns(input)
	return styled
}

func parseStyledRun(input string, i int) (Run, int, bool) {
	switch {
	case input[i] == '`':
		return parseDelimitedRun(input, i, "`", codeMarkerLength, func(r *Run) { r.Code = true })
	case strings.HasPrefix(input[i:], "**"):
		return parseDelimitedRun(input, i, "**", boldMarkerLength, func(r *Run) { r.Bold = true })
	case input[i] == '*':
		return parseDelimitedRun(input, i, "*", italicMarkerLength, func(r *Run) { r.Italic = true })
	default:
		return Run{}, i, false
	}
}

// parseDelimitedRun reads one span between matching markers. An empty span or a
// missing closing marker leaves the text alone.
func parseDelimitedRun(input string, i int, marker string, markerLen int, style func(*Run)) (Run, int, bool) {
	closeIdx := strings.Index(input[i+markerLen:], marker)
	if closeIdx < 0 {
		return Run{}, i, false
	}
	end := i + markerLen + closeIdx
	if end <= i+markerLen {
		return Run{}, i, false
	}
	run := Run{Text: input[i+markerLen : end]}
	style(&run)
	return run, end + markerLen, true
}

func nextInlineMarkerOffset(input string) int {
	return strings.IndexAny(input, inlineMarkerLiteral)
}
