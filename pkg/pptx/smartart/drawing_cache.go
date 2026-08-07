package smartart

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// A diagram's drawing part (ppt/diagrams/drawingN.xml) is the layout PowerPoint
// itself computed for the diagram: every shape with its final position, size,
// preset geometry, fill, outline and text. PowerPoint writes it as a cache so a
// reader that cannot run the DiagramML layout engine still has a picture to
// draw.
//
// Renderers that draw a diagram themselves — the native PDF exporter — get far
// closer to PowerPoint by painting this cache than by guessing a layout from the
// layout URI, which is what they had to do before. The types below are the
// subset of dsp:drawing a renderer needs.

// DrawingShape is one shape of a diagram's cached drawing.
//
// Its geometry is in EMU, relative to the diagram frame's own origin: a
// renderer adds the frame's X and Y to place it on the slide.
type DrawingShape struct {
	// ModelID ties the shape back to the point in the data model it draws.
	ModelID string

	X, Y, CX, CY int64

	// RotationDeg is the shape's rotation in degrees, clockwise.
	RotationDeg float64
	FlipH       bool
	FlipV       bool

	// PresetGeom is the OOXML preset name ("roundRect", "ellipse", …). Empty
	// means the shape stated no geometry and should not be drawn as a box.
	PresetGeom string
	// Adjustments are the preset's a:avLst guides, by name, in the preset's own
	// 1/100000 units.
	Adjustments map[string]int

	Fill ColorRef
	Line ColorRef
	// LineWidthEMU is the outline width. Zero means the shape stated none.
	LineWidthEMU int64

	// Anchor is the text body's vertical anchor: "t", "ctr" or "b".
	Anchor string
	// Insets are the text body's left, top, right and bottom insets in EMU.
	// A nil entry means the body stated none and the OOXML default applies.
	InsetLeft, InsetTop, InsetRight, InsetBottom *int64

	Paragraphs []DrawingParagraph
}

// HasText reports whether the shape carries any text to draw.
func (s DrawingShape) HasText() bool {
	for _, paragraph := range s.Paragraphs {
		for _, run := range paragraph.Runs {
			if run.Text != "" {
				return true
			}
		}
	}
	return false
}

// DrawingParagraph is one paragraph of a cached shape's text.
type DrawingParagraph struct {
	// Align is the OOXML algn value ("l", "ctr", "r", "just"), or empty.
	Align string
	Runs  []DrawingRun
}

// DrawingRun is one run of a cached paragraph, with the formatting the cache
// pinned on it. PowerPoint resolves SmartArt text size when it lays the diagram
// out, so SizePt here is the size it actually drew.
type DrawingRun struct {
	Text   string
	SizePt float64
	Bold   bool
	Italic bool
	Color  ColorRef
	Font   string
}

// ColorRef is a colour as the cache states it: either a literal RGB or a theme
// slot with the transforms the diagram's colour style applies to it. Resolving a
// scheme slot needs the deck's theme, which the cache does not carry, so that is
// left to the renderer.
type ColorRef struct {
	// SRGB is a literal colour as six hex digits, without a leading hash.
	SRGB string
	// Scheme is a theme slot name ("accent1", "lt1", …).
	Scheme string

	// Offsets are additive HSL adjustments, as OOXML states them: hue in
	// 1/60000 of a degree, saturation and luminance in 1/100000 (so 20000 is
	// +20 percentage points). The colourful SmartArt styles use these to spread
	// one accent across a diagram's nodes.
	HueOff, SatOff, LumOff int

	// Mods are multiplicative, in 1/100000. Zero means "not stated".
	LumMod, SatMod, Shade, Tint int

	// AlphaPct is the opacity in 1/100000; zero means fully opaque because the
	// cache omits a:alpha for an opaque colour.
	AlphaPct int
}

// IsSet reports whether the reference names a colour at all.
func (c ColorRef) IsSet() bool {
	return c.SRGB != "" || c.Scheme != ""
}

// ParseDrawingShapes reads the shapes out of a dsp:drawing part. A part that
// cannot be parsed, or that holds no shapes, yields no shapes and no error: the
// cache is an optimisation, and a renderer falls back to laying the diagram out
// itself.
func ParseDrawingShapes(data []byte) []DrawingShape {
	var drawing dspDrawing
	if err := xml.Unmarshal(data, &drawing); err != nil {
		return nil
	}
	shapes := make([]DrawingShape, 0, len(drawing.SpTree.Sp))
	for _, sp := range drawing.SpTree.Sp {
		shape, ok := convertDrawingShape(sp)
		if !ok {
			continue
		}
		shapes = append(shapes, shape)
	}
	if len(shapes) == 0 {
		return nil
	}
	return shapes
}

