package text

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
)

func TestTextRunValidation(t *testing.T) {
	tests := []struct {
		name    string
		run     Run
		wantErr bool
	}{
		{
			name:    "valid run",
			run:     NewRun("hello").WithBold(true).WithColor("FF0000"),
			wantErr: false,
		},
		{
			name:    "negative size",
			run:     NewRun("hello").WithSizePt(-1),
			wantErr: true,
		},
		{
			name:    "invalid color",
			run:     NewRun("hello").WithColor("red"),
			wantErr: true,
		},
		{
			name:    "conflicting baseline",
			run:     Run{Text: "hello", Subscript: true, Superscript: true},
			wantErr: true,
		},
		{
			name: "valid hyperlink",
			run:  NewRun("link").WithHyperlink(action.NewHyperlink(action.HyperlinkURL("http://example.com"))),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestTextParagraphStyleValidation(t *testing.T) {
	tests := []struct {
		name    string
		style   ParagraphStyle
		wantErr bool
	}{
		{
			name:    "valid style",
			style:   NewParagraphStyle().WithAlignCenter().WithSpaceBeforePt(10),
			wantErr: false,
		},
		{
			name:    "invalid align",
			style:   ParagraphStyle{Align: "diagonal"},
			wantErr: true,
		},
		{
			name:    "negative spacing",
			style:   NewParagraphStyle().WithSpaceAfterPt(-1),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.style.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestNormalizeTextRuns(t *testing.T) {
	runs := []Run{
		{Text: "a", Bold: true},
		{Text: "b", Bold: true},
		{Text: "", Bold: true},
		{Text: "c", Bold: false},
	}
	normalized := NormalizeRuns(runs)
	if len(normalized) != 2 {
		t.Errorf("expected 2 runs, got %d", len(normalized))
	}
	if normalized[0].Text != "ab" {
		t.Errorf("expected 'ab', got %q", normalized[0].Text)
	}
}

func TestNormalizeBulletStyle(t *testing.T) {
	checks := []struct {
		input string
		want  string
	}{
		{"Numbered", BulletStyleNumber},
		{"letter-lower", BulletStyleLetterLower},
		{"roman", BulletStyleRomanUpper},
		{"NONE", BulletStyleNone},
		{"custom", BulletStyleCustom},
	}
	for _, c := range checks {
		if got := NormalizeBulletStyle(c.input); got != c.want {
			t.Errorf("NormalizeBulletStyle(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
