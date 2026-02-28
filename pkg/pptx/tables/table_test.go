package tables

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestNewTable(t *testing.T) {
	widths := []styling.Length{styling.Inches(1), styling.Inches(2)}
	table := NewTable(widths)

	if len(table.ColumnWidths) != 2 {
		t.Errorf("expected 2 columns, got %d", len(table.ColumnWidths))
	}
	if table.X.Emu() != defaultTableX {
		t.Errorf("expected default X %d, got %d", defaultTableX, table.X.Emu())
	}
}

func TestTable_WithAltText(t *testing.T) {
	table := NewTable([]styling.Length{styling.Inches(1)}).WithAltText("test alt text")
	if table.AltText != "test alt text" {
		t.Errorf("expected alt text 'test alt text', got '%s'", table.AltText)
	}
}

func TestTable_WithDecorative(t *testing.T) {
	table := NewTable([]styling.Length{styling.Inches(1)}).WithDecorative(true)
	if !table.IsDecorative {
		t.Error("expected table to be decorative")
	}
}

func TestTable_AddRow(t *testing.T) {
	table := NewTable([]styling.Length{styling.Inches(1)})
	table = table.AddRow([]string{"cell1"})

	if len(table.Rows) != 1 {
		t.Errorf("expected 1 row, got %d", len(table.Rows))
	}
	if table.Rows[0][0] != "cell1" {
		t.Errorf("expected cell content 'cell1', got '%s'", table.Rows[0][0])
	}
}

func TestTable_AddStyledRow(t *testing.T) {
	table := NewTable([]styling.Length{styling.Inches(1)})
	cell := NewTableCell("styled").WithBold(true)
	table = table.AddStyledRow([]TableCell{cell})

	if len(table.StyledRows) != 1 {
		t.Errorf("expected 1 styled row, got %d", len(table.StyledRows))
	}
	if table.StyledRows[0][0].Text != "styled" {
		t.Errorf("expected cell text 'styled', got '%s'", table.StyledRows[0][0].Text)
	}
	if !table.StyledRows[0][0].Bold {
		t.Error("expected cell to be bold")
	}
}

func TestTable_PositionAndSize(t *testing.T) {
	table := NewTable([]styling.Length{styling.Inches(1)})
	table = table.Position(styling.Emu(100), styling.Emu(200))
	table = table.Size(styling.Emu(300), styling.Emu(400))

	if table.X.Emu() != 100 || table.Y.Emu() != 200 {
		t.Errorf("expected position (100, 200), got (%d, %d)", table.X.Emu(), table.Y.Emu())
	}
	if table.CX.Emu() != 300 || table.CY.Emu() != 400 {
		t.Errorf("expected size (300, 400), got (%d, %d)", table.CX.Emu(), table.CY.Emu())
	}
}

func TestTable_WithRowHeights(t *testing.T) {
	table := NewTable([]styling.Length{styling.Inches(1)})
	heights := []styling.Length{styling.Inches(0.5)}
	table = table.WithRowHeights(heights)

	if len(table.RowHeights) != 1 {
		t.Errorf("expected 1 row height, got %d", len(table.RowHeights))
	}
	if table.RowHeights[0].Emu() != heights[0].Emu() {
		t.Errorf("expected height %d, got %d", heights[0].Emu(), table.RowHeights[0].Emu())
	}

	table = table.WithRowHeights(nil)
	if table.RowHeights != nil {
		t.Error("expected row heights to be nil")
	}
}

func TestTable_ToTableSpec(t *testing.T) {
	table := NewTable([]styling.Length{styling.Inches(1)})
	table = table.AddRow([]string{"test"})

	spec, err := table.ToTableSpec(1)
	if err != nil {
		t.Fatalf("ToTableSpec failed: %v", err)
	}

	if spec.Rows[0][0] != "test" {
		t.Errorf("expected spec row content 'test', got '%s'", spec.Rows[0][0])
	}
	if len(spec.ColumnWidths) != 1 {
		t.Errorf("expected 1 column width, got %d", len(spec.ColumnWidths))
	}
}

func TestTable_Validate(t *testing.T) {
	table := NewTable([]styling.Length{styling.Inches(1)})
	table = table.AddRow([]string{"test"})

	err := table.Validate(1)
	if err != nil {
		t.Errorf("Validate failed: %v", err)
	}
}
