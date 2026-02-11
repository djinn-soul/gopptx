package pptx_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestSlideAddAnimation(t *testing.T) {
	slide := pptx.NewSlide("Animation Test").
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationEntranceFade))

	if len(slide.Animations) != 1 {
		t.Fatalf("expected 1 animation, got %d", len(slide.Animations))
	}
	if slide.Animations[0].Effect != pptx.AnimationEntranceFade {
		t.Errorf("expected fade effect, got %s", slide.Animations[0].Effect)
	}
}

func TestSlideAnimationXML(t *testing.T) {
	slide := pptx.NewSlide("Anim XML").
		AddShape(pptx.NewShape("rect", 100, 100, 200, 200)).
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationEntranceFade))

	data, err := pptx.CreateWithSlides("Anim Demo", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	xml := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	if !strings.Contains(xml, "<p:timing>") {
		t.Error("missing <p:timing> block")
	}
	if !strings.Contains(xml, `presetID="10"`) {
		t.Errorf("missing presetID 10 for fade animation in XML: %s", xml)
	}
	if !strings.Contains(xml, "presetClass=\"entr\"") {
		t.Error("missing presetClass entr for entrance animation")
	}
	if !strings.Contains(xml, "spid=\"3\"") {
		t.Errorf("expected shape ID 3 in animation target, but not found in XML: %s", xml)
	}
}

func TestAnimationSequence(t *testing.T) {
	slide := pptx.NewSlide("Sequence Test").
		AddShape(pptx.NewShape("rect", 0, 0, 100, 100)).
		AddShape(pptx.NewShape("ellipse", 100, 100, 100, 100)).
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationEntranceFade)).
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationExitFadeOut).WithTrigger(pptx.AnimationOnClick)).
		AddAnimation(pptx.NewAnimation(2, pptx.AnimationEmphasisPulse).WithTrigger(pptx.AnimationAfterPrevious))

	data, err := pptx.CreateWithSlides("Sequence Demo", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	xml := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	if strings.Count(xml, "<p:par>") < 3 {
		t.Errorf("expected at least 3 <p:par> animation blocks, got %d", strings.Count(xml, "<p:par>"))
	}
	if !strings.Contains(xml, "spid=\"3\"") || !strings.Contains(xml, "spid=\"4\"") {
		t.Error("missing target shape IDs 3 or 4 in animations")
	}
}

func TestTransitionAndAnimationCoexistence(t *testing.T) {
	slide := pptx.NewSlide("Coexistence").
		WithTransition(pptx.TransitionCover).
		AddShape(pptx.NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationEntranceFade))

	data, err := pptx.CreateWithSlides("Coexistence Demo", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	xml := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	if !strings.Contains(xml, "<p:transition") {
		t.Error("missing <p:transition> block")
	}
	if !strings.Contains(xml, "<p:timing>") {
		t.Error("missing <p:timing> block")
	}

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
	slide := pptx.NewSlide("Two Column").
		WithTwoColumnLayout().
		AddBullet("Only left column").
		AddShape(pptx.NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationEntranceFade))

	data, err := pptx.CreateWithSlides("Two Column Anim", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	xml := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	if !strings.Contains(xml, `spid="4"`) {
		t.Fatalf("expected shape animation target spid=4 for two-column single-bullet layout")
	}
}

func TestCreateWithSlidesRejectsUnknownAnimationEffect(t *testing.T) {
	slide := pptx.NewSlide("Invalid Animation").
		AddShape(pptx.NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationEffect("unknown_effect")))

	_, err := pptx.CreateWithSlides("Invalid Anim", []pptx.SlideContent{slide})
	if err == nil {
		t.Fatalf("expected unsupported animation effect validation error")
	}
	if !strings.Contains(err.Error(), "unsupported animation effect") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnimationTriggersRenderOpenXMLNodeTypes(t *testing.T) {
	slide := pptx.NewSlide("Trigger Mapping").
		AddShape(pptx.NewShape("rect", 0, 0, 100, 100)).
		AddShape(pptx.NewShape("rect", 200, 0, 100, 100)).
		AddShape(pptx.NewShape("rect", 400, 0, 100, 100)).
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationEntranceFade).WithTrigger(pptx.AnimationOnClick)).
		AddAnimation(pptx.NewAnimation(2, pptx.AnimationEntranceFade).WithTrigger(pptx.AnimationWithPrevious)).
		AddAnimation(pptx.NewAnimation(3, pptx.AnimationEntranceFade).WithTrigger(pptx.AnimationAfterPrevious))

	data, err := pptx.CreateWithSlides("Trigger Mapping", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	xml := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

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
	slide := pptx.NewSlide("Invalid Trigger").
		AddShape(pptx.NewShape("rect", 0, 0, 100, 100)).
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationEntranceFade).WithTrigger(pptx.AnimationTrigger("bad_trigger")))

	_, err := pptx.CreateWithSlides("Invalid Trigger", []pptx.SlideContent{slide})
	if err == nil {
		t.Fatalf("expected unsupported animation trigger validation error")
	}
	if !strings.Contains(err.Error(), "unsupported animation trigger") {
		t.Fatalf("unexpected error: %v", err)
	}
}
