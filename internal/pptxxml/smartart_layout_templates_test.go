package pptxxml_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

// A handful of layouts are named with a space or a plus sign. A directory named
// that way is skipped by the embed patterns, and the layout then quietly falls
// back to the shared template — a different diagram, with no error to show for
// it. These are the ones that have to survive the round trip from URI to
// directory name.
func TestLayoutTemplatesResolveForAwkwardNames(t *testing.T) {
	for _, uri := range []string{
		"urn:microsoft.com/office/officeart/2011/layout/Picture Frame",
		"urn:microsoft.com/office/officeart/2005/8/layout/chevronAccent+Icon",
		"urn:microsoft.com/office/officeart/2005/8/layout/rings+Icon",
		"urn:microsoft.com/office/officeart/2005/8/layout/gear1",
		"urn:microsoft.com/office/officeart/2005/8/layout/process1",
	} {
		layout := pptxxml.SmartArtLayoutXML(uri, "")
		if !strings.Contains(layout, uri) {
			t.Errorf("layout %q fell back to the shared template", uri)
		}
	}
}
