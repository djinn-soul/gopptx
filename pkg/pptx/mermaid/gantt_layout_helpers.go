package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	// ganttAxisTickCount is how many dates the axis prints.
	ganttAxisTickCount  = 5
	ganttAxisDateFormat = "2006-01-02"
	// ganttMinBarWidth keeps a one-day task visible on a multi-month chart.
	ganttMinBarWidth styling.Length = 91440 // 0.1"
	// Critical tasks are red in Mermaid's own palette, in every theme.
	ganttCriticalFill   = "F8CECC"
	ganttCriticalStroke = "B85450"
)

type ganttLayout struct {
	startX        styling.Length
	startY        styling.Length
	labelWidth    styling.Length
	chartWidth    styling.Length
	rowHeight     styling.Length
	sectionHeight styling.Length
	axisHeight    styling.Length
}

func buildGanttTitleShape(title string, theme Theme, layout ganttLayout) shapes.Shape {
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.startX,
		layout.startY-styling.Inches(0.8),
		layout.labelWidth+layout.chartWidth,
		styling.Inches(0.6),
	).WithText(title).
		WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
}

func buildGanttAxisShapes(
	gantt *GanttDiagram,
	theme Theme,
	axisY styling.Length,
	layout ganttLayout,
	dateRange ganttDateRange,
) []shapes.Shape {
	axisShapes := make([]shapes.Shape, 0, 16)
	axisLine := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.startX+layout.labelWidth,
		axisY+layout.axisHeight/2,
		layout.chartWidth,
		styling.Emu(19050),
	).WithFill(shapes.NewShapeFill(theme.PrimaryStroke))
	axisShapes = append(axisShapes, axisLine)

	taskStarts := ganttAxisLabels(gantt, dateRange)
	for i, tick := range taskStarts {
		axisShapes = append(
			axisShapes,
			buildGanttTickShapes(i, len(taskStarts), tick, theme, axisY, layout)...,
		)
	}
	return axisShapes
}

// ganttAxisLabels are the axis captions: evenly spaced calendar dates once the
// tasks carry dates, and the raw start fields when none of them parse.
func ganttAxisLabels(gantt *GanttDiagram, dateRange ganttDateRange) []string {
	if dateRange.ok {
		ticks := dateRange.tickDates(ganttAxisTickCount)
		labels := make([]string, 0, len(ticks))
		for _, tick := range ticks {
			labels = append(labels, tick.Format(ganttAxisDateFormat))
		}
		return labels
	}
	if labels := collectGanttStartLabels(gantt); len(labels) > 0 {
		return labels
	}
	return []string{"T1", "T2", "T3"}
}

func buildGanttTickShapes(
	index int,
	totalTicks int,
	tick string,
	theme Theme,
	axisY styling.Length,
	layout ganttLayout,
) []shapes.Shape {
	tickX := layout.startX + layout.labelWidth +
		(styling.Length(index) * layout.chartWidth / styling.Length(max(1, totalTicks-1)))
	tickShape := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		tickX-styling.Emu(9525),
		axisY+layout.axisHeight/2-styling.Inches(0.06),
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
	return []shapes.Shape{tickShape, tickLabel}
}

func buildGanttSectionShapes(
	section GanttSection,
	theme Theme,
	startY styling.Length,
	layout ganttLayout,
	dateRange ganttDateRange,
) ([]shapes.Shape, styling.Length) {
	sectionShapes := make([]shapes.Shape, 0, len(section.Tasks)*2+1)
	sectionShape := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.startX,
		startY,
		layout.labelWidth+layout.chartWidth,
		layout.sectionHeight,
	).WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithText(section.Name).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
	sectionShapes = append(sectionShapes, sectionShape)
	currentY := startY + layout.sectionHeight
	for index, task := range section.Tasks {
		taskShapes := buildGanttTaskShapes(task, theme, currentY, layout, dateRange, index, len(section.Tasks))
		sectionShapes = append(sectionShapes, taskShapes...)
		currentY += layout.rowHeight
	}
	return sectionShapes, currentY
}

func buildGanttTaskShapes(
	task GanttTask,
	theme Theme,
	y styling.Length,
	layout ganttLayout,
	dateRange ganttDateRange,
	index int,
	taskCount int,
) []shapes.Shape {
	labelShape := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.startX,
		y,
		layout.labelWidth,
		layout.rowHeight,
	).WithText(task.Name).
		WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithVerticalAnchor(shapes.TextAnchorMiddle).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))

	barX, barWidth := ganttBarGeometry(task, layout, dateRange, index, taskCount)
	barY := y + styling.Inches(0.05)
	barHeight := layout.rowHeight - styling.Inches(0.1)

	if task.Schedule != nil && task.Schedule.IsMilestone {
		return []shapes.Shape{labelShape, ganttMilestoneShape(task, theme, barX, barY, barHeight)}
	}

	barText := task.Duration
	if barText == "" {
		barText = task.Start
	}
	fill, stroke := ganttBarColors(task.Status, theme)
	barShape := shapes.NewShape(
		shapes.ShapeTypeRoundedRectangle,
		barX,
		barY,
		barWidth,
		barHeight,
	).WithFill(shapes.NewShapeFill(fill)).
		WithLine(shapes.NewShapeLine(stroke, theme.LineWeight)).
		WithText(barText).
		WithVerticalAnchor(shapes.TextAnchorMiddle).
		WithAutoFit(shapes.TextAutoFitNormal)
	return []shapes.Shape{labelShape, barShape}
}

// ganttBarGeometry places a bar on the date axis. Tasks with no readable date
// keep the old even spacing so a chart written without dates still draws.
func ganttBarGeometry(
	task GanttTask,
	layout ganttLayout,
	dateRange ganttDateRange,
	index int,
	taskCount int,
) (styling.Length, styling.Length) {
	chartX := layout.startX + layout.labelWidth
	if task.Schedule == nil || !dateRange.ok {
		slot := layout.chartWidth / styling.Length(max(taskCount, 1))
		return chartX + styling.Length(index)*slot, slot - styling.Inches(0.1)
	}

	startFraction := dateRange.fraction(task.Schedule.Start)
	endFraction := dateRange.fraction(task.Schedule.End)
	barX := chartX + scaleLength(layout.chartWidth, startFraction)
	width := max(scaleLength(layout.chartWidth, endFraction-startFraction), ganttMinBarWidth)
	return barX, width
}

// ganttMilestoneShape draws a zero-length task as the diamond Mermaid uses.
func ganttMilestoneShape(
	task GanttTask,
	theme Theme,
	barX, barY, barHeight styling.Length,
) shapes.Shape {
	return shapes.NewShape(
		shapes.ShapeTypeDiamond,
		barX-barHeight/2,
		barY,
		barHeight,
		barHeight,
	).WithFill(shapes.NewShapeFill(theme.PrimaryStroke)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithAltText(task.Name)
}

// ganttBarColors maps Mermaid's status tags onto the theme, so a completed or
// critical task no longer looks identical to every other bar.
func ganttBarColors(status string, theme Theme) (string, string) {
	switch {
	case strings.Contains(status, "crit"):
		return ganttCriticalFill, ganttCriticalStroke
	case strings.Contains(status, "done"):
		return theme.SecondaryFill, theme.SecondaryStroke
	case strings.Contains(status, "active"):
		return theme.PrimaryStroke, theme.PrimaryStroke
	default:
		return theme.PrimaryFill, theme.PrimaryStroke
	}
}
