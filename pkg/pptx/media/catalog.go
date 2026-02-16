package media

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

type mediaAsset struct {
	mediaName string
	extension string
	data      []byte
}

// MediaCatalog manages presentation assets (images, etc.).
type MediaCatalog struct {
	byKey   map[string]mediaAsset
	ordered []mediaAsset
}

var supportedMediaExtensions = map[string]struct{}{
	"png":  {},
	"jpg":  {},
	"jpeg": {},
	"gif":  {},
	"bmp":  {},
	"tif":  {},
	"tiff": {},
	"mp3":  {},
	"wav":  {},
	"m4a":  {},
}

// BuildMediaCatalog constructs a catalog from multiple slides.
func BuildMediaCatalog(slides []elements.SlideContent) (*MediaCatalog, error) {
	catalog := &MediaCatalog{
		byKey:   make(map[string]mediaAsset),
		ordered: make([]mediaAsset, 0),
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Helper to add image to catalog
	addImage := func(image shapes.Image) error {
		var data []byte
		var ext string
		var key string

		if image.Path != "" {
			cleanPath := filepath.Clean(image.Path)
			key = "path:" + cleanPath
			if _, exists := catalog.byKey[key]; exists {
				return nil
			}

			ext = NormalizeImageExtension(cleanPath)
			fileData, err := os.ReadFile(cleanPath)
			if err != nil {
				return fmt.Errorf("read error: %w", err)
			}
			data = fileData
		} else if len(image.Data) > 0 {
			ext = NormalizeExtension(image.Format)
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
			ext = NormalizeImageExtension(image.SourceURL)
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
			return errors.New("yielded empty data")
		}

		if _, ok := supportedMediaExtensions[ext]; !ok {
			if ext == "" {
				return errors.New("has unknown extension (cannot infer)")
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
		if slide.Background != nil && slide.Background.Type == elements.SlideBackgroundPicture &&
			slide.Background.PictureFill != nil {
			if err := addImage(*slide.Background.PictureFill); err != nil {
				return nil, fmt.Errorf("slide %d background image: %w", slideIndex+1, err)
			}
		}

		// Handle transition sounds
		if slide.Transition != nil {
			if opt, ok := slide.Transition.(transitions.TransitionOptions); ok && opt.Sound != nil &&
				strings.HasPrefix(opt.Sound.RelID, "file:") {
				path := strings.TrimPrefix(opt.Sound.RelID, "file:")
				soundMedia := shapes.Image{Path: path}
				if err := addImage(soundMedia); err != nil {
					return nil, fmt.Errorf("slide %d transition sound: %w", slideIndex+1, err)
				}

				// Find the registered media name to determine RID later?
				// Actually, BuildMediaCatalog builds the catalog but doesn't assign RIDs here usually?
				// Wait, presentation.go assigns RIDs based on order in catalog or just sequence?
				// presentation.go assigns RIDs sequence: rId2, rId3...
				// And it appends to imageTargets based on mediaName.
				// Here we just ensure it's in catalog.
				// BUT we need to update the RelID in the slide to be something we can use in presentation.go?
				// presentation.go Logic:
				//   mediaName, ok := mediaCatalog.MediaNameForImage(image)
				//   relID := ...
				//   imageTargets = append(...)

				// So here in BuildMediaCatalog, we just add it.
				// But we also need to allow presentation.go to find it.
				// The slide.Transition logic in presentation.go needs to change to LOOK UP the sound.

				// Keep the "file:" prefix or change it?
				// If we leave "file:", presentation.go can parse it and look up in catalog.
				// If we change it here, presentation.go needs to know what we changed it to.

				// Better: leave "file:" prefix, so presentation.go knows it needs resolving.
			}
		}
	}

	return catalog, nil
}

// NormalizeImageExtension sanitizes file extensions from paths or URLs.
func NormalizeImageExtension(path string) string {
	// Remove query params if URL
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	return ext
}

// NormalizeExtension sanitizes a raw extension string.
func NormalizeExtension(ext string) string {
	ext = strings.ToLower(ext)
	if strings.HasPrefix(ext, ".") {
		return ext[1:]
	}
	return ext
}

// MediaNameForImage returns the registered name of an image in the catalog.
func (c *MediaCatalog) MediaNameForImage(image shapes.Image) (string, bool) {
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

// ImageExtensions returns all unique image extensions present in the catalog.
func (c *MediaCatalog) ImageExtensions() []string {
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

// Assets returns all ordered assets in the catalog.
func (c *MediaCatalog) Assets() []mediaAsset {
	return c.ordered
}

// MediaName returns the name of a media asset.
func (a mediaAsset) MediaName() string { return a.mediaName }

// Extension returns the extension of a media asset.
func (a mediaAsset) Extension() string { return a.extension }

// Data returns the raw content of a media asset.
func (a mediaAsset) Data() []byte { return a.data }

// RegisterImage registers an image in the catalog and returns its media name.
// This is useful for registering master slide images outside of BuildMediaCatalog.
func (c *MediaCatalog) RegisterImage(image shapes.Image) (string, error) {
	var data []byte
	var ext string
	var key string

	if image.Path != "" {
		cleanPath := filepath.Clean(image.Path)
		key = "path:" + cleanPath
		if asset, exists := c.byKey[key]; exists {
			return asset.mediaName, nil
		}
		ext = NormalizeImageExtension(cleanPath)
		fileData, err := os.ReadFile(cleanPath)
		if err != nil {
			return "", fmt.Errorf("read error: %w", err)
		}
		data = fileData
	} else if len(image.Data) > 0 {
		ext = NormalizeExtension(image.Format)
		hash := sha256.Sum256(image.Data)
		key = "data:" + hex.EncodeToString(hash[:])
		if asset, exists := c.byKey[key]; exists {
			return asset.mediaName, nil
		}
		data = image.Data
	} else {
		return "", errors.New("image has no path or data")
	}

	if len(data) == 0 {
		return "", errors.New("yielded empty data")
	}
	if _, ok := supportedMediaExtensions[ext]; !ok {
		return "", fmt.Errorf("unsupported extension %q", ext)
	}

	mediaName := fmt.Sprintf("image%d.%s", len(c.ordered)+1, ext)
	asset := mediaAsset{mediaName: mediaName, extension: ext, data: data}
	c.byKey[key] = asset
	c.ordered = append(c.ordered, asset)
	return mediaName, nil
}
