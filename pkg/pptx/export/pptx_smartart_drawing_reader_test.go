package export

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const (
	nsRelationships = `xmlns="http://schemas.openxmlformats.org/package/2006/relationships"`
	relTypeSlide    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide"
	relTypeDiagram  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/diagram"
)

// smartArtSlideXML is a slide carrying one diagram frame at 1in,1in sized 6x3in.
const smartArtSlideXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
	`<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"` +
	` xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"` +
	` xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">` +
	`<p:cSld><p:spTree><p:graphicFrame><p:nvGraphicFramePr>` +
	`<p:cNvPr id="4" name="Diagram 4"/><p:cNvGraphicFramePr/><p:nvPr/></p:nvGraphicFramePr>` +
	`<p:xfrm><a:off x="914400" y="914400"/><a:ext cx="5486400" cy="2743200"/></p:xfrm>` +
	`<a:graphic><a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/diagram">` +
	`<dgm:relIds xmlns:dgm="http://schemas.openxmlformats.org/drawingml/2006/diagram"` +
	` r:dm="rId1" r:lo="rId2" r:qs="rId3" r:cs="rId4"/>` +
	`</a:graphicData></a:graphic></p:graphicFrame></p:spTree></p:cSld></p:sld>`

const smartArtLayoutPartXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
	`<dgm:layoutDef xmlns:dgm="http://schemas.openxmlformats.org/drawingml/2006/diagram"` +
	` uniqueId="urn:microsoft.com/office/officeart/2005/8/layout/process3"/>`

const smartArtDrawingPartXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
	`<dsp:drawing xmlns:dsp="http://schemas.microsoft.com/office/drawing/2008/diagram"` +
	` xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><dsp:spTree>` +
	`<dsp:sp modelId="{N1}"><dsp:spPr>` +
	`<a:xfrm><a:off x="0" y="0"/><a:ext cx="1828800" cy="914400"/></a:xfrm>` +
	`<a:prstGeom prst="roundRect"><a:avLst/></a:prstGeom>` +
	`<a:solidFill><a:schemeClr val="accent1"/></a:solidFill>` +
	`</dsp:spPr><dsp:txBody><a:bodyPr anchor="ctr"/><a:p><a:pPr algn="ctr"/>` +
	`<a:r><a:rPr sz="1800"/><a:t>Cached Step</a:t></a:r></a:p></dsp:txBody></dsp:sp>` +
	`</dsp:spTree></dsp:drawing>`

// smartArtDataPartXML holds one node so the diagram is not empty. dataModelExt
// is left out on purpose: gopptx's own writer drops it, and the reader has to
// find the drawing by the part naming convention instead.
const smartArtDataPartXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
	`<dgm:dataModel xmlns:dgm="http://schemas.openxmlformats.org/drawingml/2006/diagram"` +
	` xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><dgm:ptLst>` +
	`<dgm:pt modelId="{DOC}" type="doc"><dgm:prSet` +
	` loTypeId="urn:microsoft.com/office/officeart/2005/8/layout/process3"` +
	` csTypeId="urn:microsoft.com/office/officeart/2005/8/colors/accent1_2"/>` +
	`<dgm:t><a:p><a:endParaRPr/></a:p></dgm:t></dgm:pt>` +
	`<dgm:pt modelId="{N1}"><dgm:prSet/><dgm:t><a:p><a:r><a:t>Cached Step</a:t></a:r></a:p></dgm:t></dgm:pt>` +
	`</dgm:ptLst><dgm:cxnLst>` +
	`<dgm:cxn modelId="{C1}" srcId="{DOC}" destId="{N1}" srcOrd="0" destOrd="0"/>` +
	`</dgm:cxnLst></dgm:dataModel>`

func smartArtPackageFiles(dataPartXML, dataRels string) map[string]string {
	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
			`<Default Extension="xml" ContentType="application/xml"/></Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships ` + nsRelationships + `><Relationship Id="rId1"` +
			` Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"` +
			` Target="ppt/presentation.xml"/></Relationships>`,
		"ppt/presentation.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<p:presentation xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"` +
			` xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">` +
			`<p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst></p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships ` + nsRelationships + `><Relationship Id="rId1" Type="` + relTypeSlide +
			`" Target="slides/slide1.xml"/></Relationships>`,
		"ppt/slides/slide1.xml": smartArtSlideXML,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships ` + nsRelationships + `>` +
			`<Relationship Id="rId1" Type="` + relTypeDiagram + `" Target="../diagrams/data1.xml"/>` +
			`<Relationship Id="rId2" Type="` + relTypeDiagram + `" Target="../diagrams/layout1.xml"/>` +
			`<Relationship Id="rId3" Type="` + relTypeDiagram + `" Target="../diagrams/quickStyle1.xml"/>` +
			`<Relationship Id="rId4" Type="` + relTypeDiagram + `" Target="../diagrams/colors1.xml"/>` +
			`</Relationships>`,
		"ppt/diagrams/data1.xml":       dataPartXML,
		"ppt/diagrams/layout1.xml":     smartArtLayoutPartXML,
		"ppt/diagrams/quickStyle1.xml": `<a:styleDef xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"/>`,
		"ppt/diagrams/colors1.xml":     `<a:colorsDef xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"/>`,
		"ppt/diagrams/drawing1.xml":    smartArtDrawingPartXML,
	}
	if dataRels != "" {
		files["ppt/diagrams/_rels/data1.xml.rels"] = dataRels
	}
	return files
}

