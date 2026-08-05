package pptxxml

import "strings"

const (
	// PresPropsPartName is the package path of the presentation properties part.
	PresPropsPartName = "ppt/presProps.xml"
	// PresPropsContentType is the content type override for presProps.xml.
	PresPropsContentType = "application/vnd.openxmlformats-officedocument.presentationml.presProps+xml"
	// PresPropsRelationshipType is the relationship type from presentation.xml.
	PresPropsRelationshipType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/presProps"
	// PresPropsRelationshipTarget is the target used from ppt/_rels/presentation.xml.rels.
	PresPropsRelationshipTarget = "presProps.xml"
)

// PresentationProps renders `ppt/presProps.xml`. prnPrXML is the rendered
// `<p:prnPr/>` element, or an empty string to write the part without print
// settings.
func PresentationProps(prnPrXML string) string {
	var b strings.Builder
	b.WriteString(xmlHeader)
	b.WriteString(`
<p:presentationPr xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`)
	b.WriteString(prnPrXML)
	b.WriteString(`</p:presentationPr>`)
	return b.String()
}

// WithContentTypeOverride inserts an Override entry into an already rendered
// `[Content_Types].xml` document. The document is returned unchanged when the
// part already has an override or the closing tag is missing.
func WithContentTypeOverride(contentTypes, partName, contentType string) string {
	if partName == "" || contentType == "" {
		return contentTypes
	}
	if !strings.HasPrefix(partName, "/") {
		partName = "/" + partName
	}
	if strings.Contains(contentTypes, `PartName="`+partName+`"`) {
		return contentTypes
	}
	const closing = "</Types>"
	idx := strings.LastIndex(contentTypes, closing)
	if idx < 0 {
		return contentTypes
	}
	entry := "\n<Override PartName=\"" + Escape(partName) +
		"\" ContentType=\"" + Escape(contentType) + "\"/>"
	return contentTypes[:idx] + entry + contentTypes[idx:]
}

// WithRelationship inserts a Relationship entry into an already rendered
// `.rels` document. The document is returned unchanged when the relationship ID
// is already present or the closing tag is missing.
func WithRelationship(rels, relID, relType, target string) string {
	if relID == "" || relType == "" || target == "" {
		return rels
	}
	if strings.Contains(rels, `Id="`+relID+`"`) {
		return rels
	}
	const closing = "</Relationships>"
	idx := strings.LastIndex(rels, closing)
	if idx < 0 {
		return rels
	}
	entry := "\n<Relationship Id=\"" + Escape(relID) +
		"\" Type=\"" + Escape(relType) +
		"\" Target=\"" + Escape(target) + "\"/>"
	return rels[:idx] + entry + rels[idx:]
}
