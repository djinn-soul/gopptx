package urlfetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// WebFetcher fetches HTML from URLs using net/http.
type WebFetcher struct {
	client http.Client
	cfg    URLFetchConfig
}

// NewWebFetcher creates a WebFetcher with default config.
func NewWebFetcher() *WebFetcher {
	return NewWebFetcherWithConfig(DefaultConfig())
}

// NewWebFetcherWithConfig creates a WebFetcher with custom config.
func NewWebFetcherWithConfig(cfg URLFetchConfig) *WebFetcher {
	return &WebFetcher{
		client: http.Client{
			Timeout: time.Duration(cfg.TimeoutSecs) * time.Second,
			CheckRedirect: func(_ *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("stopped after %d redirects", len(via))
				}
				return nil
			},
		},
		cfg: cfg,
	}
}

// Fetch retrieves the HTML body from the given URL.
// Only http and https schemes are accepted.
func (f *WebFetcher) Fetch(rawURL string) (string, error) {
	_, html, err := f.fetchWithFinalURL(rawURL)
	if err != nil {
		return "", err
	}
	return html, nil
}

func (f *WebFetcher) fetchWithFinalURL(rawURL string) (string, string, error) {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", "", fmt.Errorf("unsupported scheme %q: only http and https are allowed", parsed.Scheme)
	}
	if f.cfg.MaxBodyBytes <= 0 {
		return "", "", fmt.Errorf("invalid MaxBodyBytes: %d", f.cfg.MaxBodyBytes)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, rawURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("User-Agent", f.cfg.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Referer", rawURL)

	resp, err := f.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("fetch %q: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("fetch %q: HTTP %d", rawURL, resp.StatusCode)
	}

	limited := io.LimitReader(resp.Body, f.cfg.MaxBodyBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return "", "", fmt.Errorf("read body: %w", err)
	}
	if int64(len(body)) > f.cfg.MaxBodyBytes {
		return "", "", fmt.Errorf("response too large: read > %d bytes", f.cfg.MaxBodyBytes)
	}

	finalURL := rawURL
	if resp.Request != nil && resp.Request.URL != nil {
		finalURL = resp.Request.URL.String()
	}
	return finalURL, string(body), nil
}

// FetchWithURL returns both the canonical URL (after redirects) and HTML body.
func (f *WebFetcher) FetchWithURL(rawURL string) (string, string, error) {
	finalURL, html, err := f.fetchWithFinalURL(rawURL)
	if err != nil {
		return "", "", err
	}
	return finalURL, html, nil
}
