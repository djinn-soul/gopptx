package pptxxml

import (
	"fmt"
	"strings"
	"testing"
)

// --- Baseline reference implementations (pre-optimization) for delta comparison ---

func baselineEscape(value string) string {
	return xmlEscapeReplacer.Replace(value)
}

func baselineRichSolidFill(fill SolidFillSpec) string {
	alphaXML := ""
	if fill.Transparency > 0 {
		alphaVal := int((1.0 - fill.Transparency) * transparencyBase)
		alphaXML = fmt.Sprintf(`<a:alpha val="%d"/>`, alphaVal)
	}
	return fmt.Sprintf(`<a:solidFill><a:srgbClr val="%s">%s</a:srgbClr></a:solidFill>`,
		baselineEscape(fill.Color), alphaXML)
}

func baselineRichLine(line RichShapeLineSpec) string {
	attrs := fmt.Sprintf(`w="%d"`, line.Width)
	if line.CapStyle != "" {
		attrs += fmt.Sprintf(` cap="%s"`, string(line.CapStyle))
	}
	dashXML := ""
	if line.DashStyle != "" && line.DashStyle != LineDashStyleSolid {
		dashXML = fmt.Sprintf(`<a:prstDash val="%s"/>`, string(line.DashStyle))
	}
	joinXML := ""
	switch line.JoinStyle {
	case LineJoinStyleBevel:
		joinXML = `<a:bevel/>`
	case LineJoinStyleMiter:
		joinXML = `<a:miter/>`
	case LineJoinStyleRound:
		joinXML = `<a:round/>`
	}
	alphaXML := ""
	if line.Transparency > 0 {
		alphaVal := int((1.0 - line.Transparency) * transparencyBase)
		alphaXML = fmt.Sprintf(`<a:alpha val="%d"/>`, alphaVal)
	}
	return fmt.Sprintf(`<a:ln %s><a:solidFill><a:srgbClr val="%s">%s</a:srgbClr></a:solidFill>%s%s</a:ln>`,
		attrs, baselineEscape(line.Color), alphaXML, dashXML, joinXML)
}

func baselineOuterShadow(shadow RichShapeShadowSpec) string {
	attrs := fmt.Sprintf(`blurRad="%d" dist="%d" dir="%d"`,
		shadow.BlurRadius, shadow.Distance, shadowDirEMU(shadow.Angle))
	if shadow.Alignment != "" {
		attrs += fmt.Sprintf(` algn="%s"`, baselineEscape(shadow.Alignment))
	}
	if !shadow.RotateWithShape {
		attrs += ` rotWithShape="0"`
	}
	alphaVal := shadowAlphaValue(shadow.Transparency)
	return fmt.Sprintf(`<a:outerShdw %s><a:srgbClr val="%s"><a:alpha val="%d"/></a:srgbClr></a:outerShdw>`,
		attrs, baselineEscape(shadow.Color), alphaVal)
}

func baselinePerspectiveShadow(shadow RichShapeShadowSpec) string {
	attrs := fmt.Sprintf(`dist="%d" dir="%d"`, shadow.Distance, shadowDirEMU(shadow.Angle))
	if shadow.SkewX != 0 || shadow.SkewY != 0 {
		attrs += fmt.Sprintf(` sx="%d" sy="%d"`,
			int(shadow.SkewX*shadowScaleBase), int(shadow.SkewY*shadowScaleBase))
	}
	if shadow.ScaleX != 1.0 || shadow.ScaleY != 1.0 {
		attrs += fmt.Sprintf(` kx="%d" ky="%d"`,
			int(shadow.ScaleX*shadowScaleBase), int(shadow.ScaleY*shadowScaleBase))
	}
	if shadow.Alignment != "" {
		attrs += fmt.Sprintf(` algn="%s"`, baselineEscape(shadow.Alignment))
	}
	alphaVal := shadowAlphaValue(shadow.Transparency)
	return fmt.Sprintf(
		`<a:prstShdw prst="shdw1" %s><a:srgbClr val="%s"><a:alpha val="%d"/></a:srgbClr></a:prstShdw>`,
		attrs, baselineEscape(shadow.Color), alphaVal)
}

