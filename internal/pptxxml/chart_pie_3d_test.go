package pptxxml

import (
	"strings"
	"testing"
)

func TestRenderThreeDPieChartUsesSchemaValidElements(t *testing.T) {
	xml := string(RenderChart(&ChartSpec{
		Kind:       ChartKindThreeDPie,
		Title:      "3D Pie",
		Categories: []string{"A", "B"},
		Values:     []float64{40, 60},
		SeriesName: "Share",
	}))
	for _, expected := range []string{
		`<c:view3D><c:rotX val="30"/><c:rotY val="0"/>`,
		`<c:rAngAx val="1"/><c:perspective val="30"/></c:view3D>`,
		`<c:pie3DChart>`,
		`</c:pie3DChart>`,
	} {
		if !strings.Contains(xml, expected) {
			t.Fatalf("3D-pie XML missing %q: %s", expected, xml)
		}
	}
	if strings.Contains(xml, "<c:pieChart>") {
		t.Fatalf("3D pie must not serialize as a 2D pie: %s", xml)
	}
	if strings.Index(xml, "<c:view3D>") > strings.Index(xml, "<c:plotArea>") {
		t.Fatalf("c:view3D must precede c:plotArea: %s", xml)
	}
}
