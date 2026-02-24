// Package urlfetch converts web pages or raw HTML into PPTX presentations.
//
// It is a Go port of ppt-rs/src/web2ppt with improvements:
// HTML tables are rendered as real PPTX tables rather than summary bullets.
package urlfetch

// URLFetchConfig holds options that control content extraction and slide generation.
type URLFetchConfig struct { //nolint:revive // keeping exported name for API compatibility
	// MaxSlides caps the total number of generated slides.
	MaxSlides int
	// MaxBulletsPerSlide caps how many bullet points appear on a single slide.
	MaxBulletsPerSlide int
	// IncludeImages adds image alt-text bullets when true.
	IncludeImages bool
	// IncludeTables renders HTML tables as native PPTX tables when true.
	IncludeTables bool
	// IncludeCode adds code block text as bullets when true.
	IncludeCode bool
	// ExtractLinks adds standalone hyperlink text as bullets when true.
	ExtractLinks bool
	// GroupByHeadings groups content per heading (Grouped mode);
	// set to false for linear (overflow-based) mode.
	GroupByHeadings bool
	// UserAgent used for HTTP requests.
	UserAgent string
	// TimeoutSecs is the HTTP timeout in seconds.
	TimeoutSecs int
	// MaxBodyBytes caps the fetched HTTP response body size in bytes.
	MaxBodyBytes int64
}

// Web2PptConfig is a compatibility alias for URLFetchConfig.
type Web2PptConfig = URLFetchConfig

// DefaultConfig returns a URLFetchConfig with sensible defaults.
func DefaultConfig() URLFetchConfig {
	return URLFetchConfig{
		MaxSlides:          20,
		MaxBulletsPerSlide: 6,
		IncludeImages:      true,
		IncludeTables:      true,
		IncludeCode:        true,
		ExtractLinks:       true,
		GroupByHeadings:    true,
		UserAgent:          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		TimeoutSecs:        30,
		MaxBodyBytes:       10 * 1024 * 1024,
	}
}

// WithMaxSlides sets the maximum slide count.
func (c URLFetchConfig) WithMaxSlides(n int) URLFetchConfig {
	c.MaxSlides = n
	return c
}

// WithMaxBullets sets the maximum bullets-per-slide.
func (c URLFetchConfig) WithMaxBullets(n int) URLFetchConfig {
	c.MaxBulletsPerSlide = n
	return c
}

// WithImages enables or disables image extraction.
func (c URLFetchConfig) WithImages(v bool) URLFetchConfig {
	c.IncludeImages = v
	return c
}

// WithTables enables or disables table extraction.
func (c URLFetchConfig) WithTables(v bool) URLFetchConfig {
	c.IncludeTables = v
	return c
}

// WithCode enables or disables code-block extraction.
func (c URLFetchConfig) WithCode(v bool) URLFetchConfig {
	c.IncludeCode = v
	return c
}

// WithLinks enables or disables link extraction.
func (c URLFetchConfig) WithLinks(v bool) URLFetchConfig {
	c.ExtractLinks = v
	return c
}

// WithGroupByHeadings sets the slide-grouping strategy.
func (c URLFetchConfig) WithGroupByHeadings(v bool) URLFetchConfig {
	c.GroupByHeadings = v
	return c
}

// WithUserAgent overrides the HTTP User-Agent header.
func (c URLFetchConfig) WithUserAgent(ua string) URLFetchConfig {
	c.UserAgent = ua
	return c
}

// WithTimeout sets the HTTP request timeout in seconds.
func (c URLFetchConfig) WithTimeout(secs int) URLFetchConfig {
	c.TimeoutSecs = secs
	return c
}

// WithMaxBodyBytes sets the maximum response body size in bytes.
func (c URLFetchConfig) WithMaxBodyBytes(n int64) URLFetchConfig {
	c.MaxBodyBytes = n
	return c
}

// ConversionOptions control metadata and optional features applied during conversion.
type ConversionOptions struct {
	// Title overrides the page title for the cover slide.
	Title *string
	// Author sets the author metadata.
	Author *string
	// IncludeSourceURL adds a "Source: <url>" bullet to the title slide.
	IncludeSourceURL bool
	// AddPageNumbers adds slide-number shapes (future).
	AddPageNumbers bool
}

// DefaultConversionOptions returns ConversionOptions with sensible defaults.
func DefaultConversionOptions() ConversionOptions {
	return ConversionOptions{
		IncludeSourceURL: true,
	}
}

// WithTitle sets a custom presentation title.
func (o ConversionOptions) WithTitle(title string) ConversionOptions {
	o.Title = &title
	return o
}

// WithAuthor sets the author field.
func (o ConversionOptions) WithAuthor(author string) ConversionOptions {
	o.Author = &author
	return o
}

// WithSourceURL controls whether the source URL appears on the title slide.
func (o ConversionOptions) WithSourceURL(v bool) ConversionOptions {
	o.IncludeSourceURL = v
	return o
}

// WithPageNumbers controls slide-number display.
func (o ConversionOptions) WithPageNumbers(v bool) ConversionOptions {
	o.AddPageNumbers = v
	return o
}
