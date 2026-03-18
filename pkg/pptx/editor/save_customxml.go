package editor

import (
	"fmt"
	"unsafe"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

func (e *PresentationEditor) processCustomXMLParts(out map[string][]byte) ([]string, error) {
	customXMLPropsPaths := make([]string, 0, len(e.metadata.CustomXML))
	for i, cXML := range e.metadata.CustomXML {
		itemPath, propsPath, err := e.writeCustomXMLPart(out, i, cXML)
		if err != nil {
			return nil, err
		}
		customXMLPropsPaths = append(customXMLPropsPaths, propsPath)
		e.ensureCustomXMLRelationship(itemPath)
		e.writeCustomXMLPropsRelationship(out, i+1)
	}
	return customXMLPropsPaths, nil
}

func (e *PresentationEditor) writeCustomXMLPart(
	out map[string][]byte,
	index int,
	part common.CustomXMLPart,
) (string, string, error) {
	itemContent, err := buildCustomXMLItemContent(part, index+1)
	if err != nil {
		return "", "", err
	}
	itemPath := fmt.Sprintf("customXml/item%d.xml", index+1)
	propsPath := fmt.Sprintf("customXml/itemProps%d.xml", index+1)
	out[itemPath] = []byte(itemContent)
	propsContent, err := buildCustomXMLPropsContent(part)
	if err != nil {
		return "", "", err
	}
	out[propsPath] = []byte(propsContent)
	return itemPath, propsPath, nil
}

func buildCustomXMLItemContent(part common.CustomXMLPart, index int) (string, error) {
	if part.RootElement == "" {
		return part.Content, nil
	}
	itemContent, err := editorslide.GenerateCustomXMLItem(part)
	if err != nil {
		return "", fmt.Errorf("custom XML part %d: %w", index, err)
	}
	return itemContent, nil
}

func buildCustomXMLPropsContent(part common.CustomXMLPart) (string, error) {
	schemaRefs := "<ds:schemaRefs></ds:schemaRefs>"
	if part.Namespace != "" {
		schemaRefs = fmt.Sprintf(
			`<ds:schemaRefs><ds:schemaRef ds:uri="%s"/></ds:schemaRefs>`,
			editorslide.EscapeCustomXML(part.Namespace),
		)
	}
	itemID := part.ItemID
	if itemID == "" {
		guid, err := common.NewGUID()
		if err != nil {
			return "", err
		}
		itemID = guid
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ds:datastoreItem ds:itemID="%s" xmlns:ds="http://schemas.openxmlformats.org/officeDocument/2006/customXml">
%s
</ds:datastoreItem>`, itemID, schemaRefs), nil
}

func (e *PresentationEditor) ensureCustomXMLRelationship(itemPath string) {
	itemTarget := "../" + itemPath
	for _, r := range e.nonSlideRels {
		if r.Type == common.RelTypeCustomXML && r.Target == itemTarget {
			return
		}
	}
	e.nonSlideRels = append(e.nonSlideRels, common.EditorRelationship{
		ID:     fmt.Sprintf("rId%d", e.nextRelIDNum),
		Type:   common.RelTypeCustomXML,
		Target: itemTarget,
	})
	e.nextRelIDNum++
}

func (e *PresentationEditor) writeCustomXMLPropsRelationship(out map[string][]byte, index int) {
	itemRelContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXmlProps" Target="itemProps%d.xml"/>
</Relationships>`, index)
	out[fmt.Sprintf("customXml/_rels/item%d.xml.rels", index)] = []byte(itemRelContent)
}

func (e *PresentationEditor) filterRootCustomXMLRelationships(out map[string][]byte) {
	packageRelsData, ok := e.parts.Get("_rels/.rels")
	if !ok {
		return
	}

	if len(packageRelsData) > 0 {
		ptr := uintptr(unsafe.Pointer(&packageRelsData[0])) //nolint:gosec // staleness token only, never dereferenced
		if e.packageRelsPtr == ptr {
			if e.packageRelsNeedsFilter {
				out["_rels/.rels"] = e.packageRelsFilteredXML
			}
			return
		}
		e.packageRelsPtr = ptr
	} else {
		e.packageRelsPtr = 0
		e.packageRelsNeedsFilter = false
		e.packageRelsFilteredXML = nil
		return
	}

	rels, err := parseRelationshipsXML(packageRelsData)
	if err != nil {
		return
	}
	filtered := make([]common.EditorRelationship, 0, len(rels))
	changed := false
	for _, r := range rels {
		if r.Type == common.RelTypeCustomXML || r.Type == common.RelTypeCustomXMLProps {
			changed = true
			continue
		}
		filtered = append(filtered, r)
	}
	if changed {
		filteredXML := []byte(renderRelationshipsXML(filtered))
		e.packageRelsNeedsFilter = true
		e.packageRelsFilteredXML = filteredXML
		out["_rels/.rels"] = filteredXML
		return
	}
	e.packageRelsNeedsFilter = false
	e.packageRelsFilteredXML = nil
}
