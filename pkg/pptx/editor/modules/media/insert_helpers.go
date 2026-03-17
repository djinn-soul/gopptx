package media

import (
	"errors"
	"os"
	"path"
	"strings"
)

const oleShapeIDStride = 10000

//nolint:gochecknoglobals // Reused escaper table avoids repeated allocations in hot XML-building paths.
var xmlAttrEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&apos;",
)

func ValidateMediaSlideIndex(slideIndex, slideCount int) error {
	if slideIndex < 0 || slideIndex >= slideCount {
		return errors.New("slide index out of range")
	}
	return nil
}

func RegisterPartFromDataOrPath(
	data []byte,
	filePath string,
	missingErr string,
	fromData func([]byte) (string, error),
	fromPath func(string) (string, error),
) (string, error) {
	if len(data) > 0 {
		return fromData(data)
	}
	if filePath != "" {
		return fromPath(filePath)
	}
	return "", errors.New(missingErr)
}

func RegisterVideoPart(
	videoData []byte,
	videoPath string,
	mimeType string,
	registerMedia func([]byte, string) (string, error),
) (string, error) {
	videoExt, err := ResolveVideoExtension(mimeType, videoPath, len(videoData) > 0)
	if err != nil {
		return "", err
	}
	return RegisterPartFromDataOrPath(
		videoData,
		videoPath,
		"video data or path is required",
		func(data []byte) (string, error) {
			if sizeErr := ValidateMediaPayloadSize(data, "video"); sizeErr != nil {
				return "", sizeErr
			}
			return registerMedia(data, videoExt)
		},
		func(filePath string) (string, error) {
			data, readErr := ReadMediaFileWithSizeLimit(filePath, "video")
			if readErr != nil {
				return "", readErr
			}
			return registerMedia(data, videoExt)
		},
	)
}

func RegisterEmbeddingPart(
	objectData []byte,
	objectPath string,
	registerEmbedding func([]byte, string) (string, error),
) (string, error) {
	return RegisterPartFromDataOrPath(
		objectData,
		objectPath,
		"object data or path is required",
		func(data []byte) (string, error) {
			return registerEmbedding(data, "bin")
		},
		func(filePath string) (string, error) {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return "", err
			}
			return registerEmbedding(data, strings.TrimPrefix(path.Ext(filePath), "."))
		},
	)
}

func RegisterAudioPart(
	audioData []byte,
	audioPath string,
	mimeType string,
	registerMedia func([]byte, string) (string, error),
) (string, error) {
	audioExt, err := ResolveAudioExtension(mimeType, audioPath, len(audioData) > 0)
	if err != nil {
		return "", err
	}
	return RegisterPartFromDataOrPath(
		audioData,
		audioPath,
		"audio data or path is required",
		func(data []byte) (string, error) {
			if sizeErr := ValidateMediaPayloadSize(data, "audio"); sizeErr != nil {
				return "", sizeErr
			}
			return registerMedia(data, audioExt)
		},
		func(filePath string) (string, error) {
			data, readErr := ReadMediaFileWithSizeLimit(filePath, "audio")
			if readErr != nil {
				return "", readErr
			}
			return registerMedia(data, audioExt)
		},
	)
}

func escapeXMLAttr(value string) string {
	return xmlAttrEscaper.Replace(value)
}
