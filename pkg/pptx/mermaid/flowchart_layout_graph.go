package mermaid

import (
	"sort"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

//nolint:gocognit // Graph layout uses explicit branching to preserve stable node placement heuristics.
func (s *flowchartRenderState) layoutByConnections(connections []FlowConnection) bool {
	if len(s.nodes) == 0 || len(connections) == 0 {
		return false
	}

	depth := make(map[string]int, len(s.nodes))
	for _, n := range s.nodes {
		depth[n.ID] = 0
	}

	for range len(s.nodes) {
		changed := false
		for _, c := range connections {
			fromDepth, fromOK := depth[c.From]
			toDepth, toOK := depth[c.To]
			if !fromOK || !toOK {
				continue
			}
			if toDepth < fromDepth+1 {
				depth[c.To] = fromDepth + 1
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	buckets := make(map[int][]string, len(s.nodes))
	maxDepth := 0
	for _, n := range s.nodes {
		d := depth[n.ID]
		buckets[d] = append(buckets[d], n.ID)
		if d > maxDepth {
			maxDepth = d
		}
	}

	for d := range maxDepth + 1 {
		if len(buckets[d]) == 0 {
			continue
		}
		sort.Strings(buckets[d])
	}

	colMaxWidth := make(map[int]styling.Length, maxDepth+1)
	for d := range maxDepth + 1 {
		for _, nodeID := range buckets[d] {
			node, ok := s.nodeLookup[nodeID]
			if !ok {
				continue
			}
			width := s.calculateWidth(node.Label)
			if node.Shape == NodeShapeDiamond && width < styling.Inches(3.2) {
				width = styling.Inches(3.2)
			}
			if width > colMaxWidth[d] {
				colMaxWidth[d] = width
			}
		}
	}

	colX := make(map[int]styling.Length, maxDepth+1)
	nextX := s.layout.gridStartX
	for d := range maxDepth + 1 {
		colX[d] = nextX
		colWidth := colMaxWidth[d]
		if colWidth == 0 {
			colWidth = s.layout.baseNodeWidth
		}
		nextX += colWidth + s.layout.hSpacing
	}

	for d := range maxDepth + 1 {
		for row, nodeID := range buckets[d] {
			node, ok := s.nodeLookup[nodeID]
			if !ok {
				continue
			}
			width := s.calculateWidth(node.Label)
			x := colX[d]
			y := s.layout.gridStartY + (stylingLengthFromInt(row) * s.layout.vSpacing)
			s.addNodeShape(node, x, y, width)
		}
	}
	return len(s.nodePositions) > 0
}

func stylingLengthFromInt(v int) styling.Length {
	return styling.Length(v)
}
