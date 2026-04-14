package shapes

import "github.com/djinn-soul/gopptx/pkg/pptx/styling"

func NewRectangle(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeRectangle, x, y, w, h)
}

func NewEllipse(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeEllipse, x, y, w, h)
}

func NewTextBox(text string, x, y, w, h float64) Shape {
	return NewRectangle(x, y, w, h).WithText(text)
}

func NewRoundedRectangle(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeRoundedRectangle, x, y, w, h)
}

func NewTriangle(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeTriangle, x, y, w, h)
}

func NewRightTriangle(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeRightTriangle, x, y, w, h)
}

func NewDiamond(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeDiamond, x, y, w, h)
}

func NewPentagon(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypePentagon, x, y, w, h)
}

func NewHexagon(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeHexagon, x, y, w, h)
}

func NewParallelogram(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeParallelogram, x, y, w, h)
}

func NewCloud(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeCloud, x, y, w, h)
}

func NewOctagon(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeOctagon, x, y, w, h)
}

func NewTrapezoid(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeTrapezoid, x, y, w, h)
}

func NewCube(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeCube, x, y, w, h)
}

func NewCircle(x, y, diameter float64) Shape {
	return newSquareShape(ShapeTypeEllipse, x, y, diameter)
}

func NewRibbon(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeRibbon2, x, y, w, h)
}

func NewWave(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeWave, x, y, w, h)
}

func NewSeal(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeSeal, x, y, size)
}

func NewHeart(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeHeart, x, y, size)
}

func NewBadge(text string, x, y float64, color string) Shape {
	if color == "" {
		color = styling.ColorMaterialGreen
	}
	return newShapeInches(
		ShapeTypeRoundedRectangle,
		x,
		y,
		defaultBadgeWidthInches,
		defaultBadgeHeightInches,
	).
		WithFill(NewShapeFill(color)).
		WithText(text)
}
