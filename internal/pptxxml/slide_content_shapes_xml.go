package pptxxml

import (
	"fmt"
	"strings"
)

func titleShape(title TitleSpec, width, height int64) string {
	// Standard margin is 0.5 inches (457200 EMU)
	margin := int64(457200)
	x := margin
	y := int64(274638) // Fixed top offset
	cx := width - 2*margin
	cy := int64(1143000) // Fixed height (1.25 inches)
	return titleShapeAt(title, x, y, cx, cy, "l")
}

func centeredTitleShape(title TitleSpec, width, height int64) string {
	// Standard margin is 0.5 inches (457200 EMU)
	margin := int64(457200)
	cx := width - 2*margin
	cy := int64(1371600) // 1.5 inches
	x := margin
	y := (height - cy) / 2 // Vertically centered
	return titleShapeAt(title, x, y, cx, cy, "ctr")
}

func titleShapeAt(title TitleSpec, x int64, y int64, cx int64, cy int64, align string) string {
	escaped := Escape(title.Text)
	sz := 4400
	if title.SizePt > 0 {
		sz = title.SizePt * 100
	}

	colorXML := ""
	if title.Color != "" {
		colorXML = fmt.Sprintf(`<a:solidFill><a:srgbClr val="%s"/></a:solidFill>`, Escape(title.Color))
	}

	return fmt.Sprintf(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="2" name="Title"/>
<p:cNvSpPr txBox="1"/>
<p:nvPr/>
</p:nvSpPr>
<p:spPr>
<a:xfrm>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</a:xfrm>
<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
<a:noFill/>
</p:spPr>
<p:txBody>
<a:bodyPr wrap="square" rtlCol="0" anchor="ctr"/>
<a:lstStyle/>
<a:p>
<a:pPr algn="%s"/>
<a:r>
<a:rPr lang="en-US" sz="%d" b="%s" i="%s" u="%s" dirty="0">%s</a:rPr>
<a:t>`+escaped+`</a:t>
</a:r>
</a:p>
</p:txBody>
</p:sp>`, x, y, cx, cy, Escape(align), sz, boolToFlag(title.Bold), boolToFlag(title.Italic), runUnderlineValue(title.Underline), colorXML)
}

func contentShape(bullets []string, bulletStyles []BulletParagraphSpec, bulletRuns [][]TextRunSpec, style ContentStyleSpec, shapeID int, width, height int64) string {
	margin := int64(457200)
	x := margin
	y := int64(1600200) // Fixed top offset
	cx := width - 2*margin
	cy := int64(4572000) // Fixed height
	return contentShapeAt(
		shapeID,
		"Content",
		x, y, cx, cy,
		bullets,
		bulletStyles,
		bulletRuns,
		style,
	)
}

func bigContentShape(bullets []string, bulletStyles []BulletParagraphSpec, bulletRuns [][]TextRunSpec, style ContentStyleSpec, shapeID int, width, height int64) string {
	margin := int64(457200)
	x := margin
	y := int64(1189200) // Lower top offset for big content
	cx := width - 2*margin
	cy := int64(5668800) // Taller content area
	return contentShapeAt(
		shapeID,
		"Content",
		x, y, cx, cy,
		bullets,
		bulletStyles,
		bulletRuns,
		style,
	)
}

func leftTwoColumnShape(bullets []string, bulletStyles []BulletParagraphSpec, bulletRuns [][]TextRunSpec, style ContentStyleSpec, shapeID int, width, height int64) string {
	margin := int64(457200)
	columnGap := int64(457200) // 0.5 inch gap
	x := margin
	y := int64(1189200)
	cx := (width - 2*margin - columnGap) / 2
	cy := int64(5668800)
	return contentShapeAt(
		shapeID,
		"Left Content",
		x, y, cx, cy,
		bullets,
		bulletStyles,
		bulletRuns,
		style,
	)
}

func rightTwoColumnShape(bullets []string, bulletStyles []BulletParagraphSpec, bulletRuns [][]TextRunSpec, style ContentStyleSpec, shapeID int, width, height int64) string {
	margin := int64(457200)
	columnGap := int64(457200)
	cx := (width - 2*margin - columnGap) / 2
	x := margin + cx + columnGap
	y := int64(1189200)
	cy := int64(5668800)
	return contentShapeAt(
		shapeID,
		"Right Content",
		x, y, cx, cy,
		bullets,
		bulletStyles,
		bulletRuns,
		style,
	)
}

func contentShapeAt(
	shapeID int,
	shapeName string,
	x int64,
	y int64,
	cx int64,
	cy int64,
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	style ContentStyleSpec,
) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="%d" name="%s"/>
<p:cNvSpPr txBox="1"/>
<p:nvPr/>
</p:nvSpPr>
<p:spPr>
<a:xfrm>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</a:xfrm>
<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
<a:noFill/>
</p:spPr>
<p:txBody>
<a:bodyPr wrap="square" rtlCol="0"/>
<a:lstStyle/>`, shapeID, Escape(shapeName), x, y, cx, cy))

	for i, bullet := range bullets {
		pStyle := bulletStyleAt(bulletStyles, i)
		runs := bulletRunsAt(bulletRuns, i)
		if len(runs) > 0 {
			b.WriteString(bulletParagraphRuns(runs, pStyle, style))
			continue
		}
		b.WriteString(bulletParagraph(bullet, pStyle, style))
	}

	b.WriteString(`
</p:txBody>
</p:sp>`)
	return b.String()
}

func splitBulletsForTwoColumns(bullets []string) ([]string, []string) {
	if len(bullets) == 0 {
		return nil, nil
	}
	mid := (len(bullets) + 1) / 2
	return bullets[:mid], bullets[mid:]
}

func splitBulletStylesForTwoColumns(styles []BulletParagraphSpec, leftCount int) ([]BulletParagraphSpec, []BulletParagraphSpec) {
	if len(styles) == 0 {
		return nil, nil
	}
	if leftCount > len(styles) {
		leftCount = len(styles)
	}
	left := styles[:leftCount]
	right := styles[leftCount:]
	return left, right
}

func splitBulletRunsForTwoColumns(runs [][]TextRunSpec, leftCount int) ([][]TextRunSpec, [][]TextRunSpec) {
	if len(runs) == 0 {
		return nil, nil
	}
	if leftCount > len(runs) {
		leftCount = len(runs)
	}
	left := runs[:leftCount]
	right := runs[leftCount:]
	return left, right
}

func bulletParagraph(text string, pStyle BulletParagraphSpec, style ContentStyleSpec) string {
	escaped := Escape(text)
	sz := 2800
	if style.SizePt > 0 {
		sz = style.SizePt * 100
	}

	colorXML := ""
	if style.Color != "" {
		colorXML = fmt.Sprintf(`<a:solidFill><a:srgbClr val="%s"/></a:solidFill>`, Escape(style.Color))
	}

	return fmt.Sprintf(`
<a:p>
%s
<a:r>
<a:rPr lang="en-US" sz="%d" b="%s" i="%s" u="%s" dirty="0">%s</a:rPr>
<a:t>%s</a:t>
</a:r>
</a:p>`, bulletParagraphPropsXML(pStyle), sz, boolToFlag(style.Bold), boolToFlag(style.Italic), runUnderlineValue(style.Underline), colorXML, escaped)
}
