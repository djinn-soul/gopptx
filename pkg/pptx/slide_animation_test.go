package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestSlideAddAnimation(t *testing.T) {
	slide := NewSlide("Animation Test").
		AddAnimation(NewAnimation(1, AnimationEntranceFade))

	if len(slide.Animations) != 1 {
		t.Fatalf("expected 1 animation, got %d", len(slide.Animations))
	}
	if slide.Animations[0].Effect != AnimationEntranceFade {
		t.Errorf("expected fade effect, got %s", slide.Animations[0].Effect)
	}
}

func TestSlideAnimationXML(t *testing.T) {
	// Create a slide with a shape and an animation
	slide := NewSlide("Anim XML").
		AddShape(NewShape("rect", 100, 100, 200, 200)).
		AddAnimation(NewAnimation(1, AnimationEntranceFade))

	data, err := CreateWithSlides("Anim Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	xml := readZipFileForAnim(t, zr, "ppt/slides/slide1.xml")

	// Check for <p:timing> block
	if !strings.Contains(xml, "<p:timing>") {
		t.Error("missing <p:timing> block")
	}
	if !strings.Contains(xml, `presetID="10"`) { // 10 is fade
		t.Errorf("missing presetID 10 for fade animation in XML: %s", xml)
	}
	if !strings.Contains(xml, "presetClass=\"entr\"") {
		t.Error("missing presetClass entr for entrance animation")
	}
	// Shape ID for the first custom shape should be animNextID (which is 3 for titleAndContent)
	if !strings.Contains(xml, "spid=\"3\"") {
		t.Errorf("expected shape ID 3 in animation target, but not found in XML: %s", xml)
	}
}

func readZipFileForAnim(t *testing.T, zr *zip.Reader, name string) string {
	t.Helper()
	f, err := zr.Open(name)
	if err != nil {
		t.Fatalf("failed to open zip file %s: %v", name, err)
	}
	defer func() { _ = f.Close() }()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		t.Fatalf("failed to read zip file %s: %v", name, err)
	}
	return buf.String()
}

func TestAnimationSequence(t *testing.T) {
	slide := NewSlide("Sequence Test").
		AddShape(NewShape("rect", 0, 0, 100, 100)).        // ID 3
		AddShape(NewShape("ellipse", 100, 100, 100, 100)). // ID 4
		AddAnimation(NewAnimation(1, AnimationEntranceFade)).
		AddAnimation(NewAnimation(1, AnimationExitFadeOut).WithTrigger(AnimationOnClick)).
		AddAnimation(NewAnimation(2, AnimationEmphasisPulse).WithTrigger(AnimationAfterPrevious))

	data, err := CreateWithSlides("Sequence Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	xml := readZipFileForAnim(t, zr, "ppt/slides/slide1.xml")

	// Verify all 3 animations are present
	if strings.Count(xml, "<p:par>") < 3 {
		t.Errorf("expected at least 3 <p:par> animation blocks, got %d", strings.Count(xml, "<p:par>"))
	}
	if !strings.Contains(xml, "spid=\"3\"") || !strings.Contains(xml, "spid=\"4\"") {
		t.Error("missing target shape IDs 3 or 4 in animations")
	}
}

// TestTransitionAndAnimationCoexistence verifies that adding a transition does not break animation XML structure.
func TestTransitionAndAnimationCoexistence(t *testing.T) {
	slide := NewSlide("Coexistence").
		WithTransition(TransitionCover).
		AddShape(NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(NewAnimation(1, AnimationEntranceFade))

	data, err := CreateWithSlides("Coexistence Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	xml := readZipFileForAnim(t, zr, "ppt/slides/slide1.xml")

	if !strings.Contains(xml, "<p:transition") {
		t.Error("missing <p:transition> block")
	}
	if !strings.Contains(xml, "<p:timing>") {
		t.Error("missing <p:timing> block")
	}

	// Check order: transition should be before timing (based on previous fixes, but let's verify)
	transIdx := strings.Index(xml, "<p:transition")
	timingIdx := strings.Index(xml, "<p:timing")
	if transIdx == -1 || timingIdx == -1 {
		t.Fatal("missing transition or timing")
	}
	if transIdx > timingIdx {
		t.Errorf("transition block appears after timing block, which violates schema order (should be before)")
	}
}

func TestTwoColumnSingleBulletAnimationTargetsCorrectShape(t *testing.T) {
	slide := NewSlide("Two Column").
		WithTwoColumnLayout().
		AddBullet("Only left column").
		AddShape(NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(NewAnimation(1, AnimationEntranceFade))

	data, err := CreateWithSlides("Two Column Anim", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	xml := readZipFileForAnim(t, zr, "ppt/slides/slide1.xml")

	if !strings.Contains(xml, `spid="4"`) {
		t.Fatalf("expected shape animation target spid=4 for two-column single-bullet layout")
	}
}

func TestCreateWithSlidesRejectsUnknownAnimationEffect(t *testing.T) {
	slide := NewSlide("Invalid Animation").
		AddShape(NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(NewAnimation(1, AnimationEffect("unknown_effect")))

	_, err := CreateWithSlides("Invalid Anim", []SlideContent{slide})
	if err == nil {
		t.Fatalf("expected unsupported animation effect validation error")
	}
	if !strings.Contains(err.Error(), "unsupported animation effect") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnimationTriggersRenderOpenXMLNodeTypes(t *testing.T) {
	slide := NewSlide("Trigger Mapping").
		AddShape(NewShape("rect", 0, 0, 100, 100)).
		AddShape(NewShape("rect", 200, 0, 100, 100)).
		AddShape(NewShape("rect", 400, 0, 100, 100)).
		AddAnimation(NewAnimation(1, AnimationEntranceFade).WithTrigger(AnimationOnClick)).
		AddAnimation(NewAnimation(2, AnimationEntranceFade).WithTrigger(AnimationWithPrevious)).
		AddAnimation(NewAnimation(3, AnimationEntranceFade).WithTrigger(AnimationAfterPrevious))

	data, err := CreateWithSlides("Trigger Mapping", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	xml := readZipFileForAnim(t, zr, "ppt/slides/slide1.xml")

	if !strings.Contains(xml, `nodeType="clickEffect"`) {
		t.Fatalf("expected clickEffect node type in animation XML")
	}
	if !strings.Contains(xml, `nodeType="withEffect"`) {
		t.Fatalf("expected withEffect node type in animation XML")
	}
	if !strings.Contains(xml, `nodeType="afterEffect"`) {
		t.Fatalf("expected afterEffect node type in animation XML")
	}
}

func TestCreateWithSlidesRejectsUnknownAnimationTrigger(t *testing.T) {
	slide := NewSlide("Invalid Trigger").
		AddShape(NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(NewAnimation(1, AnimationEntranceFade).WithTrigger(AnimationTrigger("bad_trigger")))

	_, err := CreateWithSlides("Invalid Trigger", []SlideContent{slide})
	if err == nil {
		t.Fatalf("expected unsupported animation trigger validation error")
	}
	if !strings.Contains(err.Error(), "unsupported animation trigger") {
		t.Fatalf("unexpected error: %v", err)
	}
}
