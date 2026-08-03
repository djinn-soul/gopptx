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
	templateMainContentType     = "application/vnd.openxmlformats-officedocument.presentationml.template.main+xml"
	macroPresentationMainType   = "application/vnd.ms-powerpoint.presentation.macroEnabled.main+xml"
	macroTemplateMainType       = "application/vnd.ms-powerpoint.template.macroEnabled.main+xml"
	macroSlideshowMainType      = "application/vnd.ms-powerpoint.slideshow.macroEnabled.main+xml"
)

// Output extensions that carry a VBA project.
const (
	extMacroPresentation = ".pptm"
	extMacroTemplate     = ".potm"
	extMacroSlideshow    = ".ppsm"
)

// mainPartContentTypes lists every content type ppt/presentation.xml can carry.
// A package opened from any of them holds the same slides; only this one
// declaration says which PowerPoint file kind it is, which is why a .potx can
// be read and then saved as a .pptx (upstream issue #1070).
func mainPartContentTypes() []string {
	return []string{
		presentationMainContentType,
		slideshowMainContentType,
		templateMainContentType,
		macroPresentationMainType,
		macroTemplateMainType,
		macroSlideshowMainType,
	}
}

// mainContentTypeForExtension maps an output extension to the content type its
// main part has to declare. An extension that is not listed keeps whatever the
// package already declares.
func mainContentTypeForExtension(ext string) string {
	switch strings.ToLower(ext) {
	case ".pptx":
		return presentationMainContentType
	case ".ppsx":
		return slideshowMainContentType
	case ".potx":
		return templateMainContentType
	case extMacroPresentation:
		return macroPresentationMainType
	case extMacroTemplate:
		return macroTemplateMainType
	case extMacroSlideshow:
		return macroSlideshowMainType
	default:
		return ""
	}
}

// isMacroEnabledExtension reports whether an output extension can carry a VBA
// project. A macro-enabled template or slideshow keeps its macros just as a
// .pptm does, so rejecting them would make those kinds unreachable.
func isMacroEnabledExtension(ext string) bool {
	switch strings.ToLower(ext) {
	case extMacroPresentation, extMacroTemplate, extMacroSlideshow:
		return true
	default:
		return false
	}
}

// applyOutputPresentationKind rewrites the main part's content type so it
// matches the file being written. Without it a deck opened from a template
// saves as a .pptx that PowerPoint still treats as a template.
func applyOutputPresentationKind(parts map[string][]byte, filePath string) error {
	wanted := mainContentTypeForExtension(filepath.Ext(filePath))
	if wanted == "" {
		return nil
	}
	contentTypes := parts[common.ContentTypesPath]
	wantedBytes := []byte(wanted)
	if bytes.Contains(contentTypes, wantedBytes) {
		return nil
	}
	for _, current := range mainPartContentTypes() {
		currentBytes := []byte(current)
		if !bytes.Contains(contentTypes, currentBytes) {
			continue
		}
		parts[common.ContentTypesPath] = bytes.Replace(
			contentTypes, currentBytes, wantedBytes, 1,
		)
		return nil
	}
	return errors.New(
		"cannot save " + filepath.Ext(filePath) +
			": the package declares no presentation main content type",
	)
}
