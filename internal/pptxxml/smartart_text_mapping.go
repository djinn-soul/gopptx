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
	presNodes := parseSmartArtPresNodes(data)

	// Prefer one "best" presentation node per data node to avoid duplicate or
	// connector-only labels stealing visible text slots.
	for dataModelID, text := range textByDataModelID {
		if text == "" {
			continue
		}
		bestModelID := ""
		bestScore := -1
		for modelID, node := range presNodes {
			if node.presAssocID != dataModelID {
				continue
			}
			score := scoreSmartArtPresTextNode(node.presName)
			if score > bestScore {
				bestScore = score
				bestModelID = modelID
			}
		}
		if bestModelID != "" {
			textByDrawingModelID[bestModelID] = text
		}
	}

	// Keep a conservative presOf fallback for templates where presAssocID
	// coverage is incomplete.
	for _, link := range parseSmartArtPresOfLinks(data) {
		if text, ok := textByDataModelID[link.srcModelID]; ok && text != "" {
			if _, exists := textByDrawingModelID[link.destModelID]; !exists {
				dest := presNodes[link.destModelID]
				if scoreSmartArtPresTextNode(dest.presName) >= 0 {
					textByDrawingModelID[link.destModelID] = text
				}
			}
		}
	}

	return textByDrawingModelID
}

func mapOrderedTextsToPreferredPresNodes(data string, orderedTexts []string) map[string]string {
	out := map[string]string{}
	if len(orderedTexts) == 0 {
		return out
	}
	dataPointTypes := parseSmartArtDataPointTypes(data)
	orderedNodes := parseSmartArtPresNodesInOrder(data)
	textIdx := 0
	for _, node := range orderedNodes {
		if textIdx >= len(orderedTexts) {
			break
		}
		assocType := dataPointTypes[node.presAssocID]
		if isSmartArtStructuralDataType(assocType) {
			continue
		}
		// High-confidence text buckets for visible user-facing labels.
		if scoreSmartArtPresTextNode(node.presName) >= 7 {
			out[node.modelID] = orderedTexts[textIdx]
			textIdx++
		}
	}
	return out
}

func preferredDataModelIDsInOrder(data string) []string {
	dataPointTypes := parseSmartArtDataPointTypes(data)
	orderedNodes := parseSmartArtPresNodesInOrder(data)
	seen := map[string]struct{}{}
	out := make([]string, 0, len(orderedNodes))
	for _, node := range orderedNodes {
		if node.presAssocID == "" {
			continue
		}
		assocType := dataPointTypes[node.presAssocID]
		if isSmartArtStructuralDataType(assocType) {
			continue
		}
		if scoreSmartArtPresTextNode(node.presName) < 7 {
			continue
		}
		if _, exists := seen[node.presAssocID]; exists {
			continue
		}
		seen[node.presAssocID] = struct{}{}
		out = append(out, node.presAssocID)
	}
	return out
}

func parseSmartArtDataPointTypes(data string) map[string]string {
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
		out[modelID] = strings.ToLower(extractXMLAttr(pt, "type"))
	}
	return out
}

func isSmartArtStructuralDataType(ptType string) bool {
	switch ptType {
	case "doc", "partrans", "sibtrans":
		return true
	default:
		return false
	}
}

type smartArtPresNode struct {
	modelID     string
	presAssocID string
	presName    string
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

func parseSmartArtPresNodes(data string) map[string]smartArtPresNode {
	out := map[string]smartArtPresNode{}
	segments := strings.Split(data, "<dgm:pt ")
	for i := 1; i < len(segments); i++ {
		pt := "<dgm:pt " + segments[i]
		if !strings.Contains(pt, `type="pres"`) {
			continue
		}
		presModelID := extractXMLAttr(pt, "modelId")
		assocModelID := extractXMLAttr(pt, "presAssocID")
		if presModelID == "" {
			continue
		}
		out[presModelID] = smartArtPresNode{
			modelID:     presModelID,
			presAssocID: assocModelID,
			presName:    strings.ToLower(extractXMLAttr(pt, "presName")),
		}
	}
	return out
}

func parseSmartArtPresNodesInOrder(data string) []smartArtPresNode {
	out := make([]smartArtPresNode, 0)
	segments := strings.Split(data, "<dgm:pt ")
	for i := 1; i < len(segments); i++ {
		pt := "<dgm:pt " + segments[i]
		if !strings.Contains(pt, `type="pres"`) {
			continue
		}
		presModelID := extractXMLAttr(pt, "modelId")
		if presModelID == "" {
			continue
		}
		out = append(out, smartArtPresNode{
			modelID:     presModelID,
			presAssocID: extractXMLAttr(pt, "presAssocID"),
			presName:    strings.ToLower(extractXMLAttr(pt, "presName")),
		})
	}
	return out
}

func scoreSmartArtPresTextNode(name string) int {
	if name == "" {
		return 0
	}
	switch {
	case strings.Contains(name, "connector"):
		return -3
	case strings.Contains(name, "sibtrans"):
		return -2
	case strings.Contains(name, "linnode"), name == "sp":
		return -1
	case strings.Contains(name, "parenttext"), strings.Contains(name, "roottext"):
		return 10
	case strings.Contains(name, "node"):
		return 8
	case strings.Contains(name, "hierchild"), strings.Contains(name, "hierroot"):
		return 7
	case strings.Contains(name, "text"):
		return 6
	default:
		return 1
	}
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
