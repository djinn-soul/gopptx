package table

import (
	"regexp"
	"strings"
	"testing"
)

const twoRowFrame = `<a:tbl><a:tblGrid><a:gridCol w="100"/><a:gridCol w="100"/></a:tblGrid>` +
	`<a:tr h="548640"><a:tc><a:txBody><a:p/></a:txBody><a:tcPr/></a:tc>` +
	`<a:tc><a:txBody><a:p/></a:txBody><a:tcPr/></a:tc></a:tr>` +
	`<a:tr h="548640"><a:tc><a:txBody><a:p/></a:txBody><a:tcPr/></a:tc>` +
	`<a:tc><a:txBody><a:p/></a:txBody><a:tcPr/></a:tc></a:tr></a:tbl>`

var trOpenPattern = regexp.MustCompile(`<a:tr\b[^>]*>`)

// CT_TableRow requires the h attribute. A row written without it makes
// PowerPoint reject the whole package as unreadable.
func TestEveryInsertedRowCarriesHeight(t *testing.T) {
	cases := map[string]func() ([]byte, error){
		"append with explicit height": func() ([]byte, error) {
			return AddTableRowInFrame([]byte(twoRowFrame), 123456)
		},
		"append with unspecified height": func() ([]byte, error) {
			return AddTableRowInFrame([]byte(twoRowFrame), 0)
		},
		"insert with unspecified height": func() ([]byte, error) {
			return InsertTableRowInFrame([]byte(twoRowFrame), 1, 0)
		},
		"insert at end with unspecified height": func() ([]byte, error) {
			return InsertTableRowInFrame([]byte(twoRowFrame), 2, 0)
		},
	}

	for name, build := range cases {
		t.Run(name, func(t *testing.T) {
			out, err := build()
			if err != nil {
				t.Fatalf("build row: %v", err)
			}
			opens := trOpenPattern.FindAllString(string(out), -1)
			if len(opens) != 3 {
				t.Fatalf("expected 3 rows, got %d: %s", len(opens), out)
			}
			for _, open := range opens {
				if !strings.Contains(open, ` h="`) {
					t.Fatalf("row tag %q has no required h attribute", open)
				}
			}
		})
	}
}

func TestUnspecifiedRowHeightInheritsLastRow(t *testing.T) {
	out, err := AddTableRowInFrame([]byte(twoRowFrame), 0)
	if err != nil {
		t.Fatalf("add row: %v", err)
	}
	if strings.Count(string(out), `<a:tr h="548640">`) != 3 {
		t.Fatalf("new row should inherit 548640: %s", out)
	}
}

func TestRowHeightFallsBackWhenNoRowsCarryHeight(t *testing.T) {
	if got := resolveNewRowHeight(nil, 0); got != defaultTableRowHeightEMU {
		t.Fatalf("resolveNewRowHeight(nil,0) = %d, want %d", got, defaultTableRowHeightEMU)
	}
	if got := resolveNewRowHeight(&XML{}, 0); got != defaultTableRowHeightEMU {
		t.Fatalf("empty table height = %d, want %d", got, defaultTableRowHeightEMU)
	}
}
