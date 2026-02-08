package pptx

// Theme represents a color palette for a presentation.
type Theme struct {
	Name       string
	Primary    string
	Secondary  string
	Accent     string
	Background string
	Text       string
	Light      string
	Dark       string
}

// Theme presets, parity with ppt-rs.
var (
	// ThemeCorporate - Professional and trustworthy
	ThemeCorporate = Theme{
		Name:       "Corporate",
		Primary:    "1565C0",
		Secondary:  "1976D2",
		Accent:     "FF6F00",
		Background: "FFFFFF",
		Text:       "212121",
		Light:      "E3F2FD",
		Dark:       "0D47A1",
	}

	// ThemeModern - Clean and simple
	ThemeModern = Theme{
		Name:       "Modern",
		Primary:    "212121",
		Secondary:  "757575",
		Accent:     "00BCD4",
		Background: "FAFAFA",
		Text:       "212121",
		Light:      "F5F5F5",
		Dark:       "424242",
	}

	// ThemeVibrant - Bold and colorful
	ThemeVibrant = Theme{
		Name:       "Vibrant",
		Primary:    "E91E63",
		Secondary:  "9C27B0",
		Accent:     "FF9800",
		Background: "FFFFFF",
		Text:       "212121",
		Light:      "FCE4EC",
		Dark:       "880E4F",
	}

	// ThemeDark - Easy on the eyes
	ThemeDark = Theme{
		Name:       "Dark",
		Primary:    "BB86FC",
		Secondary:  "03DAC6",
		Accent:     "CF6679",
		Background: "121212",
		Text:       "FFFFFF",
		Light:      "1E1E1E",
		Dark:       "000000",
	}

	// ThemeNature - Fresh and organic
	ThemeNature = Theme{
		Name:       "Nature",
		Primary:    "2E7D32",
		Secondary:  "4CAF50",
		Accent:     "8BC34A",
		Background: "FFFFFF",
		Text:       "1B5E20",
		Light:      "E8F5E9",
		Dark:       "1B5E20",
	}

	// ThemeTech - Modern technology feel
	ThemeTech = Theme{
		Name:       "Tech",
		Primary:    "0D47A1",
		Secondary:  "1976D2",
		Accent:     "00E676",
		Background: "FAFAFA",
		Text:       "263238",
		Light:      "E3F2FD",
		Dark:       "01579B",
	}

	// ThemeCarbon - IBM's design system
	ThemeCarbon = Theme{
		Name:       "Carbon",
		Primary:    "0043CE",
		Secondary:  "4589FF",
		Accent:     "24A148",
		Background: "FFFFFF",
		Text:       "161616",
		Light:      "E0E0E0",
		Dark:       "161616",
	}
)

// AllThemes returns all available theme presets.
func AllThemes() []Theme {
	return []Theme{
		ThemeCorporate,
		ThemeModern,
		ThemeVibrant,
		ThemeDark,
		ThemeNature,
		ThemeTech,
		ThemeCarbon,
	}
}
