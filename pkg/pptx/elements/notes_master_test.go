package elements

import "testing"

func TestNotesMasterValidateBackground(t *testing.T) {
	bg := SlideBackground{Type: SlideBackgroundPicture}
	m := NewNotesMaster().WithBackground(bg)
	if err := m.Validate(); err == nil {
		t.Fatalf("expected invalid notes master background error")
	}

	valid := NewNotesMaster().WithBackground(NewSolidBackground("F0F0F0"))
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid notes master background, got %v", err)
	}
}
