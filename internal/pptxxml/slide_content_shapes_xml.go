package pptxxml

import (
	"strconv"
	"strings"
)

const (
	titleHeightEmu      = 1143000 // 1.25 inches
	contentHeightEmu    = 4572000 // 5 inches
	titleTopOffset      = 274638
	contentTopOffset    = 1600200
	titleShapeGrowCap   = 1024
	contentShapeGrowCap = 2048
)

// Precomputed title shape XML segments for the standard 16:9 slide (9144000 wide)
// with default TitleSpec (no custom size/color/font/bold/italic/underline/align).
// Avoids 5 strconv allocs + 2 builder allocs per slide in the common case.
//
//nolint:gochecknoglobals // read-only precomputed constants, never mutated
var (
	defaultTitleShapePrefix = `
<p:sp>
<p:nvSpPr>
<p:cNvPr id="2" name="Title"/>
<p:cNvSpPr/>
<p:nvPr><p:ph type="title" idx="0"/></p:nvPr>
</p:nvSpPr>
<p:spPr>
<a:xfrm>
<a:off x="457200" y="274638"/>
<a:ext cx="8229600" cy="1143000"/>
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
        <a:rPr lang="en-US" sz="4400" b="0" i="0" u="none" dirty="0"></a:rPr>
        <a:t>`
	defaultTitleShapeSuffix = `</a:t>
      </a:r>
    </a:p>
  </p:txBody>
</p:sp>`
)

func contentShape(
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	style ContentStyleSpec,
	shapeID int,
	width, _ int64,
) string {
	margin := defaultMargin
	x := int64(margin)
	y := int64(contentTopOffset) // Fixed top offset
	//nolint:mnd // Bi-lateral margin factor
	cx := width - 2*int64(margin)
	cy := int64(contentHeightEmu) // Fixed height
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

//nolint:mnd // Layout constants from OOXML spec
func bigContentShape(
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	style ContentStyleSpec,
	shapeID int,
	_, _ int64,
) string {
	margin := int64(457200)
	x := margin
	y := int64(1189200) // Lower top offset for big content
	cx := int64(8230200)
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

//nolint:mnd // Layout constants from OOXML spec
func leftTwoColumnShape(
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	style ContentStyleSpec,
	shapeID int,
	width, _ int64,
) string {
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

//nolint:mnd // Layout constants from OOXML spec
func rightTwoColumnShape(
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	style ContentStyleSpec,
	shapeID int,
	width, _ int64,
) string {
	margin := int64(457200)
	columnGap := int64(457200)
	cx := (width - 2*margin - columnGap) / 2
	x := int64(4572300)
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
	b.Grow(contentShapeGrowCap)
	vAlign := style.VAlign
	if vAlign == "" {
		vAlign = "t"
	}

	b.WriteString(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="`)
	b.WriteString(strconv.Itoa(shapeID))
	b.WriteString(`" name="`)
	b.WriteString(Escape(shapeName))
	b.WriteString(`"/>
<p:cNvSpPr/>
<p:nvPr><p:ph type="body" idx="1"/></p:nvPr>
</p:nvSpPr>
<p:spPr>
<a:xfrm>
<a:off x="`)
	b.WriteString(strconv.FormatInt(x, 10))
	b.WriteString(`" y="`)
	b.WriteString(strconv.FormatInt(y, 10))
	b.WriteString(`"/>
<a:ext cx="`)
	b.WriteString(strconv.FormatInt(cx, 10))
	b.WriteString(`" cy="`)
	b.WriteString(strconv.FormatInt(cy, 10))
	b.WriteString(`"/>
</a:xfrm>
<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
<a:noFill/>
</p:spPr>
<p:txBody>
<a:bodyPr wrap="square" rtlCol="0" anchor="`)
	b.WriteString(Escape(vAlign))
	b.WriteString(`"/>
<a:lstStyle/>`)

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
