package presentation

import (
	"strings"
	"testing"
)

func TestRewriteSlideHiddenAttribute(t *testing.T) {
	const src = `<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld/></p:sld>`
	const hiddenPrefix = `<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" show="0">`
	hidden, err := rewriteSlideHiddenAttribute(src, true)
	if err != nil {
		t.Fatalf("rewriteSlideHiddenAttribute(hidden=true) failed: %v", err)
	}
	if !strings.Contains(hidden, hiddenPrefix) {
		t.Fatalf("expected show=0 on p:sld root, got %s", hidden)
	}

	visible, err := rewriteSlideHiddenAttribute(hidden, false)
	if err != nil {
		t.Fatalf("rewriteSlideHiddenAttribute(hidden=false) failed: %v", err)
	}
	if strings.Contains(visible, `show="0"`) {
		t.Fatalf("expected show attribute removed, got %s", visible)
	}
}
