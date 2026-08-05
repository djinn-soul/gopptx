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
	// Tag is the release name a commit carries, from `commit tag: "v1"`.
	Tag string
	// IsCherryPick marks a commit replayed from another branch.
	IsCherryPick bool
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
		case strings.HasPrefix(lower, "cherry-pick"):
			// A cherry-pick replays a commit onto the current branch. It used to
			// be ignored outright, so the branch came up a commit short.
			diagram.Commits = append(diagram.Commits, GitCommitInfo{
				Branch:       currentBranch,
				X:            currentX,
				Y:            diagram.Branches[currentBranch],
				Label:        gitAttributeValue(trimmed, "id:"),
				IsCherryPick: true,
			})
			currentX++
		case strings.HasPrefix(lower, "commit"):
			diagram.Commits = append(diagram.Commits, GitCommitInfo{
				Branch: currentBranch,
				X:      currentX,
				Y:      diagram.Branches[currentBranch],
				Label:  gitAttributeValue(trimmed, "id:"),
				Tag:    gitAttributeValue(trimmed, "tag:"),
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

// gitAttributeValue reads a `key: "value"` attribute off a gitGraph statement.
func gitAttributeValue(line, key string) string {
	_, after, found := strings.Cut(line, key)
	if !found {
		return ""
	}
	value := strings.TrimSpace(after)
	// Attributes may be chained: `commit id: "x" tag: "v1"`, so a quoted value
	// ends at its closing quote rather than at the end of the line.
	if quoted, _, ok := strings.Cut(strings.TrimPrefix(value, `"`), `"`); ok && strings.HasPrefix(value, `"`) {
		return quoted
	}
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
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
	maxCommitX := gitgraphMaxCommitX(diagram.Commits)
	shapesList = append(
		shapesList,
		gitgraphBranchShapes(diagram, theme, startX, startY, commitSpacing, branchSpacing, maxCommitX)...)

	for i, commit := range diagram.Commits {
		x := startX + styling.Length(commit.X)*commitSpacing
		y := startY + styling.Length(commit.Y)*branchSpacing

		color := gitgraphBranchColor(diagram.Branches[commit.Branch], theme)
		shapesList = append(shapesList, gitgraphCommitDot(x, y, color, theme))

		if commit.Label != "" {
			shapesList = append(shapesList, gitgraphCommitLabel(commit.Label, x, y, theme))
		}
		if commit.Tag != "" {
			shapesList = append(shapesList, gitgraphTagShape(commit.Tag, x, y, theme))
		}
		if commit.IsCherryPick {
			shapesList = append(shapesList, gitgraphCherryPickMark(x, y, color, theme))
		}

		if connector, ok := gitgraphPrevBranchConnector(
			diagram.Commits,
			i,
			commit.Branch,
			x,
			y,
			color,
			theme,
			startX,
			startY,
			commitSpacing,
			branchSpacing,
		); ok {
			connectors = append(connectors, connector)
		}

		if connector, ok := gitgraphMergeConnector(
			diagram,
			i,
			commit,
			x,
			y,
			theme,
			startX,
			startY,
			commitSpacing,
			branchSpacing,
		); ok {
			connectors = append(connectors, connector)
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

func gitgraphLabelWidth(text string, minInches, maxInches float64) styling.Length {
	width := minInches + float64(len(text))*0.04
	if width > maxInches {
		width = maxInches
	}
	if width < minInches {
		width = minInches
	}
	return styling.Inches(width)
}

func gitgraphMaxCommitX(commits []GitCommitInfo) int {
	maxX := 0
	for _, commit := range commits {
		if commit.X > maxX {
			maxX = commit.X
		}
	}
	return maxX
}

func gitgraphBranchColor(yIdx int, theme Theme) string {
	if yIdx > 0 {
		return theme.SecondaryStroke
	}
	return theme.PrimaryStroke
}

func gitgraphBranchShapes(
	diagram *GitGraphDiagram,
	theme Theme,
	startX styling.Length,
	startY styling.Length,
	commitSpacing styling.Length,
	branchSpacing styling.Length,
	maxCommitX int,
) []shapes.Shape {
	out := make([]shapes.Shape, 0, len(diagram.Branches)*2)
	for name, yIdx := range diagram.Branches {
		y := startY + styling.Length(yIdx)*branchSpacing
		color := gitgraphBranchColor(yIdx, theme)
		out = append(out, shapes.NewShape(
			shapes.ShapeTypeRectangle,
			startX,
			y,
			styling.Length(maxCommitX)*commitSpacing,
			styling.Emu(20000),
		).WithFill(shapes.NewShapeFill(color)))
		out = append(out, gitgraphBranchLabel(name, y, startX, theme))
	}
	return out
}

func gitgraphBranchLabel(name string, y styling.Length, startX styling.Length, theme Theme) shapes.Shape {
	branchLabelWidth := gitgraphLabelWidth(name, 0.8, 1.6)
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		startX-branchLabelWidth-styling.Inches(0.1),
		y-styling.Inches(0.15),
		branchLabelWidth,
		styling.Inches(0.3),
	).WithText(name).
		WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithAutoFit(shapes.TextAutoFitNormal)
}

func gitgraphCommitDot(x styling.Length, y styling.Length, color string, theme Theme) shapes.Shape {
	dotSize := styling.Inches(0.2)
	return shapes.NewShape(
		shapes.ShapeTypeEllipse,
		x-dotSize/2,
		y-dotSize/2,
		dotSize,
		dotSize,
	).WithFill(shapes.NewShapeFill(color)).
		WithLine(shapes.NewShapeLine(theme.Background, theme.LineWeight))
}

// gitgraphTagShape is the release flag above a tagged commit, drawn as the
// pennant Mermaid uses so it reads differently from a commit id.
func gitgraphTagShape(tag string, x styling.Length, y styling.Length, theme Theme) shapes.Shape {
	width := gitgraphLabelWidth(tag, 0.5, 1.4)
	return shapes.NewShape(
		shapes.ShapeTypeHomePlate,
		x-width/2,
		y-styling.Inches(0.5),
		width,
		styling.Inches(0.26),
	).WithText(tag).
		WithFill(shapes.NewShapeFill(theme.SecondaryFill)).
		WithLine(shapes.NewShapeLine(theme.PrimaryStroke, theme.LineWeight)).
		WithVerticalAnchor(shapes.TextAnchorMiddle).
		WithAutoFit(shapes.TextAutoFitNormal)
}

// gitgraphCherryPickMark rings a replayed commit so it is distinguishable from
// an ordinary one on the same branch.
func gitgraphCherryPickMark(x, y styling.Length, color string, theme Theme) shapes.Shape {
	size := styling.Inches(0.3)
	shape := shapes.NewShape(shapes.ShapeTypeEllipse, x-size/2, y-size/2, size, size).
		WithLine(shapes.NewShapeLine(color, theme.LineWeight*2))
	shape.Fill = nil
	return shape
}

func gitgraphCommitLabel(label string, x styling.Length, y styling.Length, theme Theme) shapes.Shape {
	commitLabelWidth := gitgraphLabelWidth(label, 0.8, 1.8)
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		x-commitLabelWidth/2,
		y+styling.Inches(0.15),
		commitLabelWidth,
		styling.Inches(0.3),
	).WithText(label).
		WithFill(shapes.NewShapeFill(theme.Background)).
		WithLine(shapes.NewShapeLine(theme.SecondaryStroke, theme.LineWeight)).
		WithAutoFit(shapes.TextAutoFitNormal)
}
