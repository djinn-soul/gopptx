package export

import "testing"

// slideXMLWithTree wraps children in the spTree scaffolding a real slide has.
func slideXMLWithTree(children string) []byte {
	return []byte(`<?xml version="1.0"?>
<p:sld xmlns:p="p" xmlns:a="a">
 <p:cSld><p:spTree>
  <p:nvGrpSpPr><p:cNvPr id="1" name=""/></p:nvGrpSpPr>
  ` + children + `
 </p:spTree></p:cSld>
</p:sld>`)
}

func TestParseSlideZOrderRecordsTreePositions(t *testing.T) {
	t.Parallel()

	order := parseSlideZOrder(slideXMLWithTree(`
  <p:sp><p:nvSpPr><p:cNvPr id="2" name="card"/></p:nvSpPr></p:sp>
  <p:pic><p:nvPicPr><p:cNvPr id="5" name="photo"/></p:nvPicPr></p:pic>
  <p:sp><p:nvSpPr><p:cNvPr id="6" name="chip"/></p:nvSpPr></p:sp>`))

	if len(order) != 3 {
		t.Fatalf("order=%v want 3 entries", order)
	}
	// The picture sits between the two shapes, which is exactly the case the
	// old layered renderer got wrong.
	if order[2] >= order[5] || order[5] >= order[6] {
		t.Fatalf("order=%v want card < photo < chip", order)
	}
}

func TestParseSlideZOrderIgnoresNestedShapeIDs(t *testing.T) {
	t.Parallel()

	// A grouped shape carries child cNvPr ids inside it. Only the group's own id
	// takes part in the slide's z-order.
	order := parseSlideZOrder(slideXMLWithTree(`
  <p:grpSp>
    <p:nvGrpSpPr><p:cNvPr id="10" name="group"/></p:nvGrpSpPr>
    <p:sp><p:nvSpPr><p:cNvPr id="11" name="inner"/></p:nvSpPr></p:sp>
  </p:grpSp>
  <p:sp><p:nvSpPr><p:cNvPr id="12" name="after"/></p:nvSpPr></p:sp>`))

	if _, nested := order[11]; nested {
		t.Fatalf("order=%v should not contain the nested shape id 11", order)
	}
	if order[10] >= order[12] {
		t.Fatalf("order=%v want group before the shape that follows it", order)
	}
}

func TestParseSlideZOrderHandlesMissingTree(t *testing.T) {
	t.Parallel()

	if order := parseSlideZOrder([]byte(`<p:sld xmlns:p="p"></p:sld>`)); len(order) != 0 {
		t.Fatalf("order=%v want empty for a slide with no shape tree", order)
	}
	if order := parseSlideZOrder([]byte("not xml at all")); len(order) != 0 {
		t.Fatalf("order=%v want empty for unparseable input", order)
	}
}
