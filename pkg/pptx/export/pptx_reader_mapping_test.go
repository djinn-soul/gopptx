package export

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestEditorShapeToShapePreservesSupportedFormatting(t *testing.T) {
	fillColor := "112233"
	fillTransparency := 0.25
	lineColor := "445566"
	lineWidth := 25400
	lineDash := "dash"
	marginLeft := 91440
	marginTop := 45720
	verticalAlign := "b"
	orientation := "vert270"
	columns := 2
	wrap := false
	autoFitType := "normal"
	shadowColor := "778899"
	shadowBlur := 32000
	shadowDist := 21000
	shapeRotation := 33.0
	glowColor := "ABCDEF"
	glowRadius := 50800
	blurRadius := 12700
	softEdgeRadius := 25400
	reflectionBlur := 9000
	reflectionDistance := 18000

	shape := editorShapeToShape(editorcommon.Shape{
		Type: "rect",
		X:    10,
		Y:    20,
		W:    30,
		H:    40,
		Text: "Hello",
		Name: "Shape 1",
		Fill: &editorcommon.ShapeFill{
			Solid:        &fillColor,
			Transparency: &fillTransparency,
		},
		Line: &editorcommon.ShapeLine{
			Color:     &lineColor,
			WidthEmu:  &lineWidth,
			DashStyle: &lineDash,
		},
		TextFrame: &editorcommon.TextFrame{
			MarginLeft:    &marginLeft,
			MarginTop:     &marginTop,
			VerticalAlign: &verticalAlign,
			Orientation:   &orientation,
			Columns:       &columns,
			WordWrap:      &wrap,
			AutoFitType:   &autoFitType,
		},
		Shadow: &editorcommon.ShapeShadow{
			Color:       &shadowColor,
			BlurEmu:     &shadowBlur,
			DistanceEmu: &shadowDist,
		},
		Glow: &editorcommon.ShapeGlow{
			Color:     &glowColor,
			RadiusEmu: &glowRadius,
		},
		Blur: &editorcommon.ShapeBlur{
			RadiusEmu: &blurRadius,
		},
		SoftEdge: &editorcommon.ShapeSoftEdge{
			RadiusEmu: &softEdgeRadius,
		},
		Reflection: &editorcommon.ShapeReflection{
			BlurEmu:     &reflectionBlur,
			DistanceEmu: &reflectionDistance,
		},
		Rotation: &shapeRotation,
	})

	if shape.Fill == nil || shape.Fill.Color != fillColor || shape.Fill.Transparency == nil ||
		*shape.Fill.Transparency != fillTransparency {
		t.Fatalf("expected solid fill with transparency, got %+v", shape.Fill)
	}
	if shape.Line == nil || shape.Line.Color != lineColor || shape.Line.Dash != lineDash {
		t.Fatalf("expected line color/dash, got %+v", shape.Line)
	}
	if shape.TextFrame == nil {
		t.Fatalf("expected text frame, got nil")
	}
	if got := int(shape.TextFrame.MarginLeft.Emu()); got != marginLeft {
		t.Fatalf("expected left margin %d, got %d", marginLeft, got)
	}
	if got := int(shape.TextFrame.MarginTop.Emu()); got != marginTop {
		t.Fatalf("expected top margin %d, got %d", marginTop, got)
	}
	if shape.TextFrame.Anchor != "b" {
		t.Fatalf("expected bottom anchor, got %s", shape.TextFrame.Anchor)
	}
	if shape.TextFrame.Wrap != "none" {
		t.Fatalf("expected no wrap, got %s", shape.TextFrame.Wrap)
	}
	if shape.TextFrame.AutoFit != "normAutoFit" {
		t.Fatalf("expected normal autofit, got %s", shape.TextFrame.AutoFit)
	}
	if shape.TextFrame.Orientation != orientation {
		t.Fatalf("expected orientation %s, got %s", orientation, shape.TextFrame.Orientation)
	}
	if shape.TextFrame.Columns != columns {
		t.Fatalf("expected columns %d, got %d", columns, shape.TextFrame.Columns)
	}
	if shape.RichShadow == nil || shape.RichShadow.Color != shadowColor ||
		shape.RichShadow.BlurRadius != shadowBlur || shape.RichShadow.Distance != shadowDist {
		t.Fatalf("expected rich shadow, got %+v", shape.RichShadow)
	}
	if shape.Effects == nil || !shape.Effects.Shadow || !shape.Effects.Glow || !shape.Effects.SoftEdges {
		t.Fatalf("expected effects flags, got %+v", shape.Effects)
	}
	if shape.Effects.BlurSpec == nil || shape.Effects.BlurSpec.RadiusEmu != blurRadius {
		t.Fatalf("expected blur radius %d, got %+v", blurRadius, shape.Effects.BlurSpec)
	}
	if shape.Effects.GlowSpec == nil || shape.Effects.GlowSpec.Color != glowColor ||
		shape.Effects.GlowSpec.RadiusEmu != glowRadius {
		t.Fatalf("expected glow detail, got %+v", shape.Effects.GlowSpec)
	}
	if shape.Effects.SoftEdgeSpec == nil || shape.Effects.SoftEdgeSpec.RadiusEmu != softEdgeRadius {
		t.Fatalf("expected soft edge radius %d, got %+v", softEdgeRadius, shape.Effects.SoftEdgeSpec)
	}
	if shape.Effects.ReflectionSpec == nil || shape.Effects.ReflectionSpec.BlurEmu != reflectionBlur ||
		shape.Effects.ReflectionSpec.DistanceEmu != reflectionDistance {
		t.Fatalf("expected reflection detail, got %+v", shape.Effects.ReflectionSpec)
	}
	if shape.RotationDeg == nil || *shape.RotationDeg != 33 {
		t.Fatalf("expected rounded shape rotation 33, got %+v", shape.RotationDeg)
	}
}

