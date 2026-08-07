package text_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func TestParseInlineRunsReadsTheMarkers(t *testing.T) {
	runs, styled := text.ParseInlineRuns("plain **bold** and `code` and *italic*")
	if !styled {
		t.Fatal("formatting was not reported")
	}

	var bold, code, italic bool
	for _, run := range runs {
		switch {
		case run.Bold && run.Text == "bold":
			bold = true
		case run.Code && run.Text == "code":
			code = true
		case run.Italic && run.Text == "italic":
			italic = true
		}
	}
	if !bold || !code || !italic {
		t.Fatalf("bold=%v code=%v italic=%v in %+v", bold, code, italic, runs)
	}
}

func TestPlainTextIsNotReportedAsStyled(t *testing.T) {
	runs, styled := text.ParseInlineRuns("nothing to see here")
	if styled {
		t.Fatal("plain text should not report formatting")
	}
	if len(runs) != 1 || runs[0].Text != "nothing to see here" {
		t.Fatalf("runs = %+v", runs)
	}
}

func TestUnmatchedMarkersStayLiteral(t *testing.T) {
	runs, styled := text.ParseInlineRuns("2 * 3 = 6")
	if styled {
		t.Fatal("a lone asterisk is multiplication, not italics")
	}
	var joined strings.Builder
	for _, run := range runs {
		joined.WriteString(run.Text)
	}
	if joined.String() != "2 * 3 = 6" {
		t.Fatalf("text round trip = %q", joined.String())
	}
}

func TestHasInlineFormatting(t *testing.T) {
	if !text.HasInlineFormatting("a **b** c") {
		t.Error("bold not detected")
	}
	if text.HasInlineFormatting("a b c") {
		t.Error("plain text reported as formatted")
	}
}
