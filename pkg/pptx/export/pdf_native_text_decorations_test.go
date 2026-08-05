package export

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func TestPDFStyledRunCarriesCharacterDecorations(t *testing.T) {
	run := text.NewRun("Hello").
		WithUnderlineStyle(text.UnderlineStyleDouble).
		WithStrikethroughStyle(text.StrikethroughStyleSingle)
	run.Superscript = true

	styled := pdfStyledRunFromTextRun(run, 18, false, false)
	if styled.Underline != text.UnderlineStyleDouble {
		t.Fatalf("underline=%q want %q", styled.Underline, text.UnderlineStyleDouble)
	}
	if styled.Strikethrough != text.StrikethroughStyleSingle {
		t.Fatalf("strikethrough=%q want %q", styled.Strikethrough, text.StrikethroughStyleSingle)
	}
	if !styled.Superscript {
		t.Fatal("superscript flag lost")
	}
}

func TestPDFStyledRunDropsDisabledDecorations(t *testing.T) {
	run := text.NewRun("Hello").WithUnderline(false).WithStrikethrough(false)
	styled := pdfStyledRunFromTextRun(run, 18, false, false)
	if styled.Underline != "" || styled.Strikethrough != "" {
		t.Fatalf("disabled decorations survived: u=%q strike=%q", styled.Underline, styled.Strikethrough)
	}
}

func TestPDFStyledRunAppliesCaps(t *testing.T) {
	allCaps := text.NewRun("Hello")
	allCaps.AllCaps = true
	styled := pdfStyledRunFromTextRun(allCaps, 18, false, false)
	if styled.Text != "HELLO" {
		t.Fatalf("all-caps text=%q want %q", styled.Text, "HELLO")
	}
	// All caps are drawn at full size, so the small-caps size reduction must not
	// also apply.
	if styled.SmallCaps {
		t.Fatal("all caps run marked small caps")
	}

	smallCaps := text.NewRun("Hello")
	smallCaps.SmallCaps = true
	styled = pdfStyledRunFromTextRun(smallCaps, 18, false, false)
	// The case change belongs to the split, which needs the original case to
	// tell capitals from small capitals.
	if styled.Text != "Hello" || !styled.SmallCaps {
		t.Fatalf("small-caps run=%+v want the original text and the flag", styled)
	}
}

func TestSplitSmallCapsRunKeepsCapitalsFullSize(t *testing.T) {
	run := pdfStyledRun{Text: "McDonald", SizePt: 20, SmallCaps: true}
	parts := splitSmallCapsRun(run)

	var joined strings.Builder
	for _, part := range parts {
		joined.WriteString(part.Text)
	}
	if joined.String() != "MCDONALD" {
		t.Fatalf("joined text=%q want %q", joined.String(), "MCDONALD")
	}

	// "M" and "D" were capitals and keep the run's size; the rest shrink.
	want := []struct {
		text  string
		small bool
	}{
		{"M", false},
		{"C", true},
		{"D", false},
		{"ONALD", true},
	}
	if len(parts) != len(want) {
		t.Fatalf("parts=%+v want %d segments", parts, len(want))
	}
	for i, w := range want {
		if parts[i].Text != w.text || parts[i].SmallCaps != w.small {
			t.Fatalf("part %d=%q small=%v want %q small=%v", i, parts[i].Text, parts[i].SmallCaps, w.text, w.small)
		}
		size := effectiveRunSizePt(nil, parts[i])
		if w.small && size >= 20 {
			t.Fatalf("small capital %q size=%d want under 20", parts[i].Text, size)
		}
		if !w.small && size != 20 {
			t.Fatalf("capital %q size=%d want 20", parts[i].Text, size)
		}
	}
}

func TestSplitSmallCapsRunLeavesOtherRunsAlone(t *testing.T) {
	run := pdfStyledRun{Text: "Hello", SizePt: 20}
	parts := splitSmallCapsRun(run)
	if len(parts) != 1 || parts[0].Text != "Hello" {
		t.Fatalf("parts=%+v want the run unchanged", parts)
	}
}

func TestDecorationWeightAndWavy(t *testing.T) {
	if decorationWeight("sng") != 1 {
		t.Fatal("plain underline should not be weighted")
	}
	if decorationWeight("sngHeavy") <= 1 {
		t.Fatal("heavy underline should be drawn thicker")
	}
	if !isWavyDecoration("wavyDouble") || isWavyDecoration("sng") {
		t.Fatal("wavy detection wrong")
	}
	// A wavy style is a stroke shape, not a line type.
	if got := decorationLineType("wavyHeavy"); got != "" {
		t.Fatalf("wavy line type=%q want solid", got)
	}
}

