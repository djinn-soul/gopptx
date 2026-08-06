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
	if !strings.Contains(xml, "<a:t>Vertical B</a:t>") {
		t.Fatal("expected second node text injected into drawing XML")
	}
	if !strings.Contains(xml, "<a:t>Vertical C</a:t>") {
		t.Fatal("expected third node text injected into drawing XML")
	}
	if strings.Contains(xml, "Bright-") || strings.Contains(xml, "Prime-") {
		t.Fatal("unexpected verifier filler text leaked into drawing XML")
	}
}

func TestSmartArtDataXMLClearsPlaceholderFlagForInjectedText(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/process1",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{Text: "Phase 1"},
			{Text: "Phase 2"},
			{Text: "Phase 3"},
		},
	}

	xml := pptxxml.SmartArtDataXML(spec)

	for _, text := range []string{"Phase 1", "Phase 2", "Phase 3"} {
		segment := pointSegmentContainingText(xml, text)
		if segment == "" {
			t.Fatalf("expected data point segment for %q", text)
		}
		if strings.Contains(segment, `phldr="1"`) {
			t.Fatalf("expected injected text point %q to clear phldr=\"1\"", text)
		}
	}
}

func TestSmartArtDataXMLOrgChartDoesNotMapFirstNodeToDocRoot(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/orgChart1",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{
				Text: "CEO",
				Children: []pptxxml.SmartArtNodeSpec{
					{Text: "VP Sales"},
					{Text: "VP Eng"},
				},
			},
		},
	}

	xml := pptxxml.SmartArtDataXML(spec)
	segment := pointSegmentContainingText(xml, "CEO")
	if segment == "" {
		t.Fatal("expected CEO in SmartArt data XML")
	}
	if strings.Contains(segment, `type="doc"`) {
		t.Fatal("expected CEO text to map to content node, not doc root")
	}
}

func TestSmartArtDataXMLHorizontalBulletMapsAcrossPrimaryColumns(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/hList1",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{Text: "Col A"},
			{Text: "Col B"},
			{Text: "Col C"},
		},
	}

	xml := pptxxml.SmartArtDataXML(spec)

	if segment := pointSegmentByModelID(xml, "{40A0009A-D31E-4203-8A68-5BC7235C8A6B}"); !strings.Contains(
		segment,
		"<a:t>Col A</a:t>",
	) {
		t.Fatal("expected first horizontal-bullet text in first primary column node")
	}
	if segment := pointSegmentByModelID(xml, "{366D710F-3C7C-4530-80E1-1F49B4505CB6}"); !strings.Contains(
		segment,
		"<a:t>Col B</a:t>",
	) {
		t.Fatal("expected second horizontal-bullet text in second primary column node")
	}
	if segment := pointSegmentByModelID(xml, "{38041CF9-B0F0-4D98-8547-E9C385390400}"); !strings.Contains(
		segment,
		"<a:t>Col C</a:t>",
	) {
		t.Fatal("expected third horizontal-bullet text in third primary column node")
	}
}

func TestSmartArtDataXMLHierarchyMapsBreadthFirstAcrossSiblingSlots(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/hierarchy1",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{
				Text: "CEO",
				Children: []pptxxml.SmartArtNodeSpec{
					{
						Text: "Finance",
						Children: []pptxxml.SmartArtNodeSpec{
							{Text: "Accounts"},
						},
					},
					{Text: "Engineering"},
				},
			},
		},
	}

	xml := pptxxml.SmartArtDataXML(spec)

	if segment := pointSegmentByModelID(
		xml,
		"{C7401706-CA2B-4D1A-BA1C-F52A374534A6}",
	); !strings.Contains(
		segment,
		"<a:t>CEO</a:t>",
	) {
		t.Fatal("expected hierarchy root text in root slot")
	}
	if segment := pointSegmentByModelID(
		xml,
		"{BF0DB955-A3BC-42E2-B6F3-76BB6EEF0C6E}",
	); !strings.Contains(
		segment,
		"<a:t>Finance</a:t>",
	) {
		t.Fatal("expected first child text in first hierarchy child slot")
	}
	if segment := pointSegmentByModelID(
		xml,
		"{1D8E8D7D-7330-4A43-BF16-A35AFE48D68D}",
	); !strings.Contains(
		segment,
		"<a:t>Engineering</a:t>",
	) {
		t.Fatal("expected second child text in second hierarchy child slot")
	}
	if segment := pointSegmentByModelID(
		xml,
		"{CDF06D8C-2EB0-4B55-84B0-DB4ADA208865}",
	); !strings.Contains(
		segment,
		"<a:t>Accounts</a:t>",
	) {
		t.Fatal("expected grandchild text in first hierarchy grandchild slot")
	}
}

