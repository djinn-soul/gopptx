package tplx

import (
	"bytes"
	"encoding/xml"
)

// mergeAdjacentRuns merges contiguous <a:r> runs within each <a:p> paragraph when
// those runs together form a Jinja token that PowerPoint split across multiple runs.
func mergeAdjacentRuns(xmlBytes []byte) []byte {
	dec := xml.NewDecoder(bytes.NewReader(xmlBytes))
	dec.Strict = false

	var out bytes.Buffer
	out.Grow(len(xmlBytes))

	if err := streamRewrite(dec, &out, ""); err != nil {
		return xmlBytes
	}
	return out.Bytes()
}

// streamRewrite recursively re-emits tokens, rewriting <a:p> blocks.
func streamRewrite(dec *xml.Decoder, out *bytes.Buffer, parentLocal string) error {
	for {
		tok, err := dec.Token()
		if err != nil {
			return handleStreamReadError(parentLocal, err)
		}
		switch t := tok.(type) {
		case xml.ProcInst:
			continue
		case xml.StartElement:
			if err = handleStreamStartElement(dec, out, t); err != nil {
				return err
			}
		case xml.EndElement:
			if shouldStop := handleStreamEndElement(out, parentLocal, t); shouldStop {
				return nil
			}
		case xml.CharData:
			handleStreamCharData(out, parentLocal, t)
		case xml.Comment:
			writeXMLComment(out, t)
		}
	}
}

func handleStreamReadError(parentLocal string, err error) error {
	if parentLocal == "" {
		return nil
	}
	return err //nolint:wrapcheck // Streaming decoder errors are returned unchanged for callers.
}

func handleStreamStartElement(dec *xml.Decoder, out *bytes.Buffer, start xml.StartElement) error {
	if isDrawingMLLocal(start, "p") {
		return writeParagraph(dec, out, start)
	}
	out.Write(reencStart(start))
	return streamRewrite(dec, out, start.Name.Local)
}

func handleStreamEndElement(out *bytes.Buffer, parentLocal string, end xml.EndElement) bool {
	out.Write(reencEnd(end))
	return end.Name.Local == parentLocal
}

func handleStreamCharData(out *bytes.Buffer, parentLocal string, charData xml.CharData) {
	if parentLocal == "" {
		return
	}
	writeEscaped(out, string(charData))
}

func writeXMLComment(out *bytes.Buffer, comment xml.Comment) {
	out.WriteString("<!--")
	out.Write(comment)
	out.WriteString("-->")
}
