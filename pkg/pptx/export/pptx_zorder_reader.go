package export

import (
	"archive/zip"
	"encoding/xml"
	"strings"
)

// PowerPoint paints a slide in shape-tree order: the first child of <p:spTree>
// is furthest back, the last is furthest forward. The native PDF renderer used
// to paint in fixed layers instead — every picture, then every shape — which
// puts a picture behind a shape that the deck placed behind it. On a slide where
// a white card sits behind a photo, the card was painted last and hid the photo
// completely.
//
// extractSlideZOrder recovers that ordering as a shape-id -> paint-index map, so
// the renderer can sort a slide's elements back into document order. Shape ids
// are unique within a slide and every paintable child of spTree carries one.

// extractSlideZOrder returns, per slide (0-based), a map of shape id to its
// position in the slide's shape tree.
func extractSlideZOrder(pptxPath string) ([]map[int]int, error) {
	zr, err := zip.OpenReader(pptxPath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	fileMap := make(map[string]*zip.File, len(zr.File))
	for _, f := range zr.File {
		fileMap[canonicalZipPath(f.Name)] = f
	}
	slideOrder := resolveSlideOrder(fileMap)

	out := make([]map[int]int, len(slideOrder))
	for i, slidePart := range slideOrder {
		data := readZipBytes(fileMap, slidePart)
		if data == nil {
			continue
		}
		out[i] = parseSlideZOrder(data)
	}
	return out, nil
}

// shapeTreeChildren are the spTree children that paint something and therefore
// take part in z-order.
//
//nolint:gochecknoglobals // Immutable lookup set.
var shapeTreeChildren = map[string]bool{
	"sp":                true,
	shapeTreePicElement: true,
	"graphicFrame":      true,
	"grpSp":             true,
	"cxnSp":             true,
	"contentPart":       true,
	"AlternateContent":  true,
}

// parseSlideZOrder walks the shape tree and records each child's paint index
// against its shape id.
func parseSlideZOrder(data []byte) map[int]int {
	order := make(map[int]int)
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	depth := 0
	inTree := false
	treeDepth := 0
	index := 0

	for {
		token, err := dec.Token()
		if err != nil {
			return order
		}
		switch element := token.(type) {
		case xml.StartElement:
			depth++
			if element.Name.Local == "spTree" {
				inTree, treeDepth = true, depth
				continue
			}
			// Only direct children of spTree are z-ordered against each other;
			// anything deeper belongs to one of them.
			if inTree && depth == treeDepth+1 && shapeTreeChildren[element.Name.Local] {
				if id, ok := firstShapeID(dec, &element); ok {
					order[id] = index
				}
				index++
				depth--
			}
		case xml.EndElement:
			if inTree && depth == treeDepth && element.Name.Local == "spTree" {
				inTree = false
			}
			depth--
		}
	}
}

// firstShapeID consumes one shape-tree child and reports its own non-visual id.
//
// The containing element differs per shape kind (nvSpPr, nvPicPr, and so on) and
// Go's encoding/xml has no wildcard for a path segment, so each is named. Only
// the child's own id is wanted, not those of anything nested inside it, which is
// why the paths stop one level down rather than matching cNvPr anywhere.
func firstShapeID(dec *xml.Decoder, start *xml.StartElement) (int, bool) {
	type idAttr struct {
		ID *int `xml:"id,attr"`
	}
	var node struct {
		Sp           []idAttr `xml:"nvSpPr>cNvPr"`
		Pic          []idAttr `xml:"nvPicPr>cNvPr"`
		GraphicFrame []idAttr `xml:"nvGraphicFramePr>cNvPr"`
		GroupShape   []idAttr `xml:"nvGrpSpPr>cNvPr"`
		Connector    []idAttr `xml:"nvCxnSpPr>cNvPr"`
	}
	if err := dec.DecodeElement(&node, start); err != nil {
		return 0, false
	}
	for _, group := range [][]idAttr{node.Sp, node.Pic, node.GraphicFrame, node.GroupShape, node.Connector} {
		for _, cNvPr := range group {
			if cNvPr.ID != nil {
				return *cNvPr.ID, true
			}
		}
	}
	return 0, false
}
