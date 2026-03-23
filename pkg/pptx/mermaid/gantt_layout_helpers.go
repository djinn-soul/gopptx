package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
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

	taskStarts := collectGanttStartLabels(gantt)
	if len(taskStarts) == 0 {
		taskStarts = []string{"T1", "T2", "T3"}
	}
	for i, tick := range taskStarts {
		axisShapes = append(
			axisShapes,
			buildGanttTickShapes(i, len(taskStarts), tick, theme, axisY, layout)...,
		)
	}
	return axisShapes
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
	for _, task := range section.Tasks {
		taskShapes := buildGanttTaskShapes(task, theme, currentY, layout)
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

	barText := task.Duration
	if barText == "" {
		barText = task.Start
	}
	barShape := shapes.NewShape(
		shapes.ShapeTypeRoundedRectangle,
		layout.startX+layout.labelWidth+styling.Inches(0.5),
		y+styling.Inches(0.05),
		styling.Inches(2.2),
		layout.rowHeight-styling.Inches(0.1),
	).WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithText(barText).
		WithVerticalAnchor(shapes.TextAnchorMiddle).
		WithAutoFit(shapes.TextAutoFitNormal)
	return []shapes.Shape{labelShape, barShape}
}
