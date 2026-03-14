package action

import (
	"testing"
)

func TestHyperlinkAction_Methods(t *testing.T) {
	t.Run("URL", func(t *testing.T) {
		a := HyperlinkURL("https://example.com")
		if a.Type != HyperlinkActionURL || a.URL != "https://example.com" {
			t.Error("HyperlinkURL failed")
		}
		if !a.IsExternal() {
			t.Error("URL should be external")
		}
		if a.RelationshipTarget() != "https://example.com" {
			t.Errorf("RelationshipTarget failed: %s", a.RelationshipTarget())
		}
	})

	t.Run("Slide", func(t *testing.T) {
		a := HyperlinkSlide(5)
		if a.Type != HyperlinkActionSlide || a.SlideNumber != 5 {
			t.Error("HyperlinkSlide failed")
		}
		if a.IsExternal() {
			t.Error("Slide should not be external")
		}
		if a.RelationshipTarget() != "slide5.xml" {
			t.Errorf("RelationshipTarget failed: %s", a.RelationshipTarget())
		}
		if a.ActionType() != "ppaction://hlinksldjump" {
			t.Errorf("ActionType failed: %s", a.ActionType())
		}
	})

	t.Run("FirstSlide", func(t *testing.T) {
		a := HyperlinkFirstSlide()
		if a.Type != HyperlinkActionFirstSlide {
			t.Error("HyperlinkFirstSlide failed")
		}
		if a.RelationshipTarget() != "ppaction://hlinkshowjump?jump=firstslide" {
			t.Error("RelationshipTarget failed")
		}
	})

	t.Run("LastSlide", func(t *testing.T) {
		a := HyperlinkLastSlide()
		if a.RelationshipTarget() != "ppaction://hlinkshowjump?jump=lastslide" {
			t.Error("RelationshipTarget failed")
		}
	})

	t.Run("NextSlide", func(t *testing.T) {
		a := HyperlinkNextSlide()
		if a.RelationshipTarget() != "ppaction://hlinkshowjump?jump=nextslide" {
			t.Error("RelationshipTarget failed")
		}
	})

	t.Run("PreviousSlide", func(t *testing.T) {
		a := HyperlinkPreviousSlide()
		if a.RelationshipTarget() != "ppaction://hlinkshowjump?jump=previousslide" {
			t.Error("RelationshipTarget failed")
		}
	})

	t.Run("EndShow", func(t *testing.T) {
		a := HyperlinkEndShow()
		if a.RelationshipTarget() != "ppaction://hlinkshowjump?jump=endshow" {
			t.Error("RelationshipTarget failed")
		}
	})

	t.Run("Email", func(t *testing.T) {
		a := HyperlinkEmail("test@example.com")
		if a.RelationshipTarget() != "mailto:test@example.com" {
			t.Errorf("RelationshipTarget failed: %s", a.RelationshipTarget())
		}

		a2 := HyperlinkEmailWithSubject("test@example.com", "Hello World")
		if a2.RelationshipTarget() != "mailto:test@example.com?subject=Hello+World" {
			t.Errorf("RelationshipTarget with subject failed: %s", a2.RelationshipTarget())
		}
	})

	t.Run("File", func(t *testing.T) {
		a := HyperlinkFile(`C:\test.txt`)
		if a.RelationshipTarget() != "file:///C:/test.txt" {
			t.Errorf("RelationshipTarget failed: %s", a.RelationshipTarget())
		}
	})

	t.Run("Program", func(t *testing.T) {
		a := HyperlinkProgram("notepad.exe")
		if a.RelationshipTarget() != "file:///notepad.exe" {
			t.Errorf("RelationshipTarget failed: %s", a.RelationshipTarget())
		}
		if a.ActionType() != "ppaction://program" {
			t.Error("ActionType failed")
		}
	})

	t.Run("Default", func(t *testing.T) {
		a := HyperlinkAction{Type: "unknown"}
		if a.RelationshipTarget() != "" {
			t.Error("RelationshipTarget should be empty for unknown")
		}
		if a.ActionType() != "" {
			t.Error("ActionType should be empty for unknown")
		}
	})
}

func TestHyperlink_Methods(t *testing.T) {
	h := NewHyperlink(HyperlinkURL("https://example.com")).
		WithTooltip("Tooltip").
		WithHighlightClick(false)

	if h.Tooltip != "Tooltip" {
		t.Error("WithTooltip failed")
	}
	if h.HighlightClick {
		t.Error("WithHighlightClick failed")
	}
}

func TestHyperlink_Validate(t *testing.T) {
	tests := []struct {
		name    string
		h       Hyperlink
		wantErr bool
	}{
		{"Valid URL", NewHyperlink(HyperlinkURL("http://x.com")), false},
		{"Empty URL", NewHyperlink(HyperlinkURL("")), true},
		{"Valid File", NewHyperlink(HyperlinkFile("x.txt")), false},
		{"Empty File", NewHyperlink(HyperlinkFile("")), true},
		{"Traversal File", NewHyperlink(HyperlinkFile("../x.txt")), true},
		{"Valid Program", NewHyperlink(HyperlinkProgram("x.exe")), false},
		{"Empty Program", NewHyperlink(HyperlinkProgram("")), true},
		{"Traversal Program", NewHyperlink(HyperlinkProgram("../x.exe")), true},
		{"Invalid File Scheme", NewHyperlink(HyperlinkFile("https://example.com/x.txt")), true},
		{"Invalid Program Scheme", NewHyperlink(HyperlinkProgram("mailto:test@example.com")), true},
		{"Valid File URI", NewHyperlink(HyperlinkFile("file:///C:/x.txt")), false},
		{"File URI With Query", NewHyperlink(HyperlinkFile("file:///C:/x.txt?x=1")), true},
		{"File URI With Fragment", NewHyperlink(HyperlinkFile("file:///C:/x.txt#frag")), true},
		{"Escaped Traversal File URI", NewHyperlink(HyperlinkFile("file:///C:/safe/%2e%2e/evil.txt")), true},
		{"Restricted Windows System Path", NewHyperlink(HyperlinkProgram(`C:\Windows\System32\cmd.exe`)), true},
		{"Restricted Unix System Path", NewHyperlink(HyperlinkFile("/etc/passwd")), true},
		{"Remote File Host Disallowed", NewHyperlink(HyperlinkFile("file://server/share/doc.txt")), true},
		{"Valid Email", NewHyperlink(HyperlinkEmail("a@b.com")), false},
		{"Empty Email", NewHyperlink(HyperlinkEmail("")), true},
		{"Valid Slide", NewHyperlink(HyperlinkSlide(1)), false},
		{"Unknown", NewHyperlink(HyperlinkAction{Type: "???"}), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.h.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
