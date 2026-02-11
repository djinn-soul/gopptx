package pptx

import (
	"testing"
)

func TestHyperlinkValidation(t *testing.T) {
	tests := []struct {
		name    string
		action  HyperlinkAction
		wantErr bool
	}{
		{"valid URL", HyperlinkURL("https://example.com"), false},
		{"empty URL", HyperlinkURL(""), true},
		{"URL without scheme", HyperlinkURL("example.com"), true},
		{"valid slide", HyperlinkSlide(1), false},
		{"invalid slide 0", HyperlinkSlide(0), true},
		{"valid email", HyperlinkEmail("test@example.com"), false},
		{"invalid email", HyperlinkEmail("invalid"), true},
		{"valid file", HyperlinkFile("C:\\docs\\file.pdf"), false},
		{"empty file", HyperlinkFile(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHyperlinkAction(tt.action, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHyperlinkAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
