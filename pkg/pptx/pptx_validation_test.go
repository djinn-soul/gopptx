package pptx

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestSlideValidation_Typography(t *testing.T) {
	// Invalid title size
	s1 := NewSlide("T")
	s1.TitleSize = 500
	if err := s1.Validate(1); err == nil || !strings.Contains(err.Error(), "title size must be between") {
		t.Errorf("Expected title size error, got %v", err)
	}

	// Invalid title color
	s2 := NewSlide("T")
	s2.TitleColor = "red"
	if err := s2.Validate(1); err == nil || !strings.Contains(err.Error(), "title color must be 6-digit RGB hex") {
		t.Errorf("Expected title color error, got %v", err)
	}

	// Invalid content size
	s3 := NewSlide("T")
	s3.ContentSize = 500
	if err := s3.Validate(1); err == nil || !strings.Contains(err.Error(), "content size must be between") {
		t.Errorf("Expected content size error, got %v", err)
	}

	// Invalid content color
	s4 := NewSlide("T")
	s4.ContentColor = "blue"
	if err := s4.Validate(1); err == nil || !strings.Contains(err.Error(), "content color must be 6-digit RGB hex") {
		t.Errorf("Expected content color error, got %v", err)
	}
}

func TestSlideValidation_Alignment(t *testing.T) {
	// Invalid title align
	s1 := NewSlide("T")
	s1.TitleAlign = "top"
	if err := s1.Validate(1); err == nil || !strings.Contains(err.Error(), "invalid title alignment") {
		t.Errorf("Expected title alignment error, got %v", err)
	}

	// Invalid content valign
	s2 := NewSlide("T")
	s2.ContentVAlign = "left"
	if err := s2.Validate(1); err == nil || !strings.Contains(err.Error(), "invalid content vertical alignment") {
		t.Errorf("Expected content valign error, got %v", err)
	}
}

func TestSlideValidation_Animations(t *testing.T) {
	// First animation trigger with/after previous
	s1 := NewSlide("T")
	// Add a real shape to pass index validation
	s1.Shapes = append(s1.Shapes, shapes.NewRectangle(0, 0, 100, 100).ToShape())

	anim1 := animations.NewAnimation(1, animations.AnimationEntranceFade).WithTrigger(animations.AnimationWithPrevious)
	s1.Animations = append(s1.Animations, anim1)
	if err := s1.Validate(1); err == nil || !strings.Contains(
		err.Error(),
		"first animation trigger cannot be with/after previous",
	) {
		t.Errorf("Expected first animation trigger error, got %v", err)
	}

	// Out of bounds shape index
	s2 := NewSlide("T")
	anim2 := animations.NewAnimation(10, animations.AnimationEntranceFade)
	s2.Animations = append(s2.Animations, anim2)
	if err := s2.Validate(1); err == nil || !strings.Contains(
		err.Error(),
		"targets shape index 10, but slide only has 0",
	) {
		t.Errorf("Expected animation shape index error, got %v", err)
	}
}

func TestSlideValidation_EmptyTitle(t *testing.T) {
	// Empty title with non-blank layout
	s1 := NewSlide("")
	if err := s1.Validate(1); err == nil || !strings.Contains(err.Error(), "title cannot be empty") {
		t.Errorf("Expected empty title error, got %v", err)
	}

	// Empty title with blank layout is OK
	s2 := NewSlide("")
	s2.Layout = SlideLayoutBlank
	if err := s2.Validate(1); err != nil {
		t.Errorf("Empty title with blank layout should be OK, got %v", err)
	}
}

func TestSlideValidation_Bullets(t *testing.T) {
	s1 := NewSlide("T").AddBullet("")
	if err := s1.Validate(1); err == nil || !strings.Contains(err.Error(), "bullet cannot be empty") {
		t.Errorf("Expected empty bullet error, got %v", err)
	}
}

func TestCreateWithMetadata_ExtraPathways(t *testing.T) {
	// Invalid Notes Master Background
	meta := Metadata{Metadata: common.Metadata{Title: "T"}}
	nm := elements.NewNotesMaster()
	bg := elements.NewSolidBackground("invalid")
	nm.Background = &bg
	meta.NotesMaster = nm

	_, err := CreateWithMetadata(meta, []SlideContent{NewSlide("S")})
	if err == nil || !strings.Contains(err.Error(), "invalid notes master background") {
		t.Errorf("Expected notes master background error, got %v", err)
	}

	// Encryption password (not empty but whitespace only should be treated as empty)
	meta.Protection.EncryptPassword = "   "
	meta.NotesMaster = nil
	_, err = CreateWithMetadata(meta, []SlideContent{NewSlide("S")})
	if err != nil {
		t.Errorf("CreateWithMetadata should ignore whitespace password, got %v", err)
	}
}

func TestValidateAndRepair_Pathways(t *testing.T) {
	// 1. Valid PPTX (covers editor success path)
	data, err := Create("Valid", 1)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	issues, err := Validate(data)
	if err != nil {
		t.Errorf("Validate failed on valid PPTX: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues on valid PPTX, got %d", len(issues))
	}

	repairedData, res, err := Repair(data)
	if err != nil {
		t.Errorf("Repair failed on valid PPTX: %v", err)
	}
	if len(res.IssuesRepaired) != 0 {
		t.Errorf("Expected 0 repaired issues on valid PPTX, got %d", len(res.IssuesRepaired))
	}
	if len(repairedData) == 0 {
		t.Error("Repair returned empty data for valid input")
	}

	// 2. Corrupted Zip (Validate failure path)
	_, err = Validate([]byte("not a zip"))
	if err == nil {
		t.Error("Expected error for non-zip data in Validate")
	}
	_, _, err = Repair([]byte("not a zip"))
	if err == nil {
		t.Error("Expected error for non-zip data in Repair")
	}

	// 3. Zip missing PPTX structure (Validate fallback path)
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("garbage.txt")
	w.Write([]byte("data"))
	zw.Close()

	// OpenPartStoreFromBytes should succeed on valid zip, but Validate will find issues
	issues, err = Validate(buf.Bytes())
	if err != nil {
		t.Errorf("Validate should not error on valid zip with missing PPTX parts, but got %v", err)
	}
	if len(issues) == 0 {
		t.Error("Expected structural issues for zip missing PPTX parts")
	}
}

func TestWriteFile_Convenience(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.pptx")
	err := WriteFile(path, "Title", []SlideContent{NewSlide("S1")})
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("File not written: %v", err)
	}
}
