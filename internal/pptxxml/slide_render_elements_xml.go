package pptxxml

import "strings"

func slideRenderImages(b *strings.Builder, images []ImageRef, nextID int) int {
	for i, image := range images {
		b.WriteString(imageShape(image, nextID+i))
	}
	return nextID + len(images)
}

func slideRenderShapes(b *strings.Builder, shapes []ShapeSpec, nextID int) ([]int, int) {
	shapeIDs, shapeXMLParts := renderCustomShapeXMLConcurrently(shapes, nextID)
	for _, part := range shapeXMLParts {
		b.WriteString(part)
	}
	return shapeIDs, nextID + len(shapes)
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
