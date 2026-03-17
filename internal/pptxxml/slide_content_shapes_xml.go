package pptxxml

import (
	"strconv"
	"strings"
	"time"
)

const (
	titleHeightEmu      = 1143000 // 1.25 inches
	contentHeightEmu    = 4572000 // 5 inches
	titleTopOffset      = 274638
	contentTopOffset    = 1600200
	titleShapeGrowCap   = 1024
	contentShapeGrowCap = 2048
)

func titleShape(title TitleSpec, width, _ int64) string {
	// Standard margin is 0.5 inches (457200 EMU)
	margin := defaultMargin
	x := int64(margin)
	y := int64(titleTopOffset) // Fixed top offset
	//nolint:mnd // Bi-lateral margin factor
	cx := width - 2*int64(margin)
	cy := int64(titleHeightEmu) // Fixed height (1.25 inches)
	align := title.Align
	if align == "" {
		align = "l"
	}
	return titleShapeAt(title, x, y, cx, cy, align)
}

//nolint:mnd // Layout constants from OOXML spec
func centeredTitleShape(title TitleSpec, width, height int64) string {
	// Standard margin is 0.5 inches (457200 EMU)
	margin := int64(457200)
	cx := width - 2*margin
	cy := int64(1371600) // 1.5 inches
	x := margin
	y := (height - cy) / 2 // Vertically centered
	align := title.Align
	if align == "" {
		align = "ctr"
	}
	return titleShapeAt(title, x, y, cx, cy, align)
}

func titleShapeAt(title TitleSpec, x int64, y int64, cx int64, cy int64, align string) string {
	escaped := Escape(title.Text)
	sz := 4400
	if title.SizePt > 0 {
		//nolint:mnd // Points to centipoints (1/100th of a point)
		sz = title.SizePt * 100
	}

	colorXML := ""
	if title.Color != "" {
		colorXML = `<a:solidFill><a:srgbClr val="` + Escape(title.Color) + `"/></a:solidFill>`
	}

	fontXML := ""
	if title.Font != "" {
		escFont := Escape(title.Font)
		fontXML = `<a:latin typeface="` + escFont + `"/><a:ea typeface="` + escFont + `"/><a:cs typeface="` + escFont + `"/>`
	}

	var b strings.Builder
	b.Grow(titleShapeGrowCap)
	b.WriteString(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="2" name="Title"/>
<p:cNvSpPr/>
<p:nvPr><p:ph type="title" idx="0"/></p:nvPr>
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
<a:bodyPr wrap="square" rtlCol="0" anchor="ctr"/>
<a:lstStyle/>
<a:p>
      <a:pPr algn="`)
	b.WriteString(Escape(align))
	b.WriteString(`"/>
      <a:r>
        <a:rPr lang="en-US" sz="`)
	b.WriteString(strconv.Itoa(sz))
	b.WriteString(`" b="`)
	b.WriteString(boolToFlag(title.Bold))
	b.WriteString(`" i="`)
	b.WriteString(boolToFlag(title.Italic))
	b.WriteString(`" u="`)
	b.WriteString(runUnderlineValue("", title.Underline))
	b.WriteString(`" dirty="0">`)
	b.WriteString(colorXML)
	b.WriteString(fontXML)
	b.WriteString(`</a:rPr>
        <a:t>`)
	b.WriteString(escaped)
	b.WriteString(`</a:t>
      </a:r>
    </a:p>
  </p:txBody>
</p:sp>`)
	return b.String()
}

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

func splitBulletsForTwoColumns(bullets []string) ([]string, []string) {
	if len(bullets) == 0 {
		return nil, nil
	}
	//nolint:mnd // Split point for two-column balancing
	mid := (len(bullets) + 1) / 2
	return bullets[:mid], bullets[mid:]
}

func splitBulletStylesForTwoColumns(
	styles []BulletParagraphSpec,
	leftCount int,
) ([]BulletParagraphSpec, []BulletParagraphSpec) {
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
		//nolint:mnd // Points to centipoints (1/100th of a point)
		sz = style.SizePt * 100
	}

	colorXML := ""
	if style.Color != "" {
		colorXML = `<a:solidFill><a:srgbClr val="` + Escape(style.Color) + `"/></a:solidFill>`
	}

	return `
<a:p>
` + bulletParagraphPropsXML(pStyle) + `
<a:r>
<a:rPr lang="en-US" sz="` + strconv.Itoa(sz) + `" b="` + boolToFlag(style.Bold) + `" i="` + boolToFlag(style.Italic) + `" u="` + runUnderlineValue("", style.Underline) + `" dirty="0">` + colorXML + `</a:rPr>
<a:t>` + escaped + `</a:t>
</a:r>
</a:p>`
}

