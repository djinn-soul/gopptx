package table

import (
	"strings"
	"testing"
)

func TestFillOnlyUpdatePreservesRichCellText(t *testing.T) {
	frame := `<a:tbl><a:tblGrid><a:gridCol w="100"/></a:tblGrid>` +
		`<a:tr h="100"><a:tc><a:txBody><a:bodyPr/><a:lstStyle/>` +
		`<a:p><a:r><a:rPr b="1"/><a:t>First</a:t></a:r></a:p>` +
		`<a:p><a:r><a:rPr i="1"/><a:t>Second</a:t></a:r></a:p>` +
		`</a:txBody><a:tcPr/></a:tc></a:tr></a:tbl>`

	out, err := UpdateTableCellContentInFrame([]byte(frame), 0, 0, CellContentUpdate{
		BackgroundColor: "FFE699",
	})
	if err != nil {
		t.Fatalf("update fill: %v", err)
	}
	got := string(out)
	for _, want := range []string{
		`<a:rPr b="1"/>`, `<a:t>First</a:t>`,
		`<a:rPr i="1"/>`, `<a:t>Second</a:t>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("fill-only update dropped %s: %s", want, got)
		}
	}
	if strings.Count(got, "<a:p>") != 2 {
		t.Fatalf("fill-only update flattened paragraphs: %s", got)
	}
}

func TestLineToBorderInfoReturnsNumericMiterLimit(t *testing.T) {
	line := &linePropertiesXML{
		Miter: &struct {
			Lim string `xml:"lim,attr"`
		}{Lim: "800000"},
	}
	border := lineToBorderInfo(line)

	if got, ok := border["miter_limit"].(int); !ok || got != 800000 {
		t.Fatalf("miter_limit = %#v, want numeric 800000", border["miter_limit"])
	}
}
