package presentation

import (
	"strings"
	"testing"

	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

func TestRewriteSlideHiddenAttribute(t *testing.T) {
	const src = `<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld/></p:sld>`
	const hiddenPrefix = `<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" show="0">`
	hiddenBytes, err := editorslide.RewriteSlideHidden([]byte(src), true)
	if err != nil {
		t.Fatalf("RewriteSlideHidden(hidden=true) failed: %v", err)
	}
	hidden := string(hiddenBytes)
	if !strings.Contains(hidden, hiddenPrefix) {
		t.Fatalf("expected show=0 on p:sld root, got %s", hidden)
	}

	visibleBytes, err := editorslide.RewriteSlideHidden(hiddenBytes, false)
	if err != nil {
		t.Fatalf("RewriteSlideHidden(hidden=false) failed: %v", err)
	}
	if strings.Contains(string(visibleBytes), `show="0"`) {
		t.Fatalf("expected show attribute removed, got %s", string(visibleBytes))
	}
}
