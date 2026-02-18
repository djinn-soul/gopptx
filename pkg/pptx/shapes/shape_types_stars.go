package shapes

// Star, banner, and scroll shape type constants (ECMA-376 ST_ShapeType).
const (
	// ShapeTypeStar5 renders a 5-pointed star.
	ShapeTypeStar5 = "star5"

	// ShapeTypeStar4 renders a 4-pointed star.
	ShapeTypeStar4 = "star4"
	// ShapeTypeStar6 renders a 6-pointed star.
	ShapeTypeStar6 = "star6"
	// ShapeTypeStar7 renders a 7-pointed star.
	ShapeTypeStar7 = "star7"
	// ShapeTypeStar8 renders an 8-pointed star.
	ShapeTypeStar8 = "star8"
	// ShapeTypeStar10 renders a 10-pointed star.
	ShapeTypeStar10 = "star10"
	// ShapeTypeStar12 renders a 12-pointed star.
	ShapeTypeStar12 = "star12"
	// ShapeTypeStar16 renders a 16-pointed star.
	ShapeTypeStar16 = "star16"
	// ShapeTypeStar24 renders a 24-pointed star.
	ShapeTypeStar24 = "star24"
	// ShapeTypeStar32 renders a 32-pointed star.
	ShapeTypeStar32 = "star32"
	// ShapeTypeIrregularSeal1 renders an explosion/burst shape 1.
	ShapeTypeIrregularSeal1 = "irregularSeal1"
	// ShapeTypeIrregularSeal2 renders an explosion/burst shape 2.
	ShapeTypeIrregularSeal2 = "irregularSeal2"

	// ShapeTypeRibbon renders a ribbon banner.
	ShapeTypeRibbon = "ribbon"
	// ShapeTypeRibbon2 renders a ribbon banner (variant 2).
	ShapeTypeRibbon2 = "ribbon2"
	// ShapeTypeEllipseRibbon renders an ellipse ribbon.
	ShapeTypeEllipseRibbon = "ellipseRibbon"
	// ShapeTypeEllipseRibbon2 renders an ellipse ribbon (variant 2).
	ShapeTypeEllipseRibbon2 = "ellipseRibbon2"

	// ShapeTypeWave renders a wave shape.
	ShapeTypeWave = "wave"
	// ShapeTypeDoubleWave renders a double wave shape.
	ShapeTypeDoubleWave = "doubleWave"
	// ShapeTypeVerticalScroll renders a vertical scroll.
	ShapeTypeVerticalScroll = "verticalScroll"
	// ShapeTypeHorizontalScroll renders a horizontal scroll.
	ShapeTypeHorizontalScroll = "horizontalScroll"

	// ShapeTypeSeal renders a seal shape.
	ShapeTypeSeal = "seal"
	// ShapeTypeSeal4 renders a 4-pointed seal.
	ShapeTypeSeal4 = "seal4"
	// ShapeTypeSeal8 renders an 8-pointed seal.
	ShapeTypeSeal8 = "seal8"
	// ShapeTypeSeal16 renders a 16-pointed seal.
	ShapeTypeSeal16 = "seal16"
	// ShapeTypeSeal32 renders a 32-pointed seal.
	ShapeTypeSeal32 = "seal32"
)

func initStarShapes() {
	for _, t := range []string{
		ShapeTypeStar5, ShapeTypeStar4, ShapeTypeStar6,
		ShapeTypeStar7, ShapeTypeStar8, ShapeTypeStar10,
		ShapeTypeStar12, ShapeTypeStar16, ShapeTypeStar24,
		ShapeTypeStar32, ShapeTypeIrregularSeal1, ShapeTypeIrregularSeal2,
		ShapeTypeRibbon, ShapeTypeRibbon2,
		ShapeTypeEllipseRibbon, ShapeTypeEllipseRibbon2,
		ShapeTypeWave, ShapeTypeDoubleWave,
		ShapeTypeVerticalScroll, ShapeTypeHorizontalScroll,
		ShapeTypeSeal, ShapeTypeSeal4, ShapeTypeSeal8,
		ShapeTypeSeal16, ShapeTypeSeal32,
	} {
		registerShapeType(t)
	}

	// Star aliases.
	registerShapeAlias("star", ShapeTypeStar5)
	registerShapeAlias("explosion", ShapeTypeIrregularSeal1)
	registerShapeAlias("burst", ShapeTypeIrregularSeal1)
}
