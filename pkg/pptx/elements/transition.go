package elements

import (
	"fmt"
	"strings"
)

func ValidateSlideTransition(s SlideContent, index int) error {
	if s.Transition == nil {
		return nil
	}
	if err := s.Transition.Validate(); err != nil {
		return fmt.Errorf("slide %d transition: %w", index, err)
	}
	xml := strings.TrimSpace(s.Transition.XML())
	if xml == "" {
		return nil
	}
	if strings.HasPrefix(xml, "<p:transition") && strings.HasSuffix(xml, "</p:transition>") {
		return nil
	}
	if strings.HasPrefix(xml, "<mc:AlternateContent") && strings.HasSuffix(xml, "</mc:AlternateContent>") {
		return nil
	}
	return fmt.Errorf(
		"slide %d transition XML must be wrapped in <p:transition>...</p:transition> or <mc:AlternateContent>...</mc:AlternateContent>",
		index,
	)
}

func SlideTransitionXML(s SlideContent) string {
	if s.Transition == nil {
		return ""
	}
	return strings.TrimSpace(s.Transition.XML())
}
