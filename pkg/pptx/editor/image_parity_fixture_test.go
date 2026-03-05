package editor

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestImageParityFixturePackagingSizingCropping(t *testing.T) {
	base := writeDeckFixture(t, "base.pptx", []elements.SlideContent{elements.NewSlide("Slide 1")})
	ed, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}

	flipH := true
	rotation := 15.0
	base64Data := "data:image/png;base64," + base64.StdEncoding.EncodeToString(testutil.TinyPNG())
	if _, err := ed.AddImageFromBase64(
		0,
		base64Data,
		"",
		111,
		222,
		333,
		444,
		&common.ShapeUpdate{
			Crop: &common.ImageCrop{
				Left:   0.05,
				Right:  0.10,
				Top:    0.15,
				Bottom: 0.20,
			},
			Rotation: &rotation,
			FlipH:    &flipH,
		},
	); err != nil {
		_ = ed.Close()
		t.Fatalf("add image from base64: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(testutil.TinyPNG())
	}))
	defer server.Close()

	if _, err := ed.AddImageFromURL(0, server.URL+"/img.png", 555, 666, 777, 888, nil); err != nil {
		_ = ed.Close()
		t.Fatalf("add image from URL: %v", err)
	}

	out := filepath.Join(t.TempDir(), "image-parity.pptx")
	if err := ed.Save(out); err != nil {
		_ = ed.Close()
		t.Fatalf("save edited pptx: %v", err)
	}
	_ = ed.Close()

	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output pptx: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("open output pptx zip: %v", err)
	}

	mediaCount := 0
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "ppt/media/image") {
			mediaCount++
		}
	}
	if mediaCount != 1 {
		t.Fatalf("expected deduped single media image part, got %d", mediaCount)
	}

	relsXML := testutil.ReadZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if strings.Count(relsXML, `/relationships/image`) != 1 {
		t.Fatalf("expected single image relationship entry, got rels: %s", relsXML)
	}
	if !strings.Contains(relsXML, `Target="../media/image1.png"`) {
		t.Fatalf("expected image target in rels, got: %s", relsXML)
	}

	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")
	needs := []string{
		`<a:off x="111" y="222"/>`,
		`<a:ext cx="333" cy="444"/>`,
		`<a:srcRect l="5000" r="10000" t="15000" b="20000"/>`,
		`<a:xfrm rot="900000" flipH="1">`,
		`<a:off x="555" y="666"/>`,
		`<a:ext cx="777" cy="888"/>`,
	}
	for _, needle := range needs {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide xml, got: %s", needle, slideXML)
		}
	}
}
