package pptxxml

import (
	"strings"
	"testing"
)

func ptr[T any](v T) *T { return &v }

func TestRenderPlaceholderOverrideGeometry(t *testing.T) {
	spec := PlaceholderOverrideSpec{
		Index: 1,
		Type:  "body",
		X:     ptr(int64(100)),
		Y:     ptr(int64(200)),
		CX:    ptr(int64(300)),
		CY:    ptr(int64(400)),
	}

	xml := placeholderShape(spec, 10)
	if !strings.Contains(xml, "<a:xfrm>") {
		t.Errorf("expected <a:xfrm>, got %s", xml)
	}
	if !strings.Contains(xml, "off x=\"100\" y=\"200\"") {
		t.Errorf("expected off x=\"100\" y=\"200\", got %s", xml)
	}
	if !strings.Contains(xml, "ext cx=\"300\" cy=\"400\"") {
		t.Errorf("expected ext cx=\"300\" cy=\"400\", got %s", xml)
	}
}

func TestRenderPlaceholderOverrideTextStyle(t *testing.T) {
	spec := PlaceholderOverrideSpec{
		Index: 1,
		Type:  "body",
		TextStyle: &PlaceholderTextStyleSpec{
			SizePt: ptr(24),
			Color:  ptr("FF0000"),
			Bold:   ptr(true),
			Font:   ptr("Arial"),
		},
	}

	xml := placeholderShape(spec, 10)
	if !strings.Contains(xml, "<p:txBody>") {
		t.Errorf("expected <p:txBody>, got %s", xml)
	}
	if !strings.Contains(xml, "sz=\"2400\"") {
		t.Errorf("expected sz=\"2400\", got %s", xml)
	}
	if !strings.Contains(xml, "val=\"FF0000\"") {
		t.Errorf("expected val=\"FF0000\", got %s", xml)
	}
	if !strings.Contains(xml, "b=\"1\"") {
		t.Errorf("expected b=\"1\", got %s", xml)
	}
	if !strings.Contains(xml, "typeface=\"Arial\"") {
		t.Errorf("expected typeface=\"Arial\", got %s", xml)
	}
}
