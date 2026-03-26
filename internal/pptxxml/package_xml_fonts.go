package pptxxml

import (
	"strconv"
	"strings"
)

// CustomProperties renders docProps/custom.xml.
func CustomProperties(markAsFinal bool) string {
	if !markAsFinal {
		return ""
	}
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/custom-properties" ` +
		`xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">
<property fmtid="{D5CDD505-2E9C-101B-9397-08002B2CF9AE}" pid="2" name="_MarkAsFinal">
<vt:bool>true</vt:bool>
</property>
</Properties>`
}

// EmbeddedFontRef describes a font for packing into the presentation XML.
type EmbeddedFontRef struct {
	Typeface    string
	Style       string // e.g. "regular", "bold", "italic", "boldItalic"
	Charset     uint8
	Panose      string
	PitchFamily uint8
	RelID       string
}

// EmbeddedFontsXML renders the <p:embeddedFontLst> block.
func EmbeddedFontsXML(fonts []EmbeddedFontRef) string {
	if len(fonts) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n<p:embeddedFontLst>")

	grouped := make(map[string][]EmbeddedFontRef)
	var order []string
	for _, f := range fonts {
		if _, exists := grouped[f.Typeface]; !exists {
			order = append(order, f.Typeface)
		}
		grouped[f.Typeface] = append(grouped[f.Typeface], f)
	}

	for _, typeface := range order {
		variants := grouped[typeface]
		if len(variants) == 0 {
			continue
		}
		first := variants[0]
		b.WriteString("\n<p:embeddedFont>")
		b.WriteString("\n  <p:font typeface=\"")
		b.WriteString(Escape(first.Typeface))
		b.WriteString("\" pitchFamily=\"")
		b.WriteString(strconv.Itoa(int(first.PitchFamily)))
		b.WriteString("\" charset=\"")
		b.WriteString(strconv.Itoa(int(first.Charset)))
		b.WriteString("\"")
		if first.Panose != "" {
			b.WriteString(" panose=\"")
			b.WriteString(Escape(first.Panose))
			b.WriteString("\"")
		}
		b.WriteString("/>")

		for _, v := range variants {
			b.WriteString("\n  <p:")
			b.WriteString(v.Style)
			b.WriteString(" r:id=\"")
			b.WriteString(v.RelID)
			b.WriteString("\"/>")
		}
		b.WriteString("\n</p:embeddedFont>")
	}

	b.WriteString("\n</p:embeddedFontLst>")
	return b.String()
}