func TestEffectiveRunSizePtShrinksScriptAndSmallCaps(t *testing.T) {
	base := pdfStyledRun{Text: "x", SizePt: 20}
	if got := effectiveRunSizePt(nil, base); got != 20 {
		t.Fatalf("plain run size=%d want 20", got)
	}

	super := base
	super.Superscript = true
	if got := effectiveRunSizePt(nil, super); got >= 20 || got <= 0 {
		t.Fatalf("superscript size=%d want smaller than 20 and positive", got)
	}

	small := base
	small.SmallCaps = true
	if got := effectiveRunSizePt(nil, small); got != 16 {
		t.Fatalf("small-caps size=%d want 16", got)
	}
}

func TestScriptBaselineOffsetDirection(t *testing.T) {
	base := pdfStyledRun{Text: "x", SizePt: 20}
	if got := scriptBaselineOffsetPt(nil, base); got != 0 {
		t.Fatalf("plain run offset=%v want 0", got)
	}

	super := base
	super.Superscript = true
	if got := scriptBaselineOffsetPt(nil, super); got >= 0 {
		t.Fatalf("superscript offset=%v want negative (raised)", got)
	}

	sub := base
	sub.Subscript = true
	if got := scriptBaselineOffsetPt(nil, sub); got <= 0 {
		t.Fatalf("subscript offset=%v want positive (dropped)", got)
	}
}

func TestNormalizeDecoration(t *testing.T) {
	cases := map[string]string{
		"":                        "",
		"  ":                      "",
		text.UnderlineStyleNone:   "",
		text.UnderlineStyleSingle: text.UnderlineStyleSingle,
		" dotted ":                "dotted",
	}
	for in, want := range cases {
		if got := normalizeDecoration(in); got != want {
			t.Fatalf("normalizeDecoration(%q)=%q want %q", in, got, want)
		}
	}
}

func TestDecorationLineType(t *testing.T) {
	if got := decorationLineType(text.UnderlineStyleDotted); got != "dotted" {
		t.Fatalf("dotted line type=%q want %q", got, "dotted")
	}
	if got := decorationLineType("dashHeavy"); got != "dashed" {
		t.Fatalf("dash line type=%q want %q", got, "dashed")
	}
	if got := decorationLineType(text.UnderlineStyleSingle); got != "" {
		t.Fatalf("single line type=%q want solid", got)
	}
}

func TestIsDoubleDecoration(t *testing.T) {
	if !isDoubleDecoration(text.UnderlineStyleDouble) || !isDoubleDecoration(text.StrikethroughStyleDouble) {
		t.Fatal("double styles not detected")
	}
	if isDoubleDecoration(text.UnderlineStyleSingle) {
		t.Fatal("single underline reported as double")
	}
}

func TestLineMetricsDecorationDefaults(t *testing.T) {
	// A font that states none of the optional fields still yields usable
	// decoration geometry.
	empty := ttfLineMetrics{UnitsPerEm: 1000}
	if empty.underlineOffsetFactor() <= 0 {
		t.Fatalf("underline offset=%v want positive (below baseline)", empty.underlineOffsetFactor())
	}
	if empty.underlineThicknessFactor() <= 0 {
		t.Fatalf("underline thickness=%v want positive", empty.underlineThicknessFactor())
	}
	if empty.strikeoutOffsetFactor() <= 0 {
		t.Fatalf("strikeout offset=%v want positive (above baseline)", empty.strikeoutOffsetFactor())
	}
	if got := empty.subscriptSizeFactor(); got <= 0 || got >= 1 {
		t.Fatalf("subscript size factor=%v want between 0 and 1", got)
	}
}

func TestLineMetricsUsesStatedDecorationValues(t *testing.T) {
	m := ttfLineMetrics{UnitsPerEm: 1000, UnderlinePosition: -150, UnderlineThickness: 60}
	if got := m.underlineOffsetFactor(); got != 0.15 {
		t.Fatalf("underline offset factor=%v want 0.15", got)
	}
	if got := m.underlineThicknessFactor(); got != 0.06 {
		t.Fatalf("underline thickness factor=%v want 0.06", got)
	}
}
