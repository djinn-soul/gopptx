package export

import (
	"math"
	"strings"

	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
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

func consumeBodyPlaceholderAsBullets(sc *elements.SlideContent, es editorcommon.Shape) bool {
	if appendParagraphBullets(sc, es.Paragraphs) {
		return true
	}
	bodyText := strings.TrimSpace(es.Text)
	if bodyText == "" {
		return false
	}
	parts := strings.Split(bodyText, "\n")
	style := editorParagraphToExportStyle(es.Paragraph)
	for _, line := range parts {
		line = strings.TrimSpace(strings.TrimPrefix(line, "•"))
		if line == "" {
			continue
		}
		sc.Bullets = append(sc.Bullets, line)
		sc.BulletRuns = append(sc.BulletRuns, nil)
		sc.BulletStyles = append(sc.BulletStyles, style)
	}
	return len(sc.Bullets) > 0
}

func appendParagraphBullets(sc *elements.SlideContent, paragraphs []editorcommon.ShapeTextParagraph) bool {
	if len(paragraphs) == 0 {
		return false
	}
	startCount := len(sc.Bullets)
	for _, paragraph := range paragraphs {
		runs := editorRunsToExportRuns(paragraph.Runs)
		textValue := strings.TrimSpace(elements.RunsToPlainText(runs))
		if textValue == "" {
			continue
		}
		style := editorParagraphToExportStyle(paragraph.Paragraph)
		sc.Bullets = append(sc.Bullets, textValue)
		if len(runs) > 0 {
			sc.BulletRuns = append(sc.BulletRuns, runs)
		} else {
			sc.BulletRuns = append(sc.BulletRuns, nil)
		}
		sc.BulletStyles = append(sc.BulletStyles, style)
	}
	return len(sc.Bullets) > startCount
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
		bg := "FFFFFF"
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

func editorTextFrameToExportTextFrame(frame *editorcommon.TextFrame) *shapes.TextFrame {
	if frame == nil {
		return nil
	}
	tf := shapes.NewTextFrame()
	has := false
	if frame.MarginLeft != nil {
		tf.MarginLeft = styling.Emu(int64(*frame.MarginLeft))
		has = true
	}
	if frame.MarginRight != nil {
		tf.MarginRight = styling.Emu(int64(*frame.MarginRight))
		has = true
	}
	if frame.MarginTop != nil {
		tf.MarginTop = styling.Emu(int64(*frame.MarginTop))
		has = true
	}
	if frame.MarginBottom != nil {
		tf.MarginBottom = styling.Emu(int64(*frame.MarginBottom))
		has = true
	}
	if frame.VerticalAlign != nil {
		switch strings.ToLower(strings.TrimSpace(*frame.VerticalAlign)) {
		case "t", "top":
			tf.Anchor = shapes.TextAnchorTop
			has = true
		case "b", "bottom":
			tf.Anchor = shapes.TextAnchorBottom
			has = true
		case "ctr", "center", "middle":
			tf.Anchor = shapes.TextAnchorMiddle
			has = true
		}
	}
	if frame.WordWrap != nil {
		if *frame.WordWrap {
			tf.Wrap = shapes.TextWrapSquare
		} else {
			tf.Wrap = shapes.TextWrapNone
		}
		has = true
	}
	if frame.AutoFitType != nil {
		switch strings.ToLower(strings.TrimSpace(*frame.AutoFitType)) {
		case "none":
			tf.AutoFit = shapes.TextAutoFitNone
			has = true
		case "normal":
			tf.AutoFit = shapes.TextAutoFitNormal
			has = true
		case "shape":
			tf.AutoFit = shapes.TextAutoFitShape
			has = true
		}
	} else if frame.AutoFit != nil {
		if *frame.AutoFit {
			tf.AutoFit = shapes.TextAutoFitShape
		} else {
			tf.AutoFit = shapes.TextAutoFitNone
		}
		has = true
	}
	if frame.Rotation != nil {
		tf.RotationDeg = frame.Rotation
		has = true
	}
	if !has {
		return nil
	}
	return &tf
}

func editorParagraphsToExportParagraphs(es editorcommon.Shape) []text.Paragraph {
	if len(es.Paragraphs) > 0 {
		out := make([]text.Paragraph, 0, len(es.Paragraphs))
		for _, paragraph := range es.Paragraphs {
			runs := editorRunsToExportRuns(paragraph.Runs)
			style := editorParagraphToExportStyle(paragraph.Paragraph)
			if len(runs) == 0 && style == elements.DefaultParagraphStyle() {
				continue
			}
			out = append(out, text.Paragraph{
				Runs:  runs,
				Style: style,
			})
		}
		if len(out) > 0 {
			return out
		}
	}
	if len(es.Runs) > 0 || es.Paragraph != nil {
		return []text.Paragraph{{
			Runs:  editorRunsToExportRuns(es.Runs),
			Style: editorParagraphToExportStyle(es.Paragraph),
		}}
	}
	return nil
}

func editorParagraphToExportStyle(paragraph *editorcommon.Paragraph) elements.ParagraphStyle {
	style := elements.DefaultParagraphStyle()
	if paragraph == nil {
		return style
	}
	if paragraph.Alignment != nil {
		style.Align = *paragraph.Alignment
	}
	if paragraph.SpaceBeforePts != nil {
		style.SpaceBeforePt = *paragraph.SpaceBeforePts / textSpacingPtsScale
	}
	if paragraph.SpaceAfterPts != nil {
		style.SpaceAfterPt = *paragraph.SpaceAfterPts / textSpacingPtsScale
	}
	if paragraph.LineSpacingPct != nil {
		style.LineSpacingPct = *paragraph.LineSpacingPct / textSpacingPctScale
	}
	if paragraph.BulletStyle != nil {
		style.BulletStyle = *paragraph.BulletStyle
	}
	if paragraph.BulletChar != nil {
		style.BulletChar = *paragraph.BulletChar
	}
	if paragraph.BulletColor != nil {
		style.BulletColor = *paragraph.BulletColor
	}
	if paragraph.BulletSizePct != nil {
		style.BulletSize = *paragraph.BulletSizePct
	}
	if paragraph.Level != nil {
		style.Level = *paragraph.Level
	}
	if paragraph.Indent != nil {
		style.LeftIndent = styling.Emu(int64(*paragraph.Indent))
	}
	if paragraph.Hanging != nil {
		style.HangingIndent = styling.Emu(int64(*paragraph.Hanging))
	}
	return style
}

func editorRunsToExportRuns(runs []editorcommon.TextRun) []elements.Run {
	if len(runs) == 0 {
		return nil
	}
	out := make([]elements.Run, 0, len(runs))
	for _, run := range runs {
		exportRun := elements.NewRun(run.Text)
		if run.Bold != nil {
			exportRun.Bold = *run.Bold
		}
		if run.Italic != nil {
			exportRun.Italic = *run.Italic
		}
		if run.Underline != nil {
			exportRun.Underline = *run.Underline
		}
		if run.Strikethrough != nil {
			exportRun.Strikethrough = *run.Strikethrough
		}
		if run.Subscript != nil {
			exportRun.Subscript = *run.Subscript
		}
		if run.Superscript != nil {
			exportRun.Superscript = *run.Superscript
		}
		if run.Color != nil {
			exportRun.Color = *run.Color
		}
		if run.Highlight != nil {
			exportRun.Highlight = *run.Highlight
		}
		if run.Font != nil {
			exportRun.Font = *run.Font
		}
		if run.SizePt != nil {
			exportRun.SizePt = *run.SizePt
		}
		if run.Code != nil {
			exportRun.Code = *run.Code
		}
		if run.AllCaps != nil {
			exportRun.AllCaps = *run.AllCaps
		}
		if run.SmallCaps != nil {
			exportRun.SmallCaps = *run.SmallCaps
		}
		if run.OutlineColor != nil {
			exportRun.OutlineColor = *run.OutlineColor
		}
		if run.OutlineWidthPt != nil {
			exportRun.OutlineWidthPt = *run.OutlineWidthPt
		}
		exportRun.Hyperlink = editorHyperlinkToExportHyperlink(run.Hyperlink)
		exportRun.HoverAction = editorHyperlinkToExportHyperlink(run.HoverAction)
		out = append(out, exportRun)
	}
	return elements.NormalizeRuns(out)
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

func editorEffectsToExportEffects(es editorcommon.Shape) *shapes.ShapeEffects {
	effects := &shapes.ShapeEffects{}
	if es.Shadow != nil && (es.Shadow.Inherit == nil || *es.Shadow.Inherit || es.Shadow.Color != nil ||
		es.Shadow.BlurEmu != nil || es.Shadow.DistanceEmu != nil || es.Shadow.AngleDeg != nil) {
		effects.Shadow = true
	}
	effects.Glow = es.Glow != nil
	effects.SoftEdges = es.SoftEdge != nil || es.Blur != nil
	effects.Reflection = es.Reflection != nil
	if !effects.Shadow && !effects.Glow && !effects.SoftEdges && !effects.Reflection {
		return nil
	}
	return effects
}

func positionPct(raw *float64) int {
	if raw == nil {
		return 0
	}
	return int(math.Round(*raw))
}
