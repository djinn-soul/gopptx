package editor

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestGetSlideBackgroundXMLReturnsActualSubtree(t *testing.T) {
	base := writeDeckFixture(t, "background-read-bridge.pptx", []elements.SlideContent{
		elements.NewSlide("Background"),
	})
	editor, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	if err := editor.SetSlideBackground(0, SlideBackground{
		Type:  "solid",
		Color: "3070B3",
	}); err != nil {
		t.Fatalf("set background: %v", err)
	}

	backgroundXML, err := editor.GetSlideBackgroundXML(0)
	if err != nil {
		t.Fatalf("get background XML: %v", err)
	}
	for _, want := range []string{
		"<p:bg>", "<a:solidFill>", `<a:srgbClr val="3070B3"/>`, "</p:bg>",
	} {
		if !strings.Contains(backgroundXML, want) {
			t.Fatalf("background XML missing %s: %s", want, backgroundXML)
		}
	}
	if strings.Contains(backgroundXML, "<p:spTree") {
		t.Fatalf("background read leaked the common slide container: %s", backgroundXML)
	}
}

func TestExtractBgXMLAcceptsBackgroundAttributes(t *testing.T) {
	source := `<p:sld><p:cSld><p:bg bwMode="gray"><p:bgPr/></p:bg><p:spTree/></p:cSld></p:sld>`
	got := extractBgXML(source)
	if got != `<p:bg bwMode="gray"><p:bgPr/></p:bg>` {
		t.Fatalf("unexpected background subtree: %s", got)
	}
}
