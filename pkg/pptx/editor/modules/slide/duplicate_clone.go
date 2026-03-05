package slide

import (
	"fmt"
	"path"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type CloneRelationshipState struct {
	NextChartNum    int
	NextExcelNum    int
	NextNotesNum    int
	ChartEmbeddings map[string]string
	NotesInventory  map[string]string
}

func DeepCloneSlideRelationships(
	srcSlideRelsBytes []byte,
	newSlidePart string,
	state CloneRelationshipState,
	getPart GetPartFn,
	setPart SetPartFn,
	parseRelationships ParseRelationshipsFn,
	renderRelationships RenderRelationshipsFn,
	rewriteChartExternalData RewriteChartExternalDataFn,
) ([]byte, CloneRelationshipState, error) {
	rels, err := parseRelationships(srcSlideRelsBytes)
	if err != nil {
		return nil, state, err
	}

	changed := false
	for i, rel := range rels {
		newTarget, handled := cloneSlideRelationshipPart(
			rel,
			newSlidePart,
			&state,
			getPart,
			setPart,
			parseRelationships,
			renderRelationships,
			rewriteChartExternalData,
		)
		if !handled {
			continue
		}
		rels[i].Target = newTarget
		changed = true
	}

	if changed {
		rendered := renderRelationships(rels)
		return []byte(rendered), state, nil
	}
	return srcSlideRelsBytes, state, nil
}

func cloneSlideRelationshipPart(
	rel common.EditorRelationship,
	newSlidePart string,
	state *CloneRelationshipState,
	getPart GetPartFn,
	setPart SetPartFn,
	parseRelationships ParseRelationshipsFn,
	renderRelationships RenderRelationshipsFn,
	rewriteChartExternalData RewriteChartExternalDataFn,
) (string, bool) {
	switch rel.Type {
	case common.RelTypeChart:
		return cloneChartPart(
			rel,
			state,
			getPart,
			setPart,
			parseRelationships,
			renderRelationships,
			rewriteChartExternalData,
		)
	case common.RelTypeNotesSlide:
		return cloneNotesSlidePart(rel, newSlidePart, state, getPart, setPart, parseRelationships, renderRelationships)
	default:
		return "", false
	}
}

func cloneChartPart(
	rel common.EditorRelationship,
	state *CloneRelationshipState,
	getPart GetPartFn,
	setPart SetPartFn,
	parseRelationships ParseRelationshipsFn,
	renderRelationships RenderRelationshipsFn,
	rewriteChartExternalData RewriteChartExternalDataFn,
) (string, bool) {
	srcChartPart := common.CanonicalPartPath(path.Join("ppt/slides", rel.Target))
	newChartPart := fmt.Sprintf("ppt/charts/chart%d.xml", state.NextChartNum)
	state.NextChartNum++

	data, chartOK := getPart(srcChartPart)
	if !chartOK {
		return "../charts/" + path.Base(newChartPart), true
	}

	newChartData := CloneBytes(data)
	updatedData, nextExcelNum, embeddingPath := CloneChartDependencies(
		srcChartPart,
		newChartPart,
		newChartData,
		state.NextExcelNum,
		getPart,
		setPart,
		parseRelationships,
		renderRelationships,
		rewriteChartExternalData,
	)
	state.NextExcelNum = nextExcelNum
	if embeddingPath != "" {
		state.ChartEmbeddings[newChartPart] = embeddingPath
	}

	setPart(newChartPart, updatedData)
	return "../charts/" + path.Base(newChartPart), true
}

func cloneNotesSlidePart(
	rel common.EditorRelationship,
	newSlidePart string,
	state *CloneRelationshipState,
	getPart GetPartFn,
	setPart SetPartFn,
	parseRelationships ParseRelationshipsFn,
	renderRelationships RenderRelationshipsFn,
) (string, bool) {
	srcNotesPart := common.CanonicalPartPath(path.Join("ppt/slides", rel.Target))
	newNotesPart := fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", state.NextNotesNum)
	state.NextNotesNum++

	data, notesOK := getPart(srcNotesPart)
	if !notesOK {
		return "../notesSlides/" + path.Base(newNotesPart), true
	}

	setPart(newNotesPart, CloneBytes(data))
	CloneNotesRelationships(
		srcNotesPart,
		newNotesPart,
		newSlidePart,
		getPart,
		setPart,
		parseRelationships,
		renderRelationships,
	)
	state.NotesInventory[newSlidePart] = newNotesPart
	return "../notesSlides/" + path.Base(newNotesPart), true
}
