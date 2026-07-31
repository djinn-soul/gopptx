package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

func parseGroupChildShapes(groupXML []byte) ([]parsedShape, error) {
	nodes, err := editorshape.ExtractGroupChildShapeNodes(
		groupXML,
		groupShapeTag,
		shapeTypePicture,
	)
	if err != nil {
		return nil, fmt.Errorf("extract group children: %w", err)
	}

	children := make([]parsedShape, 0, len(nodes))
	for _, node := range nodes {
		tag, tagErr := rootElementLocalName(node)
		if tagErr != nil {
			return nil, fmt.Errorf("read group child element: %w", tagErr)
		}
		child, parseErr := buildParsedShapeFromRange(
			node,
			0,
			int64(len(node)),
			tag,
		)
		if parseErr != nil {
			return nil, fmt.Errorf("parse group child %s: %w", tag, parseErr)
		}
		children = append(children, child)
	}
	return children, nil
}

func rootElementLocalName(content []byte) (string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			return "", errors.New("shape element is missing")
		}
		if err != nil {
			return "", err
		}
		if start, ok := token.(xml.StartElement); ok {
			return start.Name.Local, nil
		}
	}
}
