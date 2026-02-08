package pptx

import (
	"encoding/xml"
	"fmt"
)

// slideXML represents the structure of a slide for parsing.
type slideXML struct {
	XMLName xml.Name `xml:"sld"`
	CSld    struct {
		SpTree struct {
			Sp  []shapeXML `xml:"sp"`
			Pic []picXML   `xml:"pic"`
			// Groups?
		} `xml:"spTree"`
	} `xml:"cSld"`
}

type shapeXML struct {
	NvSpPr struct {
		CNvPr cNvPr `xml:"cNvPr"`
		NvPr  nvPr  `xml:"nvPr"`
	} `xml:"nvSpPr"`
	SpPr spPr `xml:"spPr"`
}

type picXML struct {
	NvPicPr struct {
		CNvPr cNvPr `xml:"cNvPr"`
		NvPr  nvPr  `xml:"nvPr"`
	} `xml:"nvPicPr"`
	SpPr spPr `xml:"spPr"`
}

type cNvPr struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type nvPr struct {
	Ph *placeholderXML `xml:"ph"`
}

type spPr struct {
	Xfrm xfrm `xml:"xfrm"`
}

type xfrm struct {
	Off offset `xml:"off"`
	Ext extent `xml:"ext"`
}

type offset struct {
	X int64 `xml:"x,attr"`
	Y int64 `xml:"y,attr"`
}

type extent struct {
	CX int64 `xml:"cx,attr"`
	CY int64 `xml:"cy,attr"`
}

type placeholderXML struct {
	Type   string `xml:"type,attr"`
	Idx    int    `xml:"idx,attr"`
	Sz     string `xml:"sz,attr"`
	Orient string `xml:"orient,attr"`
}

// ParseSlidePlaceholders parses slide XML content to extract placeholders.
func ParseSlidePlaceholders(content []byte) ([]Placeholder, error) {
	var s slideXML
	// We might need to handle namespaces? The structs should match local names.
	if err := xml.Unmarshal(content, &s); err != nil {
		return nil, fmt.Errorf("unmarshal slide xml: %w", err)
	}

	var placeholders []Placeholder

	// Shapes
	for _, sp := range s.CSld.SpTree.Sp {
		if sp.NvSpPr.NvPr.Ph != nil {
			ph := sp.NvSpPr.NvPr.Ph
			placeholders = append(placeholders, Placeholder{
				Type:  PlaceholderType(ph.Type),
				Index: ph.Idx,
				Name:  sp.NvSpPr.CNvPr.Name,
				X:     sp.SpPr.Xfrm.Off.X,
				Y:     sp.SpPr.Xfrm.Off.Y,
				CX:    sp.SpPr.Xfrm.Ext.CX,
				CY:    sp.SpPr.Xfrm.Ext.CY,
			})
		}
	}

	// Pictures
	for _, pic := range s.CSld.SpTree.Pic {
		if pic.NvPicPr.NvPr.Ph != nil {
			ph := pic.NvPicPr.NvPr.Ph
			placeholders = append(placeholders, Placeholder{
				Type:  PlaceholderType(ph.Type),
				Index: ph.Idx,
				Name:  pic.NvPicPr.CNvPr.Name,
				X:     pic.SpPr.Xfrm.Off.X,
				Y:     pic.SpPr.Xfrm.Off.Y,
				CX:    pic.SpPr.Xfrm.Ext.CX,
				CY:    pic.SpPr.Xfrm.Ext.CY,
			})
		}
	}

	return placeholders, nil
}
