package opc

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// The other half of the escape hatch: a caller who can read and write parts
// still has to produce well-formed markup for them. gopptx's own XML builder
// lives in internal/pptxxml, which the compiler makes unreachable from outside
// the module, so this is a small writer with the same shape.

// XMLWriter builds an XML document element by element, tracking open elements
// so a mismatched close is caught rather than written.
type XMLWriter struct {
	b    strings.Builder
	open []string
	err  error
}

// NewXMLWriter starts a document, with the standalone declaration OOXML parts
// carry.
func NewXMLWriter() *XMLWriter {
	w := &XMLWriter{}
	w.b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	return w
}

// NewXMLFragment starts a document with no declaration, for markup spliced into
// an existing part.
func NewXMLFragment() *XMLWriter {
	return &XMLWriter{}
}

// StartElement opens an element. Attributes are name/value pairs; an odd count
// is a programming error and is reported by Finish.
func (w *XMLWriter) StartElement(name string, attrs ...string) *XMLWriter {
	if w.err != nil {
		return w
	}
	if len(attrs)%2 != 0 {
		w.err = errOddAttributes
		return w
	}
	w.b.WriteString("<" + name)
	w.writeAttrs(attrs)
	w.b.WriteString(">")
	w.open = append(w.open, name)
	return w
}

// EmptyElement writes a self-closing element.
func (w *XMLWriter) EmptyElement(name string, attrs ...string) *XMLWriter {
	if w.err != nil {
		return w
	}
	if len(attrs)%2 != 0 {
		w.err = errOddAttributes
		return w
	}
	w.b.WriteString("<" + name)
	w.writeAttrs(attrs)
	w.b.WriteString("/>")
	return w
}

// EndElement closes the innermost open element.
func (w *XMLWriter) EndElement() *XMLWriter {
	if w.err != nil {
		return w
	}
	if len(w.open) == 0 {
		w.err = errNoOpenElement
		return w
	}
	name := w.open[len(w.open)-1]
	w.open = w.open[:len(w.open)-1]
	w.b.WriteString("</" + name + ">")
	return w
}

// Text writes escaped character data.
func (w *XMLWriter) Text(text string) *XMLWriter {
	if w.err != nil {
		return w
	}
	w.b.WriteString(EscapeXML(text))
	return w
}

// Int writes a number as text.
func (w *XMLWriter) Int(value int64) *XMLWriter {
	return w.Text(strconv.FormatInt(value, 10))
}

// Raw writes markup through untouched, for the cases a typed builder cannot
// express. The caller owns its well-formedness.
func (w *XMLWriter) Raw(markup string) *XMLWriter {
	if w.err != nil {
		return w
	}
	w.b.WriteString(markup)
	return w
}

// Element writes an element with text content in one call.
func (w *XMLWriter) Element(name, text string, attrs ...string) *XMLWriter {
	return w.StartElement(name, attrs...).Text(text).EndElement()
}

// Finish returns the document, or an error if an element was left open or an
// attribute list was malformed.
func (w *XMLWriter) Finish() (string, error) {
	if w.err != nil {
		return "", w.err
	}
	if len(w.open) > 0 {
		return "", &UnclosedElementError{Name: w.open[len(w.open)-1]}
	}
	return w.b.String(), nil
}

// String returns the document so far, ignoring any error. Finish is the
// checked form.
func (w *XMLWriter) String() string {
	return w.b.String()
}

func (w *XMLWriter) writeAttrs(attrs []string) {
	for i := 0; i < len(attrs); i += 2 {
		w.b.WriteString(" " + attrs[i] + `="` + EscapeXML(attrs[i+1]) + `"`)
	}
}

// EscapeXML escapes text for use in XML content or an attribute value.
func EscapeXML(value string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(value))
	return b.String()
}

// UnclosedElementError reports markup that ended with an element still open.
type UnclosedElementError struct {
	Name string
}

func (e *UnclosedElementError) Error() string {
	return "unclosed XML element: " + e.Name
}

type xmlWriterError string

func (e xmlWriterError) Error() string { return string(e) }

const (
	errOddAttributes xmlWriterError = "attributes must be name/value pairs"
	errNoOpenElement xmlWriterError = "no open element to close"
)
