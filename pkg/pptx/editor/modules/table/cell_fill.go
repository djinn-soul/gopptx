package table

import (
	"bytes"
	"errors"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// applyCellFill sets or clears the cell's solid fill. The fill goes last in
// <a:tcPr>, after any border lines, per CT_TableCellProperties.
func applyCellFill(cellContent []byte, backgroundColor string) ([]byte, error) {
	color := strings.TrimSpace(backgroundColor)
	if color == "" {
		return cellContent, nil
	}

	fillXML := ""
	if !strings.EqualFold(color, "none") {
		fillXML = `<a:solidFill><a:srgbClr val="` +
			common.XMLEscape(strings.TrimPrefix(color, "#")) + `"/></a:solidFill>`
	}
	return setCellFillXML(cellContent, fillXML)
}

func setCellFillXML(cellContent []byte, fillXML string) ([]byte, error) {
	tcPrStart := bytes.Index(cellContent, []byte("<a:tcPr"))
	if tcPrStart == -1 {
		closeTC := bytes.Index(cellContent, []byte("</a:tc>"))
		if closeTC == -1 {
			return nil, errors.New("invalid cell xml: missing </a:tc>")
		}
		return spliceBytes(cellContent, closeTC, closeTC, "<a:tcPr>"+fillXML+"</a:tcPr>"), nil
	}

	tagEndRel := bytes.Index(cellContent[tcPrStart:], []byte(">"))
	if tagEndRel == -1 {
		return nil, errors.New("invalid tcPr xml: missing >")
	}
	tcPrTagEnd := tcPrStart + tagEndRel

	if cellContent[tcPrTagEnd-1] == '/' {
		openTag := string(cellContent[tcPrStart:tcPrTagEnd-1]) + ">"
		return spliceBytes(
			cellContent, tcPrStart, tcPrTagEnd+1, openTag+fillXML+"</a:tcPr>",
		), nil
	}

	closeTcPr := bytes.Index(cellContent[tcPrStart:], []byte("</a:tcPr>"))
	if closeTcPr == -1 {
		return nil, errors.New("invalid tcPr xml: missing </a:tcPr>")
	}
	innerStart := tcPrTagEnd + 1
	innerEnd := tcPrStart + closeTcPr

	inner := removeCellFillElements(cellContent[innerStart:innerEnd])
	return spliceBytes(cellContent, innerStart, innerEnd, string(inner)+fillXML), nil
}

func removeCellFillElements(inner []byte) []byte {
	out := inner
	for _, tag := range fillChildren {
		out = removeSingleXMLElement(out, []byte("<a:"+tag), []byte("</a:"+tag+">"))
	}
	return out
}

func spliceBytes(content []byte, start, end int, replacement string) []byte {
	out := make([]byte, 0, len(content)-(end-start)+len(replacement))
	out = append(out, content[:start]...)
	out = append(out, []byte(replacement)...)
	out = append(out, content[end:]...)
	return out
}
