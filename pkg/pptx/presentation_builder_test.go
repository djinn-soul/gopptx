package pptx

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestPresentationBuilder_FluentAPI(t *testing.T) {
	builder := NewPresentationBuilder("Fluent Test").
		WithSlideNumbers(true).
		WithFooter("Footer Text").
		WithDateTime(true).
		WithModifyPassword("secret").
		WithMarkAsFinal(true).
		WithSignaturesEnabled(true)

	if builder.footerText != "Footer Text" {
		t.Errorf("expected footer text 'Footer Text', got %q", builder.footerText)
	}
	if !builder.showSlideNumber {
		t.Error("expected showSlideNumber to be true")
	}
	if !builder.showDateTime {
		t.Error("expected showDateTime to be true")
	}
	if builder.metadata.Protection.ModifyPassword != "secret" {
		t.Error("expected modify password to be set")
	}
	if !builder.metadata.Protection.MarkAsFinal {
		t.Error("expected mark as final to be true")
	}
	if !builder.metadata.Protection.SignaturesEnabled {
		t.Error("expected signatures enabled to be true")
	}

	// Test slide adding helpers
	builder.AddTitleSlide("Title Only").
		AddBulletSlide("Bullets", []string{"B1", "B2"}).
		AddShapesSlide("Shapes", NewShape(ShapeTypeRectangle, 0, 0, 100, 100))

	if len(builder.slides) != 3 {
		t.Errorf("expected 3 slides, got %d", len(builder.slides))
	}

	// Test metadata helpers
	builder.WithMetadata(Metadata{Metadata: MetadataFields{Creator: "Creator"}}).
		WithSlideSize(SlideSize16x9()).
		WithTheme(styling.Theme{}).
		WithMaster(&elements.SlideMaster{})

	if builder.metadata.Creator != "Creator" {
		t.Error("WithMetadata failed")
	}
	if builder.metadata.SlideSize.Width != 12192000 {
		t.Errorf("WithSlideSize failed, got width %d", builder.metadata.SlideSize.Width)
	}
	if builder.metadata.Theme == nil {
		t.Error("WithTheme failed")
	}
	if builder.metadata.Master == nil {
		t.Error("WithMaster failed")
	}
}

func TestPresentationBuilder_Build(t *testing.T) {
	builder := NewPresentationBuilder("").
		WithMetadata(Metadata{Metadata: MetadataFields{Title: "Meta Title"}}).
		AddSlide(NewSlide("Slide 1"))

	data, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("Build produced empty data")
	}

	// Title from metadata should take precedence
	if builder.metadata.Title != "Meta Title" {
		t.Errorf("expected title 'Meta Title', got %q", builder.metadata.Title)
	}
}
