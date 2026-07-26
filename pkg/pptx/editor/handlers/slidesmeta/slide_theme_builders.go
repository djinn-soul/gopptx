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
func ResolveThemeByName(name string) (styling.Theme, error) {
	switch name {
	case "Corporate":
		return styling.ThemeCorporate, nil
	case "Modern":
		return styling.ThemeModern, nil
	case "Vibrant":
		return styling.ThemeVibrant, nil
	case "Dark":
		return styling.ThemeDark, nil
	case "Nature":
		return styling.ThemeNature, nil
	case "Tech":
		return styling.ThemeTech, nil
	case "Carbon":
		return styling.ThemeCarbon, nil
	default:
		return styling.Theme{}, fmt.Errorf("%w: %q", ErrUnknownThemeName, name)
	}
}
