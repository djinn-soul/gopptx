package editor

import (
	"bytes"
	"encoding/xml"
	"html"
	"regexp"
	"sort"
	"strings"
)

// paragraphPattern bounds one <a:p> block. DrawingML paragraphs do not nest, so
// a non-greedy match is exact.
var paragraphPattern = regexp.MustCompile(`(?s)<a:p(?:\s[^>]*)?>.*?</a:p>`)

// replaceTextRuns replaces findText across the text runs of each paragraph.
//
// PowerPoint splits a typed phrase across several <a:r> runs whenever the
// formatting changes mid-word, and after autocorrect or a spellcheck edit even
// when it does not. Matching per <a:t> would therefore miss text the caller
// sees as contiguous. Matching is scoped to one <a:p> so a phrase never matches
// across a paragraph break.
func replaceTextRuns(content []byte, findText, replaceText string) ([]byte, int) {
	spans := paragraphPattern.FindAllIndex(content, -1)
	if len(spans) == 0 {
		return content, 0
	}

	var buf bytes.Buffer
	buf.Grow(len(content))
	total, pos := 0, 0
	for _, span := range spans {
		buf.Write(content[pos:span[0]])
		updated, count := replaceInParagraph(content[span[0]:span[1]], findText, replaceText)
		buf.Write(updated)
		total += count
		pos = span[1]
	}
	buf.Write(content[pos:])

	if total == 0 {
		return content, 0
	}
	return buf.Bytes(), total
}

// replaceInParagraph joins the paragraph's run texts, replaces in the joined
// string, then writes each run's share back.
func replaceInParagraph(paragraph []byte, findText, replaceText string) ([]byte, int) {
	runs := textRunPattern.FindAllSubmatchIndex(paragraph, -1)
	if len(runs) == 0 {
		return paragraph, 0
	}

	texts := make([]string, len(runs))
	starts := make([]int, len(runs))
	var joined strings.Builder
	for i, run := range runs {
		texts[i] = html.UnescapeString(string(paragraph[run[4]:run[5]]))
		starts[i] = joined.Len()
		joined.WriteString(texts[i])
	}
	if !strings.Contains(joined.String(), findText) {
		return paragraph, 0
	}

	updated, count := redistributeReplacements(joined.String(), starts, findText, replaceText)
	return rewriteRunTexts(paragraph, runs, updated), count
}

// redistributeReplacements rewrites joined text run by run. Unmatched
// characters stay in the run they came from, and each replacement lands whole
// in the run where its match started, so surrounding formatting is preserved
// and no run keeps a fragment of the old text.
func redistributeReplacements(
	joined string,
	starts []int,
	findText string,
	replaceText string,
) ([]string, int) {
	out := make([]strings.Builder, len(starts))
	runAt := func(offset int) int {
		index := sort.SearchInts(starts, offset+1) - 1
		if index < 0 {
			return 0
		}
		return index
	}
	copyRange := func(from, to int) {
		for from < to {
			run := runAt(from)
			end := to
			if run+1 < len(starts) && starts[run+1] < end {
				end = starts[run+1]
			}
			// strings.Builder.WriteString never returns an error.
			_, _ = out[run].WriteString(joined[from:end])
			from = end
		}
	}

	count, pos := 0, 0
	for {
		offset := strings.Index(joined[pos:], findText)
		if offset < 0 {
			break
		}
		start := pos + offset
		copyRange(pos, start)
		_, _ = out[runAt(start)].WriteString(replaceText)
		count++
		pos = start + len(findText)
	}
	copyRange(pos, len(joined))

	texts := make([]string, len(out))
	for i := range out {
		texts[i] = out[i].String()
	}
	return texts, count
}

// rewriteRunTexts writes each run's new text back into the paragraph XML,
// keeping the original open and close tags so run properties survive.
func rewriteRunTexts(paragraph []byte, runs [][]int, texts []string) []byte {
	var buf bytes.Buffer
	buf.Grow(len(paragraph))
	pos := 0
	for i, run := range runs {
		buf.Write(paragraph[pos:run[0]])
		buf.Write(paragraph[run[2]:run[3]]) // open tag, with its attributes
		_ = xml.EscapeText(&buf, []byte(texts[i]))
		buf.Write(paragraph[run[6]:run[7]]) // close tag
		pos = run[1]
	}
	buf.Write(paragraph[pos:])
	return buf.Bytes()
}
