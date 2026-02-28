package pptxxml

import (
	"strings"
	"testing"
)

func TestCustomShapeXML(t *testing.T) {
	tests := []struct {
		name     string
		shape    ShapeSpec
		shapeID  int
		contains []string
	}{
		{
			name: "basic shape",
			shape: ShapeSpec{
				Type: "rect",
				X:    100, Y: 200, CX: 300, CY: 400,
				Name: "My Rect",
			},
			shapeID: 10,
			contains: []string{
				"<p:sp>",
				"id=\"10\"",
				"name=\"My Rect\"",
				"prst=\"rect\"",
				"x=\"100\"", "y=\"200\"", "cx=\"300\"", "cy=\"400\"",
			},
		},
		{
			name: "shape with solid fill and transparency",
			shape: ShapeSpec{
				Type: "rect",
				Fill: &ShapeFillSpec{
					Color:        "FF0000",
					Transparency: floatPtr(0.5),
				},
			},
			shapeID: 1,
			contains: []string{
				"<a:solidFill>",
				"val=\"FF0000\"",
				"<a:alpha val=\"50000\"/>",
			},
		},
		{
			name: "shape with rotation",
			shape: ShapeSpec{
				Type:        "rect",
				RotationDeg: intPtr(45),
			},
			shapeID: 1,
			contains: []string{
				"rot=\"2700000\"",
			},
		},
		{
			name: "shape with effects",
			shape: ShapeSpec{
				Type: "rect",
				Effects: &ShapeEffectsSpec{
					Shadow:     true,
					Glow:       true,
					SoftEdges:  true,
					Reflection: true,
				},
			},
			shapeID: 1,
			contains: []string{
				"<a:effectLst>",
				"<a:outerShdw",
				"<a:glow",
				"<a:softEdge",
				"<a:ref",
			},
		},
		{
			name: "shape with text and autofit",
			shape: ShapeSpec{
				Type: "rect",
				Text: "Hello World",
				TextFrame: &TextFrameSpec{
					AutoFit: "spAutoFit",
					Wrap:    "square",
					Anchor:  "ctr",
				},
			},
			shapeID: 1,
			contains: []string{
				"<p:txBody>",
				"<a:spAutoFit/>",
				"Hello World",
				"wrap=\"square\"",
				"anchor=\"ctr\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := customShapeXML(tt.shape, tt.shapeID)
			for _, s := range tt.contains {
				if !strings.Contains(got, s) {
					t.Errorf("customShapeXML() = %v, missing %v", got, s)
				}
			}
		})
	}
}

func TestConnectorXML(t *testing.T) {
	tests := []struct {
		name         string
		connector    ConnectorSpec
		shapeID      int
		startShapeID int
		endShapeID   int
		contains     []string
	}{
		{
			name: "basic connector",
			connector: ConnectorSpec{
				Type:   "bentConnector3",
				StartX: 0, StartY: 0, EndX: 1000, EndY: 1000,
				Line: ShapeLineSpec{
					Color: "000000",
					Width: 9525,
				},
			},
			shapeID:      5,
			startShapeID: 1,
			endShapeID:   2,
			contains: []string{
				"<p:cxnSp>",
				"id=\"5\"",
				"prst=\"bentConnector3\"",
				"val=\"000000\"",
				"w=\"9525\"",
			},
		},
		{
			name: "connector with arrows and dash",
			connector: ConnectorSpec{
				Type: "straightConnector1",
				Line: ShapeLineSpec{
					Dash: "dash",
				},
				StartArrow: "triangle",
				EndArrow:   "stealth",
			},
			shapeID: 1,
			contains: []string{
				"<a:prstDash val=\"dash\"/>",
				"<a:headEnd type=\"triangle\"",
				"<a:tailEnd type=\"stealth\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := connectorXML(tt.connector, tt.shapeID, tt.startShapeID, tt.endShapeID)
			for _, s := range tt.contains {
				if !strings.Contains(got, s) {
					t.Errorf("connectorXML() = %v, missing %v", got, s)
				}
			}
		})
	}
}

func TestConnectorLabelShape(t *testing.T) {
	connector := ConnectorSpec{
		StartX: 0, StartY: 0, EndX: 1000, EndY: 1000,
		Label: "My Label",
	}
	got := connectorLabelShape(connector, 10)
	if !strings.Contains(got, "My Label") {
		t.Errorf("connectorLabelShape() missing label: %v", got)
	}
	if !strings.Contains(got, "id=\"10\"") {
		t.Errorf("connectorLabelShape() missing id: %v", got)
	}
}

func TestAlphaFromNormalizedTransparency(t *testing.T) {
	tests := []struct {
		input float64
		want  int
	}{
		{0.0, 100000},
		{0.5, 50000},
		{1.0, 0},
	}

	for _, tt := range tests {
		if got := alphaFromNormalizedTransparency(tt.input); got != tt.want {
			t.Errorf("alphaFromNormalizedTransparency(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func floatPtr(f float64) *float64 { return &f }
func intPtr(i int) *int          { return &i }
