package pptx

import "github.com/djinn-soul/gopptx/pkg/pptx/markdown"

// SlidesFromMarkdown converts a markdown document into slide content.
// This is a compatibility alias for root pptx package.
func SlidesFromMarkdown(content string) ([]SlideContent, error) {
	return markdown.SlidesFromMarkdown(content)
}