func writeTestPPTX(t *testing.T, name string, files map[string]string) string {
	t.Helper()
	outPath := filepath.Join(t.TempDir(), name)
	f, err := os.Create(outPath)
	if err != nil {
		t.Fatalf("create pptx: %v", err)
	}
	zw := zip.NewWriter(f)
	for partName, content := range files {
		w, werr := zw.Create(partName)
		if werr != nil {
			t.Fatalf("zip create %s: %v", partName, werr)
		}
		if _, werr = fmt.Fprint(w, content); werr != nil {
			t.Fatalf("zip write %s: %v", partName, werr)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zip close: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("file close: %v", err)
	}
	return outPath
}

func TestReadsSmartArtDrawingCacheByNamingConvention(t *testing.T) {
	path := writeTestPPTX(t, "smartart_cache.pptx", smartArtPackageFiles(smartArtDataPartXML, ""))

	_, slides, err := SlidesFromPPTX(path)
	if err != nil {
		t.Fatalf("SlidesFromPPTX: %v", err)
	}
	if len(slides) != 1 || len(slides[0].SmartArtDiagrams) != 1 {
		t.Fatalf("got %d slides, first with %d diagrams", len(slides), len(slides[0].SmartArtDiagrams))
	}

	drawing := slides[0].SmartArtDiagrams[0].Drawing
	if len(drawing) != 1 {
		t.Fatalf("got %d cached shapes, want 1", len(drawing))
	}
	if drawing[0].PresetGeom != "roundRect" {
		t.Errorf("PresetGeom = %q, want roundRect", drawing[0].PresetGeom)
	}
	if !drawing[0].HasText() {
		t.Error("cached shape carries no text")
	}
}

func TestReadsSmartArtDrawingCacheByDataModelExt(t *testing.T) {
	// A PowerPoint-authored deck names the drawing through the extension and a
	// relationship of the data part, under a name the convention would miss.
	dataXML := smartArtDataPartXML[:len(smartArtDataPartXML)-len(`</dgm:dataModel>`)] +
		`<dgm:extLst><a:ext uri="http://schemas.microsoft.com/office/drawing/2008/diagram">` +
		`<dsp:dataModelExt xmlns:dsp="http://schemas.microsoft.com/office/drawing/2008/diagram"` +
		` relId="rId9"/></a:ext></dgm:extLst></dgm:dataModel>`
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Relationships ` + nsRelationships + `><Relationship Id="rId9"` +
		` Type="http://schemas.microsoft.com/office/2007/relationships/diagramDrawing"` +
		` Target="cached-picture.xml"/></Relationships>`

	files := smartArtPackageFiles(dataXML, rels)
	files["ppt/diagrams/cached-picture.xml"] = smartArtDrawingPartXML
	delete(files, "ppt/diagrams/drawing1.xml")

	path := writeTestPPTX(t, "smartart_ext.pptx", files)
	_, slides, err := SlidesFromPPTX(path)
	if err != nil {
		t.Fatalf("SlidesFromPPTX: %v", err)
	}
	if len(slides[0].SmartArtDiagrams) != 1 {
		t.Fatalf("got %d diagrams, want 1", len(slides[0].SmartArtDiagrams))
	}
	if got := len(slides[0].SmartArtDiagrams[0].Drawing); got != 1 {
		t.Fatalf("got %d cached shapes via dataModelExt, want 1", got)
	}
}

func TestSmartArtDrawingCacheResolvesSchemeColors(t *testing.T) {
	path := writeTestPPTX(t, "smartart_colors.pptx", smartArtPackageFiles(smartArtDataPartXML, ""))

	_, slides, err := SlidesFromPPTX(path)
	if err != nil {
		t.Fatalf("SlidesFromPPTX: %v", err)
	}
	fill := slides[0].SmartArtDiagrams[0].Drawing[0].Fill
	if fill.Scheme != "" {
		t.Errorf("Scheme = %q, want it folded into a literal colour", fill.Scheme)
	}
	if len(fill.SRGB) != 6 {
		t.Errorf("SRGB = %q, want six hex digits", fill.SRGB)
	}
}

func TestMissingDrawingPartLeavesCacheEmpty(t *testing.T) {
	files := smartArtPackageFiles(smartArtDataPartXML, "")
	delete(files, "ppt/diagrams/drawing1.xml")

	path := writeTestPPTX(t, "smartart_nocache.pptx", files)
	_, slides, err := SlidesFromPPTX(path)
	if err != nil {
		t.Fatalf("SlidesFromPPTX: %v", err)
	}
	if got := len(slides[0].SmartArtDiagrams[0].Drawing); got != 0 {
		t.Errorf("got %d cached shapes with no drawing part, want 0", got)
	}
}
