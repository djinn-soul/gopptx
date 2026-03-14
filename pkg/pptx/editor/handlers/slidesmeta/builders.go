package slidesmeta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

var (
	ErrUnsupportedChartType = errors.New("unsupported chart type")
	ErrUnknownThemeName     = errors.New("unknown theme name")
)

func BuildChartDefinition(request editorcommand.AddChartRequest) (charts.ChartDefinition, error) {
	switch strings.ToLower(request.ChartType) {
	case "bar":
		chart := charts.NewBarChart(request.Categories, request.Values).WithTitle(request.Title)
		if request.W > 0 {
			chart = chart.Size(styling.Emu(request.W), styling.Emu(request.H)).
				Position(styling.Emu(request.X), styling.Emu(request.Y))
		}
		return chart, nil
	case "line":
		chart := charts.NewLineChart(request.Categories, request.Values).WithTitle(request.Title)
		if request.W > 0 {
			chart = chart.Size(styling.Emu(request.W), styling.Emu(request.H)).
				Position(styling.Emu(request.X), styling.Emu(request.Y))
		}
		return chart, nil
	case "pie":
		chart := charts.NewPieChart(request.Categories, request.Values).WithTitle(request.Title)
		if request.W > 0 {
			chart = chart.Size(styling.Emu(request.W), styling.Emu(request.H)).
				Position(styling.Emu(request.X), styling.Emu(request.Y))
		}
		return chart, nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedChartType, request.ChartType)
	}
}

func BuildSlideContent(request editorcommand.UpdateSlideRequest, currentTitle string) elements.SlideContent {
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
