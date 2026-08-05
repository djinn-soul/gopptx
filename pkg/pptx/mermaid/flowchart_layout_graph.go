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
	s.placeGraphBucketedNodes(buckets, s.graphRankPositions(buckets, maxDepth), maxDepth)
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

// graphRankExtent is how much room one rank needs along the layout axis: the
// widest node for a left-right graph, the tallest for a top-bottom one.
func (s *flowchartRenderState) graphRankExtent(nodeIDs []string) styling.Length {
	extent := styling.Length(0)
	for _, nodeID := range nodeIDs {
		node, ok := s.nodeLookup[nodeID]
		if !ok {
			continue
		}
		size := s.renderedNodeHeight(node.Shape)
		if s.isHorizontal {
			size = s.graphNodeWidth(node)
		}
		if size > extent {
			extent = size
		}
	}
	return extent
}

// graphNodeWidth is the width a node will actually be drawn at, diamonds
// included, so ranks are spaced by what gets painted rather than by the label.
func (s *flowchartRenderState) graphNodeWidth(node FlowNode) styling.Length {
	width := s.calculateWidth(node.Label)
	if node.Shape == NodeShapeDiamond && width < diamondMinWidth {
		width = diamondMinWidth
	}
	return width
}

// graphRankPositions places each rank along the axis the direction asks for.
//
// The rank axis used to be hard-coded to x, so `flowchart TD` was laid out
// left-to-right exactly like `flowchart LR` — the parsed Direction reached the
// connector router but never the node placement.
func (s *flowchartRenderState) graphRankPositions(
	buckets map[int][]string,
	maxDepth int,
) map[int]styling.Length {
	extents := make([]styling.Length, maxDepth+1)
	for d := range maxDepth + 1 {
		extent := s.graphRankExtent(buckets[d])
		if extent == 0 {
			extent = s.layout.baseNodeWidth
			if !s.isHorizontal {
				extent = s.layout.nodeHeight
			}
		}
		extents[d] = extent
	}

	gap := s.layout.hSpacing
	next := s.layout.gridStartX
	if !s.isHorizontal {
		gap = s.layout.vSpacing
		next = s.layout.gridStartY
	}

	// RL and BT are the same ranking read from the far end, so rank 0 goes last.
	order := make([]int, 0, maxDepth+1)
	for d := range maxDepth + 1 {
		order = append(order, d)
	}
	if s.isReversed {
		for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
			order[i], order[j] = order[j], order[i]
		}
	}

	positions := make(map[int]styling.Length, maxDepth+1)
	for _, d := range order {
		positions[d] = next
		next += extents[d] + gap
	}
	return positions
}

func (s *flowchartRenderState) placeGraphBucketedNodes(
	buckets map[int][]string,
	rankPos map[int]styling.Length,
	maxDepth int,
) {
	for d := range maxDepth + 1 {
		for index, nodeID := range buckets[d] {
			node, ok := s.nodeLookup[nodeID]
			if !ok {
				continue
			}
			x, y := s.graphNodeOrigin(rankPos[d], index)
			s.addNodeShape(node, x, y, s.calculateWidth(node.Label))
		}
	}
}

// graphNodeOrigin turns a rank offset and the node's position within that rank
// into a point, swapping the axes for a top-bottom graph.
func (s *flowchartRenderState) graphNodeOrigin(rank styling.Length, index int) (styling.Length, styling.Length) {
	offset := styling.Length(index)
	if s.isHorizontal {
		return rank, s.layout.gridStartY + offset*s.layout.vSpacing
	}
	return s.layout.gridStartX + offset*s.layout.hSpacing, rank
}
