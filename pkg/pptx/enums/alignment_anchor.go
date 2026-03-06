package enums

import (
	"fmt"
	"strings"
)

type PPAlign string

const (
	PPAlignLeft         PPAlign = "l"
	PPAlignCenter       PPAlign = "ctr"
	PPAlignRight        PPAlign = "r"
	PPAlignJustify      PPAlign = "just"
	PPAlignDistribute   PPAlign = "dist"
	PPAlignThaiDist     PPAlign = "thaiDist"
	PPAlignJustifyLow   PPAlign = "justLow"
	PPAlignCenterAcross PPAlign = "ctr"
)

func (a PPAlign) XMLValue() string {
	return string(a)
}

func ParsePPAlign(value string) (PPAlign, error) {
	switch normalizeKey(value) {
	case "l", "left":
		return PPAlignLeft, nil
	case "ctr", "center", "centre":
		return PPAlignCenter, nil
	case "r", "right":
		return PPAlignRight, nil
	case "just", "justify":
		return PPAlignJustify, nil
	case "dist", "distribute":
		return PPAlignDistribute, nil
	case "thaidist", "thai_distribute":
		return PPAlignThaiDist, nil
	case "justlow", "justify_low":
		return PPAlignJustifyLow, nil
	case "justifylow":
		return PPAlignJustifyLow, nil
	default:
		return "", fmt.Errorf("invalid PP_ALIGN value %q", value)
	}
}

type MSOAnchor string
type MSOVerticalAnchor string

const (
	MSOAnchorTop        MSOAnchor = "t"
	MSOAnchorMiddle     MSOAnchor = "ctr"
	MSOAnchorBottom     MSOAnchor = "b"
	MSOAnchorJustify    MSOAnchor = "just"
	MSOAnchorDistribute MSOAnchor = "dist"
)

const (
	MSOVerticalAnchorTop        MSOVerticalAnchor = "t"
	MSOVerticalAnchorMiddle     MSOVerticalAnchor = "ctr"
	MSOVerticalAnchorBottom     MSOVerticalAnchor = "b"
	MSOVerticalAnchorJustify    MSOVerticalAnchor = "just"
	MSOVerticalAnchorDistribute MSOVerticalAnchor = "dist"
)

func (a MSOAnchor) XMLValue() string {
	return string(a)
}

func (a MSOVerticalAnchor) XMLValue() string {
	return string(a)
}

func ParseMSOAnchor(value string) (MSOAnchor, error) {
	switch normalizeKey(value) {
	case "t", "top":
		return MSOAnchorTop, nil
	case "ctr", "center", "middle":
		return MSOAnchorMiddle, nil
	case "b", "bottom":
		return MSOAnchorBottom, nil
	case "just", "justify":
		return MSOAnchorJustify, nil
	case "dist", "distribute":
		return MSOAnchorDistribute, nil
	default:
		return "", fmt.Errorf("invalid MSO_ANCHOR value %q", value)
	}
}

func ParseMSOVerticalAnchor(value string) (MSOVerticalAnchor, error) {
	anchor, err := ParseMSOAnchor(value)
	if err != nil {
		return "", err
	}
	return MSOVerticalAnchor(anchor), nil
}

func normalizeKey(value string) string {
	key := strings.ToLower(strings.TrimSpace(value))
	key = strings.ReplaceAll(key, "-", "")
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ReplaceAll(key, " ", "")
	return key
}
