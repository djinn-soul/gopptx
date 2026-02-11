package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type (
	// Hyperlink describes a clickable link.
	Hyperlink = elements.Hyperlink
	// HyperlinkAction defines the target of a hyperlink.
	HyperlinkAction = elements.HyperlinkAction
	// HyperlinkActionType defines the type of hyperlink action.
	HyperlinkActionType = elements.HyperlinkActionType
)

const (
	HyperlinkActionURL           = elements.HyperlinkActionURL
	HyperlinkActionSlide         = elements.HyperlinkActionSlide
	HyperlinkActionFirstSlide    = elements.HyperlinkActionFirstSlide
	HyperlinkActionLastSlide     = elements.HyperlinkActionLastSlide
	HyperlinkActionNextSlide     = elements.HyperlinkActionNextSlide
	HyperlinkActionPreviousSlide = elements.HyperlinkActionPreviousSlide
	HyperlinkActionEndShow       = elements.HyperlinkActionEndShow
	HyperlinkActionEmail         = elements.HyperlinkActionEmail
	HyperlinkActionFile          = elements.HyperlinkActionFile
)

func NewHyperlink(action HyperlinkAction) Hyperlink {
	return elements.NewHyperlink(action)
}

func HyperlinkURL(urlStr string) HyperlinkAction {
	return elements.HyperlinkURL(urlStr)
}

func HyperlinkSlide(slideNum uint32) HyperlinkAction {
	return elements.HyperlinkSlide(slideNum)
}

func HyperlinkFirstSlide() HyperlinkAction {
	return elements.HyperlinkFirstSlide()
}

func HyperlinkLastSlide() HyperlinkAction {
	return elements.HyperlinkLastSlide()
}

func HyperlinkNextSlide() HyperlinkAction {
	return elements.HyperlinkNextSlide()
}

func HyperlinkPreviousSlide() HyperlinkAction {
	return elements.HyperlinkPreviousSlide()
}

func HyperlinkEndShow() HyperlinkAction {
	return elements.HyperlinkEndShow()
}

func HyperlinkEmail(address string) HyperlinkAction {
	return elements.HyperlinkEmail(address)
}

func HyperlinkEmailWithSubject(address, subject string) HyperlinkAction {
	return elements.HyperlinkEmailWithSubject(address, subject)
}

func HyperlinkFile(path string) HyperlinkAction {
	return elements.HyperlinkFile(path)
}
