package enums

import (
	"fmt"
)

type MSOThemeColor string

const (
	MSOThemeColorDark1       MSOThemeColor = "dk1"
	MSOThemeColorLight1      MSOThemeColor = "lt1"
	MSOThemeColorDark2       MSOThemeColor = "dk2"
	MSOThemeColorLight2      MSOThemeColor = "lt2"
	MSOThemeColorAccent1     MSOThemeColor = "accent1"
	MSOThemeColorAccent2     MSOThemeColor = "accent2"
	MSOThemeColorAccent3     MSOThemeColor = "accent3"
	MSOThemeColorAccent4     MSOThemeColor = "accent4"
	MSOThemeColorAccent5     MSOThemeColor = "accent5"
	MSOThemeColorAccent6     MSOThemeColor = "accent6"
	MSOThemeColorHyperlink   MSOThemeColor = "hlink"
	MSOThemeColorFollowedHL  MSOThemeColor = "folHlink"
	MSOThemeColorText1       MSOThemeColor = "tx1"
	MSOThemeColorText2       MSOThemeColor = "tx2"
	MSOThemeColorBackground1 MSOThemeColor = "bg1"
	MSOThemeColorBackground2 MSOThemeColor = "bg2"
)

func (c MSOThemeColor) XMLValue() string {
	return string(c)
}

func ParseMSOThemeColor(value string) (MSOThemeColor, error) {
	switch normalizeKey(value) {
	case "dk1", "dark1":
		return MSOThemeColorDark1, nil
	case "lt1", "light1":
		return MSOThemeColorLight1, nil
	case "dk2", "dark2":
		return MSOThemeColorDark2, nil
	case "lt2", "light2":
		return MSOThemeColorLight2, nil
	case string(MSOThemeColorAccent1):
		return MSOThemeColorAccent1, nil
	case "accent2":
		return MSOThemeColorAccent2, nil
	case "accent3":
		return MSOThemeColorAccent3, nil
	case "accent4":
		return MSOThemeColorAccent4, nil
	case "accent5":
		return MSOThemeColorAccent5, nil
	case "accent6":
		return MSOThemeColorAccent6, nil
	case "hlink", "hyperlink":
		return MSOThemeColorHyperlink, nil
	case "folhlink", "followedhyperlink":
		return MSOThemeColorFollowedHL, nil
	case "tx1", "text1":
		return MSOThemeColorText1, nil
	case "tx2", "text2":
		return MSOThemeColorText2, nil
	case "bg1", "background1":
		return MSOThemeColorBackground1, nil
	case "bg2", "background2":
		return MSOThemeColorBackground2, nil
	default:
		return "", fmt.Errorf("invalid MSO_THEME_COLOR value %q", value)
	}
}

type MSOColorType string

const (
	MSOColorTypeRGB     MSOColorType = "rgb"
	MSOColorTypeScheme  MSOColorType = "scheme"
	MSOColorTypeAuto    MSOColorType = "auto"
	MSOColorTypeUnknown MSOColorType = "unknown"
)

func (t MSOColorType) XMLValue() string {
	return string(t)
}

func ParseMSOColorType(value string) (MSOColorType, error) {
	switch normalizeKey(value) {
	case string(MSOColorTypeRGB), "srgb":
		return MSOColorTypeRGB, nil
	case "scheme", "theme":
		return MSOColorTypeScheme, nil
	case "auto":
		return MSOColorTypeAuto, nil
	case "unknown":
		return MSOColorTypeUnknown, nil
	default:
		return "", fmt.Errorf("invalid MSO_COLOR_TYPE value %q", value)
	}
}
