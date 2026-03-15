package animations

import (
	"testing"
)

func TestAnimationPresetSubtypeExtra(t *testing.T) {
	tests := []struct {
		name      string
		effect    AnimationEffect
		direction AnimationDirection
		expected  int
	}{
		{"FlyIn UpLeft", AnimationEntranceFlyIn, AnimationDirUpLeft, presetSubtypeFromBottomRight},
		{"FlyIn UpRight", AnimationEntranceFlyIn, AnimationDirUpRight, presetSubtypeFromBottomLeft},
		{"FlyIn DownRight", AnimationEntranceFlyIn, AnimationDirDownRight, presetSubtypeFromTopLeft},

		{"Wipe Down", AnimationEntranceWipe, AnimationDirDown, presetSubtypeFromBottom},
		{"Wipe Left", AnimationEntranceWipe, AnimationDirLeft, presetSubtypeFromTop},
		{"Wipe Right", AnimationEntranceWipe, AnimationDirRight, presetSubtypeFromRight},
		
		{"Split Out", AnimationEntranceSplit, AnimationDirOut, presetSubtypeSplitOut},
		{"Split Default", AnimationEntranceSplit, "invalid", presetSubtypeSplitOut},

		{"Other", AnimationEntranceFade, AnimationDirUp, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Animation{Effect: tt.effect, Direction: tt.direction}
			if got := a.PresetSubtype(); got != tt.expected {
				t.Errorf("PresetSubtype() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAnimationPresetClassExtra(t *testing.T) {
	tests := []struct {
		effect   AnimationEffect
		expected string
	}{
		{AnimationEntranceFade, classEntr},
		{AnimationExitFadeOut, "exit"},
		{AnimationEmphasisPulse, "emph"},
		{AnimationPathLines, "path"},
		{"invalid", classEntr},
	}

	for _, tt := range tests {
		a := Animation{Effect: tt.effect}
		if got := a.PresetClass(); got != tt.expected {
			t.Errorf("PresetClass(%s) = %s, want %s", tt.effect, got, tt.expected)
		}
	}
}

func TestAnimationTriggerValidateExtra(t *testing.T) {
	if err := AnimationTrigger("bad").Validate(); err == nil {
		t.Error("Expected error for bad trigger")
	}
}

func TestToAnimationExtra(t *testing.T) {
	a := NewAnimation(1, AnimationEntranceFade)
	if a.ToAnimation() != a {
		t.Error("ToAnimation should return self")
	}
}

func TestAnimationPresetIDExtra(t *testing.T) {
	// Test a few paths in PresetID
	tests := []struct {
		effect   AnimationEffect
		expected uint32
	}{
		{AnimationEntranceAppear, presetIDAppear},
		{AnimationEntranceFloat, presetIDFloat},
		{AnimationEntranceSplit, presetIDSplit},
		{AnimationEntranceShape, presetIDShape},
		{AnimationEmphasisSpin, presetIDSpin},
		{AnimationPathArcs, presetIDPathArcs},
		{"unknown", 0},
	}

	for _, tt := range tests {
		a := Animation{Effect: tt.effect}
		if got := a.PresetID(); got != tt.expected {
			t.Errorf("PresetID(%s) = %d, want %d", tt.effect, got, tt.expected)
		}
	}
}
