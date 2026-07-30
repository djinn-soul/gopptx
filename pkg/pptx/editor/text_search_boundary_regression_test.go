package editor

import (
	"strings"
	"testing"
)

func TestReplaceTextDoesNotCrossDrawingMLBreak(t *testing.T) {
	paragraph := []byte(
		`<a:p><a:r><a:t>foo</a:t></a:r><a:br/>` +
			`<a:r><a:t>bar</a:t></a:r></a:p>`,
	)

	got, count := replaceTextRuns(paragraph, "foobar", "joined")

	if count != 0 {
		t.Fatalf("replacement crossed a:br: count=%d xml=%s", count, got)
	}
	if string(got) != string(paragraph) {
		t.Fatalf("paragraph changed without a match:\n%s", got)
	}
}

func TestReplaceTextStillCrossesAdjacentRuns(t *testing.T) {
	paragraph := []byte(
		`<a:p><a:r><a:rPr b="1"/><a:t>foo</a:t></a:r>` +
			`<a:r><a:rPr i="1"/><a:t>bar</a:t></a:r></a:p>`,
	)

	got, count := replaceTextRuns(paragraph, "foobar", "joined")

	if count != 1 || !strings.Contains(string(got), "<a:t>joined</a:t>") {
		t.Fatalf("adjacent run replacement failed: count=%d xml=%s", count, got)
	}
}
