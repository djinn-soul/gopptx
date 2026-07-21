package main

import (
	"fmt"
	"log"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	smartArtOutputPath = "examples/output/24_smartart_smoke.pptx"
	stepBuild          = "Build"
	stepPlan           = "Plan"
	stepShip           = "Ship"
	smartArtX          = styling.Length(45 * 12700)
	smartArtY          = styling.Length(115 * 12700)
	smartArtCX         = styling.Length(625 * 12700)
	smartArtCY         = styling.Length(335 * 12700)
)

func main() {
	pres := pptx.NewPresentationBuilder("SmartArt Full Layout Showcase")
	pres.AddBulletSlide("SmartArt Full Layout Showcase", []string{
		"Task 24 expanded demo: all currently supported SmartArt layouts.",
		"Each slide uses native SmartArt (no fallback labels).",
		"Content/style variants are rotated per layout.",
	})

	layouts := supportedLayoutsInShowcaseOrder()
	for i, layout := range layouts {
		title := fmt.Sprintf("%02d. %s", i+1, layout.Name())
		sa := buildShowcaseSmartArt(layout, i)
		pres.AddSlide(pptx.NewSlide(title).AddSmartArt(sa))
	}

	if err := os.MkdirAll("examples/output", 0o750); err != nil {
		log.Fatal(err)
	}
	if err := pres.WriteToFile(smartArtOutputPath); err != nil {
		log.Fatal(err)
	}
	log.Println("Saved", smartArtOutputPath)
}

func supportedLayoutsInShowcaseOrder() []smartart.Layout {
	return []smartart.Layout{
		smartart.BasicBlockList,
		smartart.VerticalBlockList,
		smartart.HorizontalBulletLst,
		smartart.SquareAccentList,
		smartart.PictureAccentList,
		smartart.BasicProcess,
		smartart.AccentProcess,
		smartart.AlternatingFlow,
		smartart.ContinuousBlockProcess,
		smartart.BasicCycle,
		smartart.TextCycle,
		smartart.BlockCycle,
		smartart.OrgChart,
		smartart.Hierarchy,
		smartart.HorizontalHierarchy,
		smartart.BasicVenn,
		smartart.LinearVenn,
		smartart.StackedVenn,
		smartart.BasicRadial,
		smartart.BasicMatrix,
		smartart.TitledMatrix,
		smartart.BasicPyramid,
		smartart.InvertedPyramid,
		smartart.PictureStrips,
		smartart.PictureGrid,
	}
}

func buildShowcaseSmartArt(layout smartart.Layout, idx int) smartart.SmartArt {
	sa := smartart.NewSmartArt(layout).
		Position(smartArtX, smartArtY).
		Size(smartArtCX, smartArtCY).
		WithColorStyle("urn:microsoft.com/office/officeart/2005/8/colors/accent1_2").
		WithQuickStyle("urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1").
		WithAltText(layout.Name() + " SmartArt example")

	if isHierarchyLayout(layout) {
		return sa.AddNode(hierarchyRoot(idx, layout))
	}

	return sa.AddItems(itemsForLayout(layout, idx))
}

func isHierarchyLayout(layout smartart.Layout) bool {
	switch layout {
	case smartart.OrgChart, smartart.Hierarchy, smartart.HorizontalHierarchy:
		return true
	default:
		return false
	}
}

func hierarchyRoot(idx int, layout smartart.Layout) smartart.Node {
	root := smartart.NewNode(fmt.Sprintf("Lead %d", idx+1))

	switch layout {
	case smartart.OrgChart:
		return root.
			WithChild(smartart.NewNode("Finance")).
			WithChild(smartart.NewNode("Eng"))
	default:
		return root.
			WithChild(smartart.NewNode("Branch A").
				WithChild(smartart.NewNode("A1")).
				WithChild(smartart.NewNode("A2"))).
			WithChild(smartart.NewNode("Branch B").
				WithChild(smartart.NewNode("B1")))
	}
}

func itemsForLayout(layout smartart.Layout, idx int) []string {
	switch layout {
	case smartart.BasicBlockList:
		return []string{
			fmt.Sprintf("Topic %dA", idx+1),
			fmt.Sprintf("Topic %dB", idx+1),
			fmt.Sprintf("Topic %dC", idx+1),
			fmt.Sprintf("Topic %dD", idx+1),
			fmt.Sprintf("Topic %dE", idx+1),
		}
	case smartart.VerticalBlockList:
		return []string{
			fmt.Sprintf("Topic %dA", idx+1),
			fmt.Sprintf("Topic %dB", idx+1),
			fmt.Sprintf("Topic %dC", idx+1),
		}
	case smartart.HorizontalBulletLst:
		return []string{
			fmt.Sprintf("Pillar %dA", idx+1),
			fmt.Sprintf("Detail %dB", idx+1),
			fmt.Sprintf("Detail %dC", idx+1),
		}
	case smartart.SquareAccentList:
		return []string{
			fmt.Sprintf("Topic %dA", idx+1),
			fmt.Sprintf("Topic %dB", idx+1),
			fmt.Sprintf("Topic %dC", idx+1),
			fmt.Sprintf("Topic %dD", idx+1),
		}
	case smartart.PictureAccentList:
		return []string{
			fmt.Sprintf("Topic %dA", idx+1),
			fmt.Sprintf("Topic %dB", idx+1),
			fmt.Sprintf("Topic %dC", idx+1),
		}
	case smartart.BasicProcess:
		return []string{stepPlan, stepBuild, stepShip}
	case smartart.AccentProcess:
		return []string{stepPlan, "Design", stepBuild, "Test", stepShip}
	case smartart.AlternatingFlow, smartart.ContinuousBlockProcess:
		return []string{stepPlan, stepBuild, stepShip}
	case smartart.BasicCycle,
		smartart.TextCycle,
		smartart.BlockCycle:
		return []string{"Plan", "Develop", "Test", "Deploy", "Learn"}
	case smartart.BasicVenn:
		return []string{"People", "Process", "Platform"}
	case smartart.LinearVenn, smartart.StackedVenn:
		return []string{"People", "Process", "Platform", "Culture"}
	case smartart.BasicRadial:
		return []string{"North", "East", "South", "West"}
	case smartart.BasicMatrix,
		smartart.TitledMatrix:
		return []string{"Q1", "Q2", "Q3", "Q4"}
	case smartart.BasicPyramid,
		smartart.InvertedPyramid:
		return []string{"Level 1", "Level 2", "Level 3"}
	case smartart.PictureStrips:
		return []string{"Scene 1", "Scene 2", "Scene 3"}
	case smartart.PictureGrid:
		return []string{"Scene 1", "Scene 2", "Scene 3", "Scene 4"}
	default:
		return []string{"One", "Two", "Three"}
	}
}
