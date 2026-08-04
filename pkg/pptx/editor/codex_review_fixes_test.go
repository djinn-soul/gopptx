package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// A text_frame patch naming only the shrink amounts must keep them. They were
// dropped before, and a parsed shape almost always already has a frame, so the
// merge path is the one that matters.
func TestMergeTextFrameKeepsShrinkAmounts(t *testing.T) {
	wrap := true
	existing := &common.TextFrame{WordWrap: &wrap}
	scale, reduction := 62.5, 20.0
	patch := &common.TextFrame{FontScale: &scale, LineSpaceReduction: &reduction}

	merged := mergeTextFrame(existing, patch)

	if merged.FontScale == nil || *merged.FontScale != scale {
		t.Fatalf("font scale not merged: %+v", merged.FontScale)
	}
	if merged.LineSpaceReduction == nil || *merged.LineSpaceReduction != reduction {
		t.Fatalf("line space reduction not merged: %+v", merged.LineSpaceReduction)
	}
	if merged.WordWrap == nil || !*merged.WordWrap {
		t.Fatal("merge dropped the existing word wrap")
	}
}

// A shape with wider insets has less room for text, and the fitter has to
// measure against the frame's own margins rather than the bodyPr defaults.
func TestShapeTextInsetsUseFrameMargins(t *testing.T) {
	if x, y := shapeTextInsetsEmu(nil); x != 2*defaultTextInsetLREmu || y != 2*defaultTextInsetTBEmu {
		t.Fatalf("expected the default insets for a nil frame, got (%v,%v)", x, y)
	}

	left, right, top, bottom := 200000, 300000, 50000, 70000
	frame := &common.TextFrame{
		MarginLeft: &left, MarginRight: &right,
		MarginTop: &top, MarginBottom: &bottom,
	}
	x, y := shapeTextInsetsEmu(frame)
	if x != 500000 {
		t.Fatalf("expected 500000 EMU of horizontal inset, got %v", x)
	}
	if y != 120000 {
		t.Fatalf("expected 120000 EMU of vertical inset, got %v", y)
	}
}

// A macro-enabled template or slideshow keeps its macros just as a .pptm does.
func TestMacroEnabledExtensions(t *testing.T) {
	for _, ext := range []string{".pptm", ".potm", ".ppsm", ".PPTM"} {
		if !isMacroEnabledExtension(ext) {
			t.Errorf("%s should be accepted for a VBA-bearing deck", ext)
		}
	}
	for _, ext := range []string{".pptx", ".potx", ".ppsx", ""} {
		if isMacroEnabledExtension(ext) {
			t.Errorf("%s should not be accepted for a VBA-bearing deck", ext)
		}
	}
}

// A grandchild has to be reported in slide coordinates too. Applying the outer
// group's transform to the immediate child alone left it in the intermediate
// space of the inner group.
func TestNestedGroupTransformReachesGrandchildren(t *testing.T) {
	outer := parsedShape{
		ID: 1, Type: "group",
		X: 0, Y: 0, W: 1000, H: 1000,
		GroupChild: &common.GroupChildSpace{ExtentCx: 2000, ExtentCy: 2000},
		Shapes: []parsedShape{{
			ID: 2, Type: "group",
			X: 0, Y: 0, W: 2000, H: 2000,
			GroupChild: &common.GroupChildSpace{ExtentCx: 2000, ExtentCy: 2000},
			Shapes: []parsedShape{{
				ID: 3, Type: shapeTypeRect,
				X: 400, Y: 800, W: 200, H: 400,
			}},
		}},
	}

	shape := commonShapeFromParsed(outer)

	child := shape.Shapes[0]
	if child.W != 1000 || child.H != 1000 {
		t.Fatalf("inner group not halved by the outer transform: %dx%d", child.W, child.H)
	}
	grandchild := child.Shapes[0]
	// The inner group's space matches its own box, so it is an identity there;
	// only the outer group's halving applies.
	if grandchild.X != 200 || grandchild.Y != 400 {
		t.Fatalf("grandchild origin not mapped to the slide: (%d,%d)", grandchild.X, grandchild.Y)
	}
	if grandchild.W != 100 || grandchild.H != 200 {
		t.Fatalf("grandchild extent not mapped to the slide: %dx%d", grandchild.W, grandchild.H)
	}
}

// A placeholder shape carries no a:xfrm, so it inherits its box from the layout
// rather than being pinned wherever the default slide put it.
func TestSlidePlaceholderShapeXML(t *testing.T) {
	xml := slidePlaceholderShapeXML(
		7,
		common.PlaceholderInfo{Type: "body", Index: 2, Name: "Content 3"},
		"first\nsecond",
	)

	for _, want := range []string{
		`<p:ph type="body" idx="2"/>`,
		`<a:t>first</a:t>`,
		`<a:t>second</a:t>`,
		`name="Content 3"`,
	} {
		if !strings.Contains(xml, want) {
			t.Errorf("placeholder XML missing %s:\n%s", want, xml)
		}
	}
	if strings.Contains(xml, "<a:xfrm") {
		t.Error("a placeholder must inherit its geometry from the layout")
	}
}

// A title placeholder has idx 0, and writing idx="0" is not what PowerPoint
// does; the attribute is simply left off.
func TestSlidePlaceholderShapeXMLOmitsZeroIndex(t *testing.T) {
	xml := slidePlaceholderShapeXML(2, common.PlaceholderInfo{Type: "title"}, "")
	if !strings.Contains(xml, `<p:ph type="title"/>`) {
		t.Fatalf("expected a bare title placeholder, got:\n%s", xml)
	}
	if !strings.Contains(xml, "<a:p/>") {
		t.Errorf("an empty placeholder still needs a paragraph:\n%s", xml)
	}
}
