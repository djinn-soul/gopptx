package shapes

// Rectangle variant shape type constants (ECMA-376 ST_ShapeType).
const (
	// ShapeTypeSnip1Rect renders a rectangle with one snipped corner.
	ShapeTypeSnip1Rect = "snip1Rect"
	// ShapeTypeSnip2SameRect renders a rectangle with two same-side snipped corners.
	ShapeTypeSnip2SameRect = "snip2SameRect"
	// ShapeTypeSnip2DiagRect renders a rectangle with two diagonal snipped corners.
	ShapeTypeSnip2DiagRect = "snip2DiagRect"
	// ShapeTypeRound1Rect renders a rectangle with one rounded corner.
	ShapeTypeRound1Rect = "round1Rect"
	// ShapeTypeRound2SameRect renders a rectangle with two same-side rounded corners.
	ShapeTypeRound2SameRect = "round2SameRect"
	// ShapeTypeRound2DiagRect renders a rectangle with two diagonal rounded corners.
	ShapeTypeRound2DiagRect = "round2DiagRect"
	// ShapeTypeSnipRoundRect renders a rectangle with one snipped and one rounded corner.
	ShapeTypeSnipRoundRect = "snipRoundRect"
	// ShapeTypePlaqueTabs renders a plaque tabs shape.
	ShapeTypePlaqueTabs = "plaqueTabs"
	// ShapeTypeSquareTabs renders a square tabs shape.
	ShapeTypeSquareTabs = "squareTabs"
	// ShapeTypeCornerTabs renders a corner tabs shape.
	ShapeTypeCornerTabs = "cornerTabs"
)

func init() {
	for _, t := range []string{
		ShapeTypeSnip1Rect, ShapeTypeSnip2SameRect, ShapeTypeSnip2DiagRect,
		ShapeTypeRound1Rect, ShapeTypeRound2SameRect, ShapeTypeRound2DiagRect,
		ShapeTypeSnipRoundRect, ShapeTypePlaqueTabs, ShapeTypeSquareTabs,
		ShapeTypeCornerTabs,
	} {
		registerShapeType(t)
	}
}
