package export

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func TestParagraphRenderedLineHeightUsesPointSpacing(t *testing.T) {
	t.Parallel()

	if got := paragraphRenderedLineHeight(text.ParagraphStyle{LineSpacingPts: 18}, 12); got != 18 {
		t.Fatalf("point line spacing=%v want 18", got)
	}
	if got := paragraphRenderedLineHeight(text.ParagraphStyle{LineSpacingPct: 150}, 12); got != 18 {
		t.Fatalf("pct line spacing=%v want 18", got)
	}
}

func TestNextPDFTabAdvanceUsesStopsAndFallback(t *testing.T) {
	t.Parallel()

	if got := nextPDFTabAdvance(12, []float64{24, 48}); got != 12 {
		t.Fatalf("first tab stop advance=%v want 12", got)
	}
	if got := nextPDFTabAdvance(48, []float64{24, 48}); got != 24 {
		t.Fatalf("fallback tab advance=%v want 24", got)
	}
}

func TestBuildPDFStyledRunsPreservesHighlightAndOutline(t *testing.T) {
	t.Parallel()

	runs := buildPDFStyledRuns([]text.Run{
		text.NewRun("Hello").
			WithBold(true).
			WithItalic(true).
			WithColor("112233").
			WithHighlight("AABBCC").
			WithOutline("445566", 2),
	}, 20, false, false)
	if len(runs) != 1 {
		t.Fatalf("expected 1 styled run, got %d", len(runs))
	}
	got := runs[0]
	if !got.HasHighlight || got.HighlightColor != [3]uint8{0xAA, 0xBB, 0xCC} {
		t.Fatalf("expected highlight metadata, got %+v", got)
	}
	if !got.HasOutline || got.OutlineColor != [3]uint8{0x44, 0x55, 0x66} || got.OutlineWidthPt != 2 {
		t.Fatalf("expected outline metadata, got %+v", got)
	}
}

func TestShapeTextRotationAngleRecognizesVerticalModes(t *testing.T) {
	t.Parallel()

	frame := shapes.NewTextFrame().WithOrientation("vert270")
	angle, ok := shapeTextRotationAngle(&frame)
	if ok || angle != 0 {
		t.Fatalf("vert270 orientation should not force rotation angle, got %v ok=%v", angle, ok)
	}
	if !isVerticalShapeText(frame.Orientation) {
		t.Fatalf("expected vertical orientation detection for %q", frame.Orientation)
	}

	rotated := shapes.NewTextFrame().WithRotation(15)
	angle, ok = shapeTextRotationAngle(&rotated)
	if !ok || angle != 15 {
		t.Fatalf("rotation angle=%v ok=%v want 15,true", angle, ok)
	}
}

func TestParagraphTabStopsPtConvertsEMU(t *testing.T) {
	t.Parallel()

	style := text.NewParagraphStyle().WithTabStops(styling.Emu(914400), styling.Emu(1828800))
	got := paragraphTabStopsPt(style)
	if len(got) != 2 || got[0] != 72 || got[1] != 144 {
		t.Fatalf("tab stop conversion=%v want [72 144]", got)
	}
}

func TestResolvePDFFontAliasForRunUsesLangFallback(t *testing.T) {
	t.Parallel()

	setPDFFontAliases("SansAlias", "SerifAlias", "MonoAlias")
	setPDFCJKAlias("CJKAlias")
	if got := resolvePDFFontAliasForRun("", "ja-JP"); got != "CJKAlias" {
		t.Fatalf("expected CJK alias for ja-JP, got %q", got)
	}
	if got := resolvePDFFontAliasForRun("Consolas", "ja-JP"); got != "MonoAlias" {
		t.Fatalf("expected font hint to win over lang fallback, got %q", got)
	}
}
