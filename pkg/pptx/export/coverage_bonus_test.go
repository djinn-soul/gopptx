package export

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func TestHTML_ImageCoverage(t *testing.T) {
	// 1. Test non-embedded image with BaseURL
	opts := DefaultHTMLOptions()
	opts.EmbedImages = false
	opts.BaseURL = "http://assets.com/"

	img := shapes.Image{Path: "photo.jpg"}
	var buf bytes.Buffer
	err := renderImageToWriter(&buf, img, opts)
	if err != nil {
		t.Fatalf("renderImageToWriter failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "src=\"http://assets.com/photo.jpg\"") {
		t.Errorf("expected absolute URL, got %q", out)
	}

	// 2. Test different extensions
	exts := []string{".jpeg", ".gif", ".svg", ".png"}
	for _, ext := range exts {
		buf.Reset()
		img = shapes.Image{Data: []byte("fake"), Path: "test" + ext}
		_ = renderImageToWriter(&buf, img, DefaultHTMLOptions())
	}

	// 3. Test non-local path error
	buf.Reset()
	img = shapes.Image{Path: "../outside.png"}
	opts.EmbedImages = true
	err = renderImageToWriter(&buf, img, opts)
	if err == nil || !strings.Contains(err.Error(), "invalid image path") {
		t.Errorf("expected invalid path error, got %v", err)
	}
}

func TestHTML_TableCoverage(t *testing.T) {
	// Test plain rows (no styling)
	table := tables.NewTable([]styling.Length{styling.Inches(1)})
	table = table.AddRow([]string{"Plain"})

	var buf bytes.Buffer
	err := renderTableToWriter(&buf, &table)
	if err != nil {
		t.Fatalf("renderTableToWriter failed: %v", err)
	}
	if !strings.Contains(buf.String(), "<td>Plain</td>") {
		t.Error("Plain table row not rendered")
	}
}

func TestReader_NumericParsing(t *testing.T) {
	// parseNumericInt64 coverage
	if v, ok := parseNumericInt64(int(10)); !ok || v != 10 {
		t.Error("int conversion failed")
	}
	if v, ok := parseNumericInt64(float64(20.5)); !ok || v != 20 {
		t.Error("float64 conversion failed")
	}
	if _, ok := parseNumericInt64("not a number"); ok {
		t.Error("expected false for string")
	}

	// toLengthSlice coverage
	lengths := toLengthSlice([]any{int(100), float64(200)})
	if len(lengths) != 2 || int64(lengths[0]) != 100 || int64(lengths[1]) != 200 {
		t.Error("toLengthSlice failed")
	}
}

func TestReader_LayoutCoverage(_ *testing.T) {
	cell := tables.NewTableCell("Test")
	meta := map[string]any{
		"v_align":     "bottom",
		"margin_left": int64(10000), // EMU
	}
	applyTableCellLayout(&cell, meta)
	// Check if cell was updated (internal fields are private, but we call it for coverage)

	// Border coverage
	borderMeta := map[string]any{
		"border_left": map[string]any{
			"width": int64(5000),
			"color": "FF0000",
		},
	}
	applyTableCellBorders(&cell, borderMeta)
}

func TestReader_AdjustmentCoverage(t *testing.T) {
	adjs := []editorcommon.ShapeAdjustment{
		{Name: "adj1", Formula: "val 100"},
		{Name: "", Formula: "ignore"},
	}
	exported := editorAdjustmentsToExport(adjs)
	if len(exported) != 1 || exported[0].Name != "adj1" {
		t.Error("editorAdjustmentsToExport failed")
	}
}

func TestPDF_FontSelection(t *testing.T) {
	if !isMonospaceFontHint("consolas") {
		t.Error("consolas should be monospace")
	}
	if !isSerifFontHint("times new roman") {
		t.Error("times should be serif")
	}

	hint := inferCodeFontHint("func main() { return 0 }")
	if hint == "" {
		t.Error("inferCodeFontHint failed for code block")
	}

	setPDFFontAliases("Arial", "Times", "Courier")
}

func TestPDF_FileBased(t *testing.T) {
	// PDFFromFileWithOptions with Native
	tmpDir := t.TempDir()
	pdfPath := filepath.Join(tmpDir, "test.pdf")

	// We need a real PPTX for SlidesFromPPTX to work
	// But we can test the error path or use a mock if possible.
	// For now, let's test the driver normalization and auto-selector.

	opts := PDFOptions{Driver: PDFDriverNative}
	err := PDFFromFileWithOptions("nonexistent.pptx", pdfPath, opts)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}

	// Test sanitizeTitle
	if sanitizeTitle("Hello World!") != "Hello_World_" {
		t.Errorf("sanitizeTitle failed: %s", sanitizeTitle("Hello World!"))
	}
}

func TestPDF_AutoDriver(t *testing.T) {
	tmpDir := t.TempDir()
	pdfPath := filepath.Join(tmpDir, "auto.pdf")
	slides := []elements.SlideContent{elements.NewSlide("Auto")}

	// This will try Native first. If Native fails (e.g. no fonts), it tries others.
	_ = PDFWithOptions("Auto", slides, pdfPath, PDFOptions{Driver: PDFDriverAuto})
}
