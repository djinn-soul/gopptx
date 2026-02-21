package pptxxml

import (
	"embed"
	"fmt"
	"regexp"
	"strings"
)

//go:embed templates/smartart/*.xml templates/smartart/layouts/*/*.xml
var smartArtTemplateFS embed.FS

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
	data = injectSmartArtNodeTexts(data, flattenSmartArtNodeTexts(spec.Nodes))
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
	return injectSmartArtDrawingTexts(drawing, buildDrawingTextMapFromData(data))
}

func flattenSmartArtNodeTexts(nodes []SmartArtNodeSpec) []string {
	out := make([]string, 0, 8)
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
	if key, ok := layoutTemplateKeyByURI[layoutURI]; ok {
		candidate := "templates/smartart/layouts/" + key + "/" + fileName
		if _, err := smartArtTemplateFS.ReadFile(candidate); err == nil {
			return candidate
		}
	}
	return "templates/smartart/" + fileName
}

func injectSmartArtNodeTexts(data string, texts []string) string {
	segments := strings.Split(data, "<dgm:pt ")
	if len(segments) <= 1 {
		return clearSmartArtPlaceholderTextRuns(data)
	}

	var b strings.Builder
	b.WriteString(segments[0])

	textIndex := 0
	for i := 1; i < len(segments); i++ {
		segment := "<dgm:pt " + segments[i]
		if strings.Contains(segment, `phldrT="[Text]"`) &&
			strings.Contains(segment, `phldr="1"`) {
			segment = injectTextIntoPointSegment(segment, placeholderTextForIndex(texts, textIndex))
			textIndex++
		}
		b.WriteString(segment)
	}

	return clearSmartArtPlaceholderTextRuns(b.String())
}

func placeholderTextForIndex(texts []string, idx int) string {
	if idx < len(texts) {
		return texts[idx]
	}
	return generatedVerifierText(idx)
}

func generatedVerifierText(idx int) string {
	adjectives := [...]string{
		"Amber", "Nova", "Rapid", "Bright", "Swift",
		"Calm", "Bold", "Clear", "Prime", "Sharp",
	}
	nouns := [...]string{
		"Falcon", "River", "Orbit", "Matrix", "Beacon",
		"Vertex", "Signal", "Pulse", "Summit", "Vector",
	}

	v := uint32(idx+1)*1664525 + 1013904223
	a := adjectives[v%uint32(len(adjectives))]
	n := nouns[(v>>8)%uint32(len(nouns))]
	return fmt.Sprintf("%s-%s-%02d", a, n, idx+1)
}

func injectTextIntoPointSegment(segment, text string) string {
	escaped := Escape(text)

	if strings.Contains(segment, "<a:t>") {
		start := strings.Index(segment, "<a:t>")
		end := strings.Index(segment[start:], "</a:t>")
		if start >= 0 && end >= 0 {
			endAbs := start + end
			return segment[:start+5] + escaped + segment[endAbs:]
		}
	}

	withRun := strings.Replace(
		segment,
		"<a:p><a:endParaRPr",
		"<a:p><a:r><a:t>"+escaped+"</a:t></a:r><a:endParaRPr",
		1,
	)
	if withRun != segment {
		return withRun
	}

	if endParaIdx := strings.Index(segment, "<a:endParaRPr"); endParaIdx >= 0 {
		return segment[:endParaIdx] + "<a:r><a:t>" + escaped + "</a:t></a:r>" + segment[endParaIdx:]
	}

	return strings.Replace(
		segment,
		"<a:p/>",
		"<a:p><a:r><a:t>"+escaped+"</a:t></a:r></a:p>",
		1,
	)
}

func clearSmartArtPlaceholderTextRuns(xml string) string {
	return placeholderTextRunPattern.ReplaceAllString(xml, "<a:t></a:t>")
}

func injectSmartArtDrawingTexts(drawing string, textByModelID map[string]string) string {
	matches := drawingShapePattern.FindAllStringSubmatchIndex(drawing, -1)
	if len(matches) == 0 {
		return clearSmartArtPlaceholderTextRuns(drawing)
	}

	var b strings.Builder
	last := 0
	fallbackIndex := 0
	for _, idx := range matches {
		start := idx[0]
		end := idx[1]
		modelID := drawing[idx[2]:idx[3]]

		b.WriteString(drawing[last:start])
		shape := drawing[start:end]
		if text, ok := textByModelID[modelID]; ok {
			shape = injectTextIntoDrawingShape(shape, text)
		} else if strings.Contains(shape, "[Text]") || (strings.Contains(shape, "<dsp:txBody>") && !drawingShapeHasNonEmptyText(shape)) {
			shape = injectTextIntoDrawingShape(shape, generatedVerifierText(fallbackIndex))
			fallbackIndex++
		}
		b.WriteString(shape)
		last = end
	}
	b.WriteString(drawing[last:])
	return clearSmartArtPlaceholderTextRuns(b.String())
}

