package slide

import (
	"bytes"
	"encoding/xml"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func ParseCustomXMLInventory(ps PartLookup, partKeys []string) []common.CustomXMLPart {
	var parts []common.CustomXMLPart

	var itemPaths []string
	for _, p := range partKeys {
		if strings.HasPrefix(p, "customXml/item") && strings.HasSuffix(p, ".xml") &&
			!strings.HasPrefix(p, "customXml/itemProps") {
			itemPaths = append(itemPaths, p)
		}
	}

	for _, p := range itemPaths {
		itemData, _ := ps.Get(p)
		propsPath := strings.Replace(p, "item", "itemProps", 1)

		part := common.CustomXMLPart{Content: string(itemData)}

		if propsData, hasProps := ps.Get(propsPath); hasProps {
			part.Namespace = parseCustomXMLNamespace(propsData)
			part.ItemID = parseCustomXMLItemID(propsData)
		}

		if part.Namespace != "" {
			structuredPart, ok := parseStructuredCustomXML(itemData, part.Namespace)
			if ok {
				structuredPart.ItemID = part.ItemID
				parts = append(parts, structuredPart)
				continue
			}
		}
		parts = append(parts, part)
	}

	return parts
}

func parseCustomXMLItemID(propsData []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(propsData))
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		if start.Name.Local == "datastoreItem" {
			for _, attr := range start.Attr {
				if attr.Name.Local == "itemID" {
					return attr.Value
				}
			}
		}
	}
	return ""
}

func parseCustomXMLNamespace(propsData []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(propsData))
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		if start.Name.Local == "schemaRef" {
			for _, attr := range start.Attr {
				if attr.Name.Local == "uri" {
					return attr.Value
				}
			}
		}
	}
	return ""
}

func parseStructuredCustomXML(itemData []byte, ns string) (common.CustomXMLPart, bool) {
	decoder := xml.NewDecoder(bytes.NewReader(itemData))

	var rootName string
	var props []common.CustomXMLKV

	for {
		t, err := decoder.Token()
		if err != nil {
			return common.CustomXMLPart{}, false
		}
		if start, ok := t.(xml.StartElement); ok {
			rootName = start.Name.Local
			break
		}
	}

	var currentKey string
	var currentValue strings.Builder

	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}
		switch t := t.(type) {
		case xml.StartElement:
			currentKey = t.Name.Local
			currentValue.Reset()
		case xml.CharData:
			if currentKey != "" {
				currentValue.Write(t)
			}
		case xml.EndElement:
			if t.Name.Local == rootName {
				return common.CustomXMLPart{
					RootElement: rootName,
					Namespace:   ns,
					Properties:  props,
				}, true
			}
			if currentKey == t.Name.Local {
				props = append(props, common.CustomXMLKV{Key: currentKey, Value: currentValue.String()})
				currentKey = ""
			}
		}
	}

	return common.CustomXMLPart{
		RootElement: rootName,
		Namespace:   ns,
		Properties:  props,
	}, true
}
