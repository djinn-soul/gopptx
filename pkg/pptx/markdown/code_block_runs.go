package markdown

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

// The highlighter has always been here because the Markdown importer was the
// only thing that could reach code. A CodeBlock element placed at explicit
// coordinates needs the same tokens, so the run-building step is exported.

// HighlightedLines tokenises source and returns one slice of styled runs per
// line, ready to become paragraphs in a shape's text body. An empty or
// unrecognised language returns unhighlighted runs rather than failing.
func HighlightedLines(language, code string, fontSizePt int, font string) [][]elements.Run {
	lang := strings.ToLower(strings.TrimSpace(language))
	if lang == "" {
		lang = LanguagePlainText
	}
	tokens := newSyntaxHighlighter(lang).Tokenize(code)

	lines := make([][]elements.Run, 0, strings.Count(code, "\n")+1)
	current := make([]elements.Run, 0, defaultInlineRunsCapacity)
	flush := func() {
		lines = append(lines, current)
		current = make([]elements.Run, 0, defaultInlineRunsCapacity)
	}

	for _, token := range tokens {
		// Chroma reports newlines inside token text, so a token can close one
		// line and open the next.
		parts := strings.Split(token.Text, "\n")
		for i, part := range parts {
			if i > 0 {
				flush()
			}
			if part == "" {
				continue
			}
			current = append(current, codeRun(part, token, fontSizePt, font))
		}
	}
	if len(current) > 0 {
		flush()
	}
	return lines
}

// LanguageLabel is the caption a code block prints above its listing.
func LanguageLabel(language string) string {
	lang := strings.ToUpper(strings.TrimSpace(language))
	if lang == "" {
		lang = "TEXT"
	}
	return "[" + lang + "]"
}

func codeRun(part string, token Token, fontSizePt int, font string) elements.Run {
	color := token.Color
	if color == "" {
		color = GetColor(token.Type)
	}
	run := text.NewRun(part).
		WithCode(true).
		WithColor(color)
	if fontSizePt > 0 {
		run = run.WithSizePt(fontSizePt)
	}
	if font != "" {
		run = run.WithFont(font)
	}
	return run
}
