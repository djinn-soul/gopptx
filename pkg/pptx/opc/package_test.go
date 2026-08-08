package opc_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/opc"
)

func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o600)
}

// deckOnDisk writes a small generated deck and returns its path.
func deckOnDisk(t *testing.T) string {
	t.Helper()
	data, err := pptx.CreateWithSlides("OPC", []pptx.SlideContent{
		pptx.NewSlide("First").AddBullet("Hello"),
	})
	if err != nil {
		t.Fatalf("create deck: %v", err)
	}
	path := filepath.Join(t.TempDir(), "deck.pptx")
	if err = writeFile(path, data); err != nil {
		t.Fatalf("write deck: %v", err)
	}
	return path
}

// The only part-level API was internal, write-only, and had no read, remove or
// list, so a caller needing a part the typed API does not cover was stuck.
func TestPackageReadsAndListsParts(t *testing.T) {
	pkg, err := opc.Open(deckOnDisk(t))
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	if !pkg.Has("[Content_Types].xml") || !pkg.Has("ppt/presentation.xml") {
		t.Fatalf("expected the standard parts, got %v", pkg.PartPaths())
	}
	if pkg.PartCount() < 10 {
		t.Fatalf("part count = %d, want a real package", pkg.PartCount())
	}

	slides := pkg.PartPathsWithPrefix("ppt/slides/")
	if len(slides) == 0 {
		t.Fatal("no slide parts listed")
	}

	xmlText, err := pkg.PartString("ppt/presentation.xml")
	if err != nil {
		t.Fatalf("read presentation: %v", err)
	}
	if !strings.Contains(xmlText, "<p:presentation") {
		t.Fatalf("presentation.xml does not look like one: %s", xmlText)
	}

	if _, err = pkg.Part("ppt/nope.xml"); !errors.Is(err, opc.ErrPartNotFound) {
		t.Fatalf("error for a missing part = %v, want ErrPartNotFound", err)
	}
}

func TestPackageEditsSurviveASaveAndReopen(t *testing.T) {
	pkg, err := opc.Open(deckOnDisk(t))
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	pkg.SetPartString("customXml/item99.xml", "<root>escape hatch</root>")
	if !pkg.RemovePart("docProps/app.xml") {
		t.Fatal("expected app.xml to be present before removing it")
	}

	out := filepath.Join(t.TempDir(), "edited.pptx")
	if err = pkg.Save(out); err != nil {
		t.Fatalf("save: %v", err)
	}

	reopened, err := opc.Open(out)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	content, err := reopened.PartString("customXml/item99.xml")
	if err != nil || !strings.Contains(content, "escape hatch") {
		t.Fatalf("added part did not survive: %q, %v", content, err)
	}
	if reopened.Has("docProps/app.xml") {
		t.Fatal("removed part came back")
	}
}

// A part reaches another part only through a relationship, so the escape hatch
// has to cover them too.
func TestRelationshipsAreReadableAndEditable(t *testing.T) {
	pkg, err := opc.Open(deckOnDisk(t))
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	rels, err := pkg.Relationships("ppt/presentation.xml")
	if err != nil {
		t.Fatalf("read relationships: %v", err)
	}
	slideRels := rels.ByType("http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide")
	if len(slideRels) == 0 {
		t.Fatalf("no slide relationships found in %v", rels.Types())
	}

	added := rels.Add("http://example.com/custom", "custom/part.xml")
	external := rels.AddExternal(
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink",
		"https://example.com",
	)
	if added.ID == external.ID {
		t.Fatalf("ids collided: %s", added.ID)
	}
	pkg.SetRelationships("ppt/presentation.xml", rels)

	reread, err := pkg.Relationships("ppt/presentation.xml")
	if err != nil {
		t.Fatalf("re-read relationships: %v", err)
	}
	roundTripped, ok := reread.ByID(external.ID)
	if !ok || !roundTripped.External {
		t.Fatalf("external relationship lost: %+v", roundTripped)
	}
	if !reread.Remove(added.ID) {
		t.Fatal("expected to remove the relationship just added")
	}
}

func TestRelationshipsPathFollowsTheOPCLayout(t *testing.T) {
	cases := map[string]string{
		"":                      "_rels/.rels",
		"ppt/presentation.xml":  "ppt/_rels/presentation.xml.rels",
		"ppt/slides/slide1.xml": "ppt/slides/_rels/slide1.xml.rels",
	}
	for part, want := range cases {
		if got := opc.RelationshipsPath(part); got != want {
			t.Errorf("RelationshipsPath(%q) = %q, want %q", part, got, want)
		}
	}
}

// A package built from nothing is a legitimate use of the API.
func TestNewPackageRoundTrips(t *testing.T) {
	pkg := opc.New()
	pkg.SetPartString("[Content_Types].xml", `<Types/>`)

	data, err := pkg.Bytes()
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	reopened, err := opc.OpenBytes(data)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	if reopened.PartCount() != 1 {
		t.Fatalf("part count = %d, want 1", reopened.PartCount())
	}
}
