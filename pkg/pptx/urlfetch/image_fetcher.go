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

	"github.com/djinn-soul/gopptx/pkg/pptx/netsec"
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
// replaced with netsec.NewRestrictedTransport so IP-range checks happen at connection time.
func NewImageFetcher(client *http.Client, cfg Config, baseURL string) *ImageFetcher {
	var base *url.URL
	if baseURL != "" {
		base, _ = url.Parse(baseURL)
	}
	safeClient := &http.Client{
		Timeout:       client.Timeout,
		CheckRedirect: client.CheckRedirect,
		Jar:           client.Jar,
		Transport:     netsec.NewRestrictedTransport(cfg.AllowPrivateHosts),
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
	decodedMIME, knownFormat := imageFormatMIME(format)
	if !knownFormat || !f.isAllowedImageType(decodedMIME) {
		return nil, fmt.Errorf("decoded image type is not allowed: image/%s", format)
	}
	if !strings.EqualFold(normalizeContentType(contentType), decodedMIME) {
		return nil, fmt.Errorf(
			"image content type mismatch: header %q, decoded %q",
			contentType,
			decodedMIME,
		)
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
	contentType = normalizeContentType(contentType)

	allowedTypes := f.cfg.AllowedImageTypes
	if len(allowedTypes) == 0 {
		allowedTypes = []string{imageMIMEPNG, imageMIMEJPEG, imageMIMEGIF}
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
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", 0, 0, fmt.Errorf("decode image config: %w", err)
	}

	return format, cfg.Width, cfg.Height, nil
}

func imageFormatMIME(format string) (string, bool) {
	mimeType, ok := map[string]string{
		"gif":  imageMIMEGIF,
		"jpeg": imageMIMEJPEG,
		"png":  imageMIMEPNG,
	}[strings.ToLower(format)]
	return mimeType, ok
}

func normalizeContentType(contentType string) string {
	if semicolonIndex := strings.Index(contentType, ";"); semicolonIndex != -1 {
		contentType = contentType[:semicolonIndex]
	}
	return strings.TrimSpace(contentType)
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
