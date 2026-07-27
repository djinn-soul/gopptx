package shape

import (
	"strings"
	"testing"

	editormodcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/common"
)

const styleOrderGeometry = `<a:xfrm><a:off x="0" y="0"/><a:ext cx="10" cy="10"/></a:xfrm>` +
	`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>`

// PowerPoint drops a shape whose spPr children violate the CT_ShapeProperties
// sequence, so an effects-only update must not land ahead of an existing fill.
func TestReplaceStyleInSpPrKeepsSchemaOrder(t *testing.T) {
	original := []byte(`<p:sp><p:spPr>` + styleOrderGeometry +
		`<a:solidFill><a:srgbClr val="4472C4"/></a:solidFill>` +
		`<a:ln><a:solidFill><a:srgbClr val="1F3864"/></a:solidFill></a:ln>` +
		`</p:spPr></p:sp>`)

	updated := string(ReplaceStyleInSpPr(
		original,
		`<a:effectLst><a:outerShdw blurRad="80000"/></a:effectLst>`,
		false, false, true,
	))

	fill := strings.Index(updated, `<a:solidFill><a:srgbClr val="4472C4"/>`)
	line := strings.Index(updated, "<a:ln>")
	effect := strings.Index(updated, "<a:effectLst>")
	if fill < 0 || line < 0 || effect < 0 {
		t.Fatalf("missing style blocks: %s", updated)
	}
	if fill >= line || line >= effect {
		t.Fatalf("spPr children out of schema order: %s", updated)
	}
	if strings.Count(updated, "<a:solidFill>") != 2 {
		t.Fatalf("line's nested fill was hoisted or duplicated: %s", updated)
	}
	if !strings.Contains(updated, `<a:ln><a:solidFill><a:srgbClr val="1F3864"/></a:solidFill></a:ln>`) {
		t.Fatalf("line element was rewritten: %s", updated)
	}
}

func TestReplaceStyleInSpPrKeepsTrailingElementsLast(t *testing.T) {
	original := []byte(`<p:sp><p:spPr>` + styleOrderGeometry +
		`<a:solidFill><a:srgbClr val="FF0000"/></a:solidFill>` +
		`<a:extLst><a:ext uri="{X}"/></a:extLst>` +
		`</p:spPr></p:sp>`)

	updated := string(ReplaceStyleInSpPr(original, `<a:effectLst/>`, false, false, true))

	effect := strings.Index(updated, "<a:effectLst/>")
	ext := strings.Index(updated, "<a:extLst>")
	if effect < 0 || ext < 0 || effect > ext {
		t.Fatalf("effectLst must precede extLst: %s", updated)
	}
}

func TestSplitTopLevelElements(t *testing.T) {
	elements := editormodcommon.SplitTopLevelXMLElements(
		`<a:solidFill><a:srgbClr val="FF0000"/></a:solidFill><a:ln w="1"/><a:effectLst/>`,
	)
	if len(elements) != 3 {
		t.Fatalf("expected 3 elements, got %d (%+v)", len(elements), elements)
	}
	names := []string{elements[0].Name, elements[1].Name, elements[2].Name}
	want := []string{"a:solidFill", "a:ln", "a:effectLst"}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("element %d = %q, want %q", i, names[i], want[i])
		}
	}

	if got := editormodcommon.SplitTopLevelXMLElements(`text<a:ln/>`); got != nil {
		t.Fatalf("stray depth-0 text must bail out, got %+v", got)
	}
	if got := editormodcommon.SplitTopLevelXMLElements(`<a:ln>`); got != nil {
		t.Fatalf("unbalanced fragment must bail out, got %+v", got)
	}
}
