package editor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestImageDeduplication(t *testing.T) {
	tmpDir := t.TempDir()

	modifiedPath := filepath.Join(tmpDir, "modified.pptx")

	// 1. Create a base presentation
	slides := []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	}
	// Use writeDeckFixture from editor_test.go which is in the same package 'editor'
	basePath := writeDeckFixture(t, "base.pptx", slides)

	// 2. Open with Editor
	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("failed to open editor: %v", err)
	}
	defer func() {
		if closeErr := editor.Close(); closeErr != nil {
			t.Errorf("failed to close editor: %v", closeErr)
		}
	}()

	// 3. Register the same image twice
	imgData := []byte("fake-image-data-123")
	path1, err := editor.RegisterImage(imgData, "png")
	if err != nil {
		t.Fatalf("first registration failed: %v", err)
	}
	path2, err := editor.RegisterImage(imgData, "png")
	if err != nil {
		t.Fatalf("second registration failed: %v", err)
	}

	if path1 != path2 {
		t.Errorf("deduplication failed: expected same path, got %q and %q", path1, path2)
	}

	// 4. Save and verify parts in the archive
	if saveErr := editor.Save(modifiedPath); saveErr != nil {
		t.Fatalf("failed to save: %v", saveErr)
	}

	// We can use the editor again to open the saved file and check part list
	reopened, err := OpenPresentationEditor(modifiedPath)
	if err != nil {
		t.Fatalf("failed to reopen: %v", err)
	}
	defer func() {
		if closeErr := reopened.Close(); closeErr != nil {
			t.Errorf("failed to close reopened editor: %v", closeErr)
		}
	}()

	// Direct part store access or inventory check
	// Every RegisterImage should have populated the internal mediaInventory tracking

	// Check physical parts in reopened store
	foundImages := 0
	for _, k := range reopened.parts.Keys() {
		if strings.HasPrefix(k, "ppt/media/image") {
			foundImages++
		}
	}

	if foundImages != 1 {
		t.Errorf("expected 1 image part in archive, found %d", foundImages)
	}
}
