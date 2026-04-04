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
	reHFTag     = regexp.MustCompile(`(?s)<p:hf\b([^>]*)>.*?</p:hf>|<p:hf\b([^>]*)/\s*>`)
	reHFAttr    = regexp.MustCompile(`\b(sldNum|dt|ftr)="([^"]*)"`)
	reHFTextTag = regexp.MustCompile(`(?s)<p:ftr\b[^>]*>(.*?)</p:ftr>`)
	reHFDTTag   = regexp.MustCompile(`(?s)<p:dt\b[^>]*>(.*?)</p:dt>`)
	reATextTag  = regexp.MustCompile(`(?s)<a:t>(.*?)</a:t>`)
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
	hf.Footer = extractCombinedText(slideXML, reHFTextTag)
	hf.DateTimeText = extractCombinedText(slideXML, reHFDTTag)
	return hf
}

func extractCombinedText(slideXML string, containerPattern *regexp.Regexp) string {
	container := containerPattern.FindStringSubmatch(slideXML)
	if container == nil {
		return ""
	}
	textMatches := reATextTag.FindAllStringSubmatch(container[1], -1)
	if len(textMatches) == 0 {
		return ""
	}
	out := ""
	for _, match := range textMatches {
		if len(match) > 1 {
			out += match[1]
		}
	}
	return out
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
