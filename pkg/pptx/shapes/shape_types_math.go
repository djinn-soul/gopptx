package shapes

// Math and special shape type constants (ECMA-376 ST_ShapeType).
const (
	// ShapeTypeMathPlus renders a math plus sign.
	ShapeTypeMathPlus = "mathPlus"
	// ShapeTypeMathMinus renders a math minus sign.
	ShapeTypeMathMinus = "mathMinus"
	// ShapeTypeMathMultiply renders a math multiply sign.
	ShapeTypeMathMultiply = "mathMultiply"
	// ShapeTypeMathDivide renders a math divide sign.
	ShapeTypeMathDivide = "mathDivide"
	// ShapeTypeMathEqual renders a math equal sign.
	ShapeTypeMathEqual = "mathEqual"
	// ShapeTypeMathNotEqual renders a math not-equal sign.
	ShapeTypeMathNotEqual = "mathNotEqual"
	// ShapeTypeGear6 renders a 6-tooth gear.
	ShapeTypeGear6 = "gear6"
	// ShapeTypeGear9 renders a 9-tooth gear.
	ShapeTypeGear9 = "gear9"
	// ShapeTypeChartPlus renders a chart plus shape.
	ShapeTypeChartPlus = "chartPlus"
	// ShapeTypeChartStar renders a chart star shape.
	ShapeTypeChartStar = "chartStar"
	// ShapeTypeChartX renders a chart X shape.
	ShapeTypeChartX = "chartX"
	// ShapeTypeBracePair renders a brace pair shape.
	ShapeTypeBracePair = "bracePair"
	// ShapeTypeBracketPair renders a bracket pair shape.
	ShapeTypeBracketPair = "bracketPair"
	// ShapeTypeLeftBrace renders a left brace shape.
	ShapeTypeLeftBrace = "leftBrace"
	// ShapeTypeRightBrace renders a right brace shape.
	ShapeTypeRightBrace = "rightBrace"
	// ShapeTypeLeftBracket renders a left bracket shape.
	ShapeTypeLeftBracket = "leftBracket"
	// ShapeTypeRightBracket renders a right bracket shape.
	ShapeTypeRightBracket = "rightBracket"
	// ShapeTypeDoubleBrace renders a double brace shape.
	ShapeTypeDoubleBrace = "doubleBrace"
	// ShapeTypeDoubleBracket renders a double bracket shape.
	ShapeTypeDoubleBracket = "doubleBracket"
	// ShapeTypeLine renders a line.
	ShapeTypeLine = "line"
	// ShapeTypeLineInv renders an inverted line.
	ShapeTypeLineInv = "lineInv"
)

func initMathShapes() {
	for _, t := range []string{
		ShapeTypeMathPlus, ShapeTypeMathMinus,
		ShapeTypeMathMultiply, ShapeTypeMathDivide,
		ShapeTypeMathEqual, ShapeTypeMathNotEqual,
		ShapeTypeGear6, ShapeTypeGear9,
		ShapeTypeChartPlus, ShapeTypeChartStar, ShapeTypeChartX,
		ShapeTypeBracePair, ShapeTypeBracketPair,
		ShapeTypeLeftBrace, ShapeTypeRightBrace,
		ShapeTypeLeftBracket, ShapeTypeRightBracket,
		ShapeTypeDoubleBrace, ShapeTypeDoubleBracket,
		ShapeTypeLine, ShapeTypeLineInv,
	} {
		registerShapeType(t)
	}
}
