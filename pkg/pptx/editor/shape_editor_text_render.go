package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"strings"

	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

// ptToEMU converts points to English Metric Units (1 pt = 12700 EMU).
const ptToEMU = 12700

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
	txBody.WriteString(
		`<p:txBody xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`,
	)

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
			anchor, err := normalizeTextFrameVerticalAlign(*tf.VerticalAlign)
			if err != nil {
				return nil, err
			}
			bodyPr += fmt.Sprintf(` anchor="%s"`, escape(anchor))
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
				bodyPr += `<a:normAutoFit/>`
			case "shape":
				bodyPr += `<a:spAutoFit/>`
			case bulletStyleNone:
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
				rPr += fmt.Sprintf(` u="%s"`, escape(normalizeUnderlineValue(*r.Underline)))
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
			} else if r.SmallCaps != nil && *r.SmallCaps {
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
			if r.OutlineColor != nil && *r.OutlineColor != "" {
				widthEMU := int64(ptToEMU) // default 1pt
				if r.OutlineWidthPt != nil && *r.OutlineWidthPt > 0 {
					widthEMU = int64(*r.OutlineWidthPt * ptToEMU)
				}
				rPr += fmt.Sprintf(`<a:ln w="%d"><a:solidFill><a:srgbClr val="%s"/></a:solidFill></a:ln>`,
					widthEMU, escape(*r.OutlineColor))
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
				txBody.WriteString(
					fmt.Sprintf(`<a:p><a:r><a:rPr lang="en-US"/><a:t>%s</a:t></a:r></a:p>`, escape(s.Text)),
				)
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

func normalizeTextFrameVerticalAlign(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "t", "top":
		return "t", nil
	case alignCtr, alignCenter, "middle":
		return alignCtr, nil
	case "b", "bottom":
		return "b", nil
	case alignJust, "justify":
		return alignJust, nil
	case alignDist, "distribute":
		return alignDist, nil
	default:
		return "", fmt.Errorf("unsupported text_frame.vertical_align %q", raw)
	}
}

func normalizeTextFrameRotation(raw float64) (int64, error) {
	if math.IsNaN(raw) || math.IsInf(raw, 0) {
		return 0, errors.New("text_frame.rotation must be finite")
	}
	if raw < -360.0 || raw > 360.0 {
		return 0, errors.New("text_frame.rotation must be between -360 and 360 degrees")
	}
	return int64(math.Round(raw * editorslide.RotationDegreeToOOXML)), nil
}

func normalizeUnderlineValue(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", bulletStyleNone:
		return "none"
	case "single":
		return "sng"
	case "double":
		return "dbl"
	default:
		return raw
	}
}
