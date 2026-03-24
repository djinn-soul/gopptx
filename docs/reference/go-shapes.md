# Go Shapes and Connectors Reference

This page documents the common shape constructors, shape constants, and connector helpers exposed by `pkg/pptx`.

Primary source files:

- `pkg/pptx/shape.go`
- `pkg/pptx/shape_types.go`
- `pkg/pptx/connector.go`
- `pkg/pptx/shapes/macros.go`
- `pkg/pptx/shapes/connector.go`
- `pkg/pptx/shapes/shape_fill_rich.go`

## Common shape constructors

- `NewShape(shapeType string, x, y, cx, cy Length) Shape`
- `NewShapeFill(color string) ShapeFill`
- `NewShapeLine(color string, width Length) ShapeLine`
- `NewShapeGradientStop(positionPct int, color string) ShapeGradientStop`
- `NewShapeGradientFill(gradientType string, stops []ShapeGradientStop) ShapeGradientFill`
- `NewTextFrame() TextFrame`
- `NewSolidFill(color string) *RichShapeFill`
- `NewNoFill() *RichShapeFill`
- `NewPatternFill(pattern PatternType) *RichShapeFill`
- `NewRichShapeLine(color string, width Length) *RichShapeLine`
- `NewOuterShadow(color string) *RichShapeShadow`
- `NewInnerShadow(color string) *RichShapeShadow`
- `NewPerspectiveShadow(color string) *RichShapeShadow`
- `NewGroupShape(x, y, w, h Length) GroupShape`
- `NewGroupShapeBounds(shapesList []Shape) GroupShape`
- `NewFreeform(points []FreeformPoint) Freeform`
- `NewFreeformCoords(xCoords, yCoords []int64) (Freeform, error)`
- `NewFreeformInches(points [][2]float64) (Freeform, error)`
- `NewFreeformClosed(points []FreeformPoint) Freeform`
- `NewFreeformOpen(points []FreeformPoint) Freeform`
- `(*RichShapeFill).WithTransparency(value float64) *RichShapeFill`
- `(*TextFrame).WithRotation(degrees float64) TextFrame`

## Inches-based shape helpers

- `NewRectangle(x, y, w, h float64) Shape`
- `NewEllipse(x, y, w, h float64) Shape`
- `NewTextBox(text string, x, y, w, h float64) Shape`
- `NewRoundedRectangle(x, y, w, h float64) Shape`
- `NewTriangle(x, y, w, h float64) Shape`
- `NewRightTriangle(x, y, w, h float64) Shape`
- `NewDiamond(x, y, w, h float64) Shape`
- `NewPentagon(x, y, w, h float64) Shape`
- `NewHexagon(x, y, w, h float64) Shape`
- `NewParallelogram(x, y, w, h float64) Shape`
- `NewFlowChartProcess(x, y, w, h float64) Shape`
- `NewFlowChartDecision(x, y, w, h float64) Shape`
- `NewFlowChartTerminator(x, y, w, h float64) Shape`
- `NewRightArrow(x, y, w, h float64) Shape`
- `NewLeftArrow(x, y, w, h float64) Shape`
- `NewUpArrow(x, y, w, h float64) Shape`
- `NewDownArrow(x, y, w, h float64) Shape`
- `NewLeftRightArrow(x, y, w, h float64) Shape`
- `NewUpDownArrow(x, y, w, h float64) Shape`
- `NewQuadArrow(x, y, w, h float64) Shape`
- `NewBentArrow(x, y, w, h float64) Shape`
- `NewUturnArrow(x, y, w, h float64) Shape`
- `NewCircularArrow(x, y, w, h float64) Shape`
- `NewChevron(x, y, w, h float64) Shape`
- `NewCloud(x, y, w, h float64) Shape`
- `NewCloudCallout(x, y, w, h float64) Shape`
- `NewCircle(x, y, diameter float64) Shape`
- `NewHeart(x, y, size float64) Shape`
- `NewStar(x, y, size float64) Shape`
- `NewStar4(x, y, size float64) Shape`
- `NewStar6(x, y, size float64) Shape`
- `NewStar7(x, y, size float64) Shape`
- `NewStar8(x, y, size float64) Shape`
- `NewStar10(x, y, size float64) Shape`
- `NewStar12(x, y, size float64) Shape`
- `NewStar16(x, y, size float64) Shape`
- `NewStar24(x, y, size float64) Shape`
- `NewStar32(x, y, size float64) Shape`
- `NewRibbon(x, y, w, h float64) Shape`
- `NewWave(x, y, w, h float64) Shape`
- `NewSeal(x, y, size float64) Shape`
- `NewActionButtonHome(x, y, size float64) Shape`
- `NewActionButtonHelp(x, y, size float64) Shape`
- `NewActionButtonInformation(x, y, size float64) Shape`
- `NewActionButtonBack(x, y, size float64) Shape`
- `NewActionButtonForward(x, y, size float64) Shape`
- `NewFlowChartDocument(x, y, w, h float64) Shape`
- `NewFlowChartData(x, y, w, h float64) Shape`
- `NewOctagon(x, y, w, h float64) Shape`
- `NewTrapezoid(x, y, w, h float64) Shape`

