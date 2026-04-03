package export

import (
	"testing"

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
	wrap := false
	autoFitType := "normal"
	shadowColor := "778899"
	shadowBlur := 32000
	shadowDist := 21000
	shapeRotation := 33.0

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
			WordWrap:      &wrap,
			AutoFitType:   &autoFitType,
		},
		Shadow: &editorcommon.ShapeShadow{
			Color:       &shadowColor,
			BlurEmu:     &shadowBlur,
			DistanceEmu: &shadowDist,
		},
		Glow:     &editorcommon.ShapeGlow{},
		SoftEdge: &editorcommon.ShapeSoftEdge{},
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
	if shape.RichShadow == nil || shape.RichShadow.Color != shadowColor ||
		shape.RichShadow.BlurRadius != shadowBlur || shape.RichShadow.Distance != shadowDist {
		t.Fatalf("expected rich shadow, got %+v", shape.RichShadow)
	}
	if shape.Effects == nil || !shape.Effects.Shadow || !shape.Effects.Glow || !shape.Effects.SoftEdges {
		t.Fatalf("expected effects flags, got %+v", shape.Effects)
	}
	if shape.RotationDeg == nil || *shape.RotationDeg != 33 {
		t.Fatalf("expected rounded shape rotation 33, got %+v", shape.RotationDeg)
	}
}

func TestConsumeBodyPlaceholderAsBulletsPreservesParagraphStyle(t *testing.T) {
	level := 2
	spaceBefore := 600
	spaceAfter := 400
	lineSpacing := 125000
	bulletStyle := "roman_lower"
	bulletColor := "ABCDEF"
	bulletSize := 90
	indent := 228600
	hanging := 114300
	slide := elements.SlideContent{}

	ok := consumeBodyPlaceholderAsBullets(&slide, editorcommon.Shape{
		Text: "First\nSecond",
		Paragraph: &editorcommon.Paragraph{
			Level:          &level,
			SpaceBeforePts: &spaceBefore,
			SpaceAfterPts:  &spaceAfter,
			LineSpacingPct: &lineSpacing,
			BulletStyle:    &bulletStyle,
			BulletColor:    &bulletColor,
			BulletSizePct:  &bulletSize,
			Indent:         &indent,
			Hanging:        &hanging,
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
	if int(style.LeftIndent.Emu()) != indent || int(style.HangingIndent.Emu()) != hanging {
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
