package smartart

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"slices"
	"strconv"
	"strings"
)

const (
	defaultDataModelCapacity = 8
	defaultTopLevelCapacity  = 4
	minChildrenToSort        = 2
)

type dataModel struct {
	points      map[string]dataPoint
	pointOrder  []string
	connections []dataConnection
}

type dataPoint struct {
	modelID   string
	pointType string
	text      string
}

type dataConnection struct {
	cxnType string
	srcID   string
	destID  string
	srcOrd  int
}

type pointXML struct {
	ModelID string `xml:"modelId,attr"`
	Type    string `xml:"type,attr"`
	Inner   string `xml:",innerxml"`
}

type connectionXML struct {
	CxnType string `xml:"type,attr"`
	SrcID   string `xml:"srcId,attr"`
	DestID  string `xml:"destId,attr"`
	SrcOrd  string `xml:"srcOrd,attr"`
}

type orderedChild struct {
	ID    string
	Order int
}

// ParseDataModelNodes rebuilds semantic SmartArt nodes from a SmartArt dataX.xml part.
func ParseDataModelNodes(dataXML []byte) ([]Node, error) {
	model, err := parseDataModel(dataXML)
	if err != nil {
		return nil, err
	}
	return model.semanticNodes(), nil
}

func parseDataModel(dataXML []byte) (dataModel, error) {
	model := dataModel{
		points:      make(map[string]dataPoint),
		pointOrder:  make([]string, 0, defaultDataModelCapacity),
		connections: make([]dataConnection, 0, defaultDataModelCapacity),
	}

	decoder := xml.NewDecoder(bytes.NewReader(dataXML))
	for {
		tok, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return dataModel{}, err
		}
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		switch start.Name.Local {
		case "pt":
			var parsed pointXML
			if err := decoder.DecodeElement(&parsed, &start); err != nil {
				return dataModel{}, err
			}
			point := dataPoint{
				modelID:   parsed.ModelID,
				pointType: strings.ToLower(strings.TrimSpace(parsed.Type)),
				text:      extractPointText(parsed.Inner),
			}
			if point.modelID == "" {
				continue
			}
			model.points[point.modelID] = point
			model.pointOrder = append(model.pointOrder, point.modelID)
		case "cxn":
			var parsed connectionXML
			if err := decoder.DecodeElement(&parsed, &start); err != nil {
				return dataModel{}, err
			}
			model.connections = append(model.connections, dataConnection{
				cxnType: strings.ToLower(strings.TrimSpace(parsed.CxnType)),
				srcID:   parsed.SrcID,
				destID:  parsed.DestID,
				srcOrd:  parseOrdinal(parsed.SrcOrd),
			})
		}
	}

	return model, nil
}

func extractPointText(pointXML string) string {
	matches := reSmartArtNodeText.FindAllStringSubmatch(pointXML, -1)
	parts := make([]string, 0, len(matches))
	for _, match := range matches {
		text := strings.TrimSpace(match[1])
		if text == "" {
			continue
		}
		parts = append(parts, text)
	}
	return strings.Join(parts, "\n")
}

func parseOrdinal(value string) int {
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0
	}
	return n
}

func (m dataModel) semanticNodes() []Node {
	childrenByParent := make(map[string][]orderedChild)
	parentByChild := make(map[string]string)
	topLevel := make([]orderedChild, 0, defaultTopLevelCapacity)

	for _, cxn := range m.connections {
		if strings.HasPrefix(cxn.cxnType, "pres") {
			continue
		}
		dest := m.points[cxn.destID]
		if !isSemanticPoint(dest) {
			continue
		}
		src := m.points[cxn.srcID]
		switch {
		case src.pointType == "doc":
			topLevel = append(topLevel, orderedChild{ID: dest.modelID, Order: cxn.srcOrd})
		case isSemanticPoint(src):
			childrenByParent[src.modelID] = append(childrenByParent[src.modelID], orderedChild{
				ID:    dest.modelID,
				Order: cxn.srcOrd,
			})
			parentByChild[dest.modelID] = src.modelID
		}
	}

	nodes := make([]Node, 0, len(topLevel))
	seenRoots := make(map[string]struct{})
	for _, root := range sortChildren(topLevel, m.pointOrder) {
		if _, exists := seenRoots[root.ID]; exists {
			continue
		}
		nodes = append(nodes, m.buildNode(root.ID, childrenByParent, map[string]struct{}{}))
		seenRoots[root.ID] = struct{}{}
	}

	for _, modelID := range m.pointOrder {
		point := m.points[modelID]
		if !isSemanticPoint(point) {
			continue
		}
		if _, hasParent := parentByChild[modelID]; hasParent {
			continue
		}
		if _, exists := seenRoots[modelID]; exists {
			continue
		}
		nodes = append(nodes, m.buildNode(modelID, childrenByParent, map[string]struct{}{}))
	}

	return nodes
}

func (m dataModel) buildNode(
	modelID string,
	childrenByParent map[string][]orderedChild,
	path map[string]struct{},
) Node {
	point := m.points[modelID]
	node := NewNode(point.text)
	if _, seen := path[modelID]; seen {
		return node
	}

	nextPath := clonePath(path)
	nextPath[modelID] = struct{}{}
	for _, child := range sortChildren(childrenByParent[modelID], m.pointOrder) {
		node = node.WithChild(m.buildNode(child.ID, childrenByParent, nextPath))
	}
	return node
}

func isSemanticPoint(point dataPoint) bool {
	if point.modelID == "" || strings.TrimSpace(point.text) == "" {
		return false
	}
	switch point.pointType {
	case "doc", "pres", "partrans", "sibtrans":
		return false
	default:
		return true
	}
}

func sortChildren(children []orderedChild, pointOrder []string) []orderedChild {
	if len(children) < minChildrenToSort {
		return children
	}
	orderIndex := make(map[string]int, len(pointOrder))
	for idx, id := range pointOrder {
		orderIndex[id] = idx
	}
	out := append([]orderedChild(nil), children...)
	slices.SortFunc(out, func(a, b orderedChild) int {
		return compareChildren(a, b, orderIndex)
	})
	return out
}

func compareChildren(a, b orderedChild, orderIndex map[string]int) int {
	if a.Order != b.Order {
		return a.Order - b.Order
	}
	return orderIndex[a.ID] - orderIndex[b.ID]
}

func clonePath(path map[string]struct{}) map[string]struct{} {
	if len(path) == 0 {
		return make(map[string]struct{}, 1)
	}
	cloned := make(map[string]struct{}, len(path)+1)
	for key := range path {
		cloned[key] = struct{}{}
	}
	return cloned
}
