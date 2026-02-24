package urlfetch

// HTMLToPPTX converts a raw HTML string to a PPTX byte slice using default settings.
func HTMLToPPTX(html, url string) ([]byte, error) {
	opts := DefaultConversionOptions()
	return HTMLToPPTXWithOptions(html, url, DefaultConfig(), opts)
}

// HTMLToPPTXWithOptions converts a raw HTML string to PPTX with custom config and options.
func HTMLToPPTXWithOptions(html, url string, cfg Config, opts ConversionOptions) ([]byte, error) {
	parser := NewWebParserWithConfig(cfg)
	content, err := parser.Parse(html, url)
	if err != nil {
		return nil, err
	}
	converter := NewURLFetchConverterWithConfig(cfg)
	return converter.Convert(content, &opts)
}

// URLToPPTX fetches the page at url and converts it to a PPTX byte slice.
func URLToPPTX(url string) ([]byte, error) {
	opts := DefaultConversionOptions()
	return URLToPPTXWithOptions(url, DefaultConfig(), opts)
}

// URLToPPTXWithOptions fetches the page at url and converts it with custom config and options.
func URLToPPTXWithOptions(url string, cfg Config, opts ConversionOptions) ([]byte, error) {
	fetcher := NewWebFetcherWithConfig(cfg)
	finalURL, html, err := fetcher.FetchWithURL(url)
	if err != nil {
		return nil, err
	}
	return HTMLToPPTXWithOptions(html, finalURL, cfg, opts)
}