//nolint:mnd // Layout constants from OOXML spec
func slideNumberShape(width, height int64, shapeID int) string {
	// Standard bottom right position for slide numbers
	cx := int64(548640)
	cy := int64(396240)
	x := width - cx - int64(457200)  // margin
	y := height - cy - int64(274320) // lower margin

	return `
<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="` + strconv.Itoa(shapeID) + `" name="Slide Number Placeholder"/>
    <p:cNvSpPr>
      <a:spLocks noGrp="1"/>
    </p:cNvSpPr>
    <p:nvPr>
      <p:ph type="sldNum" sz="quarter" idx="12"/>
    </p:nvPr>
  </p:nvSpPr>
  <p:spPr>
    <a:xfrm>
      <a:off x="` + strconv.FormatInt(x, 10) + `" y="` + strconv.FormatInt(y, 10) + `"/>
      <a:ext cx="` + strconv.FormatInt(cx, 10) + `" cy="` + strconv.FormatInt(cy, 10) + `"/>
    </a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
    <a:noFill/>
  </p:spPr>
  <p:txBody>
    <a:bodyPr wrap="square" rtlCol="0" anchor="ctr"/>
    <a:lstStyle/>
    <a:p>
      <a:pPr algn="r"/>
      <a:fld type="slidenum" id="{282E2E67-0C23-4552-87C9-2C764654F79F}">
        <a:rPr lang="en-US" smtClean="0"/>
        <a:t>‹#›</a:t>
      </a:fld>
      <a:endParaRPr lang="en-US" smtClean="0"/>
    </a:p>
  </p:txBody>
</p:sp>`
}

//nolint:mnd // Layout constants from OOXML spec
func footerShape(text string, width, height int64, shapeID int) string {
	cx := int64(2133600)
	cy := int64(396240)
	x := (width - cx) / 2
	y := height - cy - int64(274320)

	return `
<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="` + strconv.Itoa(shapeID) + `" name="Footer Placeholder"/>
    <p:cNvSpPr>
      <a:spLocks noGrp="1"/>
    </p:cNvSpPr>
    <p:nvPr>
      <p:ph type="ftr" sz="quarter" idx="11"/>
    </p:nvPr>
  </p:nvSpPr>
  <p:spPr>
    <a:xfrm>
      <a:off x="` + strconv.FormatInt(x, 10) + `" y="` + strconv.FormatInt(y, 10) + `"/>
      <a:ext cx="` + strconv.FormatInt(cx, 10) + `" cy="` + strconv.FormatInt(cy, 10) + `"/>
    </a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
    <a:noFill/>
  </p:spPr>
  <p:txBody>
    <a:bodyPr wrap="square" rtlCol="0" anchor="ctr"/>
    <a:lstStyle/>
    <a:p>
      <a:pPr algn="ctr"/>
      <a:r>
        <a:rPr lang="en-US" sz="1200" dirty="0"/>
        <a:t>` + Escape(text) + `</a:t>
      </a:r>
    </a:p>
  </p:txBody>
</p:sp>`
}

//nolint:mnd // Layout constants from OOXML spec
func dateTimeShape(_ int64, height int64, shapeID int) string {
	cx := int64(2133600)
	cy := int64(396240)
	x := int64(457200)
	y := height - cy - int64(274320)

	return `
<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="` + strconv.Itoa(shapeID) + `" name="Date Placeholder"/>
    <p:cNvSpPr>
      <a:spLocks noGrp="1"/>
    </p:cNvSpPr>
    <p:nvPr>
      <p:ph type="dt" sz="quarter" idx="10"/>
    </p:nvPr>
  </p:nvSpPr>
  <p:spPr>
    <a:xfrm>
      <a:off x="` + strconv.FormatInt(x, 10) + `" y="` + strconv.FormatInt(y, 10) + `"/>
      <a:ext cx="` + strconv.FormatInt(cx, 10) + `" cy="` + strconv.FormatInt(cy, 10) + `"/>
    </a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
    <a:noFill/>
  </p:spPr>
  <p:txBody>
    <a:bodyPr wrap="square" rtlCol="0" anchor="ctr"/>
    <a:lstStyle/>
    <a:p>
      <a:pPr algn="l"/>
      <a:fld type="datetime1" id="{A1B2C3D4-E5F6-7890-ABCD-EF1234567890}">
        <a:rPr lang="en-US" dirty="0"/>
        <a:t>` + time.Now().Format("2006-01-02") + `</a:t>
      </a:fld>
      <a:endParaRPr lang="en-US" dirty="0"/>
    </a:p>
  </p:txBody>
</p:sp>`
}
