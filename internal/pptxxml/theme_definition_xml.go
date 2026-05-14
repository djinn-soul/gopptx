package pptxxml

import (
	"fmt"
	"strings"
)

const defaultThemeFont = "Calibri"

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
	majorFont, minorFont := defaultThemeFont, defaultThemeFont
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
