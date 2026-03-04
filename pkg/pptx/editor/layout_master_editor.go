package editor

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

var layoutNumPattern = regexp.MustCompile(`^slideLayout(\d+)\.xml$`)
var masterNumPattern = regexp.MustCompile(`^slideMaster(\d+)\.xml$`)
var themeNumPattern = regexp.MustCompile(`^theme(\d+)\.xml$`)

func (e *PresentationEditor) ListSlideMasters() ([]common.SlideMasterInfo, error) {
	infos := make([]common.SlideMasterInfo, 0, len(e.nonSlideRels))
	seen := make(map[string]struct{}, len(e.nonSlideRels))
	for _, rel := range e.nonSlideRels {
		if rel.Type != common.RelTypeSlideMaster {
			continue
		}
		masterPart := common.CanonicalPartPath(path.Join(path.Dir(common.PresentationXMLPath), rel.Target))
		if _, ok := seen[masterPart]; ok {
			continue
		}
		if !e.parts.Has(masterPart) {
			return nil, fmt.Errorf("slide master part not found: %s", masterPart)
		}
		seen[masterPart] = struct{}{}
		infos = append(infos, common.SlideMasterInfo{Part: masterPart})
	}
	return infos, nil
}

func (e *PresentationEditor) ListMasterLayouts(masterPart string) ([]common.SlideLayoutInfo, error) {
	masterPart = common.CanonicalPartPath(masterPart)
	if !e.parts.Has(masterPart) {
		return nil, fmt.Errorf("master part not found: %s", masterPart)
	}
	layouts, err := e.layoutsForMaster(masterPart)
	if err != nil {
		return nil, err
	}
	infos := make([]common.SlideLayoutInfo, 0, len(layouts))
	for _, part := range layouts {
		xmlData, ok := e.parts.Get(part)
		if !ok {
			return nil, fmt.Errorf("layout part not found: %s", part)
		}
		infos = append(infos, common.SlideLayoutInfo{
			Part:       part,
			Name:       editorslide.ParseLayoutName(xmlData),
			MasterPart: masterPart,
		})
	}
	return infos, nil
}

func (e *PresentationEditor) ListSlideLayouts() ([]common.SlideLayoutInfo, error) {
	layoutParts := e.parts.KeysWithPrefix("ppt/slideLayouts/slideLayout")
	infos := make([]common.SlideLayoutInfo, 0, len(layoutParts))
	for _, part := range layoutParts {
		if !strings.HasSuffix(part, ".xml") {
			continue
		}
		masterPart, err := editorslide.ResolveLayoutMasterPart(part, e.parts.Get, parseRelationshipsXML)
		if err != nil {
			return nil, err
		}
		xmlData, ok := e.parts.Get(part)
		if !ok {
			return nil, fmt.Errorf("layout part not found: %s", part)
		}
		infos = append(infos, common.SlideLayoutInfo{
			Part:       part,
			Name:       editorslide.ParseLayoutName(xmlData),
			MasterPart: masterPart,
		})
	}
	sort.Slice(infos, func(i, j int) bool { return infos[i].Part < infos[j].Part })
	return infos, nil
}

func (e *PresentationEditor) RebindSlideLayout(slideIndex int, layoutPart string) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	layoutPart = common.CanonicalPartPath(layoutPart)
	if !e.parts.Has(layoutPart) {
		return fmt.Errorf("layout part %s not found", layoutPart)
	}

	slidePart := e.slides[slideIndex].Part
	relsPath := common.RelsPathFor(slidePart)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return fmt.Errorf("slide rels part not found: %s", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return fmt.Errorf("parse slide rels: %w", err)
	}

	found := false
	for i := range rels {
		if rels[i].Type != common.RelTypeSlideLayout {
			continue
		}
		rels[i].Target = common.MakeRelativePath(slidePart, layoutPart)
		found = true
		break
	}
	if !found {
		return fmt.Errorf("slide %d has no layout relationship", slideIndex)
	}
	rendered := renderRelationshipsXML(rels)
	e.parts.Set(relsPath, []byte(rendered))
	return nil
}

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
		editorslide.NextPartNumber(e.parts.KeysWithPrefix("ppt/slideMasters/"), masterNumPattern, 2),
	)
	layoutMap := editorslide.BuildLayoutCloneMap(
		layoutFamily,
		editorslide.NextPartNumber(e.parts.KeysWithPrefix("ppt/slideLayouts/"), layoutNumPattern, 2),
	)
	masterXML, masterRels, err := e.loadMasterCloneSource(sourceMaster)
	if err != nil {
		return common.SlideMasterCloneResult{}, err
	}

	themePart, newThemePart, err := e.cloneMasterTheme(masterRels, sourceMaster)
	if err != nil {
		return common.SlideMasterCloneResult{}, err
	}
	cloneLayoutErr := e.cloneLayoutParts(layoutMap, newMaster)
	if cloneLayoutErr != nil {
		return common.SlideMasterCloneResult{}, cloneLayoutErr
	}
	e.writeClonedMaster(sourceMaster, newMaster, masterXML, masterRels, layoutMap, newThemePart)
	registerErr := e.registerClonedMaster(newMaster)
	if registerErr != nil {
		return common.SlideMasterCloneResult{}, registerErr
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

	updatedPresentationXML, err := rewritePresentationSlideMasterList([]byte(e.presentationXML), newMasterRelID)
	if err != nil {
		return err
	}
	e.presentationXML = updatedPresentationXML
	return nil
}

func (e *PresentationEditor) layoutsForMaster(masterPart string) ([]string, error) {
	masterRelsPath := common.RelsPathFor(masterPart)
	masterRelsData, ok := e.parts.Get(masterRelsPath)
	if !ok {
		return nil, fmt.Errorf("master rels part not found: %s", masterRelsPath)
	}
	rels, err := parseRelationshipsXML(masterRelsData)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", masterRelsPath, err)
	}
	out := make([]string, 0, len(rels))
	for _, rel := range rels {
		if rel.Type != common.RelTypeSlideLayout {
			continue
		}
		layoutPart := common.CanonicalPartPath(path.Join(path.Dir(masterPart), rel.Target))
		if !e.parts.Has(layoutPart) {
			return nil, fmt.Errorf("layout part not found: %s", layoutPart)
		}
		out = append(out, layoutPart)
	}
	return out, nil
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
			editorslide.NextPartNumber(e.parts.KeysWithPrefix("ppt/theme/"), themeNumPattern, 2),
		)
		e.parts.Set(newTheme, append([]byte(nil), themeXML...))
		return oldTheme, newTheme, nil
	}
	return "", "", nil
}
