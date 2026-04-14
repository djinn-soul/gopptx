package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func gitgraphPrevBranchConnector(
	commits []GitCommitInfo,
	i int,
	branch string,
	x styling.Length,
	y styling.Length,
	color string,
	theme Theme,
	startX styling.Length,
	startY styling.Length,
	commitSpacing styling.Length,
	branchSpacing styling.Length,
) (shapes.Connector, bool) {
	prevCommit, ok := gitgraphFindPreviousOnBranch(commits, i, branch)
	if !ok {
		return shapes.Connector{}, false
	}
	prevX := startX + styling.Length(prevCommit.X)*commitSpacing
	prevY := startY + styling.Length(prevCommit.Y)*branchSpacing
	return shapes.NewConnector(
		shapes.ConnectorTypeStraight,
		prevX, prevY, x, y,
	).WithLine(shapes.NewShapeLine(color, theme.LineWeight)), true
}

func gitgraphFindPreviousOnBranch(commits []GitCommitInfo, fromIndex int, branch string) (GitCommitInfo, bool) {
	for j := fromIndex - 1; j >= 0; j-- {
		if commits[j].Branch == branch {
			return commits[j], true
		}
	}
	return GitCommitInfo{}, false
}

func gitgraphMergeConnector(
	diagram *GitGraphDiagram,
	i int,
	commit GitCommitInfo,
	x styling.Length,
	y styling.Length,
	theme Theme,
	startX styling.Length,
	startY styling.Length,
	commitSpacing styling.Length,
	branchSpacing styling.Length,
) (shapes.Connector, bool) {
	if !commit.IsMerge {
		return shapes.Connector{}, false
	}
	fromCommit, ok := gitgraphFindPreviousOnBranch(diagram.Commits, i, commit.MergeFrom)
	if !ok {
		return shapes.Connector{}, false
	}
	fromX := startX + styling.Length(fromCommit.X)*commitSpacing
	fromYIdx := diagram.Branches[commit.MergeFrom]
	fromY := startY + styling.Length(fromYIdx)*branchSpacing
	fromColor := gitgraphBranchColor(fromYIdx, theme)
	return shapes.NewConnector(
		shapes.ConnectorTypeStraight,
		fromX, fromY, x, y,
	).WithLine(shapes.NewShapeLine(fromColor, theme.LineWeight)), true
}
