package export

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func TestParagraphStartGapCollapsesBeforeAfter(t *testing.T) {
	t.Parallel()

	style := text.ParagraphStyle{SpaceBeforePt: 6, SpaceAfterPt: 10}
	// PowerPoint does not lead the first paragraph in a box with space-before.
	if got := paragraphStartGap(0, 0, style, shapeTextSpacing()); got != 0 {
		t.Fatalf("first paragraph gap=%v want 0", got)
	}
	if got := paragraphAfterGap(style); got != 10 {
		t.Fatalf("after gap=%v want 10", got)
	}
	got := paragraphStartGap(1, paragraphAfterGap(style), text.ParagraphStyle{SpaceBeforePt: 4}, shapeTextSpacing())
	if got != 10 {
		t.Fatalf("collapsed inter-paragraph gap=%v want 10", got)
	}
}

func TestParagraphStartGapAppliesDefaultOnlyWhenUnset(t *testing.T) {
	t.Parallel()

	// A body placeholder inherits the Office bodyStyle default of 10pt
	// space-before when the paragraph states none.
	if got := paragraphStartGap(1, 0, text.ParagraphStyle{}, bodyPlaceholderSpacing()); got != 10 {
		t.Fatalf("default body space-before=%v want 10", got)
	}
	// An explicit value still wins over the default.
	if got := paragraphStartGap(1, 0, text.ParagraphStyle{SpaceBeforePt: 3}, bodyPlaceholderSpacing()); got != 3 {
		t.Fatalf("explicit space-before=%v want 3", got)
	}
	// A plain text box inherits otherStyle, which has no space-before.
	if got := paragraphStartGap(1, 0, text.ParagraphStyle{}, shapeTextSpacing()); got != 0 {
		t.Fatalf("default shape space-before=%v want 0", got)
	}
}

func TestParagraphLineSpacingFactorHasFloor(t *testing.T) {
	t.Parallel()

	// Body placeholders inherit the Office bodyStyle default of 90%.
	if got := paragraphLineSpacingFactor(text.ParagraphStyle{}, bodyPlaceholderSpacing()); got != 0.9 {
		t.Fatalf("default body line spacing=%v want 0.9", got)
	}
	if got := paragraphLineSpacingFactor(text.ParagraphStyle{}, shapeTextSpacing()); got != 1.0 {
		t.Fatalf("default shape line spacing=%v want 1.0", got)
	}
	// An explicit percentage overrides the default, with a floor.
	explicit := paragraphLineSpacingFactor(text.ParagraphStyle{LineSpacingPct: 150}, bodyPlaceholderSpacing())
	if explicit != 1.5 {
		t.Fatalf("explicit line spacing=%v want 1.5", explicit)
	}
	if got := paragraphLineSpacingFactor(text.ParagraphStyle{LineSpacingPct: 40}, shapeTextSpacing()); got != 0.6 {
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
