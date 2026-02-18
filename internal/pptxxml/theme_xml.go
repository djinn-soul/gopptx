package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	layoutsPerMaster = 6
)

// ColorSchemeSpec defines the 12 colors in a theme.
type ColorSchemeSpec struct {
	Name     string
	Dk1      string
	Lt1      string
	Dk2      string
	Lt2      string
	Accent1  string
	Accent2  string
	Accent3  string
	Accent4  string
	Accent5  string
	Accent6  string
	Hlink    string
	FolHlink string
}

// FontSchemeSpec defines heading and body fonts.
type FontSchemeSpec struct {
	Name      string
	MajorFont string
	MinorFont string
}

// ThemeSpec defines a theme's color and font elements.
type ThemeSpec struct {
	Name   string
	Colors ColorSchemeSpec
	Fonts  FontSchemeSpec
}

// SlideMasterSpec defines the appearance of the slide master.
type SlideMasterSpec struct {
	MasterIndex  int
	Background   *SlideBackgroundSpec
	FooterText   string
	Shapes       []ShapeSpec
	Images       []ImageRef
	ColorMapping *ColorMappingSpec
	TxStyles     *TxStylesSpec
}

// TxStylesSpec defines the default text styles for a slide master.
// Each field holds up to 9 levels of text styling (Lvl1–Lvl9).
type TxStylesSpec struct {
	TitleStyle []TextLevelStyle
	BodyStyle  []TextLevelStyle
	OtherStyle []TextLevelStyle
}

// TextLevelStyle defines default text properties for one indent level.
type TextLevelStyle struct {
	Level      int    // 0-based (0=Lvl1, 8=Lvl9)
	Font       string // Typeface override
	SizePt     int    // Size in points
	Bold       bool
	Italic     bool
	Color      string // 6-digit hex RGB
	BulletChar string // Bullet character override
	IndentEMU  int64  // Left indent in EMU
}

// ColorMappingSpec describes how theme colors map to functional roles on slides.
type ColorMappingSpec struct {
	BG1 string
	TX1 string
}

// SlideLayout renders ppt/slideLayouts/slideLayout1.xml.
func SlideLayout() string {
	return SlideLayoutTitleAndContent()
}

// SlideLayoutTitleAndContent renders a title-and-content layout.
func SlideLayoutTitleAndContent() string {
	return slideLayout("Title and Content")
}

// SlideLayoutTitleOnly renders a title-only layout.
func SlideLayoutTitleOnly() string {
	return slideLayout("Title Only")
}

// SlideLayoutBlank renders a blank layout.
func SlideLayoutBlank() string {
	return slideLayout("Blank")
}

// SlideLayoutCenteredTitle renders a centered-title layout.
func SlideLayoutCenteredTitle() string {
	return slideLayout("Centered Title")
}

// SlideLayoutTitleAndBigContent renders a title-and-big-content layout.
func SlideLayoutTitleAndBigContent() string {
	return slideLayout("Title and Big Content")
}

// SlideLayoutTwoColumn renders a two-column layout.
func SlideLayoutTwoColumn() string {
	return slideLayout("Two Column")
}

func slideLayout(name string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldLayout xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" preserve="1">
<p:cSld name="` + Escape(name) + `">
<p:spTree>
<p:nvGrpSpPr>
<p:cNvPr id="1" name=""/>
<p:cNvGrpSpPr/>
<p:nvPr/>
</p:nvGrpSpPr>
<p:grpSpPr>
<a:xfrm>
<a:off x="0" y="0"/>
<a:ext cx="0" cy="0"/>
<a:chOff x="0" y="0"/>
<a:chExt cx="0" cy="0"/>
</a:xfrm>
</p:grpSpPr>
</p:spTree>
</p:cSld>
<p:clrMapOvr>
<a:masterClrMapping/>
</p:clrMapOvr>
</p:sldLayout>`
}

func slideLayoutPartName(layoutIndex, masterIndex int) string {
	if masterIndex < 1 {
		masterIndex = 1
	}
	globalLayoutIndex := (masterIndex-1)*layoutsPerMaster + layoutIndex
	return fmt.Sprintf("slideLayout%d.xml", globalLayoutIndex)
}

// SlideLayoutRelationships renders ppt/slideLayouts/_rels/slideLayoutN.xml.rels.
func SlideLayoutRelationships(masterIndex int) string {
	if masterIndex < 1 {
		masterIndex = 1
	}
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" ` +
		`Target="../slideMasters/slideMaster` + strconv.Itoa(masterIndex) + `.xml"/>
</Relationships>`
}

