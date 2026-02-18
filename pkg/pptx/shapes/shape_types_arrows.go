package shapes

// Arrow shape type constants (ECMA-376 ST_ShapeType).
const (
	// ShapeTypeRightArrow renders a right arrow shape.
	ShapeTypeRightArrow = "rightArrow"
	// ShapeTypeLeftArrow renders a left arrow shape.
	ShapeTypeLeftArrow = "leftArrow"
	// ShapeTypeUpArrow renders an up arrow shape.
	ShapeTypeUpArrow = "upArrow"
	// ShapeTypeDownArrow renders a down arrow shape.
	ShapeTypeDownArrow = "downArrow"

	// ShapeTypeLeftRightArrow renders a left-right arrow.
	ShapeTypeLeftRightArrow = "leftRightArrow"
	// ShapeTypeUpDownArrow renders an up-down arrow.
	ShapeTypeUpDownArrow = "upDownArrow"
	// ShapeTypeQuadArrow renders a four-direction arrow.
	ShapeTypeQuadArrow = "quadArrow"
	// ShapeTypeLeftUpArrow renders a left-up arrow.
	ShapeTypeLeftUpArrow = "leftUpArrow"
	// ShapeTypeBentArrow renders a bent arrow.
	ShapeTypeBentArrow = "bentArrow"
	// ShapeTypeBentUpArrow renders a bent-up arrow.
	ShapeTypeBentUpArrow = "bentUpArrow"
	// ShapeTypeUturnArrow renders a U-turn arrow.
	ShapeTypeUturnArrow = "uturnArrow"
	// ShapeTypeCircularArrow renders a circular arrow.
	ShapeTypeCircularArrow = "circularArrow"
	// ShapeTypeNotchedRightArrow renders a notched right arrow.
	ShapeTypeNotchedRightArrow = "notchedRightArrow"
	// ShapeTypeStripedRightArrow renders a striped right arrow.
	ShapeTypeStripedRightArrow = "stripedRightArrow"
	// ShapeTypeSwooshArrow renders a swoosh arrow.
	ShapeTypeSwooshArrow = "swooshArrow"
	// ShapeTypeLeftRightUpArrow renders a left-right-up arrow.
	ShapeTypeLeftRightUpArrow = "leftRightUpArrow"
	// ShapeTypeCurvedRightArrow renders a curved right arrow.
	ShapeTypeCurvedRightArrow = "curvedRightArrow"
	// ShapeTypeCurvedLeftArrow renders a curved left arrow.
	ShapeTypeCurvedLeftArrow = "curvedLeftArrow"
	// ShapeTypeCurvedUpArrow renders a curved up arrow.
	ShapeTypeCurvedUpArrow = "curvedUpArrow"
	// ShapeTypeCurvedDownArrow renders a curved down arrow.
	ShapeTypeCurvedDownArrow = "curvedDownArrow"
	// ShapeTypeCurvedLeftRightArrow renders a curved left-right arrow.
	ShapeTypeCurvedLeftRightArrow = "curvedLeftRightArrow"
	// ShapeTypeCurvedUpDownArrow renders a curved up-down arrow.
	ShapeTypeCurvedUpDownArrow = "curvedUpDownArrow"
	// ShapeTypePentagonArrow renders a pentagon arrow.
	ShapeTypePentagonArrow = "pentArrow"
	// ShapeTypeChevronArrow renders a chevron arrow.
	ShapeTypeChevronArrow = "chevron"

	// ShapeTypeRightArrowCallout renders a right arrow callout.
	ShapeTypeRightArrowCallout = "rightArrowCallout"
	// ShapeTypeLeftArrowCallout renders a left arrow callout.
	ShapeTypeLeftArrowCallout = "leftArrowCallout"
	// ShapeTypeUpArrowCallout renders an up arrow callout.
	ShapeTypeUpArrowCallout = "upArrowCallout"
	// ShapeTypeDownArrowCallout renders a down arrow callout.
	ShapeTypeDownArrowCallout = "downArrowCallout"
	// ShapeTypeLeftRightArrowCallout renders a left-right arrow callout.
	ShapeTypeLeftRightArrowCallout = "leftRightArrowCallout"
	// ShapeTypeQuadArrowCallout renders a four-direction arrow callout.
	ShapeTypeQuadArrowCallout = "quadArrowCallout"
	// ShapeTypeUpDownArrowCallout renders an up-down arrow callout.
	ShapeTypeUpDownArrowCallout = "upDownArrowCallout"
)

func initArrowShapes() {
	for _, t := range []string{
		ShapeTypeRightArrow, ShapeTypeLeftArrow,
		ShapeTypeUpArrow, ShapeTypeDownArrow,
		ShapeTypeLeftRightArrow, ShapeTypeUpDownArrow,
		ShapeTypeQuadArrow, ShapeTypeLeftUpArrow,
		ShapeTypeBentArrow, ShapeTypeBentUpArrow,
		ShapeTypeUturnArrow, ShapeTypeCircularArrow,
		ShapeTypeNotchedRightArrow, ShapeTypeStripedRightArrow,
		ShapeTypeSwooshArrow, ShapeTypeLeftRightUpArrow,
		ShapeTypeCurvedRightArrow, ShapeTypeCurvedLeftArrow,
		ShapeTypeCurvedUpArrow, ShapeTypeCurvedDownArrow,
		ShapeTypeCurvedLeftRightArrow, ShapeTypeCurvedUpDownArrow,
		ShapeTypePentagonArrow, ShapeTypeChevronArrow,
		ShapeTypeRightArrowCallout, ShapeTypeLeftArrowCallout,
		ShapeTypeUpArrowCallout, ShapeTypeDownArrowCallout,
		ShapeTypeLeftRightArrowCallout, ShapeTypeQuadArrowCallout,
		ShapeTypeUpDownArrowCallout,
	} {
		registerShapeType(t)
	}

	// Arrow aliases.
	registerShapeAlias("rightarrow", ShapeTypeRightArrow)
	registerShapeAlias("right-arrow", ShapeTypeRightArrow)
	registerShapeAlias("right_arrow", ShapeTypeRightArrow)
	registerShapeAlias("leftarrow", ShapeTypeLeftArrow)
	registerShapeAlias("left-arrow", ShapeTypeLeftArrow)
	registerShapeAlias("left_arrow", ShapeTypeLeftArrow)
	registerShapeAlias("uparrow", ShapeTypeUpArrow)
	registerShapeAlias("up-arrow", ShapeTypeUpArrow)
	registerShapeAlias("up_arrow", ShapeTypeUpArrow)
	registerShapeAlias("downarrow", ShapeTypeDownArrow)
	registerShapeAlias("down-arrow", ShapeTypeDownArrow)
	registerShapeAlias("down_arrow", ShapeTypeDownArrow)
}
