package editor

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Re-export compatibility types and functions.
type (
	PresentationMetadata = common.PresentationMetadata
	SlideMetadata        = common.SlideMetadata
	SlideSize            = common.SlideSize
)

func SlideSize4x3() SlideSize { return common.SlideSize4x3() }
func SlideSize16x9() SlideSize {
	return common.SlideSize16x9()
}
