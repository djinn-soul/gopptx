package export

import (
	"testing"
)

const testThemeXML = `<?xml version="1.0" encoding="UTF-8"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <a:themeElements>
    <a:clrScheme name="Test">
      <a:dk1><a:sysClr val="windowText" lastClr="000000"/></a:dk1>
      <a:lt1><a:sysClr val="window" lastClr="FFFFFF"/></a:lt1>
      <a:accent1><a:srgbClr val="156082"/></a:accent1>
    </a:clrScheme>
    <a:fmtScheme name="Test">
      <a:lnStyleLst>
        <a:ln w="6350"/>
        <a:ln w="12700"/>
        <a:ln w="19050"/>
      </a:lnStyleLst>
    </a:fmtScheme>
  </a:themeElements>
</a:theme>`

func testTheme(t *testing.T) deckTheme {
	t.Helper()
	theme := deckTheme{colors: map[string]string{}}
	parseThemeXML([]byte(testThemeXML), &theme)
	return theme
}

func TestParseThemeXMLReadsColorsAndLineWidths(t *testing.T) {
	t.Parallel()

	theme := testTheme(t)
	if got, ok := theme.themeColor(themeColorAccent1); !ok || got != "156082" {
		t.Fatalf("accent1 = %q (ok=%v), want 156082", got, ok)
	}
	if got, ok := theme.themeColor("lt1"); !ok || got != "FFFFFF" {
		t.Fatalf("lt1 = %q (ok=%v), want FFFFFF", got, ok)
	}
	if got := theme.lineWidthEMU(2); got != 12700 {
		t.Fatalf("line width idx 2 = %d, want 12700", got)
	}
	// An index the theme does not carry falls back rather than panicking.
	if got := theme.lineWidthEMU(9); got != defaultThemeLineWidthEMU {
		t.Fatalf("line width idx 9 = %d, want %d", got, defaultThemeLineWidthEMU)
	}
}

// A shape PowerPoint draws carries no fill or line of its own: both are
// references into the theme, and the caption colour comes from fontRef.
const testStyledSlideXML = `<?xml version="1.0" encoding="UTF-8"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
       xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld><p:spTree>
    <p:sp>
      <p:nvSpPr><p:cNvPr id="2" name="Rectangle 1"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>
      <p:spPr><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr>
      <p:style>
        <a:lnRef idx="2"><a:schemeClr val="accent1"><a:shade val="50000"/></a:schemeClr></a:lnRef>
        <a:fillRef idx="1"><a:schemeClr val="accent1"/></a:fillRef>
        <a:effectRef idx="0"><a:schemeClr val="accent1"/></a:effectRef>
        <a:fontRef idx="minor"><a:schemeClr val="lt1"/></a:fontRef>
      </p:style>
    </p:sp>
  </p:spTree></p:cSld>
</p:sld>`

func TestParseSlideShapeStylesResolvesThemeReferences(t *testing.T) {
	t.Parallel()

	styles := parseSlideShapeStyles([]byte(testStyledSlideXML), testTheme(t))
	style, ok := styles[2]
	if !ok {
		t.Fatal("shape 2 has no resolved style")
	}
	if style.FillHex != "156082" {
		t.Errorf("fill = %q, want 156082", style.FillHex)
	}
	// A 50% shade halves the channels of the accent colour (15 60 82 → 0A 30 41).
	if style.LineHex != "0A3041" {
		t.Errorf("line = %q, want 0A3041", style.LineHex)
	}
	if style.TextHex != "FFFFFF" {
		t.Errorf("caption colour = %q, want FFFFFF", style.TextHex)
	}
	if style.LineWidthPt <= 0 {
		t.Errorf("line width = %v, want the theme's idx-2 width", style.LineWidthPt)
	}
}

// A shape that fills itself keeps its own fill: the theme only answers what the
// shape leaves unsaid.
const testNoFillSlideXML = `<?xml version="1.0" encoding="UTF-8"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
       xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld><p:spTree>
    <p:sp>
      <p:nvSpPr><p:cNvPr id="3" name="Rectangle 2"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>
      <p:spPr><a:prstGeom prst="rect"><a:avLst/></a:prstGeom><a:noFill/></p:spPr>
      <p:style>
        <a:fillRef idx="1"><a:schemeClr val="accent1"/></a:fillRef>
        <a:fontRef idx="minor"><a:schemeClr val="lt1"/></a:fontRef>
      </p:style>
    </p:sp>
  </p:spTree></p:cSld>
</p:sld>`

func TestParseSlideShapeStylesHonoursNoFill(t *testing.T) {
	t.Parallel()

	styles := parseSlideShapeStyles([]byte(testNoFillSlideXML), testTheme(t))
	if style := styles[3]; style.FillHex != "" {
		t.Fatalf("fill = %q, want none for a shape with <a:noFill/>", style.FillHex)
	}
}
