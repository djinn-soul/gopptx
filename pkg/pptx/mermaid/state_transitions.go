package mermaid

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func stateTransitionEndpoints(
	fromPos struct{ x, y styling.Length },
	toPos struct{ x, y styling.Length },
	fromWidth styling.Length,
	toWidth styling.Length,
	fromHeight styling.Length,
	toHeight styling.Length,
) stateTransitionGeometry {
	switch {
	case fromPos.x < toPos.x:
		return stateTransitionGeometry{
			startX:    fromPos.x + fromWidth,
			startY:    fromPos.y + fromHeight/2,
			endX:      toPos.x,
			endY:      toPos.y + toHeight/2,
			startSite: shapes.ConnectionSiteRight,
			endSite:   shapes.ConnectionSiteLeft,
		}
	case fromPos.x > toPos.x:
		return stateTransitionGeometry{
			startX:    fromPos.x,
			startY:    fromPos.y + fromHeight/2,
			endX:      toPos.x + toWidth,
			endY:      toPos.y + toHeight/2,
			startSite: shapes.ConnectionSiteLeft,
			endSite:   shapes.ConnectionSiteRight,
		}
	case fromPos.y < toPos.y:
		return stateTransitionGeometry{
			startX:    fromPos.x + fromWidth/2,
			startY:    fromPos.y + fromHeight,
			endX:      toPos.x + toWidth/2,
			endY:      toPos.y,
			startSite: shapes.ConnectionSiteBottom,
			endSite:   shapes.ConnectionSiteTop,
		}
	default:
		return stateTransitionGeometry{
			startX:    fromPos.x + fromWidth/2,
			startY:    fromPos.y,
			endX:      toPos.x + toWidth/2,
			endY:      toPos.y + toHeight,
			startSite: shapes.ConnectionSiteTop,
			endSite:   shapes.ConnectionSiteBottom,
		}
	}
}

func stateTransitionLabelShape(
	label string,
	startX styling.Length,
	startY styling.Length,
	endX styling.Length,
	endY styling.Length,
	theme Theme,
) shapes.Shape {
	labelWidth := styling.Inches(1.0)
	labelHeight := styling.Inches(0.4)
	midX := (startX + endX) / 2
	midY := (startY + endY) / 2
	return shapes.NewShape(
		shapes.ShapeTypeRectangle,
		midX-labelWidth/2,
		midY-labelHeight/2,
		labelWidth,
		labelHeight,
	).WithText(label).
		WithFill(shapes.NewShapeFill(theme.Background)).
		WithAutoFit(shapes.TextAutoFitNormal).
		WithTextMargins(styling.Inches(0.05), styling.Inches(0.02), styling.Inches(0.05), styling.Inches(0.02))
}
