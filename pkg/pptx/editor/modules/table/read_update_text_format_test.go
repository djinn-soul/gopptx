package table

import (
	"strings"
	"testing"
)

const styledCellFrame = `<a:tbl><a:tblGrid><a:gridCol w="100"/></a:tblGrid>` +
	`<a:tr h="370840"><a:tc><a:txBody><a:bodyPr/><a:lstStyle/>` +
	`<a:p><a:pPr algn="ctr"/><a:r>` +
	`<a:rPr lang="en-US" sz="2800" b="1"><a:solidFill><a:srgbClr val="C00000"/></a:solidFill>` +
	`<a:latin typeface="Georgia"/></a:rPr>` +
	`<a:t>Old</a:t></a:r></a:p></a:txBody>` +
	`<a:tcPr><a:solidFill><a:srgbClr val="FFE699"/></a:solidFill></a:tcPr>` +
	`</a:tc></a:tr></a:tbl>`

// Issue #1037: updating a cell's data must keep the original formatting.
func TestUpdateTableCellTextPreservesRunFormatting(t *testing.T) {
	out, err := UpdateTableCellTextInFrame([]byte(styledCellFrame), 0, 0, "New")
	if err != nil {
		t.Fatalf("update text: %v", err)
	}

	got := string(out)
	if !strings.Contains(got, "<a:t>New</a:t>") {
		t.Fatalf("text not updated: %s", got)
	}
	for _, want := range []string{
		`sz="2800"`, `b="1"`, `<a:srgbClr val="C00000"/>`, `<a:latin typeface="Georgia"/>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("run formatting %s lost: %s", want, got)
		}
	}
	if !strings.Contains(got, `<a:pPr algn="ctr"/>`) {
		t.Fatalf("paragraph properties lost: %s", got)
	}
	if !strings.Contains(got, `<a:srgbClr val="FFE699"/>`) {
		t.Fatalf("cell fill lost: %s", got)
	}
	if strings.Contains(got, "<a:t>Old</a:t>") {
		t.Fatalf("old text still present: %s", got)
	}
}

func TestUpdateTableCellTextOnPlainCellStillWorks(t *testing.T) {
	plain := `<a:tbl><a:tblGrid><a:gridCol w="100"/></a:tblGrid>` +
		`<a:tr h="100"><a:tc><a:txBody><a:bodyPr/><a:lstStyle/>` +
		`<a:p><a:r><a:rPr/><a:t>Old</a:t></a:r></a:p></a:txBody><a:tcPr/></a:tc></a:tr></a:tbl>`

	out, err := UpdateTableCellTextInFrame([]byte(plain), 0, 0, "New")
	if err != nil {
		t.Fatalf("update text: %v", err)
	}
	if !strings.Contains(string(out), "<a:t>New</a:t>") {
		t.Fatalf("text not updated: %s", out)
	}
}

func TestUpdateTableCellTextEscapes(t *testing.T) {
	out, err := UpdateTableCellTextInFrame([]byte(styledCellFrame), 0, 0, `a & b <c>`)
	if err != nil {
		t.Fatalf("update text: %v", err)
	}
	if !strings.Contains(string(out), "a &amp; b &lt;c&gt;") {
		t.Fatalf("text not escaped: %s", out)
	}
}
