package pptxxml

import (
	"fmt"
)

func placeholderShape(ph PlaceholderOverrideSpec, id int) string {
	phAttr := fmt.Sprintf(` idx="%d"`, ph.Index)
	if ph.Type != "" {
		phAttr += fmt.Sprintf(` type="%s"`, ph.Type)
	}

	if ph.Image != nil {
		// Render as Picture
		// If we have custom placement, use it, otherwise omit to inherit
		xfrm := `
  <p:spPr>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
  </p:spPr>`
		if ph.Image.X != 0 || ph.Image.Y != 0 || ph.Image.CX != 0 || ph.Image.CY != 0 {
			xfrm = fmt.Sprintf(`
  <p:spPr>
    <a:xfrm>
      <a:off x="%d" y="%d"/>
      <a:ext cx="%d" cy="%d"/>
    </a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
  </p:spPr>`, ph.Image.X, ph.Image.Y, ph.Image.CX, ph.Image.CY)
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
</p:pic>`, id, Escape(ph.Image.Name), phAttr, ph.Image.RelID, xfrm)
	}

	if ph.Table != nil {
		// Render as graphicFrame
		xfrm := ""
		if ph.Table.X != 0 || ph.Table.Y != 0 || ph.Table.CX != 0 || ph.Table.CY != 0 {
			xfrm = fmt.Sprintf(`
  <p:xfrm>
    <a:off x="%d" y="%d"/>
    <a:ext cx="%d" cy="%d"/>
  </p:xfrm>`, ph.Table.X, ph.Table.Y, ph.Table.CX, ph.Table.CY)
		}

		return fmt.Sprintf(`
<p:graphicFrame>
  <p:nvGraphicFramePr>
    <p:cNvPr id="%d" name="Placeholder Table %d"/>
    <p:cNvGraphicFramePr><a:graphicFrameLocks noGrp="1"/></p:cNvGraphicFramePr>
    <p:nvPr>
      <p:ph%s/>
    </p:nvPr>
  </p:nvGraphicFramePr>
  %s
  %s
</p:graphicFrame>`, id, ph.Index, phAttr, xfrm, tableGraphicXML(ph.Table))
	}

	if ph.Chart != nil {
		// Render as chart graphicFrame
		xfrm := ""
		if ph.Chart.X != 0 || ph.Chart.Y != 0 || ph.Chart.CX != 0 || ph.Chart.CY != 0 {
			xfrm = fmt.Sprintf(`
  <p:xfrm>
    <a:off x="%d" y="%d"/>
    <a:ext cx="%d" cy="%d"/>
  </p:xfrm>`, ph.Chart.X, ph.Chart.Y, ph.Chart.CX, ph.Chart.CY)
		}

		return fmt.Sprintf(`
<p:graphicFrame>
  <p:nvGraphicFramePr>
    <p:cNvPr id="%d" name="Placeholder Chart %d"/>
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
</p:graphicFrame>`, id, ph.Index, phAttr, xfrm, Escape(ph.Chart.RelID))
	}

	// Default to Text/Shape

	// Text body
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
