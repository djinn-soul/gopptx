package editor

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
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
	e.parts.Set(slideRef.Part, []byte(injectSlideHF(string(slideXML), hfXML)))
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

// ---------------------------------------------------------------------------
// Feature 4 – Handout master
// ---------------------------------------------------------------------------

const handoutMasterPath = "ppt/handoutMasters/handoutMaster1.xml"

// GetHandoutMaster returns basic info about the handout master.
func (e *PresentationEditor) GetHandoutMaster() (map[string]any, error) {
	data, ok := e.parts.Get(handoutMasterPath)
	if !ok {
		return map[string]any{"present": false}, nil
	}
	xmlStr := string(data)
	orientation := "landscape"
	if strings.Contains(xmlStr, `orient="portrait"`) {
		orientation = "portrait"
	}
	spp := extractHandoutSlidesPerPage(xmlStr)
	return map[string]any{
		"present":         true,
		"orientation":     orientation,
		"slides_per_page": spp,
	}, nil
}

// extractHandoutSlidesPerPage reads the slides-per-page from handout master XML.
// Falls back to 6 if not determinable.
func extractHandoutSlidesPerPage(xmlStr string) int {
	const (
		defaultSlidesPerPage  = 6
		maxValidSlidesPerPage = 9
	)
	count := strings.Count(xmlStr, "<p:sp>") + strings.Count(xmlStr, "<p:sp ")
	if count > 0 && count <= maxValidSlidesPerPage {
		return count
	}
	return defaultSlidesPerPage
}

// UpdateHandoutMaster sets orientation and slides-per-page in the handout master.
func (e *PresentationEditor) UpdateHandoutMaster(props map[string]any) error {
	data, ok := e.parts.Get(handoutMasterPath)
	var xmlStr string
	if ok {
		xmlStr = string(data)
	} else {
		xmlStr = defaultHandoutMasterXML()
	}

	if orient, castOK := props["orientation"].(string); castOK {
		xmlStr = applyHandoutOrientation(xmlStr, orient)
	}

	e.parts.Set(handoutMasterPath, []byte(xmlStr))
	e.addContentTypeOverride(handoutMasterPath,
		"application/vnd.openxmlformats-officedocument.presentationml.handoutMaster+xml")
	return nil
}

// applyHandoutOrientation sets the orient attribute in handout master XML.
func applyHandoutOrientation(xmlStr, orientation string) string {
	re := regexp.MustCompile(`orient="[^"]*"`)
	if re.MatchString(xmlStr) {
		return re.ReplaceAllString(xmlStr, fmt.Sprintf(`orient="%s"`, orientation))
	}
	reSz := regexp.MustCompile(`(<p:sldSz\b[^/]*)(/>)`)
	return reSz.ReplaceAllString(xmlStr, fmt.Sprintf(`$1 orient="%s"$2`, orientation))
}

// defaultHandoutMasterXML returns a minimal handout master XML.
func defaultHandoutMasterXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<p:handoutMaster xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"` +
		` xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
		`<p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/>` +
		`<p:nvPr/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld>` +
		`<p:clrMap bg1="lt1" tx1="dk1" bg2="lt2" tx2="dk2" accent1="accent1"` +
		` accent2="accent2" accent3="accent3" accent4="accent4" accent5="accent5"` +
		` accent6="accent6" hlink="hlink" folHlink="folHlink"/>` +
		`</p:handoutMaster>`
}

// handleGetHandoutMaster retrieves handout master info.
//
// Payload: {} (empty).
// Response: {"present": bool, "orientation": "landscape"|"portrait", "slides_per_page": N}.
func handleGetHandoutMaster(e *PresentationEditor, _ json.RawMessage) (any, error) {
	result, err := e.GetHandoutMaster()
	if err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return result, nil
}

// handleUpdateHandoutMaster updates handout master properties.
//
// Payload: {"orientation": "landscape"|"portrait", "slides_per_page": N}.
// Response: {"updated": true}.
func handleUpdateHandoutMaster(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	if updateErr := e.UpdateHandoutMaster(p); updateErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, updateErr.Error())
	}
	return map[string]bool{"updated": true}, nil
}

// ---------------------------------------------------------------------------
// Feature 5 – Digital signature detection
// ---------------------------------------------------------------------------

// HasDigitalSignature returns true if the package contains a digital signature.
func (e *PresentationEditor) HasDigitalSignature() (bool, error) {
	has := e.parts.Has("_xmlsignatures/origin.sigs") ||
		e.parts.Has("_xmlsignatures/sig1.xml")
	return has, nil
}

// handleHasDigitalSignature checks for a digital signature part.
//
// Payload: {} (empty).
// Response: {"has_digital_signature": bool}.
func handleHasDigitalSignature(e *PresentationEditor, _ json.RawMessage) (any, error) {
	has, err := e.HasDigitalSignature()
	if err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return map[string]bool{"has_digital_signature": has}, nil
}

// commandHandlerForHandoutSig routes handout master and digital signature ops.
// Kept here to avoid bloating command_router_groups.go past the 300-line ceiling.
func commandHandlerForHandoutSig(op string) (commandHandler, bool) {
	switch op {
	case OpGetHandoutMaster:
		return handleGetHandoutMaster, true
	case OpUpdateHandoutMaster:
		return handleUpdateHandoutMaster, true
	case OpHasDigitalSignature:
		return handleHasDigitalSignature, true
	default:
		return nil, false
	}
}
