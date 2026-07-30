package editor

import (
	"strings"
	"testing"
)

func TestSetAppPropertyIntPreservesNamespacePrefixAndAttributes(t *testing.T) {
	source := `<ep:Properties xmlns:ep="urn:extended">` +
		`<ep:Slides data-origin="powerpoint">1</ep:Slides>` +
		`</ep:Properties>`

	got := setAppPropertyInt(source, appSlidesPattern, "Slides", 3)

	if !strings.Contains(got, `<ep:Slides data-origin="powerpoint">3</ep:Slides>`) {
		t.Fatalf("prefixed element was not preserved: %s", got)
	}
	if strings.Contains(got, `<Slides>`) {
		t.Fatalf("unprefixed duplicate added: %s", got)
	}
}

func TestSetAppPropertyIntUsesRootPrefixForMissingElement(t *testing.T) {
	source := `<ep:Properties xmlns:ep="urn:extended"><ep:Application>PowerPoint</ep:Application></ep:Properties>`

	got := setAppPropertyInt(source, appNotesPattern, "Notes", 2)

	if !strings.Contains(got, `<ep:Notes>2</ep:Notes>`) {
		t.Fatalf("missing prefixed element not appended: %s", got)
	}
	if !strings.HasSuffix(strings.TrimSpace(got), `</ep:Properties>`) {
		t.Fatalf("element appended outside root: %s", got)
	}
}
