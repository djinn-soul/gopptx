package styling

// ColorScheme represents a set of 12 theme colors.
type ColorScheme struct {
	Name     string
	Dk1      string // Dark 1
	Lt1      string // Light 1
	Dk2      string // Dark 2
	Lt2      string // Light 2
	Accent1  string
	Accent2  string
	Accent3  string
	Accent4  string
	Accent5  string
	Accent6  string
	Hlink    string
	FolHlink string
}

// FontScheme represents a set of theme fonts.
type FontScheme struct {
	Name      string
	MajorFont string // Heading typeface
	MinorFont string // Body typeface
}

// Theme represents a complete color and font palette for a presentation.
type Theme struct {
	Name   string
	Colors ColorScheme
	Fonts  FontScheme

	// Legacy fields for backward compatibility
	Primary    string
	Secondary  string
	Accent     string
	Background string
	Text       string
	Light      string
	Dark       string
}

// Theme presets, parity with ppt-rs.
//
//nolint:gochecknoglobals // theme presets
var (
	// ThemeCorporate - Professional and trustworthy.
	ThemeCorporate = Theme{
		Name: "Corporate",
		Colors: ColorScheme{
			Name: "Corporate",
			Dk1:  "000000", Lt1: "FFFFFF", Dk2: "1F497D", Lt2: "EEECE1",
			Accent1: "4F81BD", Accent2: "C0504D", Accent3: "9BBB59",
			Accent4: "8064A2", Accent5: "4BACC6", Accent6: "F79646",
			Hlink: "0000FF", FolHlink: "800080",
		},
		Fonts: FontScheme{
			Name:      "Corporate",
			MajorFont: "Calibri",
			MinorFont: "Calibri",
		},
		Primary: "1565C0", Secondary: "1976D2", Accent: "FF6F00",
		Background: "FFFFFF", Text: "212121", Light: "E3F2FD", Dark: "0D47A1",
	}

	// ThemeModern - Clean and simple.
	ThemeModern = Theme{
		Name: "Modern",
		Colors: ColorScheme{
			Name: "Modern",
			Dk1:  "000000", Lt1: "FFFFFF", Dk2: "444444", Lt2: "F3F3F3",
			Accent1: "212121", Accent2: "757575", Accent3: "00BCD4",
			Accent4: "0097A7", Accent5: "00838F", Accent6: "006064",
			Hlink: "0000EE", FolHlink: "551A8B",
		},
		Fonts: FontScheme{
			Name:      "Modern",
			MajorFont: "Inter",
			MinorFont: "Inter",
		},
		Primary: "212121", Secondary: "757575", Accent: "00BCD4",
		Background: "FAFAFA", Text: "212121", Light: "F5F5F5", Dark: "424242",
	}

	// ThemeVibrant - Bold and colorful.
	ThemeVibrant = Theme{
		Name: "Vibrant",
		Colors: ColorScheme{
			Name: "Vibrant",
			Dk1:  "000000", Lt1: "FFFFFF", Dk2: "E91E63", Lt2: "FCE4EC",
			Accent1: "FF4081", Accent2: "9C27B0", Accent3: "FF9800",
			Accent4: "FFC107", Accent5: "FF5722", Accent6: "E91E63",
			Hlink: "3F51B5", FolHlink: "1A237E",
		},
		Fonts: FontScheme{
			Name:      "Vibrant",
			MajorFont: "Outfit",
			MinorFont: "Outfit",
		},
		Primary: "E91E63", Secondary: "9C27B0", Accent: "FF9800",
		Background: "FFFFFF", Text: "212121", Light: "FCE4EC", Dark: "880E4F",
	}

	// ThemeDark - Easy on the eyes.
	ThemeDark = Theme{
		Name: "Dark",
		Colors: ColorScheme{
			Name: "Dark",
			Dk1:  "FFFFFF", Lt1: "121212", Dk2: "BB86FC", Lt2: "1E1E1E",
			Accent1: "03DAC6", Accent2: "CF6679", Accent3: "FF9800",
			Accent4: "03DAC6", Accent5: "BB86FC", Accent6: "CF6679",
			Hlink: "03DAC6", FolHlink: "BB86FC",
		},
		Fonts: FontScheme{
			Name:      "Dark",
			MajorFont: "Roboto",
			MinorFont: "Roboto",
		},
		Primary: "BB86FC", Secondary: "03DAC6", Accent: "CF6679",
		Background: "121212", Text: "FFFFFF", Light: "1E1E1E", Dark: "000000",
	}

	// ThemeNature - Fresh and organic.
	ThemeNature = Theme{
		Name: "Nature",
		Colors: ColorScheme{
			Name: "Nature",
			Dk1:  "1B5E20", Lt1: "FFFFFF", Dk2: "2E7D32", Lt2: "E8F5E9",
			Accent1: "4CAF50", Accent2: "8BC34A", Accent3: "CDDC39",
			Accent4: "FFEB3B", Accent5: "FFC107", Accent6: "FF9800",
			Hlink: "2E7D32", FolHlink: "1B5E20",
		},
		Fonts: FontScheme{
			Name:      "Nature",
			MajorFont: "Inter",
			MinorFont: "Inter",
		},
		Primary: "2E7D32", Secondary: "4CAF50", Accent: "8BC34A",
		Background: "FFFFFF", Text: "1B5E20", Light: "E8F5E9", Dark: "1B5E20",
	}

	// ThemeTech - Modern technology feel.
	ThemeTech = Theme{
		Name: "Tech",
		Colors: ColorScheme{
			Name: "Tech",
			Dk1:  "01579B", Lt1: "FAFAFA", Dk2: "0D47A1", Lt2: "E3F2FD",
			Accent1: "1976D2", Accent2: "00E676", Accent3: "00B0FF",
			Accent4: "00B8D4", Accent5: "0091EA", Accent6: "0D47A1",
			Hlink: "00B0FF", FolHlink: "0D47A1",
		},
		Fonts: FontScheme{
			Name:      "Tech",
			MajorFont: "Roboto",
			MinorFont: "Roboto",
		},
		Primary: "0D47A1", Secondary: "1976D2", Accent: "00E676",
		Background: "FAFAFA", Text: "263238", Light: "E3F2FD", Dark: "01579B",
	}

	// ThemeCarbon - IBM's design system.
	ThemeCarbon = Theme{
		Name: "Carbon",
		Colors: ColorScheme{
			Name: "Carbon",
			Dk1:  "161616", Lt1: "FFFFFF", Dk2: "0043CE", Lt2: "E0E0E0",
			Accent1: "4589FF", Accent2: "24A148", Accent3: "08BDBA",
			Accent4: "1192E8", Accent5: "6929C4", Accent6: "0043CE",
			Hlink: "4589FF", FolHlink: "0043CE",
		},
		Fonts: FontScheme{
			Name:      "Carbon",
			MajorFont: "IBM Plex Sans",
			MinorFont: "IBM Plex Sans",
		},
		Primary: "0043CE", Secondary: "4589FF", Accent: "24A148",
		Background: "FFFFFF", Text: "161616", Light: "E0E0E0", Dark: "161616",
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
