package pptx

import (
	"fmt"
	"net/url"
	"strings"
)

func validateHyperlink(h *Hyperlink, context string) error {
	if h == nil {
		return nil
	}
	return validateHyperlinkAction(h.Action, context)
}

func validateHyperlinkAction(a HyperlinkAction, context string) error {
	switch a.Type {
	case HyperlinkActionURL:
		if a.URL == "" {
			return fmt.Errorf("%s: hyperlink URL cannot be empty", context)
		}
		parsed, err := url.Parse(a.URL)
		if err != nil {
			return fmt.Errorf("%s: invalid hyperlink URL %q: %w", context, a.URL, err)
		}
		if parsed.Scheme == "" {
			return fmt.Errorf("%s: hyperlink URL %q must have a scheme (e.g., https://)", context, a.URL)
		}

	case HyperlinkActionSlide:
		if a.SlideNumber == 0 {
			return fmt.Errorf("%s: hyperlink slide number must be >= 1", context)
		}

	case HyperlinkActionEmail:
		if a.EmailAddress == "" {
			return fmt.Errorf("%s: hyperlink email address cannot be empty", context)
		}
		if !strings.Contains(a.EmailAddress, "@") {
			return fmt.Errorf("%s: hyperlink email address %q must contain @", context, a.EmailAddress)
		}

	case HyperlinkActionFile:
		if a.FilePath == "" {
			return fmt.Errorf("%s: hyperlink file path cannot be empty", context)
		}

	case HyperlinkActionFirstSlide, HyperlinkActionLastSlide,
		HyperlinkActionNextSlide, HyperlinkActionPreviousSlide,
		HyperlinkActionEndShow:
		// No additional validation needed for navigation actions

	default:
		return fmt.Errorf("%s: unsupported hyperlink action type %q", context, a.Type)
	}
	return nil
}
