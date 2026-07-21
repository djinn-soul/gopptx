// examples/78-enum-api demonstrates the enum constants available in gopptx.
//
// Covers shape type constants (geometry, arrows, callouts, flowchart, action
// buttons, math shapes), slide layout constants, placeholder type constants,
// chart type constructors (via the chart_compat layer), animation effect
// constants, and connector type constants.
//
// Run with: go run ./examples/78-enum-api/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "78_enum_api.pptx"

	colorBlue = "4472C4"
	colorRed  = "C0504D"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	slides := []pptx.SlideContent{
		buildCoreShapeTypesSlide(),
		buildArrowTypesSlide(),
		buildLayoutConstantsSlide(),
		buildConnectorConstantsSlide(),
		buildAnimationConstantsSlide(),
		buildPlaceholderConstantsSlide(),
		buildTextFrameConstantsSlide(),
	}

	outputPath := filepath.Join(outputDir, outputFile)
	data, err := pptx.CreateWithSlides("Enum API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildCoreShapeTypesSlide() pptx.SlideContent {
	coreShapeTypes := []shapeConstant{
		{"Rect", pptx.ShapeTypeRectangle, colorBlue},
		{"RoundRect", pptx.ShapeTypeRoundedRectangle, colorRed},
		{"Ellipse", pptx.ShapeTypeEllipse, "9BBB59"},
		{"Triangle", pptx.ShapeTypeTriangle, "F79646"},
		{"Diamond", pptx.ShapeTypeDiamond, "8064A2"},
		{"Pentagon", pptx.ShapeTypePentagon, "4BACC6"},
		{"Hexagon", pptx.ShapeTypeHexagon, colorBlue},
		{"Parallelogram", pptx.ShapeTypeParallelogram, colorRed},
	}
	return buildShapeConstantsSlide("Core Shape Type Constants", coreShapeTypes)
}

func buildArrowTypesSlide() pptx.SlideContent {
	arrowTypes := []shapeConstant{
		{"RightArrow", pptx.ShapeTypeRightArrow, colorBlue},
		{"LeftArrow", pptx.ShapeTypeLeftArrow, colorRed},
		{"UpArrow", pptx.ShapeTypeUpArrow, "9BBB59"},
		{"DownArrow", pptx.ShapeTypeDownArrow, "F79646"},
		{"LR Arrow", pptx.ShapeTypeLeftRightArrow, "8064A2"},
		{"UD Arrow", pptx.ShapeTypeUpDownArrow, "4BACC6"},
		{"Quad", pptx.ShapeTypeQuadArrow, colorBlue},
		{"Bent", pptx.ShapeTypeBentArrow, colorRed},
	}
	return buildShapeConstantsSlide("Arrow Shape Constants", arrowTypes)
}

type shapeConstant struct {
	name      string
	shapeType string
	color     string
}

func buildShapeConstantsSlide(title string, shapeTypes []shapeConstant) pptx.SlideContent {
	slide := pptx.NewSlide(title)
	for i, shapeType := range shapeTypes {
		col := float64(i%4) * 2.3
		row := float64(i/4) * 1.8
		shape := pptx.NewShape(
			shapeType.shapeType,
			styling.Inches(col+0.3), styling.Inches(row+1.5),
			styling.Inches(2), styling.Inches(1.2),
		).
			WithFill(pptx.NewShapeFill(shapeType.color)).
			WithText(shapeType.name)
		slide = slide.AddShape(shape)
	}
	return slide
}

func buildLayoutConstantsSlide() pptx.SlideContent {
	return pptx.NewSlide("Slide Layout Constants").
		AddBullet(fmt.Sprintf(
			"SlideLayoutTitleAndContent    = %q", pptx.SlideLayoutTitleAndContent,
		)).
		AddBullet(fmt.Sprintf(
			"SlideLayoutTitleOnly          = %q", pptx.SlideLayoutTitleOnly,
		)).
		AddBullet(fmt.Sprintf(
			"SlideLayoutBlank              = %q", pptx.SlideLayoutBlank,
		)).
		AddBullet(fmt.Sprintf(
			"SlideLayoutCenteredTitle      = %q", pptx.SlideLayoutCenteredTitle,
		)).
		AddBullet(fmt.Sprintf(
			"SlideLayoutTitleAndBigContent = %q", pptx.SlideLayoutTitleAndBigContent,
		)).
		AddBullet(fmt.Sprintf(
			"SlideLayoutTwoColumn          = %q", pptx.SlideLayoutTwoColumn,
		))
}

