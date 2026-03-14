package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type (
	// Hyperlink describes a clickable link.
	Hyperlink = action.Hyperlink
	// HyperlinkAction defines the target of a hyperlink.
	HyperlinkAction = action.HyperlinkAction
	// HyperlinkActionType defines the type of hyperlink action.
	HyperlinkActionType = action.HyperlinkActionType
)

const (
	HyperlinkActionURL           = action.HyperlinkActionURL
	HyperlinkActionSlide         = action.HyperlinkActionSlide
	HyperlinkActionFirstSlide    = action.HyperlinkActionFirstSlide
	HyperlinkActionLastSlide     = action.HyperlinkActionLastSlide
	HyperlinkActionNextSlide     = action.HyperlinkActionNextSlide
	HyperlinkActionPreviousSlide = action.HyperlinkActionPreviousSlide
	HyperlinkActionEndShow       = action.HyperlinkActionEndShow
	HyperlinkActionEmail         = action.HyperlinkActionEmail
	HyperlinkActionFile          = action.HyperlinkActionFile
)

func NewHyperlink(act action.HyperlinkAction) Hyperlink {
	return action.NewHyperlink(act)
}

func HyperlinkURL(urlStr string) HyperlinkAction {
	return action.HyperlinkURL(urlStr)
}

func HyperlinkSlide(slideNum uint32) HyperlinkAction {
	return action.HyperlinkSlide(slideNum)
}

func HyperlinkFirstSlide() HyperlinkAction {
	return action.HyperlinkFirstSlide()
}

func HyperlinkLastSlide() HyperlinkAction {
	return action.HyperlinkLastSlide()
}

func HyperlinkNextSlide() HyperlinkAction {
	return action.HyperlinkNextSlide()
}

func HyperlinkPreviousSlide() HyperlinkAction {
	return action.HyperlinkPreviousSlide()
}

func HyperlinkEndShow() HyperlinkAction {
	return action.HyperlinkEndShow()
}

func HyperlinkEmail(address string) HyperlinkAction {
	return action.HyperlinkEmail(address)
}

func HyperlinkEmailWithSubject(address, subject string) HyperlinkAction {
	return action.HyperlinkEmailWithSubject(address, subject)
}

func HyperlinkFile(path string) HyperlinkAction {
	return action.HyperlinkFile(path)
}

func HyperlinkProgram(path string) HyperlinkAction {
	return action.HyperlinkProgram(path)
}

func validateHyperlinkAction(a HyperlinkAction, context string) error {
	return elements.ValidateHyperlinkAction(a, context)
}
