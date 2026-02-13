package action

import (
	"testing"
)

func TestHyperlinkAction_RelationshipTarget(t *testing.T) {
	tests := []struct {
		name string
		a    HyperlinkAction
		want string
	}{
		{
			name: "URL",
			a:    HyperlinkURL("https://example.com"),
			want: "https://example.com",
		},
		{
			name: "File Link Windows",
			a:    HyperlinkFile(`C:\path\to\file.txt`),
			want: "file:///C:/path/to/file.txt",
		},
		{
			name: "Program Link",
			a:    HyperlinkProgram(`C:\Windows\System32\calc.exe`),
			want: "file:///C:/Windows/System32/calc.exe",
		},
		{
			name: "Email with Subject",
			a:    HyperlinkEmailWithSubject("test@example.com", "Hello World"),
			want: "mailto:test@example.com?subject=Hello+World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.RelationshipTarget(); got != tt.want {
				t.Errorf("HyperlinkAction.RelationshipTarget() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHyperlinkAction_IsExternal(t *testing.T) {
	tests := []struct {
		name string
		a    HyperlinkAction
		want bool
	}{
		{"URL", HyperlinkURL("http://test.com"), true},
		{"Slide", HyperlinkSlide(1), false},
		{"File", HyperlinkFile("test.txt"), true},
		{"Program", HyperlinkProgram("calc.exe"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.IsExternal(); got != tt.want {
				t.Errorf("HyperlinkAction.IsExternal() = %v, want %v", got, tt.want)
			}
		})
	}
}
