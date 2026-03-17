package editor

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// SlideBackground describes a slide background to apply.
type SlideBackground struct {
	Type      string   // "solid" | "gradient" | "image" | "theme"
	Color     string   // hex RGB for solid, e.g. "FF0000"
	Colors    []string // hex RGB list for gradient
	Angle     int      // degrees for gradient
	ImagePath string   // file path for image background
	ImageData string   // base64 PNG/JPEG for image background
	ColorRef  string   // e.g. "accent1" for theme
}

// SetSlideBackground applies a background to the given slide.
func (e *PresentationEditor) SetSlideBackground(slideIndex int, bg SlideBackground) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]
	slideXML, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return fmt.Errorf("slide part %q not found", slideRef.Part)
	}
	bgXML, err := buildSlideBackgroundXML(e, slideRef.Part, bg)
	if err != nil {
		return fmt.Errorf("build background XML: %w", err)
	}
	e.parts.Set(slideRef.Part, []byte(injectSlideBg(string(slideXML), bgXML)))
	return nil
}

// buildSlideBackgroundXML creates the <p:bg> XML snippet for the given background.
func buildSlideBackgroundXML(e *PresentationEditor, slidePart string, bg SlideBackground) (string, error) {
	switch bg.Type {
	case "solid":
		return fmt.Sprintf(
			`<p:bg><p:bgPr><a:solidFill><a:srgbClr val="%s"/></a:solidFill>`+
				`<a:effectLst/></p:bgPr></p:bg>`, bg.Color), nil
	case "gradient":
		stopsXML := buildGradientStopsXML(bg.Colors)
		return fmt.Sprintf(
			`<p:bg><p:bgPr><a:gradFill><a:gsLst>%s</a:gsLst>`+
				`<a:lin ang="%d" scaled="0"/></a:gradFill>`+
				`<a:effectLst/></p:bgPr></p:bg>`,
			stopsXML, bg.Angle*60000), nil //nolint:mnd // EMU angle = degrees × 60000
	case "image":
		relID, imgErr := addBgImageRelationship(e, slidePart, bg)
		if imgErr != nil {
			return "", imgErr
		}
		return fmt.Sprintf(
			`<p:bg><p:bgPr><a:blipFill dpi="0" rotWithShape="1">`+
				`<a:blip r:embed="%s"/><a:stretch><a:fillRect/></a:stretch>`+
				`</a:blipFill><a:effectLst/></p:bgPr></p:bg>`, relID), nil
	case "theme":
		return fmt.Sprintf(
			`<p:bg><p:bgPr><a:solidFill><a:schemeClr val="%s"/></a:solidFill>`+
				`<a:effectLst/></p:bgPr></p:bg>`, bg.ColorRef), nil
	default:
		return "", fmt.Errorf("unknown background type %q", bg.Type)
	}
}

// buildGradientStopsXML builds <a:gs pos="…"> stops for gradient backgrounds.
func buildGradientStopsXML(colors []string) string {
	if len(colors) == 0 {
		return ""
	}
	const fullPct = 100000 // 100% in OOXML units
	var b strings.Builder
	step := 0
	if len(colors) > 1 {
		step = fullPct / (len(colors) - 1)
	}
	for i, c := range colors {
		pos := i * step
		if i == len(colors)-1 {
			pos = fullPct
		}
		fmt.Fprintf(&b, `<a:gs pos="%d"><a:srgbClr val="%s"/></a:gs>`, pos, c)
	}
	return b.String()
}

// addBgImageRelationship adds the image as a media part and returns the rel-ID.
func addBgImageRelationship(e *PresentationEditor, slidePart string, bg SlideBackground) (string, error) {
	imgData, ext, err := loadBgImageData(bg)
	if err != nil {
		return "", err
	}
	mediaPath, regErr := e.RegisterMedia(imgData, ext)
	if regErr != nil {
		return "", fmt.Errorf("add media part: %w", regErr)
	}
	relTarget := "../" + strings.TrimPrefix(mediaPath, "ppt/")
	relID, relErr := e.nextSlideRelID(slidePart)
	if relErr != nil {
		return "", fmt.Errorf("alloc rel id: %w", relErr)
	}
	const imgRelType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"
	rels := []common.EditorRelationship{{ID: relID, Type: imgRelType, Target: relTarget}}
	if addErr := e.addSlideRelationships(slidePart, rels); addErr != nil {
		return "", fmt.Errorf("add image rel: %w", addErr)
	}
	return relID, nil
}

// loadBgImageData reads image bytes from a path or base64 string.
func loadBgImageData(bg SlideBackground) ([]byte, string, error) {
	switch {
	case bg.ImagePath != "":
		data, err := os.ReadFile(bg.ImagePath)
		if err != nil {
			return nil, "", fmt.Errorf("read background image: %w", err)
		}
		ext := formatPNG
		lp := strings.ToLower(bg.ImagePath)
		if strings.HasSuffix(lp, ".jpg") || strings.HasSuffix(lp, ".jpeg") {
			ext = "jpeg"
		}
		return data, ext, nil
	case bg.ImageData != "":
		data, err := base64.StdEncoding.DecodeString(bg.ImageData)
		if err != nil {
			return nil, "", fmt.Errorf("decode base64 image data: %w", err)
		}
		return data, formatPNG, nil
	default:
		return nil, "", errors.New("image background requires image_path or image_data")
	}
}

// injectSlideBg removes any existing <p:bg> and inserts the new one after <p:cSld>.
func injectSlideBg(slideXML, bgXML string) string {
	reBg := regexp.MustCompile(`<p:bg>.*?</p:bg>`)
	slideXML = reBg.ReplaceAllString(slideXML, "")
	reCsld := regexp.MustCompile(`(<p:cSld[^>]*>)`)
	return reCsld.ReplaceAllString(slideXML, "$1"+bgXML)
}
