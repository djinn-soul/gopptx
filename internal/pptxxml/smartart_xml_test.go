package pptxxml_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

func TestSmartArtDataXMLContainsOrderingAndDrawingLink(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI:    "urn:microsoft.com/office/officeart/2005/8/layout/vList2",
		ColorStyleID: "urn:microsoft.com/office/officeart/2005/8/colors/accent1_2",
		QuickStyleID: "urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{Text: "Node A"},
			{Text: "Node B"},
		},
	}

	xml := pptxxml.SmartArtDataXML(spec)

	cxnCount := strings.Count(xml, "<dgm:cxn ")
	if cxnCount == 0 {
		t.Fatal("expected at least one dgm:cxn entry")
	}
	if got := strings.Count(xml, `srcOrd="`); got != cxnCount {
		t.Fatalf("expected srcOrd on all dgm:cxn entries, got %d of %d", got, cxnCount)
	}
	if got := strings.Count(xml, `destOrd="`); got != cxnCount {
		t.Fatalf("expected destOrd on all dgm:cxn entries, got %d of %d", got, cxnCount)
	}
	if !strings.Contains(xml, `<dsp:dataModelExt`) {
		t.Fatal("expected SmartArt dataModelExt link in data XML")
	}
	if !strings.Contains(xml, `relId="rId6"`) {
		t.Fatal("expected dataModelExt relId=rId6 in data XML")
	}
}

func TestSmartArtLayoutXMLContainsLayoutNode(t *testing.T) {
	xml := pptxxml.SmartArtLayoutXML(
		"urn:microsoft.com/office/officeart/2005/8/layout/vList2",
		"list",
	)
	if !strings.Contains(xml, "<dgm:layoutNode") {
		t.Fatal("expected dgm:layoutNode in SmartArt layout XML")
	}
}

func TestSmartArtStyleXMLContainsStyleLabel(t *testing.T) {
	xml := pptxxml.SmartArtStyleXML(
		"urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1",
	)
	if !strings.Contains(xml, "<dgm:styleLbl") {
		t.Fatal("expected dgm:styleLbl in SmartArt style XML")
	}
}

func TestSmartArtDataXMLUsesOrgChartCategory(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/orgChart1",
	}

	xml := pptxxml.SmartArtDataXML(spec)

	if !strings.Contains(xml, `loTypeId="urn:microsoft.com/office/officeart/2005/8/layout/orgChart1"`) {
		t.Fatal("expected orgChart layout URI in data XML")
	}
	if !strings.Contains(xml, `loCatId="hierarchy"`) {
		t.Fatal("expected hierarchy category for orgChart layout")
	}
}

func TestSmartArtDataXMLInjectsNodeTextFromSpec(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/default",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{Text: "Alpha"},
			{Text: "Beta"},
		},
	}

	xml := pptxxml.SmartArtDataXML(spec)

	if !strings.Contains(xml, "<a:t>Alpha</a:t>") {
		t.Fatal("expected first node text injected in data XML")
	}
	if !strings.Contains(xml, "<a:t>Beta</a:t>") {
		t.Fatal("expected second node text injected in data XML")
	}
}

func TestSmartArtDrawingXMLInjectsNodeTextForTitledMatrix(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/matrix1",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{Text: "Matrix A"},
			{Text: "Matrix B"},
			{Text: "Matrix C"},
		},
	}

	xml := pptxxml.SmartArtDrawingXML(spec)

	if !strings.Contains(xml, "<a:t>Matrix A</a:t>") {
		t.Fatal("expected first matrix text injected into drawing XML")
	}
	if !strings.Contains(xml, "<a:t>Matrix B</a:t>") {
		t.Fatal("expected second matrix text injected into drawing XML")
	}
	if !strings.Contains(xml, "<a:t>Matrix C</a:t>") {
		t.Fatal("expected third matrix text injected into drawing XML")
	}
}

func TestSmartArtDrawingXMLClearsPlaceholderTextForVerticalBlockList(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/vList5",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{Text: "Vertical A"},
			{Text: "Vertical B"},
			{Text: "Vertical C"},
		},
	}

	xml := pptxxml.SmartArtDrawingXML(spec)

	if strings.Contains(xml, "[Text]") {
		t.Fatal("expected no literal [Text] placeholders in drawing XML")
	}
	if !strings.Contains(xml, "<a:t>Vertical A</a:t>") {
		t.Fatal("expected node text injected into drawing XML")
	}
}
