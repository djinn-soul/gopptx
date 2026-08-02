package pptxxml

import (
	"fmt"
	"strings"
)

// HyperlinkSpec contains data needed for hyperlink XML generation.
type HyperlinkSpec struct {
	RelID          string
	Tooltip        string
	HighlightClick bool
	History        *bool
	EndSound       *bool
	Action         string // ppaction:// for internal navigation
	// UseTextColor keeps the run's own colour instead of the theme's hyperlink
	// colour. It is set when the run states a colour of its own.
	UseTextColor bool
}

// HyperlinkXML generates drawing hyperlink elements such as
// <a:hlinkClick> or <a:hlinkMouseOver>.
func HyperlinkXML(spec HyperlinkSpec, tagName string) string {
	xml := `<` + tagName

	if strings.TrimSpace(spec.RelID) != "" {
		xml += ` r:id="` + FastEscapeRID(spec.RelID) + `"`
	}

	if spec.Tooltip != "" {
		xml += fmt.Sprintf(` tooltip="%s"`, Escape(spec.Tooltip))
	}

	if spec.HighlightClick {
		xml += ` highlightClick="1"`
	}

	if spec.Action != "" {
		xml += fmt.Sprintf(` action="%s"`, Escape(spec.Action))
	}
	if spec.History != nil {
		xml += fmt.Sprintf(` history="%d"`, boolToInt(*spec.History))
	}
	if spec.EndSound != nil {
		xml += fmt.Sprintf(` endSnd="%d"`, boolToInt(*spec.EndSound))
	}

	if spec.UseTextColor {
		return xml + ">" + hyperlinkTextColorOverrideXML() + "</" + tagName + ">"
	}
	xml += "/>"
	return xml
}

// Office's hyperlink-color extension. Without it PowerPoint paints every
// hyperlink run in the theme's hlink colour and the run's own solidFill is
// never seen, which is the "impossible to set the font color" of upstream issue
// #940. val="tx" tells PowerPoint to use the run's text colour instead.
const (
	hyperlinkColorExtURI = "{A12FA001-AC4F-418D-AE19-62706E023703}"
	hyperlinkColorNS     = "http://schemas.microsoft.com/office/drawing/2018/hyperlinkcolor"
)

func hyperlinkTextColorOverrideXML() string {
	return `<a:extLst><a:ext uri="` + hyperlinkColorExtURI + `">` +
		`<ahyp:hlinkClr xmlns:ahyp="` + hyperlinkColorNS + `" val="tx"/>` +
		`</a:ext></a:extLst>`
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// HyperlinkRelationshipXML generates a relationship element for hyperlink-like targets.
func HyperlinkRelationshipXML(relID, target string, external bool, relType string) string {
	typeValue := strings.TrimSpace(relType)
	if typeValue == "" {
		typeValue = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink"
	}
	targetMode := ""
	if external {
		targetMode = ` TargetMode="External"`
	}
	return fmt.Sprintf(
		`<Relationship Id="%s" Type="%s" Target="%s"%s/>`,
		FastEscapeRID(relID),
		Escape(typeValue),
		Escape(target),
		targetMode,
	)
}
