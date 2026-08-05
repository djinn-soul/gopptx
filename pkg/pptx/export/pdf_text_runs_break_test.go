package export

import (
	"strings"
	"testing"

	"github.com/signintech/gopdf"
)

// newWrapTestPDF is a document with fonts registered and one selected.
// Measuring goes through gopdf, which dereferences the current font, so a
// document without one panics rather than returning a width.
func newWrapTestPDF(t *testing.T) *gopdf.GoPdf {
	t.Helper()
	pdf := newTestPDF(t)
	if err := configureNativePDFFont(pdf, PDFOptions{}); err != nil {
		t.Skipf("no system font available to measure with: %v", err)
	}
	setPDFTextFontWithHint(pdf, 12, false, false, "")
	return pdf
}

// The plain-text path splits an over-wide token mid-word; the rich-text path
// did not, so a long label ran straight out of its shape instead of breaking
// the way PowerPoint does.
func TestWrapStyledRunsBreaksAnOverWideToken(t *testing.T) {
	pdf := newWrapTestPDF(t)
	run := pdfStyledRun{Text: "Subroutinesubroutinesubroutine", SizePt: 12}

	lines := wrapStyledRuns(pdf, []pdfStyledRun{run}, 40, nil)
	if len(lines) < 2 {
		t.Fatalf("got %d line(s), want the token split across several", len(lines))
	}
	var rebuilt strings.Builder
	for _, line := range lines {
		for _, r := range line {
			rebuilt.WriteString(r.Text)
		}
	}
	if rebuilt.String() != run.Text {
		t.Errorf("rebuilt %q want %q -- characters were lost in the split", rebuilt.String(), run.Text)
	}
}

// Splitting must not lose the run's styling.
func TestBreakLongStyledRunKeepsStyle(t *testing.T) {
	pdf := newWrapTestPDF(t)
	run := pdfStyledRun{Text: "abcdefghijklmnopqrstuvwxyz", SizePt: 12, Bold: true, Color: [3]uint8{1, 2, 3}}

	pieces := breakLongStyledRun(pdf, run, 20, nil)
	if len(pieces) == 0 {
		t.Fatal("no pieces returned")
	}
	for i, p := range pieces {
		if !p.Bold || p.Color != run.Color || p.SizePt != run.SizePt {
			t.Errorf("piece %d lost styling: %+v", i, p)
		}
	}
}

// It must always return something, so the caller can take the last piece as
// the remainder without a bounds check.
func TestBreakLongStyledRunAlwaysReturnsAPiece(t *testing.T) {
	pdf := newWrapTestPDF(t)
	if got := breakLongStyledRun(pdf, pdfStyledRun{Text: "", SizePt: 12}, 20, nil); len(got) != 1 {
		t.Errorf("empty run gave %d pieces, want 1", len(got))
	}
	// A box too narrow for even one glyph must still terminate.
	if got := breakLongStyledRun(pdf, pdfStyledRun{Text: "ab", SizePt: 12}, 0.01, nil); len(got) == 0 {
		t.Error("a box narrower than one glyph returned nothing")
	}
}

// Normal wrapping must be untouched.
func TestWrapStyledRunsStillBreaksAtSpaces(t *testing.T) {
	pdf := newWrapTestPDF(t)
	run := pdfStyledRun{Text: "one two three four five", SizePt: 12}

	lines := wrapStyledRuns(pdf, []pdfStyledRun{run}, 60, nil)
	for _, line := range lines {
		var text strings.Builder
		for _, r := range line {
			text.WriteString(r.Text)
		}
		if strings.Contains(strings.TrimSpace(text.String()), "  ") {
			t.Errorf("line %q has collapsed spacing", text.String())
		}
	}
}
