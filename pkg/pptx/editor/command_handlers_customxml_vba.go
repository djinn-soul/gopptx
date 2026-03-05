package editor

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

func handleAddCustomXML(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	request := editorcommand.ParseCustomXMLAddRequest(
		p,
		v.OptionalString,
		func(code, message string) {
			v.setCode(code)
			v.errors = append(v.errors, message)
		},
		ErrCodeMissingField,
		ErrCodeInvalidType,
	)

	if v.HasErrors() {
		return nil, v.Error()
	}

	part := common.CustomXMLPart{
		Content:     request.Content,
		RootElement: request.RootElement,
		Namespace:   request.Namespace,
	}

	keys := make([]string, 0, len(request.Properties))
	for k := range request.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		part.Properties = append(part.Properties, common.CustomXMLKV{Key: k, Value: request.Properties[k]})
	}

	e.metadata.CustomXML = append(e.metadata.CustomXML, part)
	return map[string]int{"index": len(e.metadata.CustomXML) - 1}, nil
}

func handleListCustomXML(e *PresentationEditor, _ json.RawMessage) (any, error) {
	type CustomXMLResp struct {
		Content     string            `json:"content,omitempty"`
		RootElement string            `json:"root_element,omitempty"`
		Namespace   string            `json:"namespace,omitempty"`
		Properties  map[string]string `json:"properties,omitempty"`
	}

	out := make([]CustomXMLResp, len(e.metadata.CustomXML))
	for i, part := range e.metadata.CustomXML {
		var props map[string]string
		if len(part.Properties) > 0 {
			props = make(map[string]string, len(part.Properties))
			for _, kv := range part.Properties {
				props[kv.Key] = kv.Value
			}
		}
		out[i] = CustomXMLResp{
			Content:     part.Content,
			RootElement: part.RootElement,
			Namespace:   part.Namespace,
			Properties:  props,
		}
	}
	return map[string]any{"custom_xml": out}, nil
}

func handleRemoveCustomXML(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	index, ok := v.RequireInt(p, "index")
	if !ok {
		return nil, v.Error()
	}
	if index < 0 || index >= len(e.metadata.CustomXML) {
		v.setCode(ErrCodeInvalidIndex)
		v.errors = append(
			v.errors,
			fmt.Sprintf("custom xml index %d out of range [0,%d)", index, len(e.metadata.CustomXML)),
		)
		return nil, v.Error()
	}
	e.metadata.CustomXML = append(e.metadata.CustomXML[:index], e.metadata.CustomXML[index+1:]...)
	return map[string]bool{"removed": true}, nil
}

func handleAddVba(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	data, ok, err := editorcommand.DecodeRequiredBase64Field(
		p,
		v.RequireString,
		"data",
		"data must be a valid base64 string",
	)
	if !ok {
		return nil, v.Error()
	}
	if err != nil {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, err.Error())
		return nil, v.Error()
	}

	var project *vba.VBAProject
	if e.metadata.VBA == nil {
		project = vba.New()
		e.metadata.VBA = project
	} else {
		project, ok = e.metadata.VBA.(*vba.VBAProject)
		if !ok {
			return nil, errors.New("invalid VBA metadata type")
		}
	}
	project.SetData(data)
	return map[string]bool{"added": true}, nil
}
