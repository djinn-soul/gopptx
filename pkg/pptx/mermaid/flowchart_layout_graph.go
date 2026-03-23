package mermaid

import (
	"sort"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func (s *flowchartRenderState) layoutByConnections(connections []FlowConnection) bool {
	if len(s.nodes) == 0 || len(connections) == 0 {
		return false
	}

	depth := s.graphDepthMap(connections)
	buckets, maxDepth := s.graphDepthBuckets(depth)
	sortBuckets(buckets, maxDepth)
	colMaxWidth := s.graphColumnMaxWidths(buckets, maxDepth)
	colX := s.graphColumnPositions(colMaxWidth, maxDepth)
	s.placeGraphBucketedNodes(buckets, colX, maxDepth)
	return len(s.nodePositions) > 0
}

func (s *flowchartRenderState) graphDepthMap(connections []FlowConnection) map[string]int {
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
	return depth
}

func (s *flowchartRenderState) graphDepthBuckets(depth map[string]int) (map[int][]string, int) {
	buckets := make(map[int][]string, len(s.nodes))
	maxDepth := 0
	for _, n := range s.nodes {
		d := depth[n.ID]
		buckets[d] = append(buckets[d], n.ID)
		if d > maxDepth {
			maxDepth = d
		}
	}
	return buckets, maxDepth
}

func sortBuckets(buckets map[int][]string, maxDepth int) {
	for d := range maxDepth + 1 {
		if len(buckets[d]) == 0 {
			continue
		}
		sort.Strings(buckets[d])
	}
}

func (s *flowchartRenderState) graphColumnMaxWidths(
	buckets map[int][]string,
	maxDepth int,
) map[int]styling.Length {
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
	return colMaxWidth
}

func (s *flowchartRenderState) graphColumnPositions(
	colMaxWidth map[int]styling.Length,
	maxDepth int,
) map[int]styling.Length {
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
	return colX
}

func (s *flowchartRenderState) placeGraphBucketedNodes(
	buckets map[int][]string,
	colX map[int]styling.Length,
	maxDepth int,
) {
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
}

func stylingLengthFromInt(v int) styling.Length {
	return styling.Length(v)
}
