package editor

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

func handleAddVba(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	dataBase64, ok := v.RequireString(p, "data")
	if !ok {
		return nil, v.Error()
	}

	data, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, "data must be a valid base64 string")
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
