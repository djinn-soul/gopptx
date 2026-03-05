package editor

import (
	"fmt"
	"path"
	"regexp"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

var layoutNumPattern = regexp.MustCompile(`^slideLayout(\d+)\.xml$`)
var masterNumPattern = regexp.MustCompile(`^slideMaster(\d+)\.xml$`)
var themeNumPattern = regexp.MustCompile(`^theme(\d+)\.xml$`)

const nextPartPatternSubmatchSize = 2

func (e *PresentationEditor) CloneLayoutMasterFamily(layoutPart string) (common.SlideMasterCloneResult, error) {
	sourceMaster, layoutFamily, err := editorslide.CloneFamilyInputs(
		layoutPart,
		e.parts.Has,
		common.CanonicalPartPath,
		func(part string) (string, error) {
			return editorslide.ResolveLayoutMasterPart(part, e.parts.Get, parseRelationshipsXML)
		},
		e.layoutsForMaster,
	)
	if err != nil {
		return common.SlideMasterCloneResult{}, err
	}

	newMaster := editorslide.NextMasterPartPath(
		editorslide.NextPartNumber(
			e.parts.KeysWithPrefix("ppt/slideMasters/"),
			masterNumPattern,
			nextPartPatternSubmatchSize,
		),
	)
	layoutMap := editorslide.BuildLayoutCloneMap(
		layoutFamily,
		editorslide.NextPartNumber(
			e.parts.KeysWithPrefix("ppt/slideLayouts/"),
			layoutNumPattern,
			nextPartPatternSubmatchSize,
		),
	)
	masterXML, masterRels, err := e.loadMasterCloneSource(sourceMaster)
	if err != nil {
		return common.SlideMasterCloneResult{}, err
	}

	themePart, newThemePart, err := e.cloneMasterTheme(masterRels, sourceMaster)
	if err != nil {
		return common.SlideMasterCloneResult{}, err
	}
	if err := e.cloneLayoutParts(layoutMap, newMaster); err != nil {
		return common.SlideMasterCloneResult{}, err
	}
	e.writeClonedMaster(sourceMaster, newMaster, masterXML, masterRels, layoutMap, newThemePart)
	if err := e.registerClonedMaster(newMaster); err != nil {
		return common.SlideMasterCloneResult{}, err
	}

	return common.SlideMasterCloneResult{
		MasterPart: newMaster,
		ThemePart:  editorslide.CloneResultTheme(themePart, newThemePart),
		LayoutMap:  layoutMap,
	}, nil
}

func (e *PresentationEditor) loadMasterCloneSource(
	sourceMaster string,
) ([]byte, []common.EditorRelationship, error) {
	masterXML, ok := e.parts.Get(sourceMaster)
	if !ok {
		return nil, nil, fmt.Errorf("master part not found: %s", sourceMaster)
	}
	masterRelsPath := common.RelsPathFor(sourceMaster)
	masterRelsData, ok := e.parts.Get(masterRelsPath)
	if !ok {
		return nil, nil, fmt.Errorf("master rels part not found: %s", masterRelsPath)
	}
	masterRels, err := parseRelationshipsXML(masterRelsData)
	if err != nil {
		return nil, nil, fmt.Errorf("parse master rels: %w", err)
	}
	return masterXML, masterRels, nil
}

func (e *PresentationEditor) cloneLayoutParts(layoutMap map[string]string, newMaster string) error {
	for oldLayout, clonedLayout := range layoutMap {
		layoutXML, layoutOK := e.parts.Get(oldLayout)
		if !layoutOK {
			return fmt.Errorf("layout part not found: %s", oldLayout)
		}
		e.parts.Set(clonedLayout, append([]byte(nil), layoutXML...))

		layoutRelsPath := common.RelsPathFor(oldLayout)
		layoutRelsData, relsOK := e.parts.Get(layoutRelsPath)
		if !relsOK {
			return fmt.Errorf("layout rels missing: %s", layoutRelsPath)
		}
		layoutRels, parseErr := parseRelationshipsXML(layoutRelsData)
		if parseErr != nil {
			return fmt.Errorf("parse layout rels: %w", parseErr)
		}
		for i := range layoutRels {
			if layoutRels[i].Type == common.RelTypeSlideMaster {
				layoutRels[i].Target = common.MakeRelativePath(clonedLayout, newMaster)
			}
		}
		rendered := renderRelationshipsXML(layoutRels)
		e.parts.Set(common.RelsPathFor(clonedLayout), []byte(rendered))
	}
	return nil
}

func (e *PresentationEditor) writeClonedMaster(
	sourceMaster string,
	newMaster string,
	masterXML []byte,
	masterRels []common.EditorRelationship,
	layoutMap map[string]string,
	newThemePart string,
) {
	e.parts.Set(newMaster, append([]byte(nil), masterXML...))
	for i := range masterRels {
		switch masterRels[i].Type {
		case common.RelTypeSlideLayout:
			oldLayout := common.CanonicalPartPath(path.Join(path.Dir(sourceMaster), masterRels[i].Target))
			if newLayout, mapped := layoutMap[oldLayout]; mapped {
				masterRels[i].Target = common.MakeRelativePath(newMaster, newLayout)
			}
		case common.RelTypeTheme:
			if newThemePart != "" {
				masterRels[i].Target = common.MakeRelativePath(newMaster, newThemePart)
			}
		}
	}
	renderedMasterRels := renderRelationshipsXML(masterRels)
	e.parts.Set(common.RelsPathFor(newMaster), []byte(renderedMasterRels))
}

func (e *PresentationEditor) registerClonedMaster(newMaster string) error {
	e.recalculateNextRelIDNum()
	newMasterRelID := fmt.Sprintf("rId%d", e.nextRelIDNum)
	e.nextRelIDNum++
	e.nonSlideRels = append(e.nonSlideRels, common.EditorRelationship{
		ID:     newMasterRelID,
		Type:   common.RelTypeSlideMaster,
		Target: common.MakeRelativePath(common.PresentationXMLPath, newMaster),
	})

	updatedPresentationXML, err := editorslide.RewritePresentationSlideMasterList(
		[]byte(e.presentationXML),
		newMasterRelID,
	)
	if err != nil {
		return err
	}
	e.presentationXML = updatedPresentationXML
	return nil
}

func (e *PresentationEditor) cloneMasterTheme(
	masterRels []common.EditorRelationship,
	sourceMaster string,
) (string, string, error) {
	for _, rel := range masterRels {
		if rel.Type != common.RelTypeTheme {
			continue
		}
		oldTheme := common.CanonicalPartPath(path.Join(path.Dir(sourceMaster), rel.Target))
		themeXML, ok := e.parts.Get(oldTheme)
		if !ok {
			return "", "", fmt.Errorf("theme part not found: %s", oldTheme)
		}
		newTheme := fmt.Sprintf(
			"ppt/theme/theme%d.xml",
			editorslide.NextPartNumber(
				e.parts.KeysWithPrefix("ppt/theme/"),
				themeNumPattern,
				nextPartPatternSubmatchSize,
			),
		)
		e.parts.Set(newTheme, append([]byte(nil), themeXML...))
		return oldTheme, newTheme, nil
	}
	return "", "", nil
}
