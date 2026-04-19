package media

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/netsec"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestMedia_Normalizers(t *testing.T) {
	if NormalizeImageExtension("test.PNG?query=1") != "png" {
		t.Error("NormalizeImageExtension failed")
	}
	if NormalizeExtension(".JPG") != "jpg" {
		t.Error("NormalizeExtension failed")
	}
}

func TestMedia_Loaders(t *testing.T) {
	t.Run("Bytes", func(t *testing.T) {
		key, ext, data := loadImageFromBytes([]byte("fake"), "png")
		if !strings.HasPrefix(key, "data:") || ext != "png" || len(data) != 4 {
			t.Error("loadImageFromBytes failed")
		}
	})

	t.Run("Path", func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "test.png")
		_ = os.WriteFile(path, []byte("fake"), 0600)

		key, ext, data, err := loadImageFromPath(path)
		if err != nil {
			t.Fatalf("loadImageFromPath failed: %v", err)
		}
		if !strings.HasPrefix(key, "path:") || ext != "png" || len(data) != 4 {
			t.Error("loadImageFromPath properties failed")
		}
	})

	t.Run("URL", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "image/png")
			w.Write([]byte("fake"))
		}))
		defer server.Close()

		client := netsec.NewRestrictedHTTPClient(imageFetchTimeout, true)
		key, ext, data, err := loadImageFromURL(client, server.URL+"/test.png", true, maxCatalogImageURLBodyBytes)
		if err != nil {
			t.Fatalf("loadImageFromURL failed: %v", err)
		}
		if !strings.HasPrefix(key, "url:") || ext != "png" || len(data) != 4 {
			t.Error("loadImageFromURL properties failed")
		}
		if _, _, _, err = loadImageFromURL(
			client,
			"http://127.0.0.1/private.png",
			false,
			maxCatalogImageURLBodyBytes,
		); err == nil {
			t.Fatal("expected SSRF private-host rejection")
		}
		if _, _, _, err = loadImageFromURL(client, server.URL+"/test.png", true, 2); err == nil {
			t.Fatal("expected response too large error")
		}
	})
}

func TestMedia_Catalog(t *testing.T) {
	catalog := &Catalog{
		byKey: make(map[string]Asset),
	}

	name, err := catalog.RegisterImage(shapes.Image{Data: []byte("fake"), Format: "png"})
	if err != nil {
		t.Fatalf("RegisterImage failed: %v", err)
	}
	if name != "image1.png" {
		t.Errorf("expected image1.png, got %q", name)
	}

	if len(catalog.Assets()) != 1 {
		t.Error("Asset count mismatch")
	}
	if catalog.ImageExtensions()[0] != "png" {
		t.Error("ImageExtensions failed")
	}
}

func TestMedia_Catalog_Advanced(t *testing.T) {
	c := &Catalog{byKey: make(map[string]Asset)}

	t.Run("TransitionSound", func(t *testing.T) {
		path := "sound.wav"
		_ = os.WriteFile(path, []byte("wav"), 0600)
		defer os.Remove(path)

		name, err := c.RegisterImage(shapes.Image{Path: path})
		if err != nil {
			t.Error(err)
		}
		if !strings.HasSuffix(name, ".wav") {
			t.Errorf("expected wav, got %q", name)
		}
	})

	t.Run("ResolveAsset_Base64", func(t *testing.T) {
		img := shapes.Image{Data: []byte("b64"), Format: "png"}
		key, _, _, _ := resolveImageAsset(&http.Client{}, img)
		if !strings.HasPrefix(key, "data:") {
			t.Error("Base64 key failed")
		}
	})

	t.Run("NotesMaster", func(t *testing.T) {
		cat := &Catalog{byKey: make(map[string]Asset)}
		master := &elements.NotesMaster{
			Background: &elements.SlideBackground{
				Type:        "picture",
				PictureFill: &shapes.Image{Data: []byte("fake"), Format: "png"},
			},
		}
		_ = addNotesMasterMedia(cat, &http.Client{}, master)
		if len(cat.ordered) != 1 {
			t.Error("NotesMaster background not added")
		}
	})
}

func TestMedia_Unsupported(t *testing.T) {
	if isSupportedMediaExtension("txt") {
		t.Error("txt should not be supported")
	}
	if !isSupportedMediaExtension("png") {
		t.Error("png should be supported")
	}
	if isSupportedMediaExtension("obj") {
		t.Error("obj should not be supported without content-type mapping")
	}
	if isSupportedMediaExtension("fbx") {
		t.Error("fbx should not be supported without content-type mapping")
	}
	if isSupportedMediaExtension("stl") {
		t.Error("stl should not be supported without content-type mapping")
	}
}
