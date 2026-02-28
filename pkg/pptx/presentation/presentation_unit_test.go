package presentation

import (
	"archive/zip"
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

func TestPresentation_EffectiveMasters(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		m := Metadata{}
		res := getEffectiveMasters(m)
		if len(res) != 1 { t.Error("expected 1 default master") }
	})

	t.Run("Single", func(t *testing.T) {
		m := Metadata{Master: elements.NewMaster()}
		res := getEffectiveMasters(m)
		if len(res) != 1 { t.Error("expected 1 master") }
	})

	t.Run("Multi", func(t *testing.T) {
		m := Metadata{Masters: []*elements.SlideMaster{elements.NewMaster(), elements.NewMaster()}}
		res := getEffectiveMasters(m)
		if len(res) != 2 { t.Error("expected 2 masters") }
	})
}

func TestPresentation_NotesThemeIndex(t *testing.T) {
	if getNotesThemeIndex(false, 1) != 0 { t.Error("expected 0") }
	if getNotesThemeIndex(true, 1) != 2 { t.Error("expected 2") }
}

func TestPresentation_ConvertSections(t *testing.T) {
	secs := []Section{
		{Name: "S1", SlideIndices: []int{0}},
	}

	t.Run("Valid", func(t *testing.T) {
		res, err := convertSections(secs, 1)
		if err != nil { t.Fatalf("failed: %v", err) }
		if len(res) != 1 || res[0].Name != "S1" { t.Error("conversion failed") }
	})

	t.Run("InvalidIndex", func(t *testing.T) {
		_, err := convertSections(secs, 0)
		if err == nil { t.Error("expected error for invalid index") }
	})
}

func TestPresentation_PrepareComments(t *testing.T) {
	meta := Metadata{Metadata: common.Metadata{GeneratedDate: time.Now()}}
	slides := []elements.SlideContent{
		elements.NewSlide("S1").AddComment("Author 1", "Text 1"),
		elements.NewSlide("S2").AddComment("Author 1", "Text 2").AddComment("Author 2", "Text 3"),
	}

	authors, cms, indices := prepareComments(meta, slides)
	if len(authors) != 2 { t.Errorf("expected 2 authors, got %d", len(authors)) }
	if len(indices) != 2 { t.Error("expected 2 slides with comments") }
	if len(cms[0]) != 1 || len(cms[1]) != 2 { t.Error("comment counts failed") }
}

func TestWritePresentationPackage_Full(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	meta := Metadata{
		Metadata: common.Metadata{
			Title: "Full Test",
			CustomXML: []common.CustomXMLPart{
				{RootElement: "root", Content: "<root/>"},
			},
		},
		Sections: []Section{
			{Name: "Section 1", SlideIndices: []int{0}},
		},
		RTL: true,
		VBA: &vba.VBAProject{},
	}

	slides := []elements.SlideContent{
		elements.NewSlide("S1").
			AddImage(shapes.Image{Data: []byte("fake"), Format: "png"}).
			WithBarChart(charts.BarChart{Categories: []string{"A"}, Values: []float64{1}}).
			AddSmartArt(smartart.NewSmartArt(smartart.BasicBlockList)).
			WithTransitionSound("sound.wav"),
	}

	// Create dummy sound file
	_ = os.WriteFile("sound.wav", []byte("dummy"), 0600)
	defer os.Remove("sound.wav")

	err := WritePresentationPackage(zw, meta, slides, 1)
	if err != nil {
		t.Fatalf("WritePresentationPackage failed: %v", err)
	}
	zw.Close()
}

func TestPresentation_SlideSize_Helpers(t *testing.T) {
	if GetSlideSize4x3().Width != 9144000 { t.Error("4x3 failed") }
	if GetSlideSize16x9().Width != 12192000 { t.Error("16x9 failed") }
}
