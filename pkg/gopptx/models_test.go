package gopptx

import (
	"encoding/xml"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
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

func TestSlideNotesRich(t *testing.T) {
	slide := &Slide{}
	p1 := elements.NewParagraph()
	p1.Runs = append(p1.Runs, elements.NewRun("Line 1"))
	p2 := elements.NewParagraph()
	p2.Runs = append(p2.Runs, elements.NewRun("Line 2"))
	
	slide.SetRichNotes([]elements.Paragraph{p1, p2})
	if slide.Notes != "Line 1\nLine 2" {
		t.Errorf("Expected combined notes %q, got %q", "Line 1\nLine 2", slide.Notes)
	}
	if len(slide.NotesBody) != 2 {
		t.Errorf("Expected 2 paragraphs in NotesBody, got %d", len(slide.NotesBody))
	}

	p3 := elements.NewParagraph()
	p3.Runs = append(p3.Runs, elements.NewRun("Line 3"))
	slide.AddNoteParagraph(p3)
	if slide.Notes != "Line 1\nLine 2\nLine 3" {
		t.Errorf("Expected appended notes %q, got %q", "Line 1\nLine 2\nLine 3", slide.Notes)
	}
	if len(slide.NotesBody) != 3 {
		t.Errorf("Expected 3 paragraphs in NotesBody, got %d", len(slide.NotesBody))
	}
}

func TestSlideAddAnimationDefinition(t *testing.T) {
	slide := &Slide{}
	slide.AddAnimationDefinition(nil) // Should not panic or add
	if len(slide.Animations) != 0 {
		t.Errorf("Expected 0 animations after nil AddAnimationDefinition, got %d", len(slide.Animations))
	}

	def := animations.NewAnimation(1, animations.AnimationEmphasisGrowShrink)
	slide.AddAnimationDefinition(def)
	if len(slide.Animations) != 1 {
		t.Errorf("Expected 1 animation, got %d", len(slide.Animations))
	}
}

func TestSlideAddAnimationSequenceEdgeCases(t *testing.T) {
	slide := &Slide{}
	slide.AddAnimationSequence(100) // Empty sequence
	if len(slide.Animations) != 0 {
		t.Errorf("Expected 0 animations for empty sequence, got %d", len(slide.Animations))
	}

	def1 := animations.NewAnimation(1, animations.AnimationEntranceFade)
	slide.AddAnimationSequence(100, nil, def1) // Nil first element
	if len(slide.Animations) != 1 {
		t.Errorf("Expected 1 animation, got %d", len(slide.Animations))
	}
	// The first non-nil becomes the start of sequence, effectively index 1 for logic but 0 for non-nil elements
	if slide.Animations[0].Trigger != animations.AnimationAfterPrevious {
		t.Errorf("Expected AfterPrevious for def1 after nil, got %v", slide.Animations[0].Trigger)
	}
}

func TestPresentation_SaveErrors(t *testing.T) {
	pres := &Presentation{}
	err := pres.Save("test.pptx")
	if err == nil || err.Error() != "at least one slide is required" {
		t.Errorf("Expected 'at least one slide is required' error, got %v", err)
	}

	pres.AddSlide()
	err = pres.Save("") // Invalid path
	if err == nil {
		t.Error("Expected error for empty save path, got nil")
	}

	// Test with nil slide
	pres.Slides = append(pres.Slides, nil)
	filename := filepath.Join(t.TempDir(), "nil_slide.pptx")
	if err := pres.Save(filename); err != nil {
		t.Errorf("Failed to save with nil slide: %v", err)
	}
}

func TestPresentation_SaveFilePermissionError(t *testing.T) {
	// This might be OS dependent, but we try to trigger a permission error if possible
	// Or at least a path that doesn't exist/uncreatable
	pres := &Presentation{}
	pres.AddSlide()
	err := pres.Save("/invalid/path/to/presentation.pptx")
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}
