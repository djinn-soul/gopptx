package shapes

import (
	"testing"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func TestShapes_Creation(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		s := NewShape(ShapeTypeRectangle, 0, 0, 1, 1).
			WithFill(NewShapeFill("FF0000")).
			WithLine(ShapeLine{Color: "000000", Width: 1000}).
			WithRotation(45).
			WithAltText("Alt").
			WithDecorative(true).
			WithName("Rect")
		
		if s.Type != ShapeTypeRectangle || s.Fill.Color != "FF0000" { t.Error("Basic props failed") }
	})

	t.Run("Gradient", func(t *testing.T) {
		s := NewRectangle(0, 0, 1, 1).
			WithGradientFill(NewShapeGradientFill(ShapeGradientTypeLinear, []ShapeGradientStop{
				NewShapeGradientStop(0, "000000"),
				NewShapeGradientStop(100, "FFFFFF"),
			}))
		if s.GradientFill.Type != ShapeGradientTypeLinear { t.Error("Gradient failed") }
	})
	
	t.Run("Convenience", func(t *testing.T) {
		_ = NewRectangle(0, 0, 1, 1)
		_ = NewEllipse(0, 0, 1, 1)
		_ = NewTriangle(0, 0, 1, 1)
		_ = NewRightArrow(0, 0, 1, 1)
		_ = NewLeftArrow(0, 0, 1, 1)
		_ = NewTextBox("text", 0, 0, 1, 1)
		_ = NewShapeLine("FF0000", 1000)
	})
}

func TestShapes_Connectors(t *testing.T) {
	c := NewStraightConnector(0, 0, 1, 1).
		WithArrows(ArrowTypeStealth, ArrowTypeTriangle).
		WithLabel("L").
		ConnectStart(1, ConnectionSiteTop).
		ConnectEnd(3, ConnectionSiteBottom)
	
	if c.StartArrow != ArrowTypeStealth || c.Label != "L" { t.Error("Connector props failed") }
	
	c2 := NewElbowConnector(0, 0, 1, 1).ConnectStartAuto(1).ConnectEndAuto(2)
	if c2.Type != ConnectorTypeElbow { t.Error("Elbow failed") }
	
	c3 := NewCurvedConnector(0, 0, 1, 1).AutoReroute(nil)
	if c3.Type != ConnectorTypeCurved { t.Error("Curved failed") }
}

func TestShapes_Placeholders(t *testing.T) {
	p := &Placeholder{Type: PlaceholderTypeBody, Index: 1}
	content := p.InsertText("Hello")
	if content.Text != "Hello" || content.Type != "body" { t.Error("InsertText failed") }
	
	img := p.InsertPicture("test.png")
	if img.Path != "test.png" { t.Error("InsertPicture failed") }
	
	tbl := p.InsertTable(tables.Table{})
	if tbl.Table == nil { t.Error("InsertTable failed") }
	
	img2 := p.InsertPictureFromBytes([]byte("fake"), "png")
	if img2.Format != "png" { t.Error("InsertPictureFromBytes failed") }
	
	p2 := p.InsertPictureToSlide(shapesImage(img))
	if p2.Image == nil { t.Error("InsertPictureToSlide failed") }
}

func shapesImage(i Image) Image { return i }

func TestShapes_Validate(t *testing.T) {
	t.Run("Shape", func(t *testing.T) {
		s := NewRectangle(0, 0, -1, 1)
		if err := s.Validate(1, 1); err == nil { t.Error("expected error for negative width") }
		
		s = NewRectangle(0, 0, 1, 1).WithFill(NewShapeFill("invalid"))
		if err := s.Validate(1, 1); err == nil { t.Error("expected error for invalid color") }
	})
	
	t.Run("Gradient", func(t *testing.T) {
		g := NewShapeGradientFill("invalid", nil)
		if err := g.Validate(); err == nil { t.Error("expected error for invalid type") }
	})
}
