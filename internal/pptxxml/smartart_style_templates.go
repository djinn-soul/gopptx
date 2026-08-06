package pptxxml

import "strings"

// The quick style and colour style definitions PowerPoint writes for each of
// its styles ship with the package, one file per style, and are looked up by
// the last segment of the style URI.

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