func convertDrawingShape(sp dspSp) (DrawingShape, bool) {
	// A shape with no transform has no place on the page. PowerPoint writes one
	// for every shape it laid out, so this only drops the group wrapper and any
	// shape the cache left incomplete.
	if sp.SpPr.Xfrm == nil {
		return DrawingShape{}, false
	}
	shape := DrawingShape{
		ModelID:      sp.ModelID,
		X:            sp.SpPr.Xfrm.Off.X,
		Y:            sp.SpPr.Xfrm.Off.Y,
		CX:           sp.SpPr.Xfrm.Ext.CX,
		CY:           sp.SpPr.Xfrm.Ext.CY,
		RotationDeg:  float64(sp.SpPr.Xfrm.Rot) / ooxmlAnglePerDegree,
		FlipH:        isXMLTrue(sp.SpPr.Xfrm.FlipH),
		FlipV:        isXMLTrue(sp.SpPr.Xfrm.FlipV),
		Fill:         colorRefFrom(sp.SpPr.SolidFill),
		Paragraphs:   convertDrawingParagraphs(sp.TxBody),
		LineWidthEMU: 0,
	}
	if sp.SpPr.PrstGeom != nil {
		shape.PresetGeom = sp.SpPr.PrstGeom.Prst
		shape.Adjustments = convertAdjustments(sp.SpPr.PrstGeom.Gds)
	}
	if sp.SpPr.NoFill != nil {
		shape.Fill = ColorRef{}
	}
	if sp.SpPr.Ln != nil && sp.SpPr.Ln.NoFill == nil {
		shape.Line = colorRefFrom(sp.SpPr.Ln.SolidFill)
		shape.LineWidthEMU = sp.SpPr.Ln.W
	}
	if sp.TxBody != nil {
		shape.Anchor = sp.TxBody.BodyPr.Anchor
		shape.InsetLeft = sp.TxBody.BodyPr.LeftInset
		shape.InsetTop = sp.TxBody.BodyPr.TopInset
		shape.InsetRight = sp.TxBody.BodyPr.RightInset
		shape.InsetBottom = sp.TxBody.BodyPr.BottomInset
	}
	return shape, true
}

// ooxmlAnglePerDegree is the unit a:xfrm states rotation in: 60000ths of a
// degree.
const ooxmlAnglePerDegree = 60000.0

// ooxmlFontSizePerPoint is the unit a:rPr states font size in: hundredths of a
// point.
const ooxmlFontSizePerPoint = 100.0

func convertAdjustments(gds []aGd) map[string]int {
	if len(gds) == 0 {
		return nil
	}
	out := make(map[string]int, len(gds))
	for _, gd := range gds {
		value, ok := parseFormulaValue(gd.Fmla)
		if !ok {
			continue
		}
		out[gd.Name] = value
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// parseFormulaValue reads the "val N" formula an adjustment guide carries.
// Guides with a computed formula are left to the renderer's own preset defaults.
func parseFormulaValue(fmla string) (int, bool) {
	rest, ok := strings.CutPrefix(strings.TrimSpace(fmla), "val ")
	if !ok {
		return 0, false
	}
	value, err := strconv.Atoi(strings.TrimSpace(rest))
	if err != nil {
		return 0, false
	}
	return value, true
}

func convertDrawingParagraphs(body *aTxBody) []DrawingParagraph {
	if body == nil || len(body.Ps) == 0 {
		return nil
	}
	out := make([]DrawingParagraph, 0, len(body.Ps))
	for _, p := range body.Ps {
		paragraph := DrawingParagraph{Runs: make([]DrawingRun, 0, len(p.Rs))}
		if p.PPr != nil {
			paragraph.Align = p.PPr.Algn
		}
		for _, r := range p.Rs {
			if r.T == "" {
				continue
			}
			paragraph.Runs = append(paragraph.Runs, convertDrawingRun(r))
		}
		if len(paragraph.Runs) == 0 {
			continue
		}
		out = append(out, paragraph)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func convertDrawingRun(r aR) DrawingRun {
	run := DrawingRun{Text: r.T}
	if r.RPr == nil {
		return run
	}
	run.SizePt = float64(r.RPr.Sz) / ooxmlFontSizePerPoint
	run.Bold = isXMLTrue(r.RPr.B)
	run.Italic = isXMLTrue(r.RPr.I)
	run.Color = colorRefFrom(r.RPr.SolidFill)
	if r.RPr.Latin != nil {
		run.Font = r.RPr.Latin.Typeface
	}
	return run
}

func colorRefFrom(fill *aColorChoice) ColorRef {
	if fill == nil {
		return ColorRef{}
	}
	switch {
	case fill.SrgbClr != nil:
		ref := colorRefFromNode(*fill.SrgbClr)
		ref.SRGB = strings.TrimPrefix(fill.SrgbClr.Val, "#")
		return ref
	case fill.SchemeClr != nil:
		ref := colorRefFromNode(*fill.SchemeClr)
		ref.Scheme = fill.SchemeClr.Val
		return ref
	default:
		return ColorRef{}
	}
}

func colorRefFromNode(node aColorNode) ColorRef {
	return ColorRef{
		HueOff:   valOrZero(node.HueOff),
		SatOff:   valOrZero(node.SatOff),
		LumOff:   valOrZero(node.LumOff),
		LumMod:   valOrZero(node.LumMod),
		SatMod:   valOrZero(node.SatMod),
		Shade:    valOrZero(node.Shade),
		Tint:     valOrZero(node.Tint),
		AlphaPct: valOrZero(node.Alpha),
	}
}

func valOrZero(v *aVal) int {
	if v == nil {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(v.Val))
	if err != nil {
		return 0
	}
	return n
}

// isXMLTrue reads an OOXML boolean attribute, which is "1"/"true" for set and
// absent or "0"/"false" for unset.
func isXMLTrue(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true":
		return true
	default:
		return false
	}
}
