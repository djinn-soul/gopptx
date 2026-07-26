package pptxxml

import "strings"

// NormalizePlaceholderType returns the canonical OOXML placeholder type.
func NormalizePlaceholderType(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return "obj"
	}
	switch raw {
	case placeholderPicture, "pic":
		return "pic"
	case "title":
		return "title"
	case placeholderBody:
		return placeholderBody
	case "ctrtitle", "centeredtitle", "centered_title":
		return "ctrTitle"
	default:
		return raw
	}
}
