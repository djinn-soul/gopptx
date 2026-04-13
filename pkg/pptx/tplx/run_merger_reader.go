package tplx

import (
	"bytes"
	"encoding/xml"
	"strings"
)

// readRun reads a single <a:r>…</a:r> block, extracting rPr and text.
func readRun(dec *xml.Decoder) (runData, error) {
	var rd runData
	depth := 1
	for depth > 0 {
		tok, err := dec.Token()
		if err != nil {
			return rd, err //nolint:wrapcheck // Preserve tokenization errors while reading run content.
		}
		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			switch local := t.Name.Local; local {
			case "rPr":
				var buf bytes.Buffer
				buf.Write(reencStart(t))
				if err2 := collectRawDepth(dec, &buf); err2 != nil {
					return rd, err2
				}
				depth--
				rd.rprBytes = buf.Bytes()
			case "t":
				text, err2 := readCharData(dec)
				if err2 != nil {
					return rd, err2
				}
				depth--
				rd.text = text
			default:
				var buf bytes.Buffer
				buf.Write(reencStart(t))
				if err2 := collectRawDepth(dec, &buf); err2 != nil {
					return rd, err2
				}
				depth--
				rd.extras = append(rd.extras, buf.Bytes())
			}
		case xml.EndElement:
			depth--
		}
	}
	return rd, nil
}

// readCharData reads text content until the matching end element.
func readCharData(dec *xml.Decoder) (string, error) {
	var sb strings.Builder
	for {
		tok, err := dec.Token()
		if err != nil {
			return sb.String(), err //nolint:wrapcheck // Preserve decoder errors while reading character data.
		}
		switch t := tok.(type) {
		case xml.CharData:
			sb.Write(t)
		case xml.EndElement:
			return sb.String(), nil
		}
	}
}

// collectRawDepth reads elements until depth returns to 0, appending to buf.
func collectRawDepth(dec *xml.Decoder, buf *bytes.Buffer) error {
	depth := 1
	for depth > 0 {
		tok, err := dec.Token()
		if err != nil {
			return err //nolint:wrapcheck // Preserve raw depth parsing failures for caller handling.
		}
		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			buf.Write(reencStart(t))
		case xml.EndElement:
			depth--
			if depth >= 0 {
				buf.Write(reencEnd(t))
			}
		case xml.CharData:
			writeEscaped(buf, string(t))
		}
	}
	return nil
}
