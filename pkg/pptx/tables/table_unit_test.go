package tables

import (
	"testing"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestTable_Methods(t *testing.T) {
	colWidths := []styling.Length{styling.Inches(1), styling.Inches(2)}
	table := NewTable(colWidths).
		AddRow([]string{"A1", "B1"}).
		AddRow([]string{"A2", "B2"})

	if len(table.Rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(table.Rows))
	}
	if len(table.ColumnWidths) != 2 {
		t.Error("expected 2 columns")
	}

	table = table.Position(styling.Inches(1), styling.Inches(1)).
		Size(styling.Inches(5), styling.Inches(3)).
		WithAltText("Alt").
		WithDecorative(true)

	if table.AltText != "Alt" {
		t.Error("WithAltText failed")
	}
	if !table.IsDecorative {
		t.Error("WithDecorative failed")
	}
	
	table = table.WithRowHeights([]styling.Length{styling.Inches(0.5), styling.Inches(0.5)})
	if len(table.RowHeights) != 2 {
		t.Error("WithRowHeights failed")
	}
}

func TestTable_MergeLogic(t *testing.T) {
	table := NewTable([]styling.Length{styling.Inches(1), styling.Inches(1)})
	// Cell (0,0) spans 2 rows
	table = table.AddStyledRow([]TableCell{
		NewTableCell("A1").WithRowSpan(2),
		NewTableCell("B1"),
	})
	// Row 2: First cell is covered by A1, so we must provide an empty cell or placeholder
	table = table.AddRow([]string{"", "B2"})
	
	spec, err := table.ToTableSpec(1)
	if err != nil { t.Fatalf("ToTableSpec failed: %v", err) }
	if len(spec.StyledRows) != 2 { t.Error("Merge rows failed") }
	// Row 1, Col 0 should be a vertical merge placeholder
	if !spec.StyledRows[1][0].VMerge { t.Error("VMerge flag failed") }
}

func TestTable_Validate_Extended(t *testing.T) {
	tests := []struct {
		name    string
		table   Table
		wantErr bool
	}{
		{"Valid", NewTable([]styling.Length{styling.Inches(1)}).AddRow([]string{"A"}), false},
		{"No Rows", NewTable([]styling.Length{styling.Inches(1)}), true},
		{"No Cols", NewTable([]styling.Length{}), true},
		{"Mismatch Row", NewTable([]styling.Length{styling.Inches(1)}).AddRow([]string{"A", "B"}), true},
		{"Invalid Color", NewTable([]styling.Length{styling.Inches(1)}).AddStyledRow([]TableCell{NewTableCell("X").WithBackgroundColor("invalid")}), true},
		{"ColSpan Overflow", NewTable([]styling.Length{styling.Inches(1)}).AddStyledRow([]TableCell{NewTableCell("A").WithColSpan(2)}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.table.Validate(1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
