package export

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func TestParagraphStartGapCollapsesBeforeAfter(t *testing.T) {
	t.Parallel()

	style := text.ParagraphStyle{SpaceBeforePt: 6, SpaceAfterPt: 10}
	if got := paragraphStartGap(0, 0, style); got != 6 {
		t.Fatalf("first paragraph gap=%v want 6", got)
	}
	if got := paragraphAfterGap(style); got != 10 {
		t.Fatalf("after gap=%v want 10", got)
	}
	if got := paragraphStartGap(1, paragraphAfterGap(style), text.ParagraphStyle{SpaceBeforePt: 4}); got != 10 {
		t.Fatalf("collapsed inter-paragraph gap=%v want 10", got)
	}
}

func TestParagraphLineSpacingFactorHasFloor(t *testing.T) {
	t.Parallel()

	if got := paragraphLineSpacingFactor(text.ParagraphStyle{LineSpacingPct: 0}); got != 1.0 {
		t.Fatalf("default line spacing=%v want 1.0", got)
	}
	if got := paragraphLineSpacingFactor(text.ParagraphStyle{LineSpacingPct: 40}); got != 0.6 {
		t.Fatalf("line spacing floor=%v want 0.6", got)
	}
}

func TestBulletPrefixFormats(t *testing.T) {
	t.Parallel()

	if got := bulletPrefix(text.ParagraphStyle{BulletStyle: text.BulletStyleNumber}, 2); got != "3." {
		t.Fatalf("numbered bullet=%q want 3.", got)
	}
	if got := bulletPrefix(text.ParagraphStyle{BulletStyle: text.BulletStyleRomanUpper}, 3); got != "IV" {
		t.Fatalf("roman bullet=%q want IV", got)
	}
	if got := bulletPrefix(text.ParagraphStyle{BulletStyle: text.BulletStyleCustom, BulletChar: ">"}, 0); got != ">" {
		t.Fatalf("custom bullet=%q want >", got)
	}
}
