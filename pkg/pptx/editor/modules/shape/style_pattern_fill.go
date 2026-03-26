package shape

import (
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func renderPatternFillXML(pattern *common.PatternedFill) (string, error) {
	if pattern == nil {
		return "", nil
	}
	prst := "pct5"
	if pattern.Preset != nil && strings.TrimSpace(*pattern.Preset) != "" {
		prst = strings.TrimSpace(*pattern.Preset)
	}
	fg := "000000"
	if pattern.FgColor != nil {
		color, err := NormalizeHexColor(*pattern.FgColor)
		if err != nil {
			return "", fmt.Errorf("fill.pattern.fg_color: %w", err)
		}
		fg = color
	}
	bg := "FFFFFF"
	if pattern.BgColor != nil {
		color, err := NormalizeHexColor(*pattern.BgColor)
		if err != nil {
			return "", fmt.Errorf("fill.pattern.bg_color: %w", err)
		}
		bg = color
	}
	return fmt.Sprintf(
		`<a:pattFill prst="%s"><a:fgClr><a:srgbClr val="%s"/></a:fgClr><a:bgClr><a:srgbClr val="%s"/></a:bgClr></a:pattFill>`,
		XMLEscape(prst),
		fg,
		bg,
	), nil
}
