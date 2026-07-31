package pptxxml

import (
	"strings"
	"testing"
)

// Font sizes are carried as points and written as centipoints, so half-point
// sizes must survive rendering rather than being truncated to whole points.
func TestRichTextRunWritesFractionalPointSize(t *testing.T) {
	tests := []struct {
		name   string
		sizePt float64
		wantSz string
	}{
		{name: "whole point", sizePt: 24, wantSz: `sz="2400"`},
		{name: "half point", sizePt: 11.5, wantSz: `sz="1150"`},
		{name: "rounds to nearest centipoint", sizePt: 10.005, wantSz: `sz="1001"`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			xml := richTextRun(
				TextRunSpec{Text: "sized", SizePt: tc.sizePt},
				ContentStyleSpec{SizePt: 14},
			)
			if !strings.Contains(xml, tc.wantSz) {
				t.Fatalf("expected %s in run XML, got %s", tc.wantSz, xml)
			}
		})
	}
}

func TestRichTextRunFallsBackToContentStyleSize(t *testing.T) {
	xml := richTextRun(TextRunSpec{Text: "unsized"}, ContentStyleSpec{SizePt: 14})
	if !strings.Contains(xml, `sz="1400"`) {
		t.Fatalf("expected content-style size in run XML, got %s", xml)
	}
}
