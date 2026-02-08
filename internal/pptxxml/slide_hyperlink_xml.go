package pptxxml

import "fmt"

// HyperlinkSpec contains data needed for hyperlink XML generation.
type HyperlinkSpec struct {
	RelID          string
	Tooltip        string
	HighlightClick bool
	Action         string // ppaction:// for internal navigation
}

// HyperlinkXML generates the <a:hlinkClick> element for shapes or text runs.
func HyperlinkXML(spec HyperlinkSpec) string {
	xml := fmt.Sprintf(`<a:hlinkClick r:id="%s"`, Escape(spec.RelID))

	if spec.Tooltip != "" {
		xml += fmt.Sprintf(` tooltip="%s"`, Escape(spec.Tooltip))
	}

	if spec.HighlightClick {
		xml += ` highlightClick="1"`
	}

	if spec.Action != "" {
		xml += fmt.Sprintf(` action="%s"`, Escape(spec.Action))
	}

	xml += "/>"
	return xml
}

// HyperlinkRelationshipXML generates a relationship element for hyperlinks.
func HyperlinkRelationshipXML(relID, target string, external bool) string {
	targetMode := ""
	if external {
		targetMode = ` TargetMode="External"`
	}
	return fmt.Sprintf(
		`<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink" Target="%s"%s/>`,
		Escape(relID),
		Escape(target),
		targetMode,
	)
}
