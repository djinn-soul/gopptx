package pptxxml

import "strings"

const slideHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:cSld>
<p:bg>
<p:bgRef idx="1001">
<a:schemeClr val="bg1"/>
</p:bgRef>
</p:bg>
<p:spTree>
<p:nvGrpSpPr>
<p:cNvPr id="1" name=""/>
<p:cNvGrpSpPr/>
<p:nvPr/>
</p:nvGrpSpPr>
<p:grpSpPr>
<a:xfrm>
<a:off x="0" y="0"/>
<a:ext cx="9144000" cy="6858000"/>
<a:chOff x="0" y="0"/>
<a:chExt cx="9144000" cy="6858000"/>
</a:xfrm>
</p:grpSpPr>`

const slideFooter = `
</p:spTree>
</p:cSld>
<p:clrMapOvr>
<a:masterClrMapping/>
</p:clrMapOvr>
</p:sld>`

// SlideWithContent renders a title+bullets slide.
func SlideWithContent(title string, bullets []string) string {
	var b strings.Builder
	b.WriteString(slideHeader)
	b.WriteString(titleShape(title))
	if len(bullets) > 0 {
		b.WriteString(contentShape(bullets))
	}
	b.WriteString(slideFooter)
	return b.String()
}

// SlideRelationships renders ppt/slides/_rels/slideN.xml.rels.
func SlideRelationships() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
</Relationships>`
}

func titleShape(title string) string {
	escaped := Escape(title)
	return `
<p:sp>
<p:nvSpPr>
<p:cNvPr id="2" name="Title"/>
<p:cNvSpPr txBox="1"/>
<p:nvPr/>
</p:nvSpPr>
<p:spPr>
<a:xfrm>
<a:off x="457200" y="274638"/>
<a:ext cx="8230200" cy="1143000"/>
</a:xfrm>
<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
<a:noFill/>
</p:spPr>
<p:txBody>
<a:bodyPr wrap="square" rtlCol="0" anchor="ctr"/>
<a:lstStyle/>
<a:p>
<a:pPr algn="l"/>
<a:r>
<a:rPr lang="en-US" sz="4400" b="1" i="0" dirty="0"/>
<a:t>` + escaped + `</a:t>
</a:r>
</a:p>
</p:txBody>
</p:sp>`
}

func contentShape(bullets []string) string {
	var b strings.Builder
	b.WriteString(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="3" name="Content"/>
<p:cNvSpPr txBox="1"/>
<p:nvPr/>
</p:nvSpPr>
<p:spPr>
<a:xfrm>
<a:off x="457200" y="1600200"/>
<a:ext cx="8230200" cy="4572000"/>
</a:xfrm>
<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
<a:noFill/>
</p:spPr>
<p:txBody>
<a:bodyPr wrap="square" rtlCol="0"/>
<a:lstStyle/>`)

	for _, bullet := range bullets {
		b.WriteString(bulletParagraph(bullet))
	}

	b.WriteString(`
</p:txBody>
</p:sp>`)
	return b.String()
}

func bulletParagraph(text string) string {
	escaped := Escape(text)
	return `
<a:p>
<a:pPr lvl="0" marL="457200" indent="-457200"><a:buChar char="•"/></a:pPr>
<a:r>
<a:rPr lang="en-US" sz="2800" b="0" i="0" dirty="0"/>
<a:t>` + escaped + `</a:t>
</a:r>
</a:p>`
}
