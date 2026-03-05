package shape

import (
	"bytes"
	"regexp"
)

func ReplaceShapeTextBody(xmlData []byte, txBody []byte) []byte {
	txBodyOpenRe := regexp.MustCompile(`(?s)<p:txBody\b[^>]*>`)
	txBodyCloseTag := []byte("</p:txBody>")

	loc := txBodyOpenRe.FindIndex(xmlData)
	var startIdx, endIdx int
	found := false

	if loc != nil {
		startIdx = loc[0]
		closeIdx := bytes.Index(xmlData[loc[1]:], txBodyCloseTag)
		if closeIdx != -1 {
			endIdx = loc[1] + closeIdx
			found = true
		}
	}

	if !found {
		spPrEndRe := regexp.MustCompile(`(?s)</p:spPr>`)
		spPrLoc := spPrEndRe.FindIndex(xmlData)
		if spPrLoc != nil {
			insertIdx := spPrLoc[1]
			res := make([]byte, 0, len(xmlData)+len(txBody))
			res = append(res, xmlData[:insertIdx]...)
			res = append(res, txBody...)
			res = append(res, xmlData[insertIdx:]...)
			return res
		}
		return xmlData
	}

	res := make([]byte, 0, len(xmlData)-(endIdx+len(txBodyCloseTag)-startIdx)+len(txBody))
	res = append(res, xmlData[:startIdx]...)
	res = append(res, txBody...)
	res = append(res, xmlData[endIdx+len(txBodyCloseTag):]...)
	return res
}
