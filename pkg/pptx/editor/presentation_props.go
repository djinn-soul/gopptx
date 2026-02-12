package editor

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const themePartPath = "ppt/theme/theme1.xml"

var sldSzPattern = regexp.MustCompile(`(?s)<p:sldSz\b[^>]*/>`)

// ApplyTheme replaces the package theme with the provided palette and fonts.
func (e *PresentationEditor) ApplyTheme(theme styling.Theme) error {
	if e == nil {
		return fmt.Errorf("editor cannot be nil")
	}
	if _, err := requirePart(e.parts, themePartPath); err != nil {
		return err
	}
	e.parts[themePartPath] = []byte(pptxxml.Theme(mapEditorThemeToSpec(&theme)))
	return nil
}

// SetSlideSize updates presentation dimensions used for existing and future slides.
func (e *PresentationEditor) SetSlideSize(size common.SlideSize) error {
	if e == nil {
		return fmt.Errorf("editor cannot be nil")
	}
	if size.Width <= 0 || size.Height <= 0 {
		return fmt.Errorf("slide size must have positive dimensions")
	}

	rewritten, err := rewritePresentationSlideSize(e.presentationXML, size)
	if err != nil {
		return err
	}

	e.presentationXML = rewritten
	e.metadata.SlideSize = size
	return nil
}

func parsePresentationSlideSize(content []byte) (common.SlideSize, error) {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return common.SlideSize{}, nil
			}
			return common.SlideSize{}, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "sldSz" {
			continue
		}
		return parseSlideSizeAttrs(start.Attr)
	}
}

func parseSlideSizeAttrs(attrs []xml.Attr) (common.SlideSize, error) {
	size := common.SlideSize{}
	for _, attr := range attrs {
		switch attr.Name.Local {
		case "cx":
			width, err := strconv.ParseInt(strings.TrimSpace(attr.Value), 10, 64)
			if err != nil {
				return common.SlideSize{}, fmt.Errorf("invalid slide width %q", attr.Value)
			}
			size.Width = width
		case "cy":
			height, err := strconv.ParseInt(strings.TrimSpace(attr.Value), 10, 64)
			if err != nil {
				return common.SlideSize{}, fmt.Errorf("invalid slide height %q", attr.Value)
			}
			size.Height = height
		}
	}
	if size.Width == 0 || size.Height == 0 {
		return common.SlideSize{}, fmt.Errorf("slide size entry missing cx/cy")
	}
	return size, nil
}

func rewritePresentationSlideSize(current string, size common.SlideSize) (string, error) {
	if strings.TrimSpace(current) == "" {
		return "", fmt.Errorf("missing presentation XML content")
	}
	entry := fmt.Sprintf(`<p:sldSz cx="%d" cy="%d" type="%s"/>`, size.Width, size.Height, slideSizeType(size))

	if sldSzPattern.MatchString(current) {
		replaced := false
		return sldSzPattern.ReplaceAllStringFunc(current, func(match string) string {
			if replaced {
				return match
			}
			replaced = true
			return entry
		}), nil
	}

	if idx := strings.Index(current, "<p:notesSz"); idx >= 0 {
		return current[:idx] + entry + "\n" + current[idx:], nil
	}
	if idx := strings.Index(current, "</p:presentation>"); idx >= 0 {
		return current[:idx] + entry + "\n" + current[idx:], nil
	}
	return "", fmt.Errorf("presentation XML does not contain <p:notesSz> or </p:presentation>")
}

func slideSizeType(size common.SlideSize) string {
	if size.Width == common.SlideSize4x3.Width && size.Height == common.SlideSize4x3.Height {
		return "screen4x3"
	}
	if size.Width == common.SlideSize16x9.Width && size.Height == common.SlideSize16x9.Height {
		return "screen16x9"
	}
	return "custom"
}

func mapEditorThemeToSpec(theme *styling.Theme) *pptxxml.ThemeSpec {
	if theme == nil {
		return nil
	}
	return &pptxxml.ThemeSpec{
		Name: theme.Name,
		Colors: pptxxml.ColorSchemeSpec{
			Name:     theme.Colors.Name,
			Dk1:      theme.Colors.Dk1,
			Lt1:      theme.Colors.Lt1,
			Dk2:      theme.Colors.Dk2,
			Lt2:      theme.Colors.Lt2,
			Accent1:  theme.Colors.Accent1,
			Accent2:  theme.Colors.Accent2,
			Accent3:  theme.Colors.Accent3,
			Accent4:  theme.Colors.Accent4,
			Accent5:  theme.Colors.Accent5,
			Accent6:  theme.Colors.Accent6,
			Hlink:    theme.Colors.Hlink,
			FolHlink: theme.Colors.FolHlink,
		},
		Fonts: pptxxml.FontSchemeSpec{
			Name:      theme.Fonts.Name,
			MajorFont: theme.Fonts.MajorFont,
			MinorFont: theme.Fonts.MinorFont,
		},
	}
}
