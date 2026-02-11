package markdown

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

const (
	defaultCodeForeground = "D4D4D4"
	codeKeywordColor      = "005A9E"
	codeCommentColor      = "2E7D32"
	codeStringColor       = "A31515"
)

func addCodeBlock(slide *elements.SlideContent, lang string, code string) {
	if slide == nil {
		return
	}
	style := elements.DefaultTextParagraphStyle().WithNoBullet()
	normalizedLang := strings.ToLower(strings.TrimSpace(lang))
	if normalizedLang == "" {
		normalizedLang = "text"
	}
	header := []elements.TextRun{
		elements.NewTextRun("[" + strings.ToUpper(normalizedLang) + "]").WithBold(true).WithColor("1F4E78").WithSizePt(13),
	}
	*slide = slide.AddBulletRunsWithStyle(header, style)

	for _, line := range strings.Split(code, "\n") {
		if strings.TrimSpace(line) == "" {
			*slide = slide.AddBulletWithStyle(" ", style)
			continue
		}
		run := elements.NewTextRun(line).
			WithCode(true).
			WithColor(codeLineColor(normalizedLang, line)).
			WithSizePt(14)
		*slide = slide.AddBulletRunsWithStyle([]elements.TextRun{run}, style)
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
	return addMermaidDiagram(slide, code, lineNumber)
}

type mermaidDiagramStyle struct {
	kind      string
	subtitle  string
	shapeType string
	fillColor string
	lineColor string
}

func addMermaidDiagram(slide *elements.SlideContent, code string, lineNumber int) error {
	if slide == nil {
		return fmt.Errorf("line %d: mermaid block requires an active slide", lineNumber)
	}

	style, err := detectMermaidDiagram(code, lineNumber)
	if err != nil {
		return err
	}

	text := "Mermaid Diagram: " + style.kind
	if style.subtitle != "" {
		text += " (" + style.subtitle + ")"
	}
	if summary := mermaidSummaryLine(code); summary != "" {
		text += "\n" + summary
	}

	shape := shapes.NewShape(style.shapeType, 762000, 1524000, 8001000, 1714500).
		WithFill(shapes.NewShapeFill(style.fillColor).WithTransparency(8)).
		WithLine(shapes.NewShapeLine(style.lineColor, 15875)).
		WithText(text)
	*slide = slide.AddShape(shape)
	return nil
}

func detectMermaidDiagram(code string, lineNumber int) (mermaidDiagramStyle, error) {
	first := firstNonEmptyLine(code)
	if first == "" {
		return mermaidDiagramStyle{}, fmt.Errorf("line %d: mermaid block is empty", lineNumber)
	}

	fields := strings.Fields(first)
	if len(fields) == 0 {
		return mermaidDiagramStyle{}, fmt.Errorf("line %d: mermaid block is empty", lineNumber)
	}

	directive := strings.ToLower(strings.TrimSpace(fields[0]))
	switch directive {
	case "flowchart", "graph":
		subtitle := "Flow"
		if len(fields) > 1 {
			subtitle = strings.ToUpper(fields[1])
		}
		return mermaidDiagramStyle{
			kind:      "Flowchart",
			subtitle:  subtitle,
			shapeType: shapes.ShapeTypeFlowChartProcess,
			fillColor: "DCE6F2",
			lineColor: "2F5597",
		}, nil
	case "sequencediagram":
		return mermaidDiagramStyle{
			kind:      "Sequence Diagram",
			shapeType: shapes.ShapeTypeRectangle,
			fillColor: "E2F0D9",
			lineColor: "2E7D32",
		}, nil
	case "classdiagram":
		return mermaidDiagramStyle{
			kind:      "Class Diagram",
			shapeType: shapes.ShapeTypeRectangle,
			fillColor: "FCE4D6",
			lineColor: "A64D00",
		}, nil
	case "statediagram", "statediagram-v2":
		return mermaidDiagramStyle{
			kind:      "State Diagram",
			shapeType: shapes.ShapeTypeRoundedRectangle,
			fillColor: "EDE2F7",
			lineColor: "6A1B9A",
		}, nil
	case "erdiagram":
		return mermaidDiagramStyle{
			kind:      "Entity-Relationship Diagram",
			shapeType: shapes.ShapeTypeRectangle,
			fillColor: "E8F5E9",
			lineColor: "1B5E20",
		}, nil
	case "journey":
		return mermaidDiagramStyle{
			kind:      "User Journey",
			shapeType: shapes.ShapeTypeRoundedRectangle,
			fillColor: "FFF2CC",
			lineColor: "8A6D1A",
		}, nil
	case "gantt":
		return mermaidDiagramStyle{
			kind:      "Gantt Chart",
			shapeType: shapes.ShapeTypeRectangle,
			fillColor: "E2EFDA",
			lineColor: "2F6B2F",
		}, nil
	case "pie":
		return mermaidDiagramStyle{
			kind:      "Pie Chart",
			shapeType: shapes.ShapeTypeEllipse,
			fillColor: "FBE5D6",
			lineColor: "C65911",
		}, nil
	case "mindmap":
		return mermaidDiagramStyle{
			kind:      "Mindmap",
			shapeType: shapes.ShapeTypeEllipse,
			fillColor: "E4DFEC",
			lineColor: "5B4B8A",
		}, nil
	case "quadrantchart":
		return mermaidDiagramStyle{
			kind:      "Quadrant Chart",
			shapeType: shapes.ShapeTypeRectangle,
			fillColor: "D9E1F2",
			lineColor: "203864",
		}, nil
	case "timeline":
		return mermaidDiagramStyle{
			kind:      "Timeline",
			shapeType: shapes.ShapeTypeRightArrow,
			fillColor: "DEEAF6",
			lineColor: "2F75B5",
		}, nil
	case "gitgraph":
		return mermaidDiagramStyle{
			kind:      "Git Graph",
			shapeType: shapes.ShapeTypeParallelogram,
			fillColor: "EDEDED",
			lineColor: "595959",
		}, nil
	default:
		return mermaidDiagramStyle{}, fmt.Errorf(
			"line %d: unsupported mermaid diagram %q",
			lineNumber,
			fields[0],
		)
	}
}

func mermaidSummaryLine(code string) string {
	lines := strings.Split(code, "\n")
	for _, line := range lines[1:] {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "%%") {
			continue
		}
		if len(trimmed) > 70 {
			return trimmed[:67] + "..."
		}
		return trimmed
	}
	return ""
}

func firstNonEmptyLine(text string) string {
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
