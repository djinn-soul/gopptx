package pptxxml

import (
	"fmt"
	"strings"
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
<p:sldLayout xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" preserve="1">
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
	if masterIndex <= 1 {
		return fmt.Sprintf("slideLayout%d.xml", layoutIndex)
	}
	return fmt.Sprintf("slideLayout%d_m%d.xml", layoutIndex, masterIndex)
}

// SlideLayoutRelationships renders ppt/slideLayouts/_rels/slideLayoutN.xml.rels.
func SlideLayoutRelationships(masterIndex int) string {
	if masterIndex < 1 {
		masterIndex = 1
	}
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="../slideMasters/slideMaster` + fmt.Sprintf("%d", masterIndex) + `.xml"/>
</Relationships>`
}

// SlideMaster renders ppt/slideMasters/slideMaster1.xml.
func SlideMaster(spec *SlideMasterSpec) string {
	bgXML := slideDefaultBackground
	if spec != nil && spec.Background != nil {
		bgXML = backgroundXML(spec.Background)
	}
	masterIndex := 1
	if spec != nil && spec.MasterIndex > 0 {
		masterIndex = spec.MasterIndex
	}
	layoutIDBase := int64(2147483649 + (masterIndex-1)*6)

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
<p:sldMaster xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
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
<p:clrMap bg1="` + bg1 + `" tx1="` + tx1 + `" bg2="lt2" tx2="dk2" accent1="accent1" accent2="accent2" accent3="accent3" accent4="accent4" accent5="accent5" accent6="accent6" hlink="hlink" folHlink="folHlink"/>
<p:sldLayoutIdLst>
<p:sldLayoutId id="` + fmt.Sprintf("%d", layoutIDBase) + `" r:id="rId1"/>
<p:sldLayoutId id="` + fmt.Sprintf("%d", layoutIDBase+1) + `" r:id="rId2"/>
<p:sldLayoutId id="` + fmt.Sprintf("%d", layoutIDBase+2) + `" r:id="rId3"/>
<p:sldLayoutId id="` + fmt.Sprintf("%d", layoutIDBase+3) + `" r:id="rId4"/>
<p:sldLayoutId id="` + fmt.Sprintf("%d", layoutIDBase+4) + `" r:id="rId5"/>
<p:sldLayoutId id="` + fmt.Sprintf("%d", layoutIDBase+5) + `" r:id="rId6"/>
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
func textLevelStylesXML(levels []TextLevelStyle) string {
	var b strings.Builder
	for _, lvl := range levels {
		lvlNum := lvl.Level + 1 // 0-based → 1-based
		if lvlNum < 1 {
			lvlNum = 1
		}
		if lvlNum > 9 {
			lvlNum = 9
		}

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

// SlideMasterRelationships renders ppt/slideMasters/_rels/slideMasterN.xml.rels.
// imageTargets are optional media paths for master images (e.g. "../media/image1.png").
func SlideMasterRelationships(imageTargets []string, masterIndex int) string {
	if masterIndex < 1 {
		masterIndex = 1
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/` + slideLayoutPartName(1, masterIndex) + `"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/` + slideLayoutPartName(2, masterIndex) + `"/>
<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/` + slideLayoutPartName(3, masterIndex) + `"/>
<Relationship Id="rId4" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/` + slideLayoutPartName(4, masterIndex) + `"/>
<Relationship Id="rId5" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/` + slideLayoutPartName(5, masterIndex) + `"/>
<Relationship Id="rId6" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/` + slideLayoutPartName(6, masterIndex) + `"/>
<Relationship Id="rId7" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme1.xml"/>`)
	for i, target := range imageTargets {
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="%s"/>`, 8+i, Escape(target)))
	}
	b.WriteString(`
</Relationships>`)
	return b.String()
}

// Theme renders ppt/theme/theme1.xml.
func Theme(spec *ThemeSpec) string {
	name := "Office Theme"
	clrName := "Office"
	fontName := "Office"

	dk1, lt1, dk2, lt2 := "windowText", "window", "1F497D", "EEECE1"
	dk1Last, lt1Last := "000000", "FFFFFF"
	accent1, accent2, accent3 := "4F81BD", "C0504D", "9BBB59"
	accent4, accent5, accent6 := "8064A2", "4BACC6", "F79646"
	hlink, folHlink := "0000FF", "800080"

	majorFont, minorFont := "Calibri", "Calibri"

	if spec != nil {
		if spec.Name != "" {
			name = spec.Name
			clrName = spec.Name
			fontName = spec.Name
		}
		c := spec.Colors
		if c.Dk1 != "" {
			dk1 = "windowText"
			dk1Last = strings.TrimPrefix(c.Dk1, "#")
		}
		if c.Lt1 != "" {
			lt1 = "window"
			lt1Last = strings.TrimPrefix(c.Lt1, "#")
		}
		if c.Dk2 != "" {
			dk2 = strings.TrimPrefix(c.Dk2, "#")
		}
		if c.Lt2 != "" {
			lt2 = strings.TrimPrefix(c.Lt2, "#")
		}
		if c.Accent1 != "" {
			accent1 = strings.TrimPrefix(c.Accent1, "#")
		}
		if c.Accent2 != "" {
			accent2 = strings.TrimPrefix(c.Accent2, "#")
		}
		if c.Accent3 != "" {
			accent3 = strings.TrimPrefix(c.Accent3, "#")
		}
		if c.Accent4 != "" {
			accent4 = strings.TrimPrefix(c.Accent4, "#")
		}
		if c.Accent5 != "" {
			accent5 = strings.TrimPrefix(c.Accent5, "#")
		}
		if c.Accent6 != "" {
			accent6 = strings.TrimPrefix(c.Accent6, "#")
		}
		if c.Hlink != "" {
			hlink = strings.TrimPrefix(c.Hlink, "#")
		}
		if c.FolHlink != "" {
			folHlink = strings.TrimPrefix(c.FolHlink, "#")
		}

		if spec.Fonts.MajorFont != "" {
			majorFont = spec.Fonts.MajorFont
		}
		if spec.Fonts.MinorFont != "" {
			minorFont = spec.Fonts.MinorFont
		}
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="`+Escape(name)+` marketing">
<a:themeElements>
<a:clrScheme name="`+Escape(clrName)+` colors">
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
</a:clrScheme>
<a:fontScheme name="`+Escape(fontName)+` fonts">
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
</a:fontScheme>
<a:fmtScheme name="Office">
<a:fillStyleLst>
<a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
<a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:tint val="50000"/><a:satMod val="300000"/></a:schemeClr></a:gs><a:gs pos="35000"><a:schemeClr val="phClr"><a:tint val="37000"/><a:satMod val="300000"/></a:schemeClr></a:gs><a:gs pos="100000"><a:schemeClr val="phClr"><a:tint val="15000"/><a:satMod val="350000"/></a:schemeClr></a:gs></a:gsLst><a:lin ang="16200000" scaled="1"/></a:gradFill>
<a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:shade val="51000"/><a:satMod val="130000"/></a:schemeClr></a:gs><a:gs pos="80000"><a:schemeClr val="phClr"><a:shade val="93000"/><a:satMod val="130000"/></a:schemeClr></a:gs><a:gs pos="100000"><a:schemeClr val="phClr"><a:shade val="94000"/><a:satMod val="135000"/></a:schemeClr></a:gs></a:gsLst><a:lin ang="16200000" scaled="0"/></a:gradFill>
</a:fillStyleLst>
<a:lnStyleLst>
<a:ln w="9525" cap="flat" cmpd="sng" algn="ctr"><a:solidFill><a:schemeClr val="phClr"><a:shade val="95000"/><a:satMod val="105000"/></a:schemeClr></a:solidFill><a:prstDash val="solid"/></a:ln>
<a:ln w="25400" cap="flat" cmpd="sng" algn="ctr"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:prstDash val="solid"/></a:ln>
<a:ln w="38100" cap="flat" cmpd="sng" algn="ctr"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:prstDash val="solid"/></a:ln>
</a:lnStyleLst>
<a:effectStyleLst>
<a:effectStyle><a:effectLst/></a:effectStyle>
<a:effectStyle><a:effectLst/></a:effectStyle>
<a:effectStyle><a:effectLst/></a:effectStyle>
</a:effectStyleLst>
<a:bgFillStyleLst>
<a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
<a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:tint val="40000"/><a:satMod val="350000"/></a:schemeClr></a:gs><a:gs pos="40000"><a:schemeClr val="phClr"><a:tint val="45000"/><a:shade val="99000"/><a:satMod val="350000"/></a:schemeClr></a:gs><a:gs pos="100000"><a:schemeClr val="phClr"><a:shade val="20000"/><a:satMod val="255000"/></a:schemeClr></a:gs></a:gsLst><a:path path="circle"><a:fillToRect l="50000" t="-80000" r="50000" b="180000"/></a:path></a:gradFill>
<a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:tint val="80000"/><a:satMod val="300000"/></a:schemeClr></a:gs><a:gs pos="100000"><a:schemeClr val="phClr"><a:shade val="30000"/><a:satMod val="200000"/></a:schemeClr></a:gs></a:gsLst><a:path path="circle"><a:fillToRect l="50000" t="50000" r="50000" b="50000"/></a:path></a:gradFill>
</a:bgFillStyleLst>
</a:fmtScheme>
</a:themeElements>
<a:objectDefaults/>
<a:extraClrSchemeLst/>
</a:theme>`,
		dk1, dk1Last, lt1, lt1Last, dk2, lt2,
		accent1, accent2, accent3, accent4, accent5, accent6,
		hlink, folHlink,
		majorFont, minorFont)
}
