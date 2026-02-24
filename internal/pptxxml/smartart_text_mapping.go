package pptxxml

import "strings"

const smartArtPresOfLinksInitCap = 16

type smartArtPresOfLink struct {
	srcModelID  string
	destModelID string
}

func buildDrawingTextMapFromData(data string) map[string]string {
	textByDataModelID := parseSmartArtDataPointTexts(data)
	if len(textByDataModelID) == 0 {
		return map[string]string{}
	}

	textByDrawingModelID := map[string]string{}

	for presModelID, assocModelID := range parseSmartArtPresAssocIDs(data) {
		if text, ok := textByDataModelID[assocModelID]; ok && text != "" {
			textByDrawingModelID[presModelID] = text
		}
	}

	for _, link := range parseSmartArtPresOfLinks(data) {
		if text, ok := textByDataModelID[link.srcModelID]; ok && text != "" {
			if _, exists := textByDrawingModelID[link.destModelID]; !exists {
				textByDrawingModelID[link.destModelID] = text
			}
		}
	}

	return textByDrawingModelID
}

func parseSmartArtDataPointTexts(data string) map[string]string {
	out := map[string]string{}
	segments := strings.Split(data, "<dgm:pt ")
	for i := 1; i < len(segments); i++ {
		pt := "<dgm:pt " + segments[i]
		if strings.Contains(pt, `type="pres"`) {
			continue
		}
		modelID := extractXMLAttr(pt, "modelId")
		if modelID == "" {
			continue
		}
		textStart := strings.Index(pt, "<a:t>")
		if textStart < 0 {
			continue
		}
		textStart += len("<a:t>")
		textEnd := strings.Index(pt[textStart:], "</a:t>")
		if textEnd < 0 {
			continue
		}
		text := pt[textStart : textStart+textEnd]
		if text != "" {
			out[modelID] = text
		}
	}
	return out
}

func parseSmartArtPresAssocIDs(data string) map[string]string {
	out := map[string]string{}
	segments := strings.Split(data, "<dgm:pt ")
	for i := 1; i < len(segments); i++ {
		pt := "<dgm:pt " + segments[i]
		if !strings.Contains(pt, `type="pres"`) {
			continue
		}
		presModelID := extractXMLAttr(pt, "modelId")
		assocModelID := extractXMLAttr(pt, "presAssocID")
		if presModelID == "" || assocModelID == "" {
			continue
		}
		out[presModelID] = assocModelID
	}
	return out
}

func parseSmartArtPresOfLinks(data string) []smartArtPresOfLink {
	links := make([]smartArtPresOfLink, 0, smartArtPresOfLinksInitCap)
	segments := strings.Split(data, "<dgm:cxn ")
	for i := 1; i < len(segments); i++ {
		cxn := "<dgm:cxn " + segments[i]
		if extractXMLAttr(cxn, "type") != "presOf" {
			continue
		}
		srcModelID := extractXMLAttr(cxn, "srcId")
		destModelID := extractXMLAttr(cxn, "destId")
		if srcModelID == "" || destModelID == "" {
			continue
		}
		links = append(links, smartArtPresOfLink{
			srcModelID:  srcModelID,
			destModelID: destModelID,
		})
	}
	return links
}

func extractXMLAttr(segment, attr string) string {
	token := attr + `="`
	start := strings.Index(segment, token)
	if start < 0 {
		return ""
	}
	start += len(token)
	end := strings.Index(segment[start:], `"`)
	if end < 0 {
		return ""
	}
	return segment[start : start+end]
}
