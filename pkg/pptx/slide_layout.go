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
	// SlideLayoutCenteredTitle renders a vertically-centered title shape.
	SlideLayoutCenteredTitle = "centered_title"
	// SlideLayoutTitleAndBigContent renders a shorter title with larger content area.
	SlideLayoutTitleAndBigContent = "title_and_big_content"
	// SlideLayoutTwoColumn renders two content columns with bullet splitting.
	SlideLayoutTwoColumn = "two_column"
)

const (
	slideLayoutXMLTitleAndContent = "titleAndContent"
	slideLayoutXMLTitleOnly       = "titleOnly"
	slideLayoutXMLBlank           = "blank"
	slideLayoutXMLCenteredTitle   = "centeredTitle"
	slideLayoutXMLTitleBigContent = "titleAndBigContent"
	slideLayoutXMLTwoColumn       = "twoColumn"
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

// WithCenteredTitleLayout sets the slide layout to centered-title.
func (s SlideContent) WithCenteredTitleLayout() SlideContent {
	return s.WithLayout(SlideLayoutCenteredTitle)
}

// WithTitleAndBigContentLayout sets the slide layout to title-and-big-content.
func (s SlideContent) WithTitleAndBigContentLayout() SlideContent {
	return s.WithLayout(SlideLayoutTitleAndBigContent)
}

// WithTwoColumnLayout sets the slide layout to two-column.
func (s SlideContent) WithTwoColumnLayout() SlideContent {
	return s.WithLayout(SlideLayoutTwoColumn)
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
	case SlideLayoutCenteredTitle:
		if strings.TrimSpace(s.Title) == "" {
			return fmt.Errorf("slide %d title cannot be empty for centered_title layout", index)
		}
		if len(s.Bullets) > 0 {
			return fmt.Errorf("slide %d centered_title layout does not support bullets", index)
		}
	case SlideLayoutTitleAndBigContent:
		if strings.TrimSpace(s.Title) == "" {
			return fmt.Errorf("slide %d title cannot be empty for title_and_big_content layout", index)
		}
	case SlideLayoutTwoColumn:
		if strings.TrimSpace(s.Title) == "" {
			return fmt.Errorf("slide %d title cannot be empty for two_column layout", index)
		}
	default:
		return fmt.Errorf(
			"slide %d layout must be one of title_and_content|title_only|blank|centered_title|title_and_big_content|two_column",
			index,
		)
	}
	return nil
}

func slideLayoutXMLMode(layout string) string {
	switch normalizeSlideLayout(layout) {
	case SlideLayoutTitleOnly:
		return slideLayoutXMLTitleOnly
	case SlideLayoutBlank:
		return slideLayoutXMLBlank
	case SlideLayoutCenteredTitle:
		return slideLayoutXMLCenteredTitle
	case SlideLayoutTitleAndBigContent:
		return slideLayoutXMLTitleBigContent
	case SlideLayoutTwoColumn:
		return slideLayoutXMLTwoColumn
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
	case SlideLayoutCenteredTitle:
		return "../slideLayouts/slideLayout4.xml"
	case SlideLayoutTitleAndBigContent:
		return "../slideLayouts/slideLayout5.xml"
	case SlideLayoutTwoColumn:
		return "../slideLayouts/slideLayout6.xml"
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
	case SlideLayoutCenteredTitle:
		return SlideLayoutCenteredTitle
	case "centeredtitle", "centered-title":
		return SlideLayoutCenteredTitle
	case SlideLayoutTitleAndBigContent:
		return SlideLayoutTitleAndBigContent
	case "titleandbigcontent", "title-and-big-content", "big_content":
		return SlideLayoutTitleAndBigContent
	case SlideLayoutTwoColumn:
		return SlideLayoutTwoColumn
	case "twocolumn", "two-column":
		return SlideLayoutTwoColumn
	default:
		return normalized
	}
}
