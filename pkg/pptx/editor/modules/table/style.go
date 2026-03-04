package table

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func SetTableStyleInFrame(frame []byte, styleGUID string) ([]byte, error) {
	// Extract the table XML to locate the tableStyleId element
	tblStart := bytes.Index(frame, []byte("<a:tbl"))
	if tblStart == -1 {
		return nil, errors.New("graphicFrame does not contain a table")
	}

	// Look for existing tableStyleId element
	tblPrStart := bytes.Index(frame[tblStart:], []byte("<a:tblPr"))
	if tblPrStart == -1 {
		return nil, errors.New("table has no tblPr element")
	}
	tblPrStart += tblStart
	tblPrTagEnd := bytes.Index(frame[tblPrStart:], []byte(">"))
	if tblPrTagEnd == -1 {
		return nil, errors.New("invalid table tblPr element")
	}
	tblPrTagEnd += tblPrStart
	tblPrOpenTag := frame[tblPrStart : tblPrTagEnd+1]
	tblPrSelfClosing := bytes.HasSuffix(bytes.TrimSpace(tblPrOpenTag), []byte("/>"))

	var tblPrEnd int
	if tblPrSelfClosing {
		tblPrEnd = tblPrTagEnd + 1
	} else {
		tblPrEnd = bytes.Index(frame[tblPrStart:], []byte("</a:tblPr>"))
		if tblPrEnd == -1 {
			return nil, errors.New("invalid table tblPr element")
		}
		tblPrEnd += tblPrStart + len("</a:tblPr>")
	}

	newStyleTag := fmt.Sprintf(`<a:tableStyleId>%s</a:tableStyleId>`, common.XMLEscape(styleGUID))
	if tblPrSelfClosing {
		openTag := string(tblPrOpenTag)
		tagCloseIdx := strings.LastIndex(openTag, "/>")
		if tagCloseIdx == -1 {
			return nil, errors.New("invalid tblPr element")
		}
		expandedTblPr := openTag[:tagCloseIdx] + ">" + newStyleTag + `</a:tblPr>`
		updatedFrame := make([]byte, 0, len(frame)+len(expandedTblPr))
		updatedFrame = append(updatedFrame, frame[:tblPrStart]...)
		updatedFrame = append(updatedFrame, []byte(expandedTblPr)...)
		updatedFrame = append(updatedFrame, frame[tblPrEnd:]...)
		return updatedFrame, nil
	}

	// Check if tableStyleId already exists in tblPr
	tblPrSection := frame[tblPrStart:tblPrEnd]
	styleIDStart := bytes.Index(tblPrSection, []byte("<a:tableStyleId"))
	var updatedFrame []byte

	if styleIDStart != -1 {
		// Update existing tableStyleId
		styleIDStartInFrame := tblPrStart + styleIDStart
		styleIDEnd := bytes.Index(frame[styleIDStartInFrame:], []byte("</a:tableStyleId>"))
		if styleIDEnd == -1 {
			return nil, errors.New("invalid tableStyleId element")
		}
		styleIDEnd += styleIDStartInFrame + len("</a:tableStyleId>")
		updatedFrame = make([]byte, 0, len(frame))
		updatedFrame = append(updatedFrame, frame[:styleIDStartInFrame]...)
		updatedFrame = append(updatedFrame, []byte(newStyleTag)...)
		updatedFrame = append(updatedFrame, frame[styleIDEnd:]...)
		return updatedFrame, nil
	}

	// Insert new tableStyleId at start of tblPr
	insertPos := bytes.Index(tblPrSection, []byte(">"))
	if insertPos == -1 {
		return nil, errors.New("invalid tblPr element")
	}
	insertPos += tblPrStart + 1
	updatedFrame = make([]byte, 0, len(frame)+len(newStyleTag))
	updatedFrame = append(updatedFrame, frame[:insertPos]...)
	updatedFrame = append(updatedFrame, []byte(newStyleTag)...)
	updatedFrame = append(updatedFrame, frame[insertPos:]...)
	return updatedFrame, nil
}
