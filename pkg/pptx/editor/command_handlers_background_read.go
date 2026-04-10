// Package editor provides the PresentationEditor for reading and modifying PPTX files.
package editor

import (
	"fmt"
	"math"
	"path"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

var (
	reBgSrgbClr   = regexp.MustCompile(`<a:srgbClr val="([0-9A-Fa-f]{6})"`)
	reBgGsStop    = regexp.MustCompile(`<a:gs pos="(\d+)"><a:srgbClr val="([0-9A-Fa-f]{6})"`)
	reBgLinAngle  = regexp.MustCompile(`<a:lin ang="(\d+)"`)
	reBgBlipEmbed = regexp.MustCompile(`r:embed="([^"]+)"`)
)

// GetSlideBackground reads the background of an existing slide.
// Returns nil, nil when the slide has no explicit background defined.
func (e *PresentationEditor) GetSlideBackground(slideIndex int) (*elements.SlideBackground, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]
	slideXML, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return nil, fmt.Errorf("slide part %q not found", slideRef.Part)
	}

	raw := string(slideXML)
	bgXML := extractBgXML(raw)
	return parseBgXML(bgXML, e, slideRef.Part), nil
}

func extractBgXML(raw string) string {
	start := strings.Index(raw, "<p:bg>")
	if start < 0 {
		return ""
	}
	end := strings.Index(raw[start:], "</p:bg>")
	if end < 0 {
		return ""
	}
	return raw[start : start+end+len("</p:bg>")]
}

func parseBgXML(bgXML string, e *PresentationEditor, slidePart string) *elements.SlideBackground {
	if strings.Contains(bgXML, "<a:solidFill>") {
		if color := extractSolidBgColor(bgXML); color != "" {
			bg := elements.NewSolidBackground(color)
			return &bg
		}
	}
	if strings.Contains(bgXML, "<a:gradFill>") {
		if grad := extractGradientBgFill(bgXML); grad != nil {
			bg := elements.NewGradientBackground(*grad)
			return &bg
		}
	}
	if strings.Contains(bgXML, "<a:blipFill>") {
		if img := extractPictureBgFill(bgXML, e, slidePart); img != nil {
			bg := elements.NewPictureBackground(*img)
			return &bg
		}
	}
	return nil
}

func extractSolidBgColor(bgXML string) string {
	m := reBgSrgbClr.FindStringSubmatch(bgXML)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

func extractGradientBgFill(bgXML string) *shapes.ShapeGradientFill {
	matches := reBgGsStop.FindAllStringSubmatch(bgXML, -1)
	if len(matches) < 2 {
		return nil
	}
	stops := make([]shapes.ShapeGradientStop, 0, len(matches))
	for _, m := range matches {
		posRaw, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		// OOXML gradient stop positions are in 1/1000 of a percent (0–100000).
		pct := int(math.Round(float64(posRaw) / 1000.0)) //nolint:mnd // scale: 1000 units per percent
		stops = append(stops, shapes.NewShapeGradientStop(pct, m[2]))
	}
	if len(stops) < 2 {
		return nil
	}
	grad := shapes.NewShapeGradientFill("linear", stops)
	if angM := reBgLinAngle.FindStringSubmatch(bgXML); len(angM) >= 2 {
		if angRaw, err := strconv.Atoi(angM[1]); err == nil {
			grad = grad.WithLinearAngle(angRaw / 60000) //nolint:mnd // OOXML angle unit: 60000 units = 1 degree
		}
	}
	return &grad
}

func extractPictureBgFill(bgXML string, e *PresentationEditor, slidePart string) *shapes.Image {
	m := reBgBlipEmbed.FindStringSubmatch(bgXML)
	if len(m) < 2 {
		return nil
	}
	relID := m[1]
	relsPath := common.RelsPathFor(slidePart)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return nil
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return nil
	}
	var mediaTarget string
	for _, rel := range rels {
		if rel.ID == relID {
			mediaTarget = rel.Target
			break
		}
	}
	if mediaTarget == "" {
		return nil
	}
	mediaPath := path.Clean(path.Join(path.Dir(slidePart), mediaTarget))
	data, ok := e.parts.Get(mediaPath)
	if !ok {
		return nil
	}
	img := shapes.Image{Data: data}
	return &img
}
