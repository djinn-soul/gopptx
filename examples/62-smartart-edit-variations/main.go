// examples/62-smartart-edit-variations/main.go demonstrates many SmartArt edit variations.
//
// Run with:
//   go run ./examples/62-smartart-edit-variations/main.go
//
// Output:
//   examples/output/62_smartart_edit_variations.pptx
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir      = "examples/output"
	outputFile     = "examples/output/62_smartart_edit_variations.pptx"
	smartArtXInch  = 0.6
	smartArtYInch  = 2.0
	smartArtWInch  = 8.8
	smartArtHInch  = 3.3
)

type smartArtVariation struct {
	name      string
	seedCount int
	apply     func(e *editor.PresentationEditor, slideIndex, shapeID int) error
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("example failed: %v", err)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	variations := smartArtVariations()
	baseDeck, err := buildBaseDeck(variations)
	if err != nil {
		return err
	}

	e, err := editor.OpenPresentationEditorFromBytes(baseDeck)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}

	for i, v := range variations {
		shapeID, err := findFirstGraphicFrameID(e, i)
		if err != nil {
			return fmt.Errorf("variation %02d (%s): %w", i+1, v.name, err)
		}
		if err := v.apply(e, i, shapeID); err != nil {
			return fmt.Errorf("apply variation %02d (%s): %w", i+1, v.name, err)
		}
		fmt.Printf("[%02d] %s\n", i+1, v.name)
	}

	out, err := e.SaveToBytes()
	if err != nil {
		return fmt.Errorf("serialize presentation: %w", err)
	}
	if err := os.WriteFile(outputFile, out, 0o600); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	fmt.Printf("Saved %d SmartArt edit variations -> %s\n", len(variations), outputFile)
	return nil
}

func buildBaseDeck(variations []smartArtVariation) ([]byte, error) {
	slides := make([]elements.SlideContent, 0, len(variations))
	for i, v := range variations {
		if v.seedCount < 2 || v.seedCount > 5 {
			return nil, fmt.Errorf("variation %02d invalid seedCount=%d (allowed 2..5)", i+1, v.seedCount)
		}
		sa := smartart.NewSmartArt(smartart.AccentProcess).
			Position(styling.Inches(smartArtXInch), styling.Inches(smartArtYInch)).
			Size(styling.Inches(smartArtWInch), styling.Inches(smartArtHInch)).
			WithAltText(fmt.Sprintf("Variation %02d: %s", i+1, v.name)).
			AddItems(seedItems(v.seedCount, i))

		slide := elements.NewSlide(fmt.Sprintf("SmartArt Variation %02d", i+1))
		slide.SmartArtDiagrams = append(slide.SmartArtDiagrams, sa)
		slides = append(slides, slide)
	}

	data, err := pptx.CreateWithSlides("SmartArt Edit Variations (15)", slides)
	if err != nil {
		return nil, fmt.Errorf("create base deck: %w", err)
	}
	return data, nil
}

func findFirstGraphicFrameID(e *editor.PresentationEditor, slideIndex int) (int, error) {
	shapes, err := e.GetShapes(slideIndex)
	if err != nil {
		return 0, fmt.Errorf("GetShapes(%d): %w", slideIndex, err)
	}
	for _, sh := range shapes {
		if sh.Type == "graphicFrame" {
			return sh.ID, nil
		}
	}
	return 0, fmt.Errorf("no graphicFrame found on slide %d", slideIndex)
}

