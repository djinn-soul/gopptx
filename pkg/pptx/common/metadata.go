package common

// Forward declarations or interface if needed?
// Actually, SlideMaster is in elements. elements imports common.
// styling doesn't import common.
// common should probably not import styling or elements if they import common.

// SlideSize describes the dimensions of slides in a presentation in EMUs.
type SlideSize struct {
	Width  int64
	Height int64
}

var (
	// SlideSize4x3 is the standard 4:3 slide size (10x7.5 inches).
	SlideSize4x3 = SlideSize{Width: 9144000, Height: 6858000}
	// SlideSize16x9 is the standard 16:9 widescreen slide size (13.33x7.5 inches).
	SlideSize16x9 = SlideSize{Width: 12192000, Height: 6858000}
)

// CustomXMLPart represents an embedded custom XML document in the PPTX package.
// The Content field must be a valid XML string.
type CustomXMLPart struct {
	Content string `json:"content" xml:",innerxml"`
}

// PresentationMetadata describes summary information for a PPTX package.
type PresentationMetadata struct {
	Title          string
	Subject        string
	Creator        string
	Description    string
	FooterText     string
	ShowDateTime   bool
	SlideSize      SlideSize
	SlideCount     int
	CustomXML      []CustomXMLPart
	CoreProperties CoreProperties
}

// CoreProperties represents the docProps/core.xml metadata.
type CoreProperties struct {
	Title          string `json:"title"          xml:"http://purl.org/dc/elements/1.1/ title,omitempty"`
	Subject        string `json:"subject"        xml:"http://purl.org/dc/elements/1.1/ subject,omitempty"`
	Creator        string `json:"creator"        xml:"http://purl.org/dc/elements/1.1/ creator,omitempty"`
	Keywords       string `json:"keywords"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties keywords,omitempty"`
	Description    string `json:"description"    xml:"http://purl.org/dc/elements/1.1/ description,omitempty"`
	LastModifiedBy string `json:"lastModifiedBy" xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties lastModifiedBy,omitempty"`
	Revision       string `json:"revision"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties revision,omitempty"`
	Created        string `json:"created"        xml:"http://purl.org/dc/terms/ created,omitempty"`
	Modified       string `json:"modified"       xml:"http://purl.org/dc/terms/ modified,omitempty"`
	Category       string `json:"category"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties category,omitempty"`
	ContentStatus  string `json:"contentStatus"  xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties contentStatus,omitempty"`
}
