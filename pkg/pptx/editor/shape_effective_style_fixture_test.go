package editor

import (
	"path/filepath"
	"strings"
	"testing"
)

// writeInheritanceFixture builds a minimal deck whose title placeholder states
// nothing but its text, so every visible property has to come from the layout,
// the master or the theme.
func writeInheritanceFixture(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	files := map[string]string{
		"[Content_Types].xml": xmlHeader +
			`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
			`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
			`<Default Extension="xml" ContentType="application/xml"/>` +
			`<Override PartName="/ppt/presentation.xml" ` +
			`ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>` +
			`<Override PartName="/ppt/slides/slide1.xml" ` +
			`ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
			`<Override PartName="/ppt/slideLayouts/slideLayout1.xml" ` +
			`ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>` +
			`<Override PartName="/ppt/slideMasters/slideMaster1.xml" ` +
			`ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"/>` +
			`<Override PartName="/ppt/theme/theme1.xml" ` +
			`ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>` +
			`</Types>`,
		"_rels/.rels": relsDoc(
			`<Relationship Id="rId1" Type="` + relTypeOfficeDocument + `" Target="ppt/presentation.xml"/>`,
		),
		"ppt/presentation.xml": xmlHeader + presentationNS +
			`<p:sldMasterIdLst><p:sldMasterId id="2147483648" r:id="rId2"/></p:sldMasterIdLst>` +
			`<p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst></p:presentation>`,
		"ppt/_rels/presentation.xml.rels": relsDoc(
			`<Relationship Id="rId1" Type="` + relTypeSlide + `" Target="slides/slide1.xml"/>` +
				`<Relationship Id="rId2" Type="` + relTypeSlideMaster + `" Target="slideMasters/slideMaster1.xml"/>`,
		),
		"ppt/slides/slide1.xml": inheritanceSlideXML,
		"ppt/slides/_rels/slide1.xml.rels": relsDoc(
			`<Relationship Id="rId1" Type="` + relTypeSlideLayout + `" Target="../slideLayouts/slideLayout1.xml"/>`,
		),
		"ppt/slideLayouts/slideLayout1.xml": inheritanceLayoutXML,
		"ppt/slideLayouts/_rels/slideLayout1.xml.rels": relsDoc(
			`<Relationship Id="rId1" Type="` + relTypeSlideMaster + `" Target="../slideMasters/slideMaster1.xml"/>`,
		),
		"ppt/slideMasters/slideMaster1.xml": inheritanceMasterXML,
		"ppt/slideMasters/_rels/slideMaster1.xml.rels": relsDoc(
			`<Relationship Id="rId1" Type="` + relTypeSlideLayout + `" Target="../slideLayouts/slideLayout1.xml"/>` +
				`<Relationship Id="rId2" Type="` + relTypeTheme + `" Target="../theme/theme1.xml"/>`,
		),
		"ppt/theme/theme1.xml": inheritanceThemeXML,
	}
	if err := writeZipFixture(path, files); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

// replaceOnce fails the test when the text it is asked to rewrite is not there,
// so a fixture change cannot silently turn a test into a no-op.
func replaceOnce(t *testing.T, source, old, replacement string) string {
	t.Helper()
	if !strings.Contains(source, old) {
		t.Fatalf("fixture does not contain %q", old)
	}
	return strings.Replace(source, old, replacement, 1)
}

func relsDoc(body string) string {
	return xmlHeader +
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
		body + `</Relationships>`
}

const (
	xmlHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`

	drawingNS      = `xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"`
	relationshipNS = `xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`
	presentationML = `xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"`

	presentationNS = `<p:presentation ` + drawingNS + ` ` + relationshipNS + ` ` + presentationML + `>`

	relTypeOfficeDocument = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"
	relTypeSlide          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide"
	relTypeSlideLayout    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout"
	relTypeSlideMaster    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster"
	relTypeTheme          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme"

	// The slide's title placeholder states only its text.
	inheritanceSlideXML = xmlHeader +
		`<p:sld ` + drawingNS + ` ` + relationshipNS + ` ` + presentationML + `>` +
		`<p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Title 1"/><p:cNvSpPr/>` +
		`<p:nvPr><p:ph type="title" idx="0"/></p:nvPr></p:nvSpPr>` +
		`<p:spPr/>` +
		`<p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr lang="en-US"/><a:t>Title</a:t></a:r></a:p>` +
		`</p:txBody></p:sp>` +
		`</p:spTree></p:cSld></p:sld>`

	// The layout's title placeholder carries the position, the scheme colour,
	// the size and a theme font reference.
	inheritanceLayoutXML = xmlHeader +
		`<p:sldLayout ` + drawingNS + ` ` + relationshipNS + ` ` + presentationML + `>` +
		`<p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Title Placeholder 1"/><p:cNvSpPr/>` +
		`<p:nvPr><p:ph type="title" idx="0"/></p:nvPr></p:nvSpPr>` +
		`<p:spPr><a:xfrm><a:off x="838200" y="365125"/><a:ext cx="8000000" cy="1325563"/></a:xfrm></p:spPr>` +
		`<p:txBody><a:bodyPr/><a:lstStyle><a:lvl1pPr><a:defRPr sz="4000">` +
		`<a:solidFill><a:schemeClr val="accent1"/></a:solidFill>` +
		`<a:latin typeface="+mj-lt"/></a:defRPr></a:lvl1pPr></a:lstStyle>` +
		`<a:p><a:endParaRPr lang="en-US"/></a:p></p:txBody></p:sp>` +
		`</p:spTree></p:cSld></p:sldLayout>`

	// The master states only bold, through its title text style.
	inheritanceMasterXML = xmlHeader +
		`<p:sldMaster ` + drawingNS + ` ` + relationshipNS + ` ` + presentationML + `>` +
		`<p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`</p:spTree></p:cSld>` +
		`<p:sldLayoutIdLst><p:sldLayoutId id="2147483649" r:id="rId1"/></p:sldLayoutIdLst>` +
		`<p:txStyles><p:titleStyle><a:lvl1pPr><a:defRPr b="1"/></a:lvl1pPr></p:titleStyle>` +
		`<p:bodyStyle/><p:otherStyle/></p:txStyles>` +
		`</p:sldMaster>`

	inheritanceThemeXML = xmlHeader +
		`<a:theme ` + drawingNS + ` name="Test">` +
		`<a:themeElements><a:clrScheme name="Test">` +
		`<a:dk1><a:srgbClr val="000000"/></a:dk1><a:lt1><a:srgbClr val="FFFFFF"/></a:lt1>` +
		`<a:dk2><a:srgbClr val="44546A"/></a:dk2><a:lt2><a:srgbClr val="E7E6E6"/></a:lt2>` +
		`<a:accent1><a:srgbClr val="4472C4"/></a:accent1><a:accent2><a:srgbClr val="ED7D31"/></a:accent2>` +
		`<a:accent3><a:srgbClr val="A5A5A5"/></a:accent3><a:accent4><a:srgbClr val="FFC000"/></a:accent4>` +
		`<a:accent5><a:srgbClr val="5B9BD5"/></a:accent5><a:accent6><a:srgbClr val="70AD47"/></a:accent6>` +
		`<a:hlink><a:srgbClr val="0563C1"/></a:hlink><a:folHlink><a:srgbClr val="954F72"/></a:folHlink>` +
		`</a:clrScheme>` +
		`<a:fontScheme name="Test"><a:majorFont><a:latin typeface="Georgia"/></a:majorFont>` +
		`<a:minorFont><a:latin typeface="Verdana"/></a:minorFont></a:fontScheme>` +
		`<a:fmtScheme name="Test"/></a:themeElements></a:theme>`
)
