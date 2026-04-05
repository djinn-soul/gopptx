//nolint:mnd // Hierarchy layout constants are tuned for SmartArt readability in native PDF.
package export

import (
	"math"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

type smartArtHierarchyNode struct {
	Node     smartart.Node
	Parent   int
	Depth    int
	Children []int
}

func layoutSmartArtHierarchy(diagram smartart.SmartArt) ([]smartArtBox, []smartArtLink) {
	nodes, levels := buildSmartArtHierarchy(diagram.Nodes)
	if len(nodes) == 0 {
		return nil, nil
	}
	x, y, w, h := emuToPt(
		int64(diagram.X),
	), emuToPt(
		int64(diagram.Y),
	), emuToPt(
		int64(diagram.CX),
	), emuToPt(
		int64(diagram.CY),
	)
	maxCols := 1
	for _, level := range levels {
		if len(level) > maxCols {
			maxCols = len(level)
		}
	}
	rowGap := 14.0
	colGap := 12.0
	boxW := math.Max(1, math.Min(150, (w-colGap*float64(maxCols-1))/float64(maxCols)))
	boxH := math.Max(1, math.Min(44, (h-rowGap*float64(len(levels)-1))/float64(len(levels))))
	boxes := make([]smartArtBox, len(nodes))
	links := make([]smartArtLink, 0, len(nodes)-1)
	centers := make([][2]float64, len(nodes))

	for depth, level := range levels {
		totalW := float64(len(level))*boxW + float64(max(len(level)-1, 0))*colGap
		left := x + (w-totalW)/2
		top := y + float64(depth)*(boxH+rowGap)
		for idx, nodeIndex := range level {
			node := nodes[nodeIndex]
			bx := left + float64(idx)*(boxW+colGap)
			by := top
			boxes[nodeIndex] = smartArtBox{
				X:         bx,
				Y:         by,
				W:         boxW,
				H:         boxH,
				Text:      node.Node.Text,
				ShapeType: shapesShapeRectangle,
				Fill:      smartArtPalette(node.Depth),
			}
			centers[nodeIndex] = [2]float64{bx + boxW/2, by + boxH/2}
			if node.Parent >= 0 {
				parentCenter := centers[node.Parent]
				links = append(links, smartArtLink{
					StartX: parentCenter[0],
					StartY: parentCenter[1] + boxH/2,
					EndX:   bx + boxW/2,
					EndY:   by,
				})
			}
		}
	}
	return boxes, links
}

func buildSmartArtHierarchy(roots []smartart.Node) ([]smartArtHierarchyNode, [][]int) {
	if len(roots) == 0 {
		return nil, nil
	}
	nodes := make([]smartArtHierarchyNode, 0, len(roots))
	levels := make([][]int, 0, 4)
	var addNode func(parent int, depth int, node smartart.Node) int
	addNode = func(parent int, depth int, node smartart.Node) int {
		index := len(nodes)
		nodes = append(nodes, smartArtHierarchyNode{Node: node, Parent: parent, Depth: depth})
		for len(levels) <= depth {
			levels = append(levels, nil)
		}
		levels[depth] = append(levels[depth], index)
		for _, child := range node.Children {
			childIndex := addNode(index, depth+1, child)
			nodes[index].Children = append(nodes[index].Children, childIndex)
		}
		return index
	}
	for _, root := range roots {
		addNode(-1, 0, root)
	}
	return nodes, levels
}
