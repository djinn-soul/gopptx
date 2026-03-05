package shape

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"regexp"
	"strconv"
)

func ExtractGroupChildShapeNodes(groupXML []byte, groupShapeTag string, pictureShapeType string) ([][]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(groupXML))
	rootDepth := 0
	rootSeen := false
	children := make([][]byte, 0)

	for {
		startOffset := decoder.InputOffset()
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			nextDepth, sawRoot, childNode, startErr := processGroupStartElement(
				decoder,
				groupXML,
				startOffset,
				rootDepth,
				rootSeen,
				t,
				groupShapeTag,
				pictureShapeType,
			)
			if startErr != nil {
				return nil, startErr
			}
			rootDepth = nextDepth
			rootSeen = sawRoot
			if len(childNode) > 0 {
				children = append(children, childNode)
			}
		case xml.EndElement:
			var done bool
			rootDepth, done = processGroupEndElement(rootDepth, rootSeen)
			if done {
				return children, nil
			}
		}
	}

	return children, nil
}

func processGroupStartElement(
	decoder *xml.Decoder,
	groupXML []byte,
	startOffset int64,
	rootDepth int,
	rootSeen bool,
	start xml.StartElement,
	groupShapeTag string,
	pictureShapeType string,
) (int, bool, []byte, error) {
	if !rootSeen {
		if start.Name.Local == groupShapeTag {
			return 1, true, nil, nil
		}
		return rootDepth, false, nil, nil
	}
	if rootDepth == 1 && isShapeElementLocal(start.Name.Local, groupShapeTag, pictureShapeType) {
		endOffset, err := findElementEndOffset(decoder, start.Name.Local)
		if err != nil {
			return rootDepth, rootSeen, nil, err
		}
		child := bytes.TrimSpace(groupXML[startOffset:endOffset])
		return rootDepth, rootSeen, child, nil
	}
	return rootDepth + 1, rootSeen, nil, nil
}

func processGroupEndElement(rootDepth int, rootSeen bool) (int, bool) {
	if !rootSeen {
		return rootDepth, false
	}
	rootDepth--
	return rootDepth, rootDepth == 0
}

func isShapeElementLocal(name, groupShapeTag, pictureShapeType string) bool {
	switch name {
	case "sp", pictureShapeType, "graphicFrame", groupShapeTag, "cxnSp":
		return true
	default:
		return false
	}
}

func findElementEndOffset(decoder *xml.Decoder, stopTag string) (int64, error) {
	depth := 1
	for {
		token, err := decoder.Token()
		if err != nil {
			return 0, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == stopTag {
				depth++
			}
		case xml.EndElement:
			if t.Name.Local != stopTag {
				continue
			}
			depth--
			if depth == 0 {
				return decoder.InputOffset(), nil
			}
		}
	}
}

func FirstShapeIDInXML(shapesXML [][]byte, cNvPrIDPattern *regexp.Regexp, cNvPrSubmatchSize int) int {
	for _, shapeXML := range shapesXML {
		match := cNvPrIDPattern.FindSubmatch(shapeXML)
		if len(match) < cNvPrSubmatchSize {
			continue
		}
		id, err := strconv.Atoi(string(match[1]))
		if err == nil {
			return id
		}
	}
	return 0
}
