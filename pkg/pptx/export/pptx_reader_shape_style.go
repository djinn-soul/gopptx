package export

import (
	"archive/zip"
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// shapeThemeStyle is what a shape's <p:style> resolves to once the theme has
// been applied: the colours and line width PowerPoint would draw it with when
// its <p:spPr> states none of its own.
type shapeThemeStyle struct {
	FillHex     string
	LineHex     string
	LineWidthPt float64
	TextHex     string
}

// styleColorRef is one <a:fillRef>/<a:lnRef>/<a:fontRef>: a slot index into the
// theme's format scheme plus the colour to draw that slot with.
type styleColorRef struct {
	Idx    int
	Scheme string
	SRGB   string
	Shade  int
	Tint   int
	LumMod int
	LumOff int
}

// extractSlideShapeStyles returns, per slide (0-based), the resolved theme style
// of every shape that carries a <p:style>, keyed by shape id.
func extractSlideShapeStyles(pptxPath string, theme deckTheme) []map[int]shapeThemeStyle {
	zr, err := zip.OpenReader(pptxPath)
	if err != nil {
		return nil
	}
	defer zr.Close()

	fileMap := make(map[string]*zip.File, len(zr.File))
	for _, f := range zr.File {
		fileMap[canonicalZipPath(f.Name)] = f
	}
	slideOrder := resolveSlideOrder(fileMap)

	out := make([]map[int]shapeThemeStyle, len(slideOrder))
	for i, slidePart := range slideOrder {
		data := readZipBytes(fileMap, slidePart)
		if data == nil {
			continue
		}
		out[i] = parseSlideShapeStyles(data, theme)
	}
	return out
}

// styledShape is one <p:sp> as far as its style reference goes.
type styledShape struct {
	NvSpPr struct {
		CNvPr struct {
			ID string `xml:"id,attr"`
		} `xml:"cNvPr"`
	} `xml:"nvSpPr"`
	SpPr struct {
		Raw []byte `xml:",innerxml"`
	} `xml:"spPr"`
	Style *struct {
		FillRef styleColorRefXML `xml:"fillRef"`
		LnRef   styleColorRefXML `xml:"lnRef"`
		FontRef styleColorRefXML `xml:"fontRef"`
	} `xml:"style"`
}

type styleColorRefXML struct {
	Idx       string `xml:"idx,attr"`
	SchemeClr *struct {
		Val   string `xml:"val,attr"`
		Shade *struct {
			Val string `xml:"val,attr"`
		} `xml:"shade"`
		Tint *struct {
			Val string `xml:"val,attr"`
		} `xml:"tint"`
		LumMod *struct {
			Val string `xml:"val,attr"`
		} `xml:"lumMod"`
		LumOff *struct {
			Val string `xml:"val,attr"`
		} `xml:"lumOff"`
	} `xml:"schemeClr"`
	SrgbClr *struct {
		Val string `xml:"val,attr"`
	} `xml:"srgbClr"`
}

func (r styleColorRefXML) ref() styleColorRef {
	out := styleColorRef{Idx: atoiOrZero(r.Idx)}
	if r.SrgbClr != nil {
		out.SRGB = strings.ToUpper(strings.TrimPrefix(r.SrgbClr.Val, "#"))
	}
	if r.SchemeClr != nil {
		out.Scheme = r.SchemeClr.Val
		if r.SchemeClr.Shade != nil {
			out.Shade = atoiOrZero(r.SchemeClr.Shade.Val)
		}
		if r.SchemeClr.Tint != nil {
			out.Tint = atoiOrZero(r.SchemeClr.Tint.Val)
		}
		if r.SchemeClr.LumMod != nil {
			out.LumMod = atoiOrZero(r.SchemeClr.LumMod.Val)
		}
		if r.SchemeClr.LumOff != nil {
			out.LumOff = atoiOrZero(r.SchemeClr.LumOff.Val)
		}
	}
	return out
}

// parseSlideShapeStyles walks every <p:sp> in the slide, groups included, and
// resolves the ones that carry a <p:style>.
func parseSlideShapeStyles(data []byte, theme deckTheme) map[int]shapeThemeStyle {
	styles := make(map[int]shapeThemeStyle)
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	for {
		token, err := dec.Token()
		if err != nil {
			return styles
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "sp" {
			continue
		}
		var shape styledShape
		if decErr := dec.DecodeElement(&shape, &start); decErr != nil {
			continue
		}
		if shape.Style == nil {
			continue
		}
		id, convErr := strconv.Atoi(strings.TrimSpace(shape.NvSpPr.CNvPr.ID))
		if convErr != nil {
			continue
		}
		styles[id] = resolveShapeThemeStyle(*shape.Style, theme, string(shape.SpPr.Raw))
	}
}

// resolveShapeThemeStyle turns a shape's style references into concrete colours.
//
// Only the references the shape's own <p:spPr> leaves unanswered matter, so the
// caller is told what the theme would draw and decides; spPrRaw is inspected
// solely to skip work when the shape fills and outlines itself.
func resolveShapeThemeStyle(
	style struct {
		FillRef styleColorRefXML `xml:"fillRef"`
		LnRef   styleColorRefXML `xml:"lnRef"`
		FontRef styleColorRefXML `xml:"fontRef"`
	},
	theme deckTheme,
	spPrRaw string,
) shapeThemeStyle {
	out := shapeThemeStyle{}
	fill := style.FillRef.ref()
	// A fillRef of idx 0 means "no fill", which is a statement in its own right
	// and not something to paint.
	if fill.Idx > 0 && !strings.Contains(spPrRaw, "<a:noFill") {
		out.FillHex = resolveStyleRefColor(fill, theme)
	}
	line := style.LnRef.ref()
	if line.Idx > 0 {
		out.LineHex = resolveStyleRefColor(line, theme)
		out.LineWidthPt = emuToPt(int64(theme.lineWidthEMU(line.Idx)))
	}
	out.TextHex = resolveStyleRefColor(style.FontRef.ref(), theme)
	return out
}

// resolveStyleRefColor applies a style reference's colour transforms to the
// theme colour it names.
func resolveStyleRefColor(ref styleColorRef, theme deckTheme) string {
	base := ref.SRGB
	if base == "" && ref.Scheme != "" {
		if resolved, ok := theme.themeColor(resolveColorAlias(normalizeColorName(ref.Scheme))); ok {
			base = resolved
		}
	}
	if base == "" {
		return ""
	}
	parsed, ok := parseHexRGB(base)
	if !ok {
		return ""
	}
	parsed = applyColorTransforms(parsed, colorTransforms{
		tint:   ref.Tint,
		shade:  ref.Shade,
		lumMod: ref.LumMod,
		lumOff: ref.LumOff,
	})
	return rgbToHex(parsed)
}

// applyShapeThemeStyle fills in the look a shape inherits from the theme,
// leaving anything the shape states for itself alone.
//
// A shape drawn in PowerPoint carries no fill or outline of its own — both come
// from its <p:style>, and without this it rendered as bare text with no card
// behind it and no border around it.
func applyShapeThemeStyle(shape *shapes.Shape, style shapeThemeStyle) {
	if shape.Fill == nil && shape.GradientFill == nil && shape.RichFill == nil && style.FillHex != "" {
		fill := shapes.NewShapeFill(style.FillHex)
		shape.Fill = &fill
	}
	if shape.Line == nil && shape.RichLine == nil && style.LineHex != "" {
		width := style.LineWidthPt
		if width <= 0 {
			width = emuToPt(defaultThemeLineWidthEMU)
		}
		line := shapes.NewShapeLine(style.LineHex, styling.Points(width))
		shape.Line = &line
	}
	if style.TextHex != "" {
		applyShapeThemeTextColor(shape, style.TextHex)
	}
}

// applyShapeThemeTextColor gives the shape's runs the caption colour its
// fontRef names, for the runs that state no colour of their own.
func applyShapeThemeTextColor(shape *shapes.Shape, hexColor string) {
	for i := range shape.TextParagraphs {
		for j := range shape.TextParagraphs[i].Runs {
			if shape.TextParagraphs[i].Runs[j].Color == "" {
				shape.TextParagraphs[i].Runs[j].Color = hexColor
			}
		}
	}
}

func atoiOrZero(raw string) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0
	}
	return value
}
