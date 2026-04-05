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

// parseMasterTitleStyle parses the slide master XML using a token scanner,
// navigating to txStyles > titleStyle > lvl1pPr to read alignment and default font size.
// This approach is namespace-safe for OOXML documents.
func parseMasterTitleStyle(data []byte) masterTitleStyle {
	dec := xml.NewDecoder(bytes.NewReader(data))

	const (
		elemTxStyles    = "txStyles"
		elemTitleStyle  = "titleStyle"
		elemLvl1pPr     = "lvl1pPr"
		elemDefRPr      = "defRPr"
		attrAlgn        = "algn"
		attrSz          = "sz"
		centiToPoints   = 100
	)

	inTxStyles   := false
	inTitleStyle := false
	result       := masterTitleStyle{}
	done         := false

	for !done {
		tok, err := dec.Token()
		if err != nil {
			break
		}

		switch t := tok.(type) {
		case xml.StartElement:
			name := t.Name.Local
			switch {
			case name == elemTxStyles:
				inTxStyles = true
			case inTxStyles && name == elemTitleStyle:
				inTitleStyle = true
			case inTitleStyle && name == elemLvl1pPr:
				// Read paragraph properties (alignment)
				for _, a := range t.Attr {
					if a.Name.Local == attrAlgn {
						result.Align = a.Value
					}
				}
			case inTitleStyle && name == elemDefRPr && result.SizePt == 0:
				// Read default run properties (font size) — take first occurrence
				for _, a := range t.Attr {
					if a.Name.Local == attrSz {
						if sz, err2 := strconv.Atoi(a.Value); err2 == nil && sz > 0 {
							result.SizePt = sz / centiToPoints
						}
					}
				}
			}
		case xml.EndElement:
			name := t.Name.Local
			switch {
			case name == elemTitleStyle && inTitleStyle:
				// Finished reading titleStyle — we have what we need.
				done = true
			case name == elemTxStyles:
				done = true
			}
		}
	}
	return result
}
