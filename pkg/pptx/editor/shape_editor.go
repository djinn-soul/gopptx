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

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const textRunFontSizeScale = 100

// parsedShape represents a shape found in the slide XML.
// It contains the parsed properties and the byte range of the shape node.
type parsedShape struct {
	ID          int
	Name        string
	Type        string // "sp" or "pic"
	Text        string
	Runs        []common.TextRun
	TextFrame   *common.TextFrame
	ClickAction *common.Hyperlink
	X, Y        int
	W, H        int
	PhIndex     int    // Placeholder index, -1 if not a placeholder
	PhType      string // Placeholder type (e.g. "title", "body")
	Start       int64  // Byte offset of the start of the node
	End         int64  // Byte offset of the end of the node
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

		if se.Name.Local == "sp" || se.Name.Local == shapeTypePicture || se.Name.Local == "graphicFrame" ||
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
	pShape.Start = startOffset
	pShape.End = endOffset
	pShape.Type = stopTag
	return pShape, nil
}

// Minimal structs for parsing shape properties.
type solidFillXML struct {
	SrgbClr struct {
		Val string `xml:"val,attr"`
	} `xml:"srgbClr"`
}

type runPropsXML struct {
	Bold          *bool        `xml:"b,attr"`
	Italic        *bool        `xml:"i,attr"`
	Underline     *string      `xml:"u,attr"`
	Strikethrough *string      `xml:"strike,attr"`
	Baseline      *string      `xml:"baseline,attr"`
	Caps          *string      `xml:"caps,attr"`
	SmallCaps     *string      `xml:"smCaps,attr"`
	SolidFill     solidFillXML `xml:"solidFill"`
	Highlight     solidFillXML `xml:"highlight"`
}

type shapeXML struct {
	NvSpPr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
		NvPr struct {
			Ph *struct {
				Idx  *int   `xml:"idx,attr"`
				Type string `xml:"type,attr"`
			} `xml:"ph"`
		} `xml:"nvPr"`
	} `xml:"nvSpPr"`
	NvPicPr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
		NvPr struct {
			Ph *struct {
				Idx  *int   `xml:"idx,attr"`
				Type string `xml:"type,attr"`
			} `xml:"ph"`
		} `xml:"nvPr"`
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
				RPr *runPropsXML `xml:"rPr"`
				T   string       `xml:"t"`
			} `xml:"r"`
		} `xml:"p"`
	} `xml:"txBody"`
}