// SlideMaster renders ppt/slideMasters/slideMaster1.xml.
//
//nolint:mnd // Contains specific ID base calculations required by the spec.
func SlideMaster(spec *SlideMasterSpec) string {
	bgXML := slideDefaultBackground
	if spec != nil && spec.Background != nil {
		bgXML = backgroundXML(spec.Background)
	}
	masterIndex := 1
	if spec != nil && spec.MasterIndex > 0 {
		masterIndex = spec.MasterIndex
	}
	masterAttrs := ``
	if masterIndex > 1 {
		// PowerPoint-authored packages mark additional master families as preserved.
		masterAttrs = ` preserve="1"`
	}
	// Allocate in blocks to avoid overlap with master IDs across multiple masters.
	layoutIDBase := int64(2147483649 + (masterIndex-1)*7)

	footerXML := ""
	if spec != nil && spec.FooterText != "" {
		footerXML = fmt.Sprintf(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="10" name="Footer Placeholder"/>
<p:cNvSpPr>
<a:spLocks noGrp="1"/>
</p:cNvSpPr>
<p:nvPr>
<p:ph type="ftr" sz="quarter" idx="10"/>
</p:nvPr>
</p:nvSpPr>
<p:spPr>
<a:xfrm>
<a:off x="0" y="6350000"/>
<a:ext cx="9144000" cy="508000"/>
</a:xfrm>
</p:spPr>
<p:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p>
<a:r>
<a:rPr lang="en-US" smtClean="0"/>
<a:t>%s</a:t>
</a:r>
</a:p>
</p:txBody>
</p:sp>`, Escape(spec.FooterText))
	}

	masterElementsXML := masterShapesAndImagesXML(spec)

	bg1, tx1 := "lt1", "dk1"
	if spec != nil && spec.ColorMapping != nil {
		if spec.ColorMapping.BG1 != "" {
			bg1 = spec.ColorMapping.BG1
		}
		if spec.ColorMapping.TX1 != "" {
			tx1 = spec.ColorMapping.TX1
		}
	}

	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldMaster xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"` + masterAttrs + `>
<p:cSld>
` + bgXML + `
<p:spTree>
<p:nvGrpSpPr>
<p:cNvPr id="1" name=""/>
<p:cNvGrpSpPr/>
<p:nvPr/>
</p:nvGrpSpPr>
<p:grpSpPr>
<a:xfrm>
<a:off x="0" y="0"/>
<a:ext cx="0" cy="0"/>
<a:chOff x="0" y="0"/>
<a:chExt cx="0" cy="0"/>
</a:xfrm>
</p:grpSpPr>
` + footerXML + masterElementsXML + `
</p:spTree>
</p:cSld>
<p:clrMap bg1="` + bg1 + `" tx1="` + tx1 + `" bg2="lt2" tx2="dk2" accent1="accent1" accent2="accent2" ` +
		`accent3="accent3" accent4="accent4" accent5="accent5" accent6="accent6" hlink="hlink" folHlink="folHlink"/>
<p:sldLayoutIdLst>
<p:sldLayoutId id="` + strconv.FormatInt(layoutIDBase, 10) + `" r:id="rId1"/>
<p:sldLayoutId id="` + strconv.FormatInt(layoutIDBase+1, 10) + `" r:id="rId2"/>
<p:sldLayoutId id="` + strconv.FormatInt(layoutIDBase+2, 10) + `" r:id="rId3"/>
<p:sldLayoutId id="` + strconv.FormatInt(layoutIDBase+3, 10) + `" r:id="rId4"/>
<p:sldLayoutId id="` + strconv.FormatInt(layoutIDBase+4, 10) + `" r:id="rId5"/>
<p:sldLayoutId id="` + strconv.FormatInt(layoutIDBase+5, 10) + `" r:id="rId6"/>
</p:sldLayoutIdLst>` + txStylesXML(spec) + `
</p:sldMaster>`
}

// masterShapesAndImagesXML renders shapes and images within the master spTree.
func masterShapesAndImagesXML(spec *SlideMasterSpec) string {
	if spec == nil {
		return ""
	}
	var b strings.Builder
	// Shape IDs start at 20 to avoid footer placeholder (id=10)
	nextID := 20
	for _, shape := range spec.Shapes {
		b.WriteString(customShapeXML(shape, nextID))
		nextID++
	}
	for _, img := range spec.Images {
		b.WriteString(imageShape(img, nextID))
		nextID++
	}
	return b.String()
}

// txStylesXML renders the <p:txStyles> block for a slide master.
func txStylesXML(spec *SlideMasterSpec) string {
	if spec == nil || spec.TxStyles == nil {
		return ""
	}
	tx := spec.TxStyles
	var b strings.Builder
	b.WriteString(`
<p:txStyles>`)
	if len(tx.TitleStyle) > 0 {
		b.WriteString(`
<p:titleStyle>`)
		b.WriteString(textLevelStylesXML(tx.TitleStyle))
		b.WriteString(`
</p:titleStyle>`)
	}
	if len(tx.BodyStyle) > 0 {
		b.WriteString(`
<p:bodyStyle>`)
		b.WriteString(textLevelStylesXML(tx.BodyStyle))
		b.WriteString(`
</p:bodyStyle>`)
	}
	if len(tx.OtherStyle) > 0 {
		b.WriteString(`
<p:otherStyle>`)
		b.WriteString(textLevelStylesXML(tx.OtherStyle))
		b.WriteString(`
</p:otherStyle>`)
	}
	b.WriteString(`
</p:txStyles>`)
	return b.String()
}

// textLevelStylesXML renders <a:lvlNpPr> elements for each text level.
//
//nolint:mnd // Contains specific level limits and point conversion factors.
func textLevelStylesXML(levels []TextLevelStyle) string {
	var b strings.Builder
	for _, lvl := range levels {
		lvlNum := max(
			// 0-based → 1-based
			lvl.Level+1, 1)
		lvlNum = min(lvlNum, 9)

		attrs := ""
		if lvl.IndentEMU > 0 {
			attrs += fmt.Sprintf(` indent="%d"`, lvl.IndentEMU)
		}

		b.WriteString(fmt.Sprintf(`
<a:lvl%dpPr%s>`, lvlNum, attrs))

		if lvl.BulletChar != "" {
			b.WriteString(fmt.Sprintf(`<a:buChar char="%s"/>`, Escape(lvl.BulletChar)))
		}

		// Default text run properties
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
			b.WriteString(fmt.Sprintf(`<a:latin typeface="%s"/>`, Escape(lvl.Font)))
		}
		b.WriteString(`</a:defRPr>`)
		b.WriteString(fmt.Sprintf(`
</a:lvl%dpPr>`, lvlNum))
	}
	return b.String()
}

const imageRidStart = 8

// SlideMasterRelationships renders ppt/slideMasters/_rels/slideMasterN.xml.rels.
// imageTargets are optional media paths for master images (e.g. "../media/image1.png").
func SlideMasterRelationships(imageTargets []string, masterIndex int, themeIndex int) string {
	if masterIndex < 1 {
		masterIndex = 1
	}
	if themeIndex < 1 {
		themeIndex = 1
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	// Layouts 1-6
	for i := 1; i <= 6; i++ {
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/%s"/>`,
			i, slideLayoutPartName(i, masterIndex)))
	}
	// Theme is rId7
	b.WriteString(fmt.Sprintf(`
<Relationship Id="rId7" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme%d.xml"/>`,
		themeIndex))
	// Images start at rId8
	for i, target := range imageTargets {
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="%s"/>`,
			imageRidStart+i, Escape(target)))
	}
	b.WriteString(`
</Relationships>`)
	return b.String()
}

// Theme renders ppt/theme/theme1.xml.
func Theme(spec *ThemeSpec) string {
	name := "Office Theme"
	if spec != nil && spec.Name != "" {
		name = spec.Name
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="`+
		Escape(name)+` marketing">
<a:themeElements>
%s
%s
%s
</a:themeElements>
<a:objectDefaults/>
<a:extraClrSchemeLst/>
</a:theme>`,
		themeColorsXML(spec),
		themeFontsXML(spec),
		themeFmtSchemeXML())
}

