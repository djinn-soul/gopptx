package urlfetch

import "errors"

// Sentinel errors for URL fetch conversion operations.
var (
	// ErrNoContent is returned when no meaningful content is found on the page.
	ErrNoContent = errors.New("no content found on page")
)
