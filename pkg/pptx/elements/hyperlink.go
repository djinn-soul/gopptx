package elements

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
)

// ValidateHyperlinkAction check action for validity.
func ValidateHyperlinkAction(a action.HyperlinkAction, context string) error {
	return action.ValidateHyperlinkAction(a, context)
}
