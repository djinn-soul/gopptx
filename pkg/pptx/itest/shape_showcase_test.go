package itest

import (
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestShapeShowcase(t *testing.T) {
	builder := pptx.NewPresentationBuilder("Shape Showcase")

	// List of all shapes we want to showcase
	shapes := []struct {
		Category string
		Types    []string
	}{
		{
			Category: "Basic Shapes",
			Types: []string{
				pptx.ShapeTypeRectangle, pptx.ShapeTypeRoundedRectangle, pptx.ShapeTypeEllipse,
				pptx.ShapeTypeTriangle, pptx.ShapeTypeRightTriangle, pptx.ShapeTypeDiamond,
				pptx.ShapeTypePentagon, pptx.ShapeTypeHexagon, pptx.ShapeTypeOctagon,
				pptx.ShapeTypeParallelogram, pptx.ShapeTypeTrapezoid, pptx.ShapeTypeNonIsoscelesTrapezoid,
				pptx.ShapeTypeIsoscelesTrapezoid, pptx.ShapeTypeCube, pptx.ShapeTypeCan,
				pptx.ShapeTypeCone, pptx.ShapeTypeCylinder, pptx.ShapeTypeFoldedCorner,
				pptx.ShapeTypeBevel, pptx.ShapeTypeDonut, pptx.ShapeTypeNoSmoking,
				pptx.ShapeTypeBlockArc, pptx.ShapeTypeSmileyFace, pptx.ShapeTypeArc,
				pptx.ShapeTypeChord, pptx.ShapeTypePie, pptx.ShapeTypeTeardrop,
				pptx.ShapeTypePlaque, pptx.ShapeTypeMusicNote, pptx.ShapeTypePictureFrame,
				pptx.ShapeTypeHeart, pptx.ShapeTypeLightningBolt, pptx.ShapeTypeSun,
				pptx.ShapeTypeMoon, pptx.ShapeTypeCloud, pptx.ShapeTypeMathPlus, pptx.ShapeTypeMathMinus,
			},
		},
		{
			Category: "Arrows",
			Types: []string{
				pptx.ShapeTypeRightArrow, pptx.ShapeTypeLeftArrow, pptx.ShapeTypeUpArrow, pptx.ShapeTypeDownArrow,
				pptx.ShapeTypeLeftRightArrow, pptx.ShapeTypeUpDownArrow, pptx.ShapeTypeQuadArrow,
				pptx.ShapeTypeLeftRightUpArrow, pptx.ShapeTypeBentArrow, pptx.ShapeTypeUturnArrow,
				pptx.ShapeTypeCurvedRightArrow, pptx.ShapeTypeCurvedLeftArrow, pptx.ShapeTypeCurvedUpArrow,
				pptx.ShapeTypeCurvedDownArrow, pptx.ShapeTypeCurvedLeftRightArrow, pptx.ShapeTypeCurvedUpDownArrow,
				pptx.ShapeTypeStripedRightArrow, pptx.ShapeTypeNotchedRightArrow, pptx.ShapeTypePentagonArrow,
				pptx.ShapeTypeChevronArrow, pptx.ShapeTypeCircularArrow,
			},
		},
		{
			Category: "Callouts",
			Types: []string{
				pptx.ShapeTypeRightArrowCallout, pptx.ShapeTypeLeftArrowCallout, pptx.ShapeTypeUpArrowCallout,
				pptx.ShapeTypeDownArrowCallout, pptx.ShapeTypeLeftRightArrowCallout, pptx.ShapeTypeUpDownArrowCallout,
				pptx.ShapeTypeWedgeRectCallout, pptx.ShapeTypeWedgeEllipseCallout, pptx.ShapeTypeCloudCallout,
			},
		},
		{
			Category: "Flow Chart",
			Types: []string{
				pptx.ShapeTypeFlowChartProcess, pptx.ShapeTypeFlowChartDecision, pptx.ShapeTypeFlowChartTerminator,
				pptx.ShapeTypeFlowChartDocument, pptx.ShapeTypeFlowChartPredefinedProcess, pptx.ShapeTypeFlowChartInternalStorage,
				pptx.ShapeTypeFlowChartData, pptx.ShapeTypeFlowChartInputOutput, pptx.ShapeTypeFlowChartManualInput,
				pptx.ShapeTypeFlowChartManualOperation, pptx.ShapeTypeFlowChartConnector, pptx.ShapeTypeFlowChartOffpageConnector,
				pptx.ShapeTypeFlowChartPunchedCard, pptx.ShapeTypeFlowChartPunchedTape, pptx.ShapeTypeFlowChartSummingJunction,
				pptx.ShapeTypeFlowChartOr, pptx.ShapeTypeFlowChartCollate, pptx.ShapeTypeFlowChartSort,
				pptx.ShapeTypeFlowChartExtract, pptx.ShapeTypeFlowChartMerge, pptx.ShapeTypeFlowChartOnlineStorage,
				pptx.ShapeTypeFlowChartDelay, pptx.ShapeTypeFlowChartMagneticTape, pptx.ShapeTypeFlowChartMagneticDisk,
				pptx.ShapeTypeFlowChartMagneticDrum, pptx.ShapeTypeFlowChartDisplay, pptx.ShapeTypeFlowChartPreparation,
			},
		},
		{
			Category: "Stars & Banners",
			Types: []string{
				pptx.ShapeTypeStar4, pptx.ShapeTypeStar5, pptx.ShapeTypeStar6, pptx.ShapeTypeStar8,
				pptx.ShapeTypeStar10, pptx.ShapeTypeStar12, pptx.ShapeTypeStar16, pptx.ShapeTypeStar24,
				pptx.ShapeTypeStar32, pptx.ShapeTypeRibbon, pptx.ShapeTypeWave, pptx.ShapeTypeSeal,
				pptx.ShapeTypeSeal4, pptx.ShapeTypeSeal8, pptx.ShapeTypeSeal16, pptx.ShapeTypeSeal32,
			},
		},
		{
			Category: "Action Buttons",
			Types: []string{
				pptx.ShapeTypeActionButtonBlank, pptx.ShapeTypeActionButtonHome, pptx.ShapeTypeActionButtonHelp,
				pptx.ShapeTypeActionButtonInformation, pptx.ShapeTypeActionButtonForwardNext, pptx.ShapeTypeActionButtonBackPrevious,
				pptx.ShapeTypeActionButtonBeginning, pptx.ShapeTypeActionButtonEnd, pptx.ShapeTypeActionButtonReturn,
				pptx.ShapeTypeActionButtonDocument, pptx.ShapeTypeActionButtonSound, pptx.ShapeTypeActionButtonMovie,
			},
		},
	}

	for _, cat := range shapes {
		slide := pptx.NewSlide(cat.Category)
		tbox := pptx.NewTextBox(cat.Category, 0.5, 0.2, 9.0, 0.5)
		tbox = tbox.WithAutoFit(pptx.TextAutoFitNone)
		slide = slide.AddShape(tbox)

		// Grid layout
		cols := 6
		startX, startY := 0.5, 1.0
		cw, ch := 1.2, 1.2
		gapX, gapY := 0.2, 0.2

		for i, shapeType := range cat.Types {
			row := i / cols
			col := i % cols

			x := startX + float64(col)*(cw+gapX)
			y := startY + float64(row)*(ch+gapY)

			// Create shape
			shape := pptx.NewShape(shapeType, pptx.Inches(x), pptx.Inches(y), pptx.Inches(cw), pptx.Inches(ch))
			shape = shape.WithText(shapeType)                   // Label it
			shape = shape.WithFill(pptx.NewShapeFill("4472C4")) // Standard blue
			shape = shape.WithLine(pptx.NewShapeLine("000000", pptx.Points(1)))

			// Small text for label
			tf := pptx.NewTextFrame()
			tf.AutoFit = pptx.TextAutoFitNormal // Shrink text to fit
			shape = shape.WithTextFrame(tf)

			slide = slide.AddShape(shape)
		}
		builder.AddSlide(slide)
	}

	outPath := filepath.Join(t.TempDir(), "shape_showcase.pptx")
	if err := builder.WriteToFile(outPath); err != nil {
		t.Fatalf("failed to save presentation: %s", err)
	}
}
