package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// GanttTask represents a task in a Gantt chart.
type GanttTask struct {
	Name     string
	ID       string
	Status   string
	Start    string
	Duration string
}

// GanttSection represents a section in a Gantt chart.
type GanttSection struct {
	Name  string
	Tasks []GanttTask
}

// GanttDiagram represents the parsed structure of a Mermaid Gantt chart.
type GanttDiagram struct {
	Title    string
	Sections []GanttSection
}

// renderGantt parses and renders a Mermaid Gantt chart into PowerPoint elements.
func renderGantt(code string, theme Theme) DiagramElements {
	gantt := parseGantt(code)
	return generateGanttElements(gantt, theme)
}

func parseGantt(code string) *GanttDiagram {
	lines := ParseLines(code)
	gantt := &GanttDiagram{}
	var currentSection *GanttSection

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		if strings.HasPrefix(lower, "gantt") {
			continue
		}

		if strings.HasPrefix(lower, "title") {
			gantt.Title = strings.TrimSpace(trimmed[5:])
			continue
		}

		if strings.HasPrefix(lower, "section") {
			if currentSection != nil {
				gantt.Sections = append(gantt.Sections, *currentSection)
			}
			currentSection = &GanttSection{
				Name:  strings.TrimSpace(trimmed[7:]),
				Tasks: []GanttTask{},
			}
			continue
		}

		if strings.Contains(trimmed, ":") {
			parts := strings.Split(trimmed, ":")
			taskName := strings.TrimSpace(parts[0])
			taskDetails := strings.TrimSpace(parts[1])
			
			details := strings.Split(taskDetails, ",")
			task := GanttTask{Name: taskName}
			
			// Very basic parsing of task details
			// Format: [status,] [id,] [start,] duration
			for i, detail := range details {
				detail = strings.TrimSpace(detail)
				if i == len(details)-1 {
					task.Duration = detail
				} else if i == 0 {
					// Could be status, id, or start
					task.ID = detail
				} else if i == 1 {
					task.Start = detail
				}
			}

			if currentSection == nil {
				currentSection = &GanttSection{Name: "Default", Tasks: []GanttTask{}}
			}
			currentSection.Tasks = append(currentSection.Tasks, task)
		}
	}

	if currentSection != nil {
		gantt.Sections = append(gantt.Sections, *currentSection)
	}

	return gantt
}

func generateGanttElements(gantt *GanttDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape
	
	if len(gantt.Sections) == 0 {
		return createPlaceholder("gantt (no data)", theme)
	}

	// Layout parameters
	startX := styling.Inches(1)
	startY := styling.Inches(1.5)
	labelWidth := styling.Inches(2.5)
	chartWidth := styling.Inches(6)
	rowHeight := styling.Inches(0.5)
	sectionHeight := styling.Inches(0.6)
	
	currentY := startY

	// Title
	if gantt.Title != "" {
		titleShape := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			startX,
			startY-styling.Inches(0.8),
			labelWidth+chartWidth,
			styling.Inches(0.6),
		).WithText(gantt.Title).
			WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, titleShape)
	}

	for _, section := range gantt.Sections {
		// Section header
		sectionShape := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			startX,
			currentY,
			labelWidth+chartWidth,
			sectionHeight,
		).WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
			WithText(section.Name).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, sectionShape)
		currentY += sectionHeight

		for _, task := range section.Tasks {
			// Task label
			labelShape := shapes.NewShape(
				shapes.ShapeTypeRectangle,
				startX,
				currentY,
				labelWidth,
				rowHeight,
			).WithText(task.Name).
				WithFill(shapes.NewShapeFill(theme.Background)).
				WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
				WithAutoFit(shapes.TextAutoFitNormal).
				WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
			shapesList = append(shapesList, labelShape)

			// Task bar (simplified layout: all tasks same width for now as we don't parse dates fully)
			barWidth := styling.Inches(2)
			barX := startX + labelWidth + styling.Inches(0.5)
			
			barShape := shapes.NewShape(
				shapes.ShapeTypeRoundedRectangle,
				barX,
				currentY+styling.Inches(0.05),
				barWidth,
				rowHeight-styling.Inches(0.1),
			).WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
				WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
				WithAutoFit(shapes.TextAutoFitNormal)
			shapesList = append(shapesList, barShape)

			currentY += rowHeight
		}
		currentY += styling.Inches(0.2) // Gap between sections
	}

	return DiagramElements{
		Shapes:  shapesList,
		Grouped: true,
		Bounds: &DiagramBounds{
			X:  startX,
			Y:  startY - styling.Inches(0.8),
			CX: labelWidth + chartWidth,
			CY: currentY - (startY - styling.Inches(0.8)),
		},
	}
}
