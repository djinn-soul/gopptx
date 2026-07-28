package table

import (
	"bytes"

	editormodcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/common"
)

// Border lines come before the cell fill; appending one after an existing
// <a:solidFill> produces XML PowerPoint rejects, so a new element has to be
// spliced into position.
const (
	tcPrRankLnL = iota
	tcPrRankLnR
	tcPrRankLnT
	tcPrRankLnB
	tcPrRankLnTlToBr
	tcPrRankLnBlToTr
	tcPrRankCell3D
	tcPrRankFill
	tcPrRankHeaders
	tcPrRankTrailing
)

func tcPrChildRank(name string) int {
	switch name {
	case "a:lnL":
		return tcPrRankLnL
	case "a:lnR":
		return tcPrRankLnR
	case "a:lnT":
		return tcPrRankLnT
	case "a:lnB":
		return tcPrRankLnB
	case "a:lnTlToBr":
		return tcPrRankLnTlToBr
	case "a:lnBlToTr":
		return tcPrRankLnBlToTr
	case "a:cell3D":
		return tcPrRankCell3D
	case "a:noFill", "a:solidFill", "a:gradFill", "a:blipFill", "a:pattFill", "a:grpFill":
		return tcPrRankFill
	case "a:headers":
		return tcPrRankHeaders
	default:
		return tcPrRankTrailing
	}
}

// insertBorderInSchemaOrder places borderXML among the existing tcPr children so
// the sequence stays schema-valid. It falls back to appending when the content
// cannot be split into plain elements.
func insertBorderInSchemaOrder(inner []byte, tag string, borderXML string) []byte {
	rank := tcPrChildRank("a:" + tag)
	insertAt := -1
	for _, element := range editormodcommon.SplitTopLevelXMLElements(string(inner)) {
		if tcPrChildRank(element.Name) > rank {
			insertAt = element.Start
			break
		}
	}
	if insertAt < 0 {
		return append(inner, []byte(borderXML)...)
	}
	result := make([]byte, 0, len(inner)+len(borderXML))
	result = append(result, inner[:insertAt]...)
	result = append(result, []byte(borderXML)...)
	result = append(result, inner[insertAt:]...)
	return result
}
func removeSingleXMLElement(content []byte, openTag []byte, closeTag []byte) []byte {
	start := bytes.Index(content, openTag)
	if start == -1 {
		return content
	}
	tagEndRel := bytes.Index(content[start:], []byte(">"))
	if tagEndRel == -1 {
		return content
	}
	tagEnd := start + tagEndRel
	var removeEnd int
	if content[tagEnd-1] == '/' {
		removeEnd = tagEnd + 1
	} else {
		closeRel := bytes.Index(content[start:], closeTag)
		if closeRel == -1 {
			return content
		}
		removeEnd = start + closeRel + len(closeTag)
	}
	result := make([]byte, 0, len(content)-(removeEnd-start))
	result = append(result, content[:start]...)
	result = append(result, content[removeEnd:]...)
	return result
}
