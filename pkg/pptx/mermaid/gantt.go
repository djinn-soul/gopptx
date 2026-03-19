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
		currentSection = consumeGanttLine(gantt, currentSection, strings.TrimSpace(line))
	}

	if currentSection != nil {
		gantt.Sections = append(gantt.Sections, *currentSection)
	}

	return gantt
}

func consumeGanttLine(gantt *GanttDiagram, currentSection *GanttSection, trimmed string) *GanttSection {
	if trimmed == "" {
		return currentSection
	}
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "gantt") {
		return currentSection
	}
	if title, ok := parseGanttTitle(trimmed, lower); ok {
		gantt.Title = title
		return currentSection
	}
	if sectionName, ok := parseGanttSectionName(trimmed, lower); ok {
		if currentSection != nil {
			gantt.Sections = append(gantt.Sections, *currentSection)
		}
		return &GanttSection{Name: sectionName, Tasks: []GanttTask{}}
	}
	task, ok := parseGanttTask(trimmed)
	if !ok {
		return currentSection
	}
	currentSection = ensureGanttSection(currentSection)
	currentSection.Tasks = append(currentSection.Tasks, task)
	return currentSection
}

func parseGanttTitle(trimmed string, lower string) (string, bool) {
	if !strings.HasPrefix(lower, "title") {
		return "", false
	}
	return strings.TrimSpace(trimmed[5:]), true
}

func parseGanttSectionName(trimmed string, lower string) (string, bool) {
	if !strings.HasPrefix(lower, "section") {
		return "", false
	}
	return strings.TrimSpace(trimmed[7:]), true
}

func parseGanttTask(trimmed string) (GanttTask, bool) {
	taskName, taskDetails, ok := strings.Cut(trimmed, ":")
	if !ok {
		return GanttTask{}, false
	}
	task := GanttTask{Name: strings.TrimSpace(taskName)}
	details := strings.Split(taskDetails, ",")
	for i, detail := range details {
		detail = strings.TrimSpace(detail)
		switch {
		case i == len(details)-1:
			task.Duration = detail
		case i == 0:
			task.ID = detail
		case i == 1:
			task.Start = detail
		}
	}
	return task, true
}

func ensureGanttSection(currentSection *GanttSection) *GanttSection {
	if currentSection != nil {
		return currentSection
	}
	return &GanttSection{Name: "Default", Tasks: []GanttTask{}}
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
	axisHeight := styling.Inches(0.35)

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

	// Timeline axis and coarse tick labels so bars have temporal context.
	axisY := currentY
	axisLine := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		startX+labelWidth,
		axisY+axisHeight/2,
		chartWidth,
		styling.Emu(19050),
	).WithFill(shapes.NewShapeFill(theme.PrimaryStroke))
	shapesList = append(shapesList, axisLine)
	taskStarts := collectGanttStartLabels(gantt)
	if len(taskStarts) == 0 {
		taskStarts = []string{"T1", "T2", "T3"}
	}
	for i, tick := range taskStarts {
		tickX := startX + labelWidth + (styling.Length(i) * chartWidth / styling.Length(max(1, len(taskStarts)-1)))
		tickShape := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			tickX-styling.Emu(9525),
			axisY+axisHeight/2-styling.Inches(0.06),
			styling.Emu(19050),
			styling.Inches(0.12),
		).WithFill(shapes.NewShapeFill(theme.PrimaryStroke))
		tickLabel := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			tickX-styling.Inches(0.42),
			axisY,
			styling.Inches(0.84),
			styling.Inches(0.22),
		).WithText(tick).
			WithFill(shapes.NewShapeFill(theme.Background)).
			WithLine(shapes.NewShapeLine(theme.Background, styling.Emu(0))).
			WithAutoFit(shapes.TextAutoFitNormal)
		shapesList = append(shapesList, tickShape, tickLabel)
	}
	currentY += axisHeight + styling.Inches(0.1)

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
				WithVerticalAnchor(shapes.TextAnchorMiddle).
				WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
			shapesList = append(shapesList, labelShape)

			// Task bar with duration text shown in-bar.
			barWidth := styling.Inches(2.2)
			barX := startX + labelWidth + styling.Inches(0.5)
			barText := task.Duration
			if barText == "" {
				barText = task.Start
			}

			barShape := shapes.NewShape(
				shapes.ShapeTypeRoundedRectangle,
				barX,
				currentY+styling.Inches(0.05),
				barWidth,
				rowHeight-styling.Inches(0.1),
			).WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
				WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
				WithText(barText).
				WithVerticalAnchor(shapes.TextAnchorMiddle).
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

func collectGanttStartLabels(gantt *GanttDiagram) []string {
	seen := make(map[string]struct{})
	labels := make([]string, 0, 6)
	for _, section := range gantt.Sections {
		for _, task := range section.Tasks {
			if task.Start == "" {
				continue
			}
			if _, ok := seen[task.Start]; ok {
				continue
			}
			seen[task.Start] = struct{}{}
			labels = append(labels, task.Start)
			if len(labels) >= 6 {
				return labels
			}
		}
	}
	return labels
}
