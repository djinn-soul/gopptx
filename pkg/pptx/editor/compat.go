package editor

import (
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Metadata re-exports common.Metadata.
type Metadata = common.Metadata

// PresentationMetadata re-exports common.Metadata (for backward compatibility).
type PresentationMetadata = common.Metadata

// SlideMetadata re-exports common.SlideMetadata.
type SlideMetadata = common.SlideMetadata

// SlideSize re-exports common.SlideSize.
type SlideSize = common.SlideSize

// GetSlideSizeName returns the string representation of a common.SlideSize.
func GetSlideSizeName(size common.SlideSize) string {
	switch size {
	case common.SlideSize4x3():
		return "screen4x3"
	case common.SlideSize16x9():
		return "screen16x9"
	default:
		return ""
	}
}

func SlideSize4x3() SlideSize { return common.SlideSize4x3() }
func SlideSize16x9() SlideSize {
	return common.SlideSize16x9()
}
