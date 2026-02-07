package pptxxml

import (
	"fmt"
	"strings"
)

const (
	slideLayoutTitleAndContent = "titleAndContent"
	slideLayoutTitleOnly       = "titleOnly"
	slideLayoutBlank           = "blank"
)

const slideHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
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

// SlideWithContent renders a title+bullets slide with optional table, chart, and images.
func SlideWithContent(
	title string,
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	table *TableSpec,
	chart *ChartFrame,
	images []ImageRef,
) string {
	return SlideWithLayout(slideLayoutTitleAndContent, title, bullets, bulletStyles, bulletRuns, table, chart, images)
}

// SlideWithLayout renders a slide using an explicit layout mode.
func SlideWithLayout(
	layout string,
	title string,
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	table *TableSpec,
	chart *ChartFrame,
	images []ImageRef,
) string {
	var b strings.Builder
	layoutMode := normalizeSlideLayoutMode(layout)
	b.WriteString(slideHeader)

	nextID := 2
	if layoutMode != slideLayoutBlank {
		b.WriteString(titleShape(title))
		nextID = 3
	}

	if table != nil {
		b.WriteString(tableShape(table, nextID))
		nextID++
	} else if layoutMode == slideLayoutTitleAndContent && len(bullets) > 0 {
		b.WriteString(contentShape(bullets, bulletStyles, bulletRuns, nextID))
		nextID++
	}

	if chart != nil {
		b.WriteString(chartFrameShape(chart, nextID))
		nextID++
	}

	for i, image := range images {
		b.WriteString(imageShape(image, nextID+i))
	}
	b.WriteString(slideFooter)
	return b.String()
}

// SlideRelationships renders ppt/slides/_rels/slideN.xml.rels.
type ChartRel struct {
	RID    string
	Target string
}

func SlideRelationships(imageTargets []string, chartRel *ChartRel) string {
	return SlideRelationshipsWithLayout("../slideLayouts/slideLayout1.xml", imageTargets, chartRel)
}

func SlideRelationshipsWithLayout(layoutTarget string, imageTargets []string, chartRel *ChartRel) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="` + Escape(layoutTarget) + `"/>`)
	for i, target := range imageTargets {
		rid := i + 2
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="%s"/>`, rid, Escape(target)))
	}
	if chartRel != nil {
		b.WriteString(fmt.Sprintf(`
<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart" Target="%s"/>`,
			Escape(chartRel.RID),
			Escape(chartRel.Target),
		))
	}
	b.WriteString(`
</Relationships>`)
	return b.String()
}

func normalizeSlideLayoutMode(layout string) string {
	switch strings.ToLower(strings.TrimSpace(layout)) {
	case slideLayoutTitleOnly, "title_only", "title-only", "titleonly":
		return slideLayoutTitleOnly
	case slideLayoutBlank:
		return slideLayoutBlank
	default:
		return slideLayoutTitleAndContent
	}
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

func contentShape(bullets []string, bulletStyles []BulletParagraphSpec, bulletRuns [][]TextRunSpec, shapeID int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="%d" name="Content"/>
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
<a:lstStyle/>`, shapeID))

	for i, bullet := range bullets {
		style := bulletStyleAt(bulletStyles, i)
		runs := bulletRunsAt(bulletRuns, i)
		if len(runs) > 0 {
			b.WriteString(bulletParagraphRuns(runs, style))
			continue
		}
		b.WriteString(bulletParagraph(bullet, style))
	}

	b.WriteString(`
</p:txBody>
</p:sp>`)
	return b.String()
}

func bulletParagraph(text string, style BulletParagraphSpec) string {
	escaped := Escape(text)
	return `
<a:p>
` + bulletParagraphPropsXML(style) + `
<a:r>
<a:rPr lang="en-US" sz="2800" b="0" i="0" dirty="0"/>
<a:t>` + escaped + `</a:t>
</a:r>
</a:p>`
}
