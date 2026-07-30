package elements_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// Two builders derived from one SlideContent must not share backing arrays. The
// source below has len 3, cap 4, so an unclone slice would make both appends
// land in the same slot and the first one would be lost.
func TestBuildFromDoesNotShareSliceStorage(t *testing.T) {
	base := elements.NewSlide("t").AddBullet("a").AddBullet("b").AddBullet("c")

	b1 := elements.BuildFrom(base)
	b2 := elements.BuildFrom(base)
	b1.AddBullet("FROM-B1")
	b2.AddBullet("FROM-B2")

	first := b1.Build().Bullets
	second := b2.Build().Bullets
	if len(first) != 4 || first[3] != "FROM-B1" {
		t.Errorf("builder 1 bullets = %v, want [a b c FROM-B1]", first)
	}
	if len(second) != 4 || second[3] != "FROM-B2" {
		t.Errorf("builder 2 bullets = %v, want [a b c FROM-B2]", second)
	}
	if got := base.Bullets; len(got) != 3 {
		t.Errorf("source slide was mutated: %v", got)
	}
}

func TestBuildFromClonesNestedBulletRuns(t *testing.T) {
	base := elements.NewSlide("t").AddBullet("a")
	base.BulletRuns[0] = append(base.BulletRuns[0], elements.Run{Text: "original"})

	derived := elements.BuildFrom(base).Build()
	derived.BulletRuns[0][0].Text = "changed"

	if base.BulletRuns[0][0].Text != "original" {
		t.Errorf("source run text = %q, want %q",
			base.BulletRuns[0][0].Text, "original")
	}
}
