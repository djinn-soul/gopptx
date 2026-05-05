package media

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

type Asset struct {
	mediaName string
	extension string
	data      []byte
}

// Catalog manages presentation assets (images, etc.).
type Catalog struct {
	byKey   map[string]Asset
	ordered []Asset
}

const (
	defaultImageExt   = "png"
	jpegImageExt      = "jpg"
	jpegAliasImageExt = "jpeg"
	gifImageExt       = "gif"
	bmpImageExt       = "bmp"
	tifImageExt       = "tif"
	tiffImageExt      = "tiff"
	mp3AudioExt       = "mp3"
	wavAudioExt       = "wav"
	m4aAudioExt       = "m4a"
)

const imageFetchTimeout = 30 * time.Second

const maxCatalogImageURLBodyBytes = 20 * 1024 * 1024

// MediaName returns the name of a media asset.
func (a Asset) MediaName() string { return a.mediaName }

// Extension returns the extension of a media asset.
func (a Asset) Extension() string { return a.extension }

// Data returns the raw content of a media asset.
func (a Asset) Data() []byte { return a.data }

// Assets returns all ordered assets in the catalog.
func (c *Catalog) Assets() []Asset {
	return c.ordered
}

// ImageExtensions returns all unique image extensions present in the catalog.
func (c *Catalog) ImageExtensions() []string {
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

// MediaNameForImage returns the registered name of an image in the catalog.
func (c *Catalog) MediaNameForImage(image shapes.Image) (string, bool) {
	var key string
	switch {
	case image.SourceURL != "":
		key = "url:" + image.SourceURL
	case len(image.Data) > 0:
		hash := sha256.Sum256(image.Data)
		key = "data:" + hex.EncodeToString(hash[:])
	default:
		key = "path:" + filepath.Clean(image.Path)
	}

	asset, ok := c.byKey[key]
	if !ok {
		return "", false
	}
	return asset.mediaName, true
}

// RegisterImage registers an image in the catalog and returns its media name.
// This is useful for registering master slide images outside of BuildMediaCatalog.
func (c *Catalog) RegisterImage(image shapes.Image) (string, error) {
	var data []byte
	var ext string
	var key string

	switch {
	case image.Path != "":
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
	case len(image.Data) > 0:
		ext = NormalizeExtension(image.Format)
		hash := sha256.Sum256(image.Data)
		key = "data:" + hex.EncodeToString(hash[:])
		if asset, exists := c.byKey[key]; exists {
			return asset.mediaName, nil
		}
		data = image.Data
	default:
		return "", errors.New("image has no path or data")
	}

	if len(data) == 0 {
		return "", errors.New("yielded empty data")
	}
	if !isSupportedMediaExtension(ext) {
		return "", fmt.Errorf("unsupported extension %q", ext)
	}

	mediaName := fmt.Sprintf("image%d.%s", len(c.ordered)+1, ext)
	asset := Asset{mediaName: mediaName, extension: ext, data: data}
	c.byKey[key] = asset
	c.ordered = append(c.ordered, asset)
	return mediaName, nil
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

func isSupportedMediaExtension(ext string) bool {
	switch ext {
	case defaultImageExt,
		jpegImageExt,
		jpegAliasImageExt,
		gifImageExt,
		bmpImageExt,
		tifImageExt,
		tiffImageExt,
		mp3AudioExt,
		wavAudioExt,
		m4aAudioExt:
		return true
	default:
		return false
	}
}
