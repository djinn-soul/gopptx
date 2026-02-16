package styling

import (
	"errors"
	"fmt"
	"regexp"
)

var hexColorRE = regexp.MustCompile(`^[0-9A-Fa-f]{6}$`)

// ValidateHexColor checks that a string is a 6-digit RGB hex color.
func ValidateHexColor(color, label string) error {
	if !hexColorRE.MatchString(color) {
		return fmt.Errorf("%s must be a 6-digit RGB hex color, got %q", label, color)
	}
	return nil
}

// ValidateColorScheme validates all 12 colors in a ColorScheme.
func ValidateColorScheme(cs ColorScheme) error {
	checks := []struct {
		val   string
		label string
	}{
		{cs.Dk1, "Dk1"},
		{cs.Lt1, "Lt1"},
		{cs.Dk2, "Dk2"},
		{cs.Lt2, "Lt2"},
		{cs.Accent1, "Accent1"},
		{cs.Accent2, "Accent2"},
		{cs.Accent3, "Accent3"},
		{cs.Accent4, "Accent4"},
		{cs.Accent5, "Accent5"},
		{cs.Accent6, "Accent6"},
		{cs.Hlink, "Hlink"},
		{cs.FolHlink, "FolHlink"},
	}
	for _, c := range checks {
		if c.val == "" {
			return fmt.Errorf("color scheme field %s is required", c.label)
		}
		if err := ValidateHexColor(c.val, c.label); err != nil {
			return err
		}
	}
	return nil
}

// ValidateTheme validates a Theme's required fields.
func ValidateTheme(t Theme) error {
	if t.Name == "" {
		return errors.New("theme name is required")
	}
	if err := ValidateColorScheme(t.Colors); err != nil {
		return fmt.Errorf("theme %q colors: %w", t.Name, err)
	}
	if t.Fonts.MajorFont == "" {
		return fmt.Errorf("theme %q: major font is required", t.Name)
	}
	if t.Fonts.MinorFont == "" {
		return fmt.Errorf("theme %q: minor font is required", t.Name)
	}
	return nil
}
