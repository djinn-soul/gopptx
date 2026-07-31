package editor

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type styleColor struct {
	rgb        string
	scheme     string
	transforms []colorTransform
}

type colorTransform struct {
	kind  string
	value float64
}

const (
	colorTransformAlpha = "alpha"
	colorTransformScale = 100000
	rgbHexLength        = 6
	rgbChannelMask      = 0xff
	rgbChannelMax       = 255
	rgbRedShift         = 16
	rgbGreenShift       = 8
	hueSectorCount      = 6
	blueHueOffset       = 4
	greenBlueBoundary   = 3
	blueRedBoundary     = 5
)

var (
	colorTransformPattern = regexp.MustCompile(
		`<a:(lumMod|lumOff|tint|shade|` + colorTransformAlpha +
			`)\b[^>]*\bval="(\d+)"[^>]*/?>`,
	)
	firstRunPropertiesPattern = regexp.MustCompile(
		`(?s)<a:rPr\b[^>]*>.*?</a:rPr>|<a:rPr\b[^>]*/>`,
	)
)

func parseColorTransforms(colorXML string) []colorTransform {
	matches := colorTransformPattern.FindAllStringSubmatch(colorXML, -1)
	out := make([]colorTransform, 0, len(matches))
	for _, match := range matches {
		raw, err := strconv.Atoi(match[2])
		if err != nil {
			continue
		}
		out = append(out, colorTransform{
			kind:  match[1],
			value: float64(raw) / colorTransformScale,
		})
	}
	return out
}

func effectiveColor(
	color styleColor,
	source string,
	theme themeStyleContext,
) *common.EffectiveColor {
	rgb := color.rgb
	if rgb == "" && color.scheme != "" {
		rgb, _ = theme.resolveSchemeSlot(color.scheme)
	}
	if rgb != "" {
		rgb = applyColorTransforms(rgb, color.transforms)
	}
	return &common.EffectiveColor{
		RGB:        rgb,
		SchemeSlot: color.scheme,
		Source:     source,
	}
}

func applyColorTransforms(rgb string, transforms []colorTransform) string {
	red, green, blue, ok := parseRGB(rgb)
	if !ok {
		return rgb
	}
	for _, transform := range transforms {
		switch transform.kind {
		case "lumMod", "lumOff":
			hue, saturation, lightness := rgbToHSL(red, green, blue)
			if transform.kind == "lumMod" {
				lightness *= transform.value
			} else {
				lightness += transform.value
			}
			red, green, blue = hslToRGB(hue, saturation, clampUnit(lightness))
		case "tint":
			red = 1 - ((1 - red) * transform.value)
			green = 1 - ((1 - green) * transform.value)
			blue = 1 - ((1 - blue) * transform.value)
		case "shade":
			red *= transform.value
			green *= transform.value
			blue *= transform.value
		case colorTransformAlpha:
			// Alpha changes opacity, not the concrete RGB requested here.
		}
	}
	return formatRGB(red, green, blue)
}

func parseRGB(rgb string) (float64, float64, float64, bool) {
	if len(rgb) != rgbHexLength {
		return 0, 0, 0, false
	}
	value, err := strconv.ParseUint(rgb, 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	return float64((value>>rgbRedShift)&rgbChannelMask) / rgbChannelMax,
		float64((value>>rgbGreenShift)&rgbChannelMask) / rgbChannelMax,
		float64(value&rgbChannelMask) / rgbChannelMax,
		true
}

func formatRGB(red, green, blue float64) string {
	channel := func(value float64) int {
		return int(math.Round(clampUnit(value) * rgbChannelMax))
	}
	return strings.ToUpper(fmt.Sprintf("%02X%02X%02X", channel(red), channel(green), channel(blue)))
}

func rgbToHSL(red, green, blue float64) (float64, float64, float64) {
	maximum := math.Max(red, math.Max(green, blue))
	minimum := math.Min(red, math.Min(green, blue))
	lightness := (maximum + minimum) / 2
	if maximum == minimum {
		return 0, 0, lightness
	}
	delta := maximum - minimum
	saturation := delta / (1 - math.Abs((2*lightness)-1))
	var hue float64
	switch maximum {
	case red:
		hue = math.Mod((green-blue)/delta, hueSectorCount)
	case green:
		hue = ((blue - red) / delta) + 2
	default:
		hue = ((red - green) / delta) + blueHueOffset
	}
	hue /= hueSectorCount
	if hue < 0 {
		hue++
	}
	return hue, saturation, lightness
}

func hslToRGB(hue, saturation, lightness float64) (float64, float64, float64) {
	chroma := (1 - math.Abs((2*lightness)-1)) * saturation
	scaledHue := hue * hueSectorCount
	secondary := chroma * (1 - math.Abs(math.Mod(scaledHue, 2)-1))
	var red, green, blue float64
	switch {
	case scaledHue < 1:
		red, green = chroma, secondary
	case scaledHue < 2:
		red, green = secondary, chroma
	case scaledHue < greenBlueBoundary:
		green, blue = chroma, secondary
	case scaledHue < blueHueOffset:
		green, blue = secondary, chroma
	case scaledHue < blueRedBoundary:
		red, blue = secondary, chroma
	default:
		red, blue = chroma, secondary
	}
	offset := lightness - (chroma / 2)
	return red + offset, green + offset, blue + offset
}

func clampUnit(value float64) float64 {
	return math.Max(0, math.Min(1, value))
}
