package styling_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestSlideStyleValidation(t *testing.T) {
	tests := []struct {
		name    string
		slide   pptx.SlideContent
		wantErr string
	}{
		{
			name:    "Invalid Title Size",
			slide:   pptx.NewSlide("Test").WithTitleSize(500),
			wantErr: "title size must be between 1 and 400 pt",
		},
		{
			name:    "Invalid Title Color",
			slide:   pptx.NewSlide("Test").WithTitleColor("invalid"),
			wantErr: "title color must be 6-digit RGB hex",
		},
		{
			name:    "Invalid Content Size",
			slide:   pptx.NewSlide("Test").WithContentSize(-1),
			wantErr: "content size must be between 1 and 400 pt",
		},
		{
			name:    "Invalid Content Color",
			slide:   pptx.NewSlide("Test").WithContentColor("GG0000"),
			wantErr: "content color must be 6-digit RGB hex",
		},
		{
			name: "Valid Styles",
			slide: pptx.NewSlide("Test").
				WithTitleSize(24).
				WithTitleColor("FF0000").
				WithContentSize(18).
				WithContentColor("00FF00"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.slide.Validate(1)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("Validate() expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
