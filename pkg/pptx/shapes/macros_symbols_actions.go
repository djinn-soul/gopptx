package shapes

func NewStar4(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeStar4, x, y, size)
}

func NewStar6(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeStar6, x, y, size)
}

func NewStar7(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeStar7, x, y, size)
}

func NewStar8(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeStar8, x, y, size)
}

func NewStar10(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeStar10, x, y, size)
}

func NewStar12(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeStar12, x, y, size)
}

func NewStar16(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeStar16, x, y, size)
}

func NewStar24(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeStar24, x, y, size)
}

func NewStar32(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeStar32, x, y, size)
}

func NewStar(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeStar5, x, y, size)
}

func NewActionButtonHome(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeActionButtonHome, x, y, size)
}

func NewActionButtonHelp(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeActionButtonHelp, x, y, size)
}

func NewActionButtonInformation(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeActionButtonInformation, x, y, size)
}

func NewActionButtonBack(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeActionButtonBackPrevious, x, y, size)
}

func NewActionButtonForward(x, y, size float64) Shape {
	return newSquareShape(ShapeTypeActionButtonForwardNext, x, y, size)
}
