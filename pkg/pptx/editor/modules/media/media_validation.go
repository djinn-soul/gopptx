package media

import (
	"fmt"
	"os"
	"path"
	"strings"
)

const maxMediaPayloadBytes = 512 * 1024 * 1024 // 512 MiB

var (
	videoMimeToExt = map[string]string{
		"video/mp4":      "mp4",
		"video/webm":     "webm",
		"video/x-msvideo": "avi",
		"video/avi":      "avi",
		"video/quicktime": "mov",
	}
	audioMimeToExt = map[string]string{
		"audio/mpeg": "mp3",
		"audio/mp3":  "mp3",
		"audio/wav":  "wav",
		"audio/x-wav": "wav",
		"audio/m4a":  "m4a",
		"audio/mp4":  "m4a",
	}
)

func ResolveVideoExtension(mimeType string, filePath string, hasInlineData bool) (string, error) {
	return resolveMediaExtension("video", mimeType, filePath, hasInlineData, videoMimeToExt)
}

func ResolveAudioExtension(mimeType string, filePath string, hasInlineData bool) (string, error) {
	return resolveMediaExtension("audio", mimeType, filePath, hasInlineData, audioMimeToExt)
}

func resolveMediaExtension(
	kind string,
	mimeType string,
	filePath string,
	hasInlineData bool,
	mimeToExt map[string]string,
) (string, error) {
	normalizedMime := normalizeMIMEType(mimeType)
	if normalizedMime != "" {
		if ext, ok := mimeToExt[normalizedMime]; ok {
			return ext, nil
		}
		return "", fmt.Errorf("unsupported %s mime type: %s", kind, mimeType)
	}

	if hasInlineData {
		return "", fmt.Errorf("%s mime type is required when using in-memory data", kind)
	}

	ext := strings.TrimPrefix(strings.ToLower(path.Ext(filePath)), ".")
	if ext == "" {
		return "", fmt.Errorf("%s file extension is required when mime type is empty", kind)
	}
	for _, allowed := range mimeToExt {
		if ext == allowed {
			return ext, nil
		}
	}
	return "", fmt.Errorf("unsupported %s file extension: %s", kind, ext)
}

func ValidateMediaPayloadSize(data []byte, kind string) error {
	if len(data) == 0 {
		return fmt.Errorf("%s data cannot be empty", kind)
	}
	if len(data) > maxMediaPayloadBytes {
		return fmt.Errorf("%s exceeds max size of %d bytes", kind, maxMediaPayloadBytes)
	}
	return nil
}

func ReadMediaFileWithSizeLimit(filePath string, kind string) ([]byte, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	if info.Size() > maxMediaPayloadBytes {
		return nil, fmt.Errorf("%s exceeds max size of %d bytes", kind, maxMediaPayloadBytes)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if err := ValidateMediaPayloadSize(data, kind); err != nil {
		return nil, err
	}
	return data, nil
}

func normalizeMIMEType(mimeType string) string {
	normalized := strings.TrimSpace(strings.ToLower(mimeType))
	if normalized == "" {
		return ""
	}
	if idx := strings.Index(normalized, ";"); idx >= 0 {
		normalized = strings.TrimSpace(normalized[:idx])
	}
	return normalized
}
