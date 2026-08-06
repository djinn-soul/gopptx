package pptxxml

import (
	"fmt"
	"strings"
)

// SmartArt templates carry a fixed number of slots. When a diagram has more
// nodes than the template can hold, filling slots silently reshapes the data:
// a five-step process poured into a three-step template turned steps four and
// five into sub-steps of steps one and two.
//
// PowerPoint lays a diagram out from its data model, not from the cached
// presentation nodes, so a generated data model with the right shape renders
// correctly and round-trips through the reader unchanged.

const (
	smartArtGeneratedIDFormat = "{%08X-1D1D-4E4E-9A9A-%012X}"
	smartArtDocIDSeed         = 1
	// minSmartArtPeerRoots is the number of roots a template needs before its
	// root slots count as a capacity limit rather than a single absorbing root.
	minSmartArtPeerRoots = 2
)

// smartArtSpecFitsTemplate reports whether every level of the spec has a slot in
// the template's data tree.
func smartArtSpecFitsTemplate(spec SmartArtSpec, data string) bool {
	capacityByParent, rootID := smartArtTemplateCapacity(data)
	if rootID == "" {
		return true
	}
	// A template with a single root — a radial hub, a matrix centre, a hierarchy
	// root — is meant to absorb a flat list underneath that root, so its one slot
	// is not a capacity limit. Only templates whose roots are peers (the steps of
	// a process, the entries of a list) run out of room.
	if len(capacityByParent[rootID]) < minSmartArtPeerRoots {
		return true
	}
	var fits func(nodes []SmartArtNodeSpec, slots []string) bool
	fits = func(nodes []SmartArtNodeSpec, slots []string) bool {
		if len(nodes) > len(slots) {
			return false
		}
		for i, node := range nodes {
			if len(node.Children) == 0 {
				continue
			}
			if !fits(node.Children, capacityByParent[slots[i]]) {
				return false
			}
		}
		return true
	}
	return fits(spec.Nodes, capacityByParent[rootID])
}

func smartArtSpecHasChildren(nodes []SmartArtNodeSpec) bool {
	for _, node := range nodes {
		if len(node.Children) > 0 {
			return true
		}
	}
	return false
}

// smartArtTemplateCapacity maps every data point of the template to its child
// slots, in ordinal order, and returns the document point that roots them.
func smartArtTemplateCapacity(data string) (map[string][]string, string) {
	points := parseSmartArtDataPoints(data)
	cxns := parseSmartArtDataConnections(data)

	semantic := make(map[string]bool, len(points))
	rootID := ""
	for _, point := range points {
		if point.isPres || point.modelID == "" {
			continue
		}
		if point.pointType == smartArtDataPointDoc {
			rootID = point.modelID
			continue
		}
		// An assistant hangs off the side of an org chart rather than taking a
		// child's place, so it is not a slot ordinary children can be poured into.
		semantic[point.modelID] = !isSmartArtStructuralDataType(point.pointType) &&
			point.pointType != smartArtDataPointAssistant
	}

	childrenByParent := make(map[string][]smartArtChildLink)
	for _, cxn := range cxns {
		if strings.HasPrefix(cxn.cxnType, "pres") || cxn.srcID == "" || cxn.destID == "" {
			continue
		}
		if !semantic[cxn.destID] {
			continue
		}
		childrenByParent[cxn.srcID] = append(childrenByParent[cxn.srcID], smartArtChildLink{
			modelID: cxn.destID,
			ord:     cxn.srcOrd,
		})
	}
	sortSmartArtChildLinks(childrenByParent)

	out := make(map[string][]string, len(childrenByParent))
	for parent, children := range childrenByParent {
		slots := make([]string, 0, len(children))
		for _, child := range children {
			slots = append(slots, child.modelID)
		}
		out[parent] = slots
	}
	return out, rootID
}

// renderGeneratedSmartArtData writes a data model that mirrors the spec exactly.
func renderGeneratedSmartArtData(spec SmartArtSpec) string {
	var points, cxns strings.Builder
	seq := smartArtDocIDSeed
	docID := smartArtGeneratedID(&seq)

	quickStyleID := defaultQuickStyleID(spec.QuickStyleID)
	colorStyleID := defaultColorStyleID(spec.ColorStyleID)
	points.WriteString(`<dgm:pt modelId="` + docID + `" type="doc"><dgm:prSet loTypeId="` +
		Escape(layoutURIOrDefault(spec.LayoutURI)) + `" loCatId="" qsTypeId="` +
		Escape(quickStyleID) + `" qsCatId="` + Escape(smartArtQuickStyleCategory(quickStyleID)) +
		`" csTypeId="` + Escape(colorStyleID) +
		`" csCatId="` + Escape(smartArtColorCategory(colorStyleID)) + `" phldr="0"/>` +
		`<dgm:spPr/>` + smartArtEmptyTextBody() + `</dgm:pt>`)

	var walk func(parentID string, nodes []SmartArtNodeSpec)
	walk = func(parentID string, nodes []SmartArtNodeSpec) {
		for i, node := range nodes {
			nodeID := smartArtGeneratedID(&seq)
			parTransID := smartArtGeneratedID(&seq)
			sibTransID := smartArtGeneratedID(&seq)
			cxnID := smartArtGeneratedID(&seq)

			points.WriteString(`<dgm:pt modelId="` + nodeID + `"><dgm:prSet/>` +
				smartArtNodeShapeProperties(node) + smartArtTextBody(node.Text) + `</dgm:pt>`)
			points.WriteString(smartArtTransitionPoint(parTransID, "parTrans", cxnID))
			points.WriteString(smartArtTransitionPoint(sibTransID, "sibTrans", cxnID))

			fmt.Fprintf(
				&cxns,
				`<dgm:cxn modelId="%s" srcId="%s" destId="%s" srcOrd="%d"`+
					` destOrd="0" parTransId="%s" sibTransId="%s"/>`,
				cxnID, parentID, nodeID, i, parTransID, sibTransID,
			)

			walk(nodeID, node.Children)
		}
	}
	walk(docID, spec.Nodes)

	return xmlHeader +
		`<dgm:dataModel xmlns:dgm="http://schemas.openxmlformats.org/drawingml/2006/diagram"` +
		` xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
		`<dgm:ptLst>` + points.String() + `</dgm:ptLst>` +
		`<dgm:cxnLst>` + cxns.String() + `</dgm:cxnLst>` +
		`<dgm:bg/><dgm:whole/></dgm:dataModel>`
}

func smartArtTransitionPoint(modelID, pointType, cxnID string) string {
	return `<dgm:pt modelId="` + modelID + `" type="` + pointType + `" cxnId="` + cxnID + `">` +
		`<dgm:prSet/><dgm:spPr/>` + smartArtEmptyTextBody() + `</dgm:pt>`
}

func smartArtTextBody(text string) string {
	if text == "" {
		return smartArtEmptyTextBody()
	}
	return `<dgm:t><a:bodyPr/><a:lstStyle/><a:p>` + injectedSmartArtRunXML(Escape(text)) +
		`<a:endParaRPr lang="en-US"/></a:p></dgm:t>`
}

func smartArtEmptyTextBody() string {
	return `<dgm:t><a:bodyPr/><a:lstStyle/><a:p><a:endParaRPr lang="en-US"/></a:p></dgm:t>`
}

// smartArtGeneratedID hands out deterministic model IDs so the same diagram
// always produces byte-identical XML.
func smartArtGeneratedID(seq *int) string {
	id := fmt.Sprintf(smartArtGeneratedIDFormat, *seq, *seq)
	*seq++
	return id
}
