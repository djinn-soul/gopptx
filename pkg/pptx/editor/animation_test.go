package editor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestPresentationEditorUpdateSlideRendersAnimations(t *testing.T) {
	path := writeDeckFixture(t, "base.pptx", []elements.SlideContent{
		elements.NewSlide("Base").AddBullet("seed"),
	})

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	updated := elements.NewSlide("Animated").
		AddShape(shapes.NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(animations.NewAnimation(1, animations.AnimationEntranceFade))
	if updateErr := editor.UpdateSlide(0, updated); updateErr != nil {
		t.Fatalf("update slide: %v", updateErr)
	}

	outPath := filepath.Join(t.TempDir(), "edited-animated.pptx")
	if saveErr := editor.Save(outPath); saveErr != nil {
		t.Fatalf("save edited deck: %v", saveErr)
	}

	slideXML := string(readZipFileBytes(t, outPath, "ppt/slides/slide1.xml"))
	if !strings.Contains(slideXML, "<p:timing>") {
		t.Fatalf("expected animation timing XML in updated slide")
	}
	if !strings.Contains(slideXML, `spid="3"`) {
		t.Fatalf("expected animation target shape id 3 in updated slide XML")
	}
}
