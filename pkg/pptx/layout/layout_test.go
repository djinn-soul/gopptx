package layout

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestCenter(t *testing.T) {
	// Slide is 9144000 x 6858000
	x, y := Center(3657600, 1828800)

	expectedX := (SlideWidth - 3657600) / 2
	expectedY := (SlideHeight - 1828800) / 2

	if x != expectedX || y != expectedY {
		t.Errorf("Center(%d, %d) = (%d, %d); expected (%d, %d)", 3657600, 1828800, x, y, expectedX, expectedY)
	}
}

func TestGridInBox(t *testing.T) {
	bounds := common.Box{X: 0, Y: 0, CX: 1000, CY: 1000}
	rows, cols := 2, 2
	margin := styling.Emu(100)

	boxes, err := GridInBox(rows, cols, margin, bounds)
	if err != nil {
		t.Fatalf("GridInBox failed: %v", err)
	}

	if len(boxes) != 4 {
		t.Errorf("expected 4 boxes, got %d", len(boxes))
	}

	// Element size should be (1000 - 100) / 2 = 450
	expectedCX := styling.Emu(450)
	if boxes[0].CX != expectedCX {
		t.Errorf("expected CX %d, got %d", expectedCX, boxes[0].CX)
	}

	// Check second element in first row (c=1, r=0)
	if boxes[1].X != expectedCX+margin {
		t.Errorf("expected X %d for second box, got %d", expectedCX+margin, boxes[1].X)
	}
}

func TestStack(t *testing.T) {
	elements := []common.Size{{CX: 100, CY: 100}, {CX: 200, CY: 200}}
	start := common.Point{X: 0, Y: 0}
	gap := styling.Emu(50)

	// Horizontal stack
	points, err := Stack(OrientationHorizontal, start, gap, elements...)
	if err != nil {
		t.Fatalf("Horizontal Stack failed: %v", err)
	}

	if points[1].X != 100+50 {
		t.Errorf("expected X 150 for second point, got %d", points[1].X)
	}

	// Vertical stack
	points, err = Stack(OrientationVertical, start, gap, elements...)
	if err != nil {
		t.Fatalf("Vertical Stack failed: %v", err)
	}

	if points[1].Y != 100+50 {
		t.Errorf("expected Y 150 for second point, got %d", points[1].Y)
	}
}

func TestDistribute(t *testing.T) {
	bounds := common.Box{X: 0, Y: 0, CX: 1000, CY: 1000}
	count := 3
	elSize := styling.Emu(200)

	// Distribute horizontally
	coords, err := DistributeUniform(OrientationHorizontal, bounds, count, elSize)
	if err != nil {
		t.Fatalf("Distribute failed: %v", err)
	}

	// Total elements = 600. Remaining = 400. Gaps = 2. Gap size = 200.
	// Coords: 0, 400, 800
	if coords[1] != 400 {
		t.Errorf("expected second coord 400, got %d", coords[1])
	}
	if coords[2] != 800 {
		t.Errorf("expected third coord 800, got %d", coords[2])
	}
}

func TestCenterInBox(t *testing.T) {
	bounds := common.Box{X: 100, Y: 100, CX: 1000, CY: 1000}
	x, y := CenterInBox(400, 200, bounds)

	expectedX := styling.Emu(100 + (1000-400)/2) // 100 + 300 = 400
	expectedY := styling.Emu(100 + (1000-200)/2) // 100 + 400 = 500

	if x != expectedX || y != expectedY {
		t.Errorf("CenterInBox(%d, %d, bounds) = (%d, %d); expected (%d, %d)", 400, 200, x, y, expectedX, expectedY)
	}
}

func TestDistributeVertical(t *testing.T) {
	bounds := common.Box{X: 0, Y: 0, CX: 1000, CY: 1000}
	count := 3
	elSize := styling.Emu(200)

	// Distribute vertically
	coords, err := DistributeUniform(OrientationVertical, bounds, count, elSize)
	if err != nil {
		t.Fatalf("Vertical Distribute failed: %v", err)
	}

	if coords[1] != 400 {
		t.Errorf("expected second coord 400, got %d", coords[1])
	}
}

func TestDistributeSingleElement(t *testing.T) {
	bounds := common.Box{X: 0, Y: 0, CX: 1000, CY: 1000}
	elSize := styling.Emu(200)

	// Horizontal
	coords, err := DistributeUniform(OrientationHorizontal, bounds, 1, elSize)
	if err != nil {
		t.Fatalf("Horizontal Single Distribute failed: %v", err)
	}
	if coords[0] != 400 {
		t.Errorf("expected coord 400, got %d", coords[0])
	}

	// Vertical
	coords, err = DistributeUniform(OrientationVertical, bounds, 1, elSize)
	if err != nil {
		t.Fatalf("Vertical Single Distribute failed: %v", err)
	}
	if coords[0] != 400 {
		t.Errorf("expected coord 400, got %d", coords[0])
	}
}

func TestDistributeNonUniform(t *testing.T) {
	bounds := common.Box{X: 0, Y: 0, CX: 1000, CY: 1000}
	sizes := []styling.Length{styling.Emu(100), styling.Emu(200), styling.Emu(300)}

	// Horizontal: total size = 600. Remaining = 400. Gaps = 2. Gap size = 200.
	// Coords: 0, 100+200=300, 300+200+200=700
	coords, err := DistributeNonUniform(OrientationHorizontal, bounds, sizes)
	if err != nil {
		t.Fatalf("DistributeNonUniform failed: %v", err)
	}
	if coords[0] != 0 {
		t.Errorf("expected first coord 0, got %d", coords[0])
	}
	if coords[1] != 300 {
		t.Errorf("expected second coord 300, got %d", coords[1])
	}
	if coords[2] != 700 {
		t.Errorf("expected third coord 700, got %d", coords[2])
	}

	// Single element: centered.
	single, err := DistributeNonUniform(
		OrientationVertical, bounds, []styling.Length{styling.Emu(200)},
	)
	if err != nil {
		t.Fatalf("single DistributeNonUniform failed: %v", err)
	}
	if single[0] != 400 {
		t.Errorf("expected single coord 400, got %d", single[0])
	}
}

func TestDistributeNonUniformErrors(t *testing.T) {
	bounds := common.Box{X: 0, Y: 0, CX: 1000, CY: 1000}

	if _, err := DistributeNonUniform(OrientationHorizontal, bounds, nil); err == nil {
		t.Error("expected error for empty sizes")
	}
	if _, err := DistributeNonUniform(
		"diagonal", bounds, []styling.Length{100, 200},
	); err == nil {
		t.Error("expected error for invalid orientation")
	}
	if _, err := DistributeNonUniform(
		OrientationHorizontal, bounds, []styling.Length{600, 600},
	); err == nil {
		t.Error("expected error for elements exceeding space")
	}
	if _, err := DistributeNonUniform(
		OrientationHorizontal, bounds, []styling.Length{100, -50},
	); err == nil {
		t.Error("expected error for negative size")
	}
}

func TestLayoutHelpersErrors(t *testing.T) {
	_, err := Grid(0, 1, 0)
	if err == nil {
		t.Error("expected error for 0 rows in Grid")
	}

	_, err = Stack("invalid", common.Point{X: 0, Y: 0}, 0, common.Size{CX: 10, CY: 10})
	if err == nil {
		t.Error("expected error for invalid orientation in Stack")
	}

	_, err = DistributeUniform(OrientationHorizontal, common.Box{X: 0, Y: 0, CX: 100, CY: 100}, 2, 60)
	if err == nil {
		t.Error("expected error when elements exceed space in Distribute")
	}
}
