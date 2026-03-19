package pptxxml

import (
	"strconv"
	"time"
)

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
