package pptxxml

import "sort"

const (
	smartArtDataPointDoc       = "doc"
	smartArtDataPointAssistant = "asst"
)

// smartArtOrderedNodesForLayout flattens the spec into the order the layout's
// slots are filled. Per-node properties travel with the node so that a colour or
// a picture lands on the same point its text does.
func smartArtOrderedNodesForLayout(layoutURI string, nodes []SmartArtNodeSpec) []SmartArtNodeSpec {
	if prefersBreadthFirstSmartArtTextOrder(layoutURI) {
		return flattenSmartArtNodesBreadthFirst(nodes)
	}
	return flattenSmartArtNodesDepthFirst(nodes)
}

func smartArtOrderedTextsForLayout(layoutURI string, nodes []SmartArtNodeSpec) []string {
	ordered := smartArtOrderedNodesForLayout(layoutURI, nodes)
	out := make([]string, 0, len(ordered))
	for _, node := range ordered {
		out = append(out, node.Text)
	}
	return out
}

func prefersBreadthFirstSmartArtTextOrder(layoutURI string) bool {
	switch layoutURI {
	case "urn:microsoft.com/office/officeart/2005/8/layout/hierarchy1",
		"urn:microsoft.com/office/officeart/2005/8/layout/hierarchy2":
		return true
	default:
		return false
	}
}

func flattenSmartArtNodesDepthFirst(nodes []SmartArtNodeSpec) []SmartArtNodeSpec {
	out := make([]SmartArtNodeSpec, 0, flattenSmartArtTextsInitCap)
	var walk func([]SmartArtNodeSpec)
	walk = func(items []SmartArtNodeSpec) {
		for _, n := range items {
			out = append(out, n)
			if len(n.Children) > 0 {
				walk(n.Children)
			}
		}
	}
	walk(nodes)
	return out
}

func flattenSmartArtNodesBreadthFirst(nodes []SmartArtNodeSpec) []SmartArtNodeSpec {
	if len(nodes) == 0 {
		return nil
	}
	out := make([]SmartArtNodeSpec, 0, flattenSmartArtTextsInitCap)
	queue := append([]SmartArtNodeSpec(nil), nodes...)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		out = append(out, current)
		queue = append(queue, current.Children...)
	}
	return out
}

func preferredDataModelIDsForLayout(layoutURI, data string) []string {
	if prefersBreadthFirstSmartArtTextOrder(layoutURI) {
		if ids := semanticDataModelIDsInBreadthFirstOrder(data); len(ids) > 0 {
			return ids
		}
	}
	return preferredDataModelIDsInOrder(data)
}

//nolint:gocognit
func semanticDataModelIDsInBreadthFirstOrder(data string) []string {
	points := parseSmartArtDataPoints(data)
	cxns := parseSmartArtDataConnections(data)
	if len(points) == 0 || len(cxns) == 0 {
		return nil
	}

	pointByID := make(map[string]smartArtDataPoint, len(points))
	for _, point := range points {
		if point.modelID == "" {
			continue
		}
		pointByID[point.modelID] = point
	}

	childrenByParent := make(map[string][]smartArtChildLink)
	rootDocID := ""
	for _, point := range points {
		if point.isPres || point.pointType != smartArtDataPointDoc {
			continue
		}
		rootDocID = point.modelID
		break
	}
	if rootDocID == "" {
		return nil
	}

	for _, cxn := range cxns {
		if cxn.srcID == "" || cxn.destID == "" || cxn.cxnType == "presof" || cxn.cxnType == "presparof" {
			continue
		}
		destPoint, ok := pointByID[cxn.destID]
		if !ok || destPoint.isPres || isSmartArtStructuralDataType(destPoint.pointType) {
			continue
		}
		childrenByParent[cxn.srcID] = append(childrenByParent[cxn.srcID], smartArtChildLink{
			modelID: cxn.destID,
			ord:     cxn.srcOrd,
		})
	}
	sortSmartArtChildLinks(childrenByParent)

	queue := append([]smartArtChildLink(nil), childrenByParent[rootDocID]...)
	if len(queue) == 0 {
		return nil
	}
	out := make([]string, 0, len(queue))
	seen := make(map[string]struct{}, len(queue))
	for len(queue) > 0 {
		link := queue[0]
		queue = queue[1:]
		if _, exists := seen[link.modelID]; exists {
			continue
		}
		seen[link.modelID] = struct{}{}
		out = append(out, link.modelID)
		queue = append(queue, childrenByParent[link.modelID]...)
	}
	return out
}

// smartArtChildLink is one parent->child edge of the data tree, carrying the
// ordinal that fixes sibling order.
type smartArtChildLink struct {
	modelID string
	ord     int
}

func sortSmartArtChildLinks(childrenByParent map[string][]smartArtChildLink) {
	for _, children := range childrenByParent {
		sort.SliceStable(children, func(i, j int) bool {
			return children[i].ord < children[j].ord
		})
	}
}