func themeColorsXML(spec *ThemeSpec) string {
	c := resolveThemeColors(spec)

	return fmt.Sprintf(`<a:clrScheme name="`+Escape(c.clrName)+` colors">
<a:dk1><a:sysClr val="%s" lastClr="%s"/></a:dk1>
<a:lt1><a:sysClr val="%s" lastClr="%s"/></a:lt1>
<a:dk2><a:srgbClr val="%s"/></a:dk2>
<a:lt2><a:srgbClr val="%s"/></a:lt2>
<a:accent1><a:srgbClr val="%s"/></a:accent1>
<a:accent2><a:srgbClr val="%s"/></a:accent2>
<a:accent3><a:srgbClr val="%s"/></a:accent3>
<a:accent4><a:srgbClr val="%s"/></a:accent4>
<a:accent5><a:srgbClr val="%s"/></a:accent5>
<a:accent6><a:srgbClr val="%s"/></a:accent6>
<a:hlink><a:srgbClr val="%s"/></a:hlink>
<a:folHlink><a:srgbClr val="%s"/></a:folHlink>
</a:clrScheme>`,
		c.dk1, c.dk1Last, c.lt1, c.lt1Last, c.dk2, c.lt2,
		c.accent1, c.accent2, c.accent3, c.accent4, c.accent5, c.accent6,
		c.hlink, c.folHlink)
}

