package shapes

import "testing"

func TestToXMLTextFrameSpec_MapsOrientationAndColumns(t *testing.T) {
	rotation := 45.0
	spec := toXMLTextFrameSpec(&TextFrame{
		Orientation: "vert270",
		Columns:     2,
		RotationDeg: &rotation,
	})

	if spec == nil {
		t.Fatal("expected text frame spec, got nil")
	}
	if spec.Orientation != "vert270" {
		t.Fatalf("expected orientation vert270, got %q", spec.Orientation)
	}
	if spec.NumCol != 2 {
		t.Fatalf("expected numCol 2, got %d", spec.NumCol)
	}
	if spec.Rotation == nil || *spec.Rotation != 2700000 {
		t.Fatalf("expected OOXML rotation 2700000, got %+v", spec.Rotation)
	}
}
