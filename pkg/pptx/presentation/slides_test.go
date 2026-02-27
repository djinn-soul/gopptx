package presentation

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestLayoutTargetForMaster(t *testing.T) {
	tests := []struct {
		baseTarget string
		masterNum  int
		expected   string
	}{
		{"../slideLayouts/slideLayout1.xml", 1, "../slideLayouts/slideLayout1.xml"},
		{"../slideLayouts/slideLayout1.xml", 2, "../slideLayouts/slideLayout7.xml"},
		{"../slideLayouts/slideLayout2.xml", 2, "../slideLayouts/slideLayout8.xml"},
		{"../slideLayouts/slideLayout1.xml", 3, "../slideLayouts/slideLayout13.xml"},
		{"invalid", 2, "invalid"},
	}
	for _, tt := range tests {
		if got := layoutTargetForMaster(tt.baseTarget, tt.masterNum); got != tt.expected {
			t.Errorf("layoutTargetForMaster(%q, %d) = %q, want %q", tt.baseTarget, tt.masterNum, got, tt.expected)
		}
	}
}

func TestWritePresentationPackage_MultiMaster(t *testing.T) {
	meta := Metadata{
		Masters: []*elements.SlideMaster{
			elements.NewMaster(),
			elements.NewMaster(),
		},
	}

	slides := []elements.SlideContent{
		elements.NewSlide("Slide 1"),
		elements.NewSlide("Slide 2"),
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	err := WritePresentationPackage(zw, meta, slides, len(slides))
	if err != nil {
		t.Fatalf("WritePresentationPackage failed: %v", err)
	}
	zw.Close()

	validatePPTX(t, buf.Bytes())
}

func TestWritePresentationPackage_Comments(t *testing.T) {
	meta := Metadata{}
	slides := []elements.SlideContent{
		elements.NewSlide("Slide 1").AddComment("Author 1", "Comment 1"),
		elements.NewSlide("Slide 2").AddComment("Author 2", "Comment 2").AddComment("Author 1", "Comment 3"),
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	err := WritePresentationPackage(zw, meta, slides, len(slides))
	if err != nil {
		t.Fatalf("WritePresentationPackage failed: %v", err)
	}
	zw.Close()

	validatePPTX(t, buf.Bytes())
}

func TestWritePresentationPackage_Notes(t *testing.T) {
	meta := Metadata{}
	slides := []elements.SlideContent{
		elements.NewSlide("Slide 1").WithNotes("Speaker notes"),
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	err := WritePresentationPackage(zw, meta, slides, len(slides))
	if err != nil {
		t.Fatalf("WritePresentationPackage failed: %v", err)
	}
	zw.Close()

	validatePPTX(t, buf.Bytes())
}

func TestWritePresentationPackage_Footer(t *testing.T) {
	meta := Metadata{
		Metadata: common.Metadata{
			FooterText: "Global Footer",
		},
	}
	slides := []elements.SlideContent{
		elements.NewSlide("Slide 1"),
		elements.NewSlide("Slide 2").WithLayout(elements.SlideLayoutBlank),
	}
	slides[1].FooterText = "Override Footer"

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	err := WritePresentationPackage(zw, meta, slides, len(slides))
	if err != nil {
		t.Fatalf("WritePresentationPackage failed: %v", err)
	}
	zw.Close()

	validatePPTX(t, buf.Bytes())
}
