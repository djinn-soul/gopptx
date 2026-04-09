package export

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestFoldGeneratedConnectorLabelsAssignsConnectorLabel(t *testing.T) {
	startX := styling.Emu(connectorLabelWidthEMU)
	startY := styling.Emu(connectorLabelHeightEMU)
	endX := styling.Emu(connectorLabelWidthEMU * 3)
	endY := styling.Emu(connectorLabelHeightEMU * 3)
	labelX := (startX+endX)/2 - styling.Emu(connectorLabelWidthEMU)/2
	labelY := (startY+endY)/2 - styling.Emu(connectorLabelHeightEMU)/2
	sc := elements.SlideContent{
		Shapes: []shapes.Shape{
			shapes.NewShape(
				shapes.ShapeTypeRectangle,
				labelX,
				labelY,
				styling.Emu(connectorLabelWidthEMU),
				styling.Emu(connectorLabelHeightEMU),
			).WithName("Connector Label 7").WithText("Link"),
		},
		Connectors: []shapes.Connector{
			shapes.NewStraightConnector(
				startX,
				startY,
				endX,
				endY,
			),
		},
	}

	foldGeneratedConnectorLabels(&sc)

	if len(sc.Shapes) != 0 {
		t.Fatalf("expected generated label shape to be removed, got %d shapes", len(sc.Shapes))
	}
	if len(sc.Connectors) != 1 || sc.Connectors[0].Label != "Link" {
		t.Fatalf("expected connector label to be assigned, got %+v", sc.Connectors)
	}
}

func TestFoldGeneratedConnectorLabelsPreservesUnmatchedShape(t *testing.T) {
	sc := elements.SlideContent{
		Shapes: []shapes.Shape{
			shapes.NewShape(
				shapes.ShapeTypeRectangle,
				0,
				0,
				styling.Emu(connectorLabelWidthEMU),
				styling.Emu(connectorLabelHeightEMU),
			).WithName("Connector Label 7").WithText("Standalone"),
		},
		Connectors: []shapes.Connector{
			shapes.NewStraightConnector(styling.Inches(1), styling.Inches(1), styling.Inches(2), styling.Inches(2)),
		},
	}

	foldGeneratedConnectorLabels(&sc)

	if len(sc.Shapes) != 1 {
		t.Fatalf("expected unmatched label-like shape to be preserved, got %d shapes", len(sc.Shapes))
	}
	if sc.Connectors[0].Label != "" {
		t.Fatalf("expected connector label to stay empty, got %q", sc.Connectors[0].Label)
	}
}
