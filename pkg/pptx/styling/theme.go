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

// Theme preset names. Each appears on the Theme, its ColorScheme and its
// FontScheme, so they are named once here.
const (
	themeNameCorporate = "Corporate"
	themeNameModern    = "Modern"
	themeNameVibrant   = "Vibrant"
	themeNameDark      = "Dark"
	themeNameNature    = "Nature"
	themeNameTech      = "Tech"
	themeNameCarbon    = "Carbon"
)

// Font families shared across theme presets.
const (
	fontInter  = "Inter"
	fontRoboto = "Roboto"
)

// Palette colors repeated across presets that have no entry in colors.go.
const (
	colorTechBlue900    = "0D47A1"
	colorTechBlue700    = "1976D2"
	colorNeutral50      = "FAFAFA"
	colorNatureGreen900 = "1B5E20"
	colorDarkTeal       = "03DAC6"
	colorTechBlue50     = "E3F2FD"
	colorDarkPurple     = "BB86FC"
	colorDarkPink       = "CF6679"
)

// Theme presets, parity with ppt-rs.
//
//nolint:gochecknoglobals // theme presets
var (
	// ThemeCorporate - Professional and trustworthy.
	ThemeCorporate = Theme{
		Name: themeNameCorporate,
		Colors: ColorScheme{
			Name: themeNameCorporate,
			Dk1:  "000000", Lt1: ColorWhite, Dk2: "1F497D", Lt2: "EEECE1",
			Accent1: "4F81BD", Accent2: "C0504D", Accent3: "9BBB59",
			Accent4: "8064A2", Accent5: "4BACC6", Accent6: "F79646",
			Hlink: "0000FF", FolHlink: "800080",
		},
		Fonts: FontScheme{
			Name:      themeNameCorporate,
			MajorFont: "Calibri",
			MinorFont: "Calibri",
		},
		Primary: "1565C0", Secondary: colorTechBlue700, Accent: "FF6F00",
		Background: ColorWhite, Text: "212121", Light: colorTechBlue50, Dark: colorTechBlue900,
	}

	// ThemeModern - Clean and simple.
	ThemeModern = Theme{
		Name: themeNameModern,
		Colors: ColorScheme{
			Name: themeNameModern,
			Dk1:  "000000", Lt1: ColorWhite, Dk2: "444444", Lt2: "F3F3F3",
			Accent1: "212121", Accent2: "757575", Accent3: "00BCD4",
			Accent4: "0097A7", Accent5: "00838F", Accent6: "006064",
			Hlink: "0000EE", FolHlink: "551A8B",
		},
		Fonts: FontScheme{
			Name:      themeNameModern,
			MajorFont: fontInter,
			MinorFont: fontInter,
		},
		Primary: "212121", Secondary: "757575", Accent: "00BCD4",
		Background: colorNeutral50, Text: "212121", Light: "F5F5F5", Dark: "424242",
	}

	// ThemeVibrant - Bold and colorful.
	ThemeVibrant = Theme{
		Name: themeNameVibrant,
		Colors: ColorScheme{
			Name: themeNameVibrant,
			Dk1:  "000000", Lt1: ColorWhite, Dk2: ColorMaterialPink, Lt2: "FCE4EC",
			Accent1: "FF4081", Accent2: "9C27B0", Accent3: ColorMaterialOrange,
			Accent4: "FFC107", Accent5: "FF5722", Accent6: ColorMaterialPink,
			Hlink: "3F51B5", FolHlink: "1A237E",
		},
		Fonts: FontScheme{
			Name:      themeNameVibrant,
			MajorFont: "Outfit",
			MinorFont: "Outfit",
		},
		Primary: ColorMaterialPink, Secondary: "9C27B0", Accent: ColorMaterialOrange,
		Background: ColorWhite, Text: "212121", Light: "FCE4EC", Dark: "880E4F",
	}

	// ThemeDark - Easy on the eyes.
	ThemeDark = Theme{
		Name: themeNameDark,
		Colors: ColorScheme{
			Name: themeNameDark,
			Dk1:  ColorWhite, Lt1: "121212", Dk2: colorDarkPurple, Lt2: "1E1E1E",
			Accent1: colorDarkTeal, Accent2: colorDarkPink, Accent3: ColorMaterialOrange,
			Accent4: colorDarkTeal, Accent5: colorDarkPurple, Accent6: colorDarkPink,
			Hlink: colorDarkTeal, FolHlink: colorDarkPurple,
		},
		Fonts: FontScheme{
			Name:      themeNameDark,
			MajorFont: fontRoboto,
			MinorFont: fontRoboto,
		},
		Primary: colorDarkPurple, Secondary: colorDarkTeal, Accent: colorDarkPink,
		Background: "121212", Text: ColorWhite, Light: "1E1E1E", Dark: "000000",
	}

	// ThemeNature - Fresh and organic.
	ThemeNature = Theme{
		Name: themeNameNature,
		Colors: ColorScheme{
			Name: themeNameNature,
			Dk1:  colorNatureGreen900, Lt1: ColorWhite, Dk2: ColorCorporateGreen, Lt2: "E8F5E9",
			Accent1: "4CAF50", Accent2: "8BC34A", Accent3: "CDDC39",
			Accent4: "FFEB3B", Accent5: "FFC107", Accent6: ColorMaterialOrange,
			Hlink: ColorCorporateGreen, FolHlink: colorNatureGreen900,
		},
		Fonts: FontScheme{
			Name:      themeNameNature,
			MajorFont: fontInter,
			MinorFont: fontInter,
		},
		Primary: ColorCorporateGreen, Secondary: "4CAF50", Accent: "8BC34A",
		Background: ColorWhite, Text: colorNatureGreen900, Light: "E8F5E9", Dark: colorNatureGreen900,
	}

	// ThemeTech - Modern technology feel.
	ThemeTech = Theme{
		Name: themeNameTech,
		Colors: ColorScheme{
			Name: themeNameTech,
			Dk1:  "01579B", Lt1: colorNeutral50, Dk2: colorTechBlue900, Lt2: colorTechBlue50,
			Accent1: colorTechBlue700, Accent2: "00E676", Accent3: "00B0FF",
			Accent4: "00B8D4", Accent5: "0091EA", Accent6: colorTechBlue900,
			Hlink: "00B0FF", FolHlink: colorTechBlue900,
		},
		Fonts: FontScheme{
			Name:      themeNameTech,
			MajorFont: fontRoboto,
			MinorFont: fontRoboto,
		},
		Primary: colorTechBlue900, Secondary: colorTechBlue700, Accent: "00E676",
		Background: colorNeutral50, Text: "263238", Light: colorTechBlue50, Dark: "01579B",
	}

	// ThemeCarbon - IBM's design system.
	ThemeCarbon = Theme{
		Name: themeNameCarbon,
		Colors: ColorScheme{
			Name: themeNameCarbon,
			Dk1:  "161616", Lt1: ColorWhite, Dk2: ColorCarbonBlue60, Lt2: "E0E0E0",
			Accent1: ColorCarbonBlue40, Accent2: "24A148", Accent3: "08BDBA",
			Accent4: "1192E8", Accent5: "6929C4", Accent6: ColorCarbonBlue60,
			Hlink: ColorCarbonBlue40, FolHlink: ColorCarbonBlue60,
		},
		Fonts: FontScheme{
			Name:      themeNameCarbon,
			MajorFont: "IBM Plex Sans",
			MinorFont: "IBM Plex Sans",
		},
		Primary: ColorCarbonBlue60, Secondary: ColorCarbonBlue40, Accent: "24A148",
		Background: ColorWhite, Text: "161616", Light: "E0E0E0", Dark: "161616",
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
