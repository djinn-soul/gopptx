package export

import (
	"archive/zip"
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// This file reads one layout or master part: the placeholders it positions, the
// background it paints and the shapes it draws. What a slide then takes from
// them lives in pptx_reader_inherit.go.

type layoutPartContent struct {
	placeholders map[placeholderKey]placeholderBox
	background   *elements.SlideBackground
	shapes       []shapes.Shape
	titleSizePt  int
	bodySizePt   int
}

// layoutShapeXML is a <p:sp> of a layout or master: enough of it to know
// whether it is a placeholder, where it sits, and what it draws.
type layoutShapeXML struct {
	NvSpPr struct {
		CNvPr struct {
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
		NvPr struct {
			Ph *struct {
				Type string `xml:"type,attr"`
				Idx  string `xml:"idx,attr"`
			} `xml:"ph"`
		} `xml:"nvPr"`
	} `xml:"nvSpPr"`
	SpPr struct {
		Xfrm *struct {
			Off struct {
				X int64 `xml:"x,attr"`
				Y int64 `xml:"y,attr"`
			} `xml:"off"`
			Ext struct {
				Cx int64 `xml:"cx,attr"`
				Cy int64 `xml:"cy,attr"`
			} `xml:"ext"`
		} `xml:"xfrm"`
		PrstGeom *struct {
			Prst string `xml:"prst,attr"`
		} `xml:"prstGeom"`
		SolidFill *fillColorXML `xml:"solidFill"`
		NoFill    *struct{}     `xml:"noFill"`
		Ln        *struct {
			W         string        `xml:"w,attr"`
			SolidFill *fillColorXML `xml:"solidFill"`
			NoFill    *struct{}     `xml:"noFill"`
		} `xml:"ln"`
	} `xml:"spPr"`
	Style *struct {
		FillRef styleColorRefXML `xml:"fillRef"`
		LnRef   styleColorRefXML `xml:"lnRef"`
		FontRef styleColorRefXML `xml:"fontRef"`
	} `xml:"style"`
	TxBody *struct {
		BodyPr struct {
			Anchor string `xml:"anchor,attr"`
		} `xml:"bodyPr"`
		Paragraphs []struct {
			Runs []struct {
				RPr *struct {
					Sz string `xml:"sz,attr"`
				} `xml:"rPr"`
				T string `xml:"t"`
			} `xml:"r"`
		} `xml:"p"`
	} `xml:"txBody"`
}

// fillColorXML is a solid fill's colour, by value or by scheme slot.
type fillColorXML struct {
	SrgbClr *struct {
		Val string `xml:"val,attr"`
	} `xml:"srgbClr"`
	SchemeClr *struct {
		Val string `xml:"val,attr"`
	} `xml:"schemeClr"`
}

func (f *fillColorXML) hex(theme deckTheme) string {
	if f == nil {
		return ""
	}
	if f.SrgbClr != nil && f.SrgbClr.Val != "" {
		return strings.ToUpper(strings.TrimPrefix(f.SrgbClr.Val, "#"))
	}
	if f.SchemeClr != nil && f.SchemeClr.Val != "" {
		if hexValue, ok := theme.themeColor(resolveColorAlias(normalizeColorName(f.SchemeClr.Val))); ok {
			return hexValue
		}
	}
	return ""
}

func parseLayoutPart(data []byte, theme deckTheme) layoutPartContent {
	out := layoutPartContent{placeholders: map[placeholderKey]placeholderBox{}}
	var doc struct {
		Bg *struct {
			BgPr *struct {
				SolidFill *fillColorXML `xml:"solidFill"`
			} `xml:"bgPr"`
			BgRef *struct {
				SchemeClr *struct {
					Val string `xml:"val,attr"`
				} `xml:"schemeClr"`
			} `xml:"bgRef"`
		} `xml:"cSld>bg"`
		Shapes    []layoutShapeXML `xml:"cSld>spTree>sp"`
		TitleSize struct {
			Sz string `xml:"sz,attr"`
		} `xml:"txStyles>titleStyle>lvl1pPr>defRPr"`
		BodySize struct {
			Sz string `xml:"sz,attr"`
		} `xml:"txStyles>bodyStyle>lvl1pPr>defRPr"`
	}
	if err := xml.Unmarshal(data, &doc); err != nil {
		return out
	}
	out.titleSizePt = sizeAttrToPoints(doc.TitleSize.Sz)
	out.bodySizePt = sizeAttrToPoints(doc.BodySize.Sz)
	if doc.Bg != nil {
		out.background = layoutBackground(doc.Bg.BgPr, doc.Bg.BgRef, theme)
	}
	for i := range doc.Shapes {
		shape := doc.Shapes[i]
		if ph := shape.NvSpPr.NvPr.Ph; ph != nil {
			if box, ok := layoutPlaceholderBox(shape); ok {
				out.placeholders[placeholderKey{
					phType: strings.ToLower(strings.TrimSpace(ph.Type)),
					idx:    atoiOrZero(ph.Idx),
				}] = box
			}
			// A placeholder on a layout is a prompt, not something the slide
			// shows, so it is never painted.
			continue
		}
		if drawn, ok := layoutShapeToShape(shape, theme); ok {
			out.shapes = append(out.shapes, drawn)
		}
	}
	return out
}

func layoutBackground(
	bgPr *struct {
		SolidFill *fillColorXML `xml:"solidFill"`
	},
	bgRef *struct {
		SchemeClr *struct {
			Val string `xml:"val,attr"`
		} `xml:"schemeClr"`
	},
	theme deckTheme,
) *elements.SlideBackground {
	if bgPr != nil && bgPr.SolidFill != nil {
		if hexValue := bgPr.SolidFill.hex(theme); hexValue != "" {
			bg := elements.NewSolidBackground(hexValue)
			return &bg
		}
	}
	if bgRef != nil && bgRef.SchemeClr != nil {
		if hexValue, ok := theme.themeColor(resolveColorAlias(normalizeColorName(bgRef.SchemeClr.Val))); ok {
			bg := elements.NewSolidBackground(hexValue)
			return &bg
		}
	}
	return nil
}

func layoutPlaceholderBox(shape layoutShapeXML) (placeholderBox, bool) {
	if shape.SpPr.Xfrm == nil {
		return placeholderBox{}, false
	}
	box := placeholderBox{
		X:  shape.SpPr.Xfrm.Off.X,
		Y:  shape.SpPr.Xfrm.Off.Y,
		CX: shape.SpPr.Xfrm.Ext.Cx,
		CY: shape.SpPr.Xfrm.Ext.Cy,
	}
	if box.CX <= 0 || box.CY <= 0 {
		return placeholderBox{}, false
	}
	box.SizePt = layoutShapeFirstRunSize(shape)
	return box, true
}

// sizeAttrToPoints converts an OOXML sz attribute (hundredths of a point) to
// points, or zero when it is absent or unreadable.
func sizeAttrToPoints(raw string) int {
	const centiToPoints = 100
	hundredths, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || hundredths <= 0 {
		return 0
	}
	return hundredths / centiToPoints
}

// layoutShapeFirstRunSize is the point size the layout states for the
// placeholder's first run, or zero when it states none.
func layoutShapeFirstRunSize(shape layoutShapeXML) int {
	if shape.TxBody == nil {
		return 0
	}
	for _, paragraph := range shape.TxBody.Paragraphs {
		for _, run := range paragraph.Runs {
			if run.RPr == nil || run.RPr.Sz == "" {
				continue
			}
			if size := sizeAttrToPoints(run.RPr.Sz); size > 0 {
				return size
			}
		}
	}
	return 0
}

// layoutShapeToShape converts a layout or master shape into one the renderer
// can paint. Shapes with neither a fill nor text are skipped: they are almost
// always invisible spacers.
func layoutShapeToShape(shape layoutShapeXML, theme deckTheme) (shapes.Shape, bool) {
	if shape.SpPr.Xfrm == nil || shape.SpPr.Xfrm.Ext.Cx <= 0 || shape.SpPr.Xfrm.Ext.Cy <= 0 {
		return shapes.Shape{}, false
	}
	preset := presetRect
	if shape.SpPr.PrstGeom != nil && shape.SpPr.PrstGeom.Prst != "" {
		preset = shape.SpPr.PrstGeom.Prst
	}
	out := shapes.NewShape(
		preset,
		styling.Emu(shape.SpPr.Xfrm.Off.X),
		styling.Emu(shape.SpPr.Xfrm.Off.Y),
		styling.Emu(shape.SpPr.Xfrm.Ext.Cx),
		styling.Emu(shape.SpPr.Xfrm.Ext.Cy),
	)
	out.Name = shape.NvSpPr.CNvPr.Name
	out.Text = layoutShapeText(shape)
	if shape.TxBody != nil {
		if anchor := strings.TrimSpace(shape.TxBody.BodyPr.Anchor); anchor != "" {
			out = out.WithVerticalAnchor(shapes.TextFrameAnchor(anchor))
		}
	}
	if shape.SpPr.NoFill == nil {
		if hexValue := shape.SpPr.SolidFill.hex(theme); hexValue != "" {
			fill := shapes.NewShapeFill(hexValue)
			out.Fill = &fill
		}
	}
	if line := shape.SpPr.Ln; line != nil && line.NoFill == nil {
		if hexValue := line.SolidFill.hex(theme); hexValue != "" {
			width := emuToPt(int64(atoiOrZero(line.W)))
			if width <= 0 {
				width = emuToPt(defaultThemeLineWidthEMU)
			}
			stroke := shapes.NewShapeLine(hexValue, styling.Points(width))
			out.Line = &stroke
		}
	}
	if shape.Style != nil {
		applyShapeThemeStyle(&out, resolveShapeThemeStyle(*shape.Style, theme, ""))
	}
	if out.Fill == nil && out.Line == nil && strings.TrimSpace(out.Text) == "" {
		return shapes.Shape{}, false
	}
	return out, true
}

func layoutShapeText(shape layoutShapeXML) string {
	if shape.TxBody == nil {
		return ""
	}
	var b strings.Builder
	for i, paragraph := range shape.TxBody.Paragraphs {
		if i > 0 {
			b.WriteString("\n")
		}
		for _, run := range paragraph.Runs {
			b.WriteString(run.T)
		}
	}
	return b.String()
}

// relatedPart resolves the first relationship of the given type from a part's
// own .rels, as a package-absolute path.
func relatedPart(fileMap map[string]*zip.File, part, relType string) string {
	data := readZipBytes(fileMap, slideRelsPath(part))
	if data == nil {
		return ""
	}
	var rels struct {
		Rels []struct {
			Type   string `xml:"Type,attr"`
			Target string `xml:"Target,attr"`
		} `xml:"Relationship"`
	}
	if err := xml.Unmarshal(data, &rels); err != nil {
		return ""
	}
	for _, rel := range rels.Rels {
		if strings.HasSuffix(rel.Type, "/"+relType) {
			return resolveRelPath(part, rel.Target)
		}
	}
	return ""
}
