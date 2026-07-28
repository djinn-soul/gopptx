package editor

import (
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Upstream #1013: a placeholder that states no colour of its own reports None
// on every direct property, so the colour its layout or master defines is
// unreachable. The effective-style lookup walks that chain.
func TestGetEffectiveShapeStyleResolvesThroughLayoutMasterAndTheme(t *testing.T) {
	fixture := writeInheritanceFixture(t, "effective-style.pptx")

	ed, err := OpenPresentationEditor(fixture)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	// The title placeholder on the slide sets nothing but its text.
	shapes, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes: %v", err)
	}
	if len(shapes) != 1 {
		t.Fatalf("expected one shape, got %d", len(shapes))
	}
	if shapes[0].Fill != nil && shapes[0].Fill.Solid != nil {
		t.Fatalf("fixture shape should declare no fill of its own, got %+v", shapes[0].Fill)
	}

	style, err := ed.GetEffectiveShapeStyle(0, 2)
	if err != nil {
		t.Fatalf("effective style: %v", err)
	}

	// The layout's title placeholder carries the font colour, as a scheme
	// reference that only the theme can turn into an RGB.
	if style.FontColor == nil {
		t.Fatalf("expected an inherited font colour, got none")
	}
	if style.FontColor.Source != common.StyleSourceLayout {
		t.Fatalf("expected the font colour from the layout, got %q", style.FontColor.Source)
	}
	if style.FontColor.SchemeSlot != "accent1" {
		t.Fatalf("expected the accent1 slot, got %q", style.FontColor.SchemeSlot)
	}
	if style.FontColor.RGB != "4472C4" {
		t.Fatalf("expected accent1 resolved to 4472C4, got %q", style.FontColor.RGB)
	}

	// Size comes from the same layout placeholder.
	if style.FontSizePt == nil || style.FontSizePt.Value != 40 {
		t.Fatalf("expected 40pt from the layout, got %+v", style.FontSizePt)
	}

	// Position is not stated on the slide, so it is inherited too.
	if style.Position == nil || style.Position.Source != common.StyleSourceLayout {
		t.Fatalf("expected an inherited position, got %+v", style.Position)
	}
	if style.Position.W != 8000000 {
		t.Fatalf("expected the layout's width, got %d", style.Position.W)
	}

	// The layout defers its typeface to the theme's major font.
	if style.FontTypeface == nil || style.FontTypeface.Value != "Georgia" {
		t.Fatalf("expected +mj-lt resolved to Georgia, got %+v", style.FontTypeface)
	}
	if style.FontTypeface.Source != common.StyleSourceTheme {
		t.Fatalf("expected the typeface source to be the theme, got %q", style.FontTypeface.Source)
	}

	// Bold is stated only by the master's titleStyle.
	if style.Bold == nil || !style.Bold.Value {
		t.Fatalf("expected bold from the master text styles, got %+v", style.Bold)
	}
	if style.Bold.Source != common.StyleSourceMaster {
		t.Fatalf("expected the bold source to be the master, got %q", style.Bold.Source)
	}

	if style.LayoutPart != "ppt/slideLayouts/slideLayout1.xml" {
		t.Fatalf("unexpected layout part: %q", style.LayoutPart)
	}
	if style.MasterPart != "ppt/slideMasters/slideMaster1.xml" {
		t.Fatalf("unexpected master part: %q", style.MasterPart)
	}
}

// A value the shape states directly must win over everything it inherits.
func TestGetEffectiveShapeStyleShapeValuesWin(t *testing.T) {
	fixture := writeInheritanceFixture(t, "effective-style-override.pptx")

	ed, err := OpenPresentationEditor(fixture)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	slidePart := "ppt/slides/slide1.xml"
	content, ok := ed.parts.Get(slidePart)
	if !ok {
		t.Fatalf("slide part missing")
	}
	// Give the placeholder its own size, colour and extent.
	updated := replaceOnce(t, string(content),
		`<a:r><a:rPr lang="en-US"/><a:t>Title</a:t></a:r>`,
		`<a:r><a:rPr lang="en-US" sz="1800"><a:solidFill><a:srgbClr val="FF0000"/></a:solidFill>`+
			`</a:rPr><a:t>Title</a:t></a:r>`)
	updated = replaceOnce(t, updated,
		`<p:spPr/>`,
		`<p:spPr><a:xfrm><a:off x="10" y="20"/><a:ext cx="300" cy="400"/></a:xfrm></p:spPr>`)
	ed.parts.Set(slidePart, []byte(updated))

	style, err := ed.GetEffectiveShapeStyle(0, 2)
	if err != nil {
		t.Fatalf("effective style: %v", err)
	}
	if style.FontColor == nil || style.FontColor.RGB != "FF0000" ||
		style.FontColor.Source != common.StyleSourceShape {
		t.Fatalf("expected the shape's own colour to win, got %+v", style.FontColor)
	}
	if style.FontSizePt == nil || style.FontSizePt.Value != 18 {
		t.Fatalf("expected the shape's own 18pt to win, got %+v", style.FontSizePt)
	}
	if style.Position == nil || style.Position.Source != common.StyleSourceShape ||
		style.Position.W != 300 {
		t.Fatalf("expected the shape's own extent to win, got %+v", style.Position)
	}
}

// A shape that is not a placeholder inherits nothing from the layout chain.
func TestGetEffectiveShapeStyleIgnoresChainForNonPlaceholder(t *testing.T) {
	fixture := writeInheritanceFixture(t, "effective-style-plain.pptx")

	ed, err := OpenPresentationEditor(fixture)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	slidePart := "ppt/slides/slide1.xml"
	content, ok := ed.parts.Get(slidePart)
	if !ok {
		t.Fatalf("slide part missing")
	}
	updated := replaceOnce(t, string(content), `<p:ph type="title" idx="0"/>`, ``)
	ed.parts.Set(slidePart, []byte(updated))

	style, err := ed.GetEffectiveShapeStyle(0, 2)
	if err != nil {
		t.Fatalf("effective style: %v", err)
	}
	if style.FontColor != nil {
		t.Fatalf("a plain shape should inherit no colour, got %+v", style.FontColor)
	}
	if style.FontSizePt != nil {
		t.Fatalf("a plain shape should inherit no size, got %+v", style.FontSizePt)
	}
	// The theme default typeface still applies to any text.
	if style.FontTypeface == nil || style.FontTypeface.Value != "Verdana" {
		t.Fatalf("expected the theme's minor font, got %+v", style.FontTypeface)
	}
}
