package editor

import (
	"errors"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type rawSlideTransition struct {
	xml string
}

func (t rawSlideTransition) Validate() error {
	xml := strings.TrimSpace(t.xml)
	if !strings.HasPrefix(xml, "<p:transition") || !strings.HasSuffix(xml, "</p:transition>") {
		return errors.New("transition XML must be wrapped in <p:transition>...</p:transition>")
	}
	return nil
}

func (t rawSlideTransition) XML() string {
	return strings.TrimSpace(t.xml)
}

func preserveExistingSlideTransition(
	parts *PartStore,
	slidePart string,
	slide elements.SlideContent,
) elements.SlideContent {
	if slide.Transition != nil {
		return slide
	}
	content, ok := parts.Get(slidePart)
	if !ok {
		return slide
	}

	transitionXML := extractSlideTransitionXML(string(content))
	if transitionXML == "" {
		return slide
	}
	return slide.WithTransition(rawSlideTransition{xml: transitionXML})
}

func extractSlideTransitionXML(slideXML string) string {
	start := strings.Index(slideXML, "<p:transition")
	if start < 0 {
		return ""
	}
	endOffset := strings.Index(slideXML[start:], "</p:transition>")
	if endOffset < 0 {
		return ""
	}
	end := start + endOffset + len("</p:transition>")
	return strings.TrimSpace(slideXML[start:end])
}
