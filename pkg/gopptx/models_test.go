package gopptx

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

func TestPresentationXML(t *testing.T) {
	pres := &Presentation{}
	output, err := xml.Marshal(pres)
	if err != nil {
		t.Fatalf("Failed to marshal Presentation: %v", err)
	}

	expected := `<presentation xmlns="http://schemas.openxmlformats.org/presentationml/2006/main"><sldIdLst></sldIdLst></presentation>`
	if string(output) != expected {
		t.Errorf("Expected %s, got %s", expected, string(output))
	}
}

func TestSlideXML(t *testing.T) {
	slide := &Slide{}
	output, err := xml.Marshal(slide)
	if err != nil {
		t.Fatalf("Failed to marshal Slide: %v", err)
	}

	expected := `<sld xmlns="http://schemas.openxmlformats.org/presentationml/2006/main"></sld>`
	if string(output) != expected {
		t.Errorf("Expected %s, got %s", expected, string(output))
	}
}

func TestAddSlide(t *testing.T) {
	pres := &Presentation{}
	slide := pres.AddSlide()
	if slide == nil {
		t.Fatal("Expected new slide, got nil")
	}
	if len(pres.Slides) != 1 {
		t.Errorf("Expected 1 slide, got %d", len(pres.Slides))
	}
}

func TestSlideAddBullet(t *testing.T) {
	slide := &Slide{}
	slide.AddBullet("First bullet")
	if len(slide.Bullets) != 1 {
		t.Fatalf("Expected 1 bullet, got %d", len(slide.Bullets))
	}
	if slide.Bullets[0] != "First bullet" {
		t.Errorf("Unexpected bullet, got %q", slide.Bullets[0])
	}
}

func TestSlideContentCopiesExtras(t *testing.T) {
	slide := &Slide{}
	img := shapes.NewImage("image.png", styling.Emu(0), styling.Emu(0), styling.Emu(100), styling.Emu(100))
	slide.AddImage(img)
	rect := shapes.NewRectangle(0, 0, 1, 1)
	slide.AddShape(rect)
	conn := shapes.NewStraightConnector(styling.Emu(0), styling.Emu(0), styling.Emu(100), styling.Emu(100))
	slide.AddConnector(conn)
	slide.SetTransition(transitions.TransitionOptions{Type: transitions.TransitionFade})

	content := slide.toSlideContent(0)
	if len(content.Images) != 1 {
		t.Fatalf("Expected 1 image, got %d", len(content.Images))
	}
	if len(content.Shapes) != 1 {
		t.Fatalf("Expected 1 shape, got %d", len(content.Shapes))
	}
	if len(content.Connectors) != 1 {
		t.Fatalf("Expected 1 connector, got %d", len(content.Connectors))
	}
	if content.Transition == nil {
		t.Fatal("Expected transition to be preserved")
	}
	if xml := content.Transition.XML(); xml == "" {
		t.Fatal("Expected transition XML to be non-empty")
	}
}

func TestSlideNotesPlaceholdersAnimations(t *testing.T) {
	slide := &Slide{}
	slide.SetNotes("Speaker notes")
	slide.AddPlaceholderText(0, "Title override")
	img := shapes.NewImage("cover.png", styling.Emu(0), styling.Emu(0), styling.Emu(100), styling.Emu(100))
	slide.AddPlaceholderImage(1, img)
	anim := animations.NewAnimation(1, animations.AnimationEntranceFade)
	slide.AddAnimation(anim)

	content := slide.toSlideContent(0)
	if content.Notes != "Speaker notes" {
		t.Fatalf("Expected notes text preserved, got %q", content.Notes)
	}
	if len(content.NotesBody) != 1 {
		t.Fatalf("Expected 1 notes paragraph, got %d", len(content.NotesBody))
	}
	if len(content.PlaceholderOverrides) != 2 {
		t.Fatalf("Expected 2 placeholder overrides, got %d", len(content.PlaceholderOverrides))
	}
	if content.PlaceholderOverrides[1].Image == nil {
		t.Fatal("Expected image override payload to be preserved")
	}
	if len(content.Animations) != 1 {
		t.Fatalf("Expected 1 animation, got %d", len(content.Animations))
	}
	if content.Animations[0].Effect != animations.AnimationEntranceFade {
		t.Fatalf("Unexpected animation effect %q", content.Animations[0].Effect)
	}
}

func TestSlideTableChartSmartArtAndSequence(t *testing.T) {
	slide := &Slide{}
	table := tables.NewTable([]styling.Length{styling.Inches(3)})
	table = table.AddRow([]string{"A", "B"})
	slide.SetTable(table)
	slide.SetBarChart(charts.NewBarChart([]string{"Series"}, []float64{1}))
	slide.AddSmartArt(smartart.NewSmartArt(smartart.BasicBlockList))
	seqAnim1 := animations.NewAnimation(1, animations.AnimationEntranceFade)
	seqAnim2 := animations.NewAnimation(1, animations.AnimationEntranceZoom)
	slide.AddAnimationSequence(120, seqAnim1, seqAnim2)

	content := slide.toSlideContent(0)
	if content.Table == nil {
		t.Fatal("Expected table to survive conversion")
	}
	if content.Chart == nil {
		t.Fatal("Expected chart to survive conversion")
	}
	if len(content.SmartArtDiagrams) != 1 {
		t.Fatalf("Expected 1 SmartArt diagram, got %d", len(content.SmartArtDiagrams))
	}
	if len(content.Animations) != 2 {
		t.Fatalf("Expected 2 sequenced animations, got %d", len(content.Animations))
	}
	if content.Animations[1].Trigger != animations.AnimationAfterPrevious {
		t.Fatalf("Expected second animation trigger to be after previous, got %q", content.Animations[1].Trigger)
	}
	if content.Animations[1].DelayMS != 120 {
		t.Fatalf("Expected second animation delay 120, got %d", content.Animations[1].DelayMS)
	}
}

func TestSlideAnimationSequenceRespectsCustomDelay(t *testing.T) {
	slide := &Slide{}
	anim := animations.NewAnimation(1, animations.AnimationEntranceFlyIn)
	anim.DelayMS = 50
	slide.AddAnimationSequence(200, anim)

	content := slide.toSlideContent(0)
	if len(content.Animations) != 1 {
		t.Fatalf("Expected 1 animation, got %d", len(content.Animations))
	}
	if content.Animations[0].DelayMS != 0 {
		t.Fatalf("Expected sequence to reset delay for animation 0, got %d", content.Animations[0].DelayMS)
	}
}

func TestPresentation_SaveWithSlideData(t *testing.T) {
	pres := &Presentation{Title: "Custom Deck"}
	slide := pres.AddSlide()
	slide.Title = "Cover Page"
	slide.AddBullet("Welcome to the deck")

	filename := filepath.Join(t.TempDir(), "slide_data.pptx")
	if err := pres.Save(filename); err != nil {
		t.Fatalf("Failed to save presentation with slide data: %v", err)
	}
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("Expected file %s to exist, but got error: %v", filename, err)
	}
	if info.Size() == 0 {
		t.Fatalf("Expected file %s to have content", filename)
	}
}
