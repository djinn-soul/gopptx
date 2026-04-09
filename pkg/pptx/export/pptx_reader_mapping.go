package export

import (
	"math"

	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	textSpacingPtsScale = 100
	textSpacingPctScale = 1000
)

func editorShapeToShape(es editorcommon.Shape) shapes.Shape {
	shape := shapes.Shape{
		Type:           editorTypeToPreset(es.Type),
		X:              styling.Emu(int64(es.X)),
		Y:              styling.Emu(int64(es.Y)),
		CX:             styling.Emu(int64(es.W)),
		CY:             styling.Emu(int64(es.H)),
		Text:           es.Text,
		Name:           es.Name,
		AltText:        es.AltText,
		IsDecorative:   es.IsDecorative,
		Adjustments:    editorAdjustmentsToExport(es.Adjustments),
		TextParagraphs: editorParagraphsToExportParagraphs(es),
		ClickAction:    editorHyperlinkToExportHyperlink(es.ClickAction),
		HoverAction:    editorHyperlinkToExportHyperlink(es.HoverAction),
	}
	if fill := editorFillToExportFill(es.Fill); fill != nil {
		shape.Fill = fill
	}
	if grad := editorGradientToExportFill(es.Fill); grad != nil {
		shape.GradientFill = grad
	}
	if richFill := editorRichFillToExportFill(es.Fill); richFill != nil {
		shape.RichFill = richFill
	}
	if line := editorLineToExportLine(es.Line); line != nil {
		shape.Line = line
	}
	if richLine := editorRichLineToExportLine(es.Line); richLine != nil {
		shape.RichLine = richLine
	}
	if shadow := editorShadowToExportShadow(es.Shadow); shadow != nil {
		shape.RichShadow = shadow
	}
	if effects := editorEffectsToExportEffects(es); effects != nil {
		shape.Effects = effects
	}
	if tf := editorTextFrameToExportTextFrame(es.TextFrame); tf != nil {
		shape.TextFrame = tf
	}
	if es.Rotation != nil {
		deg := int(math.Round(*es.Rotation))
		shape.RotationDeg = &deg
	}
	return shape
}

func editorFillToExportFill(fill *editorcommon.ShapeFill) *shapes.ShapeFill {
	if fill == nil || fill.Solid == nil || *fill.Solid == "" {
		return nil
	}
	exportFill := &shapes.ShapeFill{Color: *fill.Solid}
	if fill.Transparency != nil {
		exportFill.Transparency = fill.Transparency
	}
	return exportFill
}

func editorGradientToExportFill(fill *editorcommon.ShapeFill) *shapes.ShapeGradientFill {
	if fill == nil || fill.Gradient == nil {
		return nil
	}
	stops := make([]shapes.ShapeGradientStop, 0, len(fill.Gradient.Stops))
	for _, stop := range fill.Gradient.Stops {
		shapeStop := shapes.NewShapeGradientStop(positionPct(stop.PositionPct), stop.Color)
		stops = append(stops, shapeStop)
	}
	gradient := shapes.NewShapeGradientFill(shapes.ShapeGradientTypeLinear, stops)
	if fill.Gradient.AngleDeg != nil {
		gradient = gradient.WithLinearAngle(int(math.Round(*fill.Gradient.AngleDeg)))
	}
	return &gradient
}

func editorRichFillToExportFill(fill *editorcommon.ShapeFill) *shapes.RichShapeFill {
	if fill == nil {
		return nil
	}
	if fill.Background != nil && *fill.Background {
		return shapes.NewNoFill()
	}
	if fill.Pattern != nil {
		preset := ""
		if fill.Pattern.Preset != nil {
			preset = *fill.Pattern.Preset
		}
		richFill := shapes.NewPatternFill(shapes.NormalizePatternType(preset))
		fg := "000000"
		bg := pdfTableHeaderText
		if fill.Pattern.FgColor != nil && *fill.Pattern.FgColor != "" {
			fg = *fill.Pattern.FgColor
		}
		if fill.Pattern.BgColor != nil && *fill.Pattern.BgColor != "" {
			bg = *fill.Pattern.BgColor
		}
		return richFill.WithPatternColors(fg, bg)
	}
	return nil
}

func editorLineToExportLine(line *editorcommon.ShapeLine) *shapes.ShapeLine {
	if line == nil || line.Color == nil || *line.Color == "" {
		return nil
	}
	width := styling.Emu(0)
	if line.WidthEmu != nil && *line.WidthEmu > 0 {
		width = styling.Emu(int64(*line.WidthEmu))
	}
	exportLine := shapes.NewShapeLine(*line.Color, width)
	if line.DashStyle != nil {
		exportLine.Dash = shapes.NormalizeDrawingLineDash(*line.DashStyle)
	}
	return &exportLine
}

func editorRichLineToExportLine(line *editorcommon.ShapeLine) *shapes.RichShapeLine {
	if line == nil || line.Color == nil || *line.Color == "" {
		return nil
	}
	width := styling.Emu(0)
	if line.WidthEmu != nil && *line.WidthEmu > 0 {
		width = styling.Emu(int64(*line.WidthEmu))
	}
	richLine := shapes.NewRichShapeLine(*line.Color, width)
	if line.DashStyle != nil {
		richLine.DashStyle = shapes.NormalizeLineDashStyle(*line.DashStyle)
	}
	return richLine
}

func editorAdjustmentsToExport(
	adjustments []editorcommon.ShapeAdjustment,
) []shapes.ShapeAdjustment {
	if len(adjustments) == 0 {
		return nil
	}
	out := make([]shapes.ShapeAdjustment, 0, len(adjustments))
	for _, adjustment := range adjustments {
		if adjustment.Name == "" || adjustment.Formula == "" {
			continue
		}
		out = append(out, shapes.ShapeAdjustment{
			Name:    adjustment.Name,
			Formula: adjustment.Formula,
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func editorShadowToExportShadow(shadow *editorcommon.ShapeShadow) *shapes.RichShapeShadow {
	if shadow == nil {
		return nil
	}
	if shadow.Inherit != nil && !*shadow.Inherit &&
		shadow.Color == nil && shadow.BlurEmu == nil && shadow.DistanceEmu == nil && shadow.AngleDeg == nil {
		return nil
	}
	richShadow := shapes.NewRichShapeShadow()
	if shadow.Color != nil && *shadow.Color != "" {
		richShadow.Color = *shadow.Color
	}
	if shadow.BlurEmu != nil {
		richShadow.BlurRadius = *shadow.BlurEmu
	}
	if shadow.DistanceEmu != nil {
		richShadow.Distance = *shadow.DistanceEmu
	}
	if shadow.AngleDeg != nil {
		richShadow.Angle = *shadow.AngleDeg
	}
	return richShadow
}

func positionPct(raw *float64) int {
	if raw == nil {
		return 0
	}
	return int(math.Round(*raw))
}
