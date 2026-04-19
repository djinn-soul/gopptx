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
	defaultMaxImageSizeBytes  = 5 * bytesPerMiB
	defaultMaxTotalImageSize  = 20 * bytesPerMiB
	defaultMaxImagesPerSlide  = 3
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
	// DownloadImages fetches and embeds remote images when true (replaces alt-text bullets).
	DownloadImages bool
	// MaxImageSizeBytes caps individual downloaded image size (default: 5MB).
	MaxImageSizeBytes int64
	// MaxTotalImageSizeBytes caps total downloaded image size per page (default: 20MB).
	MaxTotalImageSizeBytes int64
	// MaxImagesPerSlide limits images per slide (default: 3).
	MaxImagesPerSlide int
	// AllowedImageTypes filters by MIME type (default: ["image/png", "image/jpeg", "image/gif"]).
	AllowedImageTypes []string
	// ContentSelectors overrides default main-content CSS selectors.
	ContentSelectors []string
	// ExcludeSelectors removes matching elements from content.
	ExcludeSelectors []string
	// AllowPrivateHosts disables the SSRF guard that blocks requests to
	// loopback/private/link-local addresses. Must only be set to true in tests.
	AllowPrivateHosts bool
}

// Web2PptConfig is a compatibility alias for Config.
type Web2PptConfig = Config

// URLFetchConfig is a compatibility alias for Config.
//
//nolint:revive // Required compatibility alias for existing public API consumers.
type URLFetchConfig = Config

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxSlides:              defaultMaxSlides,
		MaxBulletsPerSlide:     defaultMaxBulletsPerSlide,
		IncludeImages:          true,
		IncludeTables:          true,
		IncludeCode:            true,
		ExtractLinks:           true,
		GroupByHeadings:        true,
		UserAgent:              "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		TimeoutSecs:            defaultTimeoutSecs,
		MaxBodyBytes:           defaultMaxBodyBytes,
		MaxImageSizeBytes:      defaultMaxImageSizeBytes,
		MaxTotalImageSizeBytes: defaultMaxTotalImageSize,
		MaxImagesPerSlide:      defaultMaxImagesPerSlide,
		AllowedImageTypes:      []string{"image/png", "image/jpeg", "image/gif"},
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

// WithDownloadImages enables or disables image downloading and embedding.
func (c Config) WithDownloadImages(v bool) Config {
	c.DownloadImages = v
	return c
}

// WithMaxImageSizeBytes sets the maximum size for individual downloaded images.
func (c Config) WithMaxImageSizeBytes(n int64) Config {
	c.MaxImageSizeBytes = n
	return c
}

// WithMaxTotalImageSizeBytes sets the maximum total size for all downloaded images.
func (c Config) WithMaxTotalImageSizeBytes(n int64) Config {
	c.MaxTotalImageSizeBytes = n
	return c
}

// WithMaxImagesPerSlide sets the maximum number of images per slide.
func (c Config) WithMaxImagesPerSlide(n int) Config {
	c.MaxImagesPerSlide = n
	return c
}

// WithAllowedImageTypes sets the allowed MIME types for downloaded images.
func (c Config) WithAllowedImageTypes(types []string) Config {
	c.AllowedImageTypes = types
	return c
}

// WithContentSelectors sets custom CSS selectors for content extraction.
func (c Config) WithContentSelectors(selectors []string) Config {
	c.ContentSelectors = selectors
	return c
}

// WithExcludeSelectors sets CSS selectors for elements to exclude from extraction.
func (c Config) WithExcludeSelectors(selectors []string) Config {
	c.ExcludeSelectors = selectors
	return c
}

// WithAllowPrivateHosts disables the SSRF guard for testing against local servers.
// Must not be used in production code.
func (c Config) WithAllowPrivateHosts(v bool) Config {
	c.AllowPrivateHosts = v
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
