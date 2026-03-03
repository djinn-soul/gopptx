package common

import "time"

// Forward declarations or interface if needed?
// Actually, SlideMaster is in elements. elements imports common.
// styling doesn't import common.
// common should probably not import styling or elements if they import common.

// SlideSize describes the dimensions of slides in a presentation in EMUs.
type SlideSize struct {
	Width  int64
	Height int64
}

const (
	width4x3  = 9144000
	height4x3 = 6858000
	width16x9 = 12192000
)

// GetSlideSize4x3 returns the standard 4:3 slide size (10x7.5 inches).
func GetSlideSize4x3() SlideSize {
	return SlideSize{Width: width4x3, Height: height4x3}
}

// GetSlideSize16x9 returns the standard 16:9 widescreen slide size (13.33x7.5 inches).
func GetSlideSize16x9() SlideSize {
	return SlideSize{Width: width16x9, Height: height4x3}
}

// CustomXMLKV is a key-value property for a CustomXMLPart.
type CustomXMLKV struct {
	Key   string
	Value string
}

// CustomXMLPart represents an embedded custom XML document in the PPTX package.
// If RootElement is populated, the XML is generated structurally.
// Otherwise, Content must be a valid XML string for legacy passthrough.
type CustomXMLPart struct {
	ItemID      string        `json:"itemID,omitempty"`
	Content     string        `json:"content,omitempty"     xml:",innerxml"`
	RootElement string        `json:"rootElement,omitempty"`
	Namespace   string        `json:"namespace,omitempty"`
	Properties  []CustomXMLKV `json:"properties,omitempty"`
}

// Metadata describes summary information for a PPTX package.
type Metadata struct {
	Title          string
	Subject        string
	Creator        string
	Description    string
	FooterText     string
	ShowDateTime   bool
	GeneratedDate  time.Time
	SlideSize      SlideSize
	SlideCount     int
	CustomXML      []CustomXMLPart
	CoreProperties CoreProperties
	Protection     Protection
}

// Protection defines write-protection and suggested read-only settings.
type Protection struct {
	ModifyPassword    string
	MarkAsFinal       bool
	SignaturesEnabled bool
}

// CoreProperties represents the docProps/core.xml metadata.
type CoreProperties struct {
	Title          string `json:"title"          xml:"http://purl.org/dc/elements/1.1/ title,omitempty"`
	Subject        string `json:"subject"        xml:"http://purl.org/dc/elements/1.1/ subject,omitempty"`
	Creator        string `json:"creator"        xml:"http://purl.org/dc/elements/1.1/ creator,omitempty"` // Also known as "author" in python-pptx
	Keywords       string `json:"keywords"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties keywords,omitempty"`
	Description    string `json:"description"    xml:"http://purl.org/dc/elements/1.1/ description,omitempty"` // Also known as "comments" in python-pptx
	LastModifiedBy string `json:"lastModifiedBy" xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties lastModifiedBy,omitempty"`
	Revision       string `json:"revision"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties revision,omitempty"` //nolint:lll // struct tags cannot be split
	Created        string `json:"created"        xml:"http://purl.org/dc/terms/ created,omitempty"`
	Modified       string `json:"modified"       xml:"http://purl.org/dc/terms/ modified,omitempty"`
	Category       string `json:"category"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties category,omitempty"`      //nolint:lll // struct tags cannot be split
	ContentStatus  string `json:"contentStatus"  xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties contentStatus,omitempty"` //nolint:lll // struct tags cannot be split
	Identifier     string `json:"identifier"     xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties identifier,omitempty"`    //nolint:lll // struct tags cannot be split
	Language       string `json:"language"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties language,omitempty"`      //nolint:lll // struct tags cannot be split
	LastPrinted    string `json:"lastPrinted"    xml:"http://purl.org/dc/terms/ lastPrinted,omitempty"`
	Version        string `json:"version"        xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties version,omitempty"` //nolint:lll // struct tags cannot be split
}
