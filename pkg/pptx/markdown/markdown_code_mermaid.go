package markdown

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/mermaid"
)

const (
	defaultCodeForeground = "D4D4D4"
	codeKeywordColor      = "005A9E"
	codeCommentColor      = "2E7D32"
	codeStringColor       = "A31515"

	codeHeaderFontSizePt = 13
	codeBodyFontSizePt   = 14
)

func addCodeBlock(slide *elements.SlideContent, lang string, code string) {
	if slide == nil {
		return
	}
	style := elements.DefaultParagraphStyle().WithNoBullet()
	normalizedLang := strings.ToLower(strings.TrimSpace(lang))
	if normalizedLang == "" {
		normalizedLang = "text"
	}
	header := []elements.Run{
		elements.NewRun("[" + strings.ToUpper(normalizedLang) + "]").
			WithBold(true).
			WithColor("1F4E78").
			WithSizePt(codeHeaderFontSizePt),
	}
	*slide = slide.AddBulletRunsWithStyle(header, style)

	for line := range strings.SplitSeq(code, "\n") {
		if strings.TrimSpace(line) == "" {
			*slide = slide.AddBulletWithStyle(" ", style)
			continue
		}
		run := elements.NewRun(line).
			WithCode(true).
			WithColor(codeLineColor(normalizedLang, line)).
			WithSizePt(codeBodyFontSizePt)
		*slide = slide.AddBulletRunsWithStyle([]elements.Run{run}, style)
	}
}

func codeLineColor(lang string, line string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return defaultCodeForeground
	}
	if isCodeComment(lang, trimmed) {
		return codeCommentColor
	}
	if strings.Contains(trimmed, `"`) || strings.Contains(trimmed, `'`) {
		return codeStringColor
	}
	if containsCodeKeyword(lang, trimmed) {
		return codeKeywordColor
	}
	return defaultCodeForeground
}

func isCodeComment(lang string, line string) bool {
	switch lang {
	case "python", "py":
		return strings.HasPrefix(line, "#")
	default:
		return strings.HasPrefix(line, "//")
	}
}

func containsCodeKeyword(lang string, line string) bool {
	keywords := codeKeywordsByLanguage(lang)
	for _, keyword := range keywords {
		if strings.Contains(line, keyword) {
			return true
		}
	}
	return false
}

func codeKeywordsByLanguage(lang string) []string {
	switch lang {
	case "rust", "rs":
		return []string{"fn ", "let ", "impl ", "struct ", "enum ", "pub ", "use "}
	case "python", "py":
		return []string{"def ", "class ", "import ", "return ", "for ", "while ", "if "}
	case "javascript", "js", "typescript", "ts":
		return []string{"function ", "const ", "let ", "class ", "return ", "if "}
	case "go", "golang":
		return []string{"func ", "type ", "struct ", "package ", "import ", "return "}
	default:
		return nil
	}
}

func addMermaidPlaceholder(slide *elements.SlideContent, code string, lineNumber int) error {
	diagram, err := mermaid.CreateDiagram(code)
	if err != nil {
		return fmt.Errorf("line %d: %w", lineNumber, err)
	}

	for _, s := range diagram.Shapes {
		*slide = slide.AddShape(s)
	}
	for _, c := range diagram.Connectors {
		*slide = slide.AddConnector(c)
	}
	return nil
}
