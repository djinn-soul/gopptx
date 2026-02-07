package pptxxml

import (
	"fmt"
	"strings"
)

func titleShape(title string) string {
	return titleShapeAt(title, 457200, 274638, 8230200, 1143000, "l")
}

func centeredTitleShape(title string) string {
	return titleShapeAt(title, 457200, 2743200, 8230200, 1371600, "ctr")
}

func titleShapeAt(title string, x int64, y int64, cx int64, cy int64, align string) string {
	escaped := Escape(title)
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
<a:rPr lang="en-US" sz="4400" b="1" i="0" dirty="0"/>
<a:t>`+escaped+`</a:t>
</a:r>
</a:p>
</p:txBody>
</p:sp>`, x, y, cx, cy, Escape(align))
}

func contentShape(bullets []string, bulletStyles []BulletParagraphSpec, bulletRuns [][]TextRunSpec, shapeID int) string {
	return contentShapeAt(
		shapeID,
		"Content",
		457200,
		1600200,
		8230200,
		4572000,
		bullets,
		bulletStyles,
		bulletRuns,
	)
}

func bigContentShape(bullets []string, bulletStyles []BulletParagraphSpec, bulletRuns [][]TextRunSpec, shapeID int) string {
	return contentShapeAt(
		shapeID,
		"Content",
		457200,
		1189200,
		8230200,
		5668800,
		bullets,
		bulletStyles,
		bulletRuns,
	)
}

func leftTwoColumnShape(bullets []string, bulletStyles []BulletParagraphSpec, bulletRuns [][]TextRunSpec, shapeID int) string {
	return contentShapeAt(
		shapeID,
		"Left Content",
		457200,
		1189200,
		4115100,
		5668800,
		bullets,
		bulletStyles,
		bulletRuns,
	)
}

func rightTwoColumnShape(bullets []string, bulletStyles []BulletParagraphSpec, bulletRuns [][]TextRunSpec, shapeID int) string {
	return contentShapeAt(
		shapeID,
		"Right Content",
		4572300,
		1189200,
		4115100,
		5668800,
		bullets,
		bulletStyles,
		bulletRuns,
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
