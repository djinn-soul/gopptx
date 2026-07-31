package pptxxml

import (
	"strings"
	"testing"
)

// A shape whose a:bodyPr omits an inset must re-render with OOXML's implicit
// value (ECMA-376 §20.1.10.44), not with the half-inch slide-layout margin.
// Using defaultMargin here reflowed the text of every edited shape.
func TestTextBodyPrXMLUsesOOXMLDefaultInsets(t *testing.T) {
	xml := TextBodyPrXML(nil)

	for _, want := range []string{
		`lIns="91440"`,
		`rIns="91440"`,
		`tIns="45720"`,
		`bIns="45720"`,
	} {
		if !strings.Contains(xml, want) {
			t.Errorf("expected %s in bodyPr, got:\n%s", want, xml)
		}
	}
	if strings.Contains(xml, "457200") {
		t.Errorf("bodyPr must not use the layout margin as a text inset:\n%s", xml)
	}
}

func TestTextBodyPrXMLExplicitInsetsWin(t *testing.T) {
	xml := TextBodyPrXML(&TextFrameSpec{
		MarginLeft:   548640,
		MarginRight:  228600,
		MarginTop:    274320,
		MarginBottom: 91440,
		Wrap:         "square",
		Anchor:       "ctr",
	})

	for _, want := range []string{
		`lIns="548640"`,
		`rIns="228600"`,
		`tIns="274320"`,
		`bIns="91440"`,
	} {
		if !strings.Contains(xml, want) {
			t.Errorf("expected %s in bodyPr, got:\n%s", want, xml)
		}
	}
}

// defaultMargin still positions title and content placeholders half an inch from
// the slide edge; splitting the constants must not have merged the two roles.
func TestLayoutMarginIsDistinctFromTextInset(t *testing.T) {
	if defaultMargin != 457200 {
		t.Errorf("slide-edge layout margin changed: %d", defaultMargin)
	}
	if defaultInsetLR == defaultMargin || defaultInsetTB == defaultMargin {
		t.Error("text insets must not reuse the slide-edge layout margin")
	}
}
