package editor

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestAddImageFromBase64DataURI(t *testing.T) {
	base := writeDeckFixture(t, "base.pptx", []elements.SlideContent{elements.NewSlide("Slide 1")})
	ed, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	dataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(testutil.TinyPNG())
	shapeID, err := ed.AddImageFromBase64(0, dataURI, "", 10, 20, 100, 80, nil)
	if err != nil {
		t.Fatalf("add image from base64: %v", err)
	}

	meta, err := ed.GetImageMetadata(0, shapeID)
	if err != nil {
		t.Fatalf("get image metadata: %v", err)
	}
	if meta.Format != "png" {
		t.Fatalf("expected png format, got %q", meta.Format)
	}
}

func TestAddImageFromURL(t *testing.T) {
	base := writeDeckFixture(t, "base.pptx", []elements.SlideContent{elements.NewSlide("Slide 1")})
	ed, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	if _, err := ed.AddImageFromURL(0, "http://127.0.0.1/pixel.png", 10, 20, 100, 80, nil); err == nil {
		t.Fatal("expected private URL to be blocked")
	}
}

func TestAddImageFromBase64_EmbedsBlip(t *testing.T) {
	base := writeDeckFixture(t, "base.pptx", []elements.SlideContent{elements.NewSlide("Slide 1")})
	ed, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	shapeID, err := ed.AddImageFromBase64(
		0,
		"data:image/png;base64,"+base64.StdEncoding.EncodeToString(testutil.TinyPNG()),
		"",
		10,
		20,
		100,
		80,
		nil,
	)
	if err != nil {
		t.Fatalf("add image from base64: %v", err)
	}
	if shapeID == 0 {
		t.Fatal("expected non-zero shape id")
	}
	slideXML, ok := ed.parts.Get("ppt/slides/slide1.xml")
	if !ok {
		t.Fatal("slide part missing")
	}
	if !strings.Contains(string(slideXML), `a:blip r:embed="`) {
		t.Fatalf("expected embedded image blip in slide xml, got: %s", string(slideXML))
	}
}
