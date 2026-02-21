package mermaid

import "github.com/djinn-soul/gopptx/pkg/pptx/styling"

// Theme defines the visual properties for a Mermaid diagram.
type Theme struct {
	Name            string
	Background      string // hex
	PrimaryFill     string // hex
	PrimaryStroke   string // hex
	SecondaryFill   string // hex
	SecondaryStroke string // hex
	TextColor       string // hex
	LineWeight      styling.Length
}

var themes = map[string]Theme{
	"default": {
		Name:            "default",
		Background:      "FFFFFF",
		PrimaryFill:     "ECECFF",
		PrimaryStroke:   "9370DB",
		SecondaryFill:   "F4F4F4",
		SecondaryStroke: "757575",
		TextColor:       "333333",
		LineWeight:      styling.Points(1),
	},
	"forest": {
		Name:            "forest",
		Background:      "FFFFFF",
		PrimaryFill:     "DDFFDD",
		PrimaryStroke:   "008000",
		SecondaryFill:   "F0FFF0",
		SecondaryStroke: "2E8B57",
		TextColor:       "004400",
		LineWeight:      styling.Points(1),
	},
	"dark": {
		Name:            "dark",
		Background:      "333333",
		PrimaryFill:     "1F2020",
		PrimaryStroke:   "81B1FF",
		SecondaryFill:   "444444",
		SecondaryStroke: "CCCCCC",
		TextColor:       "F0F0F0",
		LineWeight:      styling.Points(1),
	},
	"neutral": {
		Name:            "neutral",
		Background:      "FFFFFF",
		PrimaryFill:     "EEEEEE",
		PrimaryStroke:   "999999",
		SecondaryFill:   "F9F9F9",
		SecondaryStroke: "666666",
		TextColor:       "000000",
		LineWeight:      styling.Points(1),
	},
}

// GetTheme returns a Theme by name. If not found, returns the default theme.
func GetTheme(name string) Theme {
	if theme, ok := themes[name]; ok {
		return theme
	}
	return themes["default"]
}
