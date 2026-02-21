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

// CustomXMLPart represents an embedded custom XML document in the PPTX package.
// The Content field must be a valid XML string.
type CustomXMLPart struct {
	Content string `json:"content" xml:",innerxml"`
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
	Creator        string `json:"creator"        xml:"http://purl.org/dc/elements/1.1/ creator,omitempty"`
	Keywords       string `json:"keywords"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties keywords,omitempty"`
	Description    string `json:"description"    xml:"http://purl.org/dc/elements/1.1/ description,omitempty"`
	LastModifiedBy string `json:"lastModifiedBy" xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties lastModifiedBy,omitempty"`
	Revision       string `json:"revision"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties revision,omitempty"` //nolint:lll // struct tags cannot be split
	Created        string `json:"created"        xml:"http://purl.org/dc/terms/ created,omitempty"`
	Modified       string `json:"modified"       xml:"http://purl.org/dc/terms/ modified,omitempty"`
	Category       string `json:"category"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties category,omitempty"`      //nolint:lll // struct tags cannot be split
	ContentStatus  string `json:"contentStatus"  xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties contentStatus,omitempty"` //nolint:lll // struct tags cannot be split
}
