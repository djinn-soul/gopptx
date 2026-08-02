package editor

import (
	"errors"
	"fmt"
	"strings"
)

// transitionNone is the transition type that means "no transition".
const transitionNone = "none"

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

// Upper bounds, in milliseconds, for the "fast" and "med" ECMA-376 transition
// speeds. Anything longer is reported as "slow".
const (
	fastTransitionMaxMS   = 500
	mediumTransitionMaxMS = 1000
)

// buildTransitionXML generates a minimal p:transition element for the given type.
//
// A duration in milliseconds only exists as the PowerPoint 2010 extension
// attribute p14:dur, which is only legal inside an mc:Choice that requires p14,
// so a timed transition is emitted as an mc:AlternateContent pair with a
// spd-only fallback for readers that do not understand p14.
func buildTransitionXML(
	transitionType string,
	durationMS, advanceMS int,
	disableAdvanceOnClick bool,
) string {
	if durationMS <= 0 {
		return transitionElementXML(transitionType, 0, advanceMS, disableAdvanceOnClick, false)
	}

	var b strings.Builder
	b.WriteString(`<mc:AlternateContent xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006">`)
	b.WriteString(`<mc:Choice xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main" Requires="p14">`)
	b.WriteString(transitionElementXML(transitionType, durationMS, advanceMS, disableAdvanceOnClick, true))
	b.WriteString(`</mc:Choice><mc:Fallback>`)
	b.WriteString(transitionElementXML(transitionType, durationMS, advanceMS, disableAdvanceOnClick, false))
	b.WriteString(`</mc:Fallback></mc:AlternateContent>`)
	return b.String()
}

// transitionElementXML renders one <p:transition> element. withDuration adds the
// p14:dur attribute; the caller is responsible for placing that form inside an
// mc:Choice that declares the p14 prefix.
func transitionElementXML(
	transitionType string,
	durationMS, advanceMS int,
	disableAdvanceOnClick, withDuration bool,
) string {
	var b strings.Builder
	b.WriteString("<p:transition")
	if durationMS > 0 {
		fmt.Fprintf(&b, ` spd="%s"`, transitionSpeedForDuration(durationMS))
	}
	if withDuration {
		fmt.Fprintf(&b, ` p14:dur="%d"`, durationMS)
	}
	if advanceMS >= 0 {
		fmt.Fprintf(&b, ` advTm="%d"`, advanceMS)
	}
	if disableAdvanceOnClick {
		b.WriteString(` advClick="0"`)
	}
	b.WriteString(">")

	if transitionType == transitionNone || transitionType == "cut" || transitionType == "" {
		b.WriteString("</p:transition>")
		return b.String()
	}
	b.WriteString(resolveTransitionElement(transitionType))
	b.WriteString("</p:transition>")
	return b.String()
}

// transitionSpeedForDuration buckets a millisecond duration into the three
// ECMA-376 transition speeds, which is all a p14-unaware reader can honour.
func transitionSpeedForDuration(durationMS int) string {
	switch {
	case durationMS <= fastTransitionMaxMS:
		return "fast"
	case durationMS <= mediumTransitionMaxMS:
		return "med"
	default:
		return "slow"
	}
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

// removeExistingTransitionXML strips any existing transition block, including a
// timed transition wrapped in mc:AlternateContent. Without the wrapper case a
// second set_transition call would leave the previous Choice/Fallback pair
// behind and the slide would carry two transitions.
func removeExistingTransitionXML(slideXML string) string {
	if stripped, ok := removeTransitionAlternateContent(slideXML); ok {
		return stripped
	}
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

// removeTransitionAlternateContent removes the first mc:AlternateContent block
// that wraps a p:transition, reporting whether one was found.
func removeTransitionAlternateContent(slideXML string) (string, bool) {
	const openTag = "<mc:AlternateContent"
	const closeTag = "</mc:AlternateContent>"
	searchFrom := 0
	for {
		start := strings.Index(slideXML[searchFrom:], openTag)
		if start < 0 {
			return slideXML, false
		}
		start += searchFrom
		rel := strings.Index(slideXML[start:], closeTag)
		if rel < 0 {
			return slideXML, false
		}
		end := start + rel + len(closeTag)
		if strings.Contains(slideXML[start:end], "<p:transition") {
			return slideXML[:start] + slideXML[end:], true
		}
		searchFrom = end
	}
}
