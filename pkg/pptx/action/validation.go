package action

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"
)

var windowsDrivePathPattern = regexp.MustCompile(`^[A-Za-z]:[\\/].+`)

// ValidateHyperlinkAction ensures a hyperlink action is well-formed.
//
//nolint:gocognit // Validation enumerates many explicit policy branches for precise caller-facing errors.
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
		switch strings.ToLower(parsed.Scheme) {
		case "http", "https", "mailto", "ftp", "ftps":
			// allowed
		default:
			return fmt.Errorf("%s: hyperlink URL scheme %q is not allowed", context, parsed.Scheme)
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

func validateHyperlinkRenderableAction(action HyperlinkAction) error {
	switch action.Type {
	case HyperlinkActionURL:
		return validateRenderableHyperlinkURL(action.URL)
	case HyperlinkActionFile:
		return validateRenderableHyperlinkPath(
			action.FilePath,
			"hyperlink file path cannot be empty",
			"hyperlink file path cannot contain directory traversal (..)",
		)
	case HyperlinkActionProgram:
		return validateRenderableHyperlinkPath(
			action.ProgramPath,
			"hyperlink program path cannot be empty",
			"hyperlink program path cannot contain directory traversal (..)",
		)
	case HyperlinkActionEmail:
		return validateRenderableHyperlinkEmail(action.EmailAddress)
	case HyperlinkActionSlide, HyperlinkActionFirstSlide, HyperlinkActionLastSlide:
		return nil
	case HyperlinkActionNextSlide, HyperlinkActionPreviousSlide, HyperlinkActionEndShow:
		return nil
	default:
		return nil
	}
}

func validateRenderableHyperlinkURL(urlValue string) error {
	if urlValue == "" {
		return errors.New("hyperlink URL cannot be empty")
	}
	parsed, err := url.Parse(urlValue)
	if err != nil {
		return fmt.Errorf("invalid hyperlink URL %q: %w", urlValue, err)
	}
	if parsed.Scheme == "" {
		return fmt.Errorf("hyperlink URL %q must have a scheme (e.g., https://)", urlValue)
	}
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https", "mailto", "ftp", "ftps":
		return nil
	}
	return fmt.Errorf("hyperlink URL scheme %q is not allowed", parsed.Scheme)
}

func validateRenderableHyperlinkPath(
	pathValue string,
	emptyErr string,
	traversalErr string,
) error {
	if pathValue == "" {
		return errors.New(emptyErr)
	}
	if strings.Contains(pathValue, "..") {
		return errors.New(traversalErr)
	}
	return validateFilePathScheme(pathValue)
}

func validateRenderableHyperlinkEmail(address string) error {
	if address == "" {
		return errors.New("hyperlink email address cannot be empty")
	}
	return nil
}

func validateFilePathScheme(pathValue string) error {
	pathValue = strings.TrimSpace(pathValue)
	if pathValue == "" {
		return errors.New("file path cannot be empty")
	}
	parsed, err := url.Parse(pathValue)
	if err != nil {
		return fmt.Errorf("invalid file path %q: %w", pathValue, err)
	}
	var pathPart string
	switch {
	case parsed.Scheme == "":
		pathPart = pathValue
	case isWindowsDrivePath(pathValue, parsed.Scheme):
		// Windows drive prefix (for example `C:\...`) is not a URI scheme.
		pathPart = pathValue
	default:
		if !strings.EqualFold(parsed.Scheme, "file") {
			return fmt.Errorf("file path must use file:// scheme or no scheme, got %q", parsed.Scheme)
		}
		if parsed.RawQuery != "" || parsed.Fragment != "" {
			return errors.New("file URI must not include query or fragment components")
		}
		if parsed.Opaque != "" {
			return errors.New("file URI uses unsupported opaque path form")
		}
		host := strings.ToLower(strings.TrimSpace(parsed.Host))
		if host != "" && host != "localhost" {
			return fmt.Errorf("file URI host %q is not allowed", parsed.Host)
		}
		unescaped, unescapeErr := url.PathUnescape(parsed.Path)
		if unescapeErr != nil {
			return fmt.Errorf("invalid escaped file URI path %q: %w", parsed.Path, unescapeErr)
		}
		pathPart = unescaped
	}
	if pathPart == "" {
		return errors.New("file path cannot be empty")
	}
	if strings.ContainsRune(pathPart, '\x00') {
		return errors.New("file path cannot contain NUL characters")
	}

	normalized := strings.ReplaceAll(pathPart, "\\", "/")
	if containsTraversalSegment(normalized) {
		return errors.New("file path cannot contain directory traversal (..)")
	}
	if isSensitiveSystemPath(normalized) {
		return errors.New("file path references a restricted system location")
	}
	return nil
}

func isWindowsDrivePath(pathValue, scheme string) bool {
	if len(scheme) != 1 {
		return false
	}
	return windowsDrivePathPattern.MatchString(pathValue)
}

func containsTraversalSegment(path string) bool {
	return slices.Contains(strings.Split(path, "/"), "..")
}

func isSensitiveSystemPath(path string) bool {
	lower := strings.ToLower(strings.TrimSpace(path))
	lower = strings.TrimPrefix(lower, "file:///")
	lower = strings.TrimPrefix(lower, "file://")
	lower = strings.TrimPrefix(lower, "//localhost/")
	lower = strings.TrimPrefix(lower, "/")

	restrictedPrefixes := []string{
		"windows/", "windows/system32/",
		"programdata/",
		"etc/", "bin/", "sbin/", "usr/bin/", "usr/sbin/", "proc/", "sys/", "dev/",
	}
	for _, prefix := range restrictedPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}

	drivePrefixes := []string{
		"c:/windows/", "c:/windows/system32/", "c:/programdata/",
	}
	for _, prefix := range drivePrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}
