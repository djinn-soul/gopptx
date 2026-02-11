package animations_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

// TestAnimationXMLDeterminism verifies that concurrent animation processing produces deterministic output order.
func TestAnimationXMLDeterminism(t *testing.T) {
	slide := pptx.NewSlide("Deterministic Anim").
		AddShape(pptx.NewShape("rect", 0, 0, 100, 100))

	for i := 0; i < 60; i++ {
		effect := pptx.AnimationEntranceAppear
		if i%2 == 0 {
			effect = pptx.AnimationEntranceFade
		}
		slide = slide.AddAnimation(pptx.NewAnimation(1, effect))
	}

	data1, err := pptx.CreateWithSlides("Run 1", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("Run 1 CreateWithSlides error: %v", err)
	}

	data2, err := pptx.CreateWithSlides("Run 2", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("Run 2 CreateWithSlides error: %v", err)
	}

	xml1 := readAnimZipFile(t, data1, "ppt/slides/slide1.xml")
	xml2 := readAnimZipFile(t, data2, "ppt/slides/slide1.xml")

	if xml1 != xml2 {
		t.Errorf("Animation XML generation is not deterministic")
	}

	if !strings.Contains(xml1, `id="3"`) {
		t.Error("Missing first animation ID 3")
	}
	if !strings.Contains(xml1, `id="121"`) {
		t.Error("Missing last animation ID 121")
	}
}

func readAnimZipFile(t *testing.T, data []byte, name string) string {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip reader error: %v", err)
	}
	return testutil.ReadZipFile(t, zr, name)
}

// TestAnimationValidationConflicts verifies that invalid animation sequences are rejected.
func TestAnimationValidationConflicts(t *testing.T) {
	slide := pptx.NewSlide("Conflict 1").
		AddShape(pptx.NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationEntranceFade).WithTrigger(pptx.AnimationWithPrevious))

	_, err := pptx.CreateWithSlides("Conflict 1", []pptx.SlideContent{slide})
	if err == nil {
		t.Error("Expected error for first animation WithPrevious, got nil")
	} else if !strings.Contains(err.Error(), "trigger cannot be") {
		t.Errorf("Unexpected error message: %v", err)
	}

	slide2 := pptx.NewSlide("Conflict 2").
		AddShape(pptx.NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationEntranceFade).WithTrigger(pptx.AnimationAfterPrevious))

	_, err = pptx.CreateWithSlides("Conflict 2", []pptx.SlideContent{slide2})
	if err == nil {
		t.Error("Expected error for first animation AfterPrevious, got nil")
	}
}

// customAnim for testing AnimationDefinition interface.
type customAnim struct {
	Index int
}

func (c customAnim) ToAnimation() pptx.Animation {
	return pptx.NewAnimation(c.Index, pptx.AnimationEmphasisPulse)
}

func TestAnimationInterface(t *testing.T) {
	custom := customAnim{Index: 1}
	slide := pptx.NewSlide("Interface Test").
		AddShape(pptx.NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(custom)

	if len(slide.Animations) != 1 {
		t.Fatalf("Expected 1 animation, got %d", len(slide.Animations))
	}
	if slide.Animations[0].Effect != pptx.AnimationEmphasisPulse {
		t.Errorf("Expected Pulse effect from custom animation, got %s", slide.Animations[0].Effect)
	}
}
