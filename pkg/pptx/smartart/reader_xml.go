package smartart

import "regexp"

var (
	reSmartArtNodeText = regexp.MustCompile(`<a:t>([^<]*)</a:t>`)
	// A node's own fill sits in its shape properties, ahead of its text.
	reSmartArtNodeFillColor = regexp.MustCompile(
		`<dgm:spPr>\s*<a:solidFill>\s*<a:srgbClr val="([0-9A-Fa-f]{6})"`,
	)
	reSmartArtBlipEmbed   = regexp.MustCompile(`<a:blip[^>]*r:embed="([^"]+)"`)
	reSmartArtPresAssocID = regexp.MustCompile(`presAssocID="([^"]+)"`)
	reSmartArtUniqueID    = regexp.MustCompile(`uniqueId\s*=\s*["']([^"']+)["']`)
	reSmartArtLayoutFull  = regexp.MustCompile(`dgm:layoutDef[^>]*uniqueId\s*=\s*["']([^"']+)["']`)
)

// ExtractLayoutURI reads the SmartArt layout URI from a layout part.
func ExtractLayoutURI(layoutXML string) string {
	if match := reSmartArtUniqueID.FindStringSubmatch(layoutXML); match != nil {
		return match[1]
	}
	if match := reSmartArtLayoutFull.FindStringSubmatch(layoutXML); match != nil {
		return match[1]
	}
	return ""
}

// ExtractUniqueID reads a uniqueId attribute from SmartArt style/color/layout parts.
func ExtractUniqueID(partXML string) string {
	if match := reSmartArtUniqueID.FindStringSubmatch(partXML); match != nil {
		return match[1]
	}
	return ""
}
