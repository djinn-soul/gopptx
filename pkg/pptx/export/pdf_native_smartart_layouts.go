//nolint:mnd // SmartArt layout presets require fixed geometric constants per template family.
package export

import (
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

func layoutSmartArtLinear(diagram smartart.SmartArt) ([]smartArtBox, []smartArtLink) {
	nodes := flattenSmartArtNodes(diagram.Nodes)
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
	layoutURI := strings.ToLower(diagram.Layout.LayoutURI())
	vertical := strings.Contains(layoutURI, "vlist") || strings.Contains(layoutURI, "blocklist")
	gap := 10.0
	boxes := make([]smartArtBox, 0, len(nodes))
	links := make([]smartArtLink, 0, max(len(nodes)-1, 0))
	if vertical {
		boxH := math.Max(1, math.Min(46, (h-gap*float64(len(nodes)-1))/float64(len(nodes))))
		boxW := math.Min(w, 230.0)
		left := x + (w-boxW)/2
		for i, node := range nodes {
			top := y + float64(i)*(boxH+gap)
			boxes = append(
				boxes,
				smartArtBox{
					X:         left,
					Y:         top,
					W:         boxW,
					H:         boxH,
					Text:      node.Text,
					ShapeType: shapeTypeForLinear(layoutURI),
					Fill:      smartArtPalette(i),
				},
			)
			if i > 0 {
				links = append(
					links,
					smartArtLink{StartX: left + boxW/2, StartY: top - gap, EndX: left + boxW/2, EndY: top},
				)
			}
		}
		return boxes, links
	}

	boxW := math.Min(150, (w-gap*float64(len(nodes)-1))/float64(len(nodes)))
	boxH := math.Min(h, 52.0)
	top := y + (h-boxH)/2
	for i, node := range nodes {
		left := x + float64(i)*(boxW+gap)
		boxes = append(
			boxes,
			smartArtBox{
				X:         left,
				Y:         top,
				W:         boxW,
				H:         boxH,
				Text:      node.Text,
				ShapeType: shapeTypeForLinear(layoutURI),
				Fill:      smartArtPalette(i),
			},
		)
		if i > 0 {
			links = append(
				links,
				smartArtLink{StartX: left - gap, StartY: top + boxH/2, EndX: left, EndY: top + boxH/2},
			)
		}
	}
	return boxes, links
}

func layoutSmartArtGrid(diagram smartart.SmartArt) ([]smartArtBox, []smartArtLink) {
	nodes := flattenSmartArtNodes(diagram.Nodes)
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
	cols := int(math.Ceil(math.Sqrt(float64(len(nodes)))))
	rows := int(math.Ceil(float64(len(nodes)) / float64(cols)))
	gap := 10.0
	boxW := (w - gap*float64(max(cols-1, 0))) / float64(cols)
	boxH := (h - gap*float64(max(rows-1, 0))) / float64(rows)
	boxes := make([]smartArtBox, 0, len(nodes))
	for i, node := range nodes {
		col := i % cols
		row := i / cols
		boxes = append(boxes, smartArtBox{
			X:         x + float64(col)*(boxW+gap),
			Y:         y + float64(row)*(boxH+gap),
			W:         boxW,
			H:         boxH,
			Text:      node.Text,
			ShapeType: shapesShapeRectangle,
			Fill:      smartArtPalette(row + col),
		})
	}
	return boxes, nil
}

func layoutSmartArtPyramid(diagram smartart.SmartArt) ([]smartArtBox, []smartArtLink) {
	nodes := flattenSmartArtNodes(diagram.Nodes)
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
	inverted := strings.Contains(strings.ToLower(diagram.Layout.LayoutURI()), "inverted")
	boxH := h / float64(len(nodes))
	boxes := make([]smartArtBox, 0, len(nodes))
	for i, node := range nodes {
		ratio := float64(i+1) / float64(len(nodes))
		if inverted {
			ratio = 1 - (float64(i) / float64(len(nodes)))
		}
		boxW := w * (0.35 + 0.6*ratio)
		left := x + (w-boxW)/2
		top := y + float64(i)*boxH
		boxes = append(
			boxes,
			smartArtBox{
				X:         left,
				Y:         top,
				W:         boxW,
				H:         boxH - 6,
				Text:      node.Text,
				ShapeType: shapesShapeRectangle,
				Fill:      smartArtPalette(i),
			},
		)
	}
	return boxes, nil
}

func layoutSmartArtRadial(diagram smartart.SmartArt) ([]smartArtBox, []smartArtLink) {
	nodes := flattenSmartArtNodes(diagram.Nodes)
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
	cx, cy := x+w/2, y+h/2
	boxW := math.Min(110, w/3.2)
	boxH := math.Min(42, h/4)
	radius := math.Min(w, h)*0.32 - boxH/2
	boxes := make([]smartArtBox, 0, len(nodes))
	links := make([]smartArtLink, 0, len(nodes))
	if len(nodes) == 1 {
		return []smartArtBox{
			{
				X:         cx - boxW/2,
				Y:         cy - boxH/2,
				W:         boxW,
				H:         boxH,
				Text:      nodes[0].Text,
				ShapeType: shapesShapeEllipse,
				Fill:      smartArtPalette(0),
			},
		}, nil
	}
	center := nodes[0]
	boxes = append(
		boxes,
		smartArtBox{
			X:         cx - boxW/2,
			Y:         cy - boxH/2,
			W:         boxW,
			H:         boxH,
			Text:      center.Text,
			ShapeType: shapesShapeEllipse,
			Fill:      smartArtPalette(0),
		},
	)
	orbit := nodes[1:]
	for i, node := range orbit {
		angle := (2 * math.Pi * float64(i) / float64(len(orbit))) - math.Pi/2
		bx := cx + math.Cos(angle)*radius - boxW/2
		by := cy + math.Sin(angle)*radius - boxH/2
		boxes = append(
			boxes,
			smartArtBox{
				X:         bx,
				Y:         by,
				W:         boxW,
				H:         boxH,
				Text:      node.Text,
				ShapeType: shapesShapeEllipse,
				Fill:      smartArtPalette(i + 1),
			},
		)
		links = append(links, smartArtLink{StartX: cx, StartY: cy, EndX: bx + boxW/2, EndY: by + boxH/2})
	}
	return boxes, links
}

const (
	shapesShapeRectangle = "rect"
	shapesShapeEllipse   = "ellipse"
)

func shapeTypeForLinear(layoutURI string) string {
	if strings.Contains(layoutURI, "process") || strings.Contains(layoutURI, "picturestrips") {
		return shapesShapeRoundedRect
	}
	return shapesShapeRectangle
}

const shapesShapeRoundedRect = "roundRect"
