package export

import (
	"fmt"
	"testing"
)

// TestPaintListSortsByTreePosition pins the rule the renderer exists to follow:
// PowerPoint paints a slide in shape-tree order, whatever kind each element is.
func TestPaintListSortsByTreePosition(t *testing.T) {
	t.Parallel()

	var painted []string
	list := &paintList{}
	list.add(3, func() error { painted = append(painted, "card"); return nil })
	list.add(0, func() error { painted = append(painted, "title"); return nil })
	list.add(2, func() error { painted = append(painted, "chart"); return nil })
	list.add(1, func() error { painted = append(painted, "table"); return nil })

	if errs := list.drawAll(); len(errs) != 0 {
		t.Fatalf("drawAll errors: %v", errs)
	}
	want := []string{"title", "table", "chart", "card"}
	if fmt.Sprint(painted) != fmt.Sprint(want) {
		t.Fatalf("paint order = %v, want %v", painted, want)
	}
}

// An element whose tree position was never recovered paints above the ones that
// have one, in the order it was added — the layering the renderer gives a slide
// built in memory.
func TestPaintListKeepsUnorderedElementsOnTop(t *testing.T) {
	t.Parallel()

	var painted []string
	list := &paintList{}
	list.add(unknownZ, func() error { painted = append(painted, "first-added"); return nil })
	list.add(5, func() error { painted = append(painted, "tree-element"); return nil })
	list.add(unknownZ, func() error { painted = append(painted, "second-added"); return nil })

	if errs := list.drawAll(); len(errs) != 0 {
		t.Fatalf("drawAll errors: %v", errs)
	}
	want := []string{"tree-element", "first-added", "second-added"}
	if fmt.Sprint(painted) != fmt.Sprint(want) {
		t.Fatalf("paint order = %v, want %v", painted, want)
	}
}

func TestSlidePaintOrderReportsUnknownPositions(t *testing.T) {
	t.Parallel()

	order := newSlidePaintOrder()
	if got := order.chartZ(chartKindBar); got != unknownZ {
		t.Errorf("chart z = %d, want unknown", got)
	}
	if got := order.tableZ(0); got != unknownZ {
		t.Errorf("table z = %d, want unknown", got)
	}
	if got := order.treeZ(4); got != unknownZ {
		t.Errorf("shape z = %d on a slide with no tree order, want unknown", got)
	}

	order.known = true
	order.charts[chartKindBar] = 7
	order.tables = []int{2}
	if got := order.chartZ(chartKindBar); got != 7 {
		t.Errorf("chart z = %d, want 7", got)
	}
	if got := order.tableZ(0); got != 2 {
		t.Errorf("table z = %d, want 2", got)
	}
	if got := order.treeZ(0); got != 0 {
		t.Errorf("shape z = %d, want 0: index zero is the backmost element", got)
	}
}
