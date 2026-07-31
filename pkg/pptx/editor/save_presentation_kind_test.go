package editor

import (
	"archive/zip"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestSavePPSXWritesSlideshowMainContentType(t *testing.T) {
	source := writeDeckFixture(t, "issue-438-source.pptx", []elements.SlideContent{
		elements.NewSlide("PowerPoint Show"),
	})
	presentation, err := OpenPresentationEditor(source)
	if err != nil {
		t.Fatalf("open presentation: %v", err)
	}
	defer func() { _ = presentation.Close() }()

	output := filepath.Join(t.TempDir(), "issue-438.ppsx")
	if err := presentation.Save(output); err != nil {
		t.Fatalf("save .ppsx: %v", err)
	}
	contentTypes := readPPSXPart(t, output, "[Content_Types].xml")
	if !strings.Contains(contentTypes, slideshowMainContentType) {
		t.Fatalf("slideshow content type missing: %s", contentTypes)
	}
	if strings.Contains(contentTypes, presentationMainContentType) {
		t.Fatalf("presentation content type remained: %s", contentTypes)
	}
}

func readPPSXPart(t *testing.T, path string, partName string) string {
	t.Helper()
	archive, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer func() { _ = archive.Close() }()
	for _, file := range archive.File {
		if file.Name != partName {
			continue
		}
		reader, openErr := file.Open()
		if openErr != nil {
			t.Fatalf("open %s: %v", partName, openErr)
		}
		data, readErr := io.ReadAll(reader)
		_ = reader.Close()
		if readErr != nil {
			t.Fatalf("read %s: %v", partName, readErr)
		}
		return string(data)
	}
	t.Fatalf("part not found: %s", partName)
	return ""
}
