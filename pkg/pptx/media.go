package pptx

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type mediaAsset struct {
	mediaName string
	extension string
	data      []byte
}

type mediaCatalog struct {
	byKey   map[string]mediaAsset
	ordered []mediaAsset
}

var supportedImageExtensions = map[string]struct{}{
	"png":  {},
	"jpg":  {},
	"jpeg": {},
	"gif":  {},
	"bmp":  {},
	"tif":  {},
	"tiff": {},
}

func buildMediaCatalog(slides []SlideContent) (*mediaCatalog, error) {
	catalog := &mediaCatalog{
		byKey:   make(map[string]mediaAsset),
		ordered: make([]mediaAsset, 0),
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Helper to add image to catalog
	addImage := func(image Image) error {
		var data []byte
		var ext string
		var key string

		if image.Path != "" {
			cleanPath := filepath.Clean(image.Path)
			key = "path:" + cleanPath
			if _, exists := catalog.byKey[key]; exists {
				return nil
			}

			ext = normalizeImageExtension(cleanPath)
			fileData, err := os.ReadFile(cleanPath)
			if err != nil {
				return fmt.Errorf("read error: %w", err)
			}
			data = fileData

		} else if len(image.Data) > 0 {
			ext = normalizeExtension(image.Format)
			// key by hash of data
			hash := sha256.Sum256(image.Data)
			key = "data:" + hex.EncodeToString(hash[:])
			if _, exists := catalog.byKey[key]; exists {
				return nil
			}
			data = image.Data

		} else if image.SourceURL != "" {
			key = "url:" + image.SourceURL
			if _, exists := catalog.byKey[key]; exists {
				return nil
			}

			// Basic extension inference from URL
			ext = normalizeImageExtension(image.SourceURL)
			if ext == "" {
				ext = "png"
			}

			resp, err := client.Get(image.SourceURL)
			if err != nil {
				return fmt.Errorf("fetch error: %w", err)
			}

			if resp.StatusCode != http.StatusOK {
				if closeErr := resp.Body.Close(); closeErr != nil {
					return fmt.Errorf("fetch failed with status: %s (close body error: %w)", resp.Status, closeErr)
				}
				return fmt.Errorf("fetch failed with status: %s", resp.Status)
			}

			downloaded, readErr := io.ReadAll(resp.Body)
			closeErr := resp.Body.Close()
			if readErr != nil {
				return fmt.Errorf("read body error: %w", readErr)
			}
			if closeErr != nil {
				return fmt.Errorf("close body error: %w", closeErr)
			}
			data = downloaded

			if ext == "" {
				ct := resp.Header.Get("Content-Type")
				switch ct {
				case "image/jpeg":
					ext = "jpg"
				case "image/png":
					ext = "png"
				case "image/gif":
					ext = "gif"
				case "image/bmp":
					ext = "bmp"
				case "image/tiff":
					ext = "tiff"
				}
			}
		}

		if len(data) == 0 {
			return fmt.Errorf("yielded empty data")
		}

		if _, ok := supportedImageExtensions[ext]; !ok {
			if ext == "" {
				return fmt.Errorf("has unknown extension (cannot infer)")
			}
			return fmt.Errorf("has unsupported extension %q", ext)
		}

		mediaName := fmt.Sprintf("image%d.%s", len(catalog.ordered)+1, ext)
		asset := mediaAsset{
			mediaName: mediaName,
			extension: ext,
			data:      data,
		}
		catalog.byKey[key] = asset
		catalog.ordered = append(catalog.ordered, asset)
		return nil
	}

	for slideIndex, slide := range slides {
		for imageIndex, image := range slide.Images {
			if err := addImage(image); err != nil {
				return nil, fmt.Errorf("slide %d image %d: %w", slideIndex+1, imageIndex+1, err)
			}
		}
		for phIndex, override := range slide.PlaceholderOverrides {
			if override.Image != nil {
				if err := addImage(*override.Image); err != nil {
					return nil, fmt.Errorf("slide %d placeholder override %d image: %w", slideIndex+1, phIndex+1, err)
				}
			}
		}
	}

	return catalog, nil
}

func normalizeImageExtension(path string) string {
	// Remove query params if URL
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	return ext
}

func normalizeExtension(ext string) string {
	ext = strings.ToLower(ext)
	if strings.HasPrefix(ext, ".") {
		return ext[1:]
	}
	return ext
}

func (c *mediaCatalog) mediaNameForImage(image Image) (string, bool) {
	var key string
	if image.SourceURL != "" {
		key = "url:" + image.SourceURL
	} else if len(image.Data) > 0 {
		hash := sha256.Sum256(image.Data)
		key = "data:" + hex.EncodeToString(hash[:])
	} else {
		key = "path:" + filepath.Clean(image.Path)
	}

	asset, ok := c.byKey[key]
	if !ok {
		return "", false
	}
	return asset.mediaName, true
}

func (c *mediaCatalog) imageExtensions() []string {
	seen := make(map[string]struct{})
	out := make([]string, 0, len(c.ordered))
	for _, asset := range c.ordered {
		if _, ok := seen[asset.extension]; ok {
			continue
		}
		seen[asset.extension] = struct{}{}
		out = append(out, asset.extension)
	}
	return out
}
