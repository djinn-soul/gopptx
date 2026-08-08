package export

import (
	"math"
	"testing"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func cachedNodeShape(text string) smartart.DrawingShape {
	shape := smartart.DrawingShape{
		X: 91440, Y: 182880, CX: 1828800, CY: 914400,
		PresetGeom:  "roundRect",
		Adjustments: map[string]int{"adj": 10000},
		Fill:        smartart.ColorRef{SRGB: "4472C4"},
		Line:        smartart.ColorRef{SRGB: "FFFFFF"},
		// 0.75pt, the width PowerPoint gives a SmartArt node outline.
		LineWidthEMU: 9525,
		Anchor:       "ctr",
	}
	if text != "" {
		shape.Paragraphs = []smartart.DrawingParagraph{{
			Align: "ctr",
			Runs:  []smartart.DrawingRun{{Text: text, SizePt: 18, Color: smartart.ColorRef{SRGB: "FFFFFF"}}},
		}}
	}
	return shape
}

func cachedDiagram(shapes []smartart.DrawingShape, nodeText string) smartart.SmartArt {
	diagram := smartart.NewSmartArt(smartart.AccentProcess).
		Position(styling.Inches(1), styling.Inches(2)).
		Size(styling.Inches(6), styling.Inches(3))
	if nodeText != "" {
		diagram = diagram.AddNode(smartart.NewNode(nodeText))
	}
	return diagram.WithDrawing(shapes)
}

func TestSmartArtCachedShapeOffsetsByFrameOrigin(t *testing.T) {
	diagram := cachedDiagram([]smartart.DrawingShape{cachedNodeShape("Step")}, "Step")

	shape, ok := smartArtCachedShape(diagram, diagram.Drawing[0])
	if !ok {
		t.Fatal("cached shape was rejected")
	}
	wantX := int64(diagram.X) + diagram.Drawing[0].X
	wantY := int64(diagram.Y) + diagram.Drawing[0].Y
	if int64(shape.X) != wantX || int64(shape.Y) != wantY {
		t.Errorf("origin = %d,%d, want %d,%d", int64(shape.X), int64(shape.Y), wantX, wantY)
	}
	if int64(shape.CX) != 1828800 || int64(shape.CY) != 914400 {
		t.Errorf("size = %dx%d, want 1828800x914400", int64(shape.CX), int64(shape.CY))
	}
	if shape.Type != "roundRect" {
		t.Errorf("Type = %q, want roundRect", shape.Type)
	}
	if shape.Fill == nil || shape.Fill.Color != "4472C4" {
		t.Errorf("Fill = %+v, want 4472C4", shape.Fill)
	}
	if shape.Line == nil || int64(shape.Line.Width) != 9525 {
		t.Errorf("Line = %+v, want width 9525", shape.Line)
	}
	if shape.TextFrame == nil || shape.TextFrame.Anchor != shapes.TextAnchorMiddle {
		t.Errorf("anchor = %+v, want ctr", shape.TextFrame)
	}
	if len(shape.TextParagraphs) != 1 || shape.TextParagraphs[0].Style.Align != elements.TextAlignCenter {
		t.Errorf("paragraphs = %+v, want one centred", shape.TextParagraphs)
	}
	if shape.Text != "Step" {
		t.Errorf("Text = %q, want Step", shape.Text)
	}
}

func TestSmartArtCachedShapeSkipsEmptySpacers(t *testing.T) {
	diagram := cachedDiagram(nil, "")
	spacer := smartart.DrawingShape{X: 0, Y: 0, CX: 100, CY: 100}
	if _, ok := smartArtCachedShape(diagram, spacer); ok {
		t.Error("a shape with neither geometry nor text was drawn")
	}
	zeroSized := smartart.DrawingShape{PresetGeom: "rect"}
	if _, ok := smartArtCachedShape(diagram, zeroSized); ok {
		t.Error("a zero-sized shape was drawn")
	}
}

func TestSmartArtCacheIsStaleWhenItLostItsCaptions(t *testing.T) {
	// A writer that outgrew its template strips the cached captions; the
	// heuristic layouts know the real ones, so they have to win.
	stale := cachedDiagram([]smartart.DrawingShape{cachedNodeShape("")}, "Real caption")
	if !smartArtCacheIsStale(stale) {
		t.Error("a caption-less cache under a diagram with text was accepted")
	}

	fresh := cachedDiagram([]smartart.DrawingShape{cachedNodeShape("Real caption")}, "Real caption")
	if smartArtCacheIsStale(fresh) {
		t.Error("a cache carrying the captions was rejected")
	}

	// A diagram whose nodes have no text of their own cannot contradict the
	// cache, so a text-free cache is fine there.
	decorative := cachedDiagram([]smartart.DrawingShape{cachedNodeShape("")}, "")
	if smartArtCacheIsStale(decorative) {
		t.Error("a text-free diagram rejected its own cache")
	}
}

func TestRenderPDFSmartArtFromCacheReportsWhetherItDrew(t *testing.T) {
	pdf := newTestPDF(t)
	// Drawing text needs a font on the document, as it does in a real export.
	if err := configureNativePDFFont(pdf, PDFOptions{}); err != nil {
		t.Fatalf("configureNativePDFFont: %v", err)
	}

	if renderPDFSmartArtFromCache(pdf, cachedDiagram(nil, "Step")) {
		t.Error("reported drawing with no cache")
	}
	if !renderPDFSmartArtFromCache(pdf, cachedDiagram([]smartart.DrawingShape{cachedNodeShape("Step")}, "Step")) {
		t.Error("did not draw a usable cache")
	}
	if renderPDFSmartArtFromCache(pdf, cachedDiagram([]smartart.DrawingShape{cachedNodeShape("")}, "Step")) {
		t.Error("drew a stale cache")
	}
}

func TestCachedShapeAppliesStatedTextInsets(t *testing.T) {
	cached := cachedNodeShape("Step")
	tight := int64(45720) // 0.05in
	cached.InsetLeft = &tight
	cached.InsetTop = &tight

	shape, ok := smartArtCachedShape(cachedDiagram([]smartart.DrawingShape{cached}, "Step"), cached)
	if !ok {
		t.Fatal("the cached shape was rejected")
	}
	if shape.TextFrame == nil {
		t.Fatal("no text frame was built for a body that states insets")
	}
	if int64(shape.TextFrame.MarginLeft) != tight {
		t.Errorf("MarginLeft = %d, want the stated %d", int64(shape.TextFrame.MarginLeft), tight)
	}
	if int64(shape.TextFrame.MarginTop) != tight {
		t.Errorf("MarginTop = %d, want the stated %d", int64(shape.TextFrame.MarginTop), tight)
	}
	// The sides the body left unstated keep the OOXML defaults.
	if int64(shape.TextFrame.MarginRight) != ooxmlDefaultInsetLREMU {
		t.Errorf("MarginRight = %d, want the default %d",
			int64(shape.TextFrame.MarginRight), ooxmlDefaultInsetLREMU)
	}
	if int64(shape.TextFrame.MarginBottom) != ooxmlDefaultInsetTBEMU {
		t.Errorf("MarginBottom = %d, want the default %d",
			int64(shape.TextFrame.MarginBottom), ooxmlDefaultInsetTBEMU)
	}
}

func TestCachedShapeWithoutInsetsKeepsTheDefaults(t *testing.T) {
	cached := cachedNodeShape("Step")
	shape, ok := smartArtCachedShape(cachedDiagram([]smartart.DrawingShape{cached}, "Step"), cached)
	if !ok {
		t.Fatal("the cached shape was rejected")
	}
	if shape.TextFrame == nil {
		t.Fatal("the anchor should still have built a text frame")
	}
	if int64(shape.TextFrame.MarginLeft) != ooxmlDefaultInsetLREMU {
		t.Errorf("MarginLeft = %d, want the default %d",
			int64(shape.TextFrame.MarginLeft), ooxmlDefaultInsetLREMU)
	}
}

func TestCachedFillCarriesTheStatedOpacity(t *testing.T) {
	// 40000 of 100000 is 40% opaque, so 60% transparent.
	fill := smartArtCachedFill("4472C4", 40000)
	if fill.Transparency == nil {
		t.Fatal("a part-transparent fill came through opaque")
	}
	if math.Abs(*fill.Transparency-0.6) > 0.0001 {
		t.Errorf("Transparency = %v, want 0.6", *fill.Transparency)
	}
}

func TestCachedFillTreatsOmittedAlphaAsOpaque(t *testing.T) {
	// OOXML omits a:alpha for an opaque colour, so the parser leaves zero.
	if fill := smartArtCachedFill("4472C4", 0); fill.Transparency != nil {
		t.Errorf("Transparency = %v, want nil for an omitted alpha", *fill.Transparency)
	}
	if fill := smartArtCachedFill("4472C4", 100000); fill.Transparency != nil {
		t.Errorf("Transparency = %v, want nil for a fully opaque alpha", *fill.Transparency)
	}
}

func TestFillAlphaIsAppliedAndDefersToSoftEdges(t *testing.T) {
	pdf := newTestPDF(t)

	translucent := shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 100).
		WithFill(shapes.NewShapeFill("4472C4").WithTransparency(0.5))
	if !applyPDFShapeFillAlpha(pdf, translucent, false) {
		t.Error("a translucent fill did not set any transparency")
	}
	pdf.ClearTransparency()

	if applyPDFShapeFillAlpha(pdf, translucent, true) {
		t.Error("the fill alpha overwrote the transparency soft edges had already set")
	}

	opaque := shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 100).
		WithFill(shapes.NewShapeFill("4472C4"))
	if applyPDFShapeFillAlpha(pdf, opaque, false) {
		t.Error("an opaque fill set a transparency")
	}
}

// TestEffectsBlendModeIsAcceptedByGopdf guards the constant that silently
// disabled every transparency in this package: gopdf rejects an unknown blend
// mode, and the callers swallow that error.
func TestEffectsBlendModeIsAcceptedByGopdf(t *testing.T) {
	if _, err := gopdf.NewTransparency(0.5, shapeEffectsBlendMode); err != nil {
		t.Fatalf("gopdf rejects shapeEffectsBlendMode %q: %v", shapeEffectsBlendMode, err)
	}
}

func TestSoftEdgesActuallyApplyTransparency(t *testing.T) {
	pdf := newTestPDF(t)
	shape := shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 100).WithSoftEdges(true)
	if !applyPDFShapeSoftEdges(pdf, shape) {
		t.Error("soft edges reported no transparency set")
	}
}
