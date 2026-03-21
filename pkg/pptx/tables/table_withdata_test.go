package tables

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestTable_WithData_PopulatesRows(t *testing.T) {
	tbl := NewTable([]styling.Length{styling.Inches(1), styling.Inches(1)})

	data := [][]string{
		{"Item", "Qty"},
		{"Widgets", "50"},
		{"Gadgets", "30"},
	}

	result := tbl.WithData(data)

	if len(result.Rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(result.Rows))
	}
	if len(result.Rows[0]) != 2 {
		t.Errorf("Expected 2 columns in first row, got %d", len(result.Rows[0]))
	}
	if result.Rows[0][0] != "Item" {
		t.Errorf("Expected 'Item', got '%s'", result.Rows[0][0])
	}
	if result.Rows[1][0] != "Widgets" {
		t.Errorf("Expected 'Widgets', got '%s'", result.Rows[1][0])
	}
}

func TestTable_WithData_ClearsPreviousRows(t *testing.T) {
	tbl := NewTable([]styling.Length{styling.Inches(1), styling.Inches(1)})

	// Add some rows
	tbl = tbl.AddRow([]string{"Old1", "Old2"})
	tbl = tbl.AddRow([]string{"Old3", "Old4"})

	if len(tbl.Rows) != 2 {
		t.Errorf("Expected 2 rows before WithData, got %d", len(tbl.Rows))
	}

	// Replace with new data
	data := [][]string{
		{"New1", "New2"},
	}
	tbl = tbl.WithData(data)

	if len(tbl.Rows) != 1 {
		t.Errorf("Expected 1 row after WithData, got %d", len(tbl.Rows))
	}
	if tbl.Rows[0][0] != "New1" {
		t.Errorf("Expected 'New1', got '%s'", tbl.Rows[0][0])
	}
}

func TestTable_WithStyledData_PopulatesRows(t *testing.T) {
	tbl := NewTable([]styling.Length{styling.Inches(1), styling.Inches(1)})

	data := [][]TableCell{
		{
			NewTableCell("Header1").WithBold(true),
			NewTableCell("Header2").WithBold(true),
		},
		{
			NewTableCell("Data1"),
			NewTableCell("Data2"),
		},
	}

	result := tbl.WithStyledData(data)

	if len(result.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(result.Rows))
	}
	if !result.StyledRows[0][0].Bold {
		t.Errorf("Expected first cell to be bold")
	}
	if result.Rows[0][0] != "Header1" {
		t.Errorf("Expected 'Header1', got '%s'", result.Rows[0][0])
	}
}

func TestTable_WithData_Chaining(t *testing.T) {
	// Test that WithData works in a fluent chain
	tbl := NewTable([]styling.Length{styling.Inches(1), styling.Inches(1), styling.Inches(1)}).
		WithData([][]string{
			{"A", "B", "C"},
			{"1", "2", "3"},
		}).
		Position(styling.Inches(1), styling.Inches(2)).
		Size(styling.Inches(8), styling.Inches(4))

	if len(tbl.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(tbl.Rows))
	}
	if tbl.X.Emu() != styling.Inches(1).Emu() {
		t.Errorf("Position not set correctly")
	}
}

func TestTable_WithData_EmptyData(t *testing.T) {
	tbl := NewTable([]styling.Length{styling.Inches(1)})

	// WithData with empty slice should clear rows
	result := tbl.WithData([][]string{})

	if len(result.Rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(result.Rows))
	}
}
