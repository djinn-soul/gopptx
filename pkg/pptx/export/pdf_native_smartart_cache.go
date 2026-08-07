package export

import (
	"math"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

// renderPDFSmartArtFromCache paints the layout PowerPoint computed for the
// diagram and cached in the deck's drawing part.
//
// It reports whether it drew anything. There is nothing to draw for a diagram
// built in memory, which has never been through PowerPoint's layout engine; the
// caller falls back to laying the diagram out itself.
//
// Every cached shape is turned into an ordinary auto shape and handed to the
// normal shape renderer, so a diagram gets the same preset geometry, fills,
// outlines, rotation and rich text as any other shape on the slide — rather than
// the per-layout approximations this renderer has to fall back on.
func renderPDFSmartArtFromCache(pdf *gopdf.GoPdf, diagram smartart.SmartArt) bool {
	if len(diagram.Drawing) == 0 || smartArtCacheIsStale(diagram) {
		return false
	}
	drawn := false
	for _, cached := range diagram.Drawing {
		shape, ok := smartArtCachedShape(diagram, cached)
		if !ok {
			continue
		}
		renderPDFShape(pdf, shape)
		drawn = true
	}
	return drawn
}

// smartArtCacheIsStale reports that the cache cannot be the picture of this
// diagram, so drawing it would show the wrong thing.
//
// A writer that outgrew its template leaves the cache in place with its captions
// stripped, because PowerPoint recomputes the layout and never reads it. A cache
// with no text at all under a diagram that has some is exactly that case: the
// heuristic layouts draw the real captions, so they win.
func smartArtCacheIsStale(diagram smartart.SmartArt) bool {
	if !smartArtHasNodeText(diagram.Nodes) {
		return false
	}
	for _, cached := range diagram.Drawing {
		if cached.HasText() {
			return false
		}
	}
	return true
}

func smartArtHasNodeText(nodes []smartart.Node) bool {
	for _, node := range nodes {
		if strings.TrimSpace(node.Text) != "" || smartArtHasNodeText(node.Children) {
			return true
		}
	}
	return false
}

// smartArtCachedShape places one cached shape on the slide. The cache states its
// geometry relative to the diagram frame, so the frame's own origin is added.
func smartArtCachedShape(diagram smartart.SmartArt, cached smartart.DrawingShape) (shapes.Shape, bool) {
	if cached.CX <= 0 || cached.CY <= 0 {
		return shapes.Shape{}, false
	}
	// A shape with neither geometry nor text draws nothing; skipping it keeps
	// the invisible spacer shapes a layout emits out of the output.
	if cached.PresetGeom == "" && !cached.HasText() {
		return shapes.Shape{}, false
	}

	shape := shapes.NewShape(
		smartArtCachedGeometry(cached),
		styling.Emu(int64(diagram.X)+cached.X),
		styling.Emu(int64(diagram.Y)+cached.Y),
		styling.Emu(cached.CX),
		styling.Emu(cached.CY),
	)
	if hex := cached.Fill.SRGB; hex != "" && cached.PresetGeom != "" {
		shape = shape.WithFill(shapes.NewShapeFill(hex))
	}
	if hex := cached.Line.SRGB; hex != "" && cached.PresetGeom != "" && cached.LineWidthEMU > 0 {
		shape = shape.WithLine(shapes.NewShapeLine(hex, styling.Emu(cached.LineWidthEMU)))
	}
	shape.Adjustments = smartArtCachedAdjustments(cached)
	// The cache states a mirrored shape as a flip rather than as mirrored
	// geometry, so an arrow or chevron drew the wrong way round without this.
	shape.FlipH = cached.FlipH
	shape.FlipV = cached.FlipV
	if rotation := smartArtCachedRotation(cached); rotation != 0 {
		shape.RotationDeg = &rotation
	}
	shape = shape.WithVerticalAnchor(smartArtCachedAnchor(cached.Anchor))
	shape = smartArtCachedInsets(shape, cached)
	shape.TextParagraphs = smartArtCachedParagraphs(cached)
	shape.Text = smartArtCachedPlainText(cached)
	return shape, true
}

// smartArtCachedGeometry is the preset the cache names. A shape that names none
// carries only text, so it is drawn as an unfilled rectangle: the text renderer
// needs a box to lay the text out in.
func smartArtCachedGeometry(cached smartart.DrawingShape) string {
	if cached.PresetGeom == "" {
		return shapesShapeRectangle
	}
	return cached.PresetGeom
}

func smartArtCachedAdjustments(cached smartart.DrawingShape) []shapes.ShapeAdjustment {
	if len(cached.Adjustments) == 0 {
		return nil
	}
	out := make([]shapes.ShapeAdjustment, 0, len(cached.Adjustments))
	for name, value := range cached.Adjustments {
		out = append(out, shapes.ShapeAdjustment{Name: name, Formula: "val " + strconv.Itoa(value)})
	}
	return out
}

// smartArtCachedRotation rounds the cache's rotation to whole degrees, which is
// what Shape carries.
func smartArtCachedRotation(cached smartart.DrawingShape) int {
	return int(math.Round(cached.RotationDeg))
}

// smartArtCachedInsets applies the text insets the cached body states, falling
// back to the OOXML defaults per side. The cache was read for these and then
// thrown away, so a node whose layout tightened its insets drew its caption at
// the wrong offset.
//
// They are set even when the body states none, because the shape already has a
// text frame by this point — setting the anchor creates one — and a fresh frame
// defaults to 0.05in on every side, where OOXML puts 0.1in on the left and
// right. Left alone, every cached caption sat too close to its box's sides.
func smartArtCachedInsets(shape shapes.Shape, cached smartart.DrawingShape) shapes.Shape {
	return shape.WithTextMargins(
		styling.Emu(insetOrDefault(cached.InsetLeft, ooxmlDefaultInsetLREMU)),
		styling.Emu(insetOrDefault(cached.InsetTop, ooxmlDefaultInsetTBEMU)),
		styling.Emu(insetOrDefault(cached.InsetRight, ooxmlDefaultInsetLREMU)),
		styling.Emu(insetOrDefault(cached.InsetBottom, ooxmlDefaultInsetTBEMU)),
	)
}

func insetOrDefault(stated *int64, fallback int64) int64 {
	if stated == nil {
		return fallback
	}
	return *stated
}

func smartArtCachedAnchor(anchor string) shapes.TextFrameAnchor {
	switch strings.ToLower(strings.TrimSpace(anchor)) {
	case "t":
		return shapes.TextAnchorTop
	case "b":
		return shapes.TextAnchorBottom
	default:
		// SmartArt centres its captions unless the layout says otherwise, and
		// so does PowerPoint's own default body anchor for these shapes.
		return shapes.TextAnchorMiddle
	}
}

func smartArtCachedParagraphs(cached smartart.DrawingShape) []text.Paragraph {
	if len(cached.Paragraphs) == 0 {
		return nil
	}
	out := make([]text.Paragraph, 0, len(cached.Paragraphs))
	for _, paragraph := range cached.Paragraphs {
		runs := make([]text.Run, 0, len(paragraph.Runs))
		for _, run := range paragraph.Runs {
			runs = append(runs, text.Run{
				Text:   run.Text,
				Bold:   run.Bold,
				Italic: run.Italic,
				Color:  run.Color.SRGB,
				Font:   run.Font,
				SizePt: int(math.Round(run.SizePt)),
			})
		}
		out = append(out, text.Paragraph{
			Runs:  runs,
			Style: text.ParagraphStyle{Align: smartArtCachedAlign(paragraph.Align)},
		})
	}
	return out
}

// smartArtCachedAlign maps the OOXML algn value to the renderer's own. SmartArt
// captions are centred unless the cache says otherwise.
func smartArtCachedAlign(algn string) string {
	switch strings.ToLower(strings.TrimSpace(algn)) {
	case "l":
		return elements.TextAlignLeft
	case "r":
		return elements.TextAlignRight
	default:
		// "ctr", and anything the cache does not state.
		return elements.TextAlignCenter
	}
}

// smartArtCachedPlainText is the shape's text as one string. The shape renderer
// only draws a caption when Text is non-empty, and falls back to it for anything
// the rich-text path does not handle.
func smartArtCachedPlainText(cached smartart.DrawingShape) string {
	var b strings.Builder
	for i, paragraph := range cached.Paragraphs {
		if i > 0 {
			b.WriteByte('\n')
		}
		for _, run := range paragraph.Runs {
			b.WriteString(run.Text)
		}
	}
	return b.String()
}
