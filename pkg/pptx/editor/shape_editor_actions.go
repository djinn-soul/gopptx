package editor

import (
	"fmt"
	"path/filepath"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

func (e *PresentationEditor) getOrCreateHyperlinkRelID(partPath, address string) (string, error) {
	relsPath := common.SlideRelsPartName(partPath)
	rels := make([]common.EditorRelationship, 0)
	if data, ok := e.parts.Get(relsPath); ok {
		parsed, err := parseRelationshipsXML(data)
		if err != nil {
			return "", fmt.Errorf("parse %s: %w", relsPath, err)
		}
		rels = parsed
	}
	for _, r := range rels {
		if r.Type == common.RelTypeHyperlink && r.Target == address {
			return r.ID, nil
		}
	}
	relID, err := e.nextSlideRelID(partPath)
	if err != nil {
		return "", err
	}
	if err := e.addRelationship(partPath, relID, common.RelTypeHyperlink, address); err != nil {
		return "", err
	}
	return relID, nil
}

func (e *PresentationEditor) buildClickActionXML(partPath string, hl *common.Hyperlink) (string, error) {
	return e.buildActionXML(partPath, hl, "hlinkClick")
}

func (e *PresentationEditor) buildHoverActionXML(partPath string, hl *common.Hyperlink) (string, error) {
	return e.buildActionXML(partPath, hl, "hlinkMouseOver")
}

func (e *PresentationEditor) buildActionXML(partPath string, hl *common.Hyperlink, tag string) (string, error) {
	if hl == nil || partPath == "" {
		return "", nil
	}
	if err := editorshape.ValidateHyperlinkAction(hl); err != nil {
		return "", err
	}

	action := strings.TrimSpace(editorshape.GetStr(hl.Action))
	if action == "" {
		action = editorshape.DeriveActionURL(hl)
	}

	attrs := make([]string, 0, actionAttrCapHint)
	if hl.Address != nil && *hl.Address != "" {
		relID, err := e.getOrCreateHyperlinkRelID(partPath, *hl.Address)
		if err != nil {
			return "", fmt.Errorf("allocate hyperlink relationship id: %w", err)
		}
		attrs = append(attrs, fmt.Sprintf(`r:id="%s"`, editorshape.XMLEscape(relID)))
	} else if hl.TargetSlide != nil {
		relID, err := e.getOrCreateSlideJumpRelID(partPath, *hl.TargetSlide)
		if err != nil {
			return "", err
		}
		attrs = append(attrs, fmt.Sprintf(`r:id="%s"`, editorshape.XMLEscape(relID)))
	}
	if action != "" {
		attrs = append(attrs, fmt.Sprintf(`action="%s"`, editorshape.XMLEscape(action)))
	}
	if tooltip := strings.TrimSpace(editorshape.GetStr(hl.Tooltip)); tooltip != "" {
		attrs = append(attrs, fmt.Sprintf(`tooltip="%s"`, editorshape.XMLEscape(tooltip)))
	}
	if len(attrs) == 0 {
		return "", nil
	}
	return fmt.Sprintf(`<a:%s %s/>`, tag, strings.Join(attrs, " ")), nil
}

func (e *PresentationEditor) getOrCreateSlideJumpRelID(partPath string, targetSlideIndex int) (string, error) {
	if targetSlideIndex < 0 || targetSlideIndex >= len(e.slides) {
		return "", fmt.Errorf("target_slide index %d out of range", targetSlideIndex)
	}
	targetPart := e.slides[targetSlideIndex].Part
	relsPath := common.SlideRelsPartName(partPath)
	rels := make([]common.EditorRelationship, 0)
	if data, ok := e.parts.Get(relsPath); ok {
		parsed, err := parseRelationshipsXML(data)
		if err != nil {
			return "", fmt.Errorf("parse %s: %w", relsPath, err)
		}
		rels = parsed
	}
	sourceDir := filepath.Dir(partPath)
	targetRelPath, err := filepath.Rel(sourceDir, targetPart)
	if err != nil {
		return "", fmt.Errorf("resolve target slide relationship path: %w", err)
	}
	targetRelPath = strings.ReplaceAll(targetRelPath, "\\", "/")
	for _, r := range rels {
		if r.Type == common.RelTypeSlide && r.Target == targetRelPath {
			return r.ID, nil
		}
	}
	relID, err := e.nextSlideRelID(partPath)
	if err != nil {
		return "", err
	}
	if err := e.addRelationship(partPath, relID, common.RelTypeSlide, targetRelPath); err != nil {
		return "", err
	}
	return relID, nil
}

// replaceShapeTextBody replaces the entire <p:txBody> node with a newly constructed one based on Text/Runs.
func replaceShapeTextBody(
	e *PresentationEditor,
	partPath string,
	xmlData []byte,
	s *parsedShape,
) ([]byte, error) {
	txBody, err := renderTextBodyXML(e, partPath, s)
	if err != nil {
		return nil, err
	}
	return editorshape.ReplaceShapeTextBody(xmlData, txBody), nil
}

func replaceShapeClickAction(
	e *PresentationEditor,
	partPath string,
	xmlData []byte,
	clickAction *common.Hyperlink,
) ([]byte, error) {
	return replaceShapeActions(e, partPath, xmlData, clickAction, nil)
}

func replaceShapeStyle(
	xmlData []byte,
	fill *common.ShapeFill,
	line *common.ShapeLine,
	shadow *common.ShapeShadow,
	glow *common.ShapeGlow,
	blur *common.ShapeBlur,
	softEdge *common.ShapeSoftEdge,
	reflection *common.ShapeReflection,
) ([]byte, error) {
	return replaceShapeStyleSelective(
		xmlData,
		fill,
		line,
		shadow,
		glow,
		blur,
		softEdge,
		reflection,
		true,
		true,
		true,
	)
}

func replaceShapeStyleSelective(
	xmlData []byte,
	fill *common.ShapeFill,
	line *common.ShapeLine,
	shadow *common.ShapeShadow,
	glow *common.ShapeGlow,
	blur *common.ShapeBlur,
	softEdge *common.ShapeSoftEdge,
	reflection *common.ShapeReflection,
	applyFill bool,
	applyLine bool,
	applyEffects bool,
) ([]byte, error) {
	if !applyFill && !applyLine && !applyEffects {
		return xmlData, nil
	}

	styleXML, err := buildSelectiveStyleXML(
		fill,
		line,
		shadow,
		glow,
		blur,
		softEdge,
		reflection,
		applyFill,
		applyLine,
		applyEffects,
	)
	if err != nil {
		return nil, err
	}
	return editorshape.ReplaceStyleInSpPr(xmlData, styleXML, applyFill, applyLine, applyEffects), nil
}

func buildSelectiveStyleXML(
	fill *common.ShapeFill,
	line *common.ShapeLine,
	shadow *common.ShapeShadow,
	glow *common.ShapeGlow,
	blur *common.ShapeBlur,
	softEdge *common.ShapeSoftEdge,
	reflection *common.ShapeReflection,
	applyFill bool,
	applyLine bool,
	applyEffects bool,
) (string, error) {
	var styleBuilder strings.Builder
	if applyFill {
		fillXML, err := editorshape.RenderFillXML(fill)
		if err != nil {
			return "", err
		}
		styleBuilder.WriteString(fillXML)
	}
	if applyLine {
		lineXML, err := editorshape.RenderLineXML(line)
		if err != nil {
			return "", err
		}
		styleBuilder.WriteString(lineXML)
	}
	if applyEffects {
		effectsXML, err := editorshape.RenderEffectsXML(shadow, glow, blur, softEdge, reflection)
		if err != nil {
			return "", err
		}
		styleBuilder.WriteString(effectsXML)
	}
	return styleBuilder.String(), nil
}

func replaceShapeActions(
	e *PresentationEditor,
	partPath string,
	xmlData []byte,
	clickAction *common.Hyperlink,
	hoverAction *common.Hyperlink,
) ([]byte, error) {
	clickXML, err := e.buildClickActionXML(partPath, clickAction)
	if err != nil {
		return nil, err
	}
	hoverXML, err := e.buildHoverActionXML(partPath, hoverAction)
	if err != nil {
		return nil, err
	}
	return editorshape.ApplyCNvPrActions(xmlData, clickAction != nil, hoverAction != nil, clickXML, hoverXML)
}
