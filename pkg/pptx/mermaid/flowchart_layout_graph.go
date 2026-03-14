package mermaid

import (
	"sort"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func (s *flowchartRenderState) layoutByConnections(connections []FlowConnection) bool {
	if len(s.nodes) == 0 || len(connections) == 0 {
		return false
	}

	depth := make(map[string]int, len(s.nodes))
	for _, n := range s.nodes {
		depth[n.ID] = 0
	}

	for i := 0; i < len(s.nodes); i++ {
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

	for d := 0; d <= maxDepth; d++ {
		if len(buckets[d]) == 0 {
			continue
		}
		sort.Strings(buckets[d])
		for row, nodeID := range buckets[d] {
			node, ok := s.nodeLookup[nodeID]
			if !ok {
				continue
			}
			width := s.calculateWidth(node.Label)
			x := s.layout.gridStartX + (stylingLengthFromInt(d) * s.layout.hSpacing)
			y := s.layout.gridStartY + (stylingLengthFromInt(row) * s.layout.vSpacing)
			s.addNodeShape(node, x, y, width)
		}
	}
	return len(s.nodePositions) > 0
}

func stylingLengthFromInt(v int) styling.Length {
	return styling.Length(v)
}
