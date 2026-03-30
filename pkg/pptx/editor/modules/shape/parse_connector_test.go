package shape

import "testing"

func TestParseShapeProperties_ParsesConnectorMetadata(t *testing.T) {
	xml := []byte(`
<p:cxnSp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"
         xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <p:nvCxnSpPr>
    <p:cNvPr id="21" name="Connector 21"/>
    <p:cNvCxnSpPr>
      <a:stCxn id="2" idx="1"/>
      <a:endCxn id="3" idx="3"/>
    </p:cNvCxnSpPr>
    <p:nvPr/>
  </p:nvCxnSpPr>
  <p:spPr>
    <a:xfrm flipH="1">
      <a:off x="100" y="200"/>
      <a:ext cx="300" cy="400"/>
    </a:xfrm>
    <a:prstGeom prst="bentConnector3"><a:avLst/></a:prstGeom>
    <a:ln w="12700">
      <a:solidFill><a:srgbClr val="4472C4"/></a:solidFill>
      <a:headEnd type="triangle" w="sm" len="lg"/>
      <a:tailEnd type="stealth" w="lg" len="sm"/>
    </a:ln>
  </p:spPr>
</p:cxnSp>`)

	props, err := ParseShapeProperties(xml)
	if err != nil {
		t.Fatalf("ParseShapeProperties error: %v", err)
	}
	if props.ID != 21 || props.Name != "Connector 21" {
		t.Fatalf("expected connector identity, got %+v", props)
	}
	if props.Type != "bentConnector3" {
		t.Fatalf("expected connector type bentConnector3, got %q", props.Type)
	}
	if props.Connector == nil {
		t.Fatal("expected connector metadata")
	}
	if props.Connector.StartShapeID == nil || *props.Connector.StartShapeID != 2 {
		t.Fatalf("expected start shape id 2, got %+v", props.Connector)
	}
	if props.Connector.StartSiteIndex == nil || *props.Connector.StartSiteIndex != 1 {
		t.Fatalf("expected start site index 1, got %+v", props.Connector)
	}
	if props.Connector.EndShapeID == nil || *props.Connector.EndShapeID != 3 {
		t.Fatalf("expected end shape id 3, got %+v", props.Connector)
	}
	if props.Connector.EndSiteIndex == nil || *props.Connector.EndSiteIndex != 3 {
		t.Fatalf("expected end site index 3, got %+v", props.Connector)
	}
	if !props.Connector.FlipH || props.Connector.FlipV {
		t.Fatalf("expected flipH only, got %+v", props.Connector)
	}
	if props.Line == nil || props.Line.StartArrow == nil || *props.Line.StartArrow != "triangle" {
		t.Fatalf("expected parsed connector line arrows, got %+v", props.Line)
	}
}
