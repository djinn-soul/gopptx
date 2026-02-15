package gopptx

import "encoding/xml"

const (
	NS_P = "http://schemas.openxmlformats.org/presentationml/2006/main"
)

// Presentation represents the main presentation XML component.
type Presentation struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/presentationml/2006/main presentation"`
}

// Slide represents a single slide XML component.
type Slide struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/presentationml/2006/main sld"`
}