func smartArtVariations() []smartArtVariation {
	return []smartArtVariation{
		{name: "Update text items (3 steps)", seedCount: 3, apply: func(e *editor.PresentationEditor, s, id int) error {
			return e.UpdateSmartArt(s, id, []string{"Plan", "Build", "Ship"})
		}},
		{name: "Change layout to AccentProcess", seedCount: 2, apply: func(e *editor.PresentationEditor, s, id int) error {
			return e.ChangeSmartArtLayout(s, id, smartart.AccentProcess)
		}},
		{name: "Change layout to AlternatingFlow + update text", seedCount: 3, apply: func(e *editor.PresentationEditor, s, id int) error {
			if err := e.ChangeSmartArtLayout(s, id, smartart.AlternatingFlow); err != nil {
				return err
			}
			return e.UpdateSmartArt(s, id, []string{"Input", "Transform", "Output"})
		}},
		{name: "Change layout to VerticalBlockList + update text", seedCount: 3, apply: func(e *editor.PresentationEditor, s, id int) error {
			if err := e.ChangeSmartArtLayout(s, id, smartart.VerticalBlockList); err != nil {
				return err
			}
			return e.UpdateSmartArt(s, id, []string{"North", "Central", "South"})
		}},
		{name: "Set style simple2 + colorful2", seedCount: 4, apply: func(e *editor.PresentationEditor, s, id int) error {
			return e.SetSmartArtStyle(
				s,
				id,
				"urn:microsoft.com/office/officeart/2005/8/quickstyle/simple2",
				"urn:microsoft.com/office/officeart/2005/8/colors/colorful2",
			)
		}},
		{name: "Set style simple1 + accent1_2", seedCount: 5, apply: func(e *editor.PresentationEditor, s, id int) error {
			return e.SetSmartArtStyle(
				s,
				id,
				"urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1",
				"urn:microsoft.com/office/officeart/2005/8/colors/accent1_2",
			)
		}},
		{name: "Replace nodes with 2-item flow", seedCount: 2, apply: func(e *editor.PresentationEditor, s, id int) error {
			return e.SetSmartArtNodes(s, id, []smartart.Node{
				smartart.NewNode("Request"),
				smartart.NewNode("Response"),
			})
		}},
		{name: "Replace nodes with 3-item flow", seedCount: 3, apply: func(e *editor.PresentationEditor, s, id int) error {
			return e.SetSmartArtNodes(s, id, []smartart.Node{
				smartart.NewNode("Backlog"),
				smartart.NewNode("Sprint"),
				smartart.NewNode("Release"),
			})
		}},
		{name: "Change to BasicCycle + 5 nodes", seedCount: 3, apply: func(e *editor.PresentationEditor, s, id int) error {
			if err := e.ChangeSmartArtLayout(s, id, smartart.BasicCycle); err != nil {
				return err
			}
			return e.SetSmartArtNodes(s, id, []smartart.Node{
				smartart.NewNode("Observe"),
				smartart.NewNode("Orient"),
				smartart.NewNode("Decide"),
				smartart.NewNode("Act"),
				smartart.NewNode("Repeat"),
			})
		}},
		{name: "Change to BasicVenn + 3 nodes", seedCount: 3, apply: func(e *editor.PresentationEditor, s, id int) error {
			if err := e.ChangeSmartArtLayout(s, id, smartart.BasicVenn); err != nil {
				return err
			}
			return e.SetSmartArtNodes(s, id, []smartart.Node{
				smartart.NewNode("People"),
				smartart.NewNode("Process"),
				smartart.NewNode("Tech"),
			})
		}},
		{name: "Change to LinearVenn + 4 nodes", seedCount: 3, apply: func(e *editor.PresentationEditor, s, id int) error {
			if err := e.ChangeSmartArtLayout(s, id, smartart.LinearVenn); err != nil {
				return err
			}
			return e.SetSmartArtNodes(s, id, []smartart.Node{
				smartart.NewNode("A"),
				smartart.NewNode("B"),
				smartart.NewNode("C"),
				smartart.NewNode("D"),
			})
		}},
		{name: "Change to BasicMatrix + 4 nodes", seedCount: 3, apply: func(e *editor.PresentationEditor, s, id int) error {
			if err := e.ChangeSmartArtLayout(s, id, smartart.BasicMatrix); err != nil {
				return err
			}
			return e.SetSmartArtNodes(s, id, []smartart.Node{
				smartart.NewNode("Low Risk"),
				smartart.NewNode("High Value"),
				smartart.NewNode("Quick Wins"),
				smartart.NewNode("Big Bets"),
			})
		}},
		{name: "Change to BasicPyramid + 3 nodes", seedCount: 3, apply: func(e *editor.PresentationEditor, s, id int) error {
			if err := e.ChangeSmartArtLayout(s, id, smartart.BasicPyramid); err != nil {
				return err
			}
			return e.SetSmartArtNodes(s, id, []smartart.Node{
				smartart.NewNode("Vision"),
				smartart.NewNode("Strategy"),
				smartart.NewNode("Execution"),
			})
		}},
		{name: "Delete original and add OrgChart hierarchy", seedCount: 4, apply: func(e *editor.PresentationEditor, s, id int) error {
			if err := e.DeleteSmartArt(s, id); err != nil {
				return err
			}
			root := smartart.NewNode("Director").
				WithChild(smartart.NewNode("Eng")).
				WithChild(smartart.NewNode("Ops"))
			_, err := e.AddSmartArt(s, smartart.NewSmartArt(smartart.OrgChart).
				Position(styling.Inches(smartArtXInch), styling.Inches(smartArtYInch)).
				Size(styling.Inches(smartArtWInch), styling.Inches(smartArtHInch)).
				WithAltText("OrgChart replacement").
				AddNode(root))
			return err
		}},
		{name: "Delete original and add PictureAccentList", seedCount: 5, apply: func(e *editor.PresentationEditor, s, id int) error {
			if err := e.DeleteSmartArt(s, id); err != nil {
				return err
			}
			_, err := e.AddSmartArt(s, smartart.NewSmartArt(smartart.PictureAccentList).
				Position(styling.Inches(smartArtXInch), styling.Inches(smartArtYInch)).
				Size(styling.Inches(smartArtWInch), styling.Inches(smartArtHInch)).
				WithAltText("Replacement SmartArt").
				AddItems([]string{"Scan", "Draft", "Ship"}))
			return err
		}},
	}
}

func seedItems(count, variationIndex int) []string {
	items := make([]string, count)
	for i := 0; i < count; i++ {
		items[i] = fmt.Sprintf("Seed %02d.%d", variationIndex+1, i+1)
	}
	return items
}
