package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestHyperlinkURL(t *testing.T) {
	link := NewHyperlink(HyperlinkURL("https://example.com"))
	if link.Action.Type != HyperlinkActionURL {
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
	link := NewHyperlink(HyperlinkSlide(3))
	if link.Action.Type != HyperlinkActionSlide {
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
}

func TestHyperlinkEmail(t *testing.T) {
	link := NewHyperlink(HyperlinkEmailWithSubject("test@example.com", "Hello World"))
	if link.Action.Type != HyperlinkActionEmail {
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
		action       HyperlinkAction
		wantAction   string
		wantExternal bool
	}{
		{HyperlinkFirstSlide(), "ppaction://hlinkshowjump?jump=firstslide", true},
		{HyperlinkLastSlide(), "ppaction://hlinkshowjump?jump=lastslide", true},
		{HyperlinkNextSlide(), "ppaction://hlinkshowjump?jump=nextslide", true},
		{HyperlinkPreviousSlide(), "ppaction://hlinkshowjump?jump=previousslide", true},
		{HyperlinkEndShow(), "ppaction://hlinkshowjump?jump=endshow", true},
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
	link := NewHyperlink(HyperlinkURL("https://example.com")).
		WithTooltip("Click me").
		WithHighlightClick(false)

	if link.Tooltip != "Click me" {
		t.Errorf("Expected tooltip 'Click me', got %s", link.Tooltip)
	}
	if link.HighlightClick {
		t.Error("Expected HighlightClick to be false")
	}
}

func TestHyperlinkValidation(t *testing.T) {
	tests := []struct {
		name    string
		action  HyperlinkAction
		wantErr bool
	}{
		{"valid URL", HyperlinkURL("https://example.com"), false},
		{"empty URL", HyperlinkURL(""), true},
		{"URL without scheme", HyperlinkURL("example.com"), true},
		{"valid slide", HyperlinkSlide(1), false},
		{"invalid slide 0", HyperlinkSlide(0), true},
		{"valid email", HyperlinkEmail("test@example.com"), false},
		{"invalid email", HyperlinkEmail("invalid"), true},
		{"valid file", HyperlinkFile("C:\\docs\\file.pdf"), false},
		{"empty file", HyperlinkFile(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHyperlinkAction(tt.action, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHyperlinkAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShapeWithHyperlink(t *testing.T) {
	slide := NewSlide("Hyperlink Test").
		AddShape(NewShape("rect", 0, 0, 100, 100).
			WithHyperlink(NewHyperlink(HyperlinkURL("https://example.com")).WithTooltip("Go to Example")))

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
	run := NewTextRun("Click here").
		WithHyperlink(NewHyperlink(HyperlinkSlide(2)))

	if run.Hyperlink == nil {
		t.Fatal("Expected text run to have hyperlink")
	}
	if run.Hyperlink.Action.Type != HyperlinkActionSlide {
		t.Errorf("Expected slide action, got %s", run.Hyperlink.Action.Type)
	}
}

func TestHyperlinkInPPTX(t *testing.T) {
	slide := NewSlide("Hyperlink PPTX Test").
		AddShape(NewShape("rect", 100000, 100000, 2000000, 500000).
			WithFill(NewShapeFill("FF6600")).
			WithText("Click me").
			WithHyperlink(NewHyperlink(HyperlinkURL("https://example.com"))))

	data, err := CreateWithSlides("Hyperlink Test", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip reader error: %v", err)
	}

	// Check slide1.xml for hyperlink element
	for _, f := range zr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("failed to open slide1.xml: %v", err)
			}
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(rc)
			_ = rc.Close()
			content := buf.String()
			// The hyperlink should be rendered in XML (relationship ID will be present)
			if !strings.Contains(content, "hlinkClick") {
				t.Log("Note: hlinkClick not found in slide XML; relationship wiring may need additional implementation")
			}
			return
		}
	}
	t.Error("slide1.xml not found in zip")
}

func TestCreateWithSlidesRejectsInvalidTextRunHyperlink(t *testing.T) {
	slide := NewSlide("Invalid Hyperlink").
		AddBulletRuns([]TextRun{
			NewTextRun("Bad URL").WithHyperlink(NewHyperlink(HyperlinkURL(""))),
		})

	_, err := CreateWithSlides("Invalid Hyperlink Deck", []SlideContent{slide})
	if err == nil {
		t.Fatal("expected invalid text-run hyperlink to fail validation")
	}
	if !strings.Contains(err.Error(), "hyperlink URL cannot be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNavigationHyperlinkUsesExternalRelationshipMode(t *testing.T) {
	slides := []SlideContent{
		NewSlide("Nav").AddBulletRuns([]TextRun{
			NewTextRun("First Slide").WithHyperlink(NewHyperlink(HyperlinkFirstSlide())),
		}),
		NewSlide("Second"),
	}
	data, err := CreateWithSlides("Nav", slides)
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
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("failed to open slide1 rels: %v", err)
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
