package export

import (
	"os"
	"runtime"
)

const (
	fontFamilySans           = "sans"
	fontFamilySansBold       = "sans-bold"
	fontFamilySansItalic     = "sans-italic"
	fontFamilySansBoldItalic = "sans-bolditalic"

	fontFamilySerif           = "serif"
	fontFamilySerifBold       = "serif-bold"
	fontFamilySerifItalic     = "serif-italic"
	fontFamilySerifBoldItalic = "serif-bolditalic"

	fontFamilyMono           = "mono"
	fontFamilyMonoBold       = "mono-bold"
	fontFamilyMonoItalic     = "mono-italic"
	fontFamilyMonoBoldItalic = "mono-bolditalic"

	fontFamilyCJK = "cjk"
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
	case fontFamilyCJK:
		return []string{
			fontsDir + "msyh.ttc",
			fontsDir + "msgothic.ttc",
			fontsDir + "malgun.ttf",
			fontsDir + "simsun.ttc",
		}
	case fontFamilyMono:
		return []string{
			fontsDir + "consola.ttf",
			fontsDir + "cour.ttf",
			fontsDir + "lucon.ttf",
		}
	case fontFamilyMonoBold:
		return []string{
			fontsDir + "consolab.ttf",
			fontsDir + "courbd.ttf",
		}
	case fontFamilyMonoItalic:
		return []string{
			fontsDir + "consolai.ttf",
			fontsDir + "couri.ttf",
		}
	case fontFamilyMonoBoldItalic:
		return []string{
			fontsDir + "consolaz.ttf",
			fontsDir + "courbi.ttf",
		}
	case fontFamilySerif:
		return []string{
			fontsDir + "cambria.ttf",
			fontsDir + "times.ttf",
			fontsDir + "georgia.ttf",
		}
	case fontFamilySerifBold:
		return []string{
			fontsDir + "cambriab.ttf",
			fontsDir + "timesbd.ttf",
			fontsDir + "georgiab.ttf",
		}
	case fontFamilySerifItalic:
		return []string{
			fontsDir + "cambriai.ttf",
			fontsDir + "timesi.ttf",
			fontsDir + "georgiai.ttf",
		}
	case fontFamilySerifBoldItalic:
		return []string{
			fontsDir + "cambriaz.ttf",
			fontsDir + "timesbi.ttf",
			fontsDir + "georgiaz.ttf",
		}
	case fontFamilySansBold:
		return []string{
			fontsDir + "calibrib.ttf",
			fontsDir + "arialbd.ttf",
			fontsDir + "segoeuib.ttf",
			fontsDir + "tahomabd.ttf",
		}
	case fontFamilySansItalic:
		return []string{
			fontsDir + "calibrii.ttf",
			fontsDir + "ariali.ttf",
			fontsDir + "segoeuii.ttf",
		}
	case fontFamilySansBoldItalic:
		return []string{
			fontsDir + "calibriz.ttf",
			fontsDir + "arialbi.ttf",
			fontsDir + "segoeuiz.ttf",
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
	case fontFamilyCJK:
		return []string{
			"/System/Library/Fonts/PingFang.ttc",
			"/System/Library/Fonts/Hiragino Sans GB.ttc",
			"/System/Library/Fonts/AppleSDGothicNeo.ttc",
		}
	case fontFamilyMono:
		return []string{
			"/System/Library/Fonts/SFNSMono.ttf",
			"/System/Library/Fonts/Menlo.ttc",
			"/Library/Fonts/Courier New.ttf",
		}
	case fontFamilyMonoBold, fontFamilyMonoItalic, fontFamilyMonoBoldItalic,
		fontFamilySerifBold, fontFamilySerifItalic, fontFamilySerifBoldItalic,
		fontFamilySansBold, fontFamilySansItalic, fontFamilySansBoldItalic:
		return nil // macOS: bold/italic handled by the same .ttc file; treated as regular
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
	case fontFamilyCJK:
		return []string{
			"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/opentype/noto/NotoSansCJKSC-Regular.otf",
		}
	case fontFamilyMono:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationMono-Regular.ttf",
			"/usr/share/fonts/truetype/freefont/FreeMono.ttf",
		}
	case fontFamilyMonoBold:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSansMono-Bold.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationMono-Bold.ttf",
		}
	case fontFamilyMonoItalic:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSansMono-Oblique.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationMono-Italic.ttf",
		}
	case fontFamilyMonoBoldItalic:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSansMono-BoldOblique.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationMono-BoldItalic.ttf",
		}
	case fontFamilySerif:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSerif.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationSerif-Regular.ttf",
			"/usr/share/fonts/truetype/freefont/FreeSerif.ttf",
		}
	case fontFamilySerifBold:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSerif-Bold.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationSerif-Bold.ttf",
		}
	case fontFamilySerifItalic:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSerif-Italic.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationSerif-Italic.ttf",
		}
	case fontFamilySerifBoldItalic:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSerif-BoldItalic.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationSerif-BoldItalic.ttf",
		}
	case fontFamilySansBold:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
			"/usr/share/fonts/noto/NotoSans-Bold.ttf",
		}
	case fontFamilySansItalic:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSans-Oblique.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationSans-Italic.ttf",
			"/usr/share/fonts/noto/NotoSans-Italic.ttf",
		}
	case fontFamilySansBoldItalic:
		return []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSans-BoldOblique.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationSans-BoldItalic.ttf",
			"/usr/share/fonts/noto/NotoSans-BoldItalic.ttf",
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