func buildConnectorConstantsSlide() pptx.SlideContent {
	return pptx.NewSlide("Connector Type Constants").
		AddBullet(fmt.Sprintf("ConnectorTypeStraight = %q", pptx.ConnectorTypeStraight)).
		AddBullet(fmt.Sprintf("ConnectorTypeElbow    = %q", pptx.ConnectorTypeElbow)).
		AddBullet(fmt.Sprintf("ConnectorTypeCurved   = %q", pptx.ConnectorTypeCurved)).
		AddBullet(fmt.Sprintf("ArrowTypeNone     = %q", pptx.ArrowTypeNone)).
		AddBullet(fmt.Sprintf("ArrowTypeTriangle = %q", pptx.ArrowTypeTriangle)).
		AddBullet(fmt.Sprintf("ArrowTypeStealth  = %q", pptx.ArrowTypeStealth))
}

func buildAnimationConstantsSlide() pptx.SlideContent {
	return pptx.NewSlide("Animation Effect Constants").
		AddBullet(fmt.Sprintf("AnimationEntranceAppear  = %q", pptx.AnimationEntranceAppear)).
		AddBullet(fmt.Sprintf("AnimationEntranceFade    = %q", pptx.AnimationEntranceFade)).
		AddBullet(fmt.Sprintf("AnimationEntranceFlyIn   = %q", pptx.AnimationEntranceFlyIn)).
		AddBullet(fmt.Sprintf("AnimationEntranceZoom    = %q", pptx.AnimationEntranceZoom)).
		AddBullet(fmt.Sprintf("AnimationEntranceBounce  = %q", pptx.AnimationEntranceBounce)).
		AddBullet(fmt.Sprintf("AnimationEmphasisSpin    = %q", pptx.AnimationEmphasisSpin)).
		AddBullet(fmt.Sprintf("AnimationExitFadeOut     = %q", pptx.AnimationExitFadeOut)).
		AddBullet(fmt.Sprintf("AnimationOnClick         = %q", pptx.AnimationOnClick))
}

func buildPlaceholderConstantsSlide() pptx.SlideContent {
	return pptx.NewSlide("Placeholder Type Constants").
		AddBullet(fmt.Sprintf("PlaceholderTypeTitle    = %q", pptx.PlaceholderTypeTitle)).
		AddBullet(fmt.Sprintf("PlaceholderTypeBody     = %q", pptx.PlaceholderTypeBody)).
		AddBullet(fmt.Sprintf("PlaceholderTypeSubTitle = %q", pptx.PlaceholderTypeSubTitle)).
		AddBullet(fmt.Sprintf("PlaceholderTypeChart    = %q", pptx.PlaceholderTypeChart)).
		AddBullet(fmt.Sprintf("PlaceholderTypeTbl      = %q", pptx.PlaceholderTypeTbl)).
		AddBullet(fmt.Sprintf("PlaceholderTypePic      = %q", pptx.PlaceholderTypePic)).
		AddBullet(fmt.Sprintf("PlaceholderTypeMedia    = %q", pptx.PlaceholderTypeMedia))
}

func buildTextFrameConstantsSlide() pptx.SlideContent {
	return pptx.NewSlide("Text Frame Constants").
		AddBullet(fmt.Sprintf("TextAnchorTop    = %q", pptx.TextAnchorTop)).
		AddBullet(fmt.Sprintf("TextAnchorMiddle = %q", pptx.TextAnchorMiddle)).
		AddBullet(fmt.Sprintf("TextAnchorBottom = %q", pptx.TextAnchorBottom)).
		AddBullet(fmt.Sprintf("TextWrapNone   = %q", pptx.TextWrapNone)).
		AddBullet(fmt.Sprintf("TextWrapSquare = %q", pptx.TextWrapSquare)).
		AddBullet(fmt.Sprintf("TextAutoFitNone   = %q", pptx.TextAutoFitNone)).
		AddBullet(fmt.Sprintf("TextAutoFitShape  = %q", pptx.TextAutoFitShape)).
		AddBullet(fmt.Sprintf("TextAutoFitNormal = %q", pptx.TextAutoFitNormal))
}
