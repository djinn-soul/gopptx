package tables

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestTableValidation_EdgeCases(t *testing.T) {
	// Negative position
	t1 := NewTable([]styling.Length{styling.Inches(1)}).Position(styling.Emu(-1), 0).AddRow([]string{"A"})
	if err := t1.Validate(1); err == nil || !strings.Contains(err.Error(), "position cannot be negative") {
		t.Errorf("Expected negative position error, got %v", err)
	}

	// Zero size
	t2 := NewTable([]styling.Length{styling.Inches(1)}).Size(0, 100).AddRow([]string{"A"})
	if err := t2.Validate(1); err == nil || !strings.Contains(err.Error(), "size must be > 0") {
		t.Errorf("Expected zero size error, got %v", err)
	}

	// No columns
	t3 := Table{X: 0, Y: 0, CX: 100, CY: 100}
	if err := validateTableColumns(t3, 1); err == nil || !strings.Contains(err.Error(), "must define at least one column") {
		t.Errorf("Expected no columns error, got %v", err)
	}

	// Zero column width
	t4 := NewTable([]styling.Length{0}).AddRow([]string{"A"})
	if err := t4.Validate(1); err == nil || !strings.Contains(err.Error(), "width must be > 0") {
		t.Errorf("Expected zero width error, got %v", err)
	}

	// Mismatched row heights
	t5 := NewTable([]styling.Length{styling.Inches(1)}).AddRow([]string{"A"}).WithRowHeights([]styling.Length{1, 2})
	if err := t5.Validate(1); err == nil || !strings.Contains(err.Error(), "match row count") {
		t.Errorf("Expected row height mismatch error, got %v", err)
	}

	// Zero row height
	t6 := NewTable([]styling.Length{styling.Inches(1)}).AddRow([]string{"A"}).WithRowHeights([]styling.Length{0})
	if err := t6.Validate(1); err == nil || !strings.Contains(err.Error(), "height must be > 0") {
		t.Errorf("Expected zero row height error, got %v", err)
	}

	// Row length mismatch
	t7 := NewTable([]styling.Length{styling.Inches(1), styling.Inches(1)}).AddRow([]string{"A"})
	if err := t7.Validate(1); err == nil || !strings.Contains(err.Error(), "has 1 cells; expected 2") {
		t.Errorf("Expected row length mismatch error, got %v", err)
	}
}

func TestTableCellValidation_EdgeCases(t *testing.T) {
	// Invalid Span
	c1 := NewTableCell("A")
	c1.RowSpan = 0
	if err := validateTableCell(c1, 1, 1, 1); err == nil || !strings.Contains(err.Error(), "row span must be >= 1") {
		t.Errorf("Expected row span error, got %v", err)
	}
	c1.RowSpan = 1
	c1.ColSpan = 0
	if err := validateTableCell(c1, 1, 1, 1); err == nil || !strings.Contains(err.Error(), "col span must be >= 1") {
		t.Errorf("Expected col span error, got %v", err)
	}

	// Invalid Color
	c2 := NewTableCell("A").WithBackgroundColor("invalid")
	if err := validateTableCell(c2, 1, 1, 1); err == nil || !strings.Contains(err.Error(), "background color must be 6-digit RGB hex") {
		t.Errorf("Expected background color error, got %v", err)
	}

	// Invalid Align
	c3 := NewTableCell("A").WithAlign("top") // top is for valign
	if err := validateTableCell(c3, 1, 1, 1); err == nil || !strings.Contains(err.Error(), "align must be one of") {
		t.Errorf("Expected align error, got %v", err)
	}

	// Invalid VAlign
	c4 := NewTableCell("A").WithVAlign("left")
	if err := validateTableCell(c4, 1, 1, 1); err == nil || !strings.Contains(err.Error(), "valign must be one of") {
		t.Errorf("Expected valign error, got %v", err)
	}

	// Invalid Margins
	neg := -1.0
	c5 := NewTableCell("A")
	c5.MarginLeftPt = &neg
	if err := validateTableCell(c5, 1, 1, 1); err == nil || !strings.Contains(err.Error(), "left margin must be >= 0") {
		t.Errorf("Expected margin error, got %v", err)
	}
}

func TestTableBorderValidation_EdgeCases(t *testing.T) {
	// Negative border width
	c1 := NewTableCell("A").WithBorder(-1, "000000")
	if err := validateTableCell(c1, 1, 1, 1); err == nil || !strings.Contains(err.Error(), "width must be >= 0") {
		t.Errorf("Expected border width error, got %v", err)
	}

	// Border color without width
	c2 := NewTableCell("A").WithBorder(0, "FF0000")
	if err := validateTableCell(c2, 1, 1, 1); err == nil || !strings.Contains(err.Error(), "width must be > 0 when") {
		t.Errorf("Expected border width requirement error, got %v", err)
	}

	// Invalid border dash
	c3 := NewTableCell("A").WithBorderStyle(1, "000000", "zigzag")
	if err := validateTableCell(c3, 1, 1, 1); err == nil || !strings.Contains(err.Error(), "dash style must be one of") {
		t.Errorf("Expected border dash error, got %v", err)
	}

	// Explicit sides
	c4 := NewTableCell("A").WithTopBorder(1, "bad")
	if err := validateTableCell(c4, 1, 1, 1); err == nil || !strings.Contains(err.Error(), "top border color must be 6-digit RGB hex") {
		t.Errorf("Expected explicit top border error, got %v", err)
	}
}

func TestTableNormalizationHelpers(t *testing.T) {
	if NormalizeTableBorderDash("  DASH-DOT  ") != TableBorderDashDashDot {
		t.Error("NormalizeTableBorderDash failed for dash-dot")
	}
	if NormalizeTableBorderDash("long_dash") != TableBorderDashLongDash {
		t.Error("NormalizeTableBorderDash failed for long_dash")
	}
	if NormalizeTableBorderDash("unknown") != "unknown" {
		t.Error("NormalizeTableBorderDash should return unknown as is")
	}
}
