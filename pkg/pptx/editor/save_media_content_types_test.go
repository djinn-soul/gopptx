package editor

import (
	"archive/zip"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// A media part the editor did not add itself -- one a layout references, or one
// another tool dropped in -- still needs its extension declared. Without it
// PowerPoint refuses the package with "no content-type for partname
// '/ppt/media/...'" (upstream issue #915).
func TestSaveDeclaresContentTypeForUntrackedMedia(t *testing.T) {
	base := writeDeckFixture(t, "untracked-media.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	})

	ed, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	ed.parts.Set("ppt/media/image-1002-2.svg", []byte(`<svg xmlns="http://www.w3.org/2000/svg"/>`))

	out := filepath.Join(t.TempDir(), "with-svg.pptx")
	if err = ed.Save(out); err != nil {
		t.Fatalf("save: %v", err)
	}

	contentTypes := readPackagePart(t, out, "[Content_Types].xml")
	if !strings.Contains(contentTypes, `Extension="svg"`) {
		t.Fatalf("expected an svg Default in the content types:\n%s", contentTypes)
	}
	if !strings.Contains(contentTypes, "image/svg+xml") {
		t.Fatalf("expected the svg content type:\n%s", contentTypes)
	}
}

func readPackagePart(t *testing.T, pptxPath, partName string) string {
	t.Helper()
	reader, err := zip.OpenReader(pptxPath)
	if err != nil {
		t.Fatalf("open package: %v", err)
	}
	defer func() { _ = reader.Close() }()

	for _, file := range reader.File {
		if file.Name != partName {
			continue
		}
		rc, openErr := file.Open()
		if openErr != nil {
			t.Fatalf("open part %s: %v", partName, openErr)
		}
		defer func() { _ = rc.Close() }()
		data, readErr := io.ReadAll(rc)
		if readErr != nil {
			t.Fatalf("read part %s: %v", partName, readErr)
		}
		return string(data)
	}
	t.Fatalf("part %s not found in %s", partName, pptxPath)
	return ""
}
