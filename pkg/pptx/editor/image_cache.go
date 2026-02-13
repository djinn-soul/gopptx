package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type imagePathCacheEntry struct {
	data     []byte
	format   string
	partPath string
}

func (e *PresentationEditor) registerImageFromPath(imagePath, formatHint string) (string, string, error) {
	cleanPath := filepath.Clean(imagePath)
	e.imagePathMu.RLock()
	entry, ok := e.imagePathCache[cleanPath]
	e.imagePathMu.RUnlock()
	if ok && entry.partPath != "" {
		format := normalizeImageFormatHint(formatHint)
		if format == "" {
			format = entry.format
		}
		return entry.partPath, format, nil
	}

	data, format, err := e.loadImageFromPath(cleanPath, formatHint)
	if err != nil {
		return "", "", err
	}

	partPath, err := e.RegisterImage(data, format)
	if err != nil {
		return "", "", err
	}

	e.imagePathMu.Lock()
	e.imagePathCache[cleanPath] = imagePathCacheEntry{
		data:     nil,
		format:   format,
		partPath: partPath,
	}
	e.imagePathMu.Unlock()
	return partPath, format, nil
}

func (e *PresentationEditor) loadImageFromPath(cleanPath, formatHint string) ([]byte, string, error) {
	e.imagePathMu.RLock()
	entry, ok := e.imagePathCache[cleanPath]
	e.imagePathMu.RUnlock()
	if ok && len(entry.data) > 0 {
		format := normalizeImageFormatHint(formatHint)
		if format == "" {
			format = entry.format
		}
		return entry.data, format, nil
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, "", err
	}

	format := normalizeImageFormatHint(formatHint)
	if format == "" {
		format = normalizeImageFormatHint(filepath.Ext(cleanPath))
	}
	if format == "" {
		return nil, "", fmt.Errorf("image path %q has no detectable format", cleanPath)
	}

	e.imagePathMu.Lock()
	entry = e.imagePathCache[cleanPath]
	if len(entry.data) == 0 {
		entry.data = data
	}
	if entry.format == "" {
		entry.format = format
	}
	e.imagePathCache[cleanPath] = entry
	e.imagePathMu.Unlock()
	return entry.data, entry.format, nil
}

func normalizeImageFormatHint(format string) string {
	trimmed := strings.TrimSpace(strings.ToLower(format))
	return strings.TrimPrefix(trimmed, ".")
}
