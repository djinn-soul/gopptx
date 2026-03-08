package editor

import (
	"bytes"
	"errors"
	"regexp"
)

var shapeCloseTagPattern = regexp.MustCompile(`</p:(sp|pic|graphicFrame|grpSp|cxnSp)>`)

func insertShapeXML(content, shapeXML []byte) ([]byte, error) {
	insertAt, err := findShapeInsertOffset(content)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Grow(len(content) + len(shapeXML))
	buf.Write(content[:insertAt])
	buf.Write(shapeXML)
	buf.Write(content[insertAt:])
	return buf.Bytes(), nil
}

func findShapeInsertOffset(content []byte) (int, error) {
	matches := shapeCloseTagPattern.FindAllIndex(content, -1)
	if len(matches) > 0 {
		return matches[len(matches)-1][1], nil
	}

	idx := bytes.LastIndex(content, []byte("</p:spTree>"))
	if idx == -1 {
		return 0, errors.New("invalid slide xml: missing spTree end")
	}
	return idx, nil
}
