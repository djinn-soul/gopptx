package media

import (
	"net/http"
	"os"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

func TestBuildMediaCatalog_Coverage(t *testing.T) {
	// 1. Prepare assets
	imgFile := "test_bg.png"
	_ = os.WriteFile(imgFile, []byte("fake-png"), 0600)
	defer os.Remove(imgFile)

	soundFile := "test_sound.wav"
	_ = os.WriteFile(soundFile, []byte("fake-wav"), 0600)
	defer os.Remove(soundFile)

	// 2. Mock slides
	slides := []elements.SlideContent{
		{
			Images: []shapes.Image{{Path: imgFile}},
			Background: &elements.SlideBackground{
				Type:        elements.SlideBackgroundPicture,
				PictureFill: &shapes.Image{Data: []byte("bg-bytes"), Format: "jpg"},
			},
			Transition: transitions.TransitionOptions{
				Type:  transitions.TransitionCut,
				Sound: &transitions.TransitionSound{RelID: "file:" + soundFile},
			},
			PlaceholderOverrides: []shapes.PlaceholderContent{
				{Image: &shapes.Image{Data: []byte("ph-bytes"), Format: "png"}},
			},
		},
	}

	// 3. Mock NotesMaster
	master := &elements.NotesMaster{
		Background: &elements.SlideBackground{
			Type:        elements.SlideBackgroundPicture,
			PictureFill: &shapes.Image{Data: []byte("master-bg"), Format: "png"},
		},
	}

	// 4. Run Build
	catalog, err := BuildMediaCatalog(slides, master)
	if err != nil {
		t.Fatalf("BuildMediaCatalog failed: %v", err)
	}

	// 5. Verify counts
	// 1 image, 1 bg image, 1 sound, 1 ph image, 1 master bg image = 5 unique assets
	if len(catalog.Assets()) != 5 {
		t.Errorf("expected 5 assets, got %d", len(catalog.Assets()))
	}
}

func TestCatalog_ErrorPaths(t *testing.T) {
	client := &http.Client{}
	catalog := &Catalog{byKey: make(map[string]Asset)}

	// Unsupported extension
	img := shapes.Image{Data: []byte("fake"), Format: "txt"}
	err := addImageToCatalog(catalog, client, img)
	if err == nil {
		t.Error("expected error for unsupported extension")
	}

	// File not found
	img = shapes.Image{Path: "nonexistent.png"}
	err = addImageToCatalog(catalog, client, img)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}

	// Empty data
	err = appendMediaAsset(catalog, "key", "png", []byte{})
	if err == nil {
		t.Error("expected error for empty data")
	}
}

func TestCatalog_MediaName(t *testing.T) {
	catalog := &Catalog{byKey: make(map[string]Asset)}
	img := shapes.Image{Data: []byte("test"), Format: "png"}
	name, _ := catalog.RegisterImage(img)

	name2, ok := catalog.MediaNameForImage(img)
	if !ok || name != name2 {
		t.Errorf("expected %s, got %s", name, name2)
	}

	// Missing image
	_, ok = catalog.MediaNameForImage(shapes.Image{Path: "missing.png"})
	if ok {
		t.Error("expected false for missing image")
	}
}
