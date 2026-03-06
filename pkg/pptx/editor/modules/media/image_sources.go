package media

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	imageURLFetchTimeout = 30 * time.Second
	dataURIParts         = 2
)

const (
	mimePNG    = "image/png"
	mimeJPEG   = "image/jpeg"
	mimeGIF    = "image/gif"
	mimeBMP    = "image/bmp"
	mimeTIFF   = "image/tiff"
	formatPNG  = "png"
	formatJPEG = "jpeg"
	formatGIF  = "gif"
	formatBMP  = "bmp"
	formatTIFF = "tiff"
)

// DecodeBase64ImagePayload decodes raw base64 or data-URI image payloads.
func DecodeBase64ImagePayload(raw string) ([]byte, string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, "", errors.New("image base64 data is required")
	}

	payload := trimmed
	detectedFormat := ""
	if strings.HasPrefix(strings.ToLower(trimmed), "data:") {
		mimeType, dataPart, err := splitDataURI(trimmed)
		if err != nil {
			return nil, "", err
		}
		detectedFormat = formatFromMimeType(mimeType)
		payload = dataPart
	}

	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, "", fmt.Errorf("invalid image base64 data: %w", err)
	}
	return data, detectedFormat, nil
}

func splitDataURI(raw string) (string, string, error) {
	parts := strings.SplitN(raw, ",", dataURIParts)
	if len(parts) != dataURIParts {
		return "", "", errors.New("invalid data URI image payload")
	}
	meta := parts[0]
	if !strings.Contains(strings.ToLower(meta), ";base64") {
		return "", "", errors.New("data URI image payload must be base64 encoded")
	}
	mimeType := strings.TrimPrefix(meta, "data:")
	if idx := strings.Index(mimeType, ";"); idx >= 0 {
		mimeType = mimeType[:idx]
	}
	return strings.TrimSpace(mimeType), parts[1], nil
}

// FetchImageFromURL downloads image bytes and infers their format.
func FetchImageFromURL(sourceURL string) ([]byte, string, error) {
	parsed, err := url.ParseRequestURI(strings.TrimSpace(sourceURL))
	if err != nil {
		return nil, "", fmt.Errorf("invalid image URL: %w", err)
	}

	client := &http.Client{Timeout: imageURLFetchTimeout}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("fetch image URL: build request: %w", err)
	}
	resp, err := client.Do(req) // #nosec G107: explicit URL validation above.
	if err != nil {
		return nil, "", fmt.Errorf("fetch image URL: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("fetch image URL: unexpected status %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read image URL body: %w", err)
	}
	if len(data) == 0 {
		return nil, "", errors.New("fetch image URL: empty response body")
	}

	format := formatFromMimeType(resp.Header.Get("Content-Type"))
	if format == "" {
		format = formatFromPath(parsed.Path)
	}
	if format == "" {
		return nil, "", errors.New("fetch image URL: unable to infer image format")
	}
	return data, format, nil
}

func formatFromMimeType(contentType string) string {
	mimeType := strings.TrimSpace(strings.ToLower(contentType))
	if idx := strings.Index(mimeType, ";"); idx >= 0 {
		mimeType = strings.TrimSpace(mimeType[:idx])
	}
	switch mimeType {
	case mimePNG:
		return formatPNG
	case mimeJPEG:
		return formatJPEG
	case mimeGIF:
		return formatGIF
	case mimeBMP:
		return formatBMP
	case mimeTIFF:
		return formatTIFF
	default:
		return ""
	}
}

func formatFromPath(pathValue string) string {
	ext := strings.TrimPrefix(strings.ToLower(path.Ext(pathValue)), ".")
	switch ext {
	case formatPNG, "jpg", formatJPEG, formatGIF, formatBMP:
		return ext
	case "tif", formatTIFF:
		return formatTIFF
	default:
		return ""
	}
}
