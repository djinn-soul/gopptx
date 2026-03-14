package action

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

// ValidateHyperlinkAction ensures a hyperlink action is well-formed.
func ValidateHyperlinkAction(a HyperlinkAction, context string) error {
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
		if err := validateFilePathScheme(a.FilePath); err != nil {
			return fmt.Errorf("%s: %w", context, err)
		}

	case HyperlinkActionProgram:
		if a.ProgramPath == "" {
			return fmt.Errorf("%s: hyperlink program path cannot be empty", context)
		}
		if err := validateFilePathScheme(a.ProgramPath); err != nil {
			return fmt.Errorf("%s: %w", context, err)
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

func validateFilePathScheme(pathValue string) error {
	parsed, err := url.Parse(pathValue)
	if err != nil {
		return fmt.Errorf("invalid file path %q: %w", pathValue, err)
	}
	if parsed.Scheme == "" {
		return nil
	}
	if len(parsed.Scheme) == 1 && filepath.IsAbs(pathValue) {
		// Windows drive prefix (for example `C:\...`) is not a URI scheme.
		return nil
	}
	if !strings.EqualFold(parsed.Scheme, "file") {
		return fmt.Errorf("file path must use file:// scheme or no scheme, got %q", parsed.Scheme)
	}
	return nil
}
