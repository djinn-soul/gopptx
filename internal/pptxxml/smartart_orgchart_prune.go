package pptxxml

import (
	"regexp"
	"strings"
)

var (
	smartArtDataPointPattern = regexp.MustCompile(`(?s)<dgm:pt [^>]*>.*?</dgm:pt>`)
	smartArtDataCxnPattern   = regexp.MustCompile(`<dgm:cxn [^>]*/>`)
	smartArtTextRunPattern   = regexp.MustCompile(`(?s)<a:t>(.*?)</a:t>`)
)

const minSmartArtTextSubmatches = 2

type smartArtDataPoint struct {
	raw           string
	modelID       string
	pointType     string
	presAssocID   string
	isPres        bool
	isPlaceholder bool
	hasText       bool
}

type smartArtDataCxn struct {
	raw        string
	cxnType    string
	srcID      string
	destID     string
	parTransID string
	sibTransID string
}

func pruneUnusedOrgChartPlaceholderBranches(data string) string {
	points := parseSmartArtDataPoints(data)
	if len(points) == 0 {
		return data
	}
	cxns := parseSmartArtDataConnections(data)
	removedDataIDs := initialOrgChartRemovedDataPointIDs(points)
	if len(removedDataIDs) == 0 {
		return data
	}
	expandRemovedOrgChartTransitionIDs(removedDataIDs, cxns)
	removedPresIDs := removedOrgChartPresentationPointIDs(points, removedDataIDs, cxns)

	keptPoints := filterSmartArtDataPoints(points, removedDataIDs, removedPresIDs)
	keptCxns := filterSmartArtDataConnections(cxns, removedDataIDs, removedPresIDs)
	return rewriteSmartArtDataPointsAndConnections(data, keptPoints, keptCxns)
}

func parseSmartArtDataPoints(data string) []smartArtDataPoint {
	matches := smartArtDataPointPattern.FindAllString(data, -1)
	out := make([]smartArtDataPoint, 0, len(matches))
	for _, raw := range matches {
		pointType := strings.ToLower(extractXMLAttr(raw, "type"))
		isPres := pointType == "pres"
		out = append(out, smartArtDataPoint{
			raw:           raw,
			modelID:       extractXMLAttr(raw, "modelId"),
			pointType:     pointType,
			presAssocID:   extractXMLAttr(raw, "presAssocID"),
			isPres:        isPres,
			isPlaceholder: strings.Contains(raw, `phldr="1"`),
			hasText:       smartArtPointHasText(raw),
		})
	}
	return out
}

func parseSmartArtDataConnections(data string) []smartArtDataCxn {
	matches := smartArtDataCxnPattern.FindAllString(data, -1)
	out := make([]smartArtDataCxn, 0, len(matches))
	for _, raw := range matches {
		out = append(out, smartArtDataCxn{
			raw:        raw,
			cxnType:    strings.ToLower(extractXMLAttr(raw, "type")),
			srcID:      extractXMLAttr(raw, "srcId"),
			destID:     extractXMLAttr(raw, "destId"),
			parTransID: extractXMLAttr(raw, "parTransId"),
			sibTransID: extractXMLAttr(raw, "sibTransId"),
		})
	}
	return out
}

func initialOrgChartRemovedDataPointIDs(points []smartArtDataPoint) map[string]struct{} {
	out := make(map[string]struct{})
	for _, point := range points {
		if point.isPres || point.modelID == "" {
			continue
		}
		if !point.isPlaceholder || point.hasText {
			continue
		}
		if point.pointType == "doc" {
			continue
		}
		out[point.modelID] = struct{}{}
	}
	return out
}

func expandRemovedOrgChartTransitionIDs(
	removedDataIDs map[string]struct{},
	cxns []smartArtDataCxn,
) {
	changed := true
	for changed {
		changed = false
		for _, cxn := range cxns {
			if strings.HasPrefix(cxn.cxnType, "pres") {
				continue
			}
			if _, removeDest := removedDataIDs[cxn.destID]; !removeDest {
				continue
			}
			if addRemovedDataID(removedDataIDs, cxn.parTransID) {
				changed = true
			}
			if addRemovedDataID(removedDataIDs, cxn.sibTransID) {
				changed = true
			}
		}
	}
}

