package slidesmeta

import (
	"fmt"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// BuildSlideContent converts an update request into public slide content.
func BuildSlideContent(
	request editorcommand.UpdateSlideRequest,
	currentTitle string,
) elements.SlideContent {
	title := request.Title
	if title == "" {
		title = currentTitle
	}

	slide := elements.NewSlide(title)
	if request.Layout != "" {
		slide = slide.WithLayout(request.Layout)
	}
	for _, bullet := range request.Bullets {
		slide = slide.AddBullet(bullet)
	}
	return slide
}

// ResolveThemeByName returns a built-in theme by its public name.
//
// Accepts both the gopptx theme names ("Corporate", "Dark") and the Office
// preset names ("facet", "ion"), so apply_theme and set_global_theme_preset
// take the same vocabulary. styling.ResolveTheme owns the mapping.
func ResolveThemeByName(name string) (styling.Theme, error) {
	theme, ok := styling.ResolveTheme(name)
	if !ok {
		return styling.Theme{}, fmt.Errorf("%w: %q", ErrUnknownThemeName, name)
	}
	return theme, nil
}
