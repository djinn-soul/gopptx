package editor

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// PresentationMetadata re-exports common.PresentationMetadata.
type PresentationMetadata = common.PresentationMetadata

// SlideMetadata re-exports common.SlideMetadata.
type SlideMetadata = common.SlideMetadata

// SlideSize re-exports common.SlideSize.
type SlideSize = common.SlideSize

func SlideSize4x3() SlideSize { return common.SlideSize4x3() }
func SlideSize16x9() SlideSize {
	return common.SlideSize16x9()
}
