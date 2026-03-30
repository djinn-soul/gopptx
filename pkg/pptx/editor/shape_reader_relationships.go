package editor

import (
	"net/url"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

const (
	actionSlideJumpPrefix = "ppaction://hlinksldjump"
	actionShowJumpPrefix  = "ppaction://hlinkshowjump?jump="
	actionMacroPrefix     = "ppaction://macro?name="
)

func (e *PresentationEditor) enrichParsedShapeRelationships(partPath string, parsed []parsedShape) error {
	relsByID, err := e.slideRelationshipsByID(partPath)
	if err != nil {
		return err
	}
	slideIndexByPart := make(map[string]int, len(e.slides))
	for idx, slide := range e.slides {
		slideIndexByPart[slide.Part] = idx
	}

	for idx := range parsed {
		parsed[idx].ClickAction = resolveReaderHyperlink(
			partPath,
			parsed[idx].ClickActionRef,
			relsByID,
			slideIndexByPart,
		)
		parsed[idx].HoverAction = resolveReaderHyperlink(
			partPath,
			parsed[idx].HoverActionRef,
			relsByID,
			slideIndexByPart,
		)
		applyResolvedRunActions(&parsed[idx], partPath, relsByID, slideIndexByPart)
	}
	return nil
}

func (e *PresentationEditor) slideRelationshipsByID(partPath string) (map[string]common.EditorRelationship, error) {
	relsPath := common.SlideRelsPartName(partPath)
	data, ok := e.parts.Get(relsPath)
	if !ok {
		return map[string]common.EditorRelationship{}, nil
	}
	rels, err := parseRelationshipsXML(data)
	if err != nil {
		return nil, err
	}
	out := make(map[string]common.EditorRelationship, len(rels))
	for _, rel := range rels {
		out[rel.ID] = rel
	}
	return out, nil
}

func applyResolvedRunActions(
	shape *parsedShape,
	partPath string,
	relsByID map[string]common.EditorRelationship,
	slideIndexByPart map[string]int,
) {
	flatIndex := 0
	for pIdx := range shape.Paragraphs {
		for rIdx := range shape.Paragraphs[pIdx].Runs {
			if pIdx < len(shape.RunActions) && rIdx < len(shape.RunActions[pIdx]) {
				actions := shape.RunActions[pIdx][rIdx]
				shape.Paragraphs[pIdx].Runs[rIdx].Hyperlink = resolveReaderHyperlink(
					partPath,
					actions.ClickAction,
					relsByID,
					slideIndexByPart,
				)
				shape.Paragraphs[pIdx].Runs[rIdx].HoverAction = resolveReaderHyperlink(
					partPath,
					actions.HoverAction,
					relsByID,
					slideIndexByPart,
				)
				if flatIndex < len(shape.Runs) {
					shape.Runs[flatIndex].Hyperlink = shape.Paragraphs[pIdx].Runs[rIdx].Hyperlink
					shape.Runs[flatIndex].HoverAction = shape.Paragraphs[pIdx].Runs[rIdx].HoverAction
				}
			}
			flatIndex++
		}
	}
}

func resolveReaderHyperlink(
	partPath string,
	ref *editorshape.ReaderHyperlinkRef,
	relsByID map[string]common.EditorRelationship,
	slideIndexByPart map[string]int,
) *common.Hyperlink {
	if ref == nil {
		return nil
	}
	hl := &common.Hyperlink{
		Action:         cloneEditorString(ref.Action),
		Tooltip:        cloneEditorString(ref.Tooltip),
		History:        cloneEditorBool(ref.History),
		HighlightClick: cloneEditorBool(ref.HighlightClick),
		EndSound:       cloneEditorBool(ref.EndSound),
	}
	applyReaderActionString(hl, ref.Action)
	if rel, ok := relsByID[ref.RelID]; ok {
		applyReaderRelationshipTarget(partPath, hl, rel, slideIndexByPart)
	}
	if !hasResolvedHyperlinkData(hl) {
		return nil
	}
	return hl
}

func applyReaderActionString(hl *common.Hyperlink, actionValue *string) {
	if actionValue == nil {
		return
	}
	action := strings.TrimSpace(*actionValue)
	switch {
	case strings.HasPrefix(action, actionShowJumpPrefix):
		jump := strings.TrimPrefix(action, actionShowJumpPrefix)
		hl.TargetJump = cloneEditorString(&jump)
	case strings.HasPrefix(action, actionMacroPrefix):
		macro := strings.TrimPrefix(action, actionMacroPrefix)
		hl.Macro = cloneEditorString(&macro)
	case strings.HasPrefix(action, actionSlideJumpPrefix):
		// Resolved from the slide relationship target.
	default:
		hl.Action = cloneEditorString(&action)
	}
}

func applyReaderRelationshipTarget(
	partPath string,
	hl *common.Hyperlink,
	rel common.EditorRelationship,
	slideIndexByPart map[string]int,
) {
	switch rel.Type {
	case common.RelTypeSlide:
		targetPart := common.ResolveRelationshipTarget(partPath, rel.Target)
		if slideIndex, ok := slideIndexByPart[targetPart]; ok {
			hl.TargetSlide = &slideIndex
		}
	case common.RelTypeHyperlink:
		address := strings.TrimSpace(rel.Target)
		if address == "" {
			return
		}
		if decoded, err := url.QueryUnescape(address); err == nil {
			address = decoded
		}
		hl.Address = &address
	default:
		address := strings.TrimSpace(rel.Target)
		if address != "" {
			hl.Address = &address
		}
	}
}

func hasResolvedHyperlinkData(hl *common.Hyperlink) bool {
	return hl != nil && (hl.Address != nil || hl.Action != nil || hl.Tooltip != nil ||
		hl.TargetSlide != nil || hl.TargetJump != nil || hl.Macro != nil || hl.History != nil ||
		hl.HighlightClick != nil || hl.EndSound != nil)
}

func cloneEditorString(src *string) *string {
	if src == nil {
		return nil
	}
	value := strings.TrimSpace(*src)
	return &value
}

func cloneEditorBool(src *bool) *bool {
	if src == nil {
		return nil
	}
	value := *src
	return &value
}