func TestSmartArtDataXMLHorizontalHierarchyMapsBreadthFirstAcrossSiblingSlots(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/hierarchy2",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{
				Text: "CEO",
				Children: []pptxxml.SmartArtNodeSpec{
					{
						Text: "Finance",
						Children: []pptxxml.SmartArtNodeSpec{
							{Text: "Accounts"},
						},
					},
					{Text: "Engineering"},
				},
			},
		},
	}

	xml := pptxxml.SmartArtDataXML(spec)

	if segment := pointSegmentByModelID(
		xml,
		"{8426535C-54F3-46EA-A53B-BA4CC18FD7E0}",
	); !strings.Contains(
		segment,
		"<a:t>CEO</a:t>",
	) {
		t.Fatal("expected horizontal hierarchy root text in root slot")
	}
	if segment := pointSegmentByModelID(
		xml,
		"{CE0AAF88-E00D-4698-B726-CFDB9A8B02FB}",
	); !strings.Contains(
		segment,
		"<a:t>Finance</a:t>",
	) {
		t.Fatal("expected first semantic child text in first horizontal hierarchy child slot")
	}
	if segment := pointSegmentByModelID(
		xml,
		"{1AEDA9FF-FEC4-4109-BC1E-C7FF7AA02EA0}",
	); !strings.Contains(
		segment,
		"<a:t>Engineering</a:t>",
	) {
		t.Fatal("expected second semantic child text in second horizontal hierarchy child slot")
	}
	if segment := pointSegmentByModelID(
		xml,
		"{88A38E3E-B1A9-4D56-B46C-71706701582B}",
	); !strings.Contains(
		segment,
		"<a:t>Accounts</a:t>",
	) {
		t.Fatal("expected grandchild text in first horizontal hierarchy grandchild slot")
	}
}

func TestSmartArtDrawingXMLOrgChartHidesUnfilledPlaceholderShapes(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/orgChart1",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{
				Text: "CEO",
				Children: []pptxxml.SmartArtNodeSpec{
					{Text: "Finance"},
					{Text: "Engineering"},
				},
			},
		},
	}

	drawing := pptxxml.SmartArtDrawingXML(spec)

	if strings.Contains(drawing, `modelId="{02DEC6B8-6043-47C0-A1CA-36E0D9C9344A}"`) {
		t.Fatal("expected assistant placeholder shape to be omitted when no assistant text is provided")
	}
	if strings.Contains(drawing, `modelId="{0F1F6AAE-361E-4C23-8CC3-9AADE342BC57}"`) {
		t.Fatal("expected unused third child placeholder shape to be omitted")
	}
	if strings.Contains(drawing, `modelId="{F3D5C6F3-90CA-45EE-AA32-5EE9940677BF}"`) {
		t.Fatal("expected assistant connector branch shape to be omitted")
	}
	if strings.Contains(drawing, `modelId="{7552B3CA-53DB-4D69-9562-354F6F7EBFAA}"`) {
		t.Fatal("expected unused third-child connector branch shape to be omitted")
	}
	if !strings.Contains(drawing, "<a:t>CEO</a:t>") {
		t.Fatal("expected root text in org chart drawing")
	}
	if !strings.Contains(drawing, "<a:t>Finance</a:t>") {
		t.Fatal("expected first child text in org chart drawing")
	}
	if !strings.Contains(drawing, "<a:t>Engineering</a:t>") {
		t.Fatal("expected second child text in org chart drawing")
	}
}

func TestSmartArtDataXMLOrgChartPrunesUnusedAssistantAndChildBranches(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/orgChart1",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{
				Text: "CEO",
				Children: []pptxxml.SmartArtNodeSpec{
					{Text: "Finance"},
					{Text: "Engineering"},
				},
			},
		},
	}

	data := pptxxml.SmartArtDataXML(spec)

	for _, removed := range []string{
		`modelId="{197A29FD-B529-4802-9A22-4AB325650FFC}"`, // assistant data node
		`modelId="{7B78C063-41AF-46CB-9560-1E5EFBA80343}"`, // unused third child data node
		`modelId="{02DEC6B8-6043-47C0-A1CA-36E0D9C9344A}"`, // assistant text pres node
		`modelId="{0F1F6AAE-361E-4C23-8CC3-9AADE342BC57}"`, // unused child text pres node
	} {
		if strings.Contains(data, removed) {
			t.Fatalf("expected %s to be pruned from org chart data", removed)
		}
	}
}

