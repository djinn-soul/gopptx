package pptx

import (
	"fmt"
	"strings"
)

const (
	// SlideLayoutTitleAndContent renders title and optional bullet content shape.
	SlideLayoutTitleAndContent = "title_and_content"
	// SlideLayoutTitleOnly renders only the title shape.
	SlideLayoutTitleOnly = "title_only"
	// SlideLayoutBlank renders neither title nor content placeholders.
	SlideLayoutBlank = "blank"
)

const (
	slideLayoutXMLTitleAndContent = "titleAndContent"
	slideLayoutXMLTitleOnly       = "titleOnly"
	slideLayoutXMLBlank           = "blank"
)

// WithLayout sets the slide layout mode.
func (s SlideContent) WithLayout(layout string) SlideContent {
	s.Layout = normalizeSlideLayout(layout)
	return s
}

// WithTitleAndContentLayout sets the slide layout to title-and-content.
func (s SlideContent) WithTitleAndContentLayout() SlideContent {
	return s.WithLayout(SlideLayoutTitleAndContent)
}

// WithTitleOnlyLayout sets the slide layout to title-only.
func (s SlideContent) WithTitleOnlyLayout() SlideContent {
	return s.WithLayout(SlideLayoutTitleOnly)
}

// WithBlankLayout sets the slide layout to blank.
func (s SlideContent) WithBlankLayout() SlideContent {
	return s.WithLayout(SlideLayoutBlank)
}

func validateSlideLayout(s SlideContent, index int) error {
	switch normalizeSlideLayout(s.Layout) {
	case SlideLayoutTitleAndContent:
		if strings.TrimSpace(s.Title) == "" {
			return fmt.Errorf("slide %d title cannot be empty for title_and_content layout", index)
		}
	case SlideLayoutTitleOnly:
		if strings.TrimSpace(s.Title) == "" {
			return fmt.Errorf("slide %d title cannot be empty for title_only layout", index)
		}
		if len(s.Bullets) > 0 {
			return fmt.Errorf("slide %d title_only layout does not support bullets", index)
		}
	case SlideLayoutBlank:
		if strings.TrimSpace(s.Title) != "" {
			return fmt.Errorf("slide %d blank layout requires empty title", index)
		}
		if len(s.Bullets) > 0 {
			return fmt.Errorf("slide %d blank layout does not support bullets", index)
		}
	default:
		return fmt.Errorf("slide %d layout must be one of title_and_content|title_only|blank", index)
	}
	return nil
}

func slideLayoutXMLMode(layout string) string {
	switch normalizeSlideLayout(layout) {
	case SlideLayoutTitleOnly:
		return slideLayoutXMLTitleOnly
	case SlideLayoutBlank:
		return slideLayoutXMLBlank
	default:
		return slideLayoutXMLTitleAndContent
	}
}

func slideLayoutTarget(layout string) string {
	switch normalizeSlideLayout(layout) {
	case SlideLayoutTitleOnly:
		return "../slideLayouts/slideLayout2.xml"
	case SlideLayoutBlank:
		return "../slideLayouts/slideLayout3.xml"
	default:
		return "../slideLayouts/slideLayout1.xml"
	}
}

func normalizeSlideLayout(layout string) string {
	normalized := strings.ToLower(strings.TrimSpace(layout))
	switch normalized {
	case "", SlideLayoutTitleAndContent:
		return SlideLayoutTitleAndContent
	case "titleandcontent", "title-and-content":
		return SlideLayoutTitleAndContent
	case SlideLayoutTitleOnly:
		return SlideLayoutTitleOnly
	case "titleonly", "title-only":
		return SlideLayoutTitleOnly
	case SlideLayoutBlank:
		return SlideLayoutBlank
	default:
		return normalized
	}
}
