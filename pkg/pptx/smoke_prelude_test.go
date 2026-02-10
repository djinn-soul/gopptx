package pptx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratePreludeSmokeSample(t *testing.T) {
	builder := NewPresentationBuilder("GoPPTX Prelude Helpers Showcase").
		WithMetadata(PresentationMetadata{
			Title:       "Prelude Showcase",
			Subject:     "Testing new ergonomic features",
			Creator:     "GoPPTX Agent",
			Description: "Showcase for expanded prelude helpers, unit converters, and shape macros.",
		})

	// 1. Title Slide
	builder.AddTitleSlide("Ergonomic Prelude Helpers")

	// 2. Bullet Slide
	builder.AddBulletSlide("New Ergonomic API Features", []string{
		"Fluent PresentationBuilder API",
		"One-line Title/Bullet/Shapes slides",
		"Overflow-protected unit converters",
		"Hundredthsndths of a point font size helper",
	})

	// 3. New Shape Macros Showcase
	builder.AddShapesSlide("Expanded Shape Macros",
		NewStar(1, 1, 2).WithFill(NewShapeFill(ColorMaterialAmber)).WithText("New Star"),
		NewHeart(4, 1, 2).WithFill(NewShapeFill(ColorMaterialRed)).WithText("New Heart"),
		NewCircle(7, 1, 2).WithFill(NewShapeFill(ColorMaterialBlue)).WithText("New Circle"),
		NewBadge("NEW FEATURE", 1, 4, ColorMaterialGreen),
		NewBadge("IBM CARBON", 4, 4, ColorCarbonBlue60),
	)

	// 4. FlowChart Showcase
	builder.AddShapesSlide("FlowChart Components",
		NewFlowChartProcess(1, 1, 2, 1).WithText("Start Process"),
		NewFlowChartDecision(4, 0.75, 1.5, 1.5).WithText("Is Valid?"),
		NewFlowChartDocument(7, 1, 2, 1.5).WithText("Output Doc"),
		NewFlowChartData(1, 4, 2, 1).WithText("Input Data"),
		NewFlowChartTerminator(4, 4, 2, 1).WithText("End"),
	)

	// 5. Unit Conversion Showcase
	builder.AddShapesSlide("Robust Unit Conversions",
		NewRectangle(0.5, 1, 8, 0.5).WithText("Full-width Shape (8 inches)"),
		NewRectangle(1, 3, 2, 2).WithText("Standard 2x2 Inch Shape"),
	)

	data, err := builder.Build()
	if err != nil {
		t.Fatalf("failed to build showcase deck: %v", err)
	}

	outputPath := filepath.Join("..", "..", "smoke_samples", "16_prelude_helpers.pptx")
	err = os.WriteFile(outputPath, data, 0o644)
	if err != nil {
		t.Fatalf("failed to write smoke sample: %v", err)
	}

	t.Logf("Successfully generated %s", outputPath)
}
