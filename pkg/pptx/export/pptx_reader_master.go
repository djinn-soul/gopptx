package export

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"strconv"
)

// masterTitleStyle holds the inherited title defaults from the slide master's txStyles.
type masterTitleStyle struct {
	Align  string // e.g. "ctr", "l", "r"
	SizePt int    // font size in points (0 = unset)
}

// extractMasterTitleStyle reads the primary slide master's txStyles.titleStyle
// to determine the inherited alignment and font size for title placeholders.
// Returns a zero-value masterTitleStyle on any error.
func extractMasterTitleStyle(pptxPath string) masterTitleStyle {
	zr, err := zip.OpenReader(pptxPath)
	if err != nil {
		return masterTitleStyle{}
	}
	defer zr.Close()

	fileMap := make(map[string]*zip.File, len(zr.File))
	for _, f := range zr.File {
		fileMap[canonicalZipPath(f.Name)] = f
	}

	data := readZipBytes(fileMap, "ppt/slideMasters/slideMaster1.xml")
	if data == nil {
		return masterTitleStyle{}
	}

	return parseMasterTitleStyle(data)
}

// masterStyleParser holds incremental state while scanning the slide master XML.
type masterStyleParser struct {
	inTxStyles   bool
	inTitleStyle bool
	done         bool
	result       masterTitleStyle
}

func (p *masterStyleParser) onStart(t xml.StartElement) {
	name := t.Name.Local
	switch {
	case name == "txStyles":
		p.inTxStyles = true
	case p.inTxStyles && name == "titleStyle":
		p.inTitleStyle = true
	case p.inTitleStyle && name == "lvl1pPr":
		p.result.Align = attrValue(t.Attr, "algn")
	case p.inTitleStyle && name == "defRPr" && p.result.SizePt == 0:
		p.result.SizePt = parseSzAttr(t.Attr)
	}
}

func (p *masterStyleParser) onEnd(t xml.EndElement) {
	switch t.Name.Local {
	case "titleStyle":
		if p.inTitleStyle {
			p.done = true
		}
	case "txStyles":
		p.done = true
	}
}

// attrValue returns the value of the named attribute, or "" if absent.
func attrValue(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}

// parseSzAttr reads the "sz" attribute (OOXML hundredths-of-a-point) and converts to points.
func parseSzAttr(attrs []xml.Attr) int {
	const centiToPoints = 100
	for _, a := range attrs {
		if a.Name.Local == "sz" {
			if sz, err := strconv.Atoi(a.Value); err == nil && sz > 0 {
				return sz / centiToPoints
			}
		}
	}
	return 0
}

// parseMasterTitleStyle parses the slide master XML using a token scanner,
// navigating to txStyles > titleStyle > lvl1pPr to read alignment and default font size.
// This approach is namespace-safe for OOXML documents.
func parseMasterTitleStyle(data []byte) masterTitleStyle {
	dec := xml.NewDecoder(bytes.NewReader(data))
	p := &masterStyleParser{}
	for !p.done {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			p.onStart(t)
		case xml.EndElement:
			p.onEnd(t)
		}
	}
	return p.result
}
