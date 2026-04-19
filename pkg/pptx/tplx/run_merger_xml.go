package tplx

import (
	"bytes"
	"encoding/xml"
	"strings"
)

// isDrawingMLLocal reports whether s is a DrawingML element with the given local name.
func isDrawingMLLocal(s xml.StartElement, local string) bool {
	if s.Name.Local != local {
		return false
	}
	return s.Name.Space == "a" ||
		s.Name.Space == "http://schemas.openxmlformats.org/drawingml/2006/main" ||
		s.Name.Space == ""
}

// reencStart re-encodes a StartElement, preserving namespace prefixes.
func reencStart(s xml.StartElement) []byte {
	var b bytes.Buffer
	b.WriteByte('<')
	writeQName(&b, s.Name)
	for _, a := range s.Attr {
		b.WriteByte(' ')
		writeQName(&b, a.Name)
		b.WriteString(`="`)
		if err := xml.EscapeText(&b, []byte(a.Value)); err != nil {
			b.WriteString(a.Value)
		}
		b.WriteByte('"')
	}
	b.WriteByte('>')
	return b.Bytes()
}

// reencEnd re-encodes an EndElement.
func reencEnd(e xml.EndElement) []byte {
	var b bytes.Buffer
	b.WriteString("</")
	writeQName(&b, e.Name)
	b.WriteByte('>')
	return b.Bytes()
}

func writeQName(b *bytes.Buffer, name xml.Name) {
	if name.Space != "" {
		if prefix := namespacePrefix(name.Space); prefix != "" {
			b.WriteString(prefix)
			b.WriteByte(':')
		}
	}
	b.WriteString(name.Local)
}

func namespacePrefix(space string) string {
	switch space {
	case "http://schemas.openxmlformats.org/presentationml/2006/main":
		return "p"
	case "http://schemas.openxmlformats.org/drawingml/2006/main":
		return "a"
	case "http://schemas.openxmlformats.org/officeDocument/2006/relationships":
		return "r"
	case "http://schemas.openxmlformats.org/markup-compatibility/2006":
		return "mc"
	case "http://schemas.microsoft.com/office/drawing/2010/main":
		return "a14"
	case "http://schemas.microsoft.com/office/drawing/2012/main":
		return "a15"
	case "http://schemas.openxmlformats.org/drawingml/2006/chart":
		return "c"
	case "http://schemas.openxmlformats.org/drawingml/2006/diagram":
		return "dgm"
	case "http://schemas.openxmlformats.org/drawingml/2006/picture":
		return "pic"
	case "http://schemas.microsoft.com/office/powerpoint/2010/main":
		return "p14"
	case "http://schemas.microsoft.com/office/powerpoint/2012/main":
		return "p15"
	case "http://schemas.openxmlformats.org/package/2006/relationships":
		return "rel"
	case "http://www.w3.org/XML/1998/namespace":
		return "xml"
	case "xmlns":
		return "xmlns"
	default:
		if strings.Contains(space, "://") {
			return ""
		}
		return space
	}
}

// writeEscaped writes XML-escaped character data.
func writeEscaped(out *bytes.Buffer, s string) {
	if err := xml.EscapeText(out, []byte(s)); err != nil {
		out.WriteString(s)
	}
}

// emitRun writes a complete <a:r>…</a:r> element to out.
func emitRun(out *bytes.Buffer, rd runData) {
	out.WriteString("<a:r>")
	if len(rd.rprBytes) > 0 {
		out.Write(rd.rprBytes)
	}
	for _, x := range rd.extras {
		out.Write(x)
	}
	out.WriteString("<a:t>")
	writeEscaped(out, rd.text)
	out.WriteString("</a:t>")
	out.WriteString("</a:r>")
}
