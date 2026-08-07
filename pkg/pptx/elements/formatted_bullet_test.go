package elements_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// AddBullet treats ** and backticks as literal text; the parser that reads them
// was reachable only through the Markdown importer.
func TestAddFormattedBulletParsesMarkers(t *testing.T) {
	slide := elements.NewSlide("Formatting").AddFormattedBullet("**bold** and `code`")

	if len(slide.BulletRuns) != 1 || len(slide.BulletRuns[0]) == 0 {
		t.Fatalf("expected rich runs, got %+v", slide.BulletRuns)
	}
	var bold, code bool
	for _, run := range slide.BulletRuns[0] {
		bold = bold || (run.Bold && run.Text == "bold")
		code = code || (run.Code && run.Text == "code")
	}
	if !bold || !code {
		t.Fatalf("bold=%v code=%v in %+v", bold, code, slide.BulletRuns[0])
	}
	if slide.Bullets[0] != "bold and code" {
		t.Fatalf("plain text = %q, want the markers stripped", slide.Bullets[0])
	}
}

func TestAddFormattedBulletLeavesPlainTextAlone(t *testing.T) {
	slide := elements.NewSlide("Formatting").AddFormattedBullet("just text")
	if slide.Bullets[0] != "just text" {
		t.Fatalf("bullet = %q", slide.Bullets[0])
	}
	if slide.BulletRuns[0] != nil {
		t.Fatalf("plain text should not build runs: %+v", slide.BulletRuns[0])
	}
}
