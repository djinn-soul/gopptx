package export

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

//go:embed html_style.css
var htmlStyle string

// HTML exports the presentation slides to a single HTML document.
func HTML(title string, slides []elements.SlideContent) string {
	var b strings.Builder

	// Header
	b.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	b.WriteString("<meta charset=\"UTF-8\">\n")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	b.WriteString(fmt.Sprintf("<title>%s</title>\n", escapeHTML(title)))

	// CSS
	b.WriteString("<style>\n")
	b.WriteString(htmlStyle)
	b.WriteString("</style>\n")

	b.WriteString("</head>\n<body>\n")

	// Title Slide (Presentation Title)
	b.WriteString("<div class=\"slide title-slide\">\n")
	b.WriteString(fmt.Sprintf("<h1>%s</h1>\n", escapeHTML(title)))
	b.WriteString("</div>\n")

	// Slides
	for i, slide := range slides {
		renderSlide(&b, slide, i+1)
	}

	b.WriteString("</body>\n</html>")
	return b.String()
}

func renderSlide(b *strings.Builder, slide elements.SlideContent, index int) {
	b.WriteString(fmt.Sprintf("<div class=\"slide\" id=\"slide-%d\">\n", index))

	// Slide Number
	b.WriteString(fmt.Sprintf("<div class=\"slide-number\">%d</div>\n", index))

	// Title
	slideTitle := slide.Title
	if slideTitle == "" {
		// Try to find title in shapes? ppt-rs relies on 'slide.title' field.
		// gopptx elements.SlideContent has explicit Title field.
		slideTitle = "Slide " + fmt.Sprintf("%d", index)
	}
	b.WriteString(fmt.Sprintf("<h2>%s</h2>\n", escapeHTML(slideTitle)))

	// Content Container
	b.WriteString("<div class=\"content\">\n")

	// Bullets / Content
	if len(slide.Bullets) > 0 {
		b.WriteString("<ul>\n")
		for _, bullet := range slide.Bullets {
			b.WriteString(fmt.Sprintf("<li>%s</li>\n", escapeHTML(bullet)))
		}
		b.WriteString("</ul>\n")
	}

	// Images
	for _, img := range slide.Images {
		if imgHTML := renderImage(img); imgHTML != "" {
			b.WriteString(imgHTML)
		}
	}

	// Extra content from shapes (if not bullets)
	// Iterate through shapes to find text that isn't the title.
	for _, shape := range slide.Shapes {
		if shape.Text != "" {
			b.WriteString(fmt.Sprintf("<p>%s</p>\n", escapeHTML(shape.Text)))
		}
	}

	b.WriteString("</div>\n") // content
	b.WriteString("</div>\n") // slide
}

func renderImage(img shapes.Image) string {
	if img.Path == "" {
		return ""
	}

	data, err := os.ReadFile(img.Path)
	if err != nil {
		// Fallback for missing images
		return fmt.Sprintf("<div class=\"image-container\">[Missing Image: %s]</div>\n", escapeHTML(img.Path))
	}

	mimeType := "application/octet-stream"
	ext := strings.ToLower(filepath.Ext(img.Path))
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".svg":
		mimeType = "image/svg+xml"
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf(
		"<div class=\"image-container\"><img src=\"data:%s;base64,%s\" alt=\"%s\" /></div>\n",
		mimeType, b64, escapeHTML(filepath.Base(img.Path)),
	)
}

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	"\"", "&quot;",
	"'", "&#39;",
)

func escapeHTML(s string) string {
	return htmlReplacer.Replace(s)
}
