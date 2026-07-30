package editor

import (
	"regexp"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Master/layout XML scanning patterns. Fixed literals, compiled once at init.
var (
	shapeNamePattern          = regexp.MustCompile(`(?s)<p:nvSpPr>.*?<p:cNvPr[^>]*name="([^"]+)"`)
	placeholderElementPattern = regexp.MustCompile(`(?s)<p:ph[^>]*>`)
)

// GetMasterShapes returns the shapes in a slide master.
func (e *PresentationEditor) GetMasterShapes(masterPart string) []string {
	xmlData, ok := e.parts.Get(masterPart)
	if !ok {
		return nil
	}
	return parseShapesFromMasterLayoutXML(xmlData)
}

// GetMasterPlaceholders returns the placeholders in a slide master.
func (e *PresentationEditor) GetMasterPlaceholders(masterPart string) []common.PlaceholderInfo {
	xmlData, ok := e.parts.Get(masterPart)
	if !ok {
		return nil
	}
	return parsePlaceholdersFromMasterLayoutXML(xmlData)
}

// GetLayoutShapes returns the shapes in a slide layout.
func (e *PresentationEditor) GetLayoutShapes(layoutPart string) []string {
	xmlData, ok := e.parts.Get(layoutPart)
	if !ok {
		return nil
	}
	return parseShapesFromMasterLayoutXML(xmlData)
}

// GetLayoutPlaceholders returns the placeholders in a slide layout.
func (e *PresentationEditor) GetLayoutPlaceholders(layoutPart string) []common.PlaceholderInfo {
	xmlData, ok := e.parts.Get(layoutPart)
	if !ok {
		return nil
	}
	return parsePlaceholdersFromMasterLayoutXML(xmlData)
}

// parseShapesFromMasterLayoutXML extracts shape names from master/layout XML.
func parseShapesFromMasterLayoutXML(content []byte) []string {
	var shapes []string
	matches := shapeNamePattern.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) > 1 {
			shapes = append(shapes, match[1])
		}
	}
	return shapes
}

// parsePlaceholdersFromMasterLayoutXML extracts placeholder info from master/layout XML.
func parsePlaceholdersFromMasterLayoutXML(content []byte) []common.PlaceholderInfo {
	var placeholders []common.PlaceholderInfo
	matches := placeholderElementPattern.FindAllString(string(content), -1)
	for _, match := range matches {
		placeholders = append(placeholders, common.PlaceholderInfo{
			Type:  parsePlaceholderAttrString(match, phTypeAttrPattern),
			Index: parsePlaceholderAttrIndex(match),
			Name:  parsePlaceholderAttrString(match, phNameAttrPattern),
		})
	}
	return placeholders
}
