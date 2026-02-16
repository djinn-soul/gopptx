package editor

import (
	"errors"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// SearchShapes scans all slides and returns shapes matching the query.
func (e *PresentationEditor) SearchShapes(query common.ShapeSearchQuery) ([]common.ShapeSearchResult, error) {
	if e == nil {
		return nil, errors.New("editor cannot be nil")
	}

	results := make([]common.ShapeSearchResult, 0)
	for slideIndex := range e.slides {
		shapes, err := e.GetShapes(slideIndex)
		if err != nil {
			return nil, err
		}
		for _, shape := range shapes {
			if !shapeMatchesQuery(shape, query) {
				continue
			}
			results = append(results, common.ShapeSearchResult{
				SlideIndex: slideIndex,
				Shape:      shape,
			})
		}
	}
	return results, nil
}

func shapeMatchesQuery(shape common.Shape, query common.ShapeSearchQuery) bool {
	name := shape.Name
	typ := shape.Type
	text := shape.Text
	qName := query.NameContains
	qType := query.TypeEquals
	qText := query.TextContains

	if !query.CaseSensitive {
		name = strings.ToLower(name)
		typ = strings.ToLower(typ)
		text = strings.ToLower(text)
		qName = strings.ToLower(qName)
		qType = strings.ToLower(qType)
		qText = strings.ToLower(qText)
	}

	if qName != "" && !strings.Contains(name, qName) {
		return false
	}
	if qType != "" && typ != qType {
		return false
	}
	if qText != "" && !strings.Contains(text, qText) {
		return false
	}
	return true
}
