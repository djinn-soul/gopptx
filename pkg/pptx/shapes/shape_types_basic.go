package shapes

// Basic geometry shape type constants (ECMA-376 ST_ShapeType).
const (
	// ShapeTypeRectangle renders a rectangle shape.
	ShapeTypeRectangle = "rect"
	// ShapeTypeRoundedRectangle renders a rounded rectangle.
	ShapeTypeRoundedRectangle = "roundRect"
	// ShapeTypeEllipse renders an ellipse shape.
	ShapeTypeEllipse = "ellipse"
	// ShapeTypeTriangle renders a triangle shape.
	ShapeTypeTriangle = "triangle"
	// ShapeTypeRightTriangle renders a right triangle shape.
	ShapeTypeRightTriangle = "rtTriangle"
	// ShapeTypeDiamond renders a diamond shape.
	ShapeTypeDiamond = "diamond"
	// ShapeTypePentagon renders a pentagon shape.
	ShapeTypePentagon = "pentagon"
	// ShapeTypeHexagon renders a hexagon shape.
	ShapeTypeHexagon = "hexagon"
	// ShapeTypeParallelogram renders a parallelogram shape.
	ShapeTypeParallelogram = "parallelogram"
	// ShapeTypeCloud renders a cloud shape.
	ShapeTypeCloud = "cloud"
	// ShapeTypeHeart renders a heart shape.
	ShapeTypeHeart = "heart"

	ShapeTypeOctagon = "octagon"
	// ShapeTypeDecagon renders a decagon.
	ShapeTypeDecagon = "decagon"
	// ShapeTypeDodecagon renders a dodecagon.
	ShapeTypeDodecagon = "dodecagon"
	// ShapeTypeHeptagon renders a heptagon.
	ShapeTypeHeptagon = "heptagon"
	// ShapeTypeTrapezoid renders a trapezoid.
	ShapeTypeTrapezoid = "trapezoid"
	// ShapeTypeNonIsoscelesTrapezoid renders a non-isosceles trapezoid.
	ShapeTypeNonIsoscelesTrapezoid = "nonIsoscelesTrapezoid"
	// ShapeTypeIsoscelesTrapezoid renders an isosceles trapezoid.
	ShapeTypeIsoscelesTrapezoid = "isoTrapezoid"
	// ShapeTypeCube renders a cube.
	ShapeTypeCube = "cube"
	// ShapeTypeCan renders a cylinder/can.
	ShapeTypeCan = "can"
	// ShapeTypeCone renders a cone.
	ShapeTypeCone = "cone"
	// ShapeTypeCylinder renders a cylinder.
	ShapeTypeCylinder = "cylinder"
	// ShapeTypeFoldedCorner renders a folded corner shape.
	ShapeTypeFoldedCorner = "foldedCorner"
	// ShapeTypeFrame renders a frame shape.
	ShapeTypeFrame = "frame"
	// ShapeTypeHalfFrame renders a half-frame shape.
	ShapeTypeHalfFrame = "halfFrame"
	// ShapeTypeCorner renders a corner shape.
	ShapeTypeCorner = "corner"
	// ShapeTypeLShape renders an L-shape.
	ShapeTypeLShape = "lShape"
	// ShapeTypeDiagStripe renders a diagonal stripe.
	ShapeTypeDiagStripe = "diagStripe"
	// ShapeTypePlus renders a plus/cross shape.
	ShapeTypePlus = "plus"
	// ShapeTypeDonut renders a donut shape.
	ShapeTypeDonut = "donut"
	// ShapeTypeNoSmoking renders a no-smoking sign.
	ShapeTypeNoSmoking = "noSmoking"
	// ShapeTypeBevel renders a bevel shape.
	ShapeTypeBevel = "bevel"
	// ShapeTypeFunnel renders a funnel shape.
	ShapeTypeFunnel = "funnel"
	// ShapeTypeTeardrop renders a teardrop shape.
	ShapeTypeTeardrop = "teardrop"
	// ShapeTypePie renders a pie shape.
	ShapeTypePie = "pie"
	// ShapeTypePieWedge renders a pie-wedge shape.
	ShapeTypePieWedge = "pieWedge"
	// ShapeTypeChord renders a chord shape.
	ShapeTypeChord = "chord"
	// ShapeTypeArc renders an arc shape.
	ShapeTypeArc = "arc"
	// ShapeTypeBlockArc renders a block arc shape.
	ShapeTypeBlockArc = "blockArc"
	// ShapeTypeSmileyFace renders a smiley face.
	ShapeTypeSmileyFace = "smileyFace"
	// ShapeTypeSun renders a sun shape.
	ShapeTypeSun = "sun"
	// ShapeTypeMoon renders a moon shape.
	ShapeTypeMoon = "moon"
	// ShapeTypeLightningBolt renders a lightning bolt.
	ShapeTypeLightningBolt = "lightningBolt"
	// ShapeTypeMusicNote renders a music note.
	ShapeTypeMusicNote = "musicNote"
	// ShapeTypePictureFrame renders a picture frame.
	ShapeTypePictureFrame = "frame"
	// ShapeTypePlaque renders a plaque shape.
	ShapeTypePlaque = "plaque"
	// ShapeTypeHomePlate renders a home plate shape.
	ShapeTypeHomePlate = "homePlate"
	// ShapeTypeChevron renders a chevron shape.
	ShapeTypeChevron = "chevron"
)

