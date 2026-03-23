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

	layout := ganttLayout{
		startX:        styling.Inches(1),
		startY:        styling.Inches(1.5),
		labelWidth:    styling.Inches(2.5),
		chartWidth:    styling.Inches(6),
		rowHeight:     styling.Inches(0.5),
		sectionHeight: styling.Inches(0.6),
		axisHeight:    styling.Inches(0.35),
	}
	axisY := layout.startY
	currentY := axisY + layout.axisHeight + styling.Inches(0.1)
	if gantt.Title != "" {
		shapesList = append(shapesList, buildGanttTitleShape(gantt.Title, theme, layout))
	}

	shapesList = append(shapesList, buildGanttAxisShapes(gantt, theme, axisY, layout)...)

	for _, section := range gantt.Sections {
		sectionShapes, nextY := buildGanttSectionShapes(section, theme, currentY, layout)
		shapesList = append(shapesList, sectionShapes...)
		currentY = nextY + styling.Inches(0.2)
	}

	return DiagramElements{
		Shapes:  shapesList,
		Grouped: true,
		Bounds: &DiagramBounds{
			X:  layout.startX,
			Y:  layout.startY - styling.Inches(0.8),
			CX: layout.labelWidth + layout.chartWidth,
			CY: currentY - (layout.startY - styling.Inches(0.8)),
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
