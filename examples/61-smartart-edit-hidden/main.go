// examples/61-smartart-edit-hidden/main.go demonstrates:
//   - Creating a presentation with hidden slides
//   - Adding SmartArt to a slide
//   - Editing SmartArt: updating text, changing layout, changing style, replacing nodes
//   - Deleting a SmartArt diagram
//
// Run with: go run ./examples/61-smartart-edit-hidden/main.go
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir  = "examples/output"
	outputFile = "examples/output/61_smartart_edit_hidden.pptx"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("example failed: %v", err)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	data, err := buildPresentation()
	if err != nil {
		return err
	}

	e, err := editor.OpenPresentationEditorFromBytes(data)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}

	printSlideInventory(e.Slides())

	if err := e.SetSlideHidden(0, false); err != nil {
		return fmt.Errorf("SetSlideHidden(0): %w", err)
	}
	printFmtf("\nSlide 1 hidden after load: %v\n", e.Slides()[1].Hidden)

	shapes0, err := e.GetShapes(0)
	if err != nil {
		return fmt.Errorf("GetShapes(0): %w", err)
	}
	if err := editSlide0SmartArt(e, findGraphicFrameID(shapes0)); err != nil {
		return err
	}

	shapes2, err := e.GetShapes(2)
	if err != nil {
		return fmt.Errorf("GetShapes(2): %w", err)
	}
	if err := editSlide2SmartArt(e, findGraphicFrameID(shapes2)); err != nil {
		return err
	}

	if err := addAndDeleteTempSmartArt(e); err != nil {
		return err
	}

	out, err := e.SaveToBytes()
	if err != nil {
		return fmt.Errorf("serialize: %w", err)
	}
	if err := os.WriteFile(outputFile, out, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	printFmtf("\nSaved → %s\n", outputFile)
	return nil
}

// buildPresentation creates the initial three-slide deck in memory.
func buildPresentation() ([]byte, error) {
	sa1 := smartart.NewSmartArt(smartart.BasicProcess).
		Position(styling.Inches(0.5), styling.Inches(1.5)).
		Size(styling.Inches(9), styling.Inches(4)).
		WithColorStyle("urn:microsoft.com/office/officeart/2005/8/colors/accent1_2").
		WithQuickStyle("urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1").
		WithAltText("Project phases").
		AddNode(smartart.NewNode("Plan")).
		AddNode(smartart.NewNode("Build")).
		AddNode(smartart.NewNode("Ship"))

	sa3 := smartart.NewSmartArt(smartart.BasicCycle).
		Position(styling.Inches(1), styling.Inches(1.5)).
		Size(styling.Inches(8), styling.Inches(4.5)).
		WithColorStyle("urn:microsoft.com/office/officeart/2005/8/colors/colorful1").
		WithAltText("Development cycle").
		AddNode(smartart.NewNode("Discover")).
		AddNode(smartart.NewNode("Design")).
		AddNode(smartart.NewNode("Develop")).
		AddNode(smartart.NewNode("Deliver"))

	slide1 := elements.NewSlide("SmartArt: Project Phases")
	slide1.SmartArtDiagrams = append(slide1.SmartArtDiagrams, sa1)

	slide2 := elements.NewSlide("Hidden Draft Slide")
	slide2.Hidden = true
	slide2.Bullets = []string{"This slide is hidden — not shown in slideshow."}

	slide3 := elements.NewSlide("SmartArt: Development Cycle")
	slide3.SmartArtDiagrams = append(slide3.SmartArtDiagrams, sa3)

	data, err := pptx.CreateWithSlides("SmartArt Edit + Hidden Demo", []elements.SlideContent{slide1, slide2, slide3})
	if err != nil {
		return nil, fmt.Errorf("create presentation: %w", err)
	}
	return data, nil
}