// A radial diagram's hub is a text slot like any other. When it was skipped it
// stayed an unfilled placeholder, the pruner deleted it, and deleting it deleted
// every parent/child connection with it — PowerPoint then drew a single node and
// silently dropped the rest of the items.
func TestSmartArtDataXMLRadialFillsHubAndKeepsDataTree(t *testing.T) {
	spec := pptxxml.SmartArtSpec{
		LayoutURI: "urn:microsoft.com/office/officeart/2005/8/layout/radial1",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{Text: "North"},
			{Text: "East"},
			{Text: "South"},
			{Text: "West"},
		},
	}

	data := pptxxml.SmartArtDataXML(spec)

	if segment := pointSegmentByModelID(
		data,
		"{FF425370-D8F1-471A-878E-62CB5FEA877F}",
	); !strings.Contains(segment, "<a:t>North</a:t>") {
		t.Fatal("expected first radial text in the hub slot")
	}
	for _, text := range []string{"East", "South", "West"} {
		if !strings.Contains(data, "<a:t>"+text+"</a:t>") {
			t.Fatalf("expected radial satellite text %q to survive", text)
		}
	}
	if structuralConnectionCount(data) == 0 {
		t.Fatal("expected the radial data tree to keep its parent/child connections")
	}
}

// PowerPoint resolves a style by category and type together, and trusts a cached
// drawing the data model still points at. Both used to leave a 3-D quick style
// rendering as flat boxes.
func TestSmartArtQuickStyleReachesTheSlide(t *testing.T) {
	const threeD = "urn:microsoft.com/office/officeart/2005/8/quickstyle/3d1"

	style := pptxxml.SmartArtStyleXML(threeD)
	if !strings.Contains(style, "bevel") {
		t.Fatal("expected the 3d1 style definition to carry its bevels")
	}
	if style == pptxxml.SmartArtStyleXML("") {
		t.Fatal("expected 3d1 to differ from the default quick style")
	}

	spec := pptxxml.SmartArtSpec{
		LayoutURI:    "urn:microsoft.com/office/officeart/2005/8/layout/process1",
		QuickStyleID: threeD,
		ColorStyleID: "urn:microsoft.com/office/officeart/2005/8/colors/colorful1",
		Nodes: []pptxxml.SmartArtNodeSpec{
			{Text: "Plan"}, {Text: "Build"}, {Text: "Ship"},
		},
	}
	data := pptxxml.SmartArtDataXML(spec)

	if !strings.Contains(data, `qsCatId="3D"`) {
		t.Error("expected the 3-D quick style to carry its own category")
	}
	if !strings.Contains(data, `csCatId="colorful"`) {
		t.Error("expected the colourful colour style to carry its own category")
	}
	if strings.Contains(data, "<dgm:extLst>") {
		t.Error("expected the stale drawing cache to be unhooked so PowerPoint relays out")
	}

	// The default style is what the cached drawing was captured under, so it
	// keeps the cache.
	spec.QuickStyleID = ""
	if !strings.Contains(pptxxml.SmartArtDataXML(spec), "<dgm:extLst>") {
		t.Error("expected the default quick style to keep the drawing cache")
	}
}

func TestSmartArtColorStyleDefinitionsUseTheirOwnAccents(t *testing.T) {
	accent4 := pptxxml.SmartArtColorsXML("urn:microsoft.com/office/officeart/2005/8/colors/accent4_2")
	if !strings.Contains(accent4, `<a:schemeClr val="accent4"`) {
		t.Fatal("expected the accent4 colour style to be defined against accent4")
	}
	colorful := pptxxml.SmartArtColorsXML("urn:microsoft.com/office/officeart/2005/8/colors/colorful1")
	if colorful == accent4 {
		t.Fatal("expected a colourful colour style to differ from an accent one")
	}
}

// structuralConnectionCount counts parent/child connections, which carry no type
// attribute, unlike the presOf/presParOf connections of the layout cache.
func structuralConnectionCount(xml string) int {
	count := 0
	for _, segment := range strings.Split(xml, "<dgm:cxn ")[1:] {
		attrs, _, found := strings.Cut(segment, ">")
		if !found {
			continue
		}
		if !strings.Contains(attrs, " type=") {
			count++
		}
	}
	return count
}

func pointSegmentContainingText(xml, text string) string {
	needle := "<a:t>" + text + "</a:t>"
	segments := strings.Split(xml, "<dgm:pt ")
	for i := 1; i < len(segments); i++ {
		segment := "<dgm:pt " + segments[i]
		if strings.Contains(segment, needle) {
			return segment
		}
	}
	return ""
}

func pointSegmentByModelID(xml, modelID string) string {
	needle := `modelId="` + modelID + `"`
	segments := strings.Split(xml, "<dgm:pt ")
	for i := 1; i < len(segments); i++ {
		segment := "<dgm:pt " + segments[i]
		if strings.Contains(segment, needle) {
			return segment
		}
	}
	return ""
}
