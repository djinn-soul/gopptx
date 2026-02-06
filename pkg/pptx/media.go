package pptx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type mediaAsset struct {
	mediaName string
	extension string
	data      []byte
}

type mediaCatalog struct {
	byPath  map[string]mediaAsset
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
		byPath:  make(map[string]mediaAsset),
		ordered: make([]mediaAsset, 0),
	}

	for slideIndex, slide := range slides {
		for imageIndex, image := range slide.Images {
			cleanPath := filepath.Clean(image.Path)
			if _, exists := catalog.byPath[cleanPath]; exists {
				continue
			}

			ext := normalizeImageExtension(cleanPath)
			if _, ok := supportedImageExtensions[ext]; !ok {
				return nil, fmt.Errorf(
					"slide %d image %d has unsupported extension %q",
					slideIndex+1,
					imageIndex+1,
					ext,
				)
			}

			data, err := os.ReadFile(cleanPath)
			if err != nil {
				return nil, fmt.Errorf(
					"slide %d image %d read error: %w",
					slideIndex+1,
					imageIndex+1,
					err,
				)
			}
			if len(data) == 0 {
				return nil, fmt.Errorf("slide %d image %d is empty", slideIndex+1, imageIndex+1)
			}

			mediaName := fmt.Sprintf("image%d.%s", len(catalog.ordered)+1, ext)
			asset := mediaAsset{
				mediaName: mediaName,
				extension: ext,
				data:      data,
			}
			catalog.byPath[cleanPath] = asset
			catalog.ordered = append(catalog.ordered, asset)
		}
	}

	return catalog, nil
}

func normalizeImageExtension(path string) string {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	return ext
}

func (c *mediaCatalog) mediaNameForPath(path string) (string, bool) {
	asset, ok := c.byPath[filepath.Clean(path)]
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
