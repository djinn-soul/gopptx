package enums

import "testing"

func TestParseMSOShape(t *testing.T) {
	shape, err := ParseMSOShape("rounded rectangle")
	if err != nil {
		t.Fatalf("parse mso shape: %v", err)
	}
	if shape.XMLValue() != "roundRect" {
		t.Fatalf("expected roundRect, got %s", shape.XMLValue())
	}
	if _, err := ParseMSOShape("not-a-shape"); err == nil {
		t.Fatal("expected invalid shape error")
	}
}

func TestParsePPAlign(t *testing.T) {
	align, err := ParsePPAlign("justify_low")
	if err != nil {
		t.Fatalf("parse pp align: %v", err)
	}
	if align.XMLValue() != "justLow" {
		t.Fatalf("expected justLow, got %s", align.XMLValue())
	}
	if _, err := ParsePPAlign("diagonal"); err == nil {
		t.Fatal("expected invalid alignment error")
	}
}

func TestParseAnchorEnums(t *testing.T) {
	anchor, err := ParseMSOAnchor("middle")
	if err != nil {
		t.Fatalf("parse mso anchor: %v", err)
	}
	if anchor.XMLValue() != "ctr" {
		t.Fatalf("expected ctr, got %s", anchor.XMLValue())
	}
	vAnchor, err := ParseMSOVerticalAnchor("distribute")
	if err != nil {
		t.Fatalf("parse vertical anchor: %v", err)
	}
	if vAnchor.XMLValue() != "dist" {
		t.Fatalf("expected dist, got %s", vAnchor.XMLValue())
	}
}

func TestParseThemeColorAndColorType(t *testing.T) {
	color, err := ParseMSOThemeColor("followed_hyperlink")
	if err != nil {
		t.Fatalf("parse theme color: %v", err)
	}
	if color.XMLValue() != "folHlink" {
		t.Fatalf("expected folHlink, got %s", color.XMLValue())
	}
	colorType, err := ParseMSOColorType("theme")
	if err != nil {
		t.Fatalf("parse color type: %v", err)
	}
	if colorType.XMLValue() != "scheme" {
		t.Fatalf("expected scheme, got %s", colorType.XMLValue())
	}
}

func TestParseXLChartType(t *testing.T) {
	chartType, err := ParseXLChartType("donut")
	if err != nil {
		t.Fatalf("parse xl chart type: %v", err)
	}
	if chartType.XMLValue() != "doughnut" {
		t.Fatalf("expected doughnut, got %s", chartType.XMLValue())
	}
	if _, err := ParseXLChartType("surface3D"); err == nil {
		t.Fatal("expected invalid chart type error")
	}
}
