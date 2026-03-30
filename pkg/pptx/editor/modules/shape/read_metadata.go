package shape

import (
	"encoding/xml"
	"strings"
)

// ReaderHyperlinkRef captures the raw OOXML hyperlink/action reference before
// relationship targets are resolved against the slide rels part.
type ReaderHyperlinkRef struct {
	RelID          string
	Action         *string
	Tooltip        *string
	History        *bool
	HighlightClick *bool
	EndSound       *bool
}

// ReaderRunActions captures click/hover actions for one text run.
type ReaderRunActions struct {
	ClickAction *ReaderHyperlinkRef
	HoverAction *ReaderHyperlinkRef
}

// ParsedShapeReaderMetadata captures shape metadata that depends on cNvPr or
// run-level hyperlink tags and is resolved later by the editor reader.
type ParsedShapeReaderMetadata struct {
	AltText      string
	HasAltText   bool
	IsDecorative bool
	ClickAction  *ReaderHyperlinkRef
	HoverAction  *ReaderHyperlinkRef
	RunActions   [][]ReaderRunActions
}

type readerHyperlinkXML struct {
	Action         *string    `xml:"action,attr"`
	Tooltip        *string    `xml:"tooltip,attr"`
	History        *bool      `xml:"history,attr"`
	HighlightClick *bool      `xml:"highlightClick,attr"`
	EndSnd         *bool      `xml:"endSnd,attr"`
	Attrs          []xml.Attr `xml:",any,attr"`
}

type cNvPrReaderXML struct {
	Descr          *string             `xml:"descr,attr"`
	Title          *string             `xml:"title,attr"`
	HlinkClick     *readerHyperlinkXML `xml:"hlinkClick"`
	HlinkHover     *readerHyperlinkXML `xml:"hlinkHover"`
	HlinkMouseOver *readerHyperlinkXML `xml:"hlinkMouseOver"`
}

type shapeReaderXML struct {
	NvSpPr struct {
		CNvPr cNvPrReaderXML `xml:"cNvPr"`
	} `xml:"nvSpPr"`
	NvPicPr struct {
		CNvPr cNvPrReaderXML `xml:"cNvPr"`
	} `xml:"nvPicPr"`
	NvCxnSpPr struct {
		CNvPr cNvPrReaderXML `xml:"cNvPr"`
	} `xml:"nvCxnSpPr"`
	NvGrpSpPr struct {
		CNvPr cNvPrReaderXML `xml:"cNvPr"`
	} `xml:"nvGrpSpPr"`
	NvGraphicFramePr struct {
		CNvPr cNvPrReaderXML `xml:"cNvPr"`
	} `xml:"nvGraphicFramePr"`
	TxBody struct {
		P []struct {
			R []struct {
				RPr *struct {
					HlinkClick     *readerHyperlinkXML `xml:"hlinkClick"`
					HlinkMouseOver *readerHyperlinkXML `xml:"hlinkMouseOver"`
				} `xml:"rPr"`
			} `xml:"r"`
		} `xml:"p"`
	} `xml:"txBody"`
}

func ParseShapeReaderMetadata(content []byte) (ParsedShapeReaderMetadata, error) {
	var src shapeReaderXML
	if err := xml.Unmarshal(content, &src); err != nil {
		return ParsedShapeReaderMetadata{}, err
	}

	meta := ParsedShapeReaderMetadata{}
	cNvPr := firstReaderCNvPr(&src)
	applyReaderAltText(&meta, cNvPr)
	if cNvPr != nil {
		meta.ClickAction = readerHyperlinkRef(cNvPr.HlinkClick)
		meta.HoverAction = readerHyperlinkRef(firstReaderHyperlink(cNvPr.HlinkHover, cNvPr.HlinkMouseOver))
	}
	meta.RunActions = make([][]ReaderRunActions, 0, len(src.TxBody.P))
	for _, paragraph := range src.TxBody.P {
		runActions := make([]ReaderRunActions, 0, len(paragraph.R))
		for _, run := range paragraph.R {
			actions := ReaderRunActions{}
			if run.RPr != nil {
				actions.ClickAction = readerHyperlinkRef(run.RPr.HlinkClick)
				actions.HoverAction = readerHyperlinkRef(run.RPr.HlinkMouseOver)
			}
			runActions = append(runActions, actions)
		}
		meta.RunActions = append(meta.RunActions, runActions)
	}
	return meta, nil
}

func firstReaderCNvPr(src *shapeReaderXML) *cNvPrReaderXML {
	switch {
	case src.NvSpPr.CNvPr != (cNvPrReaderXML{}):
		return &src.NvSpPr.CNvPr
	case src.NvPicPr.CNvPr != (cNvPrReaderXML{}):
		return &src.NvPicPr.CNvPr
	case src.NvCxnSpPr.CNvPr != (cNvPrReaderXML{}):
		return &src.NvCxnSpPr.CNvPr
	case src.NvGrpSpPr.CNvPr != (cNvPrReaderXML{}):
		return &src.NvGrpSpPr.CNvPr
	case src.NvGraphicFramePr.CNvPr != (cNvPrReaderXML{}):
		return &src.NvGraphicFramePr.CNvPr
	default:
		return nil
	}
}

func applyReaderAltText(meta *ParsedShapeReaderMetadata, cNvPr *cNvPrReaderXML) {
	if cNvPr == nil {
		return
	}
	if cNvPr.Descr != nil {
		if descr := strings.TrimSpace(*cNvPr.Descr); descr != "" {
			meta.AltText = descr
			meta.HasAltText = true
			return
		}
		meta.IsDecorative = true
		return
	}
	if cNvPr.Title != nil {
		if title := strings.TrimSpace(*cNvPr.Title); title != "" {
			meta.AltText = title
			meta.HasAltText = true
		}
	}
}

func firstReaderHyperlink(primary, secondary *readerHyperlinkXML) *readerHyperlinkXML {
	if primary != nil {
		return primary
	}
	return secondary
}

func readerHyperlinkRef(src *readerHyperlinkXML) *ReaderHyperlinkRef {
	if src == nil {
		return nil
	}
	ref := &ReaderHyperlinkRef{
		Action:         cloneStringPtr(src.Action),
		Tooltip:        cloneStringPtr(src.Tooltip),
		History:        cloneBoolPtr(src.History),
		HighlightClick: cloneBoolPtr(src.HighlightClick),
		EndSound:       cloneBoolPtr(src.EndSnd),
	}
	for _, attr := range src.Attrs {
		if attr.Name.Local == "id" {
			ref.RelID = strings.TrimSpace(attr.Value)
			break
		}
	}
	if ref.RelID == "" && ref.Action == nil && ref.Tooltip == nil && ref.History == nil &&
		ref.HighlightClick == nil && ref.EndSound == nil {
		return nil
	}
	return ref
}

func cloneStringPtr(src *string) *string {
	if src == nil {
		return nil
	}
	value := strings.TrimSpace(*src)
	return &value
}

func cloneBoolPtr(src *bool) *bool {
	if src == nil {
		return nil
	}
	value := *src
	return &value
}
