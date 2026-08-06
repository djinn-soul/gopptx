package pptxxml

import (
	"sort"
	"strings"
)

// slideRenderPictures writes the pictures and the custom shapes in paint order.
//
// A deck read from a PPTX records each element's position in the shape tree, and
// that order is what PowerPoint paints from — a shape behind a picture has to be
// written before it. Elements the caller built carry no order, so with none set
// the emission stays as it was: every picture, then every shape.
//
// Element ids are allocated as before, independently of emission order: they
// only have to be unique within the slide, and connectors resolve their anchors
// through the returned shape ids.
func slideRenderPictures(
	b *strings.Builder,
	images []ImageRef,
	shapes []ShapeSpec,
	nextID int,
) ([]int, int) {
	imageXMLParts := make([]string, len(images))
	for i, image := range images {
		imageXMLParts[i] = imageShape(image, nextID+i)
	}
	shapeIDs, shapeXMLParts := renderCustomShapeXMLConcurrently(shapes, nextID+len(images))

	for _, part := range orderedTreeXML(images, shapes, imageXMLParts, shapeXMLParts) {
		b.WriteString(part)
	}
	return shapeIDs, nextID + len(images) + len(shapes)
}

// treeEntry is one element awaiting emission, with the tree position it claims.
type treeEntry struct {
	zOrder int
	// pass keeps the pre-existing pictures-then-shapes order for ties and for
	// the common case where nothing states a position.
	pass int
	xml  string
}

func orderedTreeXML(
	images []ImageRef,
	shapes []ShapeSpec,
	imageXMLParts, shapeXMLParts []string,
) []string {
	entries := make([]treeEntry, 0, len(imageXMLParts)+len(shapeXMLParts))
	ordered := false
	for i, part := range imageXMLParts {
		entries = append(entries, treeEntry{zOrder: images[i].ZOrder, pass: 0, xml: part})
		ordered = ordered || images[i].ZOrder != 0
	}
	for i, part := range shapeXMLParts {
		entries = append(entries, treeEntry{zOrder: shapes[i].ZOrder, pass: 1, xml: part})
		ordered = ordered || shapes[i].ZOrder != 0
	}

	if ordered {
		sort.SliceStable(entries, func(i, j int) bool {
			if entries[i].zOrder != entries[j].zOrder {
				return entries[i].zOrder < entries[j].zOrder
			}
			return entries[i].pass < entries[j].pass
		})
	}

	parts := make([]string, 0, len(entries))
	for _, entry := range entries {
		parts = append(parts, entry.xml)
	}
	return parts
}

func slideRenderConnectors(b *strings.Builder, connectors []ConnectorSpec, shapeIDs []int, nextID int) int {
	currentID := nextID
	for _, connector := range connectors {
		startShapeID := shapeAnchorID(shapeIDs, connector.StartShapeIndex)
		endShapeID := shapeAnchorID(shapeIDs, connector.EndShapeIndex)
		b.WriteString(connectorXML(connector, currentID, startShapeID, endShapeID))
		currentID++
		if labelXML := connectorLabelShape(connector, currentID); labelXML != "" {
			b.WriteString(labelXML)
			currentID++
		}
	}
	return currentID
}

func slideRenderPlaceholders(b *strings.Builder, placeholders []PlaceholderOverrideSpec, nextID int) int {
	for i, ph := range placeholders {
		b.WriteString(PlaceholderShape(ph, nextID+i))
	}
	return nextID + len(placeholders)
}
