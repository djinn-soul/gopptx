package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func generatePreludeHelpers() ([]byte, error) {
	builder := pptx.NewPresentationBuilder("GoPPTX Prelude Helpers Showcase").
		WithMetadata(pptx.PresentationMetadata{
			Title:       "Prelude Showcase",
			Subject:     "Testing new ergonomic features",
			Creator:     "GoPPTX Agent",
			Description: "Showcase for expanded prelude helpers, unit converters, and shape macros.",
		})

	builder.AddTitleSlide("Ergonomic Prelude Helpers")

	builder.AddBulletSlide("New Ergonomic API Features", []string{
		"Fluent PresentationBuilder API",
		"One-line Title/Bullet/Shapes slides",
		"Overflow-protected unit converters",
		"Hundredths of a point font size helper",
	})

	builder.AddShapesSlide("Expanded Shape Macros",
		pptx.NewStar(1, 1, 2).WithFill(pptx.NewShapeFill(pptx.ColorMaterialAmber)).WithText("New Star"),
		pptx.NewHeart(4, 1, 2).WithFill(pptx.NewShapeFill(pptx.ColorMaterialRed)).WithText("New Heart"),
		pptx.NewCircle(7, 1, 2).WithFill(pptx.NewShapeFill(pptx.ColorMaterialBlue)).WithText("New Circle"),
		pptx.NewBadge("NEW FEATURE", 1, 4, pptx.ColorMaterialGreen),
		pptx.NewBadge("IBM CARBON", 4, 4, pptx.ColorCarbonBlue60),
	)

	builder.AddShapesSlide("FlowChart Components",
		pptx.NewFlowChartProcess(1, 1, 2, 1).WithText("Start Process"),
		pptx.NewFlowChartDecision(4, 0.75, 1.5, 1.5).WithText("Is Valid?"),
		pptx.NewFlowChartDocument(7, 1, 2, 1.5).WithText("Output Doc"),
		pptx.NewFlowChartData(1, 4, 2, 1).WithText("Input Data"),
		pptx.NewFlowChartTerminator(4, 4, 2, 1).WithText("End"),
	)

	builder.AddShapesSlide("Robust Unit Conversions",
		pptx.NewRectangle(0.5, 1, 8, 0.5).WithText("Full-width Shape (8 inches)"),
		pptx.NewRectangle(1, 3, 2, 2).WithText("Standard 2x2 Inch Shape"),
	)

	return builder.Build()
}
