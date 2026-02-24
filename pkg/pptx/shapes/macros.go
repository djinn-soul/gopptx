package shapes

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	defaultBadgeWidthInches  = 1.5
	defaultBadgeHeightInches = 0.4
)

// NewRectangle creates a rectangle shape with given inch dimensions.
func NewRectangle(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRectangle, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewEllipse creates an ellipse shape with given inch dimensions.
func NewEllipse(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeEllipse, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewTextBox creates a text box shape with given text and inch dimensions.
func NewTextBox(text string, x, y, w, h float64) Shape {
	return NewRectangle(x, y, w, h).
		WithText(text)
}

// NewRoundedRectangle creates a rounded rectangle shape with given inch dimensions.
func NewRoundedRectangle(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeRoundedRectangle,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

// NewTriangle creates a triangle shape with given inch dimensions.
func NewTriangle(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeTriangle, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewRightTriangle creates a right triangle shape with given inch dimensions.
func NewRightTriangle(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRightTriangle, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewDiamond creates a diamond shape with given inch dimensions.
func NewDiamond(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeDiamond, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewPentagon creates a pentagon shape with given inch dimensions.
func NewPentagon(x, y, w, h float64) Shape {
	return NewShape(ShapeTypePentagon, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewHexagon creates a hexagon shape with given inch dimensions.
func NewHexagon(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeHexagon, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewParallelogram creates a parallelogram shape with given inch dimensions.
func NewParallelogram(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeParallelogram, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewFlowChartProcess creates a flowchart process shape with given inch dimensions.
func NewFlowChartProcess(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartProcess,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

// NewFlowChartDecision creates a flowchart decision shape with given inch dimensions.
func NewFlowChartDecision(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartDecision,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

// NewFlowChartTerminator creates a flowchart terminator shape with given inch dimensions.
func NewFlowChartTerminator(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartTerminator,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

// NewRightArrow creates a right arrow shape with given inch dimensions.
func NewRightArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRightArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewLeftArrow creates a left arrow shape with given inch dimensions.
func NewLeftArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeLeftArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewUpArrow creates an up arrow shape with given inch dimensions.
func NewUpArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeUpArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewDownArrow creates a down arrow shape with given inch dimensions.
func NewDownArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeDownArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewLeftRightArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeLeftRightArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewUpDownArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeUpDownArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewQuadArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeQuadArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewBentArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeBentArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewUturnArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeUturnArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewCircularArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeCircularArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewChevron(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeChevronArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewWedgeRectCallout(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeWedgeRectCallout,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewWedgeRRectCallout(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeWedgeRRectCallout,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewWedgeEllipseCallout(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeWedgeEllipseCallout,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewCloudCallout(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeCloudCallout, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewCloud creates a cloud shape with given inch dimensions.

func NewCloud(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeCloud, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewOctagon(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeOctagon, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewTrapezoid(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeTrapezoid, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewCube(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeCube, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewFlowChartPredefinedProcess(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartPredefinedProcess,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartInternalStorage(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartInternalStorage,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartManualInput(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartManualInput,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartManualOperation(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartManualOperation,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartConnector(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartConnector,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartOffpageConnector(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartOffpageConnector,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartPunchedCard(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartPunchedCard,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartPunchedTape(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartPunchedTape,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartSummingJunction(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartSummingJunction,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartOr(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartOr, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewFlowChartCollate(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartCollate,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartSort(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartSort, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewFlowChartExtract(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartExtract,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartMerge(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartMerge, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewFlowChartOnlineStorage(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartOnlineStorage,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartDelay(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartDelay, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewFlowChartMagneticTape(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartMagneticTape,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartMagneticDisk(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartMagneticDisk,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartMagneticDrum(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartMagneticDrum,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartDisplay(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartDisplay,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func NewFlowChartPreparation(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartPreparation,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

// NewCircle creates a circle (ellipse with equal width and height) with given diameter in inches.

func NewCircle(x, y, diameter float64) Shape {
	return NewShape(
		ShapeTypeEllipse,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(diameter),
		styling.Inches(diameter),
	)
}

func NewStar4(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar4, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

func NewStar6(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar6, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

func NewStar7(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar7, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

func NewStar8(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar8, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

func NewStar10(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar10, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

func NewStar12(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar12, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

func NewStar16(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar16, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

func NewStar24(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar24, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

func NewStar32(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar32, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

func NewRibbon(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRibbon2, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewWave(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeWave, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

func NewSeal(x, y, size float64) Shape {
	return NewShape(ShapeTypeSeal, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

func NewActionButtonHome(x, y, size float64) Shape {
	return NewShape(
		ShapeTypeActionButtonHome,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(size),
		styling.Inches(size),
	)
}

func NewActionButtonHelp(x, y, size float64) Shape {
	return NewShape(
		ShapeTypeActionButtonHelp,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(size),
		styling.Inches(size),
	)
}

func NewActionButtonInformation(x, y, size float64) Shape {
	return NewShape(
		ShapeTypeActionButtonInformation,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(size),
		styling.Inches(size),
	)
}

func NewActionButtonBack(x, y, size float64) Shape {
	return NewShape(
		ShapeTypeActionButtonBackPrevious,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(size),
		styling.Inches(size),
	)
}

func NewActionButtonForward(x, y, size float64) Shape {
	return NewShape(
		ShapeTypeActionButtonForwardNext,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(size),
		styling.Inches(size),
	)
}

// NewStar creates a 5-pointed star shape with given size in inches.

func NewStar(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar5, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

// NewHeart creates a heart shape with given size in inches.
func NewHeart(x, y, size float64) Shape {
	return NewShape(ShapeTypeHeart, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

// NewFlowChartDocument creates a flowchart document shape with given inch dimensions.
func NewFlowChartDocument(x, y, w, h float64) Shape {
	return NewShape(
		ShapeTypeFlowChartDocument,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

// NewFlowChartData creates a flowchart data shape (parallelogram) with given inch dimensions.
func NewFlowChartData(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartData, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewBadge creates a badge (rounded rectangle with text) at a default size (1.5x0.4 inches).
func NewBadge(text string, x, y float64, color string) Shape {
	if color == "" {
		color = styling.ColorMaterialGreen
	}
	return NewShape(
		ShapeTypeRoundedRectangle,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(defaultBadgeWidthInches),
		styling.Inches(defaultBadgeHeightInches),
	).
		WithFill(NewShapeFill(color)).
		WithText(text)
}
