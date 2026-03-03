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

func findTableFrame(slideContent []byte, shapeID int) (int, int, []byte, error) {
	idStr := fmt.Appendf(nil, ` id="%d"`, shapeID)
	idIdx := bytes.Index(slideContent, idStr)
	if idIdx == -1 {
		return 0, 0, nil, fmt.Errorf("shape id %d not found", shapeID)
	}

	frameStart := bytes.LastIndex(slideContent[:idIdx], []byte("<p:graphicFrame"))
	if frameStart == -1 {
		return 0, 0, nil, fmt.Errorf("shape %d is not a graphicFrame", shapeID)
	}

	relEnd := bytes.Index(slideContent[idIdx:], []byte("</p:graphicFrame>"))
	if relEnd == -1 {
		return 0, 0, nil, errors.New("invalid graphicFrame xml")
	}
	frameEnd := idIdx + relEnd + len("</p:graphicFrame>")
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
	tblStart := bytes.Index(frame, []byte("<a:tbl"))
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
	string,
	[]byte,
	int,
	int,
	[]byte,
	error,
) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return "", nil, 0, 0, nil, fmt.Errorf("slide index %d out of range", slideIndex)
	}
	partPath := e.slides[slideIndex].Part
	var ok bool
	var slideContent []byte
	slideContent, ok = e.parts.Get(partPath)
	if !ok {
		return "", nil, 0, 0, nil, errors.New("slide part not found")
	}
	frameStart, frameEnd, frame, err := findTableFrame(slideContent, shapeID)
	if err != nil {
		return "", nil, 0, 0, nil, err
	}
	return partPath, slideContent, frameStart, frameEnd, frame, nil
}

// SetTableStyle sets the table style for the specified table on a slide.
// The styleGuid must be a valid PowerPoint table style GUID, e.g.:
//
//	"{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}" - Medium Style 2 - Accent 1
//	"{B9AC3A68-259E-4EED-9050-4AE35E7F2B2D}" - Light Style 1
//	"{5940675A-B579-460E-94D1-54222C63F5DA}" - Medium Style 1 - Accent 1
func (e *PresentationEditor) SetTableStyle(slideIndex, shapeID int, styleGuid string) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}

	// Extract the table XML to locate the tableStyleId element
	tblStart := bytes.Index(frame, []byte("<a:tbl"))
	if tblStart == -1 {
		return errors.New("graphicFrame does not contain a table")
	}

	// Look for existing tableStyleId element
	tblPrStart := bytes.Index(frame[tblStart:], []byte("<a:tblPr"))
	if tblPrStart == -1 {
		return errors.New("table has no tblPr element")
	}
	tblPrStart += tblStart
	tblPrEnd := bytes.Index(frame[tblPrStart:], []byte("</a:tblPr>"))
	if tblPrEnd == -1 {
		return errors.New("invalid table tblPr element")
	}
	tblPrEnd += tblPrStart + len("</a:tblPr>")

	// Check if tableStyleId already exists in tblPr
	tblPrSection := frame[tblPrStart:tblPrEnd]
	styleIdStart := bytes.Index(tblPrSection, []byte("<a:tableStyleId"))
	var updatedFrame []byte

	if styleIdStart != -1 {
		// Update existing tableStyleId
		styleIdStartInFrame := tblPrStart + styleIdStart
		styleIdEnd := bytes.Index(frame[styleIdStartInFrame:], []byte("</a:tableStyleId>"))
		if styleIdEnd == -1 {
			return errors.New("invalid tableStyleId element")
		}
		styleIdEnd += styleIdStartInFrame + len("</a:tableStyleId>")
		newStyleTag := fmt.Sprintf(`<a:tableStyleId>%s</a:tableStyleId>`, styleGuid)
		updatedFrame = make([]byte, 0, len(frame))
		updatedFrame = append(updatedFrame, frame[:styleIdStartInFrame]...)
		updatedFrame = append(updatedFrame, []byte(newStyleTag)...)
		updatedFrame = append(updatedFrame, frame[styleIdEnd:]...)
	} else {
		// Insert new tableStyleId after firstRow or at start of tblPr
		insertPos := bytes.Index(tblPrSection, []byte(">"))
		if insertPos == -1 {
			return errors.New("invalid tblPr element")
		}
		insertPos += tblPrStart + 1
		newStyleTag := fmt.Sprintf(`<a:tableStyleId>%s</a:tableStyleId>`, styleGuid)
		updatedFrame = make([]byte, 0, len(frame)+len(newStyleTag))
		updatedFrame = append(updatedFrame, frame[:insertPos]...)
		updatedFrame = append(updatedFrame, []byte(newStyleTag)...)
		updatedFrame = append(updatedFrame, frame[insertPos:]...)
	}

	// Replace the entire frame in the slide content
	updatedSlide := replaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame)
	e.parts.Set(partPath, updatedSlide)
	return nil
}
