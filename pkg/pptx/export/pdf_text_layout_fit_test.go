package export

import (
	"strings"
	"testing"

	"github.com/signintech/gopdf"
)

// linearFitReference is the one-point-at-a-time scan the bisecting fitter
// replaced. The two must agree: it is the definition of the right answer.
func linearFitReference(
	pdf *gopdf.GoPdf,
	text string,
	initialSize, minSize int,
	maxWidth, maxHeight float64,
) int {
	size := initialSize
	for size > minSize {
		setPDFTextFontWithHint(pdf, size, false, false, "")
		lines := wrapPDFTextWithMetrics(pdf, text, maxWidth)
		if float64(len(lines))*pdfLineHeight(size) <= maxHeight {
			return size
		}
		size--
	}
	return size
}

func newLayoutTestPDF(t *testing.T) *gopdf.GoPdf {
	t.Helper()
	pdf := newTestPDF(t)
	if err := configureNativePDFFont(pdf, PDFOptions{}); err != nil {
		t.Fatalf("configureNativePDFFont: %v", err)
	}
	return pdf
}

func TestBisectingFitAgreesWithTheLinearScan(t *testing.T) {
	pdf := newLayoutTestPDF(t)

	texts := []string{
		"Short",
		"A rather longer caption that has to wrap at least once inside its box",
		strings.Repeat("wrap me ", 40),
		"日本語のテキストはスペースがないので一語として折り返される",
	}
	boxes := []struct{ w, h float64 }{
		{200, 40}, {120, 200}, {300, 18}, {80, 500},
	}

	for _, text := range texts {
		for _, box := range boxes {
			want := linearFitReference(pdf, text, 44, 10, box.w, box.h)
			got := fitPDFTextToBoxWithMetrics(pdf, text, 44, 10, false, false, box.w, box.h, "")
			if got != want {
				t.Errorf("text %.20q in %vx%v: got %d, want %d", text, box.w, box.h, got, want)
			}
		}
	}
}

func TestBisectingFitKeepsTheDocumentOnTheChosenSize(t *testing.T) {
	pdf := newLayoutTestPDF(t)
	text := strings.Repeat("overflowing text ", 30)

	size := fitPDFTextToBoxWithMetrics(pdf, text, 44, 10, false, false, 100, 30, "")
	lines := wrapPDFTextWithMetrics(pdf, text, 100)
	if float64(len(lines))*pdfLineHeight(size) > 30 && size > 10 {
		t.Errorf("chosen size %d still overflows the box", size)
	}
}

func TestBreakLongTokenKeepsPartsInsideTheWidth(t *testing.T) {
	pdf := newLayoutTestPDF(t)
	setPDFTextFontWithHint(pdf, 18, false, false, "")

	token := strings.Repeat("日本語のテキスト", 20)
	parts := breakLongToken(pdf, token, 120)
	if len(parts) < 2 {
		t.Fatalf("got %d parts, want the token broken across lines", len(parts))
	}
	if strings.Join(parts, "") != token {
		t.Error("breaking the token lost or reordered its text")
	}
	for i, part := range parts {
		if measuredWidth(pdf, part) > 120 {
			t.Errorf("part %d is %v wide, over the 120 limit", i, measuredWidth(pdf, part))
		}
	}
}

func TestBreakLongTokenKeepsAnUnbreakableRune(t *testing.T) {
	pdf := newLayoutTestPDF(t)
	setPDFTextFontWithHint(pdf, 18, false, false, "")

	// A box narrower than one glyph still has to make progress rather than loop.
	parts := breakLongToken(pdf, "abcd", 0.0001)
	if len(parts) != 4 {
		t.Fatalf("got %d parts, want one per rune", len(parts))
	}
	if strings.Join(parts, "") != "abcd" {
		t.Errorf("got %q, want abcd", strings.Join(parts, ""))
	}
}

// legacyBreakLongToken is the prefix-re-measuring break the current one
// replaced, kept here only so the benchmark can show what the change bought.
func legacyBreakLongToken(pdf *gopdf.GoPdf, token string, maxWidth float64) []string {
	parts := make([]string, 0, 8)
	var b strings.Builder
	for _, r := range token {
		next := b.String() + string(r)
		if measuredWidth(pdf, next) <= maxWidth || b.Len() == 0 {
			b.WriteRune(r)
			continue
		}
		parts = append(parts, b.String())
		b.Reset()
		b.WriteRune(r)
	}
	if b.Len() > 0 {
		parts = append(parts, b.String())
	}
	return parts
}

func BenchmarkBreakLongTokenLegacy(b *testing.B) {
	pdf, paragraph := benchmarkWrapFixture(b)
	for b.Loop() {
		_ = legacyBreakLongToken(pdf, paragraph, 400)
	}
}

func BenchmarkBreakLongToken(b *testing.B) {
	pdf, paragraph := benchmarkWrapFixture(b)
	for b.Loop() {
		_ = breakLongToken(pdf, paragraph, 400)
	}
}

func benchmarkWrapFixture(b *testing.B) (*gopdf.GoPdf, string) {
	b.Helper()
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 720, H: 540}})
	pdf.AddPage()
	b.Cleanup(func() { releaseDocumentFonts(pdf) })
	if err := configureNativePDFFont(pdf, PDFOptions{}); err != nil {
		b.Fatalf("configureNativePDFFont: %v", err)
	}
	setPDFTextFontWithHint(pdf, 18, false, false, "")
	return pdf, strings.Repeat("日本語のテキストはスペースがない", 200)
}

func BenchmarkWrapLongCJKParagraph(b *testing.B) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 720, H: 540}})
	pdf.AddPage()
	defer releaseDocumentFonts(pdf)
	if err := configureNativePDFFont(pdf, PDFOptions{}); err != nil {
		b.Fatalf("configureNativePDFFont: %v", err)
	}
	setPDFTextFontWithHint(pdf, 18, false, false, "")
	paragraph := strings.Repeat("日本語のテキストはスペースがない", 200)

	for b.Loop() {
		_ = wrapPDFTextWithMetrics(pdf, paragraph, 400)
	}
}
