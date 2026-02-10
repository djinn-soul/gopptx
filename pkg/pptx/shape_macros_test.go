package pptx

import (
	"testing"
)

func TestShapeMacros(t *testing.T) {
	tests := []struct {
		name     string
		shape    Shape
		wantType string
	}{
		{"Rectangle", NewRectangle(1, 1, 2, 2), ShapeTypeRectangle},
		{"RoundedRectangle", NewRoundedRectangle(1, 1, 2, 2), ShapeTypeRoundedRectangle},
		{"Ellipse", NewEllipse(1, 1, 2, 2), ShapeTypeEllipse},
		{"Triangle", NewTriangle(1, 1, 2, 2), ShapeTypeTriangle},
		{"RightTriangle", NewRightTriangle(1, 1, 2, 2), ShapeTypeRightTriangle},
		{"Diamond", NewDiamond(1, 1, 2, 2), ShapeTypeDiamond},
		{"Pentagon", NewPentagon(1, 1, 2, 2), ShapeTypePentagon},
		{"Hexagon", NewHexagon(1, 1, 2, 2), ShapeTypeHexagon},
		{"Parallelogram", NewParallelogram(1, 1, 2, 2), ShapeTypeParallelogram},
		{"FlowChartProcess", NewFlowChartProcess(1, 1, 2, 2), ShapeTypeFlowChartProcess},
		{"FlowChartDecision", NewFlowChartDecision(1, 1, 2, 2), ShapeTypeFlowChartDecision},
		{"FlowChartTerminator", NewFlowChartTerminator(1, 1, 2, 2), ShapeTypeFlowChartTerminator},
		{"RightArrow", NewRightArrow(1, 1, 2, 2), ShapeTypeRightArrow},
		{"LeftArrow", NewLeftArrow(1, 1, 2, 2), ShapeTypeLeftArrow},
		{"UpArrow", NewUpArrow(1, 1, 2, 2), ShapeTypeUpArrow},
		{"DownArrow", NewDownArrow(1, 1, 2, 2), ShapeTypeDownArrow},
		{"Cloud", NewCloud(1, 1, 2, 2), ShapeTypeCloud},
		{"TextBox", NewTextBox("text", 1, 1, 2, 2), ShapeTypeRectangle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shape.Type != tt.wantType {
				t.Errorf("New%s() type = %v, want %v", tt.name, tt.shape.Type, tt.wantType)
			}
			if tt.shape.X != Inches(1) {
				t.Errorf("New%s() X = %v, want %v", tt.name, tt.shape.X, Inches(1))
			}
			if tt.shape.Y != Inches(1) {
				t.Errorf("New%s() Y = %v, want %v", tt.name, tt.shape.Y, Inches(1))
			}
			if tt.shape.CX != Inches(2) {
				t.Errorf("New%s() CX = %v, want %v", tt.name, tt.shape.CX, Inches(2))
			}
			if tt.shape.CY != Inches(2) {
				t.Errorf("New%s() CY = %v, want %v", tt.name, tt.shape.CY, Inches(2))
			}
		})
	}

	// Double check text box has text
	tb := NewTextBox("hello", 1, 1, 2, 2)
	if tb.Text != "hello" {
		t.Errorf("NewTextBox() text = %v, want 'hello'", tb.Text)
	}

	// New macros with specific dimensions
	t.Run("Circle", func(t *testing.T) {
		s := NewCircle(1, 1, 2)
		if s.Type != ShapeTypeEllipse || s.CX != Inches(2) || s.CY != Inches(2) {
			t.Errorf("NewCircle() unexpected properties: %+v", s)
		}
	})

	t.Run("Star", func(t *testing.T) {
		s := NewStar(1, 1, 2)
		if s.Type != ShapeTypeStar5 || s.CX != Inches(2) || s.CY != Inches(2) {
			t.Errorf("NewStar() unexpected properties: %+v", s)
		}
	})

	t.Run("Heart", func(t *testing.T) {
		s := NewHeart(1, 1, 2)
		if s.Type != ShapeTypeHeart || s.CX != Inches(2) || s.CY != Inches(2) {
			t.Errorf("NewHeart() unexpected properties: %+v", s)
		}
	})

	t.Run("FlowChartDocument", func(t *testing.T) {
		s := NewFlowChartDocument(1, 1, 2, 2)
		if s.Type != ShapeTypeFlowChartDocument || s.CX != Inches(2) || s.CY != Inches(2) {
			t.Errorf("NewFlowChartDocument() unexpected properties: %+v", s)
		}
	})

	t.Run("FlowChartData", func(t *testing.T) {
		s := NewFlowChartData(1, 1, 2, 2)
		if s.Type != ShapeTypeFlowChartData || s.CX != Inches(2) || s.CY != Inches(2) {
			t.Errorf("NewFlowChartData() unexpected properties: %+v", s)
		}
	})

	t.Run("Badge", func(t *testing.T) {
		s := NewBadge("NEW", 1, 1, ColorMaterialGreen)
		if s.Type != ShapeTypeRoundedRectangle || s.CX != Inches(1.5) || s.CY != Inches(0.4) {
			t.Errorf("NewBadge() unexpected properties: %+v", s)
		}
		if s.Text != "NEW" || s.Fill == nil || s.Fill.Color != ColorMaterialGreen {
			t.Errorf("NewBadge() missing text or fill: %+v", s)
		}
	})
}
