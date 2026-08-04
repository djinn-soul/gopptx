package pptxxml

import (
	"strings"
	"testing"
)

// A blank must not be written as a point at all. Writing 0 instead is what
// upstream issue #968 reports: a None reading is drawn as a real zero.
func TestRenderChartOmitsBlankValuesAndDeclaresDisplayBlanksAs(t *testing.T) {
	spec := &ChartSpec{
		Kind:            ChartKindLine,
		Title:           "Blanks",
		Categories:      []string{"Jan", "Feb", "Mar"},
		Values:          []float64{1, BlankValue(), 3},
		SeriesName:      "Series 1",
		Color:           "4472C4",
		Grouping:        "standard",
		CX:              100,
		CY:              100,
		DisplayBlanksAs: DisplayBlanksAsGap,
	}

	xml := string(RenderChart(spec))
	values := valueCacheSection(t, xml)

	if !strings.Contains(values, `<c:ptCount val="3"/>`) {
		t.Fatalf("expected the point count to keep counting every category:\n%s", values)
	}
	if strings.Contains(values, `<c:pt idx="1">`) {
		t.Fatalf("expected no value point for the blank category:\n%s", values)
	}
	for _, want := range []string{`<c:pt idx="0">`, `<c:pt idx="2">`} {
		if !strings.Contains(values, want) {
			t.Fatalf("expected %s to survive:\n%s", want, values)
		}
	}
	if !strings.Contains(xml, `<c:plotVisOnly val="1"/>`+"\n"+`<c:dispBlanksAs val="gap"/>`) {
		t.Fatalf("expected dispBlanksAs directly after plotVisOnly:\n%s", xml)
	}
}

// valueCacheSection returns the <c:val> block, so an assertion about value
// points is not answered by the category points, which share the c:pt element.
func valueCacheSection(t *testing.T, chartXML string) string {
	t.Helper()
	start := strings.Index(chartXML, "<c:val>")
	end := strings.Index(chartXML, "</c:val>")
	if start < 0 || end < start {
		t.Fatalf("chart has no value cache:\n%s", chartXML)
	}
	return chartXML[start:end]
}

func TestRenderChartOmitsDisplayBlanksAsWhenUnset(t *testing.T) {
	spec := &ChartSpec{
		Kind:       ChartKindBar,
		Title:      "Plain",
		Categories: []string{"A"},
		Values:     []float64{1},
		SeriesName: "Series 1",
		Color:      "4472C4",
		BarDir:     "col",
		Grouping:   "clustered",
		CX:         100,
		CY:         100,
	}

	if strings.Contains(string(RenderChart(spec)), "dispBlanksAs") {
		t.Fatal("expected no dispBlanksAs element when the caller named none")
	}
}

func TestNormalizedDisplayBlanksAsFallsBackToGap(t *testing.T) {
	if got := normalizedDisplayBlanksAs("nonsense"); got != DisplayBlanksAsGap {
		t.Fatalf("expected gap for an unknown token, got %q", got)
	}
	if got := normalizedDisplayBlanksAs(" SPAN "); got != DisplayBlanksAsSpan {
		t.Fatalf("expected span, got %q", got)
	}
}

// PowerPoint's "shrink text on overflow" is a normAutofit carrying the amounts
// it shrank by (upstream issue #969).
func TestTextBodyPrXMLWritesShrinkAmounts(t *testing.T) {
	fontScale := 62.5
	reduction := 10.0
	xml := TextBodyPrXML(&TextFrameSpec{
		Wrap:               "square",
		Anchor:             "ctr",
		AutoFit:            normAutoFitToken,
		FontScale:          &fontScale,
		LineSpaceReduction: &reduction,
	})

	if !strings.Contains(xml, `<a:normAutofit fontScale="62500" lnSpcReduction="10000"/>`) {
		t.Fatalf("expected the shrink amounts in thousandths of a percent:\n%s", xml)
	}
}

// A font scale below the schema floor would make PowerPoint reject the part, so
// it is clamped rather than written through.
func TestTextBodyPrXMLClampsShrinkAmounts(t *testing.T) {
	fontScale := 1.0
	reduction := 90.0
	xml := TextBodyPrXML(&TextFrameSpec{
		Wrap:               "square",
		Anchor:             "ctr",
		AutoFit:            normAutoFitToken,
		FontScale:          &fontScale,
		LineSpaceReduction: &reduction,
	})

	if !strings.Contains(xml, `fontScale="25000"`) || !strings.Contains(xml, `lnSpcReduction="20000"`) {
		t.Fatalf("expected clamped shrink amounts:\n%s", xml)
	}
}

func TestTextBodyPrXMLKeepsPlainNormAutofitWithoutAmounts(t *testing.T) {
	xml := TextBodyPrXML(&TextFrameSpec{
		Wrap:    "square",
		Anchor:  "ctr",
		AutoFit: normAutoFitToken,
	})

	if !strings.Contains(xml, `<a:normAutofit/>`) {
		t.Fatalf("expected a bare normAutofit when no amounts were given:\n%s", xml)
	}
}