func baselineTextLevelStyles(levels []TextLevelStyle) string {
	var b strings.Builder
	for _, lvl := range levels {
		lvlNum := max(lvl.Level+1, 1)
		lvlNum = min(lvlNum, 9)
		attrs := ""
		if lvl.IndentEMU > 0 {
			attrs += fmt.Sprintf(` indent="%d"`, lvl.IndentEMU)
		}
		b.WriteString(fmt.Sprintf("\n<a:lvl%dpPr%s>", lvlNum, attrs))
		if lvl.BulletChar != "" {
			b.WriteString(fmt.Sprintf(`<a:buChar char="%s"/>`, baselineEscape(lvl.BulletChar)))
		}
		rprAttrs := ""
		if lvl.SizePt > 0 {
			rprAttrs += fmt.Sprintf(` sz="%d"`, lvl.SizePt*100)
		}
		if lvl.Bold {
			rprAttrs += ` b="1"`
		}
		if lvl.Italic {
			rprAttrs += ` i="1"`
		}
		b.WriteString(fmt.Sprintf(`<a:defRPr%s>`, rprAttrs))
		if lvl.Color != "" {
			color := strings.TrimPrefix(lvl.Color, "#")
			b.WriteString(fmt.Sprintf(`<a:solidFill><a:srgbClr val="%s"/></a:solidFill>`, color))
		}
		if lvl.Font != "" {
			b.WriteString(fmt.Sprintf(`<a:latin typeface="%s"/>`, baselineEscape(lvl.Font)))
		}
		b.WriteString(`</a:defRPr>`)
		b.WriteString(fmt.Sprintf("\n</a:lvl%dpPr>", lvlNum))
	}
	return b.String()
}

// BenchmarkRichSolidFill stresses the per-shape solidFill XML emission path.
func BenchmarkRichSolidFill(b *testing.B) {
	fill := SolidFillSpec{Color: "4472C4", Transparency: 0.2}
	b.ReportAllocs()
	for b.Loop() {
		_ = richSolidFillXML(fill)
	}
}

func BenchmarkRichSolidFillBaseline(b *testing.B) {
	fill := SolidFillSpec{Color: "4472C4", Transparency: 0.2}
	b.ReportAllocs()
	for b.Loop() {
		_ = baselineRichSolidFill(fill)
	}
}

// BenchmarkRichLine stresses the per-shape <a:ln> emission path.
func BenchmarkRichLine(b *testing.B) {
	line := RichShapeLineSpec{
		Color:        "1F3864",
		Width:        12700,
		DashStyle:    LineDashStyleDash,
		CapStyle:     LineCapStyleRound,
		JoinStyle:    LineJoinStyleMiter,
		Transparency: 0.1,
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = richShapeLineXML(line)
	}
}

func BenchmarkRichLineBaseline(b *testing.B) {
	line := RichShapeLineSpec{
		Color:        "1F3864",
		Width:        12700,
		DashStyle:    LineDashStyleDash,
		CapStyle:     LineCapStyleRound,
		JoinStyle:    LineJoinStyleMiter,
		Transparency: 0.1,
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = baselineRichLine(line)
	}
}

// BenchmarkRichOuterShadow stresses the per-shape outer shadow emission path.
func BenchmarkRichOuterShadow(b *testing.B) {
	shadow := RichShapeShadowSpec{
		Type:            ShadowTypeOuter,
		Color:           "000000",
		Transparency:    0.4,
		BlurRadius:      40000,
		Distance:        20000,
		Angle:           5400000,
		Alignment:       "ctr",
		RotateWithShape: false,
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = richOuterShadowXML(shadow)
	}
}

func BenchmarkRichOuterShadowBaseline(b *testing.B) {
	shadow := RichShapeShadowSpec{
		Type:            ShadowTypeOuter,
		Color:           "000000",
		Transparency:    0.4,
		BlurRadius:      40000,
		Distance:        20000,
		Angle:           5400000,
		Alignment:       "ctr",
		RotateWithShape: false,
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = baselineOuterShadow(shadow)
	}
}

// BenchmarkRichPerspectiveShadow stresses the more complex perspective shadow path.
func BenchmarkRichPerspectiveShadow(b *testing.B) {
	shadow := RichShapeShadowSpec{
		Type:         ShadowTypePerspective,
		Color:        "000000",
		Transparency: 0.5,
		Distance:     30000,
		Angle:        2700000,
		SkewX:        0.1,
		SkewY:        -0.05,
		ScaleX:       1.2,
		ScaleY:       0.9,
		Alignment:    "br",
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = richPerspectiveShadowXML(shadow)
	}
}

