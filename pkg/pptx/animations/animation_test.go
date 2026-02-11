package animations

import (
	"strings"
	"testing"
)

func TestAnimationValidation(t *testing.T) {
	tests := []struct {
		name    string
		anim    Animation
		wantErr bool
	}{
		{
			name:    "valid animation",
			anim:    NewAnimation(1, AnimationEntranceFade),
			wantErr: false,
		},
		{
			name:    "invalid shape index",
			anim:    NewAnimation(0, AnimationEntranceFade),
			wantErr: true,
		},
		{
			name:    "invalid effect",
			anim:    NewAnimation(1, AnimationEffect("invalid")),
			wantErr: true,
		},
		{
			name:    "invalid trigger",
			anim:    NewAnimation(1, AnimationEntranceFade).WithTrigger(AnimationTrigger("invalid")),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.anim.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Animation.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnimationXML(t *testing.T) {
	anim := NewAnimation(1, AnimationEntranceFade).
		WithDuration(1000).
		WithDelay(200).
		WithRepeat(2).
		WithAutoReverse(true)

	xml := anim.XML(10, 3)

	if !strings.Contains(xml, `presetID="10"`) {
		t.Errorf("expected presetID 10 (fade), got: %s", xml)
	}
	if !strings.Contains(xml, `dur="1000"`) {
		t.Errorf("expected duration 1000, got: %s", xml)
	}
	if !strings.Contains(xml, `delay="200"`) {
		t.Errorf("expected delay 200, got: %s", xml)
	}
	if !strings.Contains(xml, `repeatCount="2000"`) {
		t.Errorf("expected repeatCount 2000, got: %s", xml)
	}
	if !strings.Contains(xml, `autoRev="1"`) {
		t.Errorf("expected autoRev 1, got: %s", xml)
	}
	if !strings.Contains(xml, `spid="3"`) {
		t.Errorf("expected spid 3, got: %s", xml)
	}
}

func TestAnimationPresetSubtype(t *testing.T) {
	tests := []struct {
		effect    AnimationEffect
		direction AnimationDirection
		want      int
	}{
		{AnimationEntranceFlyIn, AnimationDirDown, 1},
		{AnimationEntranceFlyIn, AnimationDirLeft, 2},
		{AnimationEntranceFlyIn, AnimationDirRight, 4},
		{AnimationEntranceFlyIn, AnimationDirUp, 8},
		{AnimationEntranceFlyIn, AnimationDirDownLeft, 3},
		{AnimationEntranceWipe, AnimationDirUp, 4},
		{AnimationEntranceSplit, AnimationDirIn, 1},
	}

	for _, tt := range tests {
		anim := NewAnimation(1, tt.effect)
		anim.Direction = tt.direction
		if got := anim.PresetSubtype(); got != tt.want {
			t.Errorf("Animation(%s, %s).PresetSubtype() = %d, want %d", tt.effect, tt.direction, got, tt.want)
		}
	}
}
