package export

import (
	"testing"

	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestEditorShapeToConnectorPreservesAnchorsAndStyle(t *testing.T) {
	lineColor := "4472C4"
	lineWidth := 12700
	startArrow := "triangle"
	endArrow := "stealth"
	startShapeID := 10
	endShapeID := 11
	startSiteIndex := 1
	endSiteIndex := 3

	connector, ok := editorShapeToConnector(editorcommon.Shape{
		ID:           20,
		Type:         shapes.ConnectorTypeElbow,
		X:            100,
		Y:            200,
		W:            300,
		H:            400,
		AltText:      "Connector Alt",
		IsDecorative: true,
		Line: &editorcommon.ShapeLine{
			Color:      &lineColor,
			WidthEmu:   &lineWidth,
			StartArrow: &startArrow,
			EndArrow:   &endArrow,
		},
		Connector: &editorcommon.ConnectorInfo{
			StartShapeID:   &startShapeID,
			StartSiteIndex: &startSiteIndex,
			EndShapeID:     &endShapeID,
			EndSiteIndex:   &endSiteIndex,
			FlipH:          true,
		},
	}, map[int]int{10: 1, 11: 2})
	if !ok {
		t.Fatal("expected connector mapping to succeed")
	}
	if connector.StartShapeIndex != 1 || connector.EndShapeIndex != 2 {
		t.Fatalf("expected mapped anchor indices, got %+v", connector)
	}
	if connector.StartSite != shapes.ConnectionSiteRight || connector.EndSite != shapes.ConnectionSiteLeft {
		t.Fatalf("expected mapped site names, got %+v", connector)
	}
	if connector.StartX != 400 || connector.EndX != 100 || connector.StartY != 200 || connector.EndY != 600 {
		t.Fatalf("expected flip-adjusted endpoints, got %+v", connector)
	}
	if connector.Line.Color != lineColor || int64(connector.Line.Width) != int64(lineWidth) {
		t.Fatalf("expected connector line style, got %+v", connector.Line)
	}
	if connector.StartArrow != startArrow || connector.EndArrow != endArrow {
		t.Fatalf("expected connector arrows, got %+v", connector)
	}
	if connector.AltText != "Connector Alt" || !connector.IsDecorative {
		t.Fatalf("expected connector accessibility, got %+v", connector)
	}
}
