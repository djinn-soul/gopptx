package editor

import (
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// GetEffectiveShapeStyle resolves how a shape actually looks, following the
// inheritance chain a renderer follows: the shape's own properties, then its
// layout placeholder, then the master placeholder or the master's text styles,
// then the theme.
//
// A placeholder that sets no colour of its own reports None on every direct
// property, which is the confusion behind upstream #1013: the colour is real,
// it just lives further up the chain.
func (e *PresentationEditor) GetEffectiveShapeStyle(
	slideIndex, shapeID int,
) (*common.EffectiveShapeStyle, error) {
	if e == nil || e.parts == nil {
		return nil, errors.New("editor cannot be nil")
	}
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, fmt.Errorf("slide index %d out of range", slideIndex)
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, fmt.Errorf("read slide part %s: not found", partPath)
	}
	shapes, err := parseSlideShapes(content)
	if err != nil {
		return nil, fmt.Errorf("parse shapes: %w", err)
	}

	var shape *parsedShape
	for i := range shapes {
		if shapes[i].ID == shapeID {
			shape = &shapes[i]
			break
		}
	}
	if shape == nil {
		return nil, fmt.Errorf("shape %d not found on slide %d", shapeID, slideIndex)
	}

	style := &common.EffectiveShapeStyle{}
	applyShapeOwnStyle(style, shape, content)

	chain := e.styleInheritanceChain(slideIndex, shape)
	style.LayoutPart = chain.layoutPart
	style.MasterPart = chain.masterPart
	for _, level := range chain.levels {
		applyInheritedStyle(style, level)
	}
	applyThemeStyle(style, chain.theme)

	resolveSchemeColors(style, chain.theme)
	return style, nil
}

// applyShapeOwnStyle takes everything the shape states directly. Direct
// properties always win, so nothing later overwrites what this sets.
func applyShapeOwnStyle(style *common.EffectiveShapeStyle, shape *parsedShape, content []byte) {
	if color, slot, found := solidFillColor(content[shape.Start:shape.End]); found {
		style.FillColor = &common.EffectiveColor{
			RGB:        color,
			SchemeSlot: slot,
			Source:     common.StyleSourceShape,
		}
	}
	if shape.W > 0 && shape.H > 0 {
		style.Position = &common.EffectivePosition{
			X:      shape.X,
			Y:      shape.Y,
			W:      shape.W,
			H:      shape.H,
			Source: common.StyleSourceShape,
		}
	}
	if len(shape.Runs) == 0 {
		return
	}
	run := shape.Runs[0]
	if run.Color != nil && *run.Color != "" {
		style.FontColor = &common.EffectiveColor{RGB: *run.Color, Source: common.StyleSourceShape}
	}
	if run.Font != nil && *run.Font != "" {
		style.FontTypeface = &common.EffectiveString{Value: *run.Font, Source: common.StyleSourceShape}
	}
	if run.SizePt != nil {
		style.FontSizePt = &common.EffectiveFloat{
			Value:  float64(*run.SizePt),
			Source: common.StyleSourceShape,
		}
	}
	if run.Bold != nil {
		style.Bold = &common.EffectiveBool{Value: *run.Bold, Source: common.StyleSourceShape}
	}
	if run.Italic != nil {
		style.Italic = &common.EffectiveBool{Value: *run.Italic, Source: common.StyleSourceShape}
	}
}

// applyInheritedStyle fills only the values still unset, so the nearer level in
// the chain keeps its win.
func applyInheritedStyle(style *common.EffectiveShapeStyle, level inheritedStyle) {
	if style.FillColor == nil && level.fillRGB != "" {
		style.FillColor = &common.EffectiveColor{
			RGB:        level.fillRGB,
			SchemeSlot: level.fillScheme,
			Source:     level.source,
		}
	}
	if style.FontColor == nil && (level.fontRGB != "" || level.fontScheme != "") {
		style.FontColor = &common.EffectiveColor{
			RGB:        level.fontRGB,
			SchemeSlot: level.fontScheme,
			Source:     level.source,
		}
	}
	if style.FontTypeface == nil && level.typeface != "" {
		style.FontTypeface = &common.EffectiveString{Value: level.typeface, Source: level.source}
	}
	if style.FontSizePt == nil && level.sizePt > 0 {
		style.FontSizePt = &common.EffectiveFloat{Value: level.sizePt, Source: level.source}
	}
	if style.Bold == nil && level.bold != nil {
		style.Bold = &common.EffectiveBool{Value: *level.bold, Source: level.source}
	}
	if style.Italic == nil && level.italic != nil {
		style.Italic = &common.EffectiveBool{Value: *level.italic, Source: level.source}
	}
	if style.Position == nil && level.position != nil {
		position := *level.position
		position.Source = level.source
		style.Position = &position
	}
}

// applyThemeStyle supplies the typeface a "+mj-lt"/"+mn-lt" reference stands
// for, and the theme default when nothing in the chain named a font at all.
func applyThemeStyle(style *common.EffectiveShapeStyle, theme themeStyleContext) {
	if style.FontTypeface == nil {
		if theme.minorLatin == "" {
			return
		}
		style.FontTypeface = &common.EffectiveString{
			Value:  theme.minorLatin,
			Source: common.StyleSourceTheme,
		}
		return
	}
	resolved, ok := theme.resolveTypeface(style.FontTypeface.Value)
	if !ok {
		return
	}
	style.FontTypeface = &common.EffectiveString{Value: resolved, Source: common.StyleSourceTheme}
}

// resolveSchemeColors turns a scheme reference collected from any level into
// the concrete RGB the theme gives it.
func resolveSchemeColors(style *common.EffectiveShapeStyle, theme themeStyleContext) {
	for _, color := range []*common.EffectiveColor{style.FillColor, style.FontColor} {
		if color == nil || color.SchemeSlot == "" {
			continue
		}
		if rgb, ok := theme.resolveSchemeSlot(color.SchemeSlot); ok {
			color.RGB = rgb
		}
	}
}