func TestEditorHyperlinkToExportHyperlinkPreservesActionOnlyMetadata(t *testing.T) {
	tooltip := "Hover"
	actionValue := "ppaction://macro?name=RunMacro"
	history := false
	endSound := true
	highlight := false

	got := editorHyperlinkToExportHyperlink(&editorcommon.Hyperlink{
		Action:         &actionValue,
		Tooltip:        &tooltip,
		History:        &history,
		EndSound:       &endSound,
		HighlightClick: &highlight,
	})
	if got == nil {
		t.Fatal("expected exported hyperlink")
	}
	if got.ActionType() != actionValue {
		t.Fatalf("expected raw action %q, got %q", actionValue, got.ActionType())
	}
	if got.History == nil || *got.History != history {
		t.Fatalf("expected history=%v, got %+v", history, got.History)
	}
	if got.EndSound == nil || *got.EndSound != endSound {
		t.Fatalf("expected end_sound=%v, got %+v", endSound, got.EndSound)
	}
	if got.HighlightClick != highlight || got.Tooltip != tooltip {
		t.Fatalf("expected tooltip/highlight preserved, got %+v", got)
	}
}

func TestEditorShapeToConnectorPreservesLabelAndActions(t *testing.T) {
	macroAction := "ppaction://macro?name=RunConnector"
	shapeIndexByID := map[int]int{11: 1, 22: 2}
	startID := 11
	endID := 22
	source := editorcommon.Shape{
		Type: "straightConnector1",
		Text: "Connector Label",
		ClickAction: &editorcommon.Hyperlink{
			Action: &macroAction,
		},
		Connector: &editorcommon.ConnectorInfo{
			StartShapeID: &startID,
			EndShapeID:   &endID,
		},
	}

	connector, ok := editorShapeToConnector(source, shapeIndexByID)
	if !ok {
		t.Fatal("expected connector mapping")
	}
	if connector.Label != "Connector Label" {
		t.Fatalf("expected label to round-trip, got %q", connector.Label)
	}
	if connector.ClickAction == nil || connector.ClickAction.ActionType() != macroAction {
		t.Fatalf("expected connector click action to preserve raw action, got %+v", connector.ClickAction)
	}
	if connector.StartShapeIndex != 1 || connector.EndShapeIndex != 2 {
		t.Fatalf("expected mapped connector anchors, got %+v", connector)
	}
}

func TestEditorEffectsToExportEffectsNilWhenEmpty(t *testing.T) {
	if got := editorEffectsToExportEffects(editorcommon.Shape{}); got != nil {
		t.Fatalf("expected nil effects for empty shape, got %+v", got)
	}
}

func TestEditorShapeToConnectorPreservesHyperlinkKinds(t *testing.T) {
	address := "https://example.com"
	connector, ok := editorShapeToConnector(editorcommon.Shape{
		Type: "straightConnector1",
		ClickAction: &editorcommon.Hyperlink{
			Address: &address,
		},
	}, nil)
	if !ok {
		t.Fatal("expected connector mapping")
	}
	if connector.ClickAction == nil || connector.ClickAction.Action.Type != action.HyperlinkActionURL {
		t.Fatalf("expected URL connector action, got %+v", connector.ClickAction)
	}
}

func TestEditorEffectsToExportEffectsPreservesExplicitZeroReflection(t *testing.T) {
	zero := 0
	got := editorEffectsToExportEffects(editorcommon.Shape{
		Reflection: &editorcommon.ShapeReflection{
			BlurEmu:     &zero,
			DistanceEmu: &zero,
		},
	})
	if got == nil || got.ReflectionSpec == nil {
		t.Fatalf("expected explicit reflection spec, got %+v", got)
	}
	if got.ReflectionSpec.BlurEmu != 0 || got.ReflectionSpec.DistanceEmu != 0 {
		t.Fatalf("expected zero reflection values preserved, got %+v", got.ReflectionSpec)
	}
}

