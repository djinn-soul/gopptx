package urlfetch

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"  // Register GIF decoder
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// ImageFetcher handles concurrent image downloads with size limits.
type ImageFetcher struct {
	client    *http.Client
	cfg       Config
	baseURL   *url.URL
	mu        sync.Mutex
	totalSize int64
}

// NewImageFetcher creates an ImageFetcher with the given config and base URL.
// The passed client's Timeout and CheckRedirect are preserved; its transport is
// replaced with ssrfSafeTransport so IP-range checks happen at connection time.
func NewImageFetcher(client *http.Client, cfg Config, baseURL string) *ImageFetcher {
	var base *url.URL
	if baseURL != "" {
		base, _ = url.Parse(baseURL)
	}
	safeClient := &http.Client{
		Timeout:       client.Timeout,
		CheckRedirect: client.CheckRedirect,
		Jar:           client.Jar,
		Transport:     ssrfSafeTransport(cfg.AllowPrivateHosts),
	}
	return &ImageFetcher{
		client:    safeClient,
		cfg:       cfg,
		baseURL:   base,
		totalSize: 0,
	}
}

// FetchedImage represents a downloaded image with its metadata.
type FetchedImage struct {
	Data   []byte
	Format string
	Width  int
	Height int
}

const (
	minMagicBytes      = 4
	webpMagicBytes     = 12
	defaultAspectNum   = 3
	defaultAspectDenom = 4
)

// FetchImage downloads and validates an image from the given URL.
// Returns error if the image exceeds size limits or has invalid MIME type.
func (f *ImageFetcher) FetchImage(imageURL string) (*FetchedImage, error) {
	// Resolve relative URLs
	resolvedURL, err := f.resolveURL(imageURL)
	if err != nil {
		return nil, fmt.Errorf("resolve image URL: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, resolvedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", f.cfg.UserAgent)
	req.Header.Set("Accept", "image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	if f.baseURL != nil {
		req.Header.Set("Referer", redactURL(f.baseURL))
	}

	// Execute request
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch image: HTTP %d", resp.StatusCode)
	}

	// Check Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !f.isAllowedImageType(contentType) {
		return nil, fmt.Errorf("unsupported image type: %s", contentType)
	}

	// Read response with size limit
	maxSize := f.cfg.MaxImageSizeBytes
	if maxSize <= 0 {
		maxSize = defaultMaxImageSizeBytes
	}

	limitedReader := io.LimitReader(resp.Body, maxSize+1)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("read image data: %w", err)
	}

	if int64(len(data)) > maxSize {
		return nil, fmt.Errorf("image exceeds maximum size of %d bytes", maxSize)
	}

	// Detect image dimensions and format
	format, width, height, err := f.detectImageInfo(data)
	if err != nil {
		return nil, fmt.Errorf("detect image info: %w", err)
	}

	// Check and reserve total size limit only after validation succeeds.
	f.mu.Lock()
	defer f.mu.Unlock()
	maxTotal := f.cfg.MaxTotalImageSizeBytes
	if maxTotal <= 0 {
		maxTotal = defaultMaxTotalImageSize
	}
	if f.totalSize+int64(len(data)) > maxTotal {
		return nil, fmt.Errorf("total image size would exceed maximum of %d bytes", maxTotal)
	}
	f.totalSize += int64(len(data))

	return &FetchedImage{
		Data:   data,
		Format: format,
		Width:  width,
		Height: height,
	}, nil
}

// resolveURL resolves a potentially relative URL against the base URL.
func (f *ImageFetcher) resolveURL(imageURL string) (string, error) {
	// If it's already absolute, return as-is
	if strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://") {
		return imageURL, nil
	}

	// If it's a data URI, reject it (should be handled separately)
	if strings.HasPrefix(imageURL, "data:") {
		return "", errors.New("data URIs not supported")
	}

	// If it's protocol-relative, add https
	if strings.HasPrefix(imageURL, "//") {
		return "https:" + imageURL, nil
	}

	// Resolve against base URL
	if f.baseURL == nil {
		return imageURL, nil // Can't resolve, return as-is
	}

	u, err := f.baseURL.Parse(imageURL)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

// isAllowedImageType checks if the content type is in the allowed list.
func (f *ImageFetcher) isAllowedImageType(contentType string) bool {
	// Extract main type/subtype (ignore parameters like charset)
	semicolonIdx := strings.Index(contentType, ";")
	if semicolonIdx != -1 {
		contentType = strings.TrimSpace(contentType[:semicolonIdx])
	}

	allowedTypes := f.cfg.AllowedImageTypes
	if len(allowedTypes) == 0 {
		allowedTypes = []string{"image/png", "image/jpeg", "image/gif"}
	}

	for _, allowed := range allowedTypes {
		if strings.EqualFold(contentType, allowed) {
			return true
		}
	}

	return false
}

// detectImageInfo detects the format and dimensions of image data.
func (f *ImageFetcher) detectImageInfo(data []byte) (string, int, int, error) {
	// Use image.DecodeConfig to get dimensions without full decode
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		// Try to detect format from magic bytes if DecodeConfig fails
		format = detectFormatFromMagic(data)
		if format == "" {
			return "", 0, 0, fmt.Errorf("decode image config: %w", err)
		}
		// Return with zero dimensions if we can't decode
		return format, 0, 0, nil
	}

	return format, cfg.Width, cfg.Height, nil
}

// detectFormatFromMagic detects image format from magic bytes.
func detectFormatFromMagic(data []byte) string {
	if len(data) < minMagicBytes {
		return ""
	}

	// PNG: 89 50 4E 47
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "png"
	}

	// JPEG: FF D8 FF
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "jpeg"
	}

	// GIF: GIF87a or GIF89a
	if data[0] == 'G' && data[1] == 'I' && data[2] == 'F' && data[3] == '8' {
		return "gif"
	}

	// WebP: RIFF....WEBP
	if len(data) >= webpMagicBytes && data[0] == 'R' && data[1] == 'I' && data[2] == 'F' && data[3] == 'F' {
		if data[8] == 'W' && data[9] == 'E' && data[10] == 'B' && data[11] == 'P' {
			return "webp"
		}
	}

	return ""
}

// CalculateImageDimensions calculates EMU dimensions preserving aspect ratio.
// targetWidthEMU is the desired width; height is calculated to maintain aspect ratio.
func CalculateImageDimensions(width, height int, targetWidthEMU int64) (int64, int64) {
	if width <= 0 || height <= 0 {
		// Default size if dimensions unknown
		return targetWidthEMU, targetWidthEMU * defaultAspectNum / defaultAspectDenom // 4:3 aspect ratio default
	}

	// Calculate height maintaining aspect ratio
	targetHeightEMU := targetWidthEMU * int64(height) / int64(width)

	return targetWidthEMU, targetHeightEMU
}
