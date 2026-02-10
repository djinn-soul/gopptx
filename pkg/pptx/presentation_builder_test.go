package pptx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPresentationBuilder(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "builder_test.pptx")

	// Test fluent API
	builder := NewPresentationBuilder("Fluent Presentation").
		WithMetadata(PresentationMetadata{Creator: "Test Builder"}).
		AddSlide(NewSlide("Slide 1").AddShape(NewRectangle(1, 1, 2, 2))).
		AddSlide(NewSlide("Slide 2").AddShape(NewEllipse(3, 1, 2, 2)))

	// Test Build
	data, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("Build returned empty data")
	}

	// Test WriteToFile
	if err := builder.WriteToFile(outPath); err != nil {
		t.Fatalf("WriteToFile failed: %v", err)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}

	// Test Build Error (empty slides?)
	// CreateWithMetadata requires at least one slide currently?
	// Let's check: pkg/pptx/presentation.go:406 says "at least one slide is required" or similar?
	// Actually, CreateWithMetadata checks if len(slides) == 0.

	emptyBuilder := NewPresentationBuilder("Empty")
	_, err = emptyBuilder.Build()
	if err == nil {
		t.Errorf("expected error for empty presentation, got nil")
	}
}
