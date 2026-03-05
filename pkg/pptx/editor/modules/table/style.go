package table

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func SetTableStyleInFrame(frame []byte, styleGUID string) ([]byte, error) {
	newStyleTag := fmt.Sprintf(`<a:tableStyleId>%s</a:tableStyleId>`, common.XMLEscape(styleGUID))
	tblPrStart, tblPrEnd, tblPrOpenTag, selfClosing, err := findTblPrBounds(frame)
	if err != nil {
		return nil, err
	}

	if selfClosing {
		return expandSelfClosingTblPr(frame, tblPrStart, tblPrEnd, tblPrOpenTag, newStyleTag)
	}
	return upsertTableStyleTag(frame, tblPrStart, tblPrEnd, newStyleTag)
}

func findTblPrBounds(frame []byte) (int, int, []byte, bool, error) {
	tblStart := bytes.Index(frame, []byte("<a:tbl"))
	if tblStart == -1 {
		return 0, 0, nil, false, errors.New("graphicFrame does not contain a table")
	}

	tblPrStart := bytes.Index(frame[tblStart:], []byte("<a:tblPr"))
	if tblPrStart == -1 {
		return 0, 0, nil, false, errors.New("table has no tblPr element")
	}
	tblPrStart += tblStart
	tblPrTagEnd := bytes.Index(frame[tblPrStart:], []byte(">"))
	if tblPrTagEnd == -1 {
		return 0, 0, nil, false, errors.New("invalid table tblPr element")
	}
	tblPrTagEnd += tblPrStart
	tblPrOpenTag := frame[tblPrStart : tblPrTagEnd+1]
	selfClosing := bytes.HasSuffix(bytes.TrimSpace(tblPrOpenTag), []byte("/>"))
	if selfClosing {
		return tblPrStart, tblPrTagEnd + 1, tblPrOpenTag, true, nil
	}

	tblPrEnd := bytes.Index(frame[tblPrStart:], []byte("</a:tblPr>"))
	if tblPrEnd == -1 {
		return 0, 0, nil, false, errors.New("invalid table tblPr element")
	}
	tblPrEnd += tblPrStart + len("</a:tblPr>")
	return tblPrStart, tblPrEnd, tblPrOpenTag, false, nil
}

func expandSelfClosingTblPr(frame []byte, start, end int, openTag []byte, newStyleTag string) ([]byte, error) {
	openTagText := string(openTag)
	trimmed := strings.TrimRight(openTagText, " \t\r\n")
	if !strings.HasSuffix(trimmed, "/>") {
		return nil, errors.New("invalid tblPr element")
	}
	expandedTblPr := trimmed[:len(trimmed)-2] + ">" + newStyleTag + `</a:tblPr>`
	updatedFrame := make([]byte, 0, len(frame)+len(expandedTblPr))
	updatedFrame = append(updatedFrame, frame[:start]...)
	updatedFrame = append(updatedFrame, []byte(expandedTblPr)...)
	updatedFrame = append(updatedFrame, frame[end:]...)
	return updatedFrame, nil
}

func upsertTableStyleTag(frame []byte, tblPrStart, tblPrEnd int, newStyleTag string) ([]byte, error) {
	tblPrSection := frame[tblPrStart:tblPrEnd]
	styleIDStart := bytes.Index(tblPrSection, []byte("<a:tableStyleId"))
	if styleIDStart != -1 {
		return replaceTableStyleTag(frame, tblPrStart+styleIDStart, newStyleTag)
	}
	return insertTableStyleTag(frame, tblPrSection, tblPrStart, newStyleTag)
}

func replaceTableStyleTag(frame []byte, styleIDStartInFrame int, newStyleTag string) ([]byte, error) {
	styleIDEnd := bytes.Index(frame[styleIDStartInFrame:], []byte("</a:tableStyleId>"))
	if styleIDEnd == -1 {
		return nil, errors.New("invalid tableStyleId element")
	}
	styleIDEnd += styleIDStartInFrame + len("</a:tableStyleId>")
	updatedFrame := make([]byte, 0, len(frame))
	updatedFrame = append(updatedFrame, frame[:styleIDStartInFrame]...)
	updatedFrame = append(updatedFrame, []byte(newStyleTag)...)
	updatedFrame = append(updatedFrame, frame[styleIDEnd:]...)
	return updatedFrame, nil
}

func insertTableStyleTag(frame, tblPrSection []byte, tblPrStart int, newStyleTag string) ([]byte, error) {
	insertPos := bytes.Index(tblPrSection, []byte(">"))
	if insertPos == -1 {
		return nil, errors.New("invalid tblPr element")
	}
	insertPos += tblPrStart + 1
	updatedFrame := make([]byte, 0, len(frame)+len(newStyleTag))
	updatedFrame = append(updatedFrame, frame[:insertPos]...)
	updatedFrame = append(updatedFrame, []byte(newStyleTag)...)
	updatedFrame = append(updatedFrame, frame[insertPos:]...)
	return updatedFrame, nil
}