type resolvedThemeColors struct {
	clrName                                              string
	dk1, dk1Last, lt1, lt1Last                           string
	dk2, lt2                                             string
	accent1, accent2, accent3, accent4, accent5, accent6 string
	hlink, folHlink                                      string
}

func resolveThemeColors(spec *ThemeSpec) resolvedThemeColors {
	res := resolvedThemeColors{
		clrName: "Office",
		dk1:     "windowText", lt1: "window", dk2: "1F497D", lt2: "EEECE1",
		dk1Last: "000000", lt1Last: "FFFFFF",
		accent1: "4F81BD", accent2: "C0504D", accent3: "9BBB59",
		accent4: "8064A2", accent5: "4BACC6", accent6: "F79646",
		hlink: "0000FF", folHlink: "800080",
	}

	if spec == nil {
		return res
	}

	if spec.Name != "" {
		res.clrName = spec.Name
	}

	c := spec.Colors
	if c.Dk1 != "" {
		res.dk1Last = strings.TrimPrefix(c.Dk1, "#")
	}
	if c.Lt1 != "" {
		res.lt1Last = strings.TrimPrefix(c.Lt1, "#")
	}
	if c.Dk2 != "" {
		res.dk2 = strings.TrimPrefix(c.Dk2, "#")
	}
	if c.Lt2 != "" {
		res.lt2 = strings.TrimPrefix(c.Lt2, "#")
	}
	res.accent1 = fallbackColor(c.Accent1, res.accent1)
	res.accent2 = fallbackColor(c.Accent2, res.accent2)
	res.accent3 = fallbackColor(c.Accent3, res.accent3)
	res.accent4 = fallbackColor(c.Accent4, res.accent4)
	res.accent5 = fallbackColor(c.Accent5, res.accent5)
	res.accent6 = fallbackColor(c.Accent6, res.accent6)
	res.hlink = fallbackColor(c.Hlink, res.hlink)
	res.folHlink = fallbackColor(c.FolHlink, res.folHlink)

	return res
}

func fallbackColor(val, def string) string {
	if val == "" {
		return def
	}
	return strings.TrimPrefix(val, "#")
}

