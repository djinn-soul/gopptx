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
) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]

	data, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return fmt.Errorf("slide part %s not found", slideRef.Part)
	}

	tXML := buildTransitionXML(transitionType, durationMS, advanceMS)
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
//
//nolint:cyclop,funlen // Transition type switch is a flat 25-branch lookup; complexity is inherent.
func buildTransitionXML(transitionType string, durationMS, advanceMS int) string {
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
	b.WriteString(">")

	switch transitionType {
	case "none", "cut", "":
		b.WriteString("</p:transition>")
		return b.String()
	case "fade":
		b.WriteString("<p:fade/>")
	case "push":
		b.WriteString(`<p:push dir="r"/>`)
	case "wipe":
		b.WriteString(`<p:wipe dir="r"/>`)
	case "split":
		b.WriteString(`<p:split dir="out" orient="horz"/>`)
	case "zoom":
		b.WriteString(`<p:zoom dir="in"/>`)
	case "reveal":
		b.WriteString(`<p:reveal dir="r"/>`)
	case "cover":
		b.WriteString(`<p:cover dir="r"/>`)
	case "pull":
		b.WriteString(`<p:pull dir="r"/>`)
	case "randomBar":
		b.WriteString("<p:randomBar/>")
	case "wheel":
		b.WriteString(`<p:wheel spokes="4"/>`)
	case "flash":
		b.WriteString("<p:flash/>")
	case "strips":
		b.WriteString(`<p:strips dir="ld"/>`)
	case "blinds":
		b.WriteString(`<p:blinds dir="horz"/>`)
	case "circle":
		b.WriteString("<p:circle/>")
	case "ripple":
		b.WriteString("<p:ripple/>")
	case "honeycomb":
		b.WriteString("<p:honeycomb/>")
	case "glitter":
		b.WriteString("<p:glitter/>")
	case "vortex":
		b.WriteString("<p:vortex/>")
	case "shred":
		b.WriteString("<p:shred/>")
	case "switch":
		b.WriteString("<p:switch/>")
	case "flip":
		b.WriteString("<p:flip/>")
	case "gallery":
		b.WriteString("<p:gallery/>")
	case "cube":
		b.WriteString("<p:cube/>")
	case "doors":
		b.WriteString("<p:doors/>")
	case "box":
		b.WriteString("<p:box/>")
	case "random":
		b.WriteString("<p:random/>")
	default:
		fmt.Fprintf(&b, "<p:%s/>", transitionType)
	}
	b.WriteString("</p:transition>")
	return b.String()
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