func addRemovedDataID(removed map[string]struct{}, id string) bool {
	if id == "" {
		return false
	}
	if _, exists := removed[id]; exists {
		return false
	}
	removed[id] = struct{}{}
	return true
}

//nolint:gocognit // Removal propagation uses explicit branch checks to preserve deterministic pruning behavior.
func removedOrgChartPresentationPointIDs(
	points []smartArtDataPoint,
	removedDataIDs map[string]struct{},
	cxns []smartArtDataCxn,
) map[string]struct{} {
	removedPresIDs := make(map[string]struct{})
	for _, point := range points {
		if !point.isPres || point.modelID == "" {
			continue
		}
		if _, remove := removedDataIDs[point.presAssocID]; remove {
			removedPresIDs[point.modelID] = struct{}{}
		}
	}
	changed := true
	for changed {
		changed = false
		for _, cxn := range cxns {
			if cxn.cxnType != "presparof" {
				continue
			}
			if _, remove := removedPresIDs[cxn.srcID]; !remove {
				continue
			}
			if cxn.destID == "" {
				continue
			}
			if _, exists := removedPresIDs[cxn.destID]; exists {
				continue
			}
			removedPresIDs[cxn.destID] = struct{}{}
			changed = true
		}
	}
	return removedPresIDs
}

func filterSmartArtDataPoints(
	points []smartArtDataPoint,
	removedDataIDs map[string]struct{},
	removedPresIDs map[string]struct{},
) []string {
	kept := make([]string, 0, len(points))
	for _, point := range points {
		if point.isPres {
			if _, remove := removedPresIDs[point.modelID]; remove {
				continue
			}
		} else if _, remove := removedDataIDs[point.modelID]; remove {
			continue
		}
		kept = append(kept, point.raw)
	}
	return kept
}

func filterSmartArtDataConnections(
	cxns []smartArtDataCxn,
	removedDataIDs map[string]struct{},
	removedPresIDs map[string]struct{},
) []string {
	kept := make([]string, 0, len(cxns))
	for _, cxn := range cxns {
		if strings.HasPrefix(cxn.cxnType, "pres") {
			if _, remove := removedPresIDs[cxn.srcID]; remove {
				continue
			}
			if _, remove := removedPresIDs[cxn.destID]; remove {
				continue
			}
			kept = append(kept, cxn.raw)
			continue
		}
		if _, remove := removedDataIDs[cxn.srcID]; remove {
			continue
		}
		if _, remove := removedDataIDs[cxn.destID]; remove {
			continue
		}
		if _, remove := removedDataIDs[cxn.parTransID]; remove {
			continue
		}
		if _, remove := removedDataIDs[cxn.sibTransID]; remove {
			continue
		}
		kept = append(kept, cxn.raw)
	}
	return kept
}

func rewriteSmartArtDataPointsAndConnections(data string, points []string, cxns []string) string {
	pointStart := strings.Index(data, "<dgm:ptLst>")
	pointEnd := strings.Index(data, "</dgm:ptLst>")
	if pointStart < 0 || pointEnd < 0 || pointEnd <= pointStart {
		return data
	}
	pointStart += len("<dgm:ptLst>")
	data = data[:pointStart] + strings.Join(points, "") + data[pointEnd:]

	cxnStart := strings.Index(data, "<dgm:cxnLst>")
	cxnEnd := strings.Index(data, "</dgm:cxnLst>")
	if cxnStart < 0 || cxnEnd < 0 || cxnEnd <= cxnStart {
		return data
	}
	cxnStart += len("<dgm:cxnLst>")
	return data[:cxnStart] + strings.Join(cxns, "") + data[cxnEnd:]
}

func smartArtPointHasText(pointXML string) bool {
	matches := smartArtTextRunPattern.FindAllStringSubmatch(pointXML, -1)
	for _, m := range matches {
		if len(m) < minSmartArtTextSubmatches {
			continue
		}
		if strings.TrimSpace(m[1]) != "" {
			return true
		}
	}
	return false
}
