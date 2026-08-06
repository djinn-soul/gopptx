package export

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

// slideReadContext carries the per-slide state the reader threads through shape
// classification: what has been read so far, and where each element sat in the
// slide's shape tree.
type slideReadContext struct {
	ed       *editor.PresentationEditor
	slideIdx int
	// shapeIndexByID maps a shape id to its one-based position in
	// SlideContent.Shapes, which connectors use to name their endpoints.
	shapeIndexByID map[int]int
	// treeOrder maps a shape id to its paint index in <p:spTree>. It is nil
	// when the tree could not be read.
	treeOrder map[int]int
	// shapeStyles is the look each shape inherits from the deck's theme through
	// its <p:style>, keyed by shape id.
	shapeStyles                     map[int]shapeThemeStyle
	chartShapeIDs, smartArtShapeIDs map[int]struct{}
	ids                             slideElementIDs
}

// slideElementIDs records the shape id of every element the reader stores
// somewhere that has no room for a paint order of its own.
type slideElementIDs struct {
	title      int
	body       int
	tables     []int
	connectors []int
}

func newSlideReadContext(
	ed *editor.PresentationEditor,
	slideIdx int,
	treeOrder map[int]int,
	shapeStyles map[int]shapeThemeStyle,
) *slideReadContext {
	return &slideReadContext{
		ed:             ed,
		slideIdx:       slideIdx,
		shapeIndexByID: make(map[int]int),
		treeOrder:      treeOrder,
		shapeStyles:    shapeStyles,
	}
}

// slideShapeStyles returns the theme styles for one slide, or nil when they
// could not be read.
func slideShapeStyles(all []map[int]shapeThemeStyle, idx int) map[int]shapeThemeStyle {
	if idx >= len(all) {
		return nil
	}
	return all[idx]
}

// slideTreeOrder returns the shape-id to paint-index map for one slide, or nil
// when the deck's shape trees could not be read.
func slideTreeOrder(zOrders []map[int]int, idx int) map[int]int {
	if idx >= len(zOrders) {
		return nil
	}
	return zOrders[idx]
}

// paintOrder converts the ids gathered while reading into the paint order the
// native renderer sorts by.
func (c *slideReadContext) paintOrder(
	slideCharts [][]parsedChart,
	slideSmartArt [][]parsedSmartArt,
	idx int,
) slidePaintOrder {
	order := newSlidePaintOrder()
	if c.treeOrder == nil {
		return order
	}
	order.known = true
	order.title = c.zOf(c.ids.title)
	order.body = c.zOf(c.ids.body)
	order.tables = c.zList(c.ids.tables)
	order.connectors = c.zList(c.ids.connectors)
	if idx < len(slideCharts) {
		for _, chart := range slideCharts[idx] {
			if z, ok := c.treeOrder[chart.ShapeID]; ok {
				order.charts[chart.Kind] = z
			}
		}
	}
	if idx < len(slideSmartArt) {
		for _, diagram := range slideSmartArt[idx] {
			order.smartArt = append(order.smartArt, c.zOf(diagram.ShapeID))
		}
	}
	return order
}

func (c *slideReadContext) zOf(shapeID int) int {
	if shapeID == 0 {
		return unknownZ
	}
	if z, ok := c.treeOrder[shapeID]; ok {
		return z
	}
	return unknownZ
}

func (c *slideReadContext) zList(shapeIDs []int) []int {
	if len(shapeIDs) == 0 {
		return nil
	}
	out := make([]int, len(shapeIDs))
	for i, id := range shapeIDs {
		out[i] = c.zOf(id)
	}
	return out
}
