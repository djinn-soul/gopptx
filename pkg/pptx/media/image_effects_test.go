package media_test

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestImageEffects(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "test.png")
	if err := os.WriteFile(imgPath, testutil.TinyPNG, 0o600); err != nil {
		t.Fatalf("failed to write test image: %v", err)
	}

	slides := []pptx.SlideContent{
		pptx.NewSlide("Image Effects").
			AddImage(pptx.NewImage(imgPath, 100, 100, 1000, 1000).
				WithShadow(true).
				WithReflection(true)),
	}

	data, err := pptx.CreateWithSlides("Effects Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip.NewReader failed: %v", err)
	}

	found := false
	for _, f := range zr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			found = true
			rc, _ := f.Open()
			content, _ := io.ReadAll(rc)
			if err := rc.Close(); err != nil {
				t.Errorf("failed to close rc: %v", err)
			}
			xml := string(content)

			if !strings.Contains(xml, "<a:effectLst>") {
				t.Errorf("missing effectLst in slide XML")
			}
			if !strings.Contains(xml, "<a:outerShdw") {
				t.Errorf("missing outerShdw in slide XML")
			}
			if !strings.Contains(xml, "<a:ref") {
				t.Errorf("missing reflection (a:ref) in slide XML")
			}
		}
	}
	if !found {
		t.Fatal("slide1.xml not found")
	}
}
