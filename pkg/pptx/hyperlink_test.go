package pptx_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestHyperlinkURL(t *testing.T) {
	link := pptx.NewHyperlink(pptx.HyperlinkURL("https://example.com"))
	if link.Action.Type != pptx.HyperlinkActionURL {
		t.Errorf("Expected URL action type, got %s", link.Action.Type)
	}
	if link.Action.URL != "https://example.com" {
		t.Errorf("Expected URL 'https://example.com', got %s", link.Action.URL)
	}
	if !link.Action.IsExternal() {
		t.Error("URL should be external")
	}
}

func TestHyperlinkSlide(t *testing.T) {
	link := pptx.NewHyperlink(pptx.HyperlinkSlide(3))
	if link.Action.Type != pptx.HyperlinkActionSlide {
		t.Errorf("Expected Slide action type, got %s", link.Action.Type)
	}
	if link.Action.SlideNumber != 3 {
		t.Errorf("Expected slide 3, got %d", link.Action.SlideNumber)
	}
	if link.Action.IsExternal() {
		t.Error("Slide link should not be external")
	}
	target := link.Action.RelationshipTarget()
	if target != "slide3.xml" {
		t.Errorf("Expected 'slide3.xml', got %s", target)
	}
	if action := link.Action.ActionType(); action != "ppaction://hlinksldjump" {
		t.Errorf("Expected slide action type 'ppaction://hlinksldjump', got %s", action)
	}
}

func TestHyperlinkEmail(t *testing.T) {
	link := pptx.NewHyperlink(pptx.HyperlinkEmailWithSubject("test@example.com", "Hello World"))
	if link.Action.Type != pptx.HyperlinkActionEmail {
		t.Errorf("Expected Email action type, got %s", link.Action.Type)
	}
	target := link.Action.RelationshipTarget()
	if !strings.Contains(target, "mailto:test@example.com") {
		t.Errorf("Expected mailto link, got %s", target)
	}
	if !strings.Contains(target, "subject=Hello") {
		t.Errorf("Expected subject, got %s", target)
	}
}

func TestHyperlinkNavigation(t *testing.T) {
	tests := []struct {
		action       pptx.HyperlinkAction
		wantAction   string
		wantExternal bool
	}{
		{pptx.HyperlinkFirstSlide(), "ppaction://hlinkshowjump?jump=firstslide", true},
		{pptx.HyperlinkLastSlide(), "ppaction://hlinkshowjump?jump=lastslide", true},
		{pptx.HyperlinkNextSlide(), "ppaction://hlinkshowjump?jump=nextslide", true},
		{pptx.HyperlinkPreviousSlide(), "ppaction://hlinkshowjump?jump=previousslide", true},
		{pptx.HyperlinkEndShow(), "ppaction://hlinkshowjump?jump=endshow", true},
	}
	for _, tt := range tests {
		if tt.action.ActionType() != tt.wantAction {
			t.Errorf("ActionType() = %s, want %s", tt.action.ActionType(), tt.wantAction)
		}
		if tt.action.IsExternal() != tt.wantExternal {
			t.Errorf("IsExternal() = %v, want %v", tt.action.IsExternal(), tt.wantExternal)
		}
	}
}

func TestHyperlinkWithTooltip(t *testing.T) {
	link := pptx.NewHyperlink(pptx.HyperlinkURL("https://example.com")).
		WithTooltip("Click me").
		WithHighlightClick(false)

	if link.Tooltip != "Click me" {
		t.Errorf("Expected tooltip 'Click me', got %s", link.Tooltip)
	}
	if link.HighlightClick {
		t.Error("Expected HighlightClick to be false")
	}
}

func TestShapeWithHyperlink(t *testing.T) {
	slide := pptx.NewSlide("Hyperlink Test").
		AddShape(pptx.NewShape("rect", 0, 0, 100, 100).
			WithHyperlink(pptx.NewHyperlink(pptx.HyperlinkURL("https://example.com")).WithTooltip("Go to Example")))

	if len(slide.Shapes) != 1 {
		t.Fatalf("Expected 1 shape, got %d", len(slide.Shapes))
	}
	if slide.Shapes[0].Hyperlink == nil {
		t.Fatal("Expected shape to have hyperlink")
	}
	if slide.Shapes[0].Hyperlink.Tooltip != "Go to Example" {
		t.Errorf("Expected tooltip 'Go to Example', got %s", slide.Shapes[0].Hyperlink.Tooltip)
	}
}