// findGraphicFrameID returns the shape ID of the first graphicFrame in the list, or 0.
func findGraphicFrameID(shapes []common.Shape) int {
	for _, sh := range shapes {
		if sh.Type == "graphicFrame" {
			return sh.ID
		}
	}
	return 0
}

// printSlideInventory prints index, title, and hidden flag for each slide.
func printSlideInventory(slides []common.SlideMetadata) {
	printLine("=== Slide inventory ===")
	for _, s := range slides {
		printFmtf("  [%d] %q  hidden=%v\n", s.Index, s.Title, s.Hidden)
	}
}

// editSlide0SmartArt updates text, layout, and style on the SmartArt at shapeID.
func editSlide0SmartArt(e *editor.PresentationEditor, shapeID int) error {
	if shapeID == 0 {
		return nil
	}
	printFmtf("\nFound SmartArt on slide 0: shapeID=%d\n", shapeID)

	if err := e.UpdateSmartArt(0, shapeID, []string{"Discover", "Prototype", "Launch"}); err != nil {
		return fmt.Errorf("UpdateSmartArt: %w", err)
	}
	printLine("UpdateSmartArt: replaced text items on slide 0 SmartArt")

	if err := e.ChangeSmartArtLayout(0, shapeID, smartart.AccentProcess); err != nil {
		return fmt.Errorf("ChangeSmartArtLayout: %w", err)
	}
	printLine("ChangeSmartArtLayout: switched to AccentProcess layout")

	if err := e.SetSmartArtStyle(0, shapeID,
		"urn:microsoft.com/office/officeart/2005/8/quickstyle/simple2",
		"urn:microsoft.com/office/officeart/2005/8/colors/colorful2",
	); err != nil {
		return fmt.Errorf("SetSmartArtStyle: %w", err)
	}
	printLine("SetSmartArtStyle: updated quick style and color style")
	return nil
}

// editSlide2SmartArt replaces the node tree on the SmartArt at shapeID.
func editSlide2SmartArt(e *editor.PresentationEditor, shapeID int) error {
	if shapeID == 0 {
		return nil
	}
	printFmtf("\nFound SmartArt on slide 2: shapeID=%d\n", shapeID)

	newNodes := []smartart.Node{
		smartart.NewNode("Plan"),
		smartart.NewNode("Build"),
		smartart.NewNode("Measure"),
		smartart.NewNode("Learn"),
		smartart.NewNode("Repeat"),
	}
	if err := e.SetSmartArtNodes(2, shapeID, newNodes); err != nil {
		return fmt.Errorf("SetSmartArtNodes: %w", err)
	}
	printLine("SetSmartArtNodes: replaced cycle nodes with 5-step loop")
	return nil
}

// addAndDeleteTempSmartArt adds a temporary SmartArt to slide 1 then removes it.
func addAndDeleteTempSmartArt(e *editor.PresentationEditor) error {
	tempSA := smartart.NewSmartArt(smartart.BasicBlockList).
		Position(styling.Inches(1), styling.Inches(1)).
		Size(styling.Inches(5), styling.Inches(3)).
		AddNode(smartart.NewNode("Temp A")).
		AddNode(smartart.NewNode("Temp B"))

	tempID, err := e.AddSmartArt(1, tempSA)
	if err != nil {
		return fmt.Errorf("AddSmartArt (temp): %w", err)
	}
	printFmtf("\nAdded temp SmartArt on slide 1: shapeID=%d\n", tempID)

	if err := e.DeleteSmartArt(1, tempID); err != nil {
		return fmt.Errorf("DeleteSmartArt: %w", err)
	}
	printFmtf("DeleteSmartArt: removed shapeID=%d from slide 1\n", tempID)
	return nil
}

func printLine(args ...any) {
	_, _ = io.WriteString(os.Stdout, fmt.Sprintln(args...))
}

func printFmtf(format string, args ...any) {
	_, _ = io.WriteString(os.Stdout, fmt.Sprintf(format, args...))
}
