package pptxxml

import "strings"

func placeholderDataModelIDsInOrder(data string) []string {
	out := make([]string, 0)
	segments := strings.Split(data, "<dgm:pt ")
	for i := 1; i < len(segments); i++ {
		pt := "<dgm:pt " + segments[i]
		if strings.Contains(pt, `type="pres"`) {
			continue
		}
		if !strings.Contains(pt, `phldr="1"`) {
			continue
		}
		modelID := extractXMLAttr(pt, "modelId")
		if modelID == "" {
			continue
		}
		out = append(out, modelID)
	}
	return out
}

func unfilledPlaceholderPresModelIDs(data string) map[string]struct{} {
	out := make(map[string]struct{})
	hiddenDataPoints := hiddenDataPointModelIDs(data)
	presNodes := parseSmartArtPresNodes(data)
	for modelID, node := range presNodes {
		if node.presAssocID == "" {
			continue
		}
		if _, hidden := hiddenDataPoints[node.presAssocID]; hidden {
			out[modelID] = struct{}{}
		}
	}
	return out
}

func existingPresModelIDs(data string) map[string]struct{} {
	out := make(map[string]struct{})
	presNodes := parseSmartArtPresNodes(data)
	for modelID := range presNodes {
		out[modelID] = struct{}{}
	}
	return out
}

func hiddenDataPointModelIDs(data string) map[string]struct{} {
	out := make(map[string]struct{})
	placeholderFlags := dataPointPlaceholderFlags(data)
	for modelID, isPlaceholder := range placeholderFlags {
		if isPlaceholder {
			out[modelID] = struct{}{}
		}
	}
	for modelID := range hiddenTransitionPointIDsForUnfilledDestinations(data, out) {
		out[modelID] = struct{}{}
	}
	return out
}

func hiddenTransitionPointIDsForUnfilledDestinations(
	data string,
	unfilledPlaceholderDataPoints map[string]struct{},
) map[string]struct{} {
	out := make(map[string]struct{})
	segments := strings.Split(data, "<dgm:cxn ")
	for i := 1; i < len(segments); i++ {
		cxn := "<dgm:cxn " + segments[i]
		cxnType := strings.ToLower(extractXMLAttr(cxn, "type"))
		if strings.HasPrefix(cxnType, "pres") {
			continue
		}
		destModelID := extractXMLAttr(cxn, "destId")
		if _, hiddenDest := unfilledPlaceholderDataPoints[destModelID]; !hiddenDest {
			continue
		}
		parTransID := extractXMLAttr(cxn, "parTransId")
		if parTransID != "" {
			out[parTransID] = struct{}{}
		}
		sibTransID := extractXMLAttr(cxn, "sibTransId")
		if sibTransID != "" {
			out[sibTransID] = struct{}{}
		}
	}
	return out
}

func dataPointPlaceholderFlags(data string) map[string]bool {
	out := make(map[string]bool)
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
		out[modelID] = strings.Contains(pt, `phldr="1"`)
	}
	return out
}