func BenchmarkRichPerspectiveShadowBaseline(b *testing.B) {
	shadow := RichShapeShadowSpec{
		Type:         ShadowTypePerspective,
		Color:        "000000",
		Transparency: 0.5,
		Distance:     30000,
		Angle:        2700000,
		SkewX:        0.1,
		SkewY:        -0.05,
		ScaleX:       1.2,
		ScaleY:       0.9,
		Alignment:    "br",
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = baselinePerspectiveShadow(shadow)
	}
}

// BenchmarkPlaceholderTextStyleFull renders the full placeholder text style block.
func BenchmarkPlaceholderTextStyleFull(b *testing.B) {
	bold := true
	italic := true
	sz := 18
	color := "#FF0000"
	underline := "single"
	font := "Calibri"
	align := "ctr"
	ts := &PlaceholderTextStyleSpec{
		SizePt:    &sz,
		Color:     &color,
		Bold:      &bold,
		Italic:    &italic,
		Underline: &underline,
		Align:     &align,
		Font:      &font,
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = renderPlaceholderTextStyle(ts)
	}
}

// BenchmarkTextLevelStyles renders a typical 9-level body text style block.
func BenchmarkTextLevelStyles(b *testing.B) {
	levels := makeTextLevelFixture()
	b.ReportAllocs()
	for b.Loop() {
		_ = textLevelStylesXML(levels)
	}
}

func BenchmarkTextLevelStylesBaseline(b *testing.B) {
	levels := makeTextLevelFixture()
	b.ReportAllocs()
	for b.Loop() {
		_ = baselineTextLevelStyles(levels)
	}
}

func makeTextLevelFixture() []TextLevelStyle {
	levels := make([]TextLevelStyle, 9)
	for i := range levels {
		levels[i] = TextLevelStyle{
			Level:      i,
			IndentEMU:  int64(228600 * (i + 1)),
			BulletChar: "•",
			SizePt:     18 - i,
			Bold:       i%2 == 0,
			Italic:     i%3 == 0,
			Color:      "#1F3864",
			Font:       "Calibri",
		}
	}
	return levels
}

// BenchmarkEscapePlain measures the fast-path (no replacements).
func BenchmarkEscapePlain(b *testing.B) {
	values := []string{"4472C4", "rId12", "Calibri", "en-US", "ctr", "ppt/media/image1.png"}
	b.ReportAllocs()
	for b.Loop() {
		for _, v := range values {
			_ = Escape(v)
		}
	}
}

func BenchmarkEscapePlainBaseline(b *testing.B) {
	values := []string{"4472C4", "rId12", "Calibri", "en-US", "ctr", "ppt/media/image1.png"}
	b.ReportAllocs()
	for b.Loop() {
		for _, v := range values {
			_ = baselineEscape(v)
		}
	}
}

// BenchmarkEscapeWithSpecials measures the slow-path (replacements needed).
func BenchmarkEscapeWithSpecials(b *testing.B) {
	values := []string{"A & B", "<title>", `"quoted"`, "it's"}
	b.ReportAllocs()
	for b.Loop() {
		for _, v := range values {
			_ = Escape(v)
		}
	}
}

// BenchmarkRenderShapeFullStack composes a realistic per-shape render:
// solid fill + rich line + outer shadow.
func BenchmarkRenderShapeFullStack(b *testing.B) {
	fill, line, shadow := makeFullStackFixture()
	b.ReportAllocs()
	for b.Loop() {
		_ = richSolidFillXML(fill)
		_ = richShapeLineXML(line)
		_ = richOuterShadowXML(shadow)
	}
}

func BenchmarkRenderShapeFullStackBaseline(b *testing.B) {
	fill, line, shadow := makeFullStackFixture()
	b.ReportAllocs()
	for b.Loop() {
		_ = baselineRichSolidFill(fill)
		_ = baselineRichLine(line)
		_ = baselineOuterShadow(shadow)
	}
}

func makeFullStackFixture() (SolidFillSpec, RichShapeLineSpec, RichShapeShadowSpec) {
	fill := SolidFillSpec{Color: "4472C4", Transparency: 0.2}
	line := RichShapeLineSpec{
		Color:     "1F3864",
		Width:     12700,
		DashStyle: LineDashStyleDash,
		CapStyle:  LineCapStyleRound,
		JoinStyle: LineJoinStyleMiter,
	}
	shadow := RichShapeShadowSpec{
		Type:         ShadowTypeOuter,
		Color:        "000000",
		Transparency: 0.4,
		BlurRadius:   40000,
		Distance:     20000,
		Angle:        5400000,
	}
	return fill, line, shadow
}
