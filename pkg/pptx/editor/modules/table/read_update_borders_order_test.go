package table

import (
	"strings"
	"testing"
)

func filledCellFrame() []byte {
	return []byte(`<a:tbl><a:tblGrid><a:gridCol w="100"/></a:tblGrid>` +
		`<a:tr h="100"><a:tc><a:txBody><a:p/></a:txBody>` +
		`<a:tcPr><a:solidFill><a:srgbClr val="FFE699"/></a:solidFill></a:tcPr>` +
		`</a:tc></a:tr></a:tbl>`)
}

// CT_TableCellProperties orders the border lines before the cell fill; appending
// a border after an existing <a:solidFill> makes PowerPoint reject the table.
func TestUpdateTableCellBordersKeepsSchemaOrder(t *testing.T) {
	updated, err := UpdateTableCellBordersInFrame(
		filledCellFrame(), 0, 0, "left",
		&CellBorderSideUpdate{Width: 38100, Color: "C00000", Dash: "solid"},
	)
	if err != nil {
		t.Fatalf("update border: %v", err)
	}

	out := string(updated)
	border := strings.Index(out, "<a:lnL")
	fill := strings.Index(out, `<a:solidFill><a:srgbClr val="FFE699"/>`)
	if border < 0 || fill < 0 {
		t.Fatalf("missing border or fill: %s", out)
	}
	if border > fill {
		t.Fatalf("border must precede the cell fill: %s", out)
	}
}

func TestUpdateTableCellBordersOrdersSidesAmongThemselves(t *testing.T) {
	frame := filledCellFrame()
	var err error
	for _, side := range []string{"bottom", "left", "top"} {
		frame, err = UpdateTableCellBordersInFrame(
			frame, 0, 0, side,
			&CellBorderSideUpdate{Width: 12700, Color: "000000", Dash: "solid"},
		)
		if err != nil {
			t.Fatalf("update %s border: %v", side, err)
		}
	}

	out := string(frame)
	left := strings.Index(out, "<a:lnL")
	top := strings.Index(out, "<a:lnT")
	bottom := strings.Index(out, "<a:lnB")
	fill := strings.Index(out, `<a:srgbClr val="FFE699"/>`)
	if left < 0 || top < 0 || bottom < 0 || fill < 0 {
		t.Fatalf("missing elements: %s", out)
	}
	if left > top || top > bottom || bottom > fill {
		t.Fatalf("tcPr children out of order (want lnL<lnT<lnB<fill): %s", out)
	}
}

// Issue #71 asks for cap, join and inset pen alongside colour, width and dash.
func TestBuildBorderLineXMLCarriesCapJoinAndInset(t *testing.T) {
	xml := buildBorderLineXML("lnT", &CellBorderSideUpdate{
		Width:    38100,
		Color:    "0070C0",
		Dash:     "dash",
		Cap:      "square",
		Join:     "miter",
		Compound: "double",
		Inset:    true,
	})

	for _, want := range []string{
		`w="38100"`, `cap="sq"`, `cmpd="dbl"`, `algn="in"`,
		`<a:srgbClr val="0070C0"/>`, `<a:prstDash val="dash"/>`, `<a:miter lim="800000"/>`,
	} {
		if !strings.Contains(xml, want) {
			t.Fatalf("missing %s in %s", want, xml)
		}
	}
	if strings.Index(xml, "<a:prstDash") > strings.Index(xml, "<a:miter") {
		t.Fatalf("prstDash must precede the join: %s", xml)
	}
}

func TestBuildBorderLineXMLJoinVariants(t *testing.T) {
	cases := map[string]string{
		"round": "<a:round/>",
		"bevel": "<a:bevel/>",
		"":      "",
	}
	for join, want := range cases {
		xml := buildBorderLineXML("lnB", &CellBorderSideUpdate{Width: 1, Join: join})
		if want == "" {
			if strings.Contains(xml, "<a:round/>") || strings.Contains(xml, "<a:bevel/>") ||
				strings.Contains(xml, "<a:miter") {
				t.Fatalf("join %q should emit none, got %s", join, xml)
			}
			continue
		}
		if !strings.Contains(xml, want) {
			t.Fatalf("join %q = %s, want %s", join, xml, want)
		}
	}
}