func init() {
	for _, t := range []string{
		ShapeTypeRectangle, ShapeTypeRoundedRectangle, ShapeTypeEllipse,
		ShapeTypeTriangle, ShapeTypeRightTriangle, ShapeTypeDiamond,
		ShapeTypePentagon, ShapeTypeHexagon, ShapeTypeParallelogram,
		ShapeTypeCloud, ShapeTypeHeart,
		ShapeTypeOctagon, ShapeTypeDecagon, ShapeTypeDodecagon,
		ShapeTypeHeptagon, ShapeTypeTrapezoid, ShapeTypeNonIsoscelesTrapezoid, ShapeTypeIsoscelesTrapezoid,
		ShapeTypeCube, ShapeTypeCan, ShapeTypeCone, ShapeTypeCylinder, ShapeTypeFoldedCorner,
		ShapeTypeFrame, ShapeTypeHalfFrame, ShapeTypeCorner,
		ShapeTypeLShape, ShapeTypeDiagStripe, ShapeTypePlus,
		ShapeTypeDonut, ShapeTypeNoSmoking, ShapeTypeBevel,
		ShapeTypeFunnel, ShapeTypeTeardrop, ShapeTypePie,
		ShapeTypePieWedge, ShapeTypeChord, ShapeTypeArc,
		ShapeTypeBlockArc, ShapeTypeSmileyFace, ShapeTypeSun,
		ShapeTypeMoon, ShapeTypeLightningBolt, ShapeTypeMusicNote, ShapeTypePictureFrame, ShapeTypePlaque,
		ShapeTypeHomePlate, ShapeTypeChevron,
	} {
		registerShapeType(t)
	}

	// Aliases for basic geometry shapes.
	registerShapeAlias("rectangle", ShapeTypeRectangle)
	registerShapeAlias("roundedrectangle", ShapeTypeRoundedRectangle)
	registerShapeAlias("rounded-rectangle", ShapeTypeRoundedRectangle)
	registerShapeAlias("rounded_rectangle", ShapeTypeRoundedRectangle)
	registerShapeAlias("circle", ShapeTypeEllipse)
	registerShapeAlias("righttriangle", ShapeTypeRightTriangle)
	registerShapeAlias("right-triangle", ShapeTypeRightTriangle)
	registerShapeAlias("right_triangle", ShapeTypeRightTriangle)
	registerShapeAlias("cylinder", ShapeTypeCan)
	registerShapeAlias("cross", ShapeTypePlus)
	registerShapeAlias("smiley", ShapeTypeSmileyFace)
	registerShapeAlias("lightning", ShapeTypeLightningBolt)
}
