package pptxxml

import (
	"embed"
	"fmt"
	"regexp"
	"strings"
)

//go:embed templates/smartart/*.xml templates/smartart/layouts/*/*.xml
var smartArtTemplateFS embed.FS

const (
	flattenSmartArtTextsInitCap = 8
	verifierLCGMultiplier       = int64(1664525)
	verifierLCGIncrement        = int64(1013904223)
	verifierNounIndexDivisor    = int64(256)
	verifierLCGModulus          = int64(1 << 32)
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

func injectSmartArtNodeTextsForModelIDs(data string, modelIDs []string, texts []string) string {
	if len(modelIDs) == 0 {
		return clearSmartArtPlaceholderTextRuns(data)
	}
	textByModelID := map[string]string{}
	textIndex := 0
	for _, modelID := range modelIDs {
		if textIndex >= len(texts) {
			break
		}
		textByModelID[modelID] = texts[textIndex]
		textIndex++
	}
	if textIndex < len(texts) {
		for _, modelID := range placeholderDataModelIDsInOrder(data) {
			if textIndex >= len(texts) {
				break
			}
			if _, exists := textByModelID[modelID]; exists {
				continue
			}
			textByModelID[modelID] = texts[textIndex]
			textIndex++
		}
	}
	segments := strings.Split(data, "<dgm:pt ")
	if len(segments) <= 1 {
		return clearSmartArtPlaceholderTextRuns(data)
	}
	var b strings.Builder
	b.WriteString(segments[0])
	for i := 1; i < len(segments); i++ {
		segment := "<dgm:pt " + segments[i]
		modelID := extractXMLAttr(segment, "modelId")
		if text, ok := textByModelID[modelID]; ok {
			segment = injectTextIntoPointSegment(segment, text)
		}
		b.WriteString(segment)
	}
	return clearSmartArtPlaceholderTextRuns(b.String())
}

func placeholderTextForIndex(texts []string, idx int) string {
	if idx < len(texts) {
		return texts[idx]
	}
	return ""
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

	v := ((int64(idx)+1)*verifierLCGMultiplier + verifierLCGIncrement) % verifierLCGModulus
	if v < 0 {
		v += verifierLCGModulus
	}
	a := adjectives[v%int64(len(adjectives))]
	n := nouns[(v/verifierNounIndexDivisor)%int64(len(nouns))]
	return fmt.Sprintf("%s-%s-%02d", a, n, idx+1)
}

func injectTextIntoPointSegment(segment, text string) string {
	if text != "" {
		segment = strings.Replace(segment, ` phldr="1"`, "", 1)
	}

	escaped := Escape(text)
	runXML := injectedSmartArtRunXML(escaped)

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
		"<a:p>"+runXML+"<a:endParaRPr",
		1,
	)
	if withRun != segment {
		return withRun
	}

	if endParaIdx := strings.Index(segment, "<a:endParaRPr"); endParaIdx >= 0 {
		return segment[:endParaIdx] + runXML + segment[endParaIdx:]
	}

	return strings.Replace(
		segment,
		"<a:p/>",
		"<a:p>"+runXML+"</a:p>",
		1,
	)
}

func injectedSmartArtRunXML(escapedText string) string {
	return `<a:r><a:rPr lang="en-US" sz="1800"/><a:t>` + escapedText + `</a:t></a:r>`
}

func clearSmartArtPlaceholderTextRuns(xml string) string {
	return placeholderTextRunPattern.ReplaceAllString(xml, "<a:t></a:t>")
}

func injectSmartArtDrawingTexts(
	drawing string,
	textByModelID map[string]string,
	hiddenPlaceholderModels map[string]struct{},
	allowedDrawingModels map[string]struct{},
) string {
	matches := drawingShapePattern.FindAllStringSubmatchIndex(drawing, -1)
	if len(matches) == 0 {
		return clearSmartArtPlaceholderTextRuns(drawing)
	}

	var b strings.Builder
	last := 0
	for _, idx := range matches {
		start := idx[0]
		end := idx[1]
		modelID := drawing[idx[2]:idx[3]]

		b.WriteString(drawing[last:start])
		shape := drawing[start:end]
		if allowedDrawingModels != nil {
			if _, allowed := allowedDrawingModels[modelID]; !allowed {
				last = end
				continue
			}
		}
		if text, ok := textByModelID[modelID]; ok {
			shape = injectTextIntoDrawingShape(shape, text)
		} else if _, hide := hiddenPlaceholderModels[modelID]; hide {
			last = end
			continue
		} else if strings.Contains(shape, "[Text]") || (strings.Contains(shape, "<dsp:txBody>") && !drawingShapeHasNonEmptyText(shape)) {
			shape = injectTextIntoDrawingShape(shape, "")
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
	return xml
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
	placeholderTextRunPattern = regexp.MustCompile(`<a:t>\[Text(?: [0-9]+)?\]</a:t>`)
	drawingShapePattern       = regexp.MustCompile(`(?s)<dsp:sp modelId="([^"]+)".*?</dsp:sp>`)
	paragraphPattern          = regexp.MustCompile(`(?s)<a:p(?: [^>]*)?>.*?</a:p>`)
)
