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

	presentationXML := readPresentationXMLFromZip(t, zr)

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

func readPresentationXMLFromZip(t *testing.T, zr *zip.Reader) string {
	t.Helper()
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
		return string(xmlData)
	}
	t.Fatal("ppt/presentation.xml not found")
	return ""
}

func savePresentationXML(t *testing.T, ed *PresentationEditor) string {
	t.Helper()
	out := filepath.Join(t.TempDir(), "protection-roundtrip.pptx")
	if err := ed.Save(out); err != nil {
		t.Fatalf("save deck: %v", err)
	}
	data := readFile(t, out)
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("open output zip: %v", err)
	}
	return readPresentationXMLFromZip(t, zr)
}

// A deck protected with a password to modify must stay protected when it is
// opened and saved again. The verifier is a salted hash, so preserving it never
// requires the plaintext password.
func TestEditorSave_PreservesExistingModifyVerifier(t *testing.T) {
	base := writeDeckFixture(t, "protection-preserve.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	})

	protectedPath := filepath.Join(t.TempDir(), "protected.pptx")
	first, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	first.SetModifyPassword("Secret123!")
	if err := first.Save(protectedPath); err != nil {
		t.Fatalf("save protected deck: %v", err)
	}
	_ = first.Close()

	protectedXML := func() string {
		data := readFile(t, protectedPath)
		zr, zipErr := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if zipErr != nil {
			t.Fatalf("open protected zip: %v", zipErr)
		}
		return readPresentationXMLFromZip(t, zr)
	}()
	original := extractModifyVerifierTag(protectedXML)
	if original == "" {
		t.Fatal("fixture is not protected")
	}

	// Re-open and save without touching protection.
	second, err := OpenPresentationEditor(protectedPath)
	if err != nil {
		t.Fatalf("reopen protected deck: %v", err)
	}
	defer func() { _ = second.Close() }()

	roundTripped := extractModifyVerifierTag(savePresentationXML(t, second))
	if roundTripped != original {
		t.Fatalf("modifyVerifier not preserved:\n orig: %q\n got:  %q", original, roundTripped)
	}
}

func TestEditorSave_EmptyModifyPasswordClearsProtection(t *testing.T) {
	base := writeDeckFixture(t, "protection-clear.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	})

	protectedPath := filepath.Join(t.TempDir(), "protected.pptx")
	first, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	first.SetModifyPassword("Secret123!")
	if err := first.Save(protectedPath); err != nil {
		t.Fatalf("save protected deck: %v", err)
	}
	_ = first.Close()

	second, err := OpenPresentationEditor(protectedPath)
	if err != nil {
		t.Fatalf("reopen protected deck: %v", err)
	}
	defer func() { _ = second.Close() }()

	second.SetModifyPassword("")
	if got := savePresentationXML(t, second); strings.Contains(got, "<p:modifyVerifier") {
		t.Fatal("empty password did not clear p:modifyVerifier")
	}
}

// Replacing the password must not leave the old verifier behind.
func TestEditorSave_NewModifyPasswordReplacesExistingVerifier(t *testing.T) {
	base := writeDeckFixture(t, "protection-replace.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	})

	protectedPath := filepath.Join(t.TempDir(), "protected.pptx")
	first, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	first.SetModifyPassword("Secret123!")
	if err := first.Save(protectedPath); err != nil {
		t.Fatalf("save protected deck: %v", err)
	}
	_ = first.Close()

	data := readFile(t, protectedPath)
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("open protected zip: %v", err)
	}
	original := extractModifyVerifierTag(readPresentationXMLFromZip(t, zr))

	second, err := OpenPresentationEditor(protectedPath)
	if err != nil {
		t.Fatalf("reopen protected deck: %v", err)
	}
	defer func() { _ = second.Close() }()

	second.SetModifyPassword("Different456!")
	got := savePresentationXML(t, second)
	if strings.Count(got, "<p:modifyVerifier") != 1 {
		t.Fatalf("expected exactly one p:modifyVerifier, got %d", strings.Count(got, "<p:modifyVerifier"))
	}
	if replaced := extractModifyVerifierTag(got); replaced == original {
		t.Fatal("verifier was not regenerated for the new password")
	}
}
