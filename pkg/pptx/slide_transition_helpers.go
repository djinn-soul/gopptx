package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func slideTransitionXML(s SlideContent) string {
	return elements.SlideTransitionXML(s)
}
