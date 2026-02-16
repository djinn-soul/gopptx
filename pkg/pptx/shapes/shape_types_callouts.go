package shapes

// Callout shape type constants (ECMA-376 ST_ShapeType).
const (
	// ShapeTypeCallout1 renders a line callout 1.
	ShapeTypeCallout1 = "callout1"
	// ShapeTypeCallout2 renders a line callout 2.
	ShapeTypeCallout2 = "callout2"
	// ShapeTypeCallout3 renders a line callout 3.
	ShapeTypeCallout3 = "callout3"
	// ShapeTypeBorderCallout1 renders a line callout 1 with border.
	ShapeTypeBorderCallout1 = "borderCallout1"
	// ShapeTypeBorderCallout2 renders a line callout 2 with border.
	ShapeTypeBorderCallout2 = "borderCallout2"
	// ShapeTypeBorderCallout3 renders a line callout 3 with border.
	ShapeTypeBorderCallout3 = "borderCallout3"
	// ShapeTypeAccentCallout1 renders a line callout 1 with accent.
	ShapeTypeAccentCallout1 = "accentCallout1"
	// ShapeTypeAccentCallout2 renders a line callout 2 with accent.
	ShapeTypeAccentCallout2 = "accentCallout2"
	// ShapeTypeAccentCallout3 renders a line callout 3 with accent.
	ShapeTypeAccentCallout3 = "accentCallout3"
	// ShapeTypeAccentBorderCallout1 renders a line callout 1 with accent and border.
	ShapeTypeAccentBorderCallout1 = "accentBorderCallout1"
	// ShapeTypeAccentBorderCallout2 renders a line callout 2 with accent and border.
	ShapeTypeAccentBorderCallout2 = "accentBorderCallout2"
	// ShapeTypeAccentBorderCallout3 renders a line callout 3 with accent and border.
	ShapeTypeAccentBorderCallout3 = "accentBorderCallout3"
	// ShapeTypeWedgeRectCallout renders a rectangular callout.
	ShapeTypeWedgeRectCallout = "wedgeRectCallout"
	// ShapeTypeWedgeRRectCallout renders a rounded rectangular callout.
	ShapeTypeWedgeRRectCallout = "wedgeRRectCallout"
	// ShapeTypeWedgeEllipseCallout renders an oval callout.
	ShapeTypeWedgeEllipseCallout = "wedgeEllipseCallout"
	// ShapeTypeCloudCallout renders a cloud callout.
	ShapeTypeCloudCallout = "cloudCallout"
)

func initCalloutShapes() {
	for _, t := range []string{
		ShapeTypeCallout1, ShapeTypeCallout2, ShapeTypeCallout3,
		ShapeTypeBorderCallout1, ShapeTypeBorderCallout2, ShapeTypeBorderCallout3,
		ShapeTypeAccentCallout1, ShapeTypeAccentCallout2, ShapeTypeAccentCallout3,
		ShapeTypeAccentBorderCallout1, ShapeTypeAccentBorderCallout2, ShapeTypeAccentBorderCallout3,
		ShapeTypeWedgeRectCallout, ShapeTypeWedgeRRectCallout, ShapeTypeWedgeEllipseCallout,
		ShapeTypeCloudCallout,
	} {
		registerShapeType(t)
	}

	// Callout aliases.
	registerShapeAlias("speechbubble", ShapeTypeWedgeRectCallout)
	registerShapeAlias("speech-bubble", ShapeTypeWedgeRectCallout)
	registerShapeAlias("thoughtbubble", ShapeTypeCloudCallout)
	registerShapeAlias("thought-bubble", ShapeTypeCloudCallout)
}
