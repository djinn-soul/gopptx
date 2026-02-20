package urlfetch

import "errors"

// Sentinel errors for web2ppt operations.
var (
	// ErrNoContent is returned when no meaningful content is found on the page.
	ErrNoContent = errors.New("no content found on page")
)
