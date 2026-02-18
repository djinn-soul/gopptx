package pptxxml

import (
	"strings"
	"testing"
)

func TestSlideWithLayoutKeepsPlaceholderInsideCSld(t *testing.T) {
	xml := SlideWithLayout(
		slideLayoutTitleAndContent,
		TitleSpec{Text: "Title"},
		nil,
		nil,
		nil,
		ContentStyleSpec{},
		nil,
		nil,
		nil,
		nil,
		nil,
		[]PlaceholderOverrideSpec{{Index: 1, Type: "body", Text: "Body"}},
		nil,
		nil,
		"",
		"",
		false,
		"",
		false,
		9144000,
		6858000,
	)

	cSldEnd := strings.Index(xml, "</p:cSld>")
	placeholderIdx := strings.Index(xml, `<p:ph idx="1" type="body"/>`)
	if cSldEnd == -1 || placeholderIdx == -1 {
		t.Fatalf("missing expected XML markers")
	}
	if placeholderIdx > cSldEnd {
		t.Fatalf("placeholder emitted outside p:cSld")
	}
}

func TestSlideWithLayoutPlacesTransitionOutsideCSld(t *testing.T) {
	transition := `<p:transition><p:fade/></p:transition>`
	xml := SlideWithLayout(
		slideLayoutTitleAndContent,
		TitleSpec{Text: "Title"},
		nil,
		nil,
		nil,
		ContentStyleSpec{},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		transition,
		"",
		true,
		"Footer",
		true,
		9144000,
		6858000,
	)

	cSldEnd := strings.Index(xml, "</p:cSld>")
	transitionIdx := strings.Index(xml, transition)
	slideNumIdx := strings.Index(xml, `name="Slide Number Placeholder"`)
	footerIdx := strings.Index(xml, `name="Footer Placeholder"`)
	dateIdx := strings.Index(xml, `name="Date Placeholder"`)
	if cSldEnd == -1 || transitionIdx == -1 {
		t.Fatalf("missing expected XML markers")
	}
	if transitionIdx < cSldEnd {
		t.Fatalf("transition emitted inside p:cSld")
	}
	if slideNumIdx == -1 || footerIdx == -1 || dateIdx == -1 {
		t.Fatalf("expected overlay shapes not rendered")
	}
	if slideNumIdx > cSldEnd || footerIdx > cSldEnd || dateIdx > cSldEnd {
		t.Fatalf("overlay shape emitted outside p:cSld")
	}
}
