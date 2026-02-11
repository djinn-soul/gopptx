package transitions

import (
	"testing"
)

func TestTransitionOptionsValidation(t *testing.T) {
	tests := []struct {
		name    string
		opts    TransitionOptions
		wantErr bool
	}{
		{
			name:    "valid fade",
			opts:    TransitionOptions{Type: TransitionFade},
			wantErr: false,
		},
		{
			name:    "invalid type",
			opts:    TransitionOptions{Type: TransitionType("invalid")},
			wantErr: true,
		},
		{
			name:    "invalid orientation",
			opts:    TransitionOptions{Type: TransitionFade, Orientation: TransitionOrientHorizontal},
			wantErr: true,
		},
		{
			name:    "valid orientation for split",
			opts:    TransitionOptions{Type: TransitionSplit, Orientation: TransitionOrientHorizontal},
			wantErr: false,
		},
		{
			name:    "invalid direction for push",
			opts:    TransitionOptions{Type: TransitionPush, Direction: TransitionDirection("invalid")},
			wantErr: true,
		},
		{
			name:    "valid direction for push",
			opts:    TransitionOptions{Type: TransitionPush, Direction: TransitionDirUp},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.opts.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("TransitionOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransitionXML(t *testing.T) {
	tests := []struct {
		name string
		opts TransitionOptions
		want string
	}{
		{
			name: "fade",
			opts: TransitionOptions{Type: TransitionFade, DurationMS: 2000},
			want: `<p:transition dur="2000"><p:fade/></p:transition>`,
		},
		{
			name: "push right",
			opts: TransitionOptions{Type: TransitionPush, Direction: TransitionDirRight},
			want: `<p:transition><p:push dir="r"/></p:transition>`,
		},
		{
			name: "split horizontal out",
			opts: TransitionOptions{Type: TransitionSplit, Direction: TransitionDirOut, Orientation: TransitionOrientHorizontal},
			want: `<p:transition><p:split dir="out" orient="horz"/></p:transition>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.XML()
			if got != tt.want {
				t.Errorf("TransitionOptions.XML() = %v, want %v", got, tt.want)
			}
		})
	}
}
