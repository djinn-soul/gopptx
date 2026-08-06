package pptxxml

import (
	"regexp"
	"strings"
)

// The cached drawing is captured under the template's own quick style, so it
// contradicts a diagram that asks for another one. Copying the requested style's
// 3-D scene and shape references onto the cached shapes keeps the cache honest:
// it then carries the same scene, bevels and style references PowerPoint writes
// for that style.
//
// Neither renderer that can be checked here reads the cache — PowerPoint lays
// the diagram out from the definitions, and LibreOffice renders identically with
// the cache emptied — so this changes no picture either of them draws. It exists
// so the part a package ships describes the diagram it belongs to.

const drawingSpPrOpenTag = "<dsp:spPr>"

var (
	drawingShapeOpenPattern = regexp.MustCompile(`<dsp:sp modelId="([^"]+)">`)
	styleLabelPattern       = regexp.MustCompile(`(?s)<dgm:styleLbl name="([^"]+)">(.*?)</dgm:styleLbl>`)
	presStylePointPattern   = regexp.MustCompile(`<dgm:pt modelId="([^"]+)" type="pres">.*?presStyleLbl="([^"]+)"`)
	dgmSceneOrShape3DTag    = regexp.MustCompile(`(?s)<dgm:(scene3d|sp3d)([ >].*?)</dgm:(?:scene3d|sp3d)>`)
	dgmStyleBlockPattern    = regexp.MustCompile(`(?s)<dgm:style>(.*?)</dgm:style>`)
	dspStyleBlockPattern    = regexp.MustCompile(`(?s)<dsp:style>.*?</dsp:style>`)
)

// applySmartArtQuickStyleToDrawing rewrites the cached shapes to carry the
// style the diagram asks for.
func applySmartArtQuickStyleToDrawing(drawing, data, quickStyleID string) string {
	labels := smartArtStyleLabelDefinitions(quickStyleID)
	if len(labels) == 0 {
		return drawing
	}
	labelByModelID := smartArtPresStyleLabels(data)
	if len(labelByModelID) == 0 {
		return drawing
	}

	matches := drawingShapeOpenPattern.FindAllStringSubmatchIndex(drawing, -1)
	if len(matches) == 0 {
		return drawing
	}

	var b strings.Builder
	last := 0
	for _, idx := range matches {
		shapeStart, shapeEnd := idx[0], idx[1]
		modelID := drawing[idx[2]:idx[3]]
		label, ok := labels[labelByModelID[modelID]]
		if !ok {
			continue
		}
		bodyEnd := smartArtShapeBodyEnd(drawing, shapeEnd)
		if bodyEnd < 0 {
			continue
		}
		b.WriteString(drawing[last:shapeEnd])
		b.WriteString(applySmartArtStyleLabelToShape(drawing[shapeEnd:bodyEnd], label))
		last = bodyEnd
		_ = shapeStart
	}
	b.WriteString(drawing[last:])
	return b.String()
}

// smartArtShapeBodyEnd finds where one cached shape ends, so the next one is
// not rewritten with this one's style.
func smartArtShapeBodyEnd(drawing string, from int) int {
	end := strings.Index(drawing[from:], "</dsp:sp>")
	if end < 0 {
		return -1
	}
	return from + end
}

func applySmartArtStyleLabelToShape(shape string, label smartArtStyleLabel) string {
	if label.scene3D != "" && !strings.Contains(shape, "<a:scene3d>") {
		shape = strings.Replace(shape, drawingSpPrOpenTag, drawingSpPrOpenTag+label.scene3D+label.shape3D, 1)
	}
	if label.style != "" {
		shape = dspStyleBlockPattern.ReplaceAllLiteralString(shape, "<dsp:style>"+label.style+"</dsp:style>")
	}
	return shape
}

type smartArtStyleLabel struct {
	scene3D string
	shape3D string
	style   string
}

// smartArtStyleLabelDefinitions reads each style label out of a quick style,
// with its elements renamed from the diagram namespace to the drawing one.
func smartArtStyleLabelDefinitions(quickStyleID string) map[string]smartArtStyleLabel {
	definition := smartArtStyleDefinition(defaultQuickStyleID(quickStyleID))
	out := map[string]smartArtStyleLabel{}
	for _, match := range styleLabelPattern.FindAllStringSubmatch(definition, -1) {
		name, body := match[1], match[2]
		label := smartArtStyleLabel{}
		for _, part := range dgmSceneOrShape3DTag.FindAllStringSubmatch(body, -1) {
			converted := "<a:" + part[1] + part[2] + "</a:" + part[1] + ">"
			if part[1] == "scene3d" {
				label.scene3D = converted
				continue
			}
			label.shape3D = converted
		}
		if style := dgmStyleBlockPattern.FindStringSubmatch(body); len(style) > 1 {
			label.style = style[1]
		}
		out[name] = label
	}
	return out
}

// smartArtPresStyleLabels maps each cached shape to the style label its
// presentation point was given.
func smartArtPresStyleLabels(data string) map[string]string {
	out := map[string]string{}
	for _, segment := range strings.Split(data, "<dgm:pt ")[1:] {
		pt := "<dgm:pt " + segment
		match := presStylePointPattern.FindStringSubmatch(pt)
		if len(match) < presStyleLabelSubmatches {
			continue
		}
		out[match[1]] = match[2]
	}
	return out
}

const presStyleLabelSubmatches = 3
