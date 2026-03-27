package editor

import (
	"errors"
	"fmt"
	"strings"
)

// SetSlideTransition applies a transition type to an existing slide.
// transitionType should be one of the Go transition constants (e.g. "fade", "push").
// durationMS is optional (0 = default). advanceMS is optional (-1 = disabled).
func (e *PresentationEditor) SetSlideTransition(
	slideIndex int,
	transitionType string,
	durationMS int,
	advanceMS int,
	disableAdvanceOnClick bool,
) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]

	data, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return fmt.Errorf("slide part %s not found", slideRef.Part)
	}

	tXML := buildTransitionXML(
		transitionType,
		durationMS,
		advanceMS,
		disableAdvanceOnClick,
	)
	slideXML := removeExistingTransitionXML(string(data))

	const closeSld = "</p:sld>"
	if !strings.Contains(slideXML, closeSld) {
		return errors.New("slide XML missing </p:sld> end tag")
	}
	updated := strings.Replace(slideXML, closeSld, tXML+closeSld, 1)
	e.parts.Set(slideRef.Part, []byte(updated))
	return nil
}

// buildTransitionXML generates a minimal p:transition element for the given type.
func buildTransitionXML(
	transitionType string,
	durationMS, advanceMS int,
	disableAdvanceOnClick bool,
) string {
	var b strings.Builder
	b.WriteString("<p:transition")
	if durationMS > 0 {
		fmt.Fprintf(
			&b,
			` xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main" p14:dur="%d"`,
			durationMS,
		)
	}
	if advanceMS >= 0 {
		fmt.Fprintf(&b, ` advTm="%d"`, advanceMS)
	}
	if disableAdvanceOnClick {
		b.WriteString(` advClick="0"`)
	}
	b.WriteString(">")

	if transitionType == "none" || transitionType == "cut" || transitionType == "" {
		b.WriteString("</p:transition>")
		return b.String()
	}
	b.WriteString(resolveTransitionElement(transitionType))
	b.WriteString("</p:transition>")
	return b.String()
}

func resolveTransitionElement(transitionType string) string {
	elementByType := map[string]string{
		"fade":      "<p:fade/>",
		"push":      `<p:push dir="r"/>`,
		"wipe":      `<p:wipe dir="r"/>`,
		"split":     `<p:split dir="out" orient="horz"/>`,
		"zoom":      `<p:zoom dir="in"/>`,
		"reveal":    `<p:reveal dir="r"/>`,
		"cover":     `<p:cover dir="r"/>`,
		"pull":      `<p:pull dir="r"/>`,
		"randomBar": "<p:randomBar/>",
		"wheel":     `<p:wheel spokes="4"/>`,
		"flash":     "<p:flash/>",
		"strips":    `<p:strips dir="ld"/>`,
		"blinds":    `<p:blinds dir="horz"/>`,
		"circle":    "<p:circle/>",
		"ripple":    "<p:ripple/>",
		"honeycomb": "<p:honeycomb/>",
		"glitter":   "<p:glitter/>",
		"vortex":    "<p:vortex/>",
		"shred":     "<p:shred/>",
		"switch":    "<p:switch/>",
		"flip":      "<p:flip/>",
		"gallery":   "<p:gallery/>",
		"cube":      "<p:cube/>",
		"doors":     "<p:doors/>",
		"box":       "<p:box/>",
		"random":    "<p:random/>",
	}
	if element, ok := elementByType[transitionType]; ok {
		return element
	}
	return "<p:" + transitionType + "/>"
}

// removeExistingTransitionXML strips any existing <p:transition>...</p:transition> block.
func removeExistingTransitionXML(slideXML string) string {
	start := strings.Index(slideXML, "<p:transition")
	if start < 0 {
		return slideXML
	}
	endTag := "</p:transition>"
	end := strings.Index(slideXML[start:], endTag)
	if end < 0 {
		return slideXML
	}
	return slideXML[:start] + slideXML[start+end+len(endTag):]
}
