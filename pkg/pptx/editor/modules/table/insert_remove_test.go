package table

import (
	"strings"
	"testing"
)

func TestInsertTableRowInFrame(t *testing.T) {
	_, frame := buildTestSlideAndFrame(t, 10)

	// Insert at index 0 (prepend)
	updated, err := InsertTableRowInFrame(frame, 0, 0)
	if err != nil {
		t.Fatalf("InsertTableRowInFrame(0) failed: %v", err)
	}
	parsed, _ := ParseTable(updated)
	rows, cols := Dimensions(parsed)
	if rows != 3 {
		t.Errorf("expected 3 rows after prepend insert, got %d", rows)
	}
	if cols != 2 {
		t.Errorf("expected 2 cols unchanged, got %d", cols)
	}

	// Insert at index 1 (middle)
	updated2, err := InsertTableRowInFrame(frame, 1, 914400)
	if err != nil {
		t.Fatalf("InsertTableRowInFrame(1) failed: %v", err)
	}
	if !strings.Contains(string(updated2), `h="914400"`) {
		t.Errorf("expected h=914400 in inserted row")
	}
	parsed2, _ := ParseTable(updated2)
	rows2, _ := Dimensions(parsed2)
	if rows2 != 3 {
		t.Errorf("expected 3 rows after middle insert, got %d", rows2)
	}

	// Insert at end (equivalent to append)
	updated3, err := InsertTableRowInFrame(frame, 2, 0)
	if err != nil {
		t.Fatalf("InsertTableRowInFrame(end) failed: %v", err)
	}
	parsed3, _ := ParseTable(updated3)
	rows3, _ := Dimensions(parsed3)
	if rows3 != 3 {
		t.Errorf("expected 3 rows after end insert, got %d", rows3)
	}

	// Out-of-range index
	if _, err = InsertTableRowInFrame(frame, 99, 0); err == nil {
		t.Fatal("expected error for out-of-range index")
	}
}

func TestRemoveTableRowInFrame(t *testing.T) {
	_, frame := buildTestSlideAndFrame(t, 11)

	updated, err := RemoveTableRowInFrame(frame, 0)
	if err != nil {
		t.Fatalf("RemoveTableRowInFrame(0) failed: %v", err)
	}
	parsed, _ := ParseTable(updated)
	rows, _ := Dimensions(parsed)
	if rows != 1 {
		t.Errorf("expected 1 row after removing row 0, got %d", rows)
	}

	// Remove last row
	updated2, err := RemoveTableRowInFrame(frame, 1)
	if err != nil {
		t.Fatalf("RemoveTableRowInFrame(1) failed: %v", err)
	}
	parsed2, _ := ParseTable(updated2)
	rows2, _ := Dimensions(parsed2)
	if rows2 != 1 {
		t.Errorf("expected 1 row after removing last row, got %d", rows2)
	}

	// Out-of-range
	if _, err = RemoveTableRowInFrame(frame, 99); err == nil {
		t.Fatal("expected error for out-of-range row index")
	}
}

func TestInsertTableColumnInFrame(t *testing.T) {
	_, frame := buildTestSlideAndFrame(t, 12)

	// Insert before column 0
	updated, err := InsertTableColumnInFrame(frame, 0, 914400)
	if err != nil {
		t.Fatalf("InsertTableColumnInFrame(0) failed: %v", err)
	}
	parsed, _ := ParseTable(updated)
	rows, cols := Dimensions(parsed)
	if cols != 3 {
		t.Errorf("expected 3 cols after insert at 0, got %d", cols)
	}
	if rows != 2 {
		t.Errorf("expected 2 rows unchanged, got %d", rows)
	}

	// Insert in the middle
	updated2, err := InsertTableColumnInFrame(frame, 1, 457200)
	if err != nil {
		t.Fatalf("InsertTableColumnInFrame(1) failed: %v", err)
	}
	parsed2, _ := ParseTable(updated2)
	_, cols2 := Dimensions(parsed2)
	if cols2 != 3 {
		t.Errorf("expected 3 cols after middle insert, got %d", cols2)
	}

	// Zero width must fail
	if _, err = InsertTableColumnInFrame(frame, 0, 0); err == nil {
		t.Fatal("expected error for zero width")
	}

	// Out-of-range
	if _, err = InsertTableColumnInFrame(frame, 99, 914400); err == nil {
		t.Fatal("expected error for out-of-range col index")
	}
}

func TestRemoveTableColumnInFrame(t *testing.T) {
	_, frame := buildTestSlideAndFrame(t, 13)

	// Remove first column
	updated, err := RemoveTableColumnInFrame(frame, 0)
	if err != nil {
		t.Fatalf("RemoveTableColumnInFrame(0) failed: %v", err)
	}
	parsed, _ := ParseTable(updated)
	_, cols := Dimensions(parsed)
	if cols != 1 {
		t.Errorf("expected 1 col after removing col 0, got %d", cols)
	}
	rows, _ := Dimensions(parsed)
	if rows != 2 {
		t.Errorf("expected 2 rows unchanged, got %d", rows)
	}

	// Remove last column
	updated2, err := RemoveTableColumnInFrame(frame, 1)
	if err != nil {
		t.Fatalf("RemoveTableColumnInFrame(1) failed: %v", err)
	}
	parsed2, _ := ParseTable(updated2)
	_, cols2 := Dimensions(parsed2)
	if cols2 != 1 {
		t.Errorf("expected 1 col after removing col 1, got %d", cols2)
	}

	// Cannot remove only remaining column
	singleColFrame, _ := RemoveTableColumnInFrame(frame, 0)
	if _, err = RemoveTableColumnInFrame(singleColFrame, 0); err == nil {
		t.Fatal("expected error when removing last column")
	}

	// Out-of-range
	if _, err = RemoveTableColumnInFrame(frame, 99); err == nil {
		t.Fatal("expected error for out-of-range col index")
	}
}
