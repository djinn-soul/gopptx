package media

import (
	"context"
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

const defaultImageExt = "png"

const imageFetchTimeout = 30 * time.Second

// BuildMediaCatalog constructs a catalog from multiple slides.
func BuildMediaCatalog(slides []elements.SlideContent, notesMaster *elements.NotesMaster) (*Catalog, error) {
	catalog := &Catalog{
		byKey:   make(map[string]Asset),
		ordered: make([]Asset, 0),
	}
	client := &http.Client{
		Timeout: imageFetchTimeout,
	}

	for slideIndex, slide := range slides {
		if err := addSlideMedia(catalog, client, slide, slideIndex); err != nil {
			return nil, err
		}
	}

	if notesMaster != nil {
		if err := addNotesMasterMedia(catalog, client, notesMaster); err != nil {
			return nil, err
		}
	}

	return catalog, nil
}

func addNotesMasterMedia(catalog *Catalog, client *http.Client, master *elements.NotesMaster) error {
	if master.Background == nil ||
		master.Background.Type != elements.SlideBackgroundPicture ||
		master.Background.PictureFill == nil {
		return nil
	}
	if err := addImageToCatalog(catalog, client, *master.Background.PictureFill); err != nil {
		return fmt.Errorf("notes master background image: %w", err)
	}
	return nil
}

func addSlideMedia(catalog *Catalog, client *http.Client, slide elements.SlideContent, slideIndex int) error {
	for imageIndex, image := range slide.Images {
		if err := addImageToCatalog(catalog, client, image); err != nil {
			return fmt.Errorf("slide %d image %d: %w", slideIndex+1, imageIndex+1, err)
		}
	}
	for phIndex, override := range slide.PlaceholderOverrides {
		if override.Image == nil {
			continue
		}
		if err := addImageToCatalog(catalog, client, *override.Image); err != nil {
			return fmt.Errorf("slide %d placeholder override %d image: %w", slideIndex+1, phIndex+1, err)
		}
	}
	if err := addBackgroundImageToCatalog(catalog, client, slide, slideIndex); err != nil {
		return err
	}
	return addTransitionSoundToCatalog(catalog, client, slide, slideIndex)
}

func addBackgroundImageToCatalog(
	catalog *Catalog,
	client *http.Client,
	slide elements.SlideContent,
	slideIndex int,
) error {
	if slide.Background == nil ||
		slide.Background.Type != elements.SlideBackgroundPicture ||
		slide.Background.PictureFill == nil {
		return nil
	}
	if err := addImageToCatalog(catalog, client, *slide.Background.PictureFill); err != nil {
		return fmt.Errorf("slide %d background image: %w", slideIndex+1, err)
	}
	return nil
}

func addTransitionSoundToCatalog(
	catalog *Catalog,
	client *http.Client,
	slide elements.SlideContent,
	slideIndex int,
) error {
	if slide.Transition == nil {
		return nil
	}
	opt, ok := slide.Transition.(transitions.TransitionOptions)
	if !ok || opt.Sound == nil || !strings.HasPrefix(opt.Sound.RelID, "file:") {
		return nil
	}
	path := strings.TrimPrefix(opt.Sound.RelID, "file:")
	if err := addImageToCatalog(catalog, client, shapes.Image{Path: path}); err != nil {
		return fmt.Errorf("slide %d transition sound: %w", slideIndex+1, err)
	}
	return nil
}

func addImageToCatalog(catalog *Catalog, client *http.Client, image shapes.Image) error {
	key, ext, data, err := resolveImageAsset(client, image)
	if err != nil {
		return err
	}
	if _, exists := catalog.byKey[key]; exists {
		return nil
	}
	return appendMediaAsset(catalog, key, ext, data)
}

func resolveImageAsset(client *http.Client, image shapes.Image) (string, string, []byte, error) {
	if image.Path != "" {
		return loadImageFromPath(image.Path)
	}
	if len(image.Data) > 0 {
		key, ext, data := loadImageFromBytes(image.Data, image.Format)
		return key, ext, data, nil
	}
	if image.SourceURL != "" {
		return loadImageFromURL(client, image.SourceURL)
	}
	return "", "", nil, errors.New("image has no path, data, or source URL")
}

func loadImageFromPath(path string) (string, string, []byte, error) {
	cleanPath := filepath.Clean(path)
	fileData, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", "", nil, fmt.Errorf("read error: %w", err)
	}
	return "path:" + cleanPath, NormalizeImageExtension(cleanPath), fileData, nil
}

func loadImageFromBytes(data []byte, format string) (string, string, []byte) {
	hash := sha256.Sum256(data)
	return "data:" + hex.EncodeToString(hash[:]), NormalizeExtension(format), data
}

func loadImageFromURL(client *http.Client, sourceURL string) (string, string, []byte, error) {
	req, reqErr := http.NewRequestWithContext(context.Background(), http.MethodGet, sourceURL, http.NoBody)
	if reqErr != nil {
		return "", "", nil, fmt.Errorf("build request error: %w", reqErr)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", nil, fmt.Errorf("fetch error: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", "", nil, fmt.Errorf("fetch failed with status: %s", resp.Status)
	}
	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", "", nil, fmt.Errorf("read body error: %w", readErr)
	}
	ext := resolveURLExtension(sourceURL, resp.Header.Get("Content-Type"))
	return "url:" + sourceURL, ext, data, nil
}

func resolveURLExtension(sourceURL string, contentType string) string {
	ext := NormalizeImageExtension(sourceURL)
	if ext != "" {
		return ext
	}
	switch contentType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return defaultImageExt
	case "image/gif":
		return "gif"
	case "image/bmp":
		return "bmp"
	case "image/tiff":
		return "tiff"
	default:
		return defaultImageExt
	}
}

func appendMediaAsset(catalog *Catalog, key string, ext string, data []byte) error {
	if len(data) == 0 {
		return errors.New("yielded empty data")
	}
	if !isSupportedMediaExtension(ext) {
		if ext == "" {
			return errors.New("has unknown extension (cannot infer)")
		}
		return fmt.Errorf("has unsupported extension %q", ext)
	}
	mediaName := fmt.Sprintf("image%d.%s", len(catalog.ordered)+1, ext)
	asset := Asset{
		mediaName: mediaName,
		extension: ext,
		data:      data,
	}
	catalog.byKey[key] = asset
	catalog.ordered = append(catalog.ordered, asset)
	return nil
}

func isSupportedMediaExtension(ext string) bool {
	switch ext {
	case "png", "jpg", "jpeg", "gif", "bmp", "tif", "tiff", "mp3", "wav", "m4a":
		return true
	default:
		return false
	}
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

// Assets returns all ordered assets in the catalog.
func (c *Catalog) Assets() []Asset {
	return c.ordered
}

// MediaName returns the name of a media asset.
func (a Asset) MediaName() string { return a.mediaName }

// Extension returns the extension of a media asset.
func (a Asset) Extension() string { return a.extension }

// Data returns the raw content of a media asset.
func (a Asset) Data() []byte { return a.data }

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
