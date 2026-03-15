package media

import (
	"fmt"
	"os"
	"path"
	"strings"
)

const maxMediaPayloadBytes = 512 * 1024 * 1024 // 512 MiB

func ResolveVideoExtension(mimeType string, filePath string, hasInlineData bool) (string, error) {
	return resolveMediaExtension("video", mimeType, filePath, hasInlineData, videoExtForMIME, isAllowedVideoExt)
}

func ResolveAudioExtension(mimeType string, filePath string, hasInlineData bool) (string, error) {
	return resolveMediaExtension("audio", mimeType, filePath, hasInlineData, audioExtForMIME, isAllowedAudioExt)
}

func resolveMediaExtension(
	kind string,
	mimeType string,
	filePath string,
	hasInlineData bool,
	resolveMIME func(string) (string, bool),
	isAllowedExt func(string) bool,
) (string, error) {
	normalizedMime := normalizeMIMEType(mimeType)
	if normalizedMime != "" {
		if ext, ok := resolveMIME(normalizedMime); ok {
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
	if isAllowedExt(ext) {
		return ext, nil
	}
	return "", fmt.Errorf("unsupported %s file extension: %s", kind, ext)
}

func videoExtForMIME(mimeType string) (string, bool) {
	switch mimeType {
	case "video/mp4":
		return "mp4", true
	case "video/webm":
		return "webm", true
	case "video/x-msvideo", "video/avi":
		return "avi", true
	case "video/x-ms-wmv", "video/wmv":
		return "wmv", true
	case "video/quicktime":
		return "mov", true
	case "video/x-matroska", "video/mkv":
		return "mkv", true
	case "video/x-m4v", "video/m4v":
		return "m4v", true
	default:
		return "", false
	}
}

func audioExtForMIME(mimeType string) (string, bool) {
	switch mimeType {
	case "audio/mpeg", "audio/mp3":
		return "mp3", true
	case "audio/wav", "audio/x-wav":
		return "wav", true
	case "audio/m4a", "audio/mp4":
		return "m4a", true
	case "audio/x-ms-wma", "audio/wma":
		return "wma", true
	case "audio/ogg":
		return "ogg", true
	case "audio/flac":
		return "flac", true
	case "audio/aac":
		return "aac", true
	default:
		return "", false
	}
}

func isAllowedVideoExt(ext string) bool {
	switch ext {
	case "mp4", "webm", "avi", "wmv", "mov", "mkv", "m4v":
		return true
	default:
		return false
	}
}

func isAllowedAudioExt(ext string) bool {
	switch ext {
	case "mp3", "wav", "m4a", "wma", "ogg", "flac", "aac":
		return true
	default:
		return false
	}
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
