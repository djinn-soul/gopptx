package pptx

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPresentationEditorAddUpdateRemoveSave(t *testing.T) {
	initial := []SlideContent{
		NewSlide("Intro").AddBullet("Original"),
		NewSlide("Keep").AddBullet("To be removed"),
	}
	initialPath := writeDeckFixture(t, "initial.pptx", initial)

	editor, err := OpenPresentationEditor(initialPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	if editor.SlideCount() != 2 {
		t.Fatalf("expected 2 slides, got %d", editor.SlideCount())
	}

	if _, err := editor.AddSlide(NewSlide("Added").AddBullet("new bullet")); err != nil {
		t.Fatalf("add slide: %v", err)
	}
	if err := editor.UpdateSlide(0, NewSlide("Updated Intro").AddBullet("Updated")); err != nil {
		t.Fatalf("update slide: %v", err)
	}
	if err := editor.RemoveSlide(1); err != nil {
		t.Fatalf("remove slide: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "edited.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save edited deck: %v", err)
	}

	edited, err := OpenPresentationEditor(outPath)
	if err != nil {
		t.Fatalf("reopen edited deck: %v", err)
	}
	if edited.SlideCount() != 2 {
		t.Fatalf("expected 2 slides after edit, got %d", edited.SlideCount())
	}

	slides := edited.Slides()
	if slides[0].Title != "Updated Intro" {
		t.Fatalf("unexpected slide[0] title: %q", slides[0].Title)
	}
	if slides[1].Title != "Added" {
		t.Fatalf("unexpected slide[1] title: %q", slides[1].Title)
	}
}

func TestPresentationEditorPreservesNonEditedParts(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "sample.png")
	if err := os.WriteFile(imgPath, tinyPNG, 0o600); err != nil {
		t.Fatalf("write image fixture: %v", err)
	}

	deck := []SlideContent{
		NewSlide("Image Slide").AddImage(NewImage(imgPath, Inches(1), Inches(1), Inches(2), Inches(2))),
		NewSlide("Editable").AddBullet("old"),
	}
	originalPath := writeDeckFixture(t, "original-with-image.pptx", deck)

	editor, err := OpenPresentationEditor(originalPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	if err := editor.UpdateSlide(1, NewSlide("Editable").AddBullet("updated text")); err != nil {
		t.Fatalf("update text-only slide: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "edited-with-image.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save edited deck: %v", err)
	}

	origMedia := readZipFileBytes(t, originalPath, "ppt/media/image1.png")
	editedMedia := readZipFileBytes(t, outPath, "ppt/media/image1.png")
	if !bytes.Equal(origMedia, editedMedia) {
		t.Fatalf("expected untouched media bytes to be preserved")
	}

	origRel := string(readZipFileBytes(t, originalPath, "ppt/slides/_rels/slide1.xml.rels"))
	editedRel := string(readZipFileBytes(t, outPath, "ppt/slides/_rels/slide1.xml.rels"))
	if origRel != editedRel {
		t.Fatalf("expected untouched slide relationships to be preserved")
	}
}

func TestPresentationEditorRejectsUpdateForSlideWithExternalRelationships(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "sample.png")
	if err := os.WriteFile(imgPath, tinyPNG, 0o600); err != nil {
		t.Fatalf("write image fixture: %v", err)
	}

	deck := []SlideContent{
		NewSlide("Image Slide").AddImage(NewImage(imgPath, Inches(1), Inches(1), Inches(2), Inches(2))),
	}
	path := writeDeckFixture(t, "image-only.pptx", deck)

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	err = editor.UpdateSlide(0, NewSlide("Replacement").AddBullet("text"))
	if err == nil {
		t.Fatalf("expected unsupported relationship error")
	}
	if !strings.Contains(err.Error(), "unsupported relationship type") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPresentationEditorMergeFromFile(t *testing.T) {
	destPath := writeDeckFixture(t, "dest.pptx", []SlideContent{
		NewSlide("Dest 1").AddBullet("a"),
	})
	sourcePath := writeDeckFixture(t, "source.pptx", []SlideContent{
		NewSlide("Source 1").AddBullet("b"),
		NewSlide("Source 2").AddBullet("c"),
	})

	editor, err := OpenPresentationEditor(destPath)
	if err != nil {
		t.Fatalf("open dest editor: %v", err)
	}
	if err := editor.MergeFromFile(sourcePath); err != nil {
		t.Fatalf("merge from file: %v", err)
	}
	if editor.SlideCount() != 3 {
		t.Fatalf("expected 3 slides after merge, got %d", editor.SlideCount())
	}

	outPath := filepath.Join(t.TempDir(), "merged.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save merged deck: %v", err)
	}
	reopened, err := OpenPresentationEditor(outPath)
	if err != nil {
		t.Fatalf("reopen merged deck: %v", err)
	}
	if reopened.SlideCount() != 3 {
		t.Fatalf("expected reopened merged deck to have 3 slides, got %d", reopened.SlideCount())
	}
}

func TestOpenPresentationEditorRejectsCorruptPackage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "corrupt.pptx")
	if err := writeZipFixture(path, map[string]string{
		"docProps/core.xml": `<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties"/>`,
	}); err != nil {
		t.Fatalf("write corrupt zip fixture: %v", err)
	}

	_, err := OpenPresentationEditor(path)
	if err == nil {
		t.Fatalf("expected error for missing required package parts")
	}
	if !strings.Contains(err.Error(), "missing required package part") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPresentationEditorRejectsAddSlideWithUnsupportedAssets(t *testing.T) {
	path := writeDeckFixture(t, "base.pptx", []SlideContent{NewSlide("Base")})
	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}

	slide := NewSlide("Unsupported")
	chart := NewBarChart([]string{"A"}, []float64{1})
	slide.Chart = &chart
	_, err = editor.AddSlide(slide)
	if err == nil {
		t.Fatalf("expected unsupported chart authoring error")
	}
	if !strings.Contains(err.Error(), "does not support chart authoring yet") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func writeDeckFixture(t *testing.T, name string, slides []SlideContent) string {
	t.Helper()
	data, err := CreateWithSlides("Editor Fixture", slides)
	if err != nil {
		t.Fatalf("create fixture deck: %v", err)
	}
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write fixture deck: %v", err)
	}
	return path
}

func readZipFileBytes(t *testing.T, zipPath string, entryName string) []byte {
	t.Helper()

	data, err := os.ReadFile(zipPath)
	if err != nil {
		t.Fatalf("read zip file %s: %v", zipPath, err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("open zip %s: %v", zipPath, err)
	}
	for _, file := range zr.File {
		if file.Name != entryName {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			t.Fatalf("open zip entry %s: %v", entryName, err)
		}
		defer rc.Close()
		content, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("read zip entry %s: %v", entryName, err)
		}
		return content
	}
	t.Fatalf("zip entry %s not found in %s", entryName, zipPath)
	return nil
}

func writeZipFixture(path string, files map[string]string) error {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			_ = zw.Close()
			return err
		}
		if _, err := w.Write([]byte(content)); err != nil {
			_ = zw.Close()
			return err
		}
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o600)
}
