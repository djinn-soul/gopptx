// Package urlfetch converts web pages or raw HTML into PPTX presentations.
//
// It is a Go port of ppt-rs/src/web2ppt with improvements:
// HTML tables are rendered as real PPTX tables rather than summary bullets.
package urlfetch

const (
	defaultMaxSlides          = 20
	defaultMaxBulletsPerSlide = 6
	defaultTimeoutSecs        = 30
	bytesPerMiB               = 1024 * 1024
	defaultMaxBodyBytes       = 10 * bytesPerMiB
)

// Config holds options that control content extraction and slide generation.
type Config struct {
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

// Web2PptConfig is a compatibility alias for Config.
type Web2PptConfig = Config

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxSlides:          defaultMaxSlides,
		MaxBulletsPerSlide: defaultMaxBulletsPerSlide,
		IncludeImages:      true,
		IncludeTables:      true,
		IncludeCode:        true,
		ExtractLinks:       true,
		GroupByHeadings:    true,
		UserAgent:          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		TimeoutSecs:        defaultTimeoutSecs,
		MaxBodyBytes:       defaultMaxBodyBytes,
	}
}

// WithMaxSlides sets the maximum slide count.
func (c Config) WithMaxSlides(n int) Config {
	c.MaxSlides = n
	return c
}

// WithMaxBullets sets the maximum bullets-per-slide.
func (c Config) WithMaxBullets(n int) Config {
	c.MaxBulletsPerSlide = n
	return c
}

// WithImages enables or disables image extraction.
func (c Config) WithImages(v bool) Config {
	c.IncludeImages = v
	return c
}

// WithTables enables or disables table extraction.
func (c Config) WithTables(v bool) Config {
	c.IncludeTables = v
	return c
}

// WithCode enables or disables code-block extraction.
func (c Config) WithCode(v bool) Config {
	c.IncludeCode = v
	return c
}

// WithLinks enables or disables link extraction.
func (c Config) WithLinks(v bool) Config {
	c.ExtractLinks = v
	return c
}

// WithGroupByHeadings sets the slide-grouping strategy.
func (c Config) WithGroupByHeadings(v bool) Config {
	c.GroupByHeadings = v
	return c
}

// WithUserAgent overrides the HTTP User-Agent header.
func (c Config) WithUserAgent(ua string) Config {
	c.UserAgent = ua
	return c
}

// WithTimeout sets the HTTP request timeout in seconds.
func (c Config) WithTimeout(secs int) Config {
	c.TimeoutSecs = secs
	return c
}

// WithMaxBodyBytes sets the maximum response body size in bytes.
func (c Config) WithMaxBodyBytes(n int64) Config {
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
