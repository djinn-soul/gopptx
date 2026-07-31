package editor

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	presentationMainContentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"
	slideshowMainContentType    = "application/vnd.openxmlformats-officedocument.presentationml.slideshow.main+xml"
)

func applyOutputPresentationKind(parts map[string][]byte, filePath string) error {
	if !strings.EqualFold(filepath.Ext(filePath), ".ppsx") {
		return nil
	}
	contentTypes := parts[common.ContentTypesPath]
	presentationType := []byte(presentationMainContentType)
	slideshowType := []byte(slideshowMainContentType)
	if bytes.Contains(contentTypes, slideshowType) {
		return nil
	}
	if !bytes.Contains(contentTypes, presentationType) {
		return errors.New(
			"cannot save .ppsx: presentation main content type is missing",
		)
	}
	parts[common.ContentTypesPath] = bytes.Replace(
		contentTypes,
		presentationType,
		slideshowType,
		1,
	)
	return nil
}
