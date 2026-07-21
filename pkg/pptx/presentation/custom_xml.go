package presentation

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

var customXMLNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9._-]*$`)

//nolint:gochecknoglobals // reusable immutable replacer
var customXMLEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	"\"", "&quot;",
	"'", "&apos;",
)

func writeCustomXMLParts(pw *pptxxml.PackageWriter, customXML []common.CustomXMLPart) error {
	for i, part := range customXML {
		var itemXML string
		var err error
		if part.RootElement != "" {
			itemXML, err = generateCustomXMLItem(part)
			if err != nil {
				return fmt.Errorf("custom XML part %d: %w", i+1, err)
			}
		} else {
			// Legacy passthrough: ensure part content is well-formed XML.
			if err := xml.Unmarshal([]byte(part.Content), new(any)); err != nil {
				return fmt.Errorf("custom XML part %d contains invalid XML: %w", i+1, err)
			}
			itemXML = part.Content
		}
		path := fmt.Sprintf("customXml/item%d.xml", i+1)
		pw.AddPart(path, itemXML)

		// Generate itemProps.
		itemID, err := common.NewGUID()
		if err != nil {
			return fmt.Errorf("generate custom XML itemID for part %d: %w", i+1, err)
		}
		schemaRefs := "<ds:schemaRefs></ds:schemaRefs>"
		if part.Namespace != "" {
			schemaRefs = fmt.Sprintf(
				`<ds:schemaRefs><ds:schemaRef ds:uri="%s"/></ds:schemaRefs>`,
				escapeCustomXML(part.Namespace),
			)
		}

		propsContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ds:datastoreItem ds:itemID="%s" xmlns:ds="http://schemas.openxmlformats.org/officeDocument/2006/customXml">
%s
</ds:datastoreItem>`, itemID, schemaRefs)
		propsPath := fmt.Sprintf("customXml/itemProps%d.xml", i+1)
		pw.AddPart(propsPath, propsContent)

		// Create itemN.xml.rels to link the item to its properties.
		itemRelContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXmlProps" Target="itemProps%d.xml"/>
</Relationships>`, i+1)
		pw.AddPart(fmt.Sprintf("customXml/_rels/item%d.xml.rels", i+1), itemRelContent)
	}
	return nil
}

func generateCustomXMLItem(part common.CustomXMLPart) (string, error) {
	if !customXMLNamePattern.MatchString(part.RootElement) {
		return "", fmt.Errorf("invalid root element name %q", part.RootElement)
	}

	nsAttr := ""
	if part.Namespace != "" {
		nsAttr = fmt.Sprintf(` xmlns="%s"`, escapeCustomXML(part.Namespace))
	}

	inner := part.Content
	if inner == "" && len(part.Properties) > 0 {
		var propsSb55 strings.Builder
		for j, kv := range part.Properties {
			if !customXMLNamePattern.MatchString(kv.Key) {
				return "", fmt.Errorf("invalid property element name %q", kv.Key)
			}
			if j > 0 {
				propsSb55.WriteString("\n  ")
			}
			fmt.Fprintf(&propsSb55, "<%s>%s</%s>", kv.Key, escapeCustomXML(kv.Value), kv.Key)
		}
		inner = propsSb55.String()
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<%s%s>
  %s
</%s>`, part.RootElement, nsAttr, inner, part.RootElement), nil
}

func escapeCustomXML(value string) string {
	return customXMLEscaper.Replace(value)
}