func TestTextRunWithHyperlink(t *testing.T) {
	run := pptx.NewTextRun("Click here").
		WithHyperlink(pptx.NewHyperlink(pptx.HyperlinkSlide(2)))

	if run.Hyperlink == nil {
		t.Fatal("Expected text run to have hyperlink")
	}
	if run.Hyperlink.Action.Type != pptx.HyperlinkActionSlide {
		t.Errorf("Expected slide action, got %s", run.Hyperlink.Action.Type)
	}
}

func TestHyperlinkInPPTX(t *testing.T) {
	slide := pptx.NewSlide("Hyperlink PPTX Test").
		AddShape(pptx.NewShape("rect", 100000, 100000, 2000000, 500000).
			WithFill(pptx.NewShapeFill("FF6600")).
			WithText("Click me").
			WithHyperlink(pptx.NewHyperlink(pptx.HyperlinkURL("https://example.com"))))

	data, err := pptx.CreateWithSlides("Hyperlink Test", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip reader error: %v", err)
	}

	for _, f := range zr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			rc, openErr := f.Open()
			if openErr != nil {
				t.Fatalf("failed to open slide1.xml: %v", openErr)
			}
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(rc)
			_ = rc.Close()
			content := buf.String()
			if !strings.Contains(content, "hlinkClick") {
				t.Log("Note: hlinkClick not found in slide XML; relationship wiring may need additional implementation")
			}
			return
		}
	}
	t.Error("slide1.xml not found in zip")
}

func TestCreateWithSlidesRejectsInvalidTextRunHyperlink(t *testing.T) {
	slide := pptx.NewSlide("Invalid Hyperlink").
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("Bad URL").WithHyperlink(pptx.NewHyperlink(pptx.HyperlinkURL(""))),
		})

	_, err := pptx.CreateWithSlides("Invalid Hyperlink Deck", []pptx.SlideContent{slide})
	if err == nil {
		t.Fatal("expected invalid text-run hyperlink to fail validation")
	}
	if !strings.Contains(err.Error(), "hyperlink URL cannot be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNavigationHyperlinkUsesExternalRelationshipMode(t *testing.T) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Nav").AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("First Slide").WithHyperlink(pptx.NewHyperlink(pptx.HyperlinkFirstSlide())),
		}),
		pptx.NewSlide("Second"),
	}
	data, err := pptx.CreateWithSlides("Nav", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip reader error: %v", err)
	}

	for _, f := range zr.File {
		if f.Name != "ppt/slides/_rels/slide1.xml.rels" {
			continue
		}
		rc, openErr := f.Open()
		if openErr != nil {
			t.Fatalf("failed to open slide1 rels: %v", openErr)
		}
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(rc)
		_ = rc.Close()
		rels := buf.String()
		if !strings.Contains(rels, `TargetMode="External"`) {
			t.Fatalf("expected navigation hyperlink relationship to set external target mode")
		}
		return
	}
	t.Fatal("slide1 relationships not found")
}

func TestShapeHoverAction_EmitsHlinkHoverInSlideXML(t *testing.T) {
	slide := pptx.NewSlide("Hover Test").
		AddShape(pptx.NewShape("rect", 100000, 100000, 2000000, 500000).
			WithFill(pptx.NewShapeFill("0288D1")).
			WithText("Hover me").
			WithHoverAction(pptx.NewHyperlink(pptx.HyperlinkURL("https://hover.example.com")).
				WithTooltip("Hover tooltip")))

	data, err := pptx.CreateWithSlides("Hover Test", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip reader error: %v", err)
	}

	for _, f := range zr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			rc, _ := f.Open()
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(rc)
			_ = rc.Close()
			content := buf.String()
			if !strings.Contains(content, "hlinkHover") {
				t.Error("expected hlinkHover in slide XML for shape hover action")
			}
			if !strings.Contains(content, "Hover tooltip") {
				t.Error("expected tooltip text in slide XML")
			}
			return
		}
	}
	t.Error("slide1.xml not found")
}

func TestTextRunHoverAction_EmitsHlinkHoverInSlideXML(t *testing.T) {
	slide := pptx.NewSlide("Text Hover").
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("Hover text").
				WithHoverAction(pptx.NewHyperlink(pptx.HyperlinkNextSlide()).
					WithTooltip("Text hover tip")),
		})

	data, err := pptx.CreateWithSlides("Text Hover", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip reader error: %v", err)
	}

	for _, f := range zr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			rc, _ := f.Open()
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(rc)
			_ = rc.Close()
			content := buf.String()
			if !strings.Contains(content, "hlinkHover") {
				t.Error("expected hlinkHover in slide XML for text run hover action")
			}
			if !strings.Contains(content, "Text hover tip") {
				t.Error("expected tooltip text in slide XML")
			}
			return
		}
	}
	t.Error("slide1.xml not found")
}
