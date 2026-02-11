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

var (
	SlideSize4x3  = common.SlideSize4x3
	SlideSize16x9 = common.SlideSize16x9
)