func injectTextIntoDrawingShape(shape, text string) string {
	txBodyStart := strings.Index(shape, "<dsp:txBody>")
	txBodyEnd := strings.Index(shape, "</dsp:txBody>")
	if txBodyStart < 0 || txBodyEnd < 0 || txBodyEnd <= txBodyStart {
		return shape
	}
	txBodyEnd += len("</dsp:txBody>")

	txBody := shape[txBodyStart:txBodyEnd]
	paragraphMatch := paragraphPattern.FindStringIndex(txBody)
	if paragraphMatch == nil {
		return shape
	}

	paragraph := txBody[paragraphMatch[0]:paragraphMatch[1]]
	paragraph = injectTextIntoPointSegment(paragraph, text)
	txBody = txBody[:paragraphMatch[0]] + paragraph + txBody[paragraphMatch[1]:]
	txBody = fillPlaceholderTextRuns(txBody, text)

	return shape[:txBodyStart] + txBody + shape[txBodyEnd:]
}

func fillPlaceholderTextRuns(xml, text string) string {
	escaped := Escape(text)
	if strings.Contains(xml, "<a:t>[Text") {
		xml = placeholderTextRunPattern.ReplaceAllString(xml, "<a:t>"+escaped+"</a:t>")
	}
	return strings.ReplaceAll(xml, "<a:t></a:t>", "<a:t>"+escaped+"</a:t>")
}

func drawingShapeHasNonEmptyText(shape string) bool {
	textMatches := regexp.MustCompile(`(?s)<a:t>(.*?)</a:t>`).FindAllStringSubmatch(shape, -1)
	for _, m := range textMatches {
		if len(m) > 1 && strings.TrimSpace(m[1]) != "" {
			return true
		}
	}
	return false
}

var (
	placeholderTextRunPattern = regexp.MustCompile(`<a:t>\[Text(?: [0-9]+)?\]</a:t>`)           //nolint:gochecknoglobals
	drawingShapePattern       = regexp.MustCompile(`(?s)<dsp:sp modelId="([^"]+)".*?</dsp:sp>`) //nolint:gochecknoglobals
	paragraphPattern          = regexp.MustCompile(`(?s)<a:p(?: [^>]*)?>.*?</a:p>`)             //nolint:gochecknoglobals
)

//nolint:gochecknoglobals // static mapping for embedded template selection
var layoutTemplateKeyByURI = map[string]string{
	"urn:microsoft.com/office/officeart/2005/8/layout/default":        "basic_block_list",
	"urn:microsoft.com/office/officeart/2005/8/layout/vList5":         "vertical_block_list",
	"urn:microsoft.com/office/officeart/2005/8/layout/hList1":         "horizontal_bullet_list",
	"urn:microsoft.com/office/officeart/2008/layout/SquareAccentList": "square_accent_list",
	"urn:microsoft.com/office/officeart/2005/8/layout/hList2":         "picture_accent_list",
	"urn:microsoft.com/office/officeart/2005/8/layout/process1":       "basic_process",
	"urn:microsoft.com/office/officeart/2005/8/layout/process3":       "accent_process",
	"urn:microsoft.com/office/officeart/2005/8/layout/hProcess4":      "alternating_flow",
	"urn:microsoft.com/office/officeart/2005/8/layout/hProcess9":      "continuous_block_process",
	"urn:microsoft.com/office/officeart/2005/8/layout/cycle2":         "basic_cycle",
	"urn:microsoft.com/office/officeart/2005/8/layout/cycle1":         "text_cycle",
	"urn:microsoft.com/office/officeart/2005/8/layout/cycle5":         "block_cycle",
	"urn:microsoft.com/office/officeart/2005/8/layout/orgChart1":      "org_chart",
	"urn:microsoft.com/office/officeart/2005/8/layout/hierarchy1":     "hierarchy",
	"urn:microsoft.com/office/officeart/2005/8/layout/hierarchy2":     "horizontal_hierarchy",
	"urn:microsoft.com/office/officeart/2005/8/layout/venn1":          "basic_venn",
	"urn:microsoft.com/office/officeart/2005/8/layout/venn3":          "linear_venn",
	"urn:microsoft.com/office/officeart/2005/8/layout/venn2":          "stacked_venn",
	"urn:microsoft.com/office/officeart/2005/8/layout/radial1":        "basic_radial",
	"urn:microsoft.com/office/officeart/2005/8/layout/matrix3":        "basic_matrix",
	"urn:microsoft.com/office/officeart/2005/8/layout/matrix1":        "titled_matrix",
	"urn:microsoft.com/office/officeart/2005/8/layout/pyramid1":       "basic_pyramid",
	"urn:microsoft.com/office/officeart/2005/8/layout/pyramid3":       "inverted_pyramid",
	"urn:microsoft.com/office/officeart/2008/layout/PictureStrips":    "picture_strips",
	"urn:microsoft.com/office/officeart/2008/layout/PictureGrid":      "picture_grid",
}
