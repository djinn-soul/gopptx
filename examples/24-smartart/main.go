package main

import (
	"fmt"
	"log"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func main() {
	const outputPath = "examples/output/24_smartart_smoke.pptx"

	pres := pptx.NewPresentationBuilder("SmartArt Demo")
	pres.AddBulletSlide("SmartArt Demo", []string{
		"This slide contains SmartArt diagrams.",
		"See following slides.",
	})

	// List Layout
	saList := smartart.NewSmartArt(smartart.VerticalBlockList).
		AddNode(smartart.NewNode("Block 1")).
		AddNode(smartart.NewNode("Block 2").WithColor("FF0000")).
		AddNode(smartart.NewNode("Block 3")).
		Position(styling.Points(50), styling.Points(150)).
		Size(styling.Points(400), styling.Points(300))

	pres.AddSlide(pptx.NewSlide("Vertical Block List").AddSmartArt(saList))

	// Process Layout
	saProcess := smartart.NewSmartArt(smartart.BasicProcess).
		AddNode(smartart.NewNode("Phase 1")).
		AddNode(smartart.NewNode("Phase 2")).
		AddNode(smartart.NewNode("Phase 3")).
		Position(styling.Points(50), styling.Points(150)).
		Size(styling.Points(600), styling.Points(200)).
		WithColorStyle("urn:microsoft.com/office/officeart/2005/8/colors/colorful1").
		WithQuickStyle("urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1")

	pres.AddSlide(pptx.NewSlide("Basic Process").AddSmartArt(saProcess))

	// Hierarchy Layout
	saOrg := smartart.NewSmartArt(smartart.OrgChart).
		AddNode(smartart.NewNode("CEO").
			WithChild(smartart.NewNode("VP Sales").
				WithChild(smartart.NewNode("Manager 1")).
				WithChild(smartart.NewNode("Manager 2"))).
			WithChild(smartart.NewNode("VP Eng").
				WithChild(smartart.NewNode("Dev 1")).
				WithChild(smartart.NewNode("Dev 2")))).
		Position(styling.Points(50), styling.Points(150)).
		Size(styling.Points(600), styling.Points(400))

	pres.AddSlide(pptx.NewSlide("Organization Chart").AddSmartArt(saOrg))

	// Cycle Layout
	saCycle := smartart.NewSmartArt(smartart.BasicCycle).
		AddItems([]string{"Plan", "Develop", "Test", "Deploy", "Monitor"}).
		Position(styling.Points(150), styling.Points(150)).
		Size(styling.Points(400), styling.Points(400))

	pres.AddSlide(pptx.NewSlide("Basic Cycle").AddSmartArt(saCycle))

	if err := os.MkdirAll("examples/output", 0o755); err != nil {
		log.Fatal(err)
	}
	if err := pres.WriteToFile(outputPath); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved", outputPath)
}
