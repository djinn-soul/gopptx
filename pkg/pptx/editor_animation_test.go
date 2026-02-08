package pptx

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestPresentationEditorUpdateSlideRendersAnimations(t *testing.T) {
	path := writeDeckFixture(t, "base.pptx", []SlideContent{
		NewSlide("Base").AddBullet("seed"),
	})

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}

	updated := NewSlide("Animated").
		AddShape(NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(NewAnimation(1, AnimationEntranceFade))
	if err := editor.UpdateSlide(0, updated); err != nil {
		t.Fatalf("update slide: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "edited-animated.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save edited deck: %v", err)
	}

	slideXML := string(readZipFileBytes(t, outPath, "ppt/slides/slide1.xml"))
	if !strings.Contains(slideXML, "<p:timing>") {
		t.Fatalf("expected animation timing XML in updated slide")
	}
	if !strings.Contains(slideXML, `spid="3"`) {
		t.Fatalf("expected animation target shape id 3 in updated slide XML")
	}
}
