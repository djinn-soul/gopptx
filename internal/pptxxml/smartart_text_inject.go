package pptxxml

import (
	"regexp"
	"strings"
)

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
