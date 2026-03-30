package pptxxml

import (
	"embed"
	"strings"
	"sync"
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
	orderedTexts := smartArtOrderedTextsForLayout(spec.LayoutURI, spec.Nodes)
	targetDataModelIDs := preferredDataModelIDsForLayout(spec.LayoutURI, data)
	if len(targetDataModelIDs) > 0 {
		data = injectSmartArtNodeTextsForModelIDs(data, targetDataModelIDs, orderedTexts)
	} else {
		data = injectSmartArtNodeTexts(data, orderedTexts)
	}
	data = pruneUnusedOrgChartPlaceholderBranches(data)
	return data
}

func renderSmartArtLayoutFromTemplate(layoutURI string) string {
	if v, ok := renderedLayoutCache.Load(layoutURI); ok {
		if s, ok := v.(string); ok {
			return s
		}
		panic("renderedLayoutCache contained non-string value")
	}
	layout := mustTemplate(templatePathForLayout(layoutURI, "layout.xml"))
	s := strings.Replace(layout,
		`uniqueId="urn:microsoft.com/office/officeart/2005/8/layout/default"`,
		`uniqueId="`+Escape(layoutURIOrDefault(layoutURI))+`"`,
		1,
	)
	renderedLayoutCache.Store(layoutURI, s)
	return s
}

func renderSmartArtStyleFromTemplate(quickStyleID string) string {
	if v, ok := renderedStyleCache.Load(quickStyleID); ok {
		if s, ok := v.(string); ok {
			return s
		}
		panic("renderedStyleCache contained non-string value")
	}
	style := mustTemplate("templates/smartart/quickStyle.xml")
	s := strings.Replace(style,
		`uniqueId="urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1"`,
		`uniqueId="`+Escape(defaultQuickStyleID(quickStyleID))+`"`,
		1,
	)
	renderedStyleCache.Store(quickStyleID, s)
	return s
}

func renderSmartArtColorsFromTemplate(colorStyleID string) string {
	if v, ok := renderedColorsCache.Load(colorStyleID); ok {
		if s, ok := v.(string); ok {
			return s
		}
		panic("renderedColorsCache contained non-string value")
	}
	colors := mustTemplate("templates/smartart/colors.xml")
	s := strings.Replace(colors,
		`uniqueId="urn:microsoft.com/office/officeart/2005/8/colors/accent1_2"`,
		`uniqueId="`+Escape(defaultColorStyleID(colorStyleID))+`"`,
		1,
	)
	renderedColorsCache.Store(colorStyleID, s)
	return s
}

func renderSmartArtDrawingFromTemplate(spec SmartArtSpec) string {
	drawing := mustTemplate(templatePathForLayout(spec.LayoutURI, "drawing.xml"))
	data := renderSmartArtDataFromTemplate(spec)
	orderedTexts := smartArtOrderedTextsForLayout(spec.LayoutURI, spec.Nodes)
	textByModelID := buildDrawingTextMapFromData(data)
	hiddenPlaceholderModels := unfilledPlaceholderPresModelIDs(data)
	allowedDrawingModels := existingPresModelIDs(data)
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

func layoutURIOrDefault(uri string) string {
	if uri != "" {
		return uri
	}
	return "urn:microsoft.com/office/officeart/2005/8/layout/default"
}

//nolint:gochecknoglobals // package-level cache for embedded template strings
var templateCache sync.Map

// renderedLayoutCache / renderedStyleCache / renderedColorsCache cache the final
// rendered XML for layout, style, and colors — keyed by URI/ID string.
// These renders are pure functions of their input, so caching is safe and
// eliminates repeated strings.Replace + allocs on repeated SmartArt insertions.
//
//nolint:gochecknoglobals // package-level render caches, never mutated after first Store
var (
	renderedLayoutCache sync.Map
	renderedStyleCache  sync.Map
	renderedColorsCache sync.Map
)

func mustTemplate(path string) string {
	if v, ok := templateCache.Load(path); ok {
		if s, ok := v.(string); ok {
			return s
		}
		panic("templateCache contained non-string value")
	}
	b, err := smartArtTemplateFS.ReadFile(path)
	if err != nil {
		panic(err)
	}
	s := string(b)
	templateCache.Store(path, s)
	return s
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
