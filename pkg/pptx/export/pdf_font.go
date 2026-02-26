package export

import (
	"os"
	"runtime"
)

// systemFontPaths returns candidate TTF paths for common sans-serif fonts
// on Windows, macOS, and Linux. Paths are tried in order; first hit wins.
func systemFontPaths() []string {
	var paths []string

	switch runtime.GOOS {
	case "windows":
		winDir := os.Getenv("WINDIR")
		if winDir == "" {
			winDir = `C:\Windows`
		}
		fontsDir := winDir + `\Fonts\`
		paths = []string{
			fontsDir + "calibri.ttf",
			fontsDir + "arial.ttf",
			fontsDir + "segoeui.ttf",
			fontsDir + "tahoma.ttf",
			fontsDir + "verdana.ttf",
		}
	case "darwin":
		paths = []string{
			"/System/Library/Fonts/Helvetica.ttc",
			"/System/Library/Fonts/SFPro.ttf",
			"/Library/Fonts/Arial.ttf",
			"/System/Library/Fonts/Supplemental/Arial.ttf",
		}
	default: // Linux
		paths = []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
			"/usr/share/fonts/TTF/DejaVuSans.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
			"/usr/share/fonts/truetype/freefont/FreeSans.ttf",
			"/usr/share/fonts/noto/NotoSans-Regular.ttf",
		}
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
