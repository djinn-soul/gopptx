package editor

import (
	"errors"
	"fmt"
	"path"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorgrayscale "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/grayscale"
)

// GrayscaleOptions configures selective grayscale conversion.
type GrayscaleOptions = editorgrayscale.Options

type grayscaleTargets struct {
	slides map[int]struct{}
	shapes map[string]struct{}
	text   map[string]map[int]struct{}
}

const bgEmbedRelIDSubmatchGroup = 1

// ConvertToGrayscale applies grayscale conversion to selected slide content.
func (e *PresentationEditor) ConvertToGrayscale(opts GrayscaleOptions) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if !opts.Colors && !opts.Images && !opts.Backgrounds {
		return errors.New("at least one grayscale target must be enabled")
	}
	targets, err := e.buildGrayscaleTargets(opts)
	if err != nil {
		return err
	}
	for slideIndex := range targets.slides {
		if opts.Colors {
			if err := e.convertSlideShapeColors(slideIndex, targets); err != nil {
				return err
			}
		}
		if opts.Images {
			if err := e.convertSlideImages(slideIndex, targets); err != nil {
				return err
			}
		}
		if opts.Backgrounds {
			if err := e.convertSlideBackground(slideIndex); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *PresentationEditor) buildGrayscaleTargets(opts GrayscaleOptions) (grayscaleTargets, error) {
	targets := grayscaleTargets{
		slides: make(map[int]struct{}, len(e.slides)),
		shapes: make(map[string]struct{}, len(opts.Shapes)),
		text:   make(map[string]map[int]struct{}, len(opts.Text)),
	}
	if len(opts.Slides) == 0 {
		for slideIndex := range e.slides {
			targets.slides[slideIndex] = struct{}{}
		}
	} else {
		for _, slideIndex := range opts.Slides {
			if err := e.validateSlideIndex(slideIndex); err != nil {
				return grayscaleTargets{}, err
			}
			targets.slides[slideIndex] = struct{}{}
		}
	}
	for _, ref := range opts.Shapes {
		if err := e.validateSlideIndex(ref.SlideIndex); err != nil {
			return grayscaleTargets{}, err
		}
		targets.shapes[shapeTargetKey(ref.SlideIndex, ref.ShapeID)] = struct{}{}
		targets.slides[ref.SlideIndex] = struct{}{}
	}
	for _, ref := range opts.Text {
		if err := e.validateSlideIndex(ref.SlideIndex); err != nil {
			return grayscaleTargets{}, err
		}
		key := shapeTargetKey(ref.SlideIndex, ref.ShapeID)
		runSet := make(map[int]struct{}, len(ref.RunIndices))
		for _, runIndex := range ref.RunIndices {
			if runIndex < 0 {
				return grayscaleTargets{}, fmt.Errorf("run index %d must be >= 0", runIndex)
			}
			runSet[runIndex] = struct{}{}
		}
		targets.text[key] = runSet
		targets.slides[ref.SlideIndex] = struct{}{}
	}
	if err := e.addGrayscalePlaceholderTargets(&targets, opts.Placeholders); err != nil {
		return grayscaleTargets{}, err
	}
	return targets, nil
}

func (e *PresentationEditor) convertSlideShapeColors(slideIndex int, targets grayscaleTargets) error {
	slideRef := e.slides[slideIndex]
	content, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return fmt.Errorf("slide part %q not found", slideRef.Part)
	}
	shapes, err := parseSlideShapes(content)
	if err != nil {
		return err
	}
	var replaceErr error
	updated := replaceShapeNodes(content, shapes, func(_ int, s *parsedShape) ([]byte, bool) {
		wholeShape := targets.wholeShapeSelected(slideIndex, s.ID)
		runSelection, hasRunSelection := targets.runSelection(slideIndex, s.ID)
		if !wholeShape && !hasRunSelection {
			return nil, false
		}
		if s.Type == shapeTypePicture || s.IsGroup || s.Type == shapeTypeGraphicFrame {
			return nil, false
		}
		if !applyGrayscaleToShape(s, wholeShape, runSelection) {
			return nil, false
		}
		xmlBytes, renderErr := e.renderShapeXML(slideRef.Part, s)
		if renderErr != nil {
			replaceErr = renderErr
			return nil, false
		}
		return xmlBytes, len(xmlBytes) > 0
	})
	if replaceErr != nil {
		return replaceErr
	}
	e.parts.Set(slideRef.Part, updated)
	return nil
}

func (e *PresentationEditor) convertSlideImages(slideIndex int, targets grayscaleTargets) error {
	slideRef := e.slides[slideIndex]
	content, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return fmt.Errorf("slide part %q not found", slideRef.Part)
	}
	shapes, err := parseSlideShapes(content)
	if err != nil {
		return err
	}
	for _, shape := range shapes {
		if shape.Type != shapeTypePicture {
			continue
		}
		if !targets.imageSelected(slideIndex, shape.ID) {
			continue
		}
		relID := extractEmbedRelID(content[shape.Start:shape.End])
		if relID == "" {
			continue
		}
		if err := e.grayscaleImageRelationship(slideIndex, relID); err != nil {
			return err
		}
	}
	return nil
}

func (e *PresentationEditor) convertSlideBackground(slideIndex int) error {
	slideRef := e.slides[slideIndex]
	content, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return fmt.Errorf("slide part %q not found", slideRef.Part)
	}
	updated, changed := grayscaleBackgroundXML(content)
	if changed {
		e.parts.Set(slideRef.Part, updated)
		content = updated
	}
	match := bgEmbedPattern.FindSubmatch(content)
	if len(match) <= bgEmbedRelIDSubmatchGroup {
		return nil
	}
	return e.grayscaleImageRelationship(slideIndex, string(match[bgEmbedRelIDSubmatchGroup]))
}

func (e *PresentationEditor) grayscaleImageRelationship(slideIndex int, relID string) error {
	partPath, format, err := e.resolveImagePart(slideIndex, relID)
	if err != nil {
		return err
	}
	data, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("media part %q not found", partPath)
	}
	gray, outFormat, err := editorgrayscale.ImageBytes(data, format)
	if err != nil {
		return err
	}
	return e.swapImageByRelID(slideIndex, relID, gray, outFormat)
}

func (e *PresentationEditor) resolveImagePart(slideIndex int, relID string) (string, string, error) {
	rels, err := e.slideRelationships(e.slides[slideIndex].Part)
	if err != nil {
		return "", "", err
	}
	for _, rel := range rels {
		if rel.ID != relID || rel.Type != common.RelTypeImage {
			continue
		}
		partPath := common.CanonicalPartPath(path.Join("ppt/slides", rel.Target))
		format := imageFormatFromTarget(rel.Target)
		return partPath, format, nil
	}
	return "", "", fmt.Errorf("image relationship %q not found on slide %d", relID, slideIndex)
}
