package pptx_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

func TestPresentationBuilder(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "builder_test.pptx")

	builder := pptx.NewPresentationBuilder("Fluent Presentation").
		WithMetadata(pptx.PresentationMetadata{PresentationMetadata: common.PresentationMetadata{Creator: "Test Builder"}}).
		AddSlide(pptx.NewSlide("Slide 1").AddShape(pptx.NewRectangle(1, 1, 2, 2))).
		AddSlide(pptx.NewSlide("Slide 2").AddShape(pptx.NewEllipse(3, 1, 2, 2)))

	data, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("Build returned empty data")
	}

	if err := builder.WriteToFile(outPath); err != nil {
		t.Fatalf("WriteToFile failed: %v", err)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}

	emptyBuilder := pptx.NewPresentationBuilder("Empty")
	_, err = emptyBuilder.Build()
	if err == nil {
		t.Errorf("expected error for empty presentation, got nil")
	}
}
