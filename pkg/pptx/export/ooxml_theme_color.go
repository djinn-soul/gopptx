//nolint:mnd // OOXML color math uses spec-defined numeric ranges (e.g., 0..100000, 0..255).
package export

import (
	"encoding/hex"
	"strconv"
	"strings"
)

const themeColorAccent1 = "accent1"

type rgbColor struct {
	r uint8
	g uint8
	b uint8
}

type colorTransforms struct {
	tint   int
	shade  int
	lumMod int
	lumOff int
}

func resolveOOXMLColorToken(token string) (uint8, uint8, uint8, bool) {
	baseToken, t := parseColorToken(token)
	if c, ok := parseHexRGB(baseToken); ok {
		c = applyColorTransforms(c, t)
		return c.r, c.g, c.b, true
	}
	base := normalizeColorName(baseToken)
	base = resolveColorAlias(base)
	c, ok := resolveThemeBaseColor(base)
	if !ok {
		return 0, 0, 0, false
	}
	c = applyColorTransforms(c, t)
	return c.r, c.g, c.b, true
}

func resolveColorAlias(base string) string {
	switch base {
	case "tx1":
		return themeSlotDk1
	case "bg1":
		return themeSlotLt1
	case "tx2":
		return themeSlotDk2
	case "bg2":
		return themeSlotLt2
	default:
		return base
	}
}

func resolveThemeBaseColor(name string) (rgbColor, bool) {
	switch name {
	case "dk1":
		return rgbColor{r: 0x00, g: 0x00, b: 0x00}, true
	case "lt1":
		return rgbColor{r: 0xFF, g: 0xFF, b: 0xFF}, true
	case "dk2":
		return rgbColor{r: 0x1F, g: 0x49, b: 0x7D}, true
	case "lt2":
		return rgbColor{r: 0xEE, g: 0xEC, b: 0xE1}, true
	case themeColorAccent1:
		return rgbColor{r: 0x4F, g: 0x81, b: 0xBD}, true
	case "accent2":
		return rgbColor{r: 0xC0, g: 0x50, b: 0x4D}, true
	case "accent3":
		return rgbColor{r: 0x9B, g: 0xBB, b: 0x59}, true
	case "accent4":
		return rgbColor{r: 0x80, g: 0x64, b: 0xA2}, true
	case "accent5":
		return rgbColor{r: 0x4B, g: 0xAC, b: 0xC6}, true
	case "accent6":
		return rgbColor{r: 0xF7, g: 0x96, b: 0x46}, true
	case "hlink":
		return rgbColor{r: 0x00, g: 0x00, b: 0xFF}, true
	case "folhlink":
		return rgbColor{r: 0x80, g: 0x00, b: 0x80}, true
	default:
		return rgbColor{}, false
	}
}

func parseColorToken(token string) (string, colorTransforms) {
	parts := strings.Split(strings.TrimSpace(token), "|")
	base := strings.TrimSpace(parts[0])
	t := colorTransforms{}
	for _, p := range parts[1:] {
		k, v, ok := strings.Cut(strings.TrimSpace(p), "=")
		if !ok {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(k)) {
		case "tint":
			t.tint = clampPct(n)
		case "shade":
			t.shade = clampPct(n)
		case "lummod":
			t.lumMod = clampPct(n)
		case "lumoff":
			t.lumOff = clampPct(n)
		}
	}
	return base, t
}

func parseHexRGB(raw string) (rgbColor, bool) {
	color := stripHash(strings.TrimSpace(raw))
	if len(color) != 6 {
		return rgbColor{}, false
	}
	b, err := hex.DecodeString(color)
	if err != nil || len(b) < 3 {
		return rgbColor{}, false
	}
	return rgbColor{r: b[0], g: b[1], b: b[2]}, true
}

func applyColorTransforms(c rgbColor, t colorTransforms) rgbColor {
	out := c
	if t.shade > 0 {
		out = shadeColor(out, t.shade)
	}
	if t.tint > 0 {
		out = tintColor(out, t.tint)
	}
	if t.lumMod > 0 || t.lumOff > 0 {
		out = lumModOffColor(out, t.lumMod, t.lumOff)
	}
	return out
}

func shadeColor(c rgbColor, shade int) rgbColor {
	return rgbColor{
		r: scaleColor(c.r, shade),
		g: scaleColor(c.g, shade),
		b: scaleColor(c.b, shade),
	}
}

func tintColor(c rgbColor, tint int) rgbColor {
	return rgbColor{
		r: tintComponent(c.r, tint),
		g: tintComponent(c.g, tint),
		b: tintComponent(c.b, tint),
	}
}

func lumModOffColor(c rgbColor, lumMod int, lumOff int) rgbColor {
	if lumMod == 0 {
		lumMod = 100000
	}
	return rgbColor{
		r: lumComponent(c.r, lumMod, lumOff),
		g: lumComponent(c.g, lumMod, lumOff),
		b: lumComponent(c.b, lumMod, lumOff),
	}
}

func scaleColor(v uint8, pct int) uint8 {
	out := (int(v) * pct) / 100000
	return clampByteToUint8(out)
}

func tintComponent(v uint8, pct int) uint8 {
	out := int(v) + ((255-int(v))*pct)/100000
	return clampByteToUint8(out)
}

func lumComponent(v uint8, lumMod int, lumOff int) uint8 {
	out := (int(v)*lumMod)/100000 + (255*lumOff)/100000
	return clampByteToUint8(out)
}

func clampPct(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100000 {
		return 100000
	}
	return v
}

func clampByteToUint8(v int) uint8 {
	switch {
	case v <= 0:
		return 0
	case v >= 255:
		return 255
	default:
		return uint8(v)
	}
}

func normalizeColorName(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if after, ok := strings.CutPrefix(s, "scheme:"); ok {
		return strings.TrimSpace(after)
	}
	return s
}
