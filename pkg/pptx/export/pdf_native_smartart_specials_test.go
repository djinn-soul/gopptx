package export

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestSmartArtFitScaleNeverEnlarges(t *testing.T) {
	if got := smartArtFitScale(1000, 1000, 200, 200); got != 1 {
		t.Errorf("got %v in a frame bigger than the content, want 1", got)
	}
	if got := smartArtFitScale(100, 1000, 200, 200); got != 0.5 {
		t.Errorf("got %v for a half-width frame, want 0.5", got)
	}
	if got := smartArtFitScale(1000, 50, 200, 200); got != 0.25 {
		t.Errorf("got %v for a quarter-height frame, want 0.25", got)
	}
	if got := smartArtFitScale(0, 0, 200, 200); got != 1 {
		t.Errorf("got %v for an unknown frame, want 1", got)
	}
}

// The fixed-size layouts used to draw past a small frame. They are exercised
// here through the renderer, which would panic or write nothing if the scaled
// geometry came out non-positive.
func TestFixedSizeSpecialLayoutsSurviveATinyFrame(t *testing.T) {
	pdf := newTestPDF(t)
	if err := configureNativePDFFont(pdf, PDFOptions{}); err != nil {
		t.Fatalf("configureNativePDFFont: %v", err)
	}

	for _, layout := range everySpecialLayout() {
		diagram := smartart.NewSmartArt(layout).
			Position(styling.Inches(0.2), styling.Inches(0.2)).
			Size(styling.Inches(1.2), styling.Inches(0.8)).
			AddNode(smartart.NewNode("One").WithChild(smartart.NewNode("One a"))).
			AddNode(smartart.NewNode("Two")).
			AddNode(smartart.NewNode("Three")).
			AddNode(smartart.NewNode("Four"))

		if !renderPDFSmartArtSpecial(pdf, diagram) {
			t.Errorf("%s was not handled by a special renderer", layout.Name())
		}
	}
}

// everySpecialLayout is one layout per branch of renderPDFSmartArtSpecial, so a
// new per-layout renderer is exercised at a small frame size too.
func everySpecialLayout() []smartart.Layout {
	return []smartart.Layout{
		smartart.BasicBlockList,
		smartart.VerticalBlockList,
		smartart.HorizontalBulletLst,
		smartart.PictureAccentList,
		smartart.ContinuousBlockProcess,
		smartart.HorizontalHierarchy,
		smartart.LinearVenn,
		smartart.StackedVenn,
		smartart.BasicRadial,
		smartart.BasicMatrix,
		smartart.TitledMatrix,
		smartart.BasicPyramid,
		smartart.InvertedPyramid,
		smartart.PictureGrid,
	}
}

func TestSpecialLayoutsHandleAnEmptyDiagram(t *testing.T) {
	pdf := newTestPDF(t)
	if err := configureNativePDFFont(pdf, PDFOptions{}); err != nil {
		t.Fatalf("configureNativePDFFont: %v", err)
	}
	for _, layout := range everySpecialLayout() {
		diagram := smartart.NewSmartArt(layout).
			Position(styling.Inches(1), styling.Inches(1)).
			Size(styling.Inches(4), styling.Inches(3))
		if !renderPDFSmartArtSpecial(pdf, diagram) {
			t.Errorf("%s was not handled by a special renderer", layout.Name())
		}
	}
}
