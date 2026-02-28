package export_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
)

func TestHTMLOptionsThemeColors(t *testing.T) {
	opts := export.HTMLOptions{
		Theme: &export.ThemeColors{
			TitleColor: "#AA0000",
		},
	}
	slides := []elements.SlideContent{elements.NewSlide("Test")}

	html := export.HTMLWithOptions("Custom Theme", slides, opts)
	if !strings.Contains(html, "--title-color: #AA0000;") {
		t.Error("HTML missing injected custom theme variable")
	}
}

func TestHTMLOptionsNavigation(t *testing.T) {
	opts := export.HTMLOptions{
		IncludeNavigation: false,
	}
	slides := []elements.SlideContent{elements.NewSlide("Test")}

	html := export.HTMLWithOptions("No Nav", slides, opts)
	if strings.Contains(html, "<button id=\"btn-prev\"") {
		t.Error("HTML should not include navigation script when disabled")
	}
}
