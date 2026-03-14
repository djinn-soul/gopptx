package pptx

import "github.com/djinn-soul/gopptx/pkg/pptx/markdown"

// SlidesFromMarkdown converts a markdown document into slide content.
// This is a compatibility alias for root pptx package.
func SlidesFromMarkdown(content string) ([]SlideContent, error) {
	return markdown.SlidesFromMarkdown(content)
}

// SlidesFromMarkdownWithOptions converts markdown into slide content using parser options.
func SlidesFromMarkdownWithOptions(content string, options markdown.ParseOptions) ([]SlideContent, error) {
	return markdown.SlidesFromMarkdownWithOptions(content, options)
}

// SlidesFromMarkdownFile reads markdown from a file and resolves relative assets from that file directory.
func SlidesFromMarkdownFile(path string) ([]SlideContent, error) {
	return markdown.SlidesFromMarkdownFile(path)
}
