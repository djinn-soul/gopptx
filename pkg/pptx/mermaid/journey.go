package mermaid

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// JourneyTask represents a task in a user journey.
type JourneyTask struct {
	Description string
	Score       int
	Actors      []string
}

// JourneySection represents a section in a user journey.
type JourneySection struct {
	Title string
	Tasks []JourneyTask
}

// JourneyDiagram represents the parsed structure of a Mermaid journey diagram.
type JourneyDiagram struct {
	Title    string
	Sections []JourneySection
}

// renderJourney parses and renders a Mermaid journey diagram into PowerPoint elements.
func renderJourney(code string, theme Theme) DiagramElements {
	journey := parseJourney(code)
	return generateJourneyElements(journey, theme)
}

func parseJourney(code string) *JourneyDiagram {
	lines := ParseLines(code)
	journey := &JourneyDiagram{}
	var currentSection *JourneySection

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		if lower == "journey" {
			continue
		}

		if strings.HasPrefix(lower, "title ") {
			journey.Title = strings.TrimSpace(trimmed[6:])
			continue
		}

		if strings.HasPrefix(lower, "section ") {
			if currentSection != nil {
				journey.Sections = append(journey.Sections, *currentSection)
			}
			currentSection = &JourneySection{
				Title: strings.TrimSpace(trimmed[8:]),
			}
			continue
		}

		if strings.Contains(trimmed, ":") {
			parts := strings.Split(trimmed, ":")
			if len(parts) >= 2 {
				task := JourneyTask{
					Description: strings.TrimSpace(parts[0]),
				}

				score, err := strconv.Atoi(strings.TrimSpace(parts[1]))
				if err == nil {
					task.Score = score
				}

				if len(parts) >= 3 {
					actors := strings.SplitSeq(parts[2], ",")
					for actor := range actors {
						task.Actors = append(task.Actors, strings.TrimSpace(actor))
					}
				}

				if currentSection != nil {
					currentSection.Tasks = append(currentSection.Tasks, task)
				} else {
					// Default section if none defined
					currentSection = &JourneySection{Title: "Default"}
					currentSection.Tasks = append(currentSection.Tasks, task)
				}
			}
		}
	}

	if currentSection != nil {
		journey.Sections = append(journey.Sections, *currentSection)
	}

	return journey
}

func generateJourneyElements(journey *JourneyDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape

	if len(journey.Sections) == 0 {
		return createPlaceholder("journey (no data)", theme)
	}

	startX := styling.Inches(0.5)
	startY := styling.Inches(1.5)
	sectionWidth := styling.Inches(2.35)
	spacing := styling.Inches(0.22)
	maxCharsPerLine := 22
	maxTasksPerSection := 0
	for _, section := range journey.Sections {
		if len(section.Tasks) > maxTasksPerSection {
			maxTasksPerSection = len(section.Tasks)
		}
	}

	rowHeights := make([]styling.Length, maxTasksPerSection)
	for row := range maxTasksPerSection {
		maxLines := 3
		for _, section := range journey.Sections {
			if row >= len(section.Tasks) {
				continue
			}
			lines := estimateJourneyTaskLines(section.Tasks[row], maxCharsPerLine)
			if lines > maxLines {
				maxLines = lines
			}
		}
		rowHeights[row] = styling.Inches(0.28*float64(maxLines) + 0.46)
	}

	rowTopY := make([]styling.Length, maxTasksPerSection)
	currentRowY := startY + styling.Inches(0.5) + spacing
	for row := range maxTasksPerSection {
		rowTopY[row] = currentRowY
		currentRowY += rowHeights[row] + spacing
	}

	// Title
	if journey.Title != "" {
		titleShape := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			startX,
			startY-styling.Inches(0.8),
			styling.Inches(8),
			styling.Inches(0.6),
		).WithText(journey.Title).
			WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, titleShape)
	}

	for i, section := range journey.Sections {
		x := startX + styling.Length(i)*(sectionWidth+spacing)

		// Section Header
		sectionHeader := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			x,
			startY,
			sectionWidth,
			styling.Inches(0.5),
		).WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
			WithText(section.Title).
			WithAutoFit(shapes.TextAutoFitNormal).
			WithTextMargins(styling.Inches(0.1), styling.Inches(0.05), styling.Inches(0.1), styling.Inches(0.05))
		shapesList = append(shapesList, sectionHeader)

		for j, task := range section.Tasks {
			y := rowTopY[j]
			taskHeight := rowHeights[j]
			taskBox := shapes.NewShape(
				shapes.ShapeTypeRoundedRectangle,
				x,
				y,
				sectionWidth,
				taskHeight,
			).WithFill(shapes.NewShapeFill(theme.PrimaryFill)).
				WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
				WithText(journeyTaskText(task)).
				WithAutoFit(shapes.TextAutoFitNormal).
				WithVerticalAnchor(shapes.TextAnchorTop).
				WithTextMargins(styling.Inches(0.08), styling.Inches(0.03), styling.Inches(0.08), styling.Inches(0.03))
			shapesList = append(shapesList, taskBox)
		}
	}

	return DiagramElements{
		Shapes:  shapesList,
		Grouped: true,
		Bounds: &DiagramBounds{
			X:  startX,
			Y:  startY - styling.Inches(0.8),
			CX: styling.Length(len(journey.Sections)) * (sectionWidth + spacing),
			CY: (currentRowY - spacing) - (startY - styling.Inches(0.8)),
		},
	}
}

func estimateJourneyTaskLines(task JourneyTask, maxCharsPerLine int) int {
	descLines := wrappedLineCount(task.Description, maxCharsPerLine)
	metaLine := wrappedLineCount(journeyTaskMeta(task), maxCharsPerLine)
	return descLines + metaLine
}

func wrappedLineCount(text string, maxCharsPerLine int) int {
	if maxCharsPerLine <= 0 {
		return 1
	}
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return 1
	}

	words := strings.Fields(trimmed)
	if len(words) == 0 {
		return 1
	}
	lines := 1
	current := 0
	for _, word := range words {
		wordLen := len(word)
		if current == 0 {
			current = wordLen
			continue
		}
		if current+1+wordLen > maxCharsPerLine {
			lines++
			current = wordLen
			continue
		}
		current += 1 + wordLen
	}
	return lines
}

func journeyTaskText(task JourneyTask) string {
	return fmt.Sprintf("%s\n%s", task.Description, journeyTaskMeta(task))
}

func journeyTaskMeta(task JourneyTask) string {
	actors := strings.TrimSpace(strings.Join(task.Actors, ", "))
	if actors == "" {
		return fmt.Sprintf("Score: %d", task.Score)
	}
	return fmt.Sprintf("Score: %d %s", task.Score, actors)
}