## Common constants

- `ShapeTypeRectangle`
- `ShapeTypeRoundedRectangle`
- `ShapeTypeEllipse`
- `ShapeTypeTriangle`
- `ShapeTypeRightTriangle`
- `ShapeTypeDiamond`
- `ShapeTypePentagon`
- `ShapeTypeHexagon`
- `ShapeTypeParallelogram`
- `ShapeTypeFlowChartProcess`
- `ShapeTypeFlowChartDecision`
- `ShapeTypeFlowChartTerminator`
- `ShapeTypeRightArrow`
- `ShapeTypeLeftArrow`
- `ShapeTypeUpArrow`
- `ShapeTypeDownArrow`
- `ShapeTypeCloud`
- `ShapeTypeStar5`
- `ShapeTypeHeart`
- `ShapeTypeFlowChartDocument`
- `ShapeTypeFlowChartData`
- `ShapeTypeGroup`
- `ShapeGradientTypeLinear`
- `ShapeGradientTypeRadial`
- `ShapeGradientTypeRectangular`
- `ShapeGradientTypePath`
- `TextAnchorTop`
- `TextAnchorMiddle`
- `TextAnchorBottom`
- `TextWrapNone`
- `TextWrapSquare`
- `TextAutoFitNone`
- `TextAutoFitShape`
- `TextAutoFitNormal`
- `ConnectorTypeStraight`
- `ConnectorTypeElbow`
- `ConnectorTypeCurved`
- `ArrowTypeNone`
- `ArrowTypeTriangle`
- `ArrowTypeStealth`
- `ArrowTypeDiamond`
- `ArrowTypeOval`
- `ArrowTypeOpen`
- `ArrowSizeSmall`
- `ArrowSizeMedium`
- `ArrowSizeLarge`
- `ConnectionSiteTop`
- `ConnectionSiteRight`
- `ConnectionSiteBottom`
- `ConnectionSiteLeft`
- `ConnectionSiteTopLeft`
- `ConnectionSiteTopRight`
- `ConnectionSiteBottomRight`
- `ConnectionSiteBottomLeft`
- `ConnectionSiteCenter`

## Connector helpers

- `NewConnector(connectorType string, startX, startY, endX, endY styling.Length) Connector`
- `NewStraightConnector(startX, startY, endX, endY styling.Length) Connector`
- `NewElbowConnector(startX, startY, endX, endY styling.Length) Connector`
- `NewCurvedConnector(startX, startY, endX, endY styling.Length) Connector`
- `ConnectStartAuto(c Connector, shapeIndex int) Connector`
- `ConnectEndAuto(c Connector, shapeIndex int) Connector`
- `AutoReroute(c Connector, shapes []Shape) Connector`

## Related pages

- [Go API Reference](go-api.md)
- [Go Slides Reference](go-slides.md)
