package templates

import (
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// LayoutOverrides maps template slide index -> layout override.
//
// Indices are zero-based.
type LayoutOverrides map[int]string

func applyLayoutOverrides(slides []elements.SlideContent, overrides LayoutOverrides) ([]elements.SlideContent, error) {
	if len(overrides) == 0 {
		return slides, nil
	}

	for idx, layout := range overrides {
		if idx < 0 || idx >= len(slides) {
			return nil, fmt.Errorf("layout override index %d out of range [0,%d)", idx, len(slides))
		}
		normalized := elements.NormalizeSlideLayout(layout)
		if !isSupportedTemplateLayout(normalized) {
			return nil, fmt.Errorf("unsupported layout override %q for slide index %d", layout, idx)
		}
		slides[idx] = slides[idx].WithLayout(normalized)
	}

	return slides, nil
}

func isSupportedTemplateLayout(layout string) bool {
	switch layout {
	case elements.SlideLayoutTitleAndContent,
		elements.SlideLayoutTitleOnly,
		elements.SlideLayoutBlank,
		elements.SlideLayoutCenteredTitle,
		elements.SlideLayoutTitleAndBigContent,
		elements.SlideLayoutTwoColumn,
		elements.SlideLayoutTitle,
		elements.SlideLayoutSectionHeader,
		elements.SlideLayoutTwoContent,
		elements.SlideLayoutComparison,
		elements.SlideLayoutContentCaption,
		elements.SlideLayoutPictureCaption:
		return true
	default:
		return false
	}
}
