package pptxxml

import (
	"fmt"
	"strings"
)

func placeholderShape(ph PlaceholderOverrideSpec, id int) string {
	phAttr := fmt.Sprintf(` idx="%d"`, ph.Index)
	phType := normalizePlaceholderType(ph.Type)
	if phType != "" {
		phAttr += fmt.Sprintf(` type="%s"`, phType)
	}

	if ph.Image != nil {
		return renderPlaceholderImage(ph.Image, id, phAttr)
	}

	if ph.Table != nil {
		return renderPlaceholderTable(ph.Table, id, ph.Index, phAttr)
	}

	if ph.Chart != nil {
		return renderPlaceholderChart(ph.Chart, id, ph.Index, phAttr)
	}

	return renderPlaceholderDefault(ph, id, phAttr)
}

func renderPlaceholderImage(img *ImageRef, id int, phAttr string) string {
	xfrm := `
  <p:spPr>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
  </p:spPr>`
	if img.X != 0 || img.Y != 0 || img.CX != 0 || img.CY != 0 {
		xfrm = fmt.Sprintf(`
  <p:spPr>
    <a:xfrm>
      <a:off x="%d" y="%d"/>
      <a:ext cx="%d" cy="%d"/>
    </a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
  </p:spPr>`, img.X, img.Y, img.CX, img.CY)
	}

	return fmt.Sprintf(`
<p:pic>
  <p:nvPicPr>
    <p:cNvPr id="%d" name="%s"/>
    <p:cNvPicPr>
      <a:picLocks noChangeAspect="1"/>
    </p:cNvPicPr>
    <p:nvPr>
      <p:ph%s/>
    </p:nvPr>
  </p:nvPicPr>
  <p:blipFill>
    <a:blip r:embed="%s"/>
    <a:stretch>
      <a:fillRect/>
    </a:stretch>
  </p:blipFill>
  %s
</p:pic>`, id, Escape(img.Name), phAttr, FastEscapeRID(img.RelID), xfrm)
}

func renderPlaceholderTable(tbl *TableSpec, id int, index int, phAttr string) string {
	x := tbl.X
	y := tbl.Y
	cx := tbl.CX
	cy := tbl.CY
	if cx == 0 {
		cx = 5486400
	}
	if cy == 0 {
		cy = 2743200
	}
	xfrm := fmt.Sprintf(`
  <p:xfrm>
    <a:off x="%d" y="%d"/>
    <a:ext cx="%d" cy="%d"/>
  </p:xfrm>`, x, y, cx, cy)

return fmt.Sprintf(`
<p:graphicFrame>
  <p:nvGraphicFramePr>
    <p:cNvPr id="%d" name="Placeholder Table %d"%s/>
    <p:cNvGraphicFramePr><a:graphicFrameLocks noGrp="1"/></p:cNvGraphicFramePr>
    <p:nvPr>
      <p:ph%s/>
    </p:nvPr>
  </p:nvGraphicFramePr>
  %s
  %s
</p:graphicFrame>`, id, index, makeCNvPrAttrs(tbl.AltText, tbl.IsDecorative), phAttr, xfrm, tableGraphicXML(tbl))
}

func renderPlaceholderChart(ch *ChartFrame, id int, index int, phAttr string) string {
	x := ch.X
	y := ch.Y
	cx := ch.CX
	cy := ch.CY
	if cx == 0 {
		cx = 5486400
	}
	if cy == 0 {
		cy = 2743200
	}
	xfrm := fmt.Sprintf(`
  <p:xfrm>
    <a:off x="%d" y="%d"/>
    <a:ext cx="%d" cy="%d"/>
  </p:xfrm>`, x, y, cx, cy)

return fmt.Sprintf(`
<p:graphicFrame>
  <p:nvGraphicFramePr>
    <p:cNvPr id="%d" name="Placeholder Chart %d"%s/>
    <p:cNvGraphicFramePr/>
    <p:nvPr>
      <p:ph%s/>
    </p:nvPr>
  </p:nvGraphicFramePr>
  %s
  <a:graphic>
    <a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/chart">
      <c:chart xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" r:id="%s"/>
    </a:graphicData>
  </a:graphic>
</p:graphicFrame>`, id, index, makeCNvPrAttrs(ch.AltText, ch.IsDecorative), phAttr, xfrm, FastEscapeRID(ch.RelID))
}

func renderPlaceholderDefault(ph PlaceholderOverrideSpec, id int, phAttr string) string {
	txBody := `
<p:txBody>
  <a:bodyPr/>
  <a:lstStyle/>
  <a:p/>
</p:txBody>`
	if ph.Text != "" {
		escaped := Escape(ph.Text)
		txBody = fmt.Sprintf(`
<p:txBody>
  <a:bodyPr/>
  <a:lstStyle/>
  <a:p>
    <a:r>
      <a:rPr lang="en-US" dirty="0"/>
      <a:t>%s</a:t>
    </a:r>
  </a:p>
</p:txBody>`, escaped)
	}

	return fmt.Sprintf(`
<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="%d" name="Placeholder %d"/>
    <p:cNvSpPr txBox="1"/>
    <p:nvPr>
      <p:ph%s/>
    </p:nvPr>
  </p:nvSpPr>
  <p:spPr/>
%s
</p:sp>`, id, ph.Index, phAttr, txBody)
}

func normalizePlaceholderType(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "picture", "pic":
		return "pic"
	case "title":
		return "title"
	case "body":
		return "body"
	case "ctrtitle", "centeredtitle", "centered_title":
		return "ctrTitle"
	default:
		return strings.TrimSpace(raw)
	}
}
