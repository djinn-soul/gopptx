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
// `<p:prnPr/>` element and showPrXML the rendered `<p:showPr/>`; either may be
// empty. CT_PresentationProperties orders print settings before show settings,
// and it — not CT_Presentation — is where both live.
func PresentationProps(prnPrXML, showPrXML string) string {
	var b strings.Builder
	b.WriteString(xmlHeader)
	b.WriteString(`
<p:presentationPr xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`)
	b.WriteString(prnPrXML)
	b.WriteString(showPrXML)
	// The three extensions Office writes into presProps by default. Cosmetic for
	// rendering; without them a generated deck's presProps differs from one
	// PowerPoint has round-tripped.
	b.WriteString(`<p:extLst>` +
		`<p:ext uri="{E76CE94A-603C-4142-B9EB-6D1370010A27}">` +
		`<p14:discardImageEditData xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main" val="0"/>` +
		`</p:ext>` +
		`<p:ext uri="{D31A062A-798A-4329-ABDD-BBA856620510}">` +
		`<p14:defaultImageDpi xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main" val="220"/>` +
		`</p:ext>` +
		`<p:ext uri="{FD5EFAAD-0ECE-453E-9831-46B23BE46B34}">` +
		`<p15:chartTrackingRefBased xmlns:p15="http://schemas.microsoft.com/office/powerpoint/2012/main" val="1"/>` +
		`</p:ext>` +
		`</p:extLst>`)
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

// WithContentTypeDefault inserts a Default entry for a file extension into an
// already rendered `[Content_Types].xml`, which is how media parts declare
// their type.
func WithContentTypeDefault(contentTypes, extension, contentType string) string {
	if extension == "" || contentType == "" {
		return contentTypes
	}
	if strings.Contains(contentTypes, `Extension="`+extension+`"`) {
		return contentTypes
	}
	const closing = "</Types>"
	idx := strings.LastIndex(contentTypes, closing)
	if idx < 0 {
		return contentTypes
	}
	entry := "\n<Default Extension=\"" + Escape(extension) +
		"\" ContentType=\"" + Escape(contentType) + "\"/>"
	return contentTypes[:idx] + entry + contentTypes[idx:]
}

// WithExternalRelationship inserts a Relationship whose target is outside the
// package — a hosted video, say — which OPC marks with TargetMode="External".
func WithExternalRelationship(rels, relID, relType, target string) string {
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
		"\" Target=\"" + Escape(target) +
		"\" TargetMode=\"External\"/>"
	return rels[:idx] + entry + rels[idx:]
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
