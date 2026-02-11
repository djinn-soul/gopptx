package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

// TestAnimationXMLDeterminism verifies that concurrent animation processing produces deterministic output order.
func TestAnimationXMLDeterminism(t *testing.T) {
	// Create a slide with enough animations to trigger parallel processing (> 50)
	slide := NewSlide("Deterministic Anim").
		AddShape(NewShape("rect", 0, 0, 100, 100)) // shape ID 3

	// Add 60 animations
	for i := 0; i < 60; i++ {
		// Alternate effects to ensure we're not just getting lucky with identical strings
		effect := AnimationEntranceAppear
		if i%2 == 0 {
			effect = AnimationEntranceFade
		}
		slide = slide.AddAnimation(NewAnimation(1, effect))
	}

	// Generate XML multiple times and verify they are identical
	data1, err := CreateWithSlides("Run 1", []SlideContent{slide})
	if err != nil {
		t.Fatalf("Run 1 CreateWithSlides error: %v", err)
	}

	data2, err := CreateWithSlides("Run 2", []SlideContent{slide})
	if err != nil {
		t.Fatalf("Run 2 CreateWithSlides error: %v", err)
	}

	xml1 := readZipFileForAnim(t, getZipReader(t, data1), "ppt/slides/slide1.xml")
	xml2 := readZipFileForAnim(t, getZipReader(t, data2), "ppt/slides/slide1.xml")

	if xml1 != xml2 {
		t.Errorf("Animation XML generation is not deterministic")
	}

	// Verify order is preserved (first animation should be fade, second appear, etc.)
	// This is weak parsing but sufficient for order check if we look for preset IDs
	// Fade = 10, Appear = 1
	// We expect 10, 1, 10, 1 ...

	// We can check if the XML contains the sequence of presetIDs in correct order
	// Or just check that the parallel execution didn't scramble things wildly.
	// Since we compare exact string equality above, we know it's *stable*.
	// Use manual inspection or substring index checks for *correctness* of order.

	// Check first and last to ensure no dropped items
	// The first cTn should have id="3" (root is 1, mainSeq is 2, first par is 3) or similar
	// Actually, the indices in `xml()` are `i*2+3` => 0->3, 1->5, ...
	// So first animation is ID 3, last is 59*2+3 = 121.
	if !strings.Contains(xml1, `id="3"`) {
		t.Error("Missing first animation ID 3")
	}
	if !strings.Contains(xml1, `id="121"`) {
		t.Error("Missing last animation ID 121")
	}
}

func getZipReader(t *testing.T, data []byte) *zip.Reader {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip reader error: %v", err)
	}
	return zr
}

// TestAnimationValidationConflicts verifies that invalid animation sequences are rejected.
func TestAnimationValidationConflicts(t *testing.T) {
	// 1. First animation cannot be WithPrevious
	slide := NewSlide("Conflict 1").
		AddShape(NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(NewAnimation(1, AnimationEntranceFade).WithTrigger(AnimationWithPrevious))

	_, err := CreateWithSlides("Conflict 1", []SlideContent{slide})
	if err == nil {
		t.Error("Expected error for first animation WithPrevious, got nil")
	} else if !strings.Contains(err.Error(), "trigger cannot be") {
		t.Errorf("Unexpected error message: %v", err)
	}

	// 2. First animation cannot be AfterPrevious
	slide2 := NewSlide("Conflict 2").
		AddShape(NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(NewAnimation(1, AnimationEntranceFade).WithTrigger(AnimationAfterPrevious))

	_, err = CreateWithSlides("Conflict 2", []SlideContent{slide2})
	if err == nil {
		t.Error("Expected error for first animation AfterPrevious, got nil")
	}
}

// TestAnimationInterface extensibility verification
type CustomAnim struct {
	Index int
}

func (c CustomAnim) ToAnimation() Animation {
	return NewAnimation(c.Index, AnimationEmphasisPulse)
}

func TestAnimationInterface(t *testing.T) {
	custom := CustomAnim{Index: 1}
	slide := NewSlide("Interface Test").
		AddShape(NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(custom) // Should accept interface

	if len(slide.Animations) != 1 {
		t.Fatalf("Expected 1 animation, got %d", len(slide.Animations))
	}
	if slide.Animations[0].Effect != AnimationEmphasisPulse {
		t.Errorf("Expected Pulse effect from custom animation, got %s", slide.Animations[0].Effect)
	}
}
