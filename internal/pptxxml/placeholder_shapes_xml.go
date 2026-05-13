package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

func PlaceholderShape(ph PlaceholderOverrideSpec, id int) string {
	var pb strings.Builder
	pb.Grow(40)
	pb.WriteString(` idx="`)
	pb.WriteString(strconv.Itoa(ph.Index))
	pb.WriteString(`"`)
	phType := NormalizePlaceholderType(ph.Type)
	if phType != "" && phType != "obj" {
		pb.WriteString(` type="`)
		pb.WriteString(phType)
		pb.WriteString(`"`)
	}
	phAttr := pb.String()

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
<p:pic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
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
<p:graphicFrame xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
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
<p:graphicFrame xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
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
	if ph.Text != "" || ph.TextStyle != nil {
		escaped := Escape(ph.Text)
		textStyle := renderPlaceholderTextStyle(ph.TextStyle)
		txBody = fmt.Sprintf(`
<p:txBody>
  <a:bodyPr/>
  <a:lstStyle/>
  <a:p>
    %s<a:r>
      <a:rPr lang="en-US" dirty="0"/>
      <a:t>%s</a:t>
    </a:r>
  </a:p>
</p:txBody>`, textStyle, escaped)
	}

	return fmt.Sprintf(`
<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="%d" name="Placeholder %d"/>
    <p:cNvSpPr/>
    <p:nvPr>
      <p:ph%s/>
    </p:nvPr>
  </p:nvSpPr>
  %s
%s
</p:sp>`,
		id,
		ph.Index,
		phAttr,
		renderOverrideXfrm(ph.X, ph.Y, ph.CX, ph.CY, ph.GeometryXML, ph.ForceRectGeometry),
		txBody,
	)
}

func renderOverrideXfrm(x, y, cx, cy *int64, geometryXML string, forceRect *bool) string {
	hasBounds := x != nil || y != nil || cx != nil || cy != nil
	geomXML := strings.TrimSpace(geometryXML)
	if geomXML == "" || shouldForceRect(forceRect) {
		geomXML = `<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>`
	}
	if !hasBounds {
		if geomXML == "" {
			return "<p:spPr/>"
		}
		return fmt.Sprintf(`
  <p:spPr>
    %s
  </p:spPr>`, geomXML)
	}
	// Default to 0 if not specified but others are
	xv, yv, cxv, cyv := int64(0), int64(0), int64(0), int64(0)
	if x != nil {
		xv = *x
	}
	if y != nil {
		yv = *y
	}
	if cx != nil {
		cxv = *cx
	}
	if cy != nil {
		cyv = *cy
	}
	return fmt.Sprintf(`
  <p:spPr>
    <a:xfrm>
      <a:off x="%d" y="%d"/>
      <a:ext cx="%d" cy="%d"/>
    </a:xfrm>
    %s
  </p:spPr>`, xv, yv, cxv, cyv, geomXML)
}

func shouldForceRect(forceRect *bool) bool {
	return forceRect != nil && *forceRect
}

func renderPlaceholderTextStyle(ts *PlaceholderTextStyleSpec) string {
	if ts == nil {
		return ""
	}
	var b strings.Builder
	b.Grow(192)
	b.WriteString("<a:pPr")
	if ts.Align != nil {
		b.WriteString(` algn="`)
		b.WriteString(Escape(*ts.Align))
		b.WriteString(`"`)
	}
	b.WriteString(">")
	//nolint:nestif // Attribute emission is intentionally explicit per optional style field.
	if ts.Bold != nil || ts.Italic != nil || ts.SizePt != nil || ts.Color != nil || ts.Underline != nil ||
		ts.Font != nil {
		b.WriteString("<a:defRPr")
		if ts.Bold != nil {
			b.WriteString(` b="`)
			b.WriteString(boolToFlag(*ts.Bold))
			b.WriteString(`"`)
		}
		if ts.Italic != nil {
			b.WriteString(` i="`)
			b.WriteString(boolToFlag(*ts.Italic))
			b.WriteString(`"`)
		}
		if ts.SizePt != nil {
			b.WriteString(` sz="`)
			b.WriteString(strconv.Itoa(*ts.SizePt * ptFactor))
			b.WriteString(`"`)
		}
		if ts.Underline != nil {
			b.WriteString(` u="`)
			b.WriteString(Escape(*ts.Underline))
			b.WriteString(`"`)
		}
		b.WriteString(">")
		if ts.Color != nil {
			b.WriteString(`<a:solidFill><a:srgbClr val="`)
			b.WriteString(strings.TrimPrefix(*ts.Color, "#"))
			b.WriteString(`"/></a:solidFill>`)
		}
		if ts.Font != nil {
			b.WriteString(`<a:latin typeface="`)
			b.WriteString(Escape(*ts.Font))
			b.WriteString(`"/>`)
		}
		b.WriteString("</a:defRPr>")
	}
	b.WriteString("</a:pPr>")
	return b.String()
}

func NormalizePlaceholderType(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return "obj"
	}
	switch raw {
	case "picture", "pic":
		return "pic"
	case "title":
		return "title"
	case "body":
		return "body"
	case "ctrtitle", "centeredtitle", "centered_title":
		return "ctrTitle"
	default:
		return raw
	}
}