//nolint:gocognit,funlen,nestif // XML shape parsing must handle many optional OOXML branches in one pass.
func parseShapeProperties(content []byte) (parsedShape, error) {
	var s shapeXML
	if err := xml.Unmarshal(content, &s); err != nil {
		return parsedShape{}, err
	}

	ps := parsedShape{
		PhIndex: -1,
	}
	// Extract ID/Name (handle both sp and pic variants)
	if s.NvSpPr.CNvPr.ID != 0 {
		ps.ID = s.NvSpPr.CNvPr.ID
		ps.Name = s.NvSpPr.CNvPr.Name
		if s.NvSpPr.NvPr.Ph != nil {
			ps.PhType = s.NvSpPr.NvPr.Ph.Type
			if s.NvSpPr.NvPr.Ph.Idx != nil {
				ps.PhIndex = *s.NvSpPr.NvPr.Ph.Idx
			} else {
				// idx is optional and defaults to 0 in OOXML for <p:ph/>
				ps.PhIndex = 0
			}
		}
	} else if s.NvPicPr.CNvPr.ID != 0 {
		ps.ID = s.NvPicPr.CNvPr.ID
		ps.Name = s.NvPicPr.CNvPr.Name
		if s.NvPicPr.NvPr.Ph != nil {
			ps.PhType = s.NvPicPr.NvPr.Ph.Type
			if s.NvPicPr.NvPr.Ph.Idx != nil {
				ps.PhIndex = *s.NvPicPr.NvPr.Ph.Idx
			} else {
				// idx is optional and defaults to 0 in OOXML for <p:ph/>
				ps.PhIndex = 0
			}
		}
	}

	// Transform
	ps.X = s.SpPr.Xfrm.Off.X
	ps.Y = s.SpPr.Xfrm.Off.Y
	ps.W = s.SpPr.Xfrm.Ext.Cx
	ps.H = s.SpPr.Xfrm.Ext.Cy

	// Text (simple accumulation) and Runs parsing
	var txt strings.Builder
	for pIdx, p := range s.TxBody.P {
		for _, r := range p.R {
			txt.WriteString(r.T)

			// Build TextRun from parsed properties
			if r.RPr != nil || r.T != "" {
				run := common.TextRun{Text: r.T}
				if r.RPr != nil {
					if r.RPr.Bold != nil && *r.RPr.Bold {
						run.Bold = r.RPr.Bold
					}
					if r.RPr.Italic != nil && *r.RPr.Italic {
						run.Italic = r.RPr.Italic
					}
					if r.RPr.Underline != nil && *r.RPr.Underline != "" {
						run.Underline = r.RPr.Underline
					}
					if r.RPr.Strikethrough != nil && *r.RPr.Strikethrough != "" {
						run.Strikethrough = r.RPr.Strikethrough
					}
					switch {
					case parseIntAttr(r.RPr.Baseline) < 0:
						v := true
						run.Subscript = &v
					case parseIntAttr(r.RPr.Baseline) > 0:
						v := true
						run.Superscript = &v
					}
					if r.RPr.Caps != nil {
						switch strings.ToLower(strings.TrimSpace(*r.RPr.Caps)) {
						case "all":
							v := true
							run.AllCaps = &v
						case "small":
							v := true
							run.SmallCaps = &v
						}
					}
					if parseXMLBoolAttr(r.RPr.SmallCaps) {
						v := true
						run.SmallCaps = &v
					}
					if r.RPr.SolidFill.SrgbClr.Val != "" {
						val := r.RPr.SolidFill.SrgbClr.Val
						run.Color = &val
					}
					if r.RPr.Highlight.SrgbClr.Val != "" {
						val := r.RPr.Highlight.SrgbClr.Val
						run.Highlight = &val
					}
				}
				ps.Runs = append(ps.Runs, run)
			}
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

// renderTextBodyXML constructs the <p:txBody> node based on Text or Runs.
// If a PresentationEditor is provided, it will register any hyperlink relationships.
//
//nolint:gocognit,gocyclo,cyclop,funlen,nestif // Text-body emission covers run/paragraph/hyperlink variants required for PPTX fidelity.
func renderTextBodyXML(e *PresentationEditor, partPath string, s *parsedShape) ([]byte, error) {
	escape := func(str string) string {
		var buf bytes.Buffer
		_ = xml.EscapeText(&buf, []byte(str))
		return buf.String()
	}

	var txBody bytes.Buffer
	txBody.WriteString(`<p:txBody>`)

	// Emit custom bodyPr attributes if TextFrame is provided
	bodyPr := `<a:bodyPr`
	if s.TextFrame != nil {
		tf := s.TextFrame
		if tf.MarginTop != nil {
			bodyPr += fmt.Sprintf(` tIns="%d"`, *tf.MarginTop)
		}
		if tf.MarginBottom != nil {
			bodyPr += fmt.Sprintf(` bIns="%d"`, *tf.MarginBottom)
		}
		if tf.MarginLeft != nil {
			bodyPr += fmt.Sprintf(` lIns="%d"`, *tf.MarginLeft)
		}
		if tf.MarginRight != nil {
			bodyPr += fmt.Sprintf(` rIns="%d"`, *tf.MarginRight)
		}
		if tf.WordWrap != nil {
			if *tf.WordWrap {
				bodyPr += ` wrap="square"`
			} else {
				bodyPr += ` wrap="none"`
			}
		}
		if tf.VerticalAlign != nil && *tf.VerticalAlign != "" {
			bodyPr += fmt.Sprintf(` anchor="%s"`, escape(*tf.VerticalAlign))
		}
		bodyPr += `>`

		if tf.AutoFitType != nil {
			switch *tf.AutoFitType {
			case "normal":
				bodyPr += `<a:normAutofit/>`
			case "shape":
				bodyPr += `<a:spAutoFit/>`
			case "none":
				bodyPr += `<a:noAutofit/>`
			}
		} else if tf.AutoFit != nil {
			// Backwards compatibility with boolean field
			if *tf.AutoFit {
				bodyPr += `<a:spAutoFit/>`
			} else {
				bodyPr += `<a:noAutofit/>`
			}
		}
		bodyPr += `</a:bodyPr>`
	} else {
		bodyPr += `/>`
	}

	txBody.WriteString(bodyPr)
	txBody.WriteString(`<a:lstStyle/>`)
	if len(s.Runs) > 0 {
		txBody.WriteString(`<a:p>`)
		for _, r := range s.Runs {
			rPr := `<a:rPr lang="en-US"`
			if r.Bold != nil && *r.Bold {
				rPr += ` b="1"`
			}
			if r.Italic != nil && *r.Italic {
				rPr += ` i="1"`
			}
			if r.Underline != nil && *r.Underline != "" {
				rPr += fmt.Sprintf(` u="%s"`, escape(*r.Underline))
			}
			if r.Strikethrough != nil && *r.Strikethrough != "" {
				val := *r.Strikethrough
				switch val {
				case "sng":
					val = "sngStrike"
				case "dbl":
					val = "dblStrike"
				}
				rPr += fmt.Sprintf(` strike="%s"`, escape(val))
			}
			if r.Subscript != nil && *r.Subscript {
				rPr += ` baseline="-25000"`
			}
			if r.Superscript != nil && *r.Superscript {
				rPr += ` baseline="30000"`
			}
			if r.SizePt != nil && *r.SizePt > 0 {
				rPr += fmt.Sprintf(` sz="%d"`, *r.SizePt*textRunFontSizeScale)
			}
			if r.AllCaps != nil && *r.AllCaps {
				rPr += ` caps="all"`
			}
			if r.SmallCaps != nil && *r.SmallCaps {
				rPr += ` smCaps="1"`
			}
			rPr += `>`

			if r.Color != nil && *r.Color != "" {
				rPr += fmt.Sprintf(`<a:solidFill><a:srgbClr val="%s"/></a:solidFill>`, escape(*r.Color))
			}
			if r.Highlight != nil && *r.Highlight != "" {
				rPr += fmt.Sprintf(`<a:highlight><a:srgbClr val="%s"/></a:highlight>`, escape(*r.Highlight))
			}
			if r.Font != nil && *r.Font != "" {
				rPr += fmt.Sprintf(`<a:latin typeface="%s"/><a:cs typeface="%s"/>`, escape(*r.Font), escape(*r.Font))
			}

			if r.Hyperlink != nil && e != nil && partPath != "" {
				clickXML, err := e.buildClickActionXML(partPath, r.Hyperlink)
				if err != nil {
					return nil, err
				}
				if clickXML != "" {
					rPr += clickXML
				}
			}

			rPr += `</a:rPr>`
			txBody.WriteString(fmt.Sprintf(`<a:r>%s<a:t>%s</a:t></a:r>`, rPr, escape(r.Text)))
		}
		txBody.WriteString(`</a:p>`)
	} else {
		txBody.WriteString(fmt.Sprintf(`<a:p><a:r><a:rPr lang="en-US"/><a:t>%s</a:t></a:r></a:p>`, escape(s.Text)))
	}
	txBody.WriteString(`</p:txBody>`)

	return txBody.Bytes(), nil
}

// renderShapeXML reconstructs the XML for a shape based on its parsed properties.
func (e *PresentationEditor) renderShapeXML(partPath string, s *parsedShape) ([]byte, error) {
	// Helper for XML escaping
	escape := func(s string) string {
		var buf bytes.Buffer
		if err := xml.EscapeText(&buf, []byte(s)); err != nil {
			return s
		}
		return buf.String()
	}

	if s.Type == shapeTypePicture {
		return nil, nil
	}

	// Basic preset geometry mapping (Phase 1 supports common types)
	prst := "rect"
	switch strings.ToLower(s.Type) {
	case "ellipse", "oval":
		prst = "ellipse"
	case "triangle":
		prst = "triangle"
	}

	txBody, err := renderTextBodyXML(e, partPath, s)
	if err != nil {
		return nil, err
	}

	clickXML, err := e.buildClickActionXML(partPath, s.ClickAction)
	if err != nil {
		return nil, err
	}

	return fmt.Appendf(
		nil,
		`<p:sp xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`+
			`<p:nvSpPr><p:cNvPr id="%d" name="%s">%s</p:cNvPr><p:cNvSpPr/><p:nvPr/></p:nvSpPr>`+
			`<p:spPr>`+
			`<a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm>`+
			`<a:prstGeom prst="%s"><a:avLst/></a:prstGeom>`+
			`</p:spPr>`+
			`%s`+
			`</p:sp>`,
		s.ID,
		escape(s.Name),
		clickXML,
		s.X,
		s.Y,
		s.W,
		s.H,
		prst,
		string(txBody),
	), nil
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
	// OPTIMIZATION: We only need the offsets, not the full properties.
	shapes, err := scanShapesWithOffsets(content, true) // true = skip properties parsing
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

	shapeXML, err := e.renderShapeXML(partPath, &newShape)
	if err != nil {
		return 0, err
	}

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

// UpdateShape updates an existing shape's properties.
//
//nolint:gocognit,nestif // Update path intentionally coordinates parsing, targeted mutation, and XML rewrite atomically.
func (e *PresentationEditor) UpdateShape(slideIndex, shapeID int, updates common.ShapeUpdate) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapes, err := parseSlideShapes(content) // parses basic properties and byte ranges
	if err != nil {
		return fmt.Errorf("parse shapes: %w", err)
	}

	found := false
	var updateErr error
	newXML := replaceShapeNodes(content, shapes, func(_ int, s *parsedShape) ([]byte, bool) {
		if updateErr != nil {
			return nil, false
		}
		if s.ID == shapeID || (s.PhType == "title" && shapeID == 0) { // Naive placeholder fallback
			found = true
			replace := false

			originalXML := content[s.Start:s.End]
			updatedXML := originalXML

			if updates.X != nil || updates.Y != nil || updates.W != nil || updates.H != nil {
				replace = true
				if updates.X != nil {
					s.X = *updates.X
				}
				if updates.Y != nil {
					s.Y = *updates.Y
				}
				if updates.W != nil {
					s.W = *updates.W
				}
				if updates.H != nil {
					s.H = *updates.H
				}
				updatedXML = updateShapeTransforms(updatedXML, s.X, s.Y, s.W, s.H)
			}

			if updates.Text != nil || updates.Runs != nil || updates.TextFrame != nil {
				replace = true
				if updates.Text != nil {
					s.Text = *updates.Text
					s.Runs = nil // Override runs if raw text is provided
				}
				if updates.Runs != nil {
					s.Runs = *updates.Runs
				}
				if updates.TextFrame != nil {
					s.TextFrame = updates.TextFrame
				}
				updatedXML, updateErr = replaceShapeTextBody(e, partPath, updatedXML, s)
				if updateErr != nil {
					return nil, false
				}
			}
			if updates.ClickAction != nil {
				replace = true
				s.ClickAction = updates.ClickAction
				updatedXML, updateErr = replaceShapeClickAction(e, partPath, updatedXML, s.ClickAction)
				if updateErr != nil {
					return nil, false
				}
			}

			if replace {
				return updatedXML, true
			}
		}
		return nil, false
	})
	if updateErr != nil {
		return updateErr
	}

	if !found {
		return fmt.Errorf("shape id %d not found on slide %d", shapeID, slideIndex)
	}

	e.parts.Set(partPath, newXML)
	return nil
}

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

func getStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func parseIntAttr(value *string) int {
	if value == nil {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(*value))
	if err != nil {
		return 0
	}
	return n
}

func parseXMLBoolAttr(value *string) bool {
	if value == nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(*value)) {
	case "1", "true", "on", "yes":
		return true
	default:
		return false
	}
}

func xmlEscape(value string) string {
	var buf bytes.Buffer
	if err := xml.EscapeText(&buf, []byte(value)); err != nil {
		return value
	}
	return buf.String()
}

func (e *PresentationEditor) getOrCreateHyperlinkRelID(partPath, address string) (string, error) {
	relsPath := common.SlideRelsPartName(partPath)
	rels := make([]common.EditorRelationship, 0)
	if data, ok := e.parts.Get(relsPath); ok {
		parsed, err := parseRelationshipsXML(data)
		if err != nil {
			return "", fmt.Errorf("parse %s: %w", relsPath, err)
		}
		rels = parsed
	}
	for _, r := range rels {
		if r.Type == common.RelTypeHyperlink && r.Target == address {
			return r.ID, nil
		}
	}
	relID, err := e.nextSlideRelID(partPath)
	if err != nil {
		return "", err
	}
	if err := e.addRelationship(partPath, relID, common.RelTypeHyperlink, address); err != nil {
		return "", err
	}
	return relID, nil
}

func (e *PresentationEditor) buildClickActionXML(partPath string, hl *common.Hyperlink) (string, error) {
	if hl == nil || partPath == "" {
		return "", nil
	}

	if hl.Address != nil && *hl.Address != "" {
		relID, err := e.getOrCreateHyperlinkRelID(partPath, *hl.Address)
		if err != nil {
			return "", fmt.Errorf("allocate hyperlink relationship id: %w", err)
		}
		return fmt.Sprintf(`<a:hlinkClick r:id="%s" action="%s" tooltip="%s"/>`,
			xmlEscape(relID), xmlEscape(getStr(hl.Action)), xmlEscape(getStr(hl.Tooltip))), nil
	}

	if hl.Action != nil && *hl.Action != "" {
		return fmt.Sprintf(`<a:hlinkClick action="%s" tooltip="%s"/>`,
			xmlEscape(*hl.Action), xmlEscape(getStr(hl.Tooltip))), nil
	}

	return "", nil
}

// updateShapeTransforms performs a surgical regular expression replacement of shape transforms.
func updateShapeTransforms(xmlData []byte, x, y, w, h int) []byte {
	offRe := regexp.MustCompile(`(<a:off\s+x=")[0-9]+("\s+y=")[0-9]+("\s*/>)`)
	extRe := regexp.MustCompile(`(<a:ext\s+cx=")[0-9]+("\s+cy=")[0-9]+("\s*/>)`)

	res := offRe.ReplaceAll(xmlData, fmt.Appendf(nil, "${1}%d${2}%d${3}", x, y))
	res = extRe.ReplaceAll(res, fmt.Appendf(nil, "${1}%d${2}%d${3}", w, h))
	return res
}

// replaceShapeTextBody replaces the entire <p:txBody> node with a newly constructed one based on Text/Runs.
func replaceShapeTextBody(
	e *PresentationEditor,
	partPath string,
	xmlData []byte,
	s *parsedShape,
) ([]byte, error) {
	// Robustly find <p:txBody> (it might have attributes or different spacing)
	txBodyOpenRe := regexp.MustCompile(`(?s)<p:txBody\b[^>]*>`)
	txBodyCloseTag := []byte("</p:txBody>")

	loc := txBodyOpenRe.FindIndex(xmlData)
	var startIdx, endIdx int
	found := false

	if loc != nil {
		startIdx = loc[0]
		// Search for closing tag after the opening tag
		closeIdx := bytes.Index(xmlData[loc[1]:], txBodyCloseTag)
		if closeIdx != -1 {
			endIdx = loc[1] + closeIdx
			found = true
		}
	}

	txBody, err := renderTextBodyXML(e, partPath, s)
	if err != nil {
		return nil, err
	}

	if !found {
		// If txBody doesn't exist, try to insert after spPr.
		// Use regex for spPr as well for consistency.
		spPrEndRe := regexp.MustCompile(`(?s)</p:spPr>`)
		spPrLoc := spPrEndRe.FindIndex(xmlData)
		if spPrLoc != nil {
			insertIdx := spPrLoc[1]
			res := make([]byte, 0, len(xmlData)+len(txBody))
			res = append(res, xmlData[:insertIdx]...)
			res = append(res, txBody...)
			res = append(res, xmlData[insertIdx:]...)
			return res, nil
		}
		return xmlData, nil // Give up
	}

	res := make([]byte, 0, len(xmlData)-(endIdx+len(txBodyCloseTag)-startIdx)+len(txBody))
	res = append(res, xmlData[:startIdx]...)
	res = append(res, txBody...)
	res = append(res, xmlData[endIdx+len(txBodyCloseTag):]...)

	return res, nil
}

func replaceShapeClickAction(
	e *PresentationEditor,
	partPath string,
	xmlData []byte,
	clickAction *common.Hyperlink,
) ([]byte, error) {
	clickXML, err := e.buildClickActionXML(partPath, clickAction)
	if err != nil {
		return nil, err
	}

	xmlStr := string(xmlData)
	hlinkClickPattern := regexp.MustCompile(`(?s)<a:hlinkClick\b[^>]*/>|<a:hlinkClick\b[^>]*>.*?</a:hlinkClick>`)
	cNvPrOpenClose := regexp.MustCompile(`(?s)<p:cNvPr\b([^>]*)>(.*?)</p:cNvPr>`)
	if match := cNvPrOpenClose.FindStringSubmatchIndex(xmlStr); match != nil {
		inner := xmlStr[match[4]:match[5]]
		removeHlinks := func(input string) string {
			matches := hlinkClickPattern.FindAllStringIndex(input, -1)
			if len(matches) == 0 {
				return input
			}
			var builder strings.Builder
			builder.Grow(len(input))
			last := 0
			for _, m := range matches {
				builder.WriteString(input[last:m[0]])
				last = m[1]
			}
			builder.WriteString(input[last:])
			return builder.String()
		}
		cleanInner := removeHlinks(inner)
		var replacement strings.Builder
		replacement.WriteString(cleanInner)
		if clickXML != "" {
			replacement.WriteString(clickXML)
		}
		updated := xmlStr[:match[4]] + replacement.String() + xmlStr[match[5]:]
		return []byte(updated), nil
	}

	cNvPrSelfClosing := regexp.MustCompile(`<p:cNvPr\b([^>]*)/>`)
	if match := cNvPrSelfClosing.FindStringSubmatchIndex(xmlStr); match != nil {
		if clickXML == "" {
			return xmlData, nil
		}
		attrs := xmlStr[match[2]:match[3]]
		replacement := fmt.Sprintf(`<p:cNvPr%s>%s</p:cNvPr>`, attrs, clickXML)
		updated := xmlStr[:match[0]] + replacement + xmlStr[match[1]:]
		return []byte(updated), nil
	}

	if clickAction != nil {
		return nil, errors.New("shape has no cNvPr node for click action update")
	}

	return xmlData, nil
}
