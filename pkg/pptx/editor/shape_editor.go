package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// parsedShape represents a shape found in the slide XML.
// It contains the parsed properties and the byte range of the shape node.
type parsedShape struct {
	ID    int
	Name  string
	Type  string // "sp" or "pic"
	Text  string
	X, Y  int
	W, H  int
	Start int64 // Byte offset of the start of the node
	End   int64 // Byte offset of the end of the node
}

// parseSlideShapes scans the slide XML for shape nodes and extracts their properties and byte ranges.
func parseSlideShapes(content []byte) ([]parsedShape, error) {
	return scanShapesWithOffsets(content)
}

func scanShapesWithOffsets(content []byte) ([]parsedShape, error) {
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

		if se.Name.Local == "sp" || se.Name.Local == shapeTypePicture {
			// Found a shape start.
			// We need to capture the exact bytes from `startOffset` until the end element.
			// The `decoder.InputOffset()` gives the start of the token *buffer* usually, but for Bytes.Reader it's precise enough usually
			// IF we haven't read ahead.
			// Actually `InputOffset()` returns the number of bytes read *so far*.
			// So `startOffset` is the end of the *previous* token.

			// Let's extract this node.
			shape, endOffset, extractErr := extractShapeNode(content, startOffset, decoder, se.Name.Local)
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
				pShape, parseErr := buildParsedShapeFromRange(fullContent, startOffset, endOffset, stopTag)
				if parseErr != nil {
					return parsedShape{}, 0, parseErr
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
	pShape.Start = startOffset
	pShape.End = endOffset
	pShape.Type = stopTag
	return pShape, nil
}

// Minimal structs for parsing shape properties.
type shapeXML struct {
	NvSpPr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
	} `xml:"nvSpPr"`
	NvPicPr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
	} `xml:"nvPicPr"`
	SpPr struct {
		Xfrm struct {
			Off struct {
				X int `xml:"x,attr"`
				Y int `xml:"y,attr"`
			} `xml:"off"`
			Ext struct {
				Cx int `xml:"cx,attr"`
				Cy int `xml:"cy,attr"`
			} `xml:"ext"`
		} `xml:"xfrm"`
	} `xml:"spPr"`
	TxBody struct {
		P []struct {
			R []struct {
				T string `xml:"t"`
			} `xml:"r"`
		} `xml:"p"`
	} `xml:"txBody"`
}

func parseShapeProperties(content []byte) (parsedShape, error) {
	var s shapeXML
	if err := xml.Unmarshal(content, &s); err != nil {
		return parsedShape{}, err
	}

	ps := parsedShape{}
	// Extract ID/Name (handle both sp and pic variants)
	if s.NvSpPr.CNvPr.ID != 0 {
		ps.ID = s.NvSpPr.CNvPr.ID
		ps.Name = s.NvSpPr.CNvPr.Name
	} else if s.NvPicPr.CNvPr.ID != 0 {
		ps.ID = s.NvPicPr.CNvPr.ID
		ps.Name = s.NvPicPr.CNvPr.Name
	}

	// Transform
	ps.X = s.SpPr.Xfrm.Off.X
	ps.Y = s.SpPr.Xfrm.Off.Y
	ps.W = s.SpPr.Xfrm.Ext.Cx
	ps.H = s.SpPr.Xfrm.Ext.Cy

	// Text (simple accumulation)
	var txt strings.Builder
	for pIdx, p := range s.TxBody.P {
		for _, r := range p.R {
			txt.WriteString(r.T)
		}
		if pIdx < len(s.TxBody.P)-1 {
			txt.WriteString("\n") // naive paragraph join
		}
	}
	ps.Text = txt.String()

	return ps, nil
}

// replaceShapeNodes replaces the XML at the given indices.
func replaceShapeNodes(
	content []byte,
	shapes []parsedShape,
	modFunc func(i int, p *parsedShape) ([]byte, bool),
) []byte {
	// Reconstruct the file by appending chunks.
	// Must process shapes in order of offset to keep clean.
	// Optimization: Assumed shapes are sorted by offset (scanned sequentially).

	var buf bytes.Buffer
	currentOffset := int64(0)

	for i := range shapes {
		s := &shapes[i]

		// Write untouched content before this shape
		if s.Start > currentOffset {
			buf.Write(content[currentOffset:s.Start])
		}

		// Check if modification is requested
		newXML, replace := modFunc(i, s)
		if replace {
			// Write replacement
			buf.Write(newXML)
		} else {
			// Write original shape content
			buf.Write(content[s.Start:s.End])
		}

		currentOffset = s.End
	}

	// Write remainder
	if currentOffset < int64(len(content)) {
		buf.Write(content[currentOffset:])
	}

	return buf.Bytes()
}

