package presentation

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
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

func TestWritePresentationPackage_TransitionSoundRIDMatchesRelationships(t *testing.T) {
	tmp := t.TempDir()
	soundPath := filepath.Join(tmp, "sound.wav")
	if err := os.WriteFile(soundPath, []byte("fake-audio"), 0o600); err != nil {
		t.Fatalf("write sound: %v", err)
	}

	meta := Metadata{}
	slides := []elements.SlideContent{
		elements.NewSlide("Slide 1").
			WithBarChart(charts.BarChart{Categories: []string{"A"}, Values: []float64{1}}).
			WithTransitionSound(soundPath),
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	if err := WritePresentationPackage(zw, meta, slides, len(slides)); err != nil {
		t.Fatalf("WritePresentationPackage failed: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zip close: %v", err)
	}

	slideXML := string(readZipPart(t, buf.Bytes(), "ppt/slides/slide1.xml"))
	relsXML := string(readZipPart(t, buf.Bytes(), "ppt/slides/_rels/slide1.xml.rels"))
	matches := regexp.MustCompile(`r:embed="(rId\d+)"`).FindStringSubmatch(slideXML)
	if len(matches) != 2 {
		t.Fatalf("expected transition sound embed relationship in slide xml, got: %s", slideXML)
	}
	soundRID := matches[1]
	if !strings.Contains(
		relsXML,
		fmt.Sprintf(`Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/audio"`, soundRID),
	) {
		t.Fatalf("expected %s to map to audio relationship, got rels: %s", soundRID, relsXML)
	}
}

func TestWritePresentationPackage_PlaceholderTableSpecError(t *testing.T) {
	badTable := tables.NewTable([]styling.Length{styling.Inches(1)}).AddStyledRow(
		[]tables.TableCell{
			tables.NewTableCell("A").WithRowSpan(2),
		},
	)
	slides := []elements.SlideContent{
		elements.NewSlide("Slide 1").WithPlaceholderTableAs(1, "tbl", badTable),
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	err := WritePresentationPackage(zw, Metadata{}, slides, len(slides))
	if err == nil {
		t.Fatal("expected placeholder table conversion error")
	}
}

func readZipPart(t *testing.T, blob []byte, name string) []byte {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(blob), int64(len(blob)))
	if err != nil {
		t.Fatalf("zip reader: %v", err)
	}
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		rc, openErr := f.Open()
		if openErr != nil {
			t.Fatalf("open %s: %v", name, openErr)
		}
		defer rc.Close()
		data, readErr := io.ReadAll(rc)
		if readErr != nil {
			t.Fatalf("read %s: %v", name, readErr)
		}
		return data
	}
	t.Fatalf("zip part not found: %s", name)
	return nil
}
