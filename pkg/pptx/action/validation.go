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
	pathValue = strings.TrimSpace(pathValue)
	if pathValue == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	parsed, err := url.Parse(pathValue)
	if err != nil {
		return fmt.Errorf("invalid file path %q: %w", pathValue, err)
	}
	pathPart := pathValue
	if parsed.Scheme == "" {
		pathPart = pathValue
	} else if len(parsed.Scheme) == 1 && filepath.IsAbs(pathValue) {
		// Windows drive prefix (for example `C:\...`) is not a URI scheme.
		pathPart = pathValue
	} else {
		if !strings.EqualFold(parsed.Scheme, "file") {
			return fmt.Errorf("file path must use file:// scheme or no scheme, got %q", parsed.Scheme)
		}
		if parsed.RawQuery != "" || parsed.Fragment != "" {
			return fmt.Errorf("file URI must not include query or fragment components")
		}
		if parsed.Opaque != "" {
			return fmt.Errorf("file URI uses unsupported opaque path form")
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
		return fmt.Errorf("file path cannot be empty")
	}
	if strings.ContainsRune(pathPart, '\x00') {
		return fmt.Errorf("file path cannot contain NUL characters")
	}

	normalized := strings.ReplaceAll(pathPart, "\\", "/")
	if containsTraversalSegment(normalized) {
		return fmt.Errorf("file path cannot contain directory traversal (..)")
	}
	if isSensitiveSystemPath(normalized) {
		return fmt.Errorf("file path references a restricted system location")
	}
	return nil
}

func containsTraversalSegment(path string) bool {
	for _, seg := range strings.Split(path, "/") {
		if seg == ".." {
			return true
		}
	}
	return false
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
