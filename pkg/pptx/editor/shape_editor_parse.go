package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// parsedShape represents a shape found in the slide XML.
// It contains the parsed properties and the byte range of the shape node.
type parsedShape struct {
	ID          int
	Name        string
	Type        string // "sp" or "pic"
	Text        string
	Runs        []common.TextRun
	TextFrame   *common.TextFrame
	Paragraph   *common.Paragraph
	Fill        *common.ShapeFill
	Line        *common.ShapeLine
	Shadow      *common.ShapeShadow
	Glow        *common.ShapeGlow
	Blur        *common.ShapeBlur
	SoftEdge    *common.ShapeSoftEdge
	Reflection  *common.ShapeReflection
	ClickAction *common.Hyperlink
	HoverAction *common.Hyperlink
	X, Y        int
	W, H        int
	PhIndex     int    // Placeholder index, -1 if not a placeholder
	PhType      string // Placeholder type (e.g. "title", "body")
	Adjustments []common.ShapeAdjustment
	Start       int64 // Byte offset of the start of the node
	End         int64 // Byte offset of the end of the node
	IsGroup     bool
}

func (p parsedShape) ToShape() shapes.Shape {
	return shapes.Shape{
		Type: p.Type,
		X:    styling.Emu(int64(p.X)),
		Y:    styling.Emu(int64(p.Y)),
		CX:   styling.Emu(int64(p.W)),
		CY:   styling.Emu(int64(p.H)),
		Text: p.Text,
		Name: p.Name,
		Adjustments: func() []shapes.ShapeAdjustment {
			out := make([]shapes.ShapeAdjustment, 0, len(p.Adjustments))
			for _, adj := range p.Adjustments {
				out = append(out, shapes.ShapeAdjustment{
					Name:    adj.Name,
					Formula: adj.Formula,
				})
			}
			return out
		}(),
	}
}

// parseSlideShapes scans the slide XML for shape nodes and extracts their properties and byte ranges.
func parseSlideShapes(content []byte) ([]parsedShape, error) {
	return scanShapesWithOffsets(content, false)
}

func scanShapesWithOffsets(content []byte, skipProperties bool) ([]parsedShape, error) {
	reader := bytes.NewReader(content)
	decoder := xml.NewDecoder(reader)
	var shapes []parsedShape

	// We need to track depth to know when we exit a shape
	// <p:sp> ... </p:sp>

	for {
		// handle offset before reading token
		startOffset := decoder.InputOffset()
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		se, ok := token.(xml.StartElement)
		if !ok {
			continue
		}

		if se.Name.Local == "sp" || se.Name.Local == shapeTypePicture || se.Name.Local == shapeTypeGraphicFrame ||
			se.Name.Local == groupShapeTag ||
			se.Name.Local == "cxnSp" {
			// Found a shape start.
			// We need to capture the exact bytes from `startOffset` until the end element.
			// The `decoder.InputOffset()` gives the start of the token *buffer* usually, but for Bytes.Reader it's precise enough usually
			// IF we haven't read ahead.
			// Actually `InputOffset()` returns the number of bytes read *so far*.
			// So `startOffset` is the end of the *previous* token.

			// Let's extract this node.
			shape, endOffset, extractErr := extractShapeNode(
				content,
				startOffset,
				decoder,
				se.Name.Local,
				skipProperties,
			)
			if extractErr != nil {
				return nil, extractErr
			}
			shapes = append(shapes, shape)

			// Reset/Sync decoder is tricky if we consumed bytes manually.
			// Helper `extractShapeNode` should advance the decoder one token at a time until end.
			_ = endOffset
		}
	}

	return shapes, nil
}

// extractShapeNode consumes tokens until the matching end element is found.
// It also parses the content within that range to populate parsedShape.
func extractShapeNode(
	fullContent []byte,
	startOffset int64,
	decoder *xml.Decoder,
	stopTag string,
	skipProperties bool,
) (parsedShape, int64, error) {
	depth := 1

	// To parse attributes, we can try to unmarshal the captured byte range later.
	// For now, let's just find the end offset.

	for {
		token, err := decoder.Token()
		if err != nil {
			return parsedShape{}, 0, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			depth = adjustShapeDepthForStart(depth, t.Name.Local, stopTag)
		case xml.EndElement:
			nextDepth, done := adjustShapeDepthForEnd(depth, t.Name.Local, stopTag)
			depth = nextDepth
			if done {
				endOffset := decoder.InputOffset()
				var pShape parsedShape
				var parseErr error

				if skipProperties {
					// Optimization: Just record boundaries/type
					pShape = parsedShape{
						Start: startOffset,
						End:   endOffset,
						Type:  stopTag,
					}
				} else {
					pShape, parseErr = buildParsedShapeFromRange(fullContent, startOffset, endOffset, stopTag)
					if parseErr != nil {
						return parsedShape{}, 0, parseErr
					}
				}
				return pShape, endOffset, nil
			}
		}
	}
}

func adjustShapeDepthForStart(currentDepth int, tokenName, stopTag string) int {
	if tokenName == stopTag {
		return currentDepth + 1
	}
	return currentDepth
}

func adjustShapeDepthForEnd(currentDepth int, tokenName, stopTag string) (int, bool) {
	if tokenName != stopTag {
		return currentDepth, false
	}
	nextDepth := currentDepth - 1
	return nextDepth, nextDepth == 0
}

func buildParsedShapeFromRange(
	fullContent []byte,
	startOffset, endOffset int64,
	stopTag string,
) (parsedShape, error) {
	if startOffset < 0 || startOffset >= endOffset || endOffset > int64(len(fullContent)) {
		return parsedShape{}, fmt.Errorf(
			"invalid shape offsets: start=%d end=%d size=%d",
			startOffset,
			endOffset,
			len(fullContent),
		)
	}

	shapeXML := fullContent[startOffset:endOffset]
	pShape, parseErr := parseShapeProperties(shapeXML)
	if parseErr != nil {
		return parsedShape{}, parseErr
	}
	if stopTag == shapeTypePicture || stopTag == shapeTypeGraphicFrame || stopTag == groupShapeTag {
		pShape.Type = stopTag
	} else if pShape.Type == "" {
		pShape.Type = stopTag
	}
	pShape.Start = startOffset
	pShape.End = endOffset
	pShape.IsGroup = stopTag == "grpSp"
	return pShape, nil
}

func parseShapeProperties(content []byte) (parsedShape, error) {
	props, err := editorshape.ParseShapeProperties(content)
	if err != nil {
		return parsedShape{}, err
	}
	return parsedShape{
		ID:          props.ID,
		Name:        props.Name,
		Type:        props.Type,
		Text:        props.Text,
		Runs:        props.Runs,
		Paragraph:   props.Paragraph,
		Fill:        props.Fill,
		Line:        props.Line,
		Shadow:      props.Shadow,
		Glow:        props.Glow,
		Blur:        props.Blur,
		SoftEdge:    props.SoftEdge,
		Reflection:  props.Reflection,
		X:           props.X,
		Y:           props.Y,
		W:           props.W,
		H:           props.H,
		PhIndex:     props.PhIndex,
		PhType:      props.PhType,
		Adjustments: props.Adjustments,
	}, nil
}
