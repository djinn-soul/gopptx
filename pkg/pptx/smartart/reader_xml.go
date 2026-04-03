package smartart

import "regexp"

var (
	reSmartArtNodeText   = regexp.MustCompile(`<a:t>([^<]*)</a:t>`)
	reSmartArtUniqueID   = regexp.MustCompile(`uniqueId\s*=\s*["']([^"']+)["']`)
	reSmartArtLayoutFull = regexp.MustCompile(`dgm:layoutDef[^>]*uniqueId\s*=\s*["']([^"']+)["']`)
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
