package table

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

type XML struct {
	TblPr struct {
		FirstRow string `xml:"firstRow,attr"`
		FirstCol string `xml:"firstCol,attr"`
		LastRow  string `xml:"lastRow,attr"`
		LastCol  string `xml:"lastCol,attr"`
		BandRow  string `xml:"bandRow,attr"`
		BandCol  string `xml:"bandCol,attr"`
	} `xml:"tblPr"`
	Grid struct {
		Cols []struct {
			Width int64 `xml:"w,attr"`
		} `xml:"gridCol"`
	} `xml:"tblGrid"`
	Rows []RowXML `xml:"tr"`
}

type RowXML struct {
	Height int64     `xml:"h,attr"`
	Cells  []CellXML `xml:"tc"`
}

type CellXML struct {
	RowSpan  int    `xml:"rowSpan,attr"`
	GridSpan int    `xml:"gridSpan,attr"`
	VMerge   string `xml:"vMerge,attr"`
	HMerge   string `xml:"hMerge,attr"`
	TcPr     cellPropertiesXML `xml:"tcPr"`
	TxBody   struct {
		Paragraphs []struct {
			Runs []struct {
				Text string `xml:"t"`
			} `xml:"r"`
		} `xml:"p"`
	} `xml:"txBody"`
}

type cellPropertiesXML struct {
	Anchor string `xml:"anchor,attr"`
	MarL   *int64 `xml:"marL,attr"`
	MarR   *int64 `xml:"marR,attr"`
	MarT   *int64 `xml:"marT,attr"`
	MarB   *int64 `xml:"marB,attr"`
	LnL    *linePropertiesXML `xml:"lnL"`
	LnR    *linePropertiesXML `xml:"lnR"`
	LnT    *linePropertiesXML `xml:"lnT"`
	LnB    *linePropertiesXML `xml:"lnB"`
}

type linePropertiesXML struct {
	Width int64 `xml:"w,attr"`
	NoFill *struct{} `xml:"noFill"`
	PrstDash *struct {
		Val string `xml:"val,attr"`
	} `xml:"prstDash"`
	SolidFill *struct {
		SrgbClr *struct {
			Val string `xml:"val,attr"`
		} `xml:"srgbClr"`
		SchemeClr *struct {
			Val    string `xml:"val,attr"`
			LumMod *struct {
				Val string `xml:"val,attr"`
			} `xml:"lumMod"`
			LumOff *struct {
				Val string `xml:"val,attr"`
			} `xml:"lumOff"`
			Tint *struct {
				Val string `xml:"val,attr"`
			} `xml:"tint"`
			Shade *struct {
				Val string `xml:"val,attr"`
			} `xml:"shade"`
		} `xml:"schemeClr"`
	} `xml:"solidFill"`
}

func TruthyAttr(v string) bool {
	return v == "1" || strings.EqualFold(v, "true")
}

func FindTableFrame(slideContent []byte, shapeID int) (int, int, []byte, error) {
	idStr := fmt.Sprintf(` id="%d"`, shapeID)
	cursor := 0
	var idIdx int
	for {
		relIdx := bytes.Index(slideContent[cursor:], []byte(idStr))
		if relIdx == -1 {
			return 0, 0, nil, fmt.Errorf("shape id %d not found", shapeID)
		}
		idIdx = cursor + relIdx
		nextCharIdx := idIdx + len(idStr)
		if nextCharIdx < len(slideContent) {
			nextChar := slideContent[nextCharIdx]
			if nextChar == ' ' || nextChar == '"' || nextChar == '>' || nextChar == '/' {
				break // Exact match found
			}
		} else {
			break // Exact match at end of content
		}
		cursor = idIdx + 1
	}

	frameStart := bytes.LastIndex(slideContent[:idIdx], []byte("<p:graphicFrame"))
	if frameStart == -1 {
		return 0, 0, nil, fmt.Errorf("shape %d is not a graphicFrame", shapeID)
	}

	// Verify that no other shape start tag exists between frameStart and the ID.
	lastAnyTag := bytes.LastIndex(slideContent[:idIdx], []byte("<p:"))
	if lastAnyTag > frameStart {
		tagNameEnd := bytes.IndexAny(slideContent[lastAnyTag:], " >/")
		if tagNameEnd != -1 {
			actualTag := string(slideContent[lastAnyTag : lastAnyTag+tagNameEnd])
			if actualTag != "<p:graphicFrame" && actualTag != "<p:nvGraphicFramePr" && actualTag != "<p:cNvPr" {
				return 0, 0, nil, fmt.Errorf("shape %d is a %s, not a graphicFrame", shapeID, actualTag[3:])
			}
		}
	}

	relEnd := bytes.Index(slideContent[idIdx:], []byte("</p:graphicFrame>"))
	if relEnd == -1 {
		return 0, 0, nil, errors.New("invalid graphicFrame xml")
	}
	frameEnd := idIdx + relEnd + len("</p:graphicFrame>")
	return frameStart, frameEnd, slideContent[frameStart:frameEnd], nil
}

func ReplaceTableFrame(slideContent []byte, frameStart, frameEnd int, frame []byte) []byte {
	updated := make([]byte, 0, len(slideContent)-((frameEnd-frameStart)-len(frame)))
	updated = append(updated, slideContent[:frameStart]...)
	updated = append(updated, frame...)
	updated = append(updated, slideContent[frameEnd:]...)
	return updated
}

func ExtractTableXML(frame []byte) ([]byte, error) {
	tblStart := bytes.Index(frame, []byte("<a:tbl "))
	if tblStart == -1 {
		tblStart = bytes.Index(frame, []byte("<a:tbl>"))
	}
	if tblStart == -1 {
		return nil, errors.New("graphicFrame does not contain a table")
	}
	tblEnd := bytes.Index(frame[tblStart:], []byte("</a:tbl>"))
	if tblEnd == -1 {
		return nil, errors.New("invalid table xml")
	}
	tblEnd += tblStart + len("</a:tbl>")
	return frame[tblStart:tblEnd], nil
}

func ParseTable(frame []byte) (*XML, error) {
	tblBytes, err := ExtractTableXML(frame)
	if err != nil {
		return nil, err
	}
	var parsed XML
	if err := xml.Unmarshal(tblBytes, &parsed); err != nil {
		return nil, fmt.Errorf("parse table xml: %w", err)
	}
	return &parsed, nil
}

func Dimensions(parsed *XML) (int, int) {
	rows := len(parsed.Rows)
	cols := len(parsed.Grid.Cols)
	if cols == 0 && rows > 0 {
		for _, cell := range parsed.Rows[0].Cells {
			span := cell.GridSpan
			if span <= 0 {
				span = 1
			}
			cols += span
		}
	}
	return rows, cols
}
