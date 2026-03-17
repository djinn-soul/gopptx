package export

import (
	"os"
	"runtime"
)

const (
	fontFamilySans  = "sans"
	fontFamilySerif = "serif"
	fontFamilyMono  = "mono"
)

func systemFontPathsForFamily(family string) []string {
	var paths []string

	switch runtime.GOOS {
	case "windows":
		winDir := os.Getenv("WINDIR")
		if winDir == "" {
			winDir = `C:\Windows`
		}
		fontsDir := winDir + `\Fonts\`
		paths = windowsFontCandidates(fontsDir, family)
	case "darwin":
		paths = macFontCandidates(family)
	default: // Linux
		paths = linuxFontCandidates(family)
	}

	// Filter to only paths that actually exist.
	existing := make([]string, 0, len(paths))
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			existing = append(existing, p)
		}
	}
	return existing
}

func windowsFontCandidates(fontsDir, family string) []string {
	switch family {
	case fontFamilyMono:
		return []string{
			fontsDir + "consola.ttf",
			fontsDir + "cour.ttf",
			fontsDir + "lucon.ttf",
		}
	case fontFamilySerif:
		return []string{
			fontsDir + "cambria.ttf",
			fontsDir + "times.ttf",
			fontsDir + "georgia.ttf",
		}
	default:
		return []string{
			fontsDir + "calibri.ttf",
			fontsDir + "arial.ttf",
			fontsDir + "segoeui.ttf",
			fontsDir + "tahoma.ttf",
			fontsDir + "verdana.ttf",
		}
	}
}

func macFontCandidates(family string) []string {
	switch family {
	case fontFamilyMono:
		return []string{
			"/System/Library/Fonts/SFNSMono.ttf",
			"/System/Library/Fonts/Menlo.ttc",
			"/Library/Fonts/Courier New.ttf",
		}
	case fontFamilySerif:
		return []string{
			"/System/Library/Fonts/Times.ttc",
			"/Library/Fonts/Times New Roman.ttf",
			"/System/Library/Fonts/Supplemental/Times New Roman.ttf",
		}
	default:
		return []string{
			"/System/Library/Fonts/Helvetica.ttc",
			"/System/Library/Fonts/SFPro.ttf",
			"/Library/Fonts/Arial.ttf",
			"/System/Library/Fonts/Supplemental/Arial.ttf",
		}
	}
}

func linuxFontCandidates(family string) []string {
	switch family {
	case fontFamilyMono:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationMono-Regular.ttf",
			"/usr/share/fonts/truetype/freefont/FreeMono.ttf",
		}
	case fontFamilySerif:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSerif.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationSerif-Regular.ttf",
			"/usr/share/fonts/truetype/freefont/FreeSerif.ttf",
		}
	default:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
			"/usr/share/fonts/TTF/DejaVuSans.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
			"/usr/share/fonts/truetype/freefont/FreeSans.ttf",
			"/usr/share/fonts/noto/NotoSans-Regular.ttf",
		}
	}
}
