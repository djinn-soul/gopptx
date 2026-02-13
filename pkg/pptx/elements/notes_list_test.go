package elements_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func TestNotesListParagraphs(t *testing.T) {
	slide := elements.NewSlide("Test").
		AddNoteBullet("Bullet one").
		AddNoteBullet("Bullet two").
		AddNoteNumbered("Numbered one").
		AddNoteNumbered("Numbered two").
		AddNoteSubBullet(1, "Sub bullet")

	if len(slide.NotesBody) != 5 {
		t.Fatalf("expected 5 note paragraphs, got %d", len(slide.NotesBody))
	}

	// Verify bullet styles carried through
	if slide.NotesBody[0].Style.BulletStyle != text.BulletStyleBullet {
		t.Errorf("para 0: expected bullet style %q, got %q", text.BulletStyleBullet, slide.NotesBody[0].Style.BulletStyle)
	}
	if slide.NotesBody[2].Style.BulletStyle != text.BulletStyleNumber {
		t.Errorf("para 2: expected number style %q, got %q", text.BulletStyleNumber, slide.NotesBody[2].Style.BulletStyle)
	}
	if slide.NotesBody[4].Style.Level != 1 {
		t.Errorf("para 4: expected level 1, got %d", slide.NotesBody[4].Style.Level)
	}

	// Verify XML output contains bullet markers
	xml := pptxxml.NotesSlide(slide.NotesBody)
	if !strings.Contains(xml, `<a:buChar char="•"/>`) {
		t.Errorf("expected bullet char in notes XML")
	}
	if !strings.Contains(xml, `<a:buAutoNum type="arabicPeriod"/>`) {
		t.Errorf("expected auto-numbering in notes XML")
	}
	if !strings.Contains(xml, `lvl="1"`) {
		t.Errorf("expected level 1 indent in notes XML")
	}
}

func TestNotesPlainTextSync(t *testing.T) {
	slide := elements.NewSlide("Test").
		AddNoteBullet("First").
		AddNoteNumbered("Second")

	if !strings.Contains(slide.Notes, "First") {
		t.Errorf("plain text Notes missing 'First': %q", slide.Notes)
	}
	if !strings.Contains(slide.Notes, "Second") {
		t.Errorf("plain text Notes missing 'Second': %q", slide.Notes)
	}
}
