package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

type tableXML struct {
	TblPr struct {
		FirstRow string `xml:"firstRow,attr"`
		FirstCol string `xml:"firstCol,attr"`
		LastRow  string `xml:"lastRow,attr"`
		LastCol  string `xml:"lastCol,attr"`
		BandRow  string `xml:"bandRow,attr"`
		BandCol  string `xml:"bandCol,attr"`
	} `xml:"tblPr"`
	Grid struct {
		Cols []struct{} `xml:"gridCol"`
	} `xml:"tblGrid"`
	Rows []tableRowXML `xml:"tr"`
}

type tableRowXML struct {
	Cells []tableCellXML `xml:"tc"`
}

type tableCellXML struct {
	RowSpan  int    `xml:"rowSpan,attr"`
	GridSpan int    `xml:"gridSpan,attr"`
	VMerge   string `xml:"vMerge,attr"`
	HMerge   string `xml:"hMerge,attr"`
	TxBody   struct {
		Paragraphs []struct {
			Runs []struct {
				Text string `xml:"t"`
			} `xml:"r"`
		} `xml:"p"`
	} `xml:"txBody"`
}

func truthyAttr(v string) bool {
	return v == "1" || strings.EqualFold(v, "true")
}

func findTableFrame(slideContent []byte, shapeID int) (frameStart int, frameEnd int, frame []byte, err error) {
	idStr := []byte(fmt.Sprintf(` id="%d"`, shapeID))
	idIdx := bytes.Index(slideContent, idStr)
	if idIdx == -1 {
		return 0, 0, nil, fmt.Errorf("shape id %d not found", shapeID)
	}

	frameStart = bytes.LastIndex(slideContent[:idIdx], []byte("<p:graphicFrame"))
	if frameStart == -1 {
		return 0, 0, nil, fmt.Errorf("shape %d is not a graphicFrame", shapeID)
	}

	relEnd := bytes.Index(slideContent[idIdx:], []byte("</p:graphicFrame>"))
	if relEnd == -1 {
		return 0, 0, nil, errors.New("invalid graphicFrame xml")
	}
	frameEnd = idIdx + relEnd + len("</p:graphicFrame>")
	return frameStart, frameEnd, slideContent[frameStart:frameEnd], nil
}

func replaceTableFrame(slideContent []byte, frameStart, frameEnd int, frame []byte) []byte {
	updated := make([]byte, 0, len(slideContent)-((frameEnd-frameStart)-len(frame)))
	updated = append(updated, slideContent[:frameStart]...)
	updated = append(updated, frame...)
	updated = append(updated, slideContent[frameEnd:]...)
	return updated
}

func extractTableXML(frame []byte) ([]byte, error) {
	tblStart := bytes.Index(frame, []byte("<a:tbl>"))
	if tblStart == -1 {
		tblStart = bytes.Index(frame, []byte("<a:tbl "))
		if tblStart == -1 {
			tblStart = bytes.Index(frame, []byte("<a:tbl>"))
		}
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

func parseTable(frame []byte) (*tableXML, error) {
	tblBytes, err := extractTableXML(frame)
	if err != nil {
		return nil, err
	}
	var parsed tableXML
	if err := xml.Unmarshal(tblBytes, &parsed); err != nil {
		return nil, fmt.Errorf("parse table xml: %w", err)
	}
	return &parsed, nil
}

func tableDimensions(parsed *tableXML) (int, int) {
	rows := len(parsed.Rows)
	cols := len(parsed.Grid.Cols)
	if cols == 0 && rows > 0 {
		cols = len(parsed.Rows[0].Cells)
	}
	return rows, cols
}

func getSlideTableFrame(e *PresentationEditor, slideIndex, shapeID int) (
	partPath string,
	slideContent []byte,
	frameStart int,
	frameEnd int,
	frame []byte,
	err error,
) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		err = fmt.Errorf("slide index %d out of range", slideIndex)
		return
	}
	partPath = e.slides[slideIndex].Part
	slideContent, ok := e.parts.Get(partPath)
	if !ok {
		err = errors.New("slide part not found")
		return
	}
	frameStart, frameEnd, frame, err = findTableFrame(slideContent, shapeID)
	return
}
