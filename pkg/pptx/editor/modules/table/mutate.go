package table

import (
	"bytes"
	"fmt"
)

const attrInsertOverheadBytes = 4

func SetOrInsertAttr(openingTag []byte, attrName, attrValue string) []byte {
	attrStr := []byte(" " + attrName + `="`)
	idx := bytes.Index(openingTag, attrStr)
	if idx != -1 {
		valStart := idx + len(attrStr)
		valEndRel := bytes.Index(openingTag[valStart:], []byte(`"`))
		if valEndRel != -1 {
			valEnd := valStart + valEndRel
			updated := make([]byte, 0, len(openingTag)-((valEnd-valStart)-len(attrValue)))
			updated = append(updated, openingTag[:valStart]...)
			updated = append(updated, []byte(attrValue)...)
			updated = append(updated, openingTag[valEnd:]...)
			return updated
		}
	}

	insertAt := len(openingTag) - 1
	if insertAt > 0 && openingTag[insertAt-1] == '/' {
		insertAt--
	}
	updated := make([]byte, 0, len(openingTag)+len(attrName)+len(attrValue)+attrInsertOverheadBytes)
	updated = append(updated, openingTag[:insertAt]...)
	updated = append(updated, []byte(" "+attrName+`="`+attrValue+`"`)...)
	updated = append(updated, openingTag[insertAt:]...)
	return updated
}

func SetTcAttr(tcContent []byte, attrName, attrValue string) []byte {
	tagEnd := bytes.Index(tcContent, []byte(">"))
	if tagEnd == -1 {
		return tcContent
	}
	openTag := tcContent[:tagEnd+1]
	updatedTag := SetOrInsertAttr(openTag, attrName, attrValue)
	updated := make([]byte, 0, len(tcContent)-((len(openTag))-len(updatedTag)))
	updated = append(updated, updatedTag...)
	updated = append(updated, tcContent[tagEnd+1:]...)
	return updated
}

func RemoveTcAttr(tcContent []byte, attrName string) []byte {
	tagEnd := bytes.Index(tcContent, []byte(">"))
	if tagEnd == -1 {
		return tcContent
	}
	openTag := tcContent[:tagEnd+1]
	attrStr := []byte(" " + attrName + `="`)
	idx := bytes.Index(openTag, attrStr)
	if idx == -1 {
		return tcContent
	}
	valStart := idx + len(attrStr)
	valEndRel := bytes.Index(openTag[valStart:], []byte(`"`))
	if valEndRel == -1 {
		return tcContent
	}
	valEnd := valStart + valEndRel + 1
	updatedTag := make([]byte, 0, len(openTag)-(valEnd-idx))
	updatedTag = append(updatedTag, openTag[:idx]...)
	updatedTag = append(updatedTag, openTag[valEnd:]...)

	updated := make([]byte, 0, len(tcContent)-((len(openTag))-len(updatedTag)))
	updated = append(updated, updatedTag...)
	updated = append(updated, tcContent[tagEnd+1:]...)
	return updated
}

func MutateTableRows(
	frame []byte,
	rowStart int,
	rowEnd int,
	mutator func(row int, rowContent []byte) ([]byte, error),
) ([]byte, error) {
	return MutateTableElements(frame, []byte("<a:tr"), []byte("</a:tr>"), rowStart, rowEnd, "row", mutator)
}

func MutateTableCells(
	rowContent []byte,
	colStart int,
	colEnd int,
	mutator func(col int, cellContent []byte) ([]byte, error),
) ([]byte, error) {
	return MutateTableElements(
		rowContent,
		[]byte("<a:tc"),
		[]byte("</a:tc>"),
		colStart,
		colEnd,
		"col",
		mutator,
	)
}

func MutateTableElements(
	content []byte,
	openTag []byte,
	closeTag []byte,
	start int,
	end int,
	label string,
	mutator func(index int, cellContent []byte) ([]byte, error),
) ([]byte, error) {
	var out bytes.Buffer
	cursor := 0
	index := 0

	for {
		rel := bytes.Index(content[cursor:], openTag)
		if rel == -1 {
			out.Write(content[cursor:])
			break
		}
		elementStart := cursor + rel
		elementEndRel := bytes.Index(content[elementStart:], closeTag)
		if elementEndRel == -1 {
			return nil, fmt.Errorf("invalid %s xml at %s %d", string(openTag), label, index)
		}
		elementEnd := elementStart + elementEndRel + len(closeTag)

		out.Write(content[cursor:elementStart])
		elementContent := content[elementStart:elementEnd]
		if index >= start && index <= end {
			updated, err := mutator(index, elementContent)
			if err != nil {
				return nil, err
			}
			out.Write(updated)
		} else {
			out.Write(elementContent)
		}

		cursor = elementEnd
		index++
	}

	return out.Bytes(), nil
}
