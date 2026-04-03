package editor

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// GetSlideHeaderFooter reads the <p:hf> element from a slide and returns the parsed settings.
// If no <p:hf> element is present the returned struct has all fields at their zero values.
func (e *PresentationEditor) GetSlideHeaderFooter(slideIndex int) (SlideHeaderFooter, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return SlideHeaderFooter{}, fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]
	slideXML, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return SlideHeaderFooter{}, fmt.Errorf("slide part %q not found", slideRef.Part)
	}
	return parseHeaderFooterXML(string(slideXML)), nil
}

var (
	reHFTag    = regexp.MustCompile(`(?s)<p:hf\b([^>]*)>.*?</p:hf>|<p:hf\b([^>]*)/\s*>`)
	reHFAttr   = regexp.MustCompile(`\b(sldNum|dt|ftr)="([^"]*)"`)
	reHFText   = regexp.MustCompile(`(?s)<p:ftr>.*?<a:t>(.*?)</a:t>.*?</p:ftr>`)
	reHFDTText = regexp.MustCompile(`(?s)<p:dt>.*?<a:t>(.*?)</a:t>.*?</p:dt>`)
)

func parseHeaderFooterXML(slideXML string) SlideHeaderFooter {
	m := reHFTag.FindStringSubmatch(slideXML)
	if m == nil {
		return SlideHeaderFooter{}
	}
	attrs := m[1] + m[2]
	var hf SlideHeaderFooter
	for _, a := range reHFAttr.FindAllStringSubmatch(attrs, -1) {
		switch a[1] {
		case "sldNum":
			hf.ShowSlideNum = a[2] == "1"
		case "dt":
			hf.ShowDateTime = a[2] == "1"
		case "ftr":
			hf.ShowFooter = a[2] == "1"
		}
	}
	if tm := reHFText.FindStringSubmatch(slideXML); tm != nil {
		hf.Footer = tm[1]
	}
	if dm := reHFDTText.FindStringSubmatch(slideXML); dm != nil {
		hf.DateTimeText = dm[1]
	}
	return hf
}

// handleGetSlideHeaderFooter reads header/footer settings from a slide.
//
// Payload: {"slide_index": N}.
// Response: {"footer": "...", "show_footer": bool, "show_slide_num": bool,
//
//	"show_date_time": bool, "date_time_text": "..."}.
func handleGetSlideHeaderFooter(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	hf, getErr := e.GetSlideHeaderFooter(slideIndex)
	if getErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, getErr.Error())
	}
	return map[string]any{
		"footer":         hf.Footer,
		"show_footer":    hf.ShowFooter,
		"show_slide_num": hf.ShowSlideNum,
		"show_date_time": hf.ShowDateTime,
		"date_time_text": hf.DateTimeText,
	}, nil
}
