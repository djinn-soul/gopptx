package pptxxml

import (
	"strconv"
	"strings"
)

const titleSideMarginCount = 2

func titleShape(title TitleSpec, width, _ int64) string {
	if width == 9144000 && title.SizePt == 0 && title.Color == "" && title.Font == "" &&
		!title.Bold && !title.Italic && !title.Underline && title.Align == "" {
		return defaultTitleShapePrefix + Escape(title.Text) + defaultTitleShapeSuffix
	}
	margin := defaultMargin
	x := int64(margin)
	y := int64(titleTopOffset)
	cx := width - titleSideMarginCount*int64(margin)
	cy := int64(titleHeightEmu)
	align := title.Align
	if align == "" {
		align = "l"
	}
	return titleShapeAt(title, x, y, cx, cy, align)
}

//nolint:mnd // Layout constants from OOXML spec
func centeredTitleShape(title TitleSpec, width, height int64) string {
	margin := int64(457200)
	cx := width - 2*margin
	cy := int64(1371600)
	x := margin
	y := (height - cy) / 2
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
		sz = title.SizePt * 100 //nolint:mnd // Points to centipoints
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
