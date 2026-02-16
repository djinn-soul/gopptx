package presentation

import (
	"encoding/xml"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

func writeCustomXMLParts(pw *pptxxml.PackageWriter, customXML []common.CustomXMLPart) error {
	for i, part := range customXML {
		// Ensure part content is well-formed XML.
		if err := xml.Unmarshal([]byte(part.Content), new(any)); err != nil {
			return fmt.Errorf("custom XML part %d contains invalid XML: %w", i+1, err)
		}
		path := fmt.Sprintf("customXml/item%d.xml", i+1)
		pw.AddPart(path, part.Content)

		// Generate itemProps.
		itemID, err := common.NewGUID()
		if err != nil {
			return fmt.Errorf("generate custom XML itemID for part %d: %w", i+1, err)
		}
		propsContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ds:datastoreItem ds:itemID="%s" xmlns:ds="http://schemas.openxmlformats.org/officeDocument/2006/customXml">
<ds:schemaRefs/>
</ds:datastoreItem>`, itemID)
		propsPath := fmt.Sprintf("customXml/itemProps%d.xml", i+1)
		pw.AddPart(propsPath, propsContent)
	}
	return nil
}
