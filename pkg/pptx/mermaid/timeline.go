package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// TimelineEvent represents an event in a timeline.
type TimelineEvent struct {
	Date   string
	Events []string
}

// TimelineDiagram represents the parsed structure of a Mermaid timeline diagram.
type TimelineDiagram struct {
	Title  string
	Events []TimelineEvent
}

// renderTimeline parses and renders a Mermaid timeline diagram into PowerPoint elements.
func renderTimeline(code string, theme Theme) DiagramElements {
	timeline := parseTimeline(code)
	return generateTimelineElements(timeline, theme)
}

func parseTimeline(code string) *TimelineDiagram {
	lines := ParseLines(code)
	timeline := &TimelineDiagram{}
	var currentEvent *TimelineEvent

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		if strings.HasPrefix(lower, "timeline") {
			continue
		}

		if strings.HasPrefix(lower, "title") {
			timeline.Title = strings.TrimSpace(trimmed[5:])
			continue
		}

		if strings.Contains(trimmed, ":") {
			parts := strings.Split(trimmed, ":")
			date := strings.TrimSpace(parts[0])
			event := strings.TrimSpace(parts[1])

			if date != "" {
				if currentEvent != nil {
					timeline.Events = append(timeline.Events, *currentEvent)
				}
				currentEvent = &TimelineEvent{
					Date:   date,
					Events: []string{event},
				}
			} else if currentEvent != nil {
				currentEvent.Events = append(currentEvent.Events, event)
			}
		} else if trimmed != "" {
			// Could be a date or an event
			if currentEvent != nil && !strings.Contains(trimmed, " ") {
				// If it's a single word, it might be a date
				timeline.Events = append(timeline.Events, *currentEvent)
				currentEvent = &TimelineEvent{
					Date:   trimmed,
					Events: []string{},
				}
			} else if currentEvent != nil {
				currentEvent.Events = append(currentEvent.Events, trimmed)
			}
		}
	}

	if currentEvent != nil {
		timeline.Events = append(timeline.Events, *currentEvent)
	}

	return timeline
}

func generateTimelineElements(timeline *TimelineDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape
	
	if len(timeline.Events) == 0 {
		return createPlaceholder("timeline (no data)", theme)
	}

	// Layout parameters
	startX := styling.Inches(1)
	startY := styling.Inches(3.5)
	eventWidth := styling.Inches(1.8)
	eventSpacing := styling.Inches(0.5)
	
	// Title
	if timeline.Title != "" {
		titleShape := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			startX,
			styling.Inches(0.5),
			styling.Inches(8),
			styling.Inches(0.6),
		).WithText(timeline.Title).
			WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, titleShape)
	}

	// Main timeline line
	totalWidth := styling.Length(len(timeline.Events)) * (eventWidth + eventSpacing)
	lineY := startY
	timelineLine := shapes.NewShape(
		shapes.ShapeTypeRectangle,
		startX,
		lineY,
		totalWidth,
		styling.Emu(50000),
	).WithFill(shapes.NewShapeFill(theme.PrimaryStroke))
	shapesList = append(shapesList, timelineLine)

	for i, event := range timeline.Events {
		x := startX + styling.Length(i)*(eventWidth+eventSpacing)
		
		// Date marker
		markerSize := styling.Inches(0.2)
		marker := shapes.NewShape(
			shapes.ShapeTypeEllipse,
			x+eventWidth/2-markerSize/2,
			lineY-markerSize/2,
			markerSize,
			markerSize,
		).WithFill(shapes.NewShapeFill(theme.PrimaryStroke)).
			WithLine(shapes.NewShapeLine(theme.Background, theme.LineWeight))
		shapesList = append(shapesList, marker)

		// Date label
		dateLabel := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			x,
			lineY+styling.Inches(0.2),
			eventWidth,
			styling.Inches(0.35),
		).WithText(event.Date).
			WithFill(shapes.NewShapeFill(theme.Background)).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
		shapesList = append(shapesList, dateLabel)

		// Event boxes - alternate above and below
		isAbove := i%2 == 0
		var eventY styling.Length
		if isAbove {
			eventY = lineY - styling.Inches(0.6)
		} else {
			eventY = lineY + styling.Inches(0.8)
		}

		for j, e := range event.Events {
			var y styling.Length
			if isAbove {
				y = eventY - styling.Length(j)*styling.Inches(0.45)
			} else {
				y = eventY + styling.Length(j)*styling.Inches(0.45)
			}

			if y < 0 {
				y = 0
			}

			eventBox := shapes.NewShape(
				shapes.ShapeTypeRoundedRectangle,
				x,
				y,
				eventWidth,
				styling.Inches(0.4),
			).WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
				WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
				WithText(e).
				WithAutoFit(shapes.TextAutoFitNormal).
				WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
			shapesList = append(shapesList, eventBox)
		}
	}

	return DiagramElements{
		Shapes:  shapesList,
		Grouped: true,
		Bounds: &DiagramBounds{
			X:  startX,
			Y:  styling.Inches(0.5),
			CX: totalWidth,
			CY: styling.Inches(6),
		},
	}
}
