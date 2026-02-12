package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

type (
	// Theme represents a color palette for a presentation.
	Theme = styling.Theme
)

// Font size presets in points.
const (
	FontSizeTitle    = styling.FontSizeTitle
	FontSizeSubtitle = styling.FontSizeSubtitle
	FontSizeHeading  = styling.FontSizeHeading
	FontSizeBody     = styling.FontSizeBody
	FontSizeSmall    = styling.FontSizeSmall
	FontSizeCaption  = styling.FontSizeCaption
	FontSizeCode     = styling.FontSizeCode
	FontSizeLarge    = styling.FontSizeLarge
	FontSizeXLarge   = styling.FontSizeXLarge
)

// Color constants for convenience.
const (
	ColorRed       = styling.ColorRed
	ColorGreen     = styling.ColorGreen
	ColorBlue      = styling.ColorBlue
	ColorWhite     = styling.ColorWhite
	ColorBlack     = styling.ColorBlack
	ColorGray      = styling.ColorGray
	ColorLightGray = styling.ColorLightGray
	ColorDarkGray  = styling.ColorDarkGray
	ColorYellow    = styling.ColorYellow
	ColorLightBlue = styling.ColorLightBlue
	ColorOrange    = styling.ColorOrange
	ColorPurple    = styling.ColorPurple
	ColorCyan      = styling.ColorCyan
	ColorMagenta   = styling.ColorMagenta
	ColorNavy      = styling.ColorNavy
	ColorTeal      = styling.ColorTeal
	ColorOlive     = styling.ColorOlive

	ColorCorporateBlue   = styling.ColorCorporateBlue
	ColorCorporateGreen  = styling.ColorCorporateGreen
	ColorCorporateRed    = styling.ColorCorporateRed
	ColorCorporateOrange = styling.ColorCorporateOrange

	ColorMaterialRed    = styling.ColorMaterialRed
	ColorMaterialPink   = styling.ColorMaterialPink
	ColorMaterialPurple = styling.ColorMaterialPurple
	ColorMaterialIndigo = styling.ColorMaterialIndigo
	ColorMaterialBlue   = styling.ColorMaterialBlue
	ColorMaterialCyan   = styling.ColorMaterialCyan
	ColorMaterialTeal   = styling.ColorMaterialTeal
	ColorMaterialGreen  = styling.ColorMaterialGreen
	ColorMaterialLime   = styling.ColorMaterialLime
	ColorMaterialAmber  = styling.ColorMaterialAmber
	ColorMaterialOrange = styling.ColorMaterialOrange
	ColorMaterialBrown  = styling.ColorMaterialBrown
	ColorMaterialGray   = styling.ColorMaterialGray

	ColorCarbonBlue60   = styling.ColorCarbonBlue60
	ColorCarbonBlue40   = styling.ColorCarbonBlue40
	ColorCarbonGray100  = styling.ColorCarbonGray100
	ColorCarbonGray80   = styling.ColorCarbonGray80
	ColorCarbonGray20   = styling.ColorCarbonGray20
	ColorCarbonGreen50  = styling.ColorCarbonGreen50
	ColorCarbonRed60    = styling.ColorCarbonRed60
	ColorCarbonPurple60 = styling.ColorCarbonPurple60
)

// Shared constants for line dashing and styles.
const (
	LineDashSolid       = styling.LineDashSolid
	LineDashDash        = styling.LineDashDash
	LineDashDot         = styling.LineDashDot
	LineDashDashDot     = styling.LineDashDashDot
	LineDashDashDotDot  = styling.LineDashDashDotDot
	LineDashLongDash    = styling.LineDashLongDash
	LineDashLongDashDot = styling.LineDashLongDashDot
)

// Unit conversion helpers.
func Inches(v float64) Length      { return Length(styling.Inches(v)) }
func InchesToEMU(v float64) Length { return Length(styling.InchesToEMU(v)) }
func Centimeters(v float64) Length { return Length(styling.Centimeters(v)) }
func CMToEMU(v float64) Length     { return Length(styling.CMToEMU(v)) }
func Points(v float64) Length      { return Length(styling.Points(v)) }
func PointsToEMU(v float64) Length { return Length(styling.PointsToEMU(v)) }
func Emu(v int64) Length           { return Length(styling.Emu(v)) }
func FontSize(v float64) int       { return styling.FontSize(v) }

// Theme presets.
var (
	ThemeCorporate = styling.ThemeCorporate
	ThemeModern    = styling.ThemeModern
	ThemeVibrant   = styling.ThemeVibrant
	ThemeDark      = styling.ThemeDark
	ThemeNature    = styling.ThemeNature
	ThemeTech      = styling.ThemeTech
	ThemeCarbon    = styling.ThemeCarbon
)

// AllThemes returns all available theme presets.
func AllThemes() []Theme {
	return styling.AllThemes()
}
