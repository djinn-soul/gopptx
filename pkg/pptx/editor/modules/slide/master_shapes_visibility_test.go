package slide

import (
	"strings"
	"testing"
)

func TestRewriteSlideShowMasterShapesRoundTrip(t *testing.T) {
	source := []byte(
		`<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"><p:cSld/></p:sld>`,
	)

	hidden, err := RewriteSlideShowMasterShapes(source, false)
	if err != nil {
		t.Fatalf("hide master shapes: %v", err)
	}
	if !strings.Contains(string(hidden), `showMasterSp="0"`) {
		t.Fatalf("expected exact showMasterSp=0 attribute: %s", hidden)
	}
	visible, err := ParseSlideShowMasterShapes(hidden)
	if err != nil {
		t.Fatalf("parse hidden master shapes: %v", err)
	}
	if visible {
		t.Fatal("showMasterSp=0 must parse as hidden")
	}

	shown, err := RewriteSlideShowMasterShapes(hidden, true)
	if err != nil {
		t.Fatalf("show master shapes: %v", err)
	}
	if strings.Contains(string(shown), "showMasterSp") {
		t.Fatalf("default visible state must remove showMasterSp: %s", shown)
	}
}
