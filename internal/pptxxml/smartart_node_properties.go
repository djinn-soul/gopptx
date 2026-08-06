package pptxxml

import (
	"regexp"
	"strings"
)

// A node's own fill and picture live on its data point, next to its text, as
// shape properties. Neither used to be written: a colour set on a node was
// dropped on the way to the XML, and the picture layouts had no way at all to
// fill the image placeholders they draw.

// applySmartArtNodeProperties writes each node's fill onto its own data point and
// each node's picture onto the presentation point that draws its image
// placeholder — which is where PowerPoint puts one. A picture written to the
// data point instead bleeds onto every shape that point draws, filling a
// caption band rather than the picture frame.
func applySmartArtNodeProperties(
	data string,
	propsByModelID map[string]SmartArtNodeSpec,
	pictureShapeNames map[string]struct{},
) string {
	if len(propsByModelID) == 0 {
		return data
	}
	imageByPresID := smartArtImagePresPoints(data, propsByModelID, pictureShapeNames)

	segments := strings.Split(data, "<dgm:pt ")
	if len(segments) <= 1 {
		return data
	}
	var b strings.Builder
	b.WriteString(segments[0])
	for i := 1; i < len(segments); i++ {
		segment := "<dgm:pt " + segments[i]
		modelID := extractXMLAttr(segment, "modelId")
		switch {
		case imageByPresID[modelID] != "":
			segment = strings.Replace(
				segment,
				"<dgm:spPr/>",
				"<dgm:spPr>"+smartArtPictureFillXML(imageByPresID[modelID])+"</dgm:spPr>",
				1,
			)
		default:
			if node, ok := propsByModelID[modelID]; ok && node.Color != "" {
				segment = strings.Replace(segment, "<dgm:spPr/>", smartArtNodeShapeProperties(node), 1)
			}
		}
		b.WriteString(segment)
	}
	return b.String()
}

// smartArtImagePresPoints maps the presentation point that draws each node's
// image placeholder to the relationship of the picture that fills it.
func smartArtImagePresPoints(
	data string,
	propsByModelID map[string]SmartArtNodeSpec,
	pictureShapeNames map[string]struct{},
) map[string]string {
	out := map[string]string{}
	for _, node := range parseSmartArtPresNodesInOrder(data) {
		if _, isPicture := pictureShapeNames[node.presName]; !isPicture {
			continue
		}
		source, ok := propsByModelID[node.presAssocID]
		if !ok || source.ImageRelID == "" {
			continue
		}
		out[node.modelID] = source.ImageRelID
	}
	return out
}

// smartArtNodeShapeProperties renders the point's spPr: a solid fill for a
// colour, a picture fill for an image, and the empty element when the node asks
// for neither.
func smartArtNodeShapeProperties(node SmartArtNodeSpec) string {
	fill := smartArtNodeFillXML(node)
	if fill == "" {
		return "<dgm:spPr/>"
	}
	return "<dgm:spPr>" + fill + "</dgm:spPr>"
}

func smartArtNodeFillXML(node SmartArtNodeSpec) string {
	if node.Color != "" {
		return `<a:solidFill><a:srgbClr val="` + Escape(strings.ToUpper(node.Color)) + `"/></a:solidFill>`
	}
	return ""
}

// smartArtPictureFillXML matches what PowerPoint writes for a filled SmartArt
// picture placeholder.
func smartArtPictureFillXML(relID string) string {
	return `<a:blipFill><a:blip xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"` +
		` r:embed="` + Escape(relID) + `"/><a:stretch><a:fillRect/></a:stretch></a:blipFill>`
}

// smartArtPictureShapeNames reads the layout definition for the shapes that are
// picture placeholders. The layout marks them with blipPhldr="1"; their names
// vary per layout ("image" in one, "rect1" in another), so they cannot be
// guessed from the name alone.
func smartArtPictureShapeNames(layoutURI string) map[string]struct{} {
	out := map[string]struct{}{}
	layout := mustTemplate(templatePathForLayout(layoutURI, "layout.xml"))
	for _, match := range layoutNodeShapePattern.FindAllStringSubmatch(layout, -1) {
		if strings.Contains(match[2], `blipPhldr="1"`) {
			out[strings.ToLower(match[1])] = struct{}{}
		}
	}
	return out
}

// layoutNodeShapePattern captures a layout node's name and the shape element that
// immediately follows it.
var layoutNodeShapePattern = regexp.MustCompile(
	`<dgm:layoutNode name="([^"]+)"[^>]*>(?:<dgm:alg[^>]*/?>)?\s*(<dgm:shape[^>]*>)`,
)

// smartArtNodePropertiesByModelID pairs the ordered nodes with the data points
// their text was written to, keeping only the nodes that carry properties.
func smartArtNodePropertiesByModelID(
	ordered []SmartArtNodeSpec,
	modelIDs []string,
) map[string]SmartArtNodeSpec {
	out := map[string]SmartArtNodeSpec{}
	for i, modelID := range modelIDs {
		if i >= len(ordered) {
			break
		}
		if smartArtNodeFillXML(ordered[i]) == "" && ordered[i].ImageRelID == "" {
			continue
		}
		out[modelID] = ordered[i]
	}
	return out
}

// smartArtSpecHasNodeProperties reports whether any node in the tree carries a
// fill or a picture.