func TestEditorShapeToShapePreservesActionOnlyHyperlinks(t *testing.T) {
	actionValue := "ppaction://macro?name=ShapeMacro"
	shape := editorShapeToShape(editorcommon.Shape{
		Type: "rect",
		X:    1,
		Y:    2,
		W:    3,
		H:    4,
		ClickAction: &editorcommon.Hyperlink{
			Action: &actionValue,
		},
	})
	if shape.ClickAction == nil || shape.ClickAction.ActionType() != actionValue {
		t.Fatalf("expected action-only hyperlink preserved, got %+v", shape.ClickAction)
	}
	if shape.ClickAction.RequiresRelationship() {
		t.Fatalf("expected raw action to skip relationship creation, got %+v", shape.ClickAction)
	}
}

func TestConsumeBodyPlaceholderAsBulletsPreservesParagraphStyle(t *testing.T) {
	level := 2
	spaceBefore := 600
	spaceAfter := 400
	lineSpacing := 125000
	lineSpacingPts := 1800
	bulletStyle := "roman_lower"
	bulletColor := "ABCDEF"
	bulletSize := 90
	indent := 228600
	hanging := 114300
	tabStops := []int{228600, 457200}
	slide := elements.SlideContent{}

	ok := consumeBodyPlaceholderAsBullets(&slide, editorcommon.Shape{
		Text: "First\nSecond",
		Paragraph: &editorcommon.Paragraph{
			Level:          &level,
			SpaceBeforePts: &spaceBefore,
			SpaceAfterPts:  &spaceAfter,
			LineSpacingPct: &lineSpacing,
			LineSpacingPts: &lineSpacingPts,
			BulletStyle:    &bulletStyle,
			BulletColor:    &bulletColor,
			BulletSizePct:  &bulletSize,
			Indent:         &indent,
			Hanging:        &hanging,
			TabStops:       tabStops,
		},
	})
	if !ok {
		t.Fatal("expected bullets to be consumed")
	}
	if len(slide.Bullets) != 2 || slide.Bullets[0] != "First" || slide.Bullets[1] != "Second" {
		t.Fatalf("unexpected bullets: %+v", slide.Bullets)
	}
	if len(slide.BulletStyles) != 2 {
		t.Fatalf("expected 2 bullet styles, got %d", len(slide.BulletStyles))
	}
	style := slide.BulletStyles[0]
	if style.Level != level || style.SpaceBeforePt != 6 || style.SpaceAfterPt != 4 ||
		style.LineSpacingPct != 125 || style.BulletStyle != bulletStyle ||
		style.BulletColor != bulletColor || style.BulletSize != bulletSize {
		t.Fatalf("unexpected bullet style: %+v", style)
	}
	if style.LineSpacingPts != 18 {
		t.Fatalf("expected line spacing points 18, got %+v", style)
	}
	if len(style.TabStops) != 2 || int(style.TabStops[0].Emu()) != tabStops[0] ||
		int(style.TabStops[1].Emu()) != tabStops[1] {
		t.Fatalf("unexpected tab stops: %+v", style.TabStops)
	}
	if int(style.LeftIndent.Emu()) != indent || int(style.HangingIndent.Emu()) != -hanging {
		t.Fatalf("unexpected bullet indents: %+v", style)
	}
}

func TestEditorShapeToShapePreservesTextParagraphs(t *testing.T) {
	bold := true
	color := "FF0000"
	size := 18
	level := 1
	bulletStyle := "number"

	shape := editorShapeToShape(editorcommon.Shape{
		Text: "Alpha\nBeta",
		Paragraphs: []editorcommon.ShapeTextParagraph{
			{
				Runs: []editorcommon.TextRun{{
					Text:   "Alpha",
					Bold:   &bold,
					Color:  &color,
					SizePt: &size,
				}},
				Paragraph: &editorcommon.Paragraph{
					Level:       &level,
					BulletStyle: &bulletStyle,
				},
			},
			{
				Runs: []editorcommon.TextRun{{Text: "Beta"}},
			},
		},
	})

	if len(shape.TextParagraphs) != 2 {
		t.Fatalf("expected 2 text paragraphs, got %d", len(shape.TextParagraphs))
	}
	first := shape.TextParagraphs[0]
	if len(first.Runs) != 1 || first.Runs[0].Text != "Alpha" || !first.Runs[0].Bold ||
		first.Runs[0].Color != color || first.Runs[0].SizePt != size {
		t.Fatalf("unexpected first paragraph runs: %+v", first.Runs)
	}
	if first.Style.Level != level || first.Style.BulletStyle != bulletStyle {
		t.Fatalf("unexpected first paragraph style: %+v", first.Style)
	}
	if len(shape.TextParagraphs[1].Runs) != 1 || shape.TextParagraphs[1].Runs[0].Text != "Beta" {
		t.Fatalf("unexpected second paragraph: %+v", shape.TextParagraphs[1])
	}
}
