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
	// Section is the "section <name>" heading this entry falls under, empty
	// when the timeline declares none.
	Section string
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
	state := &timelineParseState{}

	for _, line := range lines {
		state.current = consumeTimelineLine(strings.TrimSpace(line), timeline, state)
	}

	if state.current != nil {
		timeline.Events = append(timeline.Events, *state.current)
	}

	return timeline
}

// timelineParseState carries what a line needs from the lines before it: the
// entry being filled in, and the section heading currently in force.
type timelineParseState struct {
	current *TimelineEvent
	section string
}

func consumeTimelineLine(trimmed string, timeline *TimelineDiagram, state *timelineParseState) *TimelineEvent {
	currentEvent := state.current
	if trimmed == "" {
		return currentEvent
	}

	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "timeline") {
		return currentEvent
	}

	if title, ok := extractTimelineTitle(trimmed, lower); ok {
		timeline.Title = title
		return currentEvent
	}

	// A "section" heading groups the entries that follow it. Without this it
	// fell through to the event branch and was drawn as an event of its own.
	if after, ok := strings.CutPrefix(lower, "section"); ok && strings.TrimSpace(after) != "" {
		state.section = strings.TrimSpace(trimmed[len("section"):])
		return currentEvent
	}

	if date, event, ok := splitTimelineEntry(trimmed); ok {
		return consumeTimelineEntry(timeline, state, date, event)
	}

	if currentEvent == nil {
		return nil
	}

	if isLikelyTimelineDate(trimmed) {
		timeline.Events = append(timeline.Events, *currentEvent)
		return &TimelineEvent{Date: trimmed, Events: []string{}, Section: state.section}
	}

	currentEvent.Events = append(currentEvent.Events, trimmed)
	return currentEvent
}

func extractTimelineTitle(trimmed string, lower string) (string, bool) {
	if !strings.HasPrefix(lower, "title") {
		return "", false
	}
	return strings.TrimSpace(trimmed[5:]), true
}

func splitTimelineEntry(trimmed string) (string, string, bool) {
	datePart, eventPart, ok := strings.Cut(trimmed, ":")
	if !ok {
		return "", "", false
	}
	return strings.TrimSpace(datePart), strings.TrimSpace(eventPart), true
}

// splitTimelineEvents splits the part after the date. Mermaid separates several
// events on one line with a colon, so "2020 : Alpha : Beta" is two events, not
// one event named "Alpha : Beta".
func splitTimelineEvents(eventPart string) []string {
	var events []string
	for part := range strings.SplitSeq(eventPart, ":") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			events = append(events, trimmed)
		}
	}
	return events
}

func consumeTimelineEntry(
	timeline *TimelineDiagram,
	state *timelineParseState,
	date string,
	event string,
) *TimelineEvent {
	events := splitTimelineEvents(event)
	if date != "" {
		if state.current != nil {
			timeline.Events = append(timeline.Events, *state.current)
		}
		return &TimelineEvent{Date: date, Events: events, Section: state.section}
	}

	if state.current != nil {
		state.current.Events = append(state.current.Events, events...)
	}
	return state.current
}

func isLikelyTimelineDate(trimmed string) bool {
	return !strings.Contains(trimmed, " ")
}

func generateTimelineElements(timeline *TimelineDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape

	if len(timeline.Events) == 0 {
		return createPlaceholder("timeline (no data)", theme)
	}

	layout := timelineLayout{
		startX:       styling.Inches(1),
		startY:       styling.Inches(3.5),
		eventWidth:   styling.Inches(1.8),
		eventSpacing: styling.Inches(0.5),
	}
	totalWidth := timelineTotalWidth(len(timeline.Events), layout)
	lineY := layout.startY

	if timeline.Title != "" {
		shapesList = append(shapesList, timelineTitleShape(timeline.Title, layout, theme))
	}

	shapesList = append(shapesList, timelineLineShape(totalWidth, lineY, layout, theme))

	previousSection := ""
	for i, event := range timeline.Events {
		if event.Section != "" && event.Section != previousSection {
			shapesList = append(shapesList, timelineSectionShape(event.Section, i, lineY, layout, theme))
			previousSection = event.Section
		}
		shapesList = append(shapesList, timelineEventShapes(event, i, lineY, layout, theme)...)
	}

	return DiagramElements{
		Shapes:  shapesList,
		Grouped: true,
		Bounds: &DiagramBounds{
			X:  layout.startX,
			Y:  styling.Inches(0.5),
			CX: totalWidth,
			CY: styling.Inches(6),
		},
	}
}

