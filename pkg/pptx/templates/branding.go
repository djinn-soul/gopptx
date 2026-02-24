package templates

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	defaultCoverTitleSizePt  = 44
	emphasisCoverTitleSizePt = 52
)

// BrandingPreset defines visual styles for templates.
type BrandingPreset string

const (
	// PresetCorporate - Professional, blue/grey theme (default).
	PresetCorporate BrandingPreset = "corporate"
	// PresetModern - Clean, dark/white theme.
	PresetModern BrandingPreset = "modern"
	// PresetCreative - Bold, vibrant theme.
	PresetCreative BrandingPreset = "creative"
)

// BrandingSpec defines visual branding for a template.
type BrandingSpec struct {
	Preset BrandingPreset
	Theme  *styling.Theme
	Header string
	Footer string
}

// MapPreset returns the styling.Theme for a given preset.
func MapPreset(p BrandingPreset) *styling.Theme {
	return presetToTheme(p)
}

func presetToTheme(p BrandingPreset) *styling.Theme {
	switch p {
	case PresetModern:
		return &styling.ThemeModern
	case PresetCreative:
		return &styling.ThemeVibrant
	case PresetCorporate:
		fallthrough
	default:
		return &styling.ThemeCorporate
	}
}

// Apply applies branding settings to a slide.
func (b BrandingSpec) Apply(s elements.SlideContent) elements.SlideContent {
	return b.ApplyAt(s, 0)
}

// ApplyAt applies branding settings with slide index context for dynamic styling.
func (b BrandingSpec) ApplyAt(s elements.SlideContent, slideIndex int) elements.SlideContent {
	if theme := b.resolveTheme(); theme != nil {
		s = applyThemeVisuals(s, *theme, slideIndex)
	}
	if b.Footer != "" {
		s.FooterText = b.Footer
	}
	return s
}

func (b BrandingSpec) resolveTheme() *styling.Theme {
	if b.Theme != nil {
		return b.Theme
	}
	if b.Preset != "" {
		return presetToTheme(b.Preset)
	}
	// Default to corporate so templates always get baseline visual styling.
	return presetToTheme(PresetCorporate)
}

func applyThemeVisuals(s elements.SlideContent, theme styling.Theme, slideIndex int) elements.SlideContent {
	if slideIndex == 0 {
		s = applyCoverSlideVisuals(s, theme)
	} else {
		s = applyBodySlideVisuals(s, theme, slideIndex)
	}
	return s
}

func applyCoverSlideVisuals(s elements.SlideContent, theme styling.Theme) elements.SlideContent {
	s = s.WithBackgroundColor(theme.Primary)
	if strings.TrimSpace(s.TitleColor) == "" {
		s = s.WithTitleColor(theme.Colors.Lt1)
	}
	if strings.TrimSpace(s.ContentColor) == "" {
		s = s.WithContentColor(theme.Colors.Lt1)
	}
	if s.TitleSize == defaultCoverTitleSizePt {
		s = s.WithTitleSize(emphasisCoverTitleSizePt)
	}
	if strings.TrimSpace(s.TitleFont) == "" && strings.TrimSpace(theme.Fonts.MajorFont) != "" {
		s = s.WithTitleFont(theme.Fonts.MajorFont)
	}
	s.TitleBold = true
	return s
}

func applyBodySlideVisuals(s elements.SlideContent, theme styling.Theme, slideIndex int) elements.SlideContent {
	accent := themeAccent(theme, slideIndex)
	if slideIndex%2 == 0 {
		s = s.WithBackgroundColor(theme.Light)
	} else {
		s = s.WithBackgroundColor(theme.Background)
	}

	if strings.TrimSpace(s.TitleColor) == "" {
		s = s.WithTitleColor(accent)
	}
	if strings.TrimSpace(s.ContentColor) == "" {
		s = s.WithContentColor(theme.Text)
	}
	if strings.TrimSpace(s.TitleFont) == "" && strings.TrimSpace(theme.Fonts.MajorFont) != "" {
		s = s.WithTitleFont(theme.Fonts.MajorFont)
	}
	return s
}

func themeAccent(theme styling.Theme, slideIndex int) string {
	accents := []string{
		theme.Colors.Accent1,
		theme.Colors.Accent2,
		theme.Colors.Accent3,
		theme.Colors.Accent4,
		theme.Colors.Accent5,
		theme.Colors.Accent6,
	}
	return accents[slideIndex%len(accents)]
}
