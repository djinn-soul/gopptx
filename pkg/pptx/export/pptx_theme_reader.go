package export

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// PowerPoint gives most shapes their look by reference rather than by value: a
// shape drawn in the UI carries a <p:style> naming a slot in the theme's format
// scheme and a scheme colour, and no fill or line of its own. Without the theme
// those shapes render as bare text on the page — the fill, the outline and the
// white caption colour all live in theme1.xml.
//
// deckTheme is the part of a theme this renderer needs: the colour scheme, and
// the line widths of the format scheme.
type deckTheme struct {
	// colors maps a scheme slot ("accent1", "lt1", …) to its RGB hex, without
	// the leading hash.
	colors map[string]string
	// lineWidthsEMU are the widths of the format scheme's line styles, in the
	// order the theme lists them. A lnRef idx is one-based into this.
	lineWidthsEMU []int
}

// themeColor resolves a scheme slot to RGB hex, falling back to the built-in
// Office palette for a theme that does not name it.
func (t deckTheme) themeColor(slot string) (string, bool) {
	if t.colors != nil {
		if hexValue, ok := t.colors[strings.ToLower(slot)]; ok && hexValue != "" {
			return hexValue, true
		}
	}
	if c, ok := resolveThemeBaseColor(resolveColorAlias(normalizeColorName(slot))); ok {
		return rgbToHex(c), true
	}
	return "", false
}

// lineWidthEMU is the width of the format scheme's idx-th line style, one-based
// as lnRef states it.
func (t deckTheme) lineWidthEMU(idx int) int {
	if idx >= 1 && idx <= len(t.lineWidthsEMU) && t.lineWidthsEMU[idx-1] > 0 {
		return t.lineWidthsEMU[idx-1]
	}
	return defaultThemeLineWidthEMU
}

// defaultThemeLineWidthEMU is the 0.75pt line PowerPoint falls back to when the
// theme states no width for the slot.
const defaultThemeLineWidthEMU = 9525

// readDeckTheme reads the first theme part of a package. A deck can carry one
// theme per master; they share a format scheme in every template Office ships,
// and the first is the one the slides in a single-master deck use.
func readDeckTheme(pptxPath string) deckTheme {
	theme := deckTheme{colors: map[string]string{}}
	zr, err := zip.OpenReader(pptxPath)
	if err != nil {
		return theme
	}
	defer zr.Close()

	fileMap := make(map[string]*zip.File, len(zr.File))
	names := make([]string, 0, len(zr.File))
	for _, f := range zr.File {
		name := canonicalZipPath(f.Name)
		fileMap[name] = f
		if strings.HasPrefix(name, "ppt/theme/theme") && strings.HasSuffix(name, ".xml") {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return theme
	}
	sort.Strings(names)
	data := readZipBytes(fileMap, names[0])
	if data == nil {
		return theme
	}
	parseThemeXML(data, &theme)
	return theme
}

// themeColorSlot is the colour inside one <a:clrScheme> child. Which slot it
// fills is the name of the element that holds it, so the caller passes that in
// rather than the decoder capturing it.
type themeColorSlot struct {
	SRGB *themeSRGBColor   `xml:"srgbClr"`
	Sys  *themeSystemColor `xml:"sysClr"`
}

type themeSRGBColor struct {
	Val string `xml:"val,attr"`
}

type themeSystemColor struct {
	LastClr string `xml:"lastClr,attr"`
}

// hex is the slot's colour, or "" when it states none this renderer understands.
func (s themeColorSlot) hex() string {
	switch {
	case s.SRGB != nil && s.SRGB.Val != "":
		return strings.ToUpper(strings.TrimPrefix(s.SRGB.Val, "#"))
	case s.Sys != nil && s.Sys.LastClr != "":
		return strings.ToUpper(s.Sys.LastClr)
	default:
		return ""
	}
}

// parseThemeColorScheme reads <a:clrScheme> by walking tokens, because the slot
// name is the element name: dk1, lt1, accent1 and so on each hold one colour.
func parseThemeColorScheme(data []byte, theme *deckTheme) {
	dec := xml.NewDecoder(bytes.NewReader(data))
	inScheme := false
	for {
		token, err := dec.Token()
		if err != nil {
			return
		}
		switch element := token.(type) {
		case xml.StartElement:
			if element.Name.Local == "clrScheme" {
				inScheme = true
				continue
			}
			if !inScheme {
				continue
			}
			var slot themeColorSlot
			if decErr := dec.DecodeElement(&slot, &element); decErr != nil {
				continue
			}
			if hexValue := slot.hex(); hexValue != "" {
				theme.colors[strings.ToLower(element.Name.Local)] = hexValue
			}
		case xml.EndElement:
			if element.Name.Local == "clrScheme" {
				return
			}
		}
	}
}

// themeXML is the part of theme1.xml this renderer reads with the struct
// decoder; the colour scheme is walked separately, in parseThemeColorScheme.
type themeXML struct {
	XMLName    xml.Name            `xml:"theme"`
	LineStyles []themeLineStyleXML `xml:"themeElements>fmtScheme>lnStyleLst>ln"`
}

// themeLineStyleXML is one <a:ln> of the format scheme's line styles.
type themeLineStyleXML struct {
	W string `xml:"w,attr"`
}

func parseThemeXML(data []byte, theme *deckTheme) {
	var doc themeXML
	if err := xml.Unmarshal(data, &doc); err != nil {
		return
	}
	parseThemeColorScheme(data, theme)
	for _, line := range doc.LineStyles {
		width, err := strconv.Atoi(strings.TrimSpace(line.W))
		if err != nil {
			width = 0
		}
		theme.lineWidthsEMU = append(theme.lineWidthsEMU, width)
	}
}

func rgbToHex(c rgbColor) string {
	return fmt.Sprintf("%02X%02X%02X", c.r, c.g, c.b)
}
