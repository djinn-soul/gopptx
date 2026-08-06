package pptxxml

import "strings"

// PowerPoint resolves a SmartArt style by category and type together. A style ID
// paired with the wrong category is not an error — the diagram simply falls back
// to the plain look, which is what a 3-D quick style used to do when it was
// written next to the template's hardcoded qsCatId="simple".
//
// The categories below are the ones PowerPoint writes for each style family.

const (
	smartArtQuickStyleCatSimple = "simple"
	smartArtQuickStyleCat3D     = "3D"
	smartArtColorCatColorful    = "colorful"
	smartArtColorCatMainScheme  = "mainScheme"
	smartArtColorCatAccent1     = "accent1"
)

// smartArtQuickStyleCategory returns the category for a quick style URI, e.g.
// ".../quickstyle/3d1" -> "3D".
func smartArtQuickStyleCategory(quickStyleID string) string {
	switch name := smartArtStyleName(quickStyleID); {
	case strings.HasPrefix(name, "3d"):
		return smartArtQuickStyleCat3D
	default:
		return smartArtQuickStyleCatSimple
	}
}

// smartArtColorCategory returns the category for a colour style URI, e.g.
// ".../colors/accent3_2" -> "accent3", ".../colors/colorful1" -> "colorful".
func smartArtColorCategory(colorStyleID string) string {
	name := smartArtStyleName(colorStyleID)
	switch {
	case strings.HasPrefix(name, "colorful"):
		return smartArtColorCatColorful
	case strings.HasPrefix(name, "accent"):
		accent, _, _ := strings.Cut(strings.TrimPrefix(name, "accent"), "_")
		if accent == "0" {
			// accent0 is the theme's own scheme rather than a numbered accent.
			return smartArtColorCatMainScheme
		}
		if len(accent) == 1 && accent[0] >= '1' && accent[0] <= '6' {
			return "accent" + accent
		}
		return smartArtColorCatAccent1
	default:
		return smartArtColorCatAccent1
	}
}

// smartArtStyleName is the last segment of a style URI.
func smartArtStyleName(styleID string) string {
	return styleID[strings.LastIndex(styleID, "/")+1:]
}
