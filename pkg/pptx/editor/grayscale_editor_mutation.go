package editor

import (
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorgrayscale "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/grayscale"
)

func applyGrayscaleToShape(s *parsedShape, wholeShape bool, runSelection map[int]struct{}) bool {
	changed := false
	if wholeShape {
		changed = grayscaleShapeVisuals(s) || changed
	}
	if len(s.Runs) == 0 {
		return changed
	}
	for i := range s.Runs {
		if !wholeShape && !runSelected(runSelection, i) {
			continue
		}
		if grayscaleRun(&s.Runs[i]) {
			changed = true
		}
	}
	if wholeShape && s.Paragraph != nil && s.Paragraph.BulletColor != nil {
		if gray, err := editorgrayscale.HexColor(*s.Paragraph.BulletColor); err == nil {
			s.Paragraph.BulletColor = &gray
			changed = true
		}
	}
	return changed
}

func grayscaleShapeVisuals(s *parsedShape) bool {
	changed := false
	if s.Fill != nil {
		changed = grayscaleFill(s.Fill) || changed
	}
	if s.Line != nil && s.Line.Color != nil {
		if gray, err := editorgrayscale.HexColor(*s.Line.Color); err == nil {
			s.Line.Color = &gray
			changed = true
		}
	}
	if s.Shadow != nil && s.Shadow.Color != nil {
		if gray, err := editorgrayscale.HexColor(*s.Shadow.Color); err == nil {
			s.Shadow.Color = &gray
			changed = true
		}
	}
	if s.Glow != nil && s.Glow.Color != nil {
		if gray, err := editorgrayscale.HexColor(*s.Glow.Color); err == nil {
			s.Glow.Color = &gray
			changed = true
		}
	}
	return changed
}

func grayscaleFill(fill *common.ShapeFill) bool {
	changed := false
	if fill.Solid != nil {
		if gray, err := editorgrayscale.HexColor(*fill.Solid); err == nil {
			fill.Solid = &gray
			changed = true
		}
	}
	if fill.Pattern == nil {
		return changed
	}
	if gray, ok := grayscaleColorPtr(fill.Pattern.FgColor); ok {
		fill.Pattern.FgColor = &gray
		changed = true
	}
	if gray, ok := grayscaleColorPtr(fill.Pattern.BgColor); ok {
		fill.Pattern.BgColor = &gray
		changed = true
	}
	if fill.Gradient != nil {
		for i := range fill.Gradient.Stops {
			gray, err := editorgrayscale.HexColor(fill.Gradient.Stops[i].Color)
			if err != nil {
				continue
			}
			fill.Gradient.Stops[i].Color = gray
			changed = true
		}
	}
	return changed
}

func grayscaleRun(run *common.TextRun) bool {
	changed := false
	for _, field := range []*string{run.Color, run.Highlight, run.OutlineColor} {
		if field == nil {
			continue
		}
		gray, err := editorgrayscale.HexColor(*field)
		if err != nil {
			continue
		}
		*field = gray
		changed = true
	}
	return changed
}

func (e *PresentationEditor) validateSlideIndex(slideIndex int) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range [0,%d)", slideIndex, len(e.slides))
	}
	return nil
}

func (t grayscaleTargets) wholeShapeSelected(slideIndex, shapeID int) bool {
	if len(t.shapes) == 0 && len(t.text) == 0 {
		return true
	}
	_, ok := t.shapes[shapeTargetKey(slideIndex, shapeID)]
	return ok
}

func (t grayscaleTargets) runSelection(slideIndex, shapeID int) (map[int]struct{}, bool) {
	runSet, ok := t.text[shapeTargetKey(slideIndex, shapeID)]
	return runSet, ok
}

func (t grayscaleTargets) imageSelected(slideIndex, shapeID int) bool {
	if len(t.shapes) == 0 {
		return true
	}
	_, ok := t.shapes[shapeTargetKey(slideIndex, shapeID)]
	return ok
}

func shapeTargetKey(slideIndex, shapeID int) string {
	return fmt.Sprintf("%d:%d", slideIndex, shapeID)
}

func runSelected(runSet map[int]struct{}, runIndex int) bool {
	if len(runSet) == 0 {
		return true
	}
	_, ok := runSet[runIndex]
	return ok
}

func extractEmbedRelID(shapeXML []byte) string {
	match := embeddedImageRelPattern.FindSubmatch(shapeXML)
	if len(match) <= bgEmbedRelIDSubmatchGroup {
		return ""
	}
	return string(match[bgEmbedRelIDSubmatchGroup])
}

func grayscaleColorPtr(value *string) (string, bool) {
	if value == nil {
		return "", false
	}
	gray, err := editorgrayscale.HexColor(*value)
	if err != nil {
		return "", false
	}
	return gray, true
}
