package editor

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestEditorSave_WithModifyPassword_WritesModifyVerifier(t *testing.T) {
	base := writeDeckFixture(t, "protection-base.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	})

	ed, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	ed.Metadata().Protection.ModifyPassword = "Secret123!"
	out := filepath.Join(t.TempDir(), "protection-output.pptx")
	if err := ed.Save(out); err != nil {
		t.Fatalf("save protected deck: %v", err)
	}

	data := readFile(t, out)
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("open output zip: %v", err)
	}

	var presentationXML string
	for _, f := range zr.File {
		if f.Name != "ppt/presentation.xml" {
			continue
		}
		rc, openErr := f.Open()
		if openErr != nil {
			t.Fatalf("open presentation.xml: %v", openErr)
		}
		xmlData, readErr := io.ReadAll(rc)
		_ = rc.Close()
		if readErr != nil {
			t.Fatalf("read presentation.xml: %v", readErr)
		}
		presentationXML = string(xmlData)
		break
	}
	if presentationXML == "" {
		t.Fatal("ppt/presentation.xml not found")
	}

	required := []string{
		"<p:modifyVerifier",
		`cryptProviderType="rsaAES"`,
		`cryptAlgorithmClass="hash"`,
		`cryptAlgorithmSid="14"`,
		`spinCount="100000"`,
		`saltData="`,
		`hashData="`,
	}
	for _, fragment := range required {
		if !strings.Contains(presentationXML, fragment) {
			t.Fatalf("missing %q in presentation.xml", fragment)
		}
	}
}

