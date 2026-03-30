package pptxxml

type smartArtTraversalNode struct {
	text     string
	children []SmartArtNodeSpec
}

func smartArtOrderedTextsForLayout(layoutURI string, nodes []SmartArtNodeSpec) []string {
	if prefersBreadthFirstSmartArtTextOrder(layoutURI) {
		return flattenSmartArtNodeTextsBreadthFirst(nodes)
	}
	return flattenSmartArtNodeTextsDepthFirst(nodes)
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

func flattenSmartArtNodeTextsDepthFirst(nodes []SmartArtNodeSpec) []string {
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

func flattenSmartArtNodeTextsBreadthFirst(nodes []SmartArtNodeSpec) []string {
	if len(nodes) == 0 {
		return nil
	}
	out := make([]string, 0, flattenSmartArtTextsInitCap)
	queue := make([]smartArtTraversalNode, 0, len(nodes))
	for _, node := range nodes {
		queue = append(queue, smartArtTraversalNode{text: node.Text, children: node.Children})
	}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		out = append(out, current.text)
		for _, child := range current.children {
			queue = append(queue, smartArtTraversalNode{text: child.Text, children: child.Children})
		}
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

	childrenByParent := make(map[string][]string)
	rootDocID := ""
	for _, point := range points {
		if point.isPres || point.pointType != "doc" {
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
		childrenByParent[cxn.srcID] = append(childrenByParent[cxn.srcID], cxn.destID)
	}

	queue := append([]string(nil), childrenByParent[rootDocID]...)
	if len(queue) == 0 {
		return nil
	}
	out := make([]string, 0, len(queue))
	seen := make(map[string]struct{}, len(queue))
	for len(queue) > 0 {
		modelID := queue[0]
		queue = queue[1:]
		if _, exists := seen[modelID]; exists {
			continue
		}
		seen[modelID] = struct{}{}
		out = append(out, modelID)
		queue = append(queue, childrenByParent[modelID]...)
	}
	return out
}
