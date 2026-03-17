package pptxxml

import (
	"embed"
	"strings"
)

//go:embed templates/smartart/*.xml templates/smartart/layouts/*/*.xml
var smartArtTemplateFS embed.FS

const (
	flattenSmartArtTextsInitCap = 8
)

func renderSmartArtDataFromTemplate(spec SmartArtSpec) string {
	data := mustTemplate(templatePathForLayout(spec.LayoutURI, "data.xml"))
	data = strings.Replace(data,
		`loTypeId="urn:microsoft.com/office/officeart/2005/8/layout/default"`,
		`loTypeId="`+Escape(layoutURIOrDefault(spec.LayoutURI))+`"`,
		1,
	)
	data = strings.Replace(data,
		`qsTypeId="urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1"`,
		`qsTypeId="`+Escape(defaultQuickStyleID(spec.QuickStyleID))+`"`,
		1,
	)
	data = strings.Replace(data,
		`csTypeId="urn:microsoft.com/office/officeart/2005/8/colors/accent1_2"`,
		`csTypeId="`+Escape(defaultColorStyleID(spec.ColorStyleID))+`"`,
		1,
	)
	orderedTexts := flattenSmartArtNodeTexts(spec.Nodes)
	targetDataModelIDs := preferredDataModelIDsInOrder(data)
	if len(targetDataModelIDs) > 0 {
		data = injectSmartArtNodeTextsForModelIDs(data, targetDataModelIDs, orderedTexts)
	} else {
		data = injectSmartArtNodeTexts(data, orderedTexts)
	}
	if strings.Contains(spec.LayoutURI, "/orgChart1") {
		data = pruneUnusedOrgChartPlaceholderBranches(data)
	}
	return data
}

func renderSmartArtLayoutFromTemplate(layoutURI string) string {
	layout := mustTemplate(templatePathForLayout(layoutURI, "layout.xml"))
	return strings.Replace(layout,
		`uniqueId="urn:microsoft.com/office/officeart/2005/8/layout/default"`,
		`uniqueId="`+Escape(layoutURIOrDefault(layoutURI))+`"`,
		1,
	)
}

func renderSmartArtStyleFromTemplate(quickStyleID string) string {
	style := mustTemplate("templates/smartart/quickStyle.xml")
	return strings.Replace(style,
		`uniqueId="urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1"`,
		`uniqueId="`+Escape(defaultQuickStyleID(quickStyleID))+`"`,
		1,
	)
}

func renderSmartArtColorsFromTemplate(colorStyleID string) string {
	colors := mustTemplate("templates/smartart/colors.xml")
	return strings.Replace(colors,
		`uniqueId="urn:microsoft.com/office/officeart/2005/8/colors/accent1_2"`,
		`uniqueId="`+Escape(defaultColorStyleID(colorStyleID))+`"`,
		1,
	)
}

func renderSmartArtDrawingFromTemplate(spec SmartArtSpec) string {
	drawing := mustTemplate(templatePathForLayout(spec.LayoutURI, "drawing.xml"))
	data := renderSmartArtDataFromTemplate(spec)
	orderedTexts := flattenSmartArtNodeTexts(spec.Nodes)
	textByModelID := buildDrawingTextMapFromData(data)
	hiddenPlaceholderModels := unfilledPlaceholderPresModelIDs(data)
	var allowedDrawingModels map[string]struct{}
	if strings.Contains(spec.LayoutURI, "/orgChart1") {
		allowedDrawingModels = existingPresModelIDs(data)
	}
	if preferOrderedNodeMapping(spec.LayoutURI) {
		preferred := mapOrderedTextsToPreferredPresNodes(data, orderedTexts)
		if len(preferred) >= len(orderedTexts) && len(preferred) > 0 {
			textByModelID = preferred
		}
	}
	if len(textByModelID) == 0 && len(orderedTexts) > 0 {
		if preferred := mapOrderedTextsToPreferredPresNodes(data, orderedTexts); len(preferred) > 0 {
			textByModelID = preferred
		}
	}
	return injectSmartArtDrawingTexts(drawing, textByModelID, hiddenPlaceholderModels, allowedDrawingModels)
}

func preferOrderedNodeMapping(layoutURI string) bool {
	return strings.Contains(layoutURI, "/vList5")
}

func flattenSmartArtNodeTexts(nodes []SmartArtNodeSpec) []string {
	out := make([]string, 0, flattenSmartArtTextsInitCap)
	var walk func([]SmartArtNodeSpec)
	walk = func(items []SmartArtNodeSpec) {
		for _, n := range items {
			out = append(out, n.Text)
			if len(n.Children) > 0 {
				walk(n.Children)
			}
		}
	}
	walk(nodes)
	return out
}

func layoutURIOrDefault(uri string) string {
	if uri != "" {
		return uri
	}
	return "urn:microsoft.com/office/officeart/2005/8/layout/default"
}

func mustTemplate(path string) string {
	b, err := smartArtTemplateFS.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func templatePathForLayout(layoutURI, fileName string) string {
	if key, ok := layoutTemplateKey(layoutURI); ok {
		candidate := "templates/smartart/layouts/" + key + "/" + fileName
		if _, err := smartArtTemplateFS.ReadFile(candidate); err == nil {
			return candidate
		}
	}
	return "templates/smartart/" + fileName
}

func layoutTemplateKey(layoutURI string) (string, bool) {
	if key, ok := layoutTemplateKeyList(layoutURI); ok {
		return key, true
	}
	if key, ok := layoutTemplateKeyProcess(layoutURI); ok {
		return key, true
	}
	if key, ok := layoutTemplateKeyDiagram(layoutURI); ok {
		return key, true
	}
	return "", false
}

func layoutTemplateKeyList(layoutURI string) (string, bool) {
	switch layoutURI {
	case "urn:microsoft.com/office/officeart/2005/8/layout/default":
		return "basic_block_list", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/vList5":
		return "vertical_block_list", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/hList1":
		return "horizontal_bullet_list", true
	case "urn:microsoft.com/office/officeart/2008/layout/SquareAccentList":
		return "square_accent_list", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/hList2":
		return "picture_accent_list", true
	default:
		return "", false
	}
}

func layoutTemplateKeyProcess(layoutURI string) (string, bool) {
	switch layoutURI {
	case "urn:microsoft.com/office/officeart/2005/8/layout/process1":
		return "basic_process", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/process3":
		return "accent_process", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/hProcess4":
		return "alternating_flow", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/hProcess9":
		return "continuous_block_process", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/cycle2":
		return "basic_cycle", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/cycle1":
		return "text_cycle", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/cycle5":
		return "block_cycle", true
	default:
		return "", false
	}
}

func layoutTemplateKeyDiagram(layoutURI string) (string, bool) {
	switch layoutURI {
	case "urn:microsoft.com/office/officeart/2005/8/layout/orgChart1":
		return "org_chart", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/hierarchy1":
		return "hierarchy", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/hierarchy2":
		return "horizontal_hierarchy", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/venn1":
		return "basic_venn", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/venn3":
		return "linear_venn", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/venn2":
		return "stacked_venn", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/radial1":
		return "basic_radial", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/matrix3":
		return "basic_matrix", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/matrix1":
		return "titled_matrix", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/pyramid1":
		return "basic_pyramid", true
	case "urn:microsoft.com/office/officeart/2005/8/layout/pyramid3":
		return "inverted_pyramid", true
	case "urn:microsoft.com/office/officeart/2008/layout/PictureStrips":
		return "picture_strips", true
	case "urn:microsoft.com/office/officeart/2008/layout/PictureGrid":
		return "picture_grid", true
	default:
		return "", false
	}
}
