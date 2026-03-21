package editor

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Feature 3 – Slide header/footer
// ---------------------------------------------------------------------------

// SlideHeaderFooter describes header/footer settings for a slide.
type SlideHeaderFooter struct {
	Footer       string
	ShowFooter   bool
	ShowSlideNum bool
	ShowDateTime bool
	DateTimeText string
}

// SetSlideHeaderFooter sets the <p:hf> element in the slide XML.
func (e *PresentationEditor) SetSlideHeaderFooter(slideIndex int, hf SlideHeaderFooter) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]
	slideXML, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return fmt.Errorf("slide part %q not found", slideRef.Part)
	}
	hfXML := buildHeaderFooterXML(hf)
	slideXMLStr := injectSlideHF(string(slideXML), hfXML)
	slideXMLStr = injectVisibleHeaderFooterShapes(
		slideXMLStr,
		hf,
		e.metadata.SlideSize.Width,
		e.metadata.SlideSize.Height,
		slideIndex+1,
	)
	e.parts.Set(slideRef.Part, []byte(slideXMLStr))
	return nil
}

// buildHeaderFooterXML creates the <p:hf> XML snippet.
func buildHeaderFooterXML(hf SlideHeaderFooter) string {
	sn := boolAttr(hf.ShowSlideNum)
	dt := boolAttr(hf.ShowDateTime)
	ftr := boolAttr(hf.ShowFooter)
	var b strings.Builder
	fmt.Fprintf(&b, `<p:hf sldNum="%s" dt="%s" ftr="%s">`, sn, dt, ftr)
	if hf.ShowFooter && hf.Footer != "" {
		fmt.Fprintf(&b, `<p:ftr><a:r><a:t>%s</a:t></a:r></p:ftr>`, xmlEscapeSimple(hf.Footer))
	}
	if hf.ShowDateTime && hf.DateTimeText != "" {
		fmt.Fprintf(&b, `<p:dt><a:r><a:t>%s</a:t></a:r></p:dt>`, xmlEscapeSimple(hf.DateTimeText))
	}
	b.WriteString(`</p:hf>`)
	return b.String()
}

// boolAttr converts bool to OOXML attribute string ("1"/"0").
func boolAttr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// injectSlideHF removes any existing <p:hf> and inserts the new one before </p:sld>.
func injectSlideHF(slideXML, hfXML string) string {
	reHF := regexp.MustCompile(`<p:hf[^>]*/?>(?:.*?</p:hf>)?`)
	slideXML = reHF.ReplaceAllString(slideXML, "")
	return strings.ReplaceAll(slideXML, "</p:sld>", hfXML+"</p:sld>")
}

func injectVisibleHeaderFooterShapes(slideXML string, hf SlideHeaderFooter, width, height int64, slideNumber int) string {
	overlayXML := buildVisibleHeaderFooterShapes(slideXML, hf, width, height, slideNumber)
	if overlayXML == "" {
		return slideXML
	}
	return strings.Replace(slideXML, "</p:spTree>", overlayXML+"</p:spTree>", 1)
}

func buildVisibleHeaderFooterShapes(slideXML string, hf SlideHeaderFooter, width, height int64, slideNumber int) string {
	nextID := maxShapeID(slideXML) + 1
	var b strings.Builder
	if hf.ShowSlideNum {
		b.WriteString(slideNumberOverlayShape(width, height, nextID))
		nextID++
	}
	if hf.ShowFooter && hf.Footer != "" {
		b.WriteString(footerOverlayShape(hf.Footer, width, height, nextID))
		nextID++
	}
	if hf.ShowDateTime {
		text := strings.TrimSpace(hf.DateTimeText)
		if text == "" {
			text = time.Now().Format("2006-01-02")
		}
		b.WriteString(dateTimeOverlayShape(text, height, nextID))
	}
	return b.String()
}

func maxShapeID(slideXML string) int {
	re := regexp.MustCompile(`<p:cNvPr[^>]*\bid="(\d+)"`)
	matches := re.FindAllStringSubmatch(slideXML, -1)
	maxID := 1
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		if id, err := strconv.Atoi(match[1]); err == nil && id > maxID {
			maxID = id
		}
	}
	return maxID
}

func slideNumberOverlayShape(width, height int64, shapeID int) string {
	cx := int64(548640)
	cy := int64(396240)
	x := width - cx - int64(457200)
	y := height - cy - int64(274320)
	return `
<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="` + strconv.Itoa(shapeID) + `" name="Slide Number Visible"/>
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

func footerOverlayShape(text string, width, height int64, shapeID int) string {
	cx := int64(2133600)
	cy := int64(396240)
	x := (width - cx) / 2
	y := height - cy - int64(274320)
	return `
<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="` + strconv.Itoa(shapeID) + `" name="Footer Visible"/>
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
        <a:t>` + xmlEscapeSimple(text) + `</a:t>
      </a:r>
    </a:p>
  </p:txBody>
</p:sp>`
}

func dateTimeOverlayShape(text string, height int64, shapeID int) string {
	cx := int64(2133600)
	cy := int64(396240)
	x := int64(457200)
	y := height - cy - int64(274320)
	return `
<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="` + strconv.Itoa(shapeID) + `" name="Date Visible"/>
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
        <a:t>` + xmlEscapeSimple(text) + `</a:t>
      </a:fld>
      <a:endParaRPr lang="en-US" dirty="0"/>
    </a:p>
  </p:txBody>
</p:sp>`
}

// handleSetSlideHeaderFooter sets header/footer on a slide.
//
// Payload: {"slide_index": N, "footer": "...", "show_footer": bool, ...}.
// Response: {"updated": true}.
func handleSetSlideHeaderFooter(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	hf := SlideHeaderFooter{}
	hf.Footer = v.OptionalString(p, "footer")
	hf.DateTimeText = v.OptionalString(p, "date_time_text")
	if sf, sfOK := v.OptionalBool(p, "show_footer"); sfOK {
		hf.ShowFooter = sf
	}
	if sn, snOK := v.OptionalBool(p, "show_slide_num"); snOK {
		hf.ShowSlideNum = sn
	}
	if sd, sdOK := v.OptionalBool(p, "show_date_time"); sdOK {
		hf.ShowDateTime = sd
	}

	if setErr := e.SetSlideHeaderFooter(slideIndex, hf); setErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, setErr.Error())
	}
	return map[string]bool{"updated": true}, nil
}
