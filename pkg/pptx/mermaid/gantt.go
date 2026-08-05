package mermaid

import (
	"strings"
	"time"

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
	// Schedule holds the resolved dates. It is nil when the task carries no
	// parsable date, in which case the bar falls back to even spacing.
	Schedule *GanttSchedule
}

// GanttSchedule is one task's resolved position on the calendar.
type GanttSchedule struct {
	Start       time.Time
	End         time.Time
	IsMilestone bool
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

// isGanttStatusTag reports whether a field is one of the leading keywords a
// task may carry before its id.
func isGanttStatusTag(field string) bool {
	switch field {
	case "done", "active", "crit", "milestone":
		return true
	default:
		return false
	}
}

// parseGanttTask reads "Name : [status,]... [id,] [start,] duration".
//
// The fields were previously assigned by position alone, so a task written
// ":done, a1, 2026-01-01, 5d" took "done" for its id and "a1" for its start
// date, which is what put task ids on the date axis.
func parseGanttTask(trimmed string) (GanttTask, bool) {
	taskName, taskDetails, ok := strings.Cut(trimmed, ":")
	if !ok {
		return GanttTask{}, false
	}
	task := GanttTask{Name: strings.TrimSpace(taskName)}

	fields := make([]string, 0, 4)
	for detail := range strings.SplitSeq(taskDetails, ",") {
		detail = strings.TrimSpace(detail)
		if detail == "" {
			continue
		}
		if len(fields) == 0 && isGanttStatusTag(strings.ToLower(detail)) {
			task.Status = appendGanttStatus(task.Status, strings.ToLower(detail))
			continue
		}
		fields = append(fields, detail)
	}

	switch len(fields) {
	case 0:
	case 1:
		task.Duration = fields[0]
	case 2:
		// "start, duration" — an id alone never appears without a duration.
		task.Start, task.Duration = fields[0], fields[1]
	default:
		task.ID, task.Start, task.Duration = fields[0], fields[1], fields[2]
	}
	return task, true
}

func appendGanttStatus(existing, status string) string {
	if existing == "" {
		return status
	}
	return existing + " " + status
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
	resolveGanttSchedule(gantt)
	dateRange := ganttRangeOf(gantt)

	axisY := layout.startY
	currentY := axisY + layout.axisHeight + styling.Inches(0.1)
	if gantt.Title != "" {
		shapesList = append(shapesList, buildGanttTitleShape(gantt.Title, theme, layout))
	}

	shapesList = append(shapesList, buildGanttAxisShapes(gantt, theme, axisY, layout, dateRange)...)

	for _, section := range gantt.Sections {
		sectionShapes, nextY := buildGanttSectionShapes(section, theme, currentY, layout, dateRange)
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
