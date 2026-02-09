package pptx_test

import (
	"fmt"
	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func ExampleCenter() {
	// Calculate coordinates to center a 4"x3" box on a standard slide
	x, y := pptx.Center(pptx.Inches(4), pptx.Inches(3))
	fmt.Printf("X: %d, Y: %d\n", x, y)
	// Output: X: 2743200, Y: 2057400
}

func ExampleGrid() {
	// Create a 2x2 grid with 0.5" margin
	margin := pptx.Inches(0.5)
	boxes, _ := pptx.Grid(2, 2, margin)

	for i, box := range boxes {
		fmt.Printf("Box %d: X=%d, Y=%d, CX=%d, CY=%d\n", i, box.X, box.Y, box.CX, box.CY)
	}
	// Output:
	// Box 0: X=0, Y=0, CX=4343400, CY=3200400
	// Box 1: X=4800600, Y=0, CX=4343400, CY=3200400
	// Box 2: X=0, Y=3657600, CX=4343400, CY=3200400
	// Box 3: X=4800600, Y=3657600, CX=4343400, CY=3200400
}

func ExampleStack() {
	// Stack three boxes vertically with a 0.2" gap
	elements := []pptx.Size{
		{CX: pptx.Inches(2), CY: pptx.Inches(1)},
		{CX: pptx.Inches(2), CY: pptx.Inches(1)},
		{CX: pptx.Inches(2), CY: pptx.Inches(1)},
	}
	start := pptx.Point{X: pptx.Inches(1), Y: pptx.Inches(1)}
	gap := pptx.Inches(0.2)

	points, _ := pptx.Stack(pptx.OrientationVertical, start, gap, elements...)

	for i, p := range points {
		fmt.Printf("Element %d: X=%d, Y=%d\n", i, p.X, p.Y)
	}
	// Output:
	// Element 0: X=914400, Y=914400
	// Element 1: X=914400, Y=2011680
	// Element 2: X=914400, Y=3108960
}

func ExampleDistribute() {
	// Distribute 3 elements horizontally across the full slide
	bounds := pptx.Box{X: 0, Y: 0, CX: pptx.SlideWidth, CY: pptx.SlideHeight}
	count := 3
	elementWidth := pptx.Inches(2)

	coords, _ := pptx.Distribute(pptx.OrientationHorizontal, bounds, count, elementWidth)

	for i, x := range coords {
		fmt.Printf("Element %d X: %d\n", i, x)
	}
	// Output:
	// Element 0 X: 0
	// Element 1 X: 3657600
	// Element 2 X: 7315200
}