// renderShapeXML reconstructs the XML for a shape based on its parsed properties.
func renderShapeXML(s *parsedShape) []byte {
	// Helper for XML escaping
	escape := func(s string) string {
		var buf bytes.Buffer
		if err := xml.EscapeText(&buf, []byte(s)); err != nil {
			return s
		}
		return buf.String()
	}

	if s.Type == shapeTypePicture {
		return nil
	}

	// Basic preset geometry mapping (Phase 1 supports common types)
	prst := "rect"
	switch strings.ToLower(s.Type) {
	case "ellipse", "oval":
		prst = "ellipse"
	case "triangle":
		prst = "triangle"
	}

	// Reconstruct a basic Text Shape / Rectangle
	// REDUCED: Removed redundant xmlns:a and xmlns:p as they should be at slide root.
	return fmt.Appendf(nil,
		`<p:sp>`+
			`<p:nvSpPr><p:cNvPr id="%d" name="%s"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>`+
			`<p:spPr>`+
			`<a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm>`+
			`<a:prstGeom prst="%s"><a:avLst/></a:prstGeom>`+
			`</p:spPr>`+
			`<p:txBody>`+
			`<a:bodyPr/><a:lstStyle/>`+
			// TODO: Implementation should perform a surgical update of text runs within the existing TxBody
			// instead of full re-rendering to preserve styles.
			`<a:p><a:r><a:rPr lang="en-US"/><a:t>%s</a:t></a:r></a:p>`+
			`</p:txBody>`+
			`</p:sp>`,
		s.ID, escape(s.Name),
		s.X, s.Y, s.W, s.H,
		prst,
		escape(s.Text),
	)
}

// AddShape adds a new shape to the slide.
func (e *PresentationEditor) AddShape(slideIndex int, shapeType string, x, y, w, h float64) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return 0, fmt.Errorf("read slide part %s: not found", partPath)
	}

	// Parse existing shapes to find max ID and last shape position
	shapes, err := parseSlideShapes(content)
	if err != nil {
		return 0, fmt.Errorf("parse shapes: %w", err)
	}

	maxID := maxObjectID(content)
	lastShapeEnd := int64(-1)
	for _, s := range shapes {
		if s.End > lastShapeEnd {
			lastShapeEnd = s.End
		}
	}
	newID := maxID + 1

	newShape := parsedShape{
		ID:   newID,
		Name: fmt.Sprintf("%s %d", shapeType, newID),
		Type: shapeType,
		Text: "",
		X:    int(x),
		Y:    int(y),
		W:    int(w),
		H:    int(h),
	}

	shapeXML := renderShapeXML(&newShape)

	// Insertion point: After last shape if exists, else before </p:spTree>
	var buf bytes.Buffer
	if lastShapeEnd != -1 {
		buf.Write(content[:lastShapeEnd])
		buf.Write(shapeXML)
		buf.Write(content[lastShapeEnd:])
	} else {
		endTree := []byte("</p:spTree>")
		idx := bytes.LastIndex(content, endTree)
		if idx == -1 {
			return 0, errors.New("invalid slide xml: missing spTree end")
		}
		buf.Write(content[:idx])
		buf.Write(shapeXML)
		buf.Write(content[idx:])
	}

	e.parts.Set(partPath, buf.Bytes())
	return newID, nil
}

var cNvPrIDPattern = regexp.MustCompile(`\bcNvPr\b[^>]*\bid="(\d+)"`)

const cNvPrSubmatchSize = 2

func maxObjectID(content []byte) int {
	matches := cNvPrIDPattern.FindAllSubmatch(content, -1)
	maxID := 0
	for _, match := range matches {
		if len(match) < cNvPrSubmatchSize {
			continue
		}
		id, err := strconv.Atoi(string(match[1]))
		if err != nil {
			continue
		}
		if id > maxID {
			maxID = id
		}
	}
	return maxID
}
