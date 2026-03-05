package editor

import (
	"bytes"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func (e *PresentationEditor) writeRelationships(path string, rels []common.EditorRelationship) error {
	e.parts.Set(path, []byte(renderRelationshipsXML(rels)))
	return nil
}

func (e *PresentationEditor) addContentTypeOverride(partName, contentType string) {
	ctPath := "[Content_Types].xml"
	data, ok := e.parts.Get(ctPath)
	if !ok {
		return
	}

	partNameRooted := "/" + partName
	if bytes.Contains(data, []byte(`PartName="`+partNameRooted+`"`)) {
		return
	}

	override := fmt.Sprintf(`<Override PartName="%s" ContentType="%s"/>`, partNameRooted, contentType)
	replaced := bytes.Replace(data, []byte("</Types>"), []byte(override+"</Types>"), 1)
	e.parts.Set(ctPath, replaced)
}