type timelineLayout struct {
	startX       styling.Length
	startY       styling.Length
	eventWidth   styling.Length
	eventSpacing styling.Length
}

func timelineTotalWidth(eventCount int, layout timelineLayout) styling.Length {
	return styling.Length(eventCount) * (layout.eventWidth + layout.eventSpacing)
}

func timelineTitleShape(title string, layout timelineLayout, theme Theme) shapes.Shape {
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.startX,
		styling.Inches(0.5),
		styling.Inches(8),
		styling.Inches(0.6),
	).WithText(title).
		WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
}

func timelineLineShape(
	totalWidth styling.Length,
	lineY styling.Length,
	layout timelineLayout,
	theme Theme,
) shapes.Shape {
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		layout.startX,
		lineY,
		totalWidth,
		styling.Emu(50000),
	).WithFill(shapes.NewShapeFill(theme.PrimaryStroke))
}

func timelineEventShapes(
	event TimelineEvent,
	index int,
	lineY styling.Length,
	layout timelineLayout,
	theme Theme,
) []shapes.Shape {
	x := layout.startX + styling.Length(index)*(layout.eventWidth+layout.eventSpacing)
	out := make([]shapes.Shape, 0, len(event.Events)+2)
	out = append(out, timelineMarkerShape(x, lineY, layout, theme))
	out = append(out, timelineDateLabelShape(x, lineY, layout.eventWidth, event.Date, theme))

	isAbove := index%2 == 0
	baseY := timelineEventBaseY(lineY, isAbove)
	for j, text := range event.Events {
		y := timelineEventRowY(baseY, j, isAbove)
		out = append(out, timelineEventBoxShape(x, y, layout.eventWidth, text, theme))
	}
	return out
}

// timelineSectionShape labels the band a run of entries belongs to, drawn above
// the axis at the first entry of the section.
func timelineSectionShape(
	name string,
	index int,
	lineY styling.Length,
	layout timelineLayout,
	theme Theme,
) shapes.Shape {
	x := layout.startX + styling.Length(index)*(layout.eventWidth+layout.eventSpacing)
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		x,
		lineY-styling.Inches(1.9),
		layout.eventWidth,
		styling.Inches(0.32),
	).WithText(name).
		WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
}

func timelineMarkerShape(x styling.Length, lineY styling.Length, layout timelineLayout, theme Theme) shapes.Shape {
	markerSize := styling.Inches(0.2)
	return shapes.NewShape(
		shapes.ShapeTypeEllipse,
		x+layout.eventWidth/2-markerSize/2,
		lineY-markerSize/2,
		markerSize,
		markerSize,
	).WithFill(shapes.NewShapeFill(theme.PrimaryStroke)).
		WithLine(shapes.NewShapeLine(theme.Background, theme.LineWeight))
}

func timelineDateLabelShape(
	x styling.Length,
	lineY styling.Length,
	eventWidth styling.Length,
	date string,
	theme Theme,
) shapes.Shape {
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		x,
		lineY+styling.Inches(0.2),
		eventWidth,
		styling.Inches(0.35),
	).WithText(date).
		WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
}

func timelineEventBaseY(lineY styling.Length, isAbove bool) styling.Length {
	if isAbove {
		return lineY - styling.Inches(0.6)
	}
	return lineY + styling.Inches(0.8)
}

func timelineEventRowY(baseY styling.Length, rowIndex int, isAbove bool) styling.Length {
	offset := styling.Length(rowIndex) * styling.Inches(0.45)
	if isAbove {
		return max(0, baseY-offset)
	}
	return max(0, baseY+offset)
}

func timelineEventBoxShape(
	x styling.Length,
	y styling.Length,
	eventWidth styling.Length,
	text string,
	theme Theme,
) shapes.Shape {
	return shapes.NewShape(
		shapes.ShapeTypeRoundedRectangle,
		x,
		y,
		eventWidth,
		styling.Inches(0.4),
	).WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithText(text).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
}
