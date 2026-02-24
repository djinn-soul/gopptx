package mermaid

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// GitCommitInfo represents a commit in a git graph.
type GitCommitInfo struct {
	Branch    string
	X         int
	Y         int
	Label     string
	IsMerge   bool
	MergeFrom string
}

// GitGraphDiagram represents the parsed structure of a Mermaid git graph.
type GitGraphDiagram struct {
	Commits  []GitCommitInfo
	Branches map[string]int
}

// renderGitGraph parses and renders a Mermaid git graph into PowerPoint elements.
func renderGitGraph(code string, theme Theme) DiagramElements {
	diagram := parseGitGraph(code)
	return generateGitGraphElements(diagram, theme)
}

func parseGitGraph(code string) *GitGraphDiagram {
	lines := ParseLines(code)
	diagram := &GitGraphDiagram{
		Branches: make(map[string]int),
	}

	currentBranch := "main"
	diagram.Branches["main"] = 0
	nextBranchIdx := 1
	currentX := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		if lower == "gitgraph" {
			continue
		}

		switch {
		case strings.HasPrefix(lower, "commit"):
			label := ""
			if strings.Contains(trimmed, "id:") {
				parts := strings.Split(trimmed, "id:")
				label = strings.Trim(strings.TrimSpace(parts[1]), "\"")
			}

			diagram.Commits = append(diagram.Commits, GitCommitInfo{
				Branch: currentBranch,
				X:      currentX,
				Y:      diagram.Branches[currentBranch],
				Label:  label,
			})
			currentX++
		case strings.HasPrefix(lower, "branch "):
			branchName := strings.TrimSpace(trimmed[7:])
			if _, ok := diagram.Branches[branchName]; !ok {
				diagram.Branches[branchName] = nextBranchIdx
				nextBranchIdx++
			}
		case strings.HasPrefix(lower, "checkout "):
			currentBranch = strings.TrimSpace(trimmed[9:])
		case strings.HasPrefix(lower, "merge "):
			fromBranch := strings.TrimSpace(trimmed[6:])
			diagram.Commits = append(diagram.Commits, GitCommitInfo{
				Branch:    currentBranch,
				X:         currentX,
				Y:         diagram.Branches[currentBranch],
				IsMerge:   true,
				MergeFrom: fromBranch,
			})
			currentX++
		}
	}

	return diagram
}

func generateGitGraphElements(diagram *GitGraphDiagram, theme Theme) DiagramElements {
	var shapesList []shapes.Shape
	var connectors []shapes.Connector

	if len(diagram.Commits) == 0 {
		return createPlaceholder("gitGraph (no data)", theme)
	}

	startX := styling.Inches(1)
	startY := styling.Inches(2)
	commitSpacing := styling.Inches(0.8)
	branchSpacing := styling.Inches(0.6)
	labelWidth := func(text string, minInches, maxInches float64) styling.Length {
		width := minInches + float64(len(text))*0.04
		if width > maxInches {
			width = maxInches
		}
		if width < minInches {
			width = minInches
		}
		return styling.Inches(width)
	}

	// Draw branch lines
	for name, yIdx := range diagram.Branches {
		y := startY + styling.Length(yIdx)*branchSpacing
		color := theme.PrimaryStroke
		if yIdx > 0 {
			color = theme.SecondaryStroke
		}

		// Find max X for this branch
		maxX := 0
		for _, commit := range diagram.Commits {
			if commit.X > maxX {
				maxX = commit.X
			}
		}

		line := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			startX,
			y,
			styling.Length(maxX)*commitSpacing,
			styling.Emu(20000),
		).WithFill(shapes.NewShapeFill(color))
		shapesList = append(shapesList, line)

		// Branch label
		branchLabelWidth := labelWidth(name, 0.8, 1.6)
		label := shapes.NewShape(
			shapes.ShapeTypeRectangle,
			startX-branchLabelWidth-styling.Inches(0.1),
			y-styling.Inches(0.15),
			branchLabelWidth,
			styling.Inches(0.3),
		).WithText(name).
			WithFill(shapes.NewShapeFill(theme.Background)).
			WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
			WithAutoFit(shapes.TextAutoFitNormal)
		shapesList = append(shapesList, label)
	}

	// Draw commits and connections
	for i, commit := range diagram.Commits {
		x := startX + styling.Length(commit.X)*commitSpacing
		y := startY + styling.Length(commit.Y)*branchSpacing

		color := theme.PrimaryStroke
		if diagram.Branches[commit.Branch] > 0 {
			color = theme.SecondaryStroke
		}

		// Commit dot
		dotSize := styling.Inches(0.2)
		dot := shapes.NewShape(
			shapes.ShapeTypeEllipse,
			x-dotSize/2,
			y-dotSize/2,
			dotSize,
			dotSize,
		).WithFill(shapes.NewShapeFill(color)).
			WithLine(shapes.NewShapeLine(theme.Background, theme.LineWeight))
		shapesList = append(shapesList, dot)

		if commit.Label != "" {
			commitLabelWidth := labelWidth(commit.Label, 0.8, 1.8)
			label := shapes.NewShape(
				shapes.ShapeTypeRectangle,
				x-commitLabelWidth/2,
				y+styling.Inches(0.15),
				commitLabelWidth,
				styling.Inches(0.3),
			).WithText(commit.Label).
				WithFill(shapes.NewShapeFill(theme.Background)).
				WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
				WithAutoFit(shapes.TextAutoFitNormal)
			shapesList = append(shapesList, label)
		}

		// Horizontal connection to previous commit on same branch
		for j := i - 1; j >= 0; j-- {
			if diagram.Commits[j].Branch == commit.Branch {
				prevX := startX + styling.Length(diagram.Commits[j].X)*commitSpacing
				prevY := startY + styling.Length(diagram.Commits[j].Y)*branchSpacing

				connector := shapes.NewConnector(
					shapes.ConnectorTypeStraight,
					prevX, prevY, x, y,
				).WithLine(shapes.NewShapeLine(color, theme.LineWeight))
				connectors = append(connectors, connector)
				break
			}
		}

		// Merge connection
		if commit.IsMerge {
			fromYIdx := diagram.Branches[commit.MergeFrom]
			fromY := startY + styling.Length(fromYIdx)*branchSpacing

			fromColor := theme.PrimaryStroke
			if fromYIdx > 0 {
				fromColor = theme.SecondaryStroke
			}

			// Find last commit on fromBranch before current X
			var lastFromX styling.Length
			found := false
			for j := i - 1; j >= 0; j-- {
				if diagram.Commits[j].Branch == commit.MergeFrom {
					lastFromX = startX + styling.Length(diagram.Commits[j].X)*commitSpacing
					found = true
					break
				}
			}

			if found {
				connector := shapes.NewConnector(
					shapes.ConnectorTypeStraight,
					lastFromX, fromY, x, y,
				).WithLine(shapes.NewShapeLine(fromColor, theme.LineWeight))
				connectors = append(connectors, connector)
			}
		}
	}

	return DiagramElements{
		Shapes:     shapesList,
		Connectors: connectors,
		Grouped:    true,
		Bounds: &DiagramBounds{
			X:  startX - styling.Inches(1),
			Y:  startY - styling.Inches(0.5),
			CX: styling.Inches(8),
			CY: styling.Inches(4),
		},
	}
}
