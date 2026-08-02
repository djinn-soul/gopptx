package shape

import "testing"

const groupWithChildSpaceXML = `<p:grpSp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" ` +
	`xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
<p:nvGrpSpPr><p:cNvPr id="7" name="Group 7"/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr>
<p:grpSpPr><a:xfrm>
<a:off x="1000000" y="2000000"/><a:ext cx="3000000" cy="1500000"/>
<a:chOff x="500000" y="600000"/><a:chExt cx="6000000" cy="3000000"/>
</a:xfrm></p:grpSpPr>
</p:grpSp>`

// A group states the space its children are drawn in. Without reading it, a
// child's raw coordinates get reported as slide coordinates and the child lands
// in the wrong place (upstream issue #925).
func TestParseShapePropertiesReadsGroupChildSpace(t *testing.T) {
	props, err := ParseShapeProperties([]byte(groupWithChildSpaceXML))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if props.GroupChild == nil {
		t.Fatal("expected the group's child coordinate space to be parsed")
	}
	if props.GroupChild.OffsetX != 500000 || props.GroupChild.OffsetY != 600000 {
		t.Fatalf("unexpected child offset: %+v", props.GroupChild)
	}
	if props.GroupChild.ExtentCx != 6000000 || props.GroupChild.ExtentCy != 3000000 {
		t.Fatalf("unexpected child extent: %+v", props.GroupChild)
	}
}

const groupWithoutChildSpaceXML = `<p:grpSp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" ` +
	`xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
<p:nvGrpSpPr><p:cNvPr id="8" name="Group 8"/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr>
<p:grpSpPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="10" cy="10"/></a:xfrm></p:grpSpPr>
</p:grpSp>`

func TestParseShapePropertiesLeavesChildSpaceUnsetWhenAbsent(t *testing.T) {
	props, err := ParseShapeProperties([]byte(groupWithoutChildSpaceXML))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if props.GroupChild != nil {
		t.Fatalf("expected no child space, got %+v", props.GroupChild)
	}
}
