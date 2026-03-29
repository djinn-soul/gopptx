package elements

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func NormalizeSlideLayout(layout string) string {
	normalized := strings.ToLower(strings.TrimSpace(layout))
	switch normalized {
	case "", strings.ToLower(SlideLayoutTitleAndContent), "titleandcontent", "title-and-content":
		return SlideLayoutTitleAndContent
	case strings.ToLower(SlideLayoutTitle), "title slide":
		return SlideLayoutTitle
	case strings.ToLower(SlideLayoutSectionHeader), "section header":
		return SlideLayoutSectionHeader
	case strings.ToLower(SlideLayoutTwoContent), "two content":
		return SlideLayoutTwoContent
	case strings.ToLower(SlideLayoutComparison), "comparison":
		return SlideLayoutComparison
	case strings.ToLower(SlideLayoutContentCaption), "content with caption":
		return SlideLayoutContentCaption
	case strings.ToLower(SlideLayoutPictureCaption), "picture with caption":
		return SlideLayoutPictureCaption
	case strings.ToLower(SlideLayoutTitleOnly), "titleonly", "title-only":
		return SlideLayoutTitleOnly
	case strings.ToLower(SlideLayoutBlank), "blank":
		return SlideLayoutBlank
	case strings.ToLower(SlideLayoutCenteredTitle), "centeredtitle", "centered-title":
		return SlideLayoutCenteredTitle
	case strings.ToLower(SlideLayoutTitleAndBigContent), "titleandbigcontent", "title-and-big-content", "big_content":
		return SlideLayoutTitleAndBigContent
	case strings.ToLower(SlideLayoutTwoColumn), "twocolumn", "two-column", "two column":
		return SlideLayoutTwoColumn
	default:
		return normalized
	}
}

func SlideLayoutXMLMode(layout string) string {
	switch NormalizeSlideLayout(layout) {
	case SlideLayoutTitleOnly:
		return "titleOnly"
	case SlideLayoutBlank:
		return "blank"
	case SlideLayoutCenteredTitle:
		return "centeredTitle"
	case SlideLayoutTitleAndBigContent:
		return "titleAndBigContent"
	case SlideLayoutTwoColumn:
		return "twoColumn"
	default:
		return "titleAndContent"
	}
}

func SlideLayoutTarget(layout string) string {
	switch NormalizeSlideLayout(layout) {
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

// WithPlaceholderText overrides a placeholder with text using the default placeholder type.
func (s SlideContent) WithPlaceholderText(index int, text string) SlideContent {
	return s.WithPlaceholderTextAs(index, defaultPlaceholderTextType(index), text)
}

// WithPlaceholderTextAs overrides a placeholder with text and explicit placeholder type.
func (s SlideContent) WithPlaceholderTextAs(index int, placeholderType, text string) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index: index,
		Type:  placeholderType,
		Text:  text,
	})
	return s
}

// WithPlaceholderImage overrides a placeholder with an image using the default placeholder type.
func (s SlideContent) WithPlaceholderImage(index int, img shapes.Image) SlideContent {
	return s.WithPlaceholderImageAs(index, defaultPlaceholderImageType(index), img)
}

// WithPlaceholderImageAs overrides a placeholder with an image and explicit placeholder type.
func (s SlideContent) WithPlaceholderImageAs(index int, placeholderType string, img shapes.Image) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index: index,
		Type:  placeholderType,
		Image: &img,
	})
	return s
}

// WithPlaceholderTable overrides a placeholder with a table using the default placeholder type.
func (s SlideContent) WithPlaceholderTable(index int, table tables.Table) SlideContent {
	return s.WithPlaceholderTableAs(index, defaultPlaceholderTextType(index), table)
}

// WithPlaceholderTableAs overrides a placeholder with a table and explicit placeholder type.
func (s SlideContent) WithPlaceholderTableAs(index int, placeholderType string, table tables.Table) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index: index,
		Type:  placeholderType,
		Table: &table,
	})
	return s
}

// WithPlaceholderChart overrides a placeholder with a chart using the default placeholder type.
func (s SlideContent) WithPlaceholderChart(index int, chart ChartDefinition) SlideContent {
	return s.WithPlaceholderChartAs(index, defaultPlaceholderTextType(index), chart)
}

// WithPlaceholderChartAs overrides a placeholder with a chart and explicit placeholder type.
func (s SlideContent) WithPlaceholderChartAs(index int, placeholderType string, chart ChartDefinition) SlideContent {
	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index: index,
		Type:  placeholderType,
		Chart: chart,
	})
	return s
}

// WithPlaceholderOverride adds custom geometry or style overrides to a placeholder.
func (s SlideContent) WithPlaceholderOverride(
	target shapes.PlaceholderTarget,
	options shapes.PlaceholderOverrideOptions,
) SlideContent {
	// If no type/index specified, default to title (index 0)
	if target.Type == "" && target.Index == 0 && target.Name == "" {
		target.Index = 0
		target.Type = placeholderTypeTitle
	}

	s.PlaceholderOverrides = append(s.PlaceholderOverrides, shapes.PlaceholderContent{
		Index:    target.Index,
		Type:     target.Type,
		Target:   &target,
		Override: &options,
	})
	return s
}

func defaultPlaceholderTextType(index int) string {
	if index == 0 {
		return placeholderTypeTitle
	}
	return placeholderTypeBody
}

func defaultPlaceholderImageType(index int) string {
	if index == 0 {
		return placeholderTypeTitle
	}
	return placeholderTypePic
}
