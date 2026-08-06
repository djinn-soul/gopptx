package pptxxml

import (
	"embed"
	"strings"
	"sync"
)

//go:embed templates/smartart/*.xml templates/smartart/layouts/*/*.xml
//go:embed templates/smartart/quickstyles/*.xml templates/smartart/colorstyles/*.xml
var smartArtTemplateFS embed.FS

const (
	flattenSmartArtTextsInitCap = 8
)

func renderSmartArtDataFromTemplate(spec SmartArtSpec) string {
	data := mustTemplate(templatePathForLayout(spec.LayoutURI, "data.xml"))
	if !smartArtSpecFitsTemplate(spec, data) {
		return renderGeneratedSmartArtData(spec)
	}
	data = strings.Replace(data,
		`loTypeId="urn:microsoft.com/office/officeart/2005/8/layout/default"`,
		`loTypeId="`+Escape(layoutURIOrDefault(spec.LayoutURI))+`"`,
		1,
	)
	quickStyleID := defaultQuickStyleID(spec.QuickStyleID)
	colorStyleID := defaultColorStyleID(spec.ColorStyleID)
	data = strings.Replace(data,
		`qsTypeId="urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1"`,
		`qsTypeId="`+Escape(quickStyleID)+`"`,
		1,
	)
	data = strings.Replace(data,
		`csTypeId="urn:microsoft.com/office/officeart/2005/8/colors/accent1_2"`,
		`csTypeId="`+Escape(colorStyleID)+`"`,
		1,
	)
	// The categories are part of how PowerPoint resolves a style, so they have to
	// follow the style rather than stay on the template's own.
	data = strings.Replace(data,
		`qsCatId="simple"`,
		`qsCatId="`+Escape(smartArtQuickStyleCategory(quickStyleID))+`"`,
		1,
	)
	data = strings.Replace(data,
		`csCatId="accent1"`,
		`csCatId="`+Escape(smartArtColorCategory(colorStyleID))+`"`,
		1,
	)
	orderedNodes := smartArtOrderedNodesForLayout(spec.LayoutURI, spec.Nodes)
	orderedTexts := smartArtOrderedTextsForLayout(spec.LayoutURI, spec.Nodes)
	targetDataModelIDs := preferredDataModelIDsForLayout(spec.LayoutURI, data)

	// A nested spec is placed by walking both trees together; the slot heuristics
	// only know how to fill a flat list.
	if slots := assignSmartArtSlots(spec, data); smartArtSpecHasChildren(spec.Nodes) && len(slots) > 0 {
		orderedNodes = smartArtSlotNodes(slots)
		orderedTexts = smartArtSlotTexts(slots)
		targetDataModelIDs = smartArtSlotModelIDs(slots)
	}

	if len(targetDataModelIDs) > 0 {
		data = injectSmartArtNodeTextsForModelIDs(data, targetDataModelIDs, orderedTexts)
	} else {
		targetDataModelIDs = placeholderDataModelIDsInOrder(data)
		data = injectSmartArtNodeTexts(data, orderedTexts)
	}
	data = applySmartArtNodeProperties(
		data,
		smartArtNodePropertiesByModelID(orderedNodes, targetDataModelIDs),
		smartArtPictureShapeNames(spec.LayoutURI),
	)
	data = pruneUnusedOrgChartPlaceholderBranches(data)
	return dropSmartArtDrawingCacheLink(data)
}

// dropSmartArtDrawingCacheLink unhooks the cached drawing from the data model.
//
// PowerPoint trusts a cache the data model still points at, and the shipped
// caches were captured from the templates: their shapes carry the template's
// quick style and the text size that suited its placeholder captions. That cache
// drew 3-D styles flat, and split real captions across lines mid-word rather
// than shrinking them to fit. Unhooked, PowerPoint lays the
// diagram out from the data, layout and style definitions and sizes the text
// itself. The drawing part stays in the package for readers that render the
// cache rather than recompute it.
func dropSmartArtDrawingCacheLink(data string) string {
	start := strings.Index(data, "<dgm:extLst>")
	if start < 0 {
		return data
	}
	end := strings.Index(data[start:], "</dgm:extLst>")
	if end < 0 {
		return data
	}
	return data[:start] + data[start+end+len("</dgm:extLst>"):]
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
	s := smartArtStyleDefinition(defaultQuickStyleID(quickStyleID))
	renderedStyleCache.Store(quickStyleID, s)
	return s
}

// smartArtStyleDefinition returns the style definition PowerPoint itself writes
// for a quick style. The effects that separate one quick style from another —
// 3-D scenes, bevels, shadows — live in these definitions, so a style ID
// stamped onto the shipped simple1 body used to change nothing on the slide.
func smartArtStyleDefinition(quickStyleID string) string {
	if xml, ok := smartArtStyleVariant("quickstyles", quickStyleID); ok {
		return xml
	}
	style := mustTemplate("templates/smartart/quickStyle.xml")
	return strings.Replace(style,
		`uniqueId="urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1"`,
		`uniqueId="`+Escape(quickStyleID)+`"`,
		1,
	)
}

func renderSmartArtColorsFromTemplate(colorStyleID string) string {
	if v, ok := renderedColorsCache.Load(colorStyleID); ok {
		if s, ok := v.(string); ok {
			return s
		}
		panic("renderedColorsCache contained non-string value")
	}
	s := smartArtColorsDefinition(defaultColorStyleID(colorStyleID))
	renderedColorsCache.Store(colorStyleID, s)
	return s
}

// smartArtColorsDefinition returns the colour definition PowerPoint itself
// writes for a colour style. Which accents a style uses, and how it cycles them
// across nodes, is described here rather than by the style ID.
func smartArtColorsDefinition(colorStyleID string) string {
	if xml, ok := smartArtStyleVariant("colorstyles", colorStyleID); ok {
		return xml
	}
	colors := mustTemplate("templates/smartart/colors.xml")
	return strings.Replace(colors,
		`uniqueId="urn:microsoft.com/office/officeart/2005/8/colors/accent1_2"`,
		`uniqueId="`+Escape(colorStyleID)+`"`,
		1,
	)
}

// smartArtStyleVariant loads the definition whose file is named after the last
// segment of the style URI, e.g. ".../quickstyle/3d1" -> "quickstyles/3d1.xml".
func smartArtStyleVariant(dir, styleID string) (string, bool) {
	name := styleID[strings.LastIndex(styleID, "/")+1:]
	if name == "" || strings.ContainsAny(name, `\/.`) {
		return "", false
	}
	path := "templates/smartart/" + dir + "/" + name + ".xml"
	b, err := smartArtTemplateFS.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func renderSmartArtDrawingFromTemplate(spec SmartArtSpec) string {
	drawing := mustTemplate(templatePathForLayout(spec.LayoutURI, "drawing.xml"))
	templateData := mustTemplate(templatePathForLayout(spec.LayoutURI, "data.xml"))
	if !smartArtSpecFitsTemplate(spec, templateData) {
		// The cached drawing describes the template's shape, which this diagram
		// has outgrown. PowerPoint lays it out from the data model instead; the
		// stale cache is left as-is rather than filled with the wrong captions.
		return clearSmartArtPlaceholderTextRuns(drawing)
	}
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
	return smartArtDefaultLayoutURN
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
	case smartArtDefaultLayoutURN:
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