func themeFontsXML(spec *ThemeSpec) string {
	fontName := "Office"
	majorFont, minorFont := "Calibri", "Calibri"

	if spec != nil {
		if spec.Name != "" {
			fontName = spec.Name
		}
		if spec.Fonts.MajorFont != "" {
			majorFont = spec.Fonts.MajorFont
		}
		if spec.Fonts.MinorFont != "" {
			minorFont = spec.Fonts.MinorFont
		}
	}

	return fmt.Sprintf(`<a:fontScheme name="`+Escape(fontName)+` fonts">
<a:majorFont>
<a:latin typeface="%s"/>
<a:ea typeface=""/>
<a:cs typeface=""/>
</a:majorFont>
<a:minorFont>
<a:latin typeface="%s"/>
<a:ea typeface=""/>
<a:cs typeface=""/>
</a:minorFont>
</a:fontScheme>`, majorFont, minorFont)
}

func themeFmtSchemeXML() string {
	return `<a:fmtScheme name="Office">
<a:fillStyleLst>
<a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
<a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:tint val="50000"/><a:satMod val="300000"/></a:schemeClr></a:gs>` +
		`<a:gs pos="35000"><a:schemeClr val="phClr"><a:tint val="37000"/><a:satMod val="300000"/></a:schemeClr></a:gs>` +
		`<a:gs pos="100000"><a:schemeClr val="phClr"><a:tint val="15000"/><a:satMod val="350000"/></a:schemeClr></a:gs></a:gsLst>` +
		`<a:lin ang="16200000" scaled="1"/></a:gradFill>
<a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:shade val="51000"/><a:satMod val="130000"/></a:schemeClr></a:gs>` +
		`<a:gs pos="80000"><a:schemeClr val="phClr"><a:shade val="93000"/><a:satMod val="130000"/></a:schemeClr></a:gs>` +
		`<a:gs pos="100000"><a:schemeClr val="phClr"><a:shade val="94000"/><a:satMod val="135000"/></a:schemeClr></a:gs></a:gsLst>` +
		`<a:lin ang="16200000" scaled="0"/></a:gradFill>
</a:fillStyleLst>
<a:lnStyleLst>
<a:ln w="9525" cap="flat" cmpd="sng" algn="ctr">` +
		`<a:solidFill><a:schemeClr val="phClr"><a:shade val="95000"/><a:satMod val="105000"/></a:schemeClr></a:solidFill>` +
		`<a:prstDash val="solid"/></a:ln>
<a:ln w="25400" cap="flat" cmpd="sng" algn="ctr">` +
		`<a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:prstDash val="solid"/></a:ln>
<a:ln w="38100" cap="flat" cmpd="sng" algn="ctr">` +
		`<a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:prstDash val="solid"/></a:ln>
</a:lnStyleLst>
<a:effectStyleLst>
<a:effectStyle><a:effectLst/></a:effectStyle>
<a:effectStyle><a:effectLst/></a:effectStyle>
<a:effectStyle><a:effectLst/></a:effectStyle>
</a:effectStyleLst>
<a:bgFillStyleLst>
<a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
<a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:tint val="40000"/><a:satMod val="350000"/></a:schemeClr></a:gs>` +
		`<a:gs pos="40000"><a:schemeClr val="phClr"><a:tint val="45000"/><a:shade val="99000"/><a:satMod val="350000"/></a:schemeClr></a:gs>` +
		`<a:gs pos="100000"><a:schemeClr val="phClr"><a:shade val="20000"/><a:satMod val="255000"/></a:schemeClr></a:gs></a:gsLst>` +
		`<a:path path="circle"><a:fillToRect l="50000" t="-80000" r="50000" b="180000"/></a:path></a:gradFill>
<a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:tint val="80000"/><a:satMod val="300000"/></a:schemeClr></a:gs>` +
		`<a:gs pos="100000"><a:schemeClr val="phClr"><a:shade val="30000"/><a:satMod val="200000"/></a:schemeClr></a:gs></a:gsLst>` +
		`<a:path path="circle"><a:fillToRect l="50000" t="50000" r="50000" b="50000"/></a:path></a:gradFill>
</a:bgFillStyleLst>
</a:fmtScheme>`
}
