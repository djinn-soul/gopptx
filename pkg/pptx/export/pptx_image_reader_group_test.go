package export

import "testing"

// A picture inside a group states its geometry in the group's child space. When
// the group has been resized, that space no longer matches the slide, and the
// picture has to be mapped through the group's transform.
const testGroupedPicSlideXML = `<?xml version="1.0" encoding="UTF-8"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
       xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
       xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld><p:spTree>
    <p:grpSp>
      <p:nvGrpSpPr><p:cNvPr id="4" name="Group"/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr>
      <p:grpSpPr>
        <a:xfrm>
          <a:off x="1000000" y="2000000"/>
          <a:ext cx="2000000" cy="1000000"/>
          <a:chOff x="0" y="0"/>
          <a:chExt cx="1000000" cy="500000"/>
        </a:xfrm>
      </p:grpSpPr>
      <p:pic>
        <p:nvPicPr><p:cNvPr id="5" name="Picture"/><p:cNvPicPr/><p:nvPr/></p:nvPicPr>
        <p:blipFill><a:blip r:embed="rId2"/></p:blipFill>
        <p:spPr><a:xfrm><a:off x="100000" y="50000"/><a:ext cx="400000" cy="200000"/></a:xfrm></p:spPr>
      </p:pic>
    </p:grpSp>
  </p:spTree></p:cSld>
</p:sld>`

func TestParsePicElementsMapsGroupedPicturesOntoTheSlide(t *testing.T) {
	t.Parallel()

	pics := parsePicElements([]byte(testGroupedPicSlideXML))
	if len(pics) != 1 {
		t.Fatalf("got %d pictures, want the one inside the group", len(pics))
	}
	pic := pics[0]
	// The group doubles its child space, so the picture doubles in size and its
	// offset is doubled before the group's own position is added.
	if pic.X != 1200000 || pic.Y != 2100000 {
		t.Errorf("position = (%d, %d), want (1200000, 2100000)", pic.X, pic.Y)
	}
	if pic.CX != 800000 || pic.CY != 400000 {
		t.Errorf("size = (%d, %d), want (800000, 400000)", pic.CX, pic.CY)
	}
	if pic.ShapeID != 5 {
		t.Errorf("shape id = %d, want 5", pic.ShapeID)
	}
}

// A picture that is not in a group keeps the geometry it states.
const testUngroupedPicSlideXML = `<?xml version="1.0" encoding="UTF-8"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
       xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
       xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld><p:spTree>
    <p:pic>
      <p:nvPicPr><p:cNvPr id="2" name="Picture"/><p:cNvPicPr/><p:nvPr/></p:nvPicPr>
      <p:blipFill><a:blip r:embed="rId2"/></p:blipFill>
      <p:spPr><a:xfrm><a:off x="300000" y="400000"/><a:ext cx="500000" cy="600000"/></a:xfrm></p:spPr>
    </p:pic>
  </p:spTree></p:cSld>
</p:sld>`

func TestParsePicElementsLeavesUngroupedPicturesAlone(t *testing.T) {
	t.Parallel()

	pics := parsePicElements([]byte(testUngroupedPicSlideXML))
	if len(pics) != 1 {
		t.Fatalf("got %d pictures, want 1", len(pics))
	}
	if pics[0].X != 300000 || pics[0].Y != 400000 || pics[0].CX != 500000 || pics[0].CY != 600000 {
		t.Fatalf("geometry = %+v, want it unchanged", pics[0])
	}
}
