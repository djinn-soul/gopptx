package editor

import "testing"

func TestExtractAllSmartArtRelIDsMatchesExactShapeID(t *testing.T) {
	slideXML := `<p:sld xmlns:p="urn:p" xmlns:dgm="urn:dgm" xmlns:r="urn:r">` +
		`<p:sp><p:nvSpPr><p:cNvPr id="123" name="Noise"/></p:nvSpPr></p:sp>` +
		`<p:graphicFrame>` +
		`<p:nvGraphicFramePr><p:cNvPr id="1" name="SmartArt One"/></p:nvGraphicFramePr>` +
		`<a:graphic><a:graphicData><dgm:relIds r:dm="rIdDm1" r:lo="rIdLo1" r:qs="rIdQs1" r:cs="rIdCs1"/></a:graphicData></a:graphic>` +
		`</p:graphicFrame>` +
		`<p:graphicFrame>` +
		`<p:nvGraphicFramePr><p:cNvPr id="23" name="SmartArt Two"/></p:nvGraphicFramePr>` +
		`<a:graphic><a:graphicData><dgm:relIds r:dm="rIdDm23" r:lo="rIdLo23" r:qs="rIdQs23" r:cs="rIdCs23"/></a:graphicData></a:graphic>` +
		`</p:graphicFrame>` +
		`</p:sld>`

	dm, lo, qs, cs := extractAllSmartArtRelIDs(slideXML, 23)
	if dm != "rIdDm23" || lo != "rIdLo23" || qs != "rIdQs23" || cs != "rIdCs23" {
		t.Fatalf("unexpected rel IDs: dm=%q lo=%q qs=%q cs=%q", dm, lo, qs, cs)
	}
}

func TestExtractAllSmartArtRelIDsReturnsEmptyWhenShapeNotSmartArt(t *testing.T) {
	slideXML := `<p:sld xmlns:p="urn:p">` +
		`<p:graphicFrame>` +
		`<p:nvGraphicFramePr><p:cNvPr id='42' name='Not SmartArt'/></p:nvGraphicFramePr>` +
		`<a:graphic><a:graphicData/></a:graphic>` +
		`</p:graphicFrame>` +
		`</p:sld>`

	dm, lo, qs, cs := extractAllSmartArtRelIDs(slideXML, 42)
	if dm != "" || lo != "" || qs != "" || cs != "" {
		t.Fatalf("expected empty rel IDs, got: dm=%q lo=%q qs=%q cs=%q", dm, lo, qs, cs)
	}
}

func TestExtractAllSmartArtRelIDsSupportsSingleQuotedShapeID(t *testing.T) {
	slideXML := `<p:sld xmlns:p="urn:p" xmlns:dgm="urn:dgm" xmlns:r="urn:r">` +
		`<p:graphicFrame>` +
		`<p:nvGraphicFramePr><p:cNvPr id='77' name='SmartArt Single Quote'/></p:nvGraphicFramePr>` +
		`<a:graphic><a:graphicData><dgm:relIds r:dm="dm77" r:lo="lo77" r:qs="qs77" r:cs="cs77"/></a:graphicData></a:graphic>` +
		`</p:graphicFrame>` +
		`</p:sld>`

	dm, lo, qs, cs := extractAllSmartArtRelIDs(slideXML, 77)
	if dm != "dm77" || lo != "lo77" || qs != "qs77" || cs != "cs77" {
		t.Fatalf("unexpected rel IDs: dm=%q lo=%q qs=%q cs=%q", dm, lo, qs, cs)
	}
}

func TestExtractAllSmartArtRelIDsDoesNotBorrowFromNeighborFrame(t *testing.T) {
	slideXML := `<p:sld xmlns:p="urn:p" xmlns:dgm="urn:dgm" xmlns:r="urn:r">` +
		`<p:graphicFrame>` +
		`<p:nvGraphicFramePr><p:cNvPr id="112" name="Neighbor"/></p:nvGraphicFramePr>` +
		`<a:graphic><a:graphicData><dgm:relIds r:dm="dm112" r:lo="lo112" r:qs="qs112" r:cs="cs112"/></a:graphicData></a:graphic>` +
		`</p:graphicFrame>` +
		`<p:graphicFrame>` +
		`<p:nvGraphicFramePr><p:cNvPr id="12" name="Target"/></p:nvGraphicFramePr>` +
		`<a:graphic><a:graphicData/></a:graphic>` +
		`</p:graphicFrame>` +
		`</p:sld>`

	dm, lo, qs, cs := extractAllSmartArtRelIDs(slideXML, 12)
	if dm != "" || lo != "" || qs != "" || cs != "" {
		t.Fatalf("expected empty rel IDs for non-SmartArt target, got: dm=%q lo=%q qs=%q cs=%q", dm, lo, qs, cs)
	}
}

func TestIsDiagramDrawingRelTypeSupportsModernAndLegacyURIs(t *testing.T) {
	if !isDiagramDrawingRelType(relTypeDiagramDrawing) {
		t.Fatalf("expected modern diagramDrawing rel type to be recognized: %q", relTypeDiagramDrawing)
	}
	if !isDiagramDrawingRelType(relTypeDiagramDrawingLegacy) {
		t.Fatalf("expected legacy diagramDrawing rel type to be recognized: %q", relTypeDiagramDrawingLegacy)
	}
	if isDiagramDrawingRelType(relTypeDiagramData) {
		t.Fatalf("non-drawing rel type must not be recognized as diagramDrawing: %q", relTypeDiagramData)
	}
}
