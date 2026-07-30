package table

import (
	"strings"
	"testing"
)

const plainCellFrame = `<a:tbl><a:tblGrid><a:gridCol w="100"/></a:tblGrid>` +
	`<a:tr h="370840"><a:tc><a:txBody><a:bodyPr/><a:lstStyle/>` +
	`<a:p><a:r><a:rPr/><a:t>Old</a:t></a:r></a:p></a:txBody><a:tcPr/></a:tc></a:tr></a:tbl>`

func boolPtr(v bool) *bool { return &v }

// update_table_cell used to accept bold, color and background_color and drop
// them silently. Each must now reach the XML.
func TestUpdateTableCellContentAppliesRunStyle(t *testing.T) {
	out, err := UpdateTableCellContentInFrame([]byte(plainCellFrame), 0, 0, CellContentUpdate{
		SizePt:    20,
		FontName:  "Georgia",
		Bold:      boolPtr(true),
		Italic:    boolPtr(true),
		Underline: boolPtr(true),
		Color:     "C00000",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	got := string(out)
	for _, want := range []string{
		`sz="2000"`, `b="1"`, `i="1"`, `u="sng"`,
		`<a:srgbClr val="C00000"/>`, `<a:latin typeface="Georgia"/>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %s in %s", want, got)
		}
	}
}

func TestUpdateTableCellContentClearsBoldWhenFalse(t *testing.T) {
	out, err := UpdateTableCellContentInFrame([]byte(plainCellFrame), 0, 0, CellContentUpdate{
		Bold: boolPtr(false),
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if !strings.Contains(string(out), `b="0"`) {
		t.Fatalf("bold=false should write b=\"0\": %s", out)
	}
}

func TestUpdateTableCellContentSetsAndClearsCellFill(t *testing.T) {
	filled, err := UpdateTableCellContentInFrame([]byte(plainCellFrame), 0, 0, CellContentUpdate{
		BackgroundColor: "FFE699",
	})
	if err != nil {
		t.Fatalf("set fill: %v", err)
	}
	if !strings.Contains(string(filled), `<a:solidFill><a:srgbClr val="FFE699"/></a:solidFill>`) {
		t.Fatalf("cell fill not applied: %s", filled)
	}

	cleared, err := UpdateTableCellContentInFrame(filled, 0, 0, CellContentUpdate{
		BackgroundColor: "none",
	})
	if err != nil {
		t.Fatalf("clear fill: %v", err)
	}
	if strings.Contains(string(cleared), "FFE699") {
		t.Fatalf("cell fill not cleared: %s", cleared)
	}
}

func TestUpdateTableCellContentReplacesFillRatherThanStacking(t *testing.T) {
	first, err := UpdateTableCellContentInFrame([]byte(plainCellFrame), 0, 0, CellContentUpdate{
		BackgroundColor: "FFE699",
	})
	if err != nil {
		t.Fatalf("first fill: %v", err)
	}
	second, err := UpdateTableCellContentInFrame(first, 0, 0, CellContentUpdate{
		BackgroundColor: "92D050",
	})
	if err != nil {
		t.Fatalf("second fill: %v", err)
	}
	got := string(second)
	if strings.Count(got, "<a:solidFill>") != 1 {
		t.Fatalf("fill stacked instead of replaced: %s", got)
	}
	if !strings.Contains(got, "92D050") || strings.Contains(got, "FFE699") {
		t.Fatalf("fill not replaced: %s", got)
	}
}

// Upstream #1037: refreshing a table's data every day must not flatten the
// formatting the deck author applied by hand.
const formattedCellFrame = `<a:tbl><a:tblGrid><a:gridCol w="100"/></a:tblGrid>` +
	`<a:tr h="370840"><a:tc><a:txBody><a:bodyPr/><a:lstStyle/>` +
	`<a:p><a:pPr algn="ctr"/><a:r>` +
	`<a:rPr lang="en-US" sz="1800" b="1" i="1" u="sng" dirty="0">` +
	`<a:solidFill><a:srgbClr val="1F4E79"/></a:solidFill>` +
	`<a:latin typeface="Garamond"/></a:rPr>` +
	`<a:t>42</a:t></a:r></a:p></a:txBody>` +
	`<a:tcPr anchor="ctr"><a:solidFill><a:srgbClr val="D9E2F3"/></a:solidFill></a:tcPr>` +
	`</a:tc></a:tr></a:tbl>`

func TestUpdateTableCellTextPreservesExistingFormatting(t *testing.T) {
	out, err := UpdateTableCellTextInFrame([]byte(formattedCellFrame), 0, 0, "99")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	got := string(out)
	for _, want := range []string{
		`sz="1800"`, `b="1"`, `i="1"`, `u="sng"`,
		`<a:srgbClr val="1F4E79"/>`, `<a:latin typeface="Garamond"/>`,
		`<a:pPr algn="ctr"/>`, `anchor="ctr"`, `<a:srgbClr val="D9E2F3"/>`,
		`<a:t>99</a:t>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("text update dropped %s: %s", want, got)
		}
	}
	if strings.Contains(got, ">42<") {
		t.Fatalf("old text survived: %s", got)
	}
}

// A partial style update changes only what it names.
func TestUpdateTableCellContentMergesOntoExistingRunStyle(t *testing.T) {
	out, err := UpdateTableCellContentInFrame([]byte(formattedCellFrame), 0, 0, CellContentUpdate{
		Bold: boolPtr(false),
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, `b="0"`) {
		t.Fatalf("bold not cleared: %s", got)
	}
	for _, want := range []string{`sz="1800"`, `i="1"`, `<a:latin typeface="Garamond"/>`, `<a:t>42</a:t>`} {
		if !strings.Contains(got, want) {
			t.Fatalf("merge dropped %s: %s", want, got)
		}
	}
}

// The fill is a choice group: a new run colour replaces the old one instead of
// producing two fills, which PowerPoint rejects.
func TestMergedRunPropertiesKeepSchemaOrderAndSingleFill(t *testing.T) {
	out, err := UpdateTableCellContentInFrame([]byte(formattedCellFrame), 0, 0, CellContentUpdate{
		Color: "C00000",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	got := string(out)
	rPr := got[strings.Index(got, "<a:rPr "):strings.Index(got, "</a:rPr>")]
	if strings.Count(rPr, "<a:solidFill>") != 1 {
		t.Fatalf("run fill stacked: %s", rPr)
	}
	if strings.Contains(rPr, "1F4E79") {
		t.Fatalf("old run colour survived: %s", rPr)
	}
	if strings.Index(rPr, "<a:solidFill>") > strings.Index(rPr, "<a:latin") {
		t.Fatalf("solidFill must precede latin: %s", rPr)
	}
}

// An empty cell keeps the formatting parked on <a:endParaRPr>.
func TestUpdateTableCellTextAdoptsEndParaRPr(t *testing.T) {
	frame := `<a:tbl><a:tblGrid><a:gridCol w="100"/></a:tblGrid>` +
		`<a:tr h="370840"><a:tc><a:txBody><a:bodyPr/><a:lstStyle/>` +
		`<a:p><a:endParaRPr lang="en-US" sz="1400" b="1"/></a:p></a:txBody><a:tcPr/></a:tc></a:tr></a:tbl>`
	out, err := UpdateTableCellTextInFrame([]byte(frame), 0, 0, "New")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, `sz="1400"`) || !strings.Contains(got, `b="1"`) {
		t.Fatalf("endParaRPr formatting not adopted: %s", got)
	}
}

func TestHasRunStyleIgnoresBackgroundOnly(t *testing.T) {
	if (CellContentUpdate{BackgroundColor: "FFE699"}).HasRunStyle() {
		t.Fatal("a cell fill is not run-level styling")
	}
	if !(CellContentUpdate{Bold: boolPtr(false)}).HasRunStyle() {
		t.Fatal("bold=false is still a run style change")
	}
}
