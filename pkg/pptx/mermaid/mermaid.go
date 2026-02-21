package mermaid

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// DetectType identifies the Mermaid diagram type from the provided code.
func DetectType(code string) MermaidType {
	lines := strings.Split(code, "\n")
	var firstLine string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "%%") {
			firstLine = strings.ToLower(trimmed)
			break
		}
	}

	if strings.HasPrefix(firstLine, "graph") || strings.HasPrefix(firstLine, "flowchart") {
		return Flowchart
	} else if strings.HasPrefix(firstLine, "sequencediagram") || strings.HasPrefix(firstLine, "sequence") {
		return Sequence
	} else if strings.HasPrefix(firstLine, "pie") {
		return Pie
	} else if strings.HasPrefix(firstLine, "gantt") {
		return Gantt
	} else if strings.HasPrefix(firstLine, "classdiagram") || strings.HasPrefix(firstLine, "class") {
		return Class
	} else if strings.HasPrefix(firstLine, "statediagram") || strings.HasPrefix(firstLine, "state") {
		return State
	} else if strings.HasPrefix(firstLine, "erdiagram") || strings.HasPrefix(firstLine, "er") {
		return ER
	} else if strings.HasPrefix(firstLine, "mindmap") {
		return Mindmap
	} else if strings.HasPrefix(firstLine, "timeline") {
		return Timeline
	} else if strings.HasPrefix(firstLine, "journey") {
		return Journey
	} else if strings.HasPrefix(firstLine, "quadrantchart") || strings.HasPrefix(firstLine, "quadrant") {
		return Quadrant
	} else if strings.HasPrefix(firstLine, "gitgraph") || strings.HasPrefix(firstLine, "git") {
		return GitGraph
	}

	return Unknown
}

// CreateDiagram generates PowerPoint elements for the given Mermaid code.
// It dispatches to type-specific generators.
func CreateDiagram(code string) (DiagramElements, error) {
	diagramType := DetectType(code)
	themeName := DetectTheme(code)
	theme := GetTheme(themeName)
	var diagram DiagramElements

	switch diagramType {
	case Flowchart:
		diagram = renderFlowchart(code, theme)
	case Sequence:
		diagram = renderSequence(code, theme)
	case Pie:
		diagram = renderPie(code, theme)
	case Gantt:
		diagram = renderGantt(code, theme)
	case Class:
		diagram = renderClass(code, theme)
	case State:
		diagram = renderState(code, theme)
	case ER:
		diagram = renderER(code, theme)
	case Mindmap:
		diagram = renderMindmap(code, theme)
	case Timeline:
		diagram = renderTimeline(code, theme)
	case Journey:
		diagram = renderJourney(code, theme)
	case GitGraph:
		diagram = renderGitGraph(code, theme)
	case Quadrant:
		diagram = renderQuadrant(code, theme)
	case Unknown:
		diagram = createPlaceholder(code, theme)
	default:
		diagram = createPlaceholder(code, theme)
	}
	if diagramType == Journey {
		return diagram, nil
	}
	return fitDiagramToSlide(diagram), nil
}

func createPlaceholder(code string, theme Theme) DiagramElements {
	lines := strings.Split(strings.TrimSpace(code), "\n")
	firstLine := "Unknown"
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "%%") {
			firstLine = trimmed
			break
		}
	}

	// Default placeholder dimensions
	x := styling.Inches(1)
	y := styling.Inches(2)
	cx := styling.Inches(7)
	cy := styling.Inches(3)

	placeholder := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		x, y, cx, cy,
	).WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, styling.Emu(12700))).
		WithText(fmt.Sprintf("Diagram: %s", firstLine))

	return DiagramElements{
		Shapes:  []shapes.Shape{placeholder},
		Grouped: true,
		Bounds: &DiagramBounds{
			X:  x,
			Y:  y,
			CX: cx,
			CY: cy,
		},
	}
}
