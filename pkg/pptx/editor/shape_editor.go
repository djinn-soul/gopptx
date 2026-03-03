package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const textRunFontSizeScale = 100
const minTextFrameColumns = 1

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
			se.Name.Local == "grpSp" ||
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
	NvGrpSpPr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
	} `xml:"nvGrpSpPr"`
	SpPr struct {
		NoFill    *struct{}     `xml:"noFill"`
		SolidFill *solidFillXML `xml:"solidFill"`
		GradFill  *struct {
			Lin *struct {
				Ang *int `xml:"ang,attr"`
			} `xml:"lin"`
			GsLst struct {
				Gs []struct {
					Pos     *int `xml:"pos,attr"`
					SrgbClr *struct {
						Val string `xml:"val,attr"`
					} `xml:"srgbClr"`
				} `xml:"gs"`
			} `xml:"gsLst"`
		} `xml:"gradFill"`
		PattFill *struct {
			Prst  *string `xml:"prst,attr"`
			FgClr *struct {
				SrgbClr *struct {
					Val string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"fgClr"`
			BgClr *struct {
				SrgbClr *struct {
					Val string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"bgClr"`
		} `xml:"pattFill"`
		Ln *struct {
			W         *int          `xml:"w,attr"`
			SolidFill *solidFillXML `xml:"solidFill"`
			PrstDash  *struct {
				Val string `xml:"val,attr"`
			} `xml:"prstDash"`
		} `xml:"ln"`
		EffectLst *struct {
			OuterShdw *struct {
				BlurRad *int `xml:"blurRad,attr"`
				Dist    *int `xml:"dist,attr"`
				Dir     *int `xml:"dir,attr"`
				SrgbClr *struct {
					Val string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"outerShdw"`
			Glow *struct {
				Rad     *int `xml:"rad,attr"`
				SrgbClr *struct {
					Val string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"glow"`
			Blur *struct {
				Rad *int `xml:"rad,attr"`
			} `xml:"blur"`
			SoftEdge *struct {
				Rad *int `xml:"rad,attr"`
			} `xml:"softEdge"`
			Reflection *struct {
				BlurRad *int `xml:"blurRad,attr"`
				Dist    *int `xml:"dist,attr"`
			} `xml:"reflection"`
		} `xml:"effectLst"`
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
	GrpSpPr struct {
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
	} `xml:"grpSpPr"`
	TxBody struct {
		P []struct {
			PPr *struct {
				MarL   *int `xml:"marL,attr"`
				Indent *int `xml:"indent,attr"`
			} `xml:"pPr"`
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
	if s.SpPr.NoFill != nil {
		background := true
		ps.Fill = &common.ShapeFill{Background: &background}
	}
	if s.SpPr.SolidFill != nil && s.SpPr.SolidFill.SrgbClr.Val != "" {
		fillColor := s.SpPr.SolidFill.SrgbClr.Val
		ps.Fill = &common.ShapeFill{Solid: &fillColor}
	}
	if s.SpPr.GradFill != nil {
		grad := &common.GradientFill{}
		if s.SpPr.GradFill.Lin != nil && s.SpPr.GradFill.Lin.Ang != nil {
			angle := float64(*s.SpPr.GradFill.Lin.Ang) / rotationDegreeToOOXML
			grad.AngleDeg = &angle
		}
		for _, gs := range s.SpPr.GradFill.GsLst.Gs {
			if gs.SrgbClr == nil || gs.SrgbClr.Val == "" {
				continue
			}
			stop := common.GradientStop{Color: gs.SrgbClr.Val}
			if gs.Pos != nil {
				pos := float64(*gs.Pos) / 1000.0
				stop.PositionPct = &pos
			}
			grad.Stops = append(grad.Stops, stop)
		}
		ps.Fill = &common.ShapeFill{Gradient: grad}
	}
	if s.SpPr.PattFill != nil {
		pattern := &common.PatternedFill{}
		if s.SpPr.PattFill.Prst != nil {
			pattern.Preset = s.SpPr.PattFill.Prst
		}
		if s.SpPr.PattFill.FgClr != nil &&
			s.SpPr.PattFill.FgClr.SrgbClr != nil &&
			s.SpPr.PattFill.FgClr.SrgbClr.Val != "" {
			fg := s.SpPr.PattFill.FgClr.SrgbClr.Val
			pattern.FgColor = &fg
		}
		if s.SpPr.PattFill.BgClr != nil &&
			s.SpPr.PattFill.BgClr.SrgbClr != nil &&
			s.SpPr.PattFill.BgClr.SrgbClr.Val != "" {
			bg := s.SpPr.PattFill.BgClr.SrgbClr.Val
			pattern.BgColor = &bg
		}
		ps.Fill = &common.ShapeFill{Pattern: pattern}
	}
	if s.SpPr.Ln != nil {
		line := &common.ShapeLine{}
		if s.SpPr.Ln.SolidFill != nil && s.SpPr.Ln.SolidFill.SrgbClr.Val != "" {
			lineColor := s.SpPr.Ln.SolidFill.SrgbClr.Val
			line.Color = &lineColor
		}
		if s.SpPr.Ln.W != nil {
			line.WidthEmu = s.SpPr.Ln.W
		}
		if s.SpPr.Ln.PrstDash != nil && s.SpPr.Ln.PrstDash.Val != "" {
			dash := s.SpPr.Ln.PrstDash.Val
			line.DashStyle = &dash
		}
		if line.Color != nil || line.WidthEmu != nil || line.DashStyle != nil {
			ps.Line = line
		}
	}
	if s.SpPr.EffectLst != nil {
		if s.SpPr.EffectLst.OuterShdw == nil && s.SpPr.EffectLst.Glow == nil && s.SpPr.EffectLst.Blur == nil && s.SpPr.EffectLst.SoftEdge == nil && s.SpPr.EffectLst.Reflection == nil {
			inherit := false
			ps.Shadow = &common.ShapeShadow{Inherit: &inherit}
		} else if s.SpPr.EffectLst.OuterShdw != nil {
			outer := s.SpPr.EffectLst.OuterShdw
			shadow := &common.ShapeShadow{}
			if outer.SrgbClr != nil && outer.SrgbClr.Val != "" {
				color := outer.SrgbClr.Val
				shadow.Color = &color
			}
			if outer.BlurRad != nil {
				shadow.BlurEmu = outer.BlurRad
			}
			if outer.Dist != nil {
				shadow.DistanceEmu = outer.Dist
			}
			if outer.Dir != nil {
				angle := float64(*outer.Dir) / rotationDegreeToOOXML
				shadow.AngleDeg = &angle
			}
			ps.Shadow = shadow
		}
		if s.SpPr.EffectLst.Glow != nil {
			glow := &common.ShapeGlow{}
			if s.SpPr.EffectLst.Glow.SrgbClr != nil && s.SpPr.EffectLst.Glow.SrgbClr.Val != "" {
				color := s.SpPr.EffectLst.Glow.SrgbClr.Val
				glow.Color = &color
			}
			if s.SpPr.EffectLst.Glow.Rad != nil {
				glow.RadiusEmu = s.SpPr.EffectLst.Glow.Rad
			}
			ps.Glow = glow
		}
		if s.SpPr.EffectLst.Blur != nil {
			blur := &common.ShapeBlur{}
			if s.SpPr.EffectLst.Blur.Rad != nil {
				blur.RadiusEmu = s.SpPr.EffectLst.Blur.Rad
			}
			ps.Blur = blur
		}
		if s.SpPr.EffectLst.SoftEdge != nil {
			softEdge := &common.ShapeSoftEdge{}
			if s.SpPr.EffectLst.SoftEdge.Rad != nil {
				softEdge.RadiusEmu = s.SpPr.EffectLst.SoftEdge.Rad
			}
			ps.SoftEdge = softEdge
		}
		if s.SpPr.EffectLst.Reflection != nil {
			reflection := &common.ShapeReflection{}
			if s.SpPr.EffectLst.Reflection.BlurRad != nil {
				reflection.BlurEmu = s.SpPr.EffectLst.Reflection.BlurRad
			}
			if s.SpPr.EffectLst.Reflection.Dist != nil {
				reflection.DistanceEmu = s.SpPr.EffectLst.Reflection.Dist
			}
			ps.Reflection = reflection
		}
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
	} else if s.NvGrpSpPr.CNvPr.ID != 0 {
		ps.ID = s.NvGrpSpPr.CNvPr.ID
		ps.Name = s.NvGrpSpPr.CNvPr.Name
	}

	// Transform
	if s.SpPr.Xfrm.Ext.Cx != 0 || s.SpPr.Xfrm.Ext.Cy != 0 || s.SpPr.Xfrm.Off.X != 0 || s.SpPr.Xfrm.Off.Y != 0 {
		ps.X = s.SpPr.Xfrm.Off.X
		ps.Y = s.SpPr.Xfrm.Off.Y
		ps.W = s.SpPr.Xfrm.Ext.Cx
		ps.H = s.SpPr.Xfrm.Ext.Cy
	} else {
		ps.X = s.GrpSpPr.Xfrm.Off.X
		ps.Y = s.GrpSpPr.Xfrm.Off.Y
		ps.W = s.GrpSpPr.Xfrm.Ext.Cx
		ps.H = s.GrpSpPr.Xfrm.Ext.Cy
	}

	// Text (simple accumulation) and Runs parsing
	var txt strings.Builder
	for pIdx, p := range s.TxBody.P {
		if pIdx == 0 && p.PPr != nil {
			paragraph := &common.Paragraph{}
			if p.PPr.MarL != nil {
				paragraph.Indent = p.PPr.MarL
			}
			if p.PPr.Indent != nil && *p.PPr.Indent < 0 {
				hanging := -*p.PPr.Indent
				paragraph.Hanging = &hanging
			}
			if paragraph.Indent != nil || paragraph.Hanging != nil {
				ps.Paragraph = paragraph
			}
		}
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
	txBody.WriteString(`<p:txBody xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`)

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
		if tf.Orientation != nil && *tf.Orientation != "" {
			orientation, err := normalizeTextFrameOrientation(*tf.Orientation)
			if err != nil {
				return nil, err
			}
			bodyPr += fmt.Sprintf(` vert="%s"`, escape(orientation))
		}
		if tf.Columns != nil {
			if *tf.Columns < minTextFrameColumns {
				return nil, fmt.Errorf("text_frame.columns must be >= %d", minTextFrameColumns)
			}
			bodyPr += fmt.Sprintf(` numCol="%d"`, *tf.Columns)
		}
		if tf.Rotation != nil {
			rotation, err := normalizeTextFrameRotation(*tf.Rotation)
			if err != nil {
				return nil, err
			}
			bodyPr += fmt.Sprintf(` rot="%d"`, rotation)
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
		if s.Paragraph != nil {
			paragraphXML, err := renderParagraphPropsXML(s.Paragraph)
			if err != nil {
				return nil, err
			}
			if paragraphXML != "" {
				txBody.WriteString(paragraphXML)
			}
		}
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
			if r.HoverAction != nil && e != nil && partPath != "" {
				hoverXML, err := e.buildHoverActionXML(partPath, r.HoverAction)
				if err != nil {
					return nil, err
				}
				if hoverXML != "" {
					rPr += hoverXML
				}
			}

			rPr += `</a:rPr>`
			txBody.WriteString(fmt.Sprintf(`<a:r>%s<a:t>%s</a:t></a:r>`, rPr, escape(r.Text)))
		}
		txBody.WriteString(`</a:p>`)
	} else {
		if s.Paragraph != nil {
			paragraphXML, err := renderParagraphPropsXML(s.Paragraph)
			if err != nil {
				return nil, err
			}
			if paragraphXML != "" {
				txBody.WriteString(`<a:p>`)
				txBody.WriteString(paragraphXML)
				txBody.WriteString(fmt.Sprintf(`<a:r><a:rPr lang="en-US"/><a:t>%s</a:t></a:r></a:p>`, escape(s.Text)))
			} else {
				txBody.WriteString(fmt.Sprintf(`<a:p><a:r><a:rPr lang="en-US"/><a:t>%s</a:t></a:r></a:p>`, escape(s.Text)))
			}
		} else {
			txBody.WriteString(fmt.Sprintf(`<a:p><a:r><a:rPr lang="en-US"/><a:t>%s</a:t></a:r></a:p>`, escape(s.Text)))
		}
	}
	txBody.WriteString(`</p:txBody>`)

	return txBody.Bytes(), nil
}

func normalizeTextFrameOrientation(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "horz", "horizontal":
		return "horz", nil
	case "vert", "vertical":
		return "vert", nil
	case "vert270", "vertical270", "vertical_270":
		return "vert270", nil
	case "wordartvert", "word_art_vert":
		return "wordArtVert", nil
	case "eavert", "ea_vert":
		return "eaVert", nil
	case "mongolianvert", "mongolian_vert":
		return "mongolianVert", nil
	case "wordartvertrtl", "word_art_vert_rtl":
		return "wordArtVertRtl", nil
	default:
		return "", fmt.Errorf("unsupported text_frame.orientation %q", raw)
	}
}

func normalizeTextFrameRotation(raw float64) (int64, error) {
	if math.IsNaN(raw) || math.IsInf(raw, 0) {
		return 0, errors.New("text_frame.rotation must be finite")
	}
	if raw < -360.0 || raw > 360.0 {
		return 0, fmt.Errorf("text_frame.rotation must be between -360 and 360 degrees")
	}
	return int64(math.Round(raw * rotationDegreeToOOXML)), nil
}

func renderParagraphPropsXML(paragraph *common.Paragraph) (string, error) {
	if paragraph == nil {
		return "", nil
	}
	pPr := `<a:pPr`
	if paragraph.Indent != nil {
		pPr += fmt.Sprintf(` marL="%d"`, *paragraph.Indent)
	}
	if paragraph.Hanging != nil {
		if *paragraph.Hanging < 0 {
			return "", errors.New("paragraph.hanging must be >= 0")
		}
		pPr += fmt.Sprintf(` indent="%d"`, -*paragraph.Hanging)
	}
	pPr += `/>`
	if pPr == `<a:pPr/>` {
		return "", nil
	}
	return pPr, nil
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
	hoverXML, err := e.buildHoverActionXML(partPath, s.HoverAction)
	if err != nil {
		return nil, err
	}
	styleXML, err := renderShapeStyleXML(s.Fill, s.Line, s.Shadow, s.Glow, s.Blur, s.SoftEdge, s.Reflection)
	if err != nil {
		return nil, err
	}

	return fmt.Appendf(
		nil,
		`<p:sp xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`+
			`<p:nvSpPr><p:cNvPr id="%d" name="%s">%s%s</p:cNvPr><p:cNvSpPr/><p:nvPr/></p:nvSpPr>`+
			`<p:spPr>`+
			`<a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm>`+
			`%s`+
			`<a:prstGeom prst="%s"><a:avLst/></a:prstGeom>`+
			`</p:spPr>`+
			`%s`+
			`</p:sp>`,
		s.ID,
		escape(s.Name),
		clickXML,
		hoverXML,
		s.X,
		s.Y,
		s.W,
		s.H,
		styleXML,
		prst,
		string(txBody),
	), nil
}

func renderShapeStyleXML(
	fill *common.ShapeFill,
	line *common.ShapeLine,
	shadow *common.ShapeShadow,
	glow *common.ShapeGlow,
	blur *common.ShapeBlur,
	softEdge *common.ShapeSoftEdge,
	reflection *common.ShapeReflection,
) (string, error) {
	var style strings.Builder
	if fill != nil {
		fillXML, err := renderFillXML(fill)
		if err != nil {
			return "", err
		}
		style.WriteString(fillXML)
	}
	if line != nil {
		lnAttrs := ""
		if line.WidthEmu != nil {
			if *line.WidthEmu <= 0 {
				return "", errors.New("line.width_emu must be > 0")
			}
			lnAttrs = fmt.Sprintf(` w="%d"`, *line.WidthEmu)
		}
		lineColor := ""
		if line.Color != nil {
			color, err := normalizeHexColor(*line.Color)
			if err != nil {
				return "", fmt.Errorf("line.color: %w", err)
			}
			lineColor = color
		}
		lineDash := ""
		if line.DashStyle != nil {
			dash, err := normalizeLineDashStyle(*line.DashStyle)
			if err != nil {
				return "", fmt.Errorf("line.dash_style: %w", err)
			}
			lineDash = dash
		}
		style.WriteString(`<a:ln`)
		style.WriteString(lnAttrs)
		if lineColor == "" && lineDash == "" {
			style.WriteString(`/>`)
		} else {
			style.WriteString(`>`)
			if lineDash != "" {
				style.WriteString(`<a:prstDash val="`)
				style.WriteString(lineDash)
				style.WriteString(`"/>`)
			}
			if lineColor != "" {
				style.WriteString(`<a:solidFill><a:srgbClr val="`)
				style.WriteString(lineColor)
				style.WriteString(`"/></a:solidFill>`)
			}
			style.WriteString(`</a:ln>`)
		}
	}
	effectsXML, err := renderEffectsXML(shadow, glow, blur, softEdge, reflection)
	if err != nil {
		return "", err
	}
	if effectsXML != "" {
		style.WriteString(effectsXML)
	}
	return style.String(), nil
}

func renderEffectsXML(
	shadow *common.ShapeShadow,
	glow *common.ShapeGlow,
	blur *common.ShapeBlur,
	softEdge *common.ShapeSoftEdge,
	reflection *common.ShapeReflection,
) (string, error) {
	if shadow == nil && glow == nil && blur == nil && softEdge == nil && reflection == nil {
		return "", nil
	}
	if shadow != nil && shadow.Inherit != nil {
		if shadow.Color != nil || shadow.BlurEmu != nil || shadow.DistanceEmu != nil || shadow.AngleDeg != nil {
			return "", errors.New("shadow.inherit cannot be combined with explicit shadow attributes")
		}
		if glow != nil || blur != nil || softEdge != nil || reflection != nil {
			return "", errors.New("shadow.inherit cannot be combined with other explicit effects")
		}
		if *shadow.Inherit {
			return "", nil
		}
		return `<a:effectLst/>`, nil
	}
	var items strings.Builder
	if shadow != nil {
		shadowXML, err := renderShadowXML(shadow)
		if err != nil {
			return "", err
		}
		items.WriteString(shadowXML)
	}
	if glow != nil {
		glowXML, err := renderGlowXML(glow)
		if err != nil {
			return "", err
		}
		items.WriteString(glowXML)
	}
	if blur != nil {
		blurXML, err := renderBlurXML(blur)
		if err != nil {
			return "", err
		}
		items.WriteString(blurXML)
	}
	if softEdge != nil {
		softEdgeXML, err := renderSoftEdgeXML(softEdge)
		if err != nil {
			return "", err
		}
		items.WriteString(softEdgeXML)
	}
	if reflection != nil {
		reflectionXML, err := renderReflectionXML(reflection)
		if err != nil {
			return "", err
		}
		items.WriteString(reflectionXML)
	}
	if items.Len() == 0 {
		return "", nil
	}
	return `<a:effectLst>` + items.String() + `</a:effectLst>`, nil
}

func renderFillXML(fill *common.ShapeFill) (string, error) {
	if fill == nil {
		return "", nil
	}
	modeCount := 0
	if fill.Solid != nil {
		modeCount++
	}
	if fill.Gradient != nil {
		modeCount++
	}
	if fill.Pattern != nil {
		modeCount++
	}
	if fill.Background != nil {
		modeCount++
	}
	if modeCount > 1 {
		return "", errors.New("fill.solid, fill.gradient, fill.pattern, and fill.background are mutually exclusive")
	}
	if fill.Solid != nil {
		color, err := normalizeHexColor(*fill.Solid)
		if err != nil {
			return "", fmt.Errorf("fill.solid: %w", err)
		}
		return `<a:solidFill><a:srgbClr val="` + color + `"/></a:solidFill>`, nil
	}
	if fill.Background != nil {
		if !*fill.Background {
			return "", errors.New("fill.background must be true when provided")
		}
		return `<a:noFill/>`, nil
	}
	if fill.Gradient != nil {
		return renderGradientFillXML(fill.Gradient)
	}
	if fill.Pattern != nil {
		return renderPatternFillXML(fill.Pattern)
	}
	return "", nil
}

func renderGradientFillXML(gradient *common.GradientFill) (string, error) {
	if gradient == nil {
		return "", nil
	}
	stops := gradient.Stops
	if len(stops) == 0 {
		return "", errors.New("fill.gradient.stops must contain at least 1 stop")
	}
	var b strings.Builder
	b.WriteString(`<a:gradFill><a:gsLst>`)
	for i := range stops {
		stop := stops[i]
		color, err := normalizeHexColor(stop.Color)
		if err != nil {
			return "", fmt.Errorf("fill.gradient.stops[%d].color: %w", i, err)
		}
		pos := 0.0
		if stop.PositionPct != nil {
			pos = *stop.PositionPct
		} else if len(stops) > 1 {
			pos = float64(i) * (100.0 / float64(len(stops)-1))
		}
		if pos < 0.0 || pos > 100.0 {
			return "", fmt.Errorf("fill.gradient.stops[%d].position_pct must be between 0 and 100", i)
		}
		b.WriteString(fmt.Sprintf(`<a:gs pos="%d"><a:srgbClr val="%s"/></a:gs>`, int(math.Round(pos*1000.0)), color))
	}
	b.WriteString(`</a:gsLst>`)
	if gradient.AngleDeg != nil {
		rotation, err := normalizeTextFrameRotation(*gradient.AngleDeg)
		if err != nil {
			return "", fmt.Errorf("fill.gradient.angle_deg: %w", err)
		}
		b.WriteString(fmt.Sprintf(`<a:lin ang="%d" scaled="1"/>`, rotation))
	}
	b.WriteString(`</a:gradFill>`)
	return b.String(), nil
}

func renderPatternFillXML(pattern *common.PatternedFill) (string, error) {
	if pattern == nil {
		return "", nil
	}
	prst := "pct5"
	if pattern.Preset != nil && strings.TrimSpace(*pattern.Preset) != "" {
		prst = strings.TrimSpace(*pattern.Preset)
	}
	fg := "000000"
	if pattern.FgColor != nil {
		color, err := normalizeHexColor(*pattern.FgColor)
		if err != nil {
			return "", fmt.Errorf("fill.pattern.fg_color: %w", err)
		}
		fg = color
	}
	bg := "FFFFFF"
	if pattern.BgColor != nil {
		color, err := normalizeHexColor(*pattern.BgColor)
		if err != nil {
			return "", fmt.Errorf("fill.pattern.bg_color: %w", err)
		}
		bg = color
	}
	return fmt.Sprintf(
		`<a:pattFill prst="%s"><a:fgClr><a:srgbClr val="%s"/></a:fgClr><a:bgClr><a:srgbClr val="%s"/></a:bgClr></a:pattFill>`,
		xmlEscape(prst),
		fg,
		bg,
	), nil
}

func renderShadowXML(shadow *common.ShapeShadow) (string, error) {
	if shadow == nil {
		return "", nil
	}
	color := "000000"
	if shadow.Color != nil {
		normalized, err := normalizeHexColor(*shadow.Color)
		if err != nil {
			return "", fmt.Errorf("shadow.color: %w", err)
		}
		color = normalized
	}
	blur := 50800
	if shadow.BlurEmu != nil {
		if *shadow.BlurEmu < 0 {
			return "", errors.New("shadow.blur_emu must be >= 0")
		}
		blur = *shadow.BlurEmu
	}
	dist := 38100
	if shadow.DistanceEmu != nil {
		if *shadow.DistanceEmu < 0 {
			return "", errors.New("shadow.distance_emu must be >= 0")
		}
		dist = *shadow.DistanceEmu
	}
	dir := int64(0)
	if shadow.AngleDeg != nil {
		rotation, err := normalizeTextFrameRotation(*shadow.AngleDeg)
		if err != nil {
			return "", fmt.Errorf("shadow.angle_deg: %w", err)
		}
		dir = rotation
	}
	return fmt.Sprintf(
		`<a:outerShdw blurRad="%d" dist="%d" dir="%d"><a:srgbClr val="%s"/></a:outerShdw>`,
		blur,
		dist,
		dir,
		color,
	), nil
}

func renderGlowXML(glow *common.ShapeGlow) (string, error) {
	if glow == nil {
		return "", nil
	}
	color := "000000"
	if glow.Color != nil {
		normalized, err := normalizeHexColor(*glow.Color)
		if err != nil {
			return "", fmt.Errorf("glow.color: %w", err)
		}
		color = normalized
	}
	radius := 38100
	if glow.RadiusEmu != nil {
		if *glow.RadiusEmu < 0 {
			return "", errors.New("glow.radius_emu must be >= 0")
		}
		radius = *glow.RadiusEmu
	}
	return fmt.Sprintf(`<a:glow rad="%d"><a:srgbClr val="%s"/></a:glow>`, radius, color), nil
}

func renderBlurXML(blur *common.ShapeBlur) (string, error) {
	if blur == nil {
		return "", nil
	}
	radius := 50800
	if blur.RadiusEmu != nil {
		if *blur.RadiusEmu < 0 {
			return "", errors.New("blur.radius_emu must be >= 0")
		}
		radius = *blur.RadiusEmu
	}
	return fmt.Sprintf(`<a:blur rad="%d"/>`, radius), nil
}

func renderSoftEdgeXML(softEdge *common.ShapeSoftEdge) (string, error) {
	if softEdge == nil {
		return "", nil
	}
	radius := 50800
	if softEdge.RadiusEmu != nil {
		if *softEdge.RadiusEmu < 0 {
			return "", errors.New("soft_edge.radius_emu must be >= 0")
		}
		radius = *softEdge.RadiusEmu
	}
	return fmt.Sprintf(`<a:softEdge rad="%d"/>`, radius), nil
}

func renderReflectionXML(reflection *common.ShapeReflection) (string, error) {
	if reflection == nil {
		return "", nil
	}
	blur := 0
	if reflection.BlurEmu != nil {
		if *reflection.BlurEmu < 0 {
			return "", errors.New("reflection.blur_emu must be >= 0")
		}
		blur = *reflection.BlurEmu
	}
	dist := 0
	if reflection.DistanceEmu != nil {
		if *reflection.DistanceEmu < 0 {
			return "", errors.New("reflection.distance_emu must be >= 0")
		}
		dist = *reflection.DistanceEmu
	}
	return fmt.Sprintf(`<a:reflection blurRad="%d" dist="%d"/>`, blur, dist), nil
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

			if updates.Text != nil || updates.Runs != nil || updates.TextFrame != nil || updates.Paragraph != nil {
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
				if updates.Paragraph != nil {
					s.Paragraph = updates.Paragraph
				}
				updatedXML, updateErr = replaceShapeTextBody(e, partPath, updatedXML, s)
				if updateErr != nil {
					return nil, false
				}
			}
			if updates.Fill != nil || updates.Line != nil || updates.Shadow != nil || updates.Glow != nil || updates.Blur != nil || updates.SoftEdge != nil || updates.Reflection != nil {
				replace = true
				if updates.Fill != nil {
					s.Fill = updates.Fill
				}
				if updates.Line != nil {
					s.Line = updates.Line
				}
				if updates.Shadow != nil {
					s.Shadow = updates.Shadow
				}
				if updates.Glow != nil {
					s.Glow = updates.Glow
				}
				if updates.Blur != nil {
					s.Blur = updates.Blur
				}
				if updates.SoftEdge != nil {
					s.SoftEdge = updates.SoftEdge
				}
				if updates.Reflection != nil {
					s.Reflection = updates.Reflection
				}
				updatedXML, updateErr = replaceShapeStyle(
					updatedXML,
					s.Fill,
					s.Line,
					s.Shadow,
					s.Glow,
					s.Blur,
					s.SoftEdge,
					s.Reflection,
				)
				if updateErr != nil {
					return nil, false
				}
			}
			if updates.ClickAction != nil || updates.HoverAction != nil {
				replace = true
				if updates.ClickAction != nil {
					s.ClickAction = updates.ClickAction
				}
				if updates.HoverAction != nil {
					s.HoverAction = updates.HoverAction
				}
				updatedXML, updateErr = replaceShapeActions(e, partPath, updatedXML, updates.ClickAction, updates.HoverAction)
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

func normalizeHexColor(raw string) (string, error) {
	color := strings.TrimSpace(strings.TrimPrefix(raw, "#"))
	if len(color) != 6 {
		return "", fmt.Errorf("expected 6 hex digits, got %q", raw)
	}
	for _, ch := range color {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
			return "", fmt.Errorf("expected hex color, got %q", raw)
		}
	}
	return strings.ToUpper(color), nil
}

func normalizeLineDashStyle(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", errors.New("must be non-empty")
	}
	switch s {
	case "solid", "dash", "dashDot", "lgDash", "lgDashDot", "lgDashDotDot", "sysDot", "sysDash",
		"sysDashDot", "sysDashDotDot":
		return s, nil
	}
	key := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(s, "-", "_"), " ", "_"))
	key = strings.ReplaceAll(key, "__", "_")
	aliases := map[string]string{
		"dash_dot":          "dashDot",
		"dashdot":           "dashDot",
		"dash_dot_dot":      "lgDashDotDot",
		"dashdotdot":        "lgDashDotDot",
		"long_dash":         "lgDash",
		"longdash":          "lgDash",
		"long_dash_dot":     "lgDashDot",
		"longdashdot":       "lgDashDot",
		"long_dash_dot_dot": "lgDashDotDot",
		"longdashdotdot":    "lgDashDotDot",
		"round_dot":         "sysDot",
		"rounddot":          "sysDot",
		"square_dot":        "sysDash",
		"squaredot":         "sysDash",
		"sys_dash":          "sysDash",
		"sysdash":           "sysDash",
		"sys_dot":           "sysDot",
		"sysdot":            "sysDot",
		"sys_dash_dot":      "sysDashDot",
		"sysdashdot":        "sysDashDot",
		"sys_dash_dot_dot":  "sysDashDotDot",
		"sysdashdotdot":     "sysDashDotDot",
	}
	if normalized, ok := aliases[key]; ok {
		return normalized, nil
	}
	return "", fmt.Errorf("unsupported value %q", raw)
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
	return e.buildActionXML(partPath, hl, "hlinkClick")
}

func (e *PresentationEditor) buildHoverActionXML(partPath string, hl *common.Hyperlink) (string, error) {
	return e.buildActionXML(partPath, hl, "hlinkMouseOver")
}

func (e *PresentationEditor) buildActionXML(partPath string, hl *common.Hyperlink, tag string) (string, error) {
	if hl == nil || partPath == "" {
		return "", nil
	}
	if err := validateHyperlinkAction(hl); err != nil {
		return "", err
	}

	action := strings.TrimSpace(getStr(hl.Action))
	if action == "" {
		action = deriveActionURL(hl)
	}

	attrs := make([]string, 0, 3)
	if hl.Address != nil && *hl.Address != "" {
		relID, err := e.getOrCreateHyperlinkRelID(partPath, *hl.Address)
		if err != nil {
			return "", fmt.Errorf("allocate hyperlink relationship id: %w", err)
		}
		attrs = append(attrs, fmt.Sprintf(`r:id="%s"`, xmlEscape(relID)))
	} else if hl.TargetSlide != nil {
		relID, err := e.getOrCreateSlideJumpRelID(partPath, *hl.TargetSlide)
		if err != nil {
			return "", err
		}
		attrs = append(attrs, fmt.Sprintf(`r:id="%s"`, xmlEscape(relID)))
	}
	if action != "" {
		attrs = append(attrs, fmt.Sprintf(`action="%s"`, xmlEscape(action)))
	}
	if tooltip := strings.TrimSpace(getStr(hl.Tooltip)); tooltip != "" {
		attrs = append(attrs, fmt.Sprintf(`tooltip="%s"`, xmlEscape(tooltip)))
	}
	if len(attrs) == 0 {
		return "", nil
	}
	return fmt.Sprintf(`<a:%s %s/>`, tag, strings.Join(attrs, " ")), nil
}

func deriveActionURL(hl *common.Hyperlink) string {
	if hl == nil {
		return ""
	}
	if hl.TargetSlide != nil {
		return "ppaction://hlinksldjump"
	}
	if hl.TargetJump != nil && *hl.TargetJump != "" {
		return "ppaction://hlinkshowjump?jump=" + strings.TrimSpace(*hl.TargetJump)
	}
	if hl.Macro != nil && *hl.Macro != "" {
		return "ppaction://macro?name=" + strings.TrimSpace(*hl.Macro)
	}
	return ""
}

func validateHyperlinkAction(hl *common.Hyperlink) error {
	if hl == nil {
		return nil
	}
	selectorCount := 0
	hasAddress := strings.TrimSpace(getStr(hl.Address)) != ""
	hasTargetSlide := hl.TargetSlide != nil
	hasJump := strings.TrimSpace(getStr(hl.TargetJump)) != ""
	hasMacro := strings.TrimSpace(getStr(hl.Macro)) != ""
	if hasAddress {
		selectorCount++
	}
	if hasTargetSlide {
		selectorCount++
	}
	if hasJump {
		selectorCount++
	}
	if hasMacro {
		selectorCount++
	}
	if selectorCount > 1 {
		return errors.New("hyperlink selectors are mutually exclusive: use only one of address, target_slide, jump, or macro")
	}
	if hasJump {
		jump := strings.ToLower(strings.TrimSpace(*hl.TargetJump))
		switch jump {
		case "nextslide", "previousslide", "firstslide", "lastslide", "lastslideviewed", "endshow":
			return nil
		default:
			return fmt.Errorf("unsupported jump target %q", *hl.TargetJump)
		}
	}
	return nil
}

func (e *PresentationEditor) getOrCreateSlideJumpRelID(partPath string, targetSlideIndex int) (string, error) {
	if targetSlideIndex < 0 || targetSlideIndex >= len(e.slides) {
		return "", fmt.Errorf("target_slide index %d out of range", targetSlideIndex)
	}
	targetPart := e.slides[targetSlideIndex].Part
	relsPath := common.SlideRelsPartName(partPath)
	rels := make([]common.EditorRelationship, 0)
	if data, ok := e.parts.Get(relsPath); ok {
		parsed, err := parseRelationshipsXML(data)
		if err != nil {
			return "", fmt.Errorf("parse %s: %w", relsPath, err)
		}
		rels = parsed
	}
	sourceDir := filepath.Dir(partPath)
	targetRelPath, err := filepath.Rel(sourceDir, targetPart)
	if err != nil {
		return "", fmt.Errorf("resolve target slide relationship path: %w", err)
	}
	targetRelPath = strings.ReplaceAll(targetRelPath, "\\", "/")
	for _, r := range rels {
		if r.Type == common.RelTypeSlide && r.Target == targetRelPath {
			return r.ID, nil
		}
	}
	relID, err := e.nextSlideRelID(partPath)
	if err != nil {
		return "", err
	}
	if err := e.addRelationship(partPath, relID, common.RelTypeSlide, targetRelPath); err != nil {
		return "", err
	}
	return relID, nil
}

// updateShapeTransforms performs a surgical regular expression replacement of shape transforms.
func updateShapeTransforms(xmlData []byte, x, y, w, h int) []byte {
	// Robustly match <a:off> and <a:ext> tags regardless of attribute order or spacing
	offRe := regexp.MustCompile(`(?s)<a:off\b[^>]*/>`)
	extRe := regexp.MustCompile(`(?s)<a:ext\b[^>]*/>`)

	res := offRe.ReplaceAllFunc(xmlData, func(b []byte) []byte {
		return []byte(fmt.Sprintf(`<a:off x="%d" y="%d"/>`, x, y))
	})
	res = extRe.ReplaceAllFunc(res, func(b []byte) []byte {
		return []byte(fmt.Sprintf(`<a:ext cx="%d" cy="%d"/>`, w, h))
	})
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
	return replaceShapeActions(e, partPath, xmlData, clickAction, nil)
}

func replaceShapeStyle(
	xmlData []byte,
	fill *common.ShapeFill,
	line *common.ShapeLine,
	shadow *common.ShapeShadow,
	glow *common.ShapeGlow,
	blur *common.ShapeBlur,
	softEdge *common.ShapeSoftEdge,
	reflection *common.ShapeReflection,
) ([]byte, error) {
	styleXML, err := renderShapeStyleXML(fill, line, shadow, glow, blur, softEdge, reflection)
	if err != nil {
		return nil, err
	}
	spPrRe := regexp.MustCompile(`(?s)<p:spPr\b([^>]*)>(.*?)</p:spPr>`)
	match := spPrRe.FindSubmatchIndex(xmlData)
	if match == nil {
		return xmlData, nil
	}
	inner := string(xmlData[match[4]:match[5]])
	solidPattern := regexp.MustCompile(`(?s)<a:solidFill\b[^>]*>.*?</a:solidFill>|<a:solidFill\b[^>]*/>`)
	noFillPattern := regexp.MustCompile(`(?s)<a:noFill\b[^>]*/>|<a:noFill\b[^>]*>.*?</a:noFill>`)
	gradPattern := regexp.MustCompile(`(?s)<a:gradFill\b[^>]*>.*?</a:gradFill>`)
	patternPattern := regexp.MustCompile(`(?s)<a:pattFill\b[^>]*>.*?</a:pattFill>|<a:pattFill\b[^>]*/>`)
	linePattern := regexp.MustCompile(`(?s)<a:ln\b[^>]*>.*?</a:ln>|<a:ln\b[^>]*/>`)
	effectPattern := regexp.MustCompile(`(?s)<a:effectLst\b[^>]*>.*?</a:effectLst>|<a:effectLst\b[^>]*/>`)
	inner = solidPattern.ReplaceAllString(inner, "")
	inner = noFillPattern.ReplaceAllString(inner, "")
	inner = gradPattern.ReplaceAllString(inner, "")
	inner = patternPattern.ReplaceAllString(inner, "")
	inner = linePattern.ReplaceAllString(inner, "")
	inner = effectPattern.ReplaceAllString(inner, "")
	if styleXML != "" {
		if idx := strings.Index(inner, "<a:prstGeom"); idx >= 0 {
			inner = inner[:idx] + styleXML + inner[idx:]
		} else {
			inner = styleXML + inner
		}
	}
	replacement := fmt.Sprintf(`<p:spPr%s>%s</p:spPr>`, string(xmlData[match[2]:match[3]]), inner)
	updated := string(xmlData[:match[0]]) + replacement + string(xmlData[match[1]:])
	return []byte(updated), nil
}

func replaceShapeActions(
	e *PresentationEditor,
	partPath string,
	xmlData []byte,
	clickAction *common.Hyperlink,
	hoverAction *common.Hyperlink,
) ([]byte, error) {
	clickXML, err := e.buildClickActionXML(partPath, clickAction)
	if err != nil {
		return nil, err
	}
	hoverXML, err := e.buildHoverActionXML(partPath, hoverAction)
	if err != nil {
		return nil, err
	}

	xmlStr := string(xmlData)
	hlinkClickPattern := regexp.MustCompile(`(?s)<a:hlinkClick\b[^>]*/>|<a:hlinkClick\b[^>]*>.*?</a:hlinkClick>`)
	hlinkHoverPattern := regexp.MustCompile(`(?s)<a:hlinkMouseOver\b[^>]*/>|<a:hlinkMouseOver\b[^>]*>.*?</a:hlinkMouseOver>`)
	cNvPrOpenClose := regexp.MustCompile(`(?s)<p:cNvPr\b([^>]*)>(.*?)</p:cNvPr>`)
	if match := cNvPrOpenClose.FindStringSubmatchIndex(xmlStr); match != nil {
		inner := xmlStr[match[4]:match[5]]
		removeTag := func(input string, pattern *regexp.Regexp) string {
			matches := pattern.FindAllStringIndex(input, -1)
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
		cleanInner := inner
		if clickAction != nil {
			cleanInner = removeTag(cleanInner, hlinkClickPattern)
		}
		if hoverAction != nil {
			cleanInner = removeTag(cleanInner, hlinkHoverPattern)
		}
		var replacement strings.Builder
		replacement.WriteString(cleanInner)
		if clickXML != "" {
			replacement.WriteString(clickXML)
		}
		if hoverXML != "" {
			replacement.WriteString(hoverXML)
		}
		updated := xmlStr[:match[4]] + replacement.String() + xmlStr[match[5]:]
		return []byte(updated), nil
	}

	cNvPrSelfClosing := regexp.MustCompile(`<p:cNvPr\b([^>]*)/>`)
	if match := cNvPrSelfClosing.FindStringSubmatchIndex(xmlStr); match != nil {
		if clickXML == "" && hoverXML == "" {
			return xmlData, nil
		}
		attrs := xmlStr[match[2]:match[3]]
		replacement := fmt.Sprintf(`<p:cNvPr%s>%s%s</p:cNvPr>`, attrs, clickXML, hoverXML)
		updated := xmlStr[:match[0]] + replacement + xmlStr[match[1]:]
		return []byte(updated), nil
	}

	if clickAction != nil || hoverAction != nil {
		return nil, errors.New("shape has no cNvPr node for action update")
	}

	return xmlData, nil
}
