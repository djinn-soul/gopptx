package editor

import (
	"errors"
	"regexp"
	"strings"
)

var (
	themeClrSchemeBlock = regexp.MustCompile(`(?is)<(?:\w+:)?clrScheme\b[^>]*>.*?</(?:\w+:)?clrScheme>`)
	themeSrgbValue      = regexp.MustCompile(`(?i)\bval="([0-9A-Fa-f]{6})"`)
	themeSysLastValue   = regexp.MustCompile(`(?i)\blastClr="([0-9A-Fa-f]{6})"`)
)

// GetThemeColorScheme reads the 12 standard color slots out of the first theme
// part, so a caller can resolve a scheme colour reference to a concrete RGB.
//
// Slots defined with a:sysClr resolve to their lastClr, which is the value
// PowerPoint itself falls back to when the system colour is unavailable.
func (e *PresentationEditor) GetThemeColorScheme() (ThemeColorScheme, error) {
	if e == nil || e.parts == nil {
		return ThemeColorScheme{}, errors.New("editor cannot be nil")
	}
	inv, err := e.GetThemeInventory()
	if err != nil {
		return ThemeColorScheme{}, err
	}
	if len(inv.ThemeParts) == 0 {
		return ThemeColorScheme{}, errors.New("presentation has no theme part")
	}
	data, ok := e.parts.Get(inv.ThemeParts[0])
	if !ok {
		return ThemeColorScheme{}, errors.New("theme part not found: " + inv.ThemeParts[0])
	}
	block := themeClrSchemeBlock.FindString(string(data))
	if block == "" {
		return ThemeColorScheme{}, errors.New("theme part has no colour scheme")
	}

	scheme := ThemeColorScheme{}
	for _, slot := range []struct {
		name  string
		field *string
	}{
		{themeSlotDk1, &scheme.Dk1}, {themeSlotLt1, &scheme.Lt1},
		{themeSlotDk2, &scheme.Dk2}, {themeSlotLt2, &scheme.Lt2},
		{themeSlotAccent1, &scheme.Accent1}, {themeSlotAccent2, &scheme.Accent2},
		{themeSlotAccent3, &scheme.Accent3}, {themeSlotAccent4, &scheme.Accent4},
		{themeSlotAccent5, &scheme.Accent5}, {themeSlotAccent6, &scheme.Accent6},
		{themeSlotHlink, &scheme.Hlink}, {themeSlotFolHlink, &scheme.FolHlink},
	} {
		*slot.field = readThemeColorSlot(block, slot.name)
	}
	return scheme, nil
}

func readThemeColorSlot(block string, slotName string) string {
	pattern := regexp.MustCompile(
		`(?is)<(?:\w+:)?` + slotName + `\b[^>]*>.*?</(?:\w+:)?` + slotName + `>`,
	)
	slot := pattern.FindString(block)
	if slot == "" {
		return ""
	}
	if match := themeSrgbValue.FindStringSubmatch(slot); len(match) > 1 {
		return strings.ToUpper(match[1])
	}
	if match := themeSysLastValue.FindStringSubmatch(slot); len(match) > 1 {
		return strings.ToUpper(match[1])
	}
	return ""
}
