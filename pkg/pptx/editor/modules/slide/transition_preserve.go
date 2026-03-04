package slide

import (
	"errors"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type RawSlideTransition struct {
	XMLValue string
}

func (t RawSlideTransition) Validate() error {
	xml := strings.TrimSpace(t.XMLValue)
	if !strings.HasPrefix(xml, "<p:transition") || !strings.HasSuffix(xml, "</p:transition>") {
		return errors.New("transition XML must be wrapped in <p:transition>...</p:transition>")
	}
	return nil
}

func (t RawSlideTransition) XML() string {
	return strings.TrimSpace(t.XMLValue)
}

func PreserveExistingSlideTransition(
	getPart GetPartFn,
	slidePart string,
	slide elements.SlideContent,
) elements.SlideContent {
	if slide.Transition != nil {
		return slide
	}
	content, ok := getPart(slidePart)
	if !ok {
		return slide
	}

	transitionXML := ExtractSlideTransitionXML(string(content))
	if transitionXML == "" {
		return slide
	}
	return slide.WithTransition(RawSlideTransition{XMLValue: transitionXML})
}

func ExtractSlideTransitionXML(slideXML string) string {
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
