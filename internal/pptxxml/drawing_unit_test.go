package pptxxml

import (
	"strings"
	"testing"
)

func TestCustomShapeXML_Internal(t *testing.T) {
	spec := ShapeSpec{
		Type: "rect",
		X:    100, Y: 100, CX: 100, CY: 100,
		Fill:        &ShapeFillSpec{Color: "FF0000"},
		Line:        &ShapeLineSpec{Color: "000000", Width: 1000},
		Text:        "Hello",
		ClickAction: &HyperlinkSpec{RelID: "rId1"},
		HoverAction: &HyperlinkSpec{RelID: "rId2"},
		Effects:     &ShapeEffectsSpec{Shadow: true},
		Adjustments: []ConnectorAdjustmentSpec{{Name: "adj1", Formula: "val 1"}},
	}
	xml := customShapeXML(spec, 5)
	if !strings.Contains(xml, "hlinkClick") {
		t.Error("Click action missing")
	}
	if !strings.Contains(xml, "hlinkHover") {
		t.Error("Hover action missing")
	}
	if !strings.Contains(xml, "adj1") {
		t.Error("Adjustment missing")
	}
}

func TestHyperlinkXML_ActionOnlyInternal(t *testing.T) {
	xml := HyperlinkXML(HyperlinkSpec{
		Action: "ppaction://macro?name=RunOnly",
	}, "a:hlinkClick")
	if strings.Contains(xml, `r:id=`) {
		t.Fatalf("expected action-only hyperlink without relationship id, got %s", xml)
	}
	if !strings.Contains(xml, `action="ppaction://macro?name=RunOnly"`) {
		t.Fatalf("expected action-only hyperlink to preserve action, got %s", xml)
	}
}

func TestConnectorXML_Extra_Internal(t *testing.T) {
	spec := ConnectorSpec{
		Type:   "straightConnector1",
		StartX: 0, StartY: 0, EndX: 100, EndY: 100,
		Line: ShapeLineSpec{Cap: "rnd", Join: "miter"},
	}
	xml := connectorXML(spec, 1, 0, 0)
	if !strings.Contains(xml, `cap="rnd"`) {
		t.Error("Cap missing")
	}
	if !strings.Contains(xml, `miter`) {
		t.Error("Join missing")
	}
}

func TestCustomShapeTextBody_Internal(t *testing.T) {
	spec := ShapeSpec{
		Text: "Hello World",
		TextFrame: &TextFrameSpec{
			Wrap:    "square",
			Anchor:  "ctr",
			AutoFit: "spAutoFit",
		},
	}
	xml := customShapeTextBody(spec)
	if !strings.Contains(xml, "Hello World") {
		t.Error("Text missing")
	}
	if !strings.Contains(xml, "spAutoFit") {
		t.Error("AutoFit missing")
	}

	// Test normAutoFit
	spec.TextFrame.AutoFit = "normAutoFit"
	xml = customShapeTextBody(spec)
	if !strings.Contains(xml, "normAutoFit") {
		t.Error("NormAutoFit missing")
	}
	spec.TextFrame.Orientation = "vert"
	xml = customShapeTextBody(spec)
	if !strings.Contains(xml, `vert="vert"`) {
		t.Error("Orientation missing")
	}
}

func TestBulletParagraphSpecIsZero(t *testing.T) {
	if !(BulletParagraphSpec{}).IsZero() {
		t.Fatal("zero-value BulletParagraphSpec must be IsZero")
	}
	if (BulletParagraphSpec{LineSpacingPts: 1}).IsZero() {
		t.Fatal("LineSpacingPts must make BulletParagraphSpec non-zero")
	}
	if (BulletParagraphSpec{TabStops: []int64{914400}}).IsZero() {
		t.Fatal("TabStops must make BulletParagraphSpec non-zero")
	}
}

func TestBulletParagraphPropsXML_WithLineSpacingPtsAndTabStops(t *testing.T) {
	xml := BulletParagraphPropsXML(BulletParagraphSpec{
		LineSpacingPts: 18,
		TabStops:       []int64{914400},
	})
	if !strings.Contains(xml, `<a:lnSpc><a:spcPts val="1800"/></a:lnSpc>`) {
		t.Fatalf("missing line spacing points xml: %s", xml)
	}
	if !strings.Contains(xml, `<a:tabLst><a:tab pos="914400"/></a:tabLst>`) {
		t.Fatalf("missing tab stops xml: %s", xml)
	}
}

func TestSlideRelationships_Internal(t *testing.T) {
	xml := SlideRelationshipsWithMultiCharts(
		"../layout.xml",
		[]string{"media/img1.png"},
		&ChartRel{RID: "rId2", Target: "charts/chart1.xml"},
		[]ChartRel{{RID: "rId3", Target: "charts/chart2.xml"}},
		[]SmartArtRel{{RID: "rId4", Type: "urn:sa", Target: "smartArt/sa1.xml"}},
		"notes/note1.xml",
		[]HyperlinkRel{{RID: "rId5", Target: "http://x.com", External: true}},
		"comments/comment1.xml",
	)

	checks := []string{"img1.png", "chart1.xml", "chart2.xml", "sa1.xml", "note1.xml", "http://x.com", "comment1.xml"}
	for _, c := range checks {
		if !strings.Contains(xml, c) {
			t.Errorf("missing %s in Relationships", c)
		}
	}
}
