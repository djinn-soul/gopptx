package editor

func (e *PresentationEditor) themeStyleContext(masterPart string) themeStyleContext {
	ctx := themeStyleContext{}
	inv, err := e.GetThemeInventory()
	if err != nil || len(inv.ThemeParts) == 0 {
		return ctx
	}

	themePart := ""
	for _, binding := range inv.Bindings {
		if binding.OwnerPart == masterPart {
			themePart = binding.ThemePart
			break
		}
	}
	if themePart == "" {
		themePart = inv.ThemeParts[0]
	}
	if scheme, schemeErr := e.themeColorSchemeForPart(themePart); schemeErr == nil {
		ctx.scheme = scheme
	}
	data, ok := e.parts.Get(themePart)
	if !ok {
		return ctx
	}
	if match := majorLatinPattern.FindStringSubmatch(string(data)); len(match) > 1 {
		ctx.majorLatin = match[1]
	}
	if match := minorLatinPattern.FindStringSubmatch(string(data)); len(match) > 1 {
		ctx.minorLatin = match[1]
	}
	return ctx
}

// resolveTypeface expands the "+mj-lt"/"+mn-lt" references a placeholder uses
// to defer its font to the theme.
func (t themeStyleContext) resolveTypeface(value string) (string, bool) {
	switch value {
	case "+mj-lt", "+mj-ea", "+mj-cs":
		return t.majorLatin, t.majorLatin != ""
	case "+mn-lt", "+mn-ea", "+mn-cs":
		return t.minorLatin, t.minorLatin != ""
	default:
		return "", false
	}
}

// resolveSchemeSlot maps a scheme colour reference onto the theme's palette.
// tx1/bg1/tx2/bg2 are the slide-facing aliases of dk1/lt1/dk2/lt2.
func (t themeStyleContext) resolveSchemeSlot(slot string) (string, bool) {
	byName := map[string]string{
		themeSlotDk1: t.scheme.Dk1, "tx1": t.scheme.Dk1,
		themeSlotLt1: t.scheme.Lt1, "bg1": t.scheme.Lt1,
		themeSlotDk2: t.scheme.Dk2, "tx2": t.scheme.Dk2,
		themeSlotLt2: t.scheme.Lt2, "bg2": t.scheme.Lt2,
		themeSlotAccent1: t.scheme.Accent1, themeSlotAccent2: t.scheme.Accent2,
		themeSlotAccent3: t.scheme.Accent3, themeSlotAccent4: t.scheme.Accent4,
		themeSlotAccent5: t.scheme.Accent5, themeSlotAccent6: t.scheme.Accent6,
		themeSlotHlink: t.scheme.Hlink, themeSlotFolHlink: t.scheme.FolHlink,
	}
	rgb, ok := byName[slot]
	return rgb, ok && rgb != ""
}
