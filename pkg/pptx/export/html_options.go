package export

// ThemeColors represents CSS color variables for HTML export.
type ThemeColors struct {
	TitleColor      string
	BodyColor       string
	AccentColor     string
	BackgroundColor string
	SlideBackground string
}

// HTMLOptions configures how slides are exported to HTML.
type HTMLOptions struct {
	// Emit images as base64 data URIs
	EmbedImages bool
	// Prefix for sidecar asset paths (used if EmbedImages is false)
	BaseURL string
	// Override default styling
	Theme *ThemeColors
	// Include JS/buttons for previous/next slide navigation
	IncludeNavigation bool
	// Display the slide number overlay
	IncludeSlideNumbers bool
	// IncludeNotes prints the speaker notes under each slide. No export path
	// could emit notes at all before, so a deck's notes were lost on the way to
	// HTML even though the model carried them.
	IncludeNotes bool
}

// DefaultHTMLOptions returns the standard configuration for HTML export.
func DefaultHTMLOptions() HTMLOptions {
	return HTMLOptions{
		EmbedImages:         true,
		IncludeNavigation:   true,
		IncludeSlideNumbers: true,
	}
}
