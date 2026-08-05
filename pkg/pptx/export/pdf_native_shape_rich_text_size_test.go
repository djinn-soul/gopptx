package export

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func TestShapeParagraphNaturalSizeUsesTheRunSize(t *testing.T) {
	// The default size is a fallback for a run that states none, not a floor:
	// 10pt diagram labels used to be laid out at 18pt.
	small := []text.Paragraph{{Runs: []text.Run{text.NewRun("node").WithSizePt(10)}}}
	if got := shapeParagraphNaturalSize(small); got != 10 {
		t.Fatalf("natural size=%d want 10", got)
	}

	mixed := []text.Paragraph{{Runs: []text.Run{
		text.NewRun("a").WithSizePt(10),
		text.NewRun("b").WithSizePt(24),
	}}}
	if got := shapeParagraphNaturalSize(mixed); got != 24 {
		t.Fatalf("natural size=%d want the largest run (24)", got)
	}

	unsized := []text.Paragraph{{Runs: []text.Run{text.NewRun("plain")}}}
	if got := shapeParagraphNaturalSize(unsized); got != defaultFontSize {
		t.Fatalf("natural size=%d want the default %d", got, defaultFontSize)
	}
	if got := shapeParagraphNaturalSize(nil); got != defaultFontSize {
		t.Fatalf("natural size of no paragraphs=%d want the default", got)
	}
}

func TestShapeTextBoxFitsALine(t *testing.T) {
	// A 10pt label is a 12pt line. In a 15.3pt box (a diagram node) it fits, so
	// clipping stays on.
	line := []shapeParagraphLayoutLine{{lineHeight: pdfLineHeight(10)}}
	if !shapeTextBoxFitsALine(line, 15.3) {
		t.Fatal("a 12pt line should fit a 15.3pt box")
	}
	// Once the box is shorter than one line, clipping would erase the text
	// entirely, so it must be switched off.
	if shapeTextBoxFitsALine(line, 8) {
		t.Fatal("a 12pt line must not be reported as fitting an 8pt box")
	}
	if !shapeTextBoxFitsALine(nil, 1) {
		t.Fatal("no lines means nothing to clip")
	}
}

func TestMaxStyledRunsLineHeightFollowsTheRuns(t *testing.T) {
	// A 10pt line is 12pt tall. Seeding the maximum with the default size gave
	// it a 21.6pt floor, which overflowed short shapes and got clipped away.
	small := []pdfStyledRun{{Text: "node", SizePt: 10}}
	if got := maxStyledRunsLineHeight(small); got != pdfLineHeight(10) {
		t.Fatalf("line height=%v want %v", got, pdfLineHeight(10))
	}

	mixed := []pdfStyledRun{{Text: "a", SizePt: 10}, {Text: "b", SizePt: 24}}
	if got := maxStyledRunsLineHeight(mixed); got != pdfLineHeight(24) {
		t.Fatalf("line height=%v want the tallest run %v", got, pdfLineHeight(24))
	}

	if got := maxStyledRunsLineHeight(nil); got != pdfLineHeight(defaultFontSize) {
		t.Fatalf("line height with no runs=%v want the default", got)
	}
	// A run with no stated size still falls back to the default.
	unsized := []pdfStyledRun{{Text: "plain"}}
	if got := maxStyledRunsLineHeight(unsized); got != pdfLineHeight(defaultFontSize) {
		t.Fatalf("line height=%v want the default", got)
	}
}
