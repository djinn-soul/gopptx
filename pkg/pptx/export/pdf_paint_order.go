package export

import (
	"fmt"
	"sort"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// PowerPoint paints one slide in a single pass over <p:spTree>: the first child
// is furthest back, the last furthest forward, whatever kind of element it is.
// The native renderer used to paint in fixed layers instead — placeholders,
// then pictures and shapes, then SmartArt, then charts, then tables — so a
// table always covered a shape the deck had put in front of it, and a chart
// always covered the table.
//
// slidePaintOrder carries the tree position of the elements that have nowhere
// to record one themselves. Pictures and shapes keep theirs in Image.ZOrder and
// Shape.ZOrder; everything else is keyed here by its position in the slide.
type slidePaintOrder struct {
	// known reports that the slide's shape tree was read, so Image.ZOrder and
	// Shape.ZOrder hold real tree positions rather than the zero value every
	// in-memory element carries.
	known bool
	// title and body are the placeholders the reader consumes into
	// SlideContent.Title and .Bullets.
	title int
	body  int
	// tables[0] is SlideContent.Table; tables[i] is Tables[i-1].
	tables []int
	// charts is keyed by the reader's chart kind, because each kind lands in
	// its own SlideContent field.
	charts     map[string]int
	smartArt   []int
	connectors []int
}

// unknownZ marks an element whose tree position was not recovered — an
// in-memory slide, or a part the z-order reader could not match.
const unknownZ = -1

// unknownZBase pushes elements with no recovered tree position above every
// element that has one, so a slide built in memory keeps the layering the
// renderer has always given it.
const unknownZBase = 1 << 20

func newSlidePaintOrder() slidePaintOrder {
	return slidePaintOrder{title: unknownZ, body: unknownZ, charts: map[string]int{}}
}

func (o slidePaintOrder) chartZ(kind string) int {
	if o.charts == nil {
		return unknownZ
	}
	if z, ok := o.charts[kind]; ok {
		return z
	}
	return unknownZ
}

func (o slidePaintOrder) tableZ(index int) int {
	if index < len(o.tables) {
		return o.tables[index]
	}
	return unknownZ
}

func (o slidePaintOrder) smartArtZ(index int) int {
	if index < len(o.smartArt) {
		return o.smartArt[index]
	}
	return unknownZ
}

func (o slidePaintOrder) connectorZ(index int) int {
	if index < len(o.connectors) {
		return o.connectors[index]
	}
	return unknownZ
}

// paintItem is one thing to draw, with the sort key that decides when.
type paintItem struct {
	key  int
	draw func() error
}

// paintList collects a slide's elements and hands them back in paint order.
type paintList struct {
	items []paintItem
	// fallback counts elements with no recovered tree position, so they keep
	// the order they were added in.
	fallback int
}

// add queues one element. z is its shape-tree index, or unknownZ.
func (l *paintList) add(z int, draw func() error) {
	key := z
	if z == unknownZ {
		key = unknownZBase + l.fallback
	}
	l.fallback++
	l.items = append(l.items, paintItem{key: key, draw: draw})
}

// drawAll paints every queued element back to front and joins what failed.
func (l *paintList) drawAll() []error {
	sort.SliceStable(l.items, func(a, b int) bool { return l.items[a].key < l.items[b].key })
	var errs []error
	for _, item := range l.items {
		if err := item.draw(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// buildSlidePaintList queues every element of one slide against its shape-tree
// position. The categories are added in the renderer's historical layer order,
// which is what decides the outcome for a slide whose order is unknown.
func buildSlidePaintList(
	pdf *gopdf.GoPdf,
	slide elements.SlideContent,
	order slidePaintOrder,
	page pageSize,
) *paintList {
	list := &paintList{}
	addPlaceholderPaintItems(list, pdf, slide, order, page)
	addShapePaintItems(list, pdf, slide, order)
	addDiagramPaintItems(list, pdf, slide, order)
	addChartPaintItems(list, pdf, slide, order)
	addTablePaintItems(list, pdf, slide, order)
	return list
}

func addPlaceholderPaintItems(
	list *paintList,
	pdf *gopdf.GoPdf,
	slide elements.SlideContent,
	order slidePaintOrder,
	page pageSize,
) {
	if slide.Title != "" {
		list.add(order.title, func() error { renderPDFTitle(pdf, slide, page); return nil })
	}
	if len(slide.Bullets) > 0 {
		list.add(order.body, func() error { renderPDFBullets(pdf, slide, page); return nil })
	}
}

func addShapePaintItems(list *paintList, pdf *gopdf.GoPdf, slide elements.SlideContent, order slidePaintOrder) {
	for i := range slide.Images {
		image := slide.Images[i]
		number := i + 1
		list.add(order.treeZ(image.ZOrder), func() error {
			if err := renderPDFImageWithEffects(pdf, image); err != nil {
				return fmt.Errorf("image %d: %w", number, err)
			}
			return nil
		})
	}
	for i := range slide.Shapes {
		shape := slide.Shapes[i]
		list.add(order.treeZ(shape.ZOrder), func() error { renderPDFShape(pdf, shape); return nil })
	}
	for i := range slide.Connectors {
		connector := slide.Connectors[i]
		list.add(order.connectorZ(i), func() error { renderPDFConnector(pdf, connector); return nil })
	}
}

func addDiagramPaintItems(list *paintList, pdf *gopdf.GoPdf, slide elements.SlideContent, order slidePaintOrder) {
	for i := range slide.SmartArtDiagrams {
		diagram := slide.SmartArtDiagrams[i]
		list.add(order.smartArtZ(i), func() error { renderPDFSmartArt(pdf, diagram); return nil })
	}
}

func addTablePaintItems(list *paintList, pdf *gopdf.GoPdf, slide elements.SlideContent, order slidePaintOrder) {
	if slide.Table != nil {
		table := *slide.Table
		list.add(order.tableZ(0), func() error { renderPDFTable(pdf, table); return nil })
	}
	for i := range slide.Tables {
		table := slide.Tables[i]
		list.add(order.tableZ(i+1), func() error { renderPDFTable(pdf, table); return nil })
	}
}

// treeZ reads the ZOrder an element stores itself. Index zero is a real tree
// position — the backmost element — so it can only be trusted on a slide whose
// tree was actually read; everywhere else it is the zero value of a field
// nobody set.
func (o slidePaintOrder) treeZ(z int) int {
	if !o.known {
		return unknownZ
	}
	return z
}
