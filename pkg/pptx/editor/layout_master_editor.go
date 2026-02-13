package editor

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	layoutNumPattern = regexp.MustCompile(`^slideLayout(\d+)\.xml$`)
	masterNumPattern = regexp.MustCompile(`^slideMaster(\d+)\.xml$`)
	themeNumPattern  = regexp.MustCompile(`^theme(\d+)\.xml$`)
)

func (e *PresentationEditor) ListSlideLayouts() ([]common.SlideLayoutInfo, error) {
	layoutParts := e.parts.KeysWithPrefix("ppt/slideLayouts/slideLayout")
	infos := make([]common.SlideLayoutInfo, 0, len(layoutParts))
	for _, part := range layoutParts {
		if !strings.HasSuffix(part, ".xml") {
			continue
		}
		masterPart, err := e.layoutMasterPart(part)
		if err != nil {
			return nil, err
		}
		xmlData, ok := e.parts.Get(part)
		if !ok {
			return nil, fmt.Errorf("layout part not found: %s", part)
		}
		infos = append(infos, common.SlideLayoutInfo{
			Part:       part,
			Name:       parseLayoutName(xmlData),
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
	rendered, err := renderRelationshipsXML(rels)
	if err != nil {
		return err
	}
	e.parts.Set(relsPath, []byte(rendered))
	return nil
}

func (e *PresentationEditor) CloneLayoutMasterFamily(layoutPart string) (common.SlideMasterCloneResult, error) {
	layoutPart = common.CanonicalPartPath(layoutPart)
	if !e.parts.Has(layoutPart) {
		return common.SlideMasterCloneResult{}, fmt.Errorf("layout part %s not found", layoutPart)
	}
	sourceMaster, err := e.layoutMasterPart(layoutPart)
	if err != nil {
		return common.SlideMasterCloneResult{}, err
	}

	layoutFamily, err := e.layoutsForMaster(sourceMaster)
	if err != nil {
		return common.SlideMasterCloneResult{}, err
	}
	if len(layoutFamily) == 0 {
		return common.SlideMasterCloneResult{}, fmt.Errorf("no layouts found for master %s", sourceMaster)
	}

	newMaster := fmt.Sprintf("ppt/slideMasters/slideMaster%d.xml", e.nextPartNumber(masterNumPattern, "ppt/slideMasters"))
	layoutMap := make(map[string]string, len(layoutFamily))
	nextLayoutNum := e.nextPartNumber(layoutNumPattern, "ppt/slideLayouts")
	for _, oldLayout := range layoutFamily {
		layoutMap[oldLayout] = fmt.Sprintf("ppt/slideLayouts/slideLayout%d.xml", nextLayoutNum)
		nextLayoutNum++
	}

	masterXML, ok := e.parts.Get(sourceMaster)
	if !ok {
		return common.SlideMasterCloneResult{}, fmt.Errorf("master part not found: %s", sourceMaster)
	}
	masterRelsPath := common.RelsPathFor(sourceMaster)
	masterRelsData, ok := e.parts.Get(masterRelsPath)
	if !ok {
		return common.SlideMasterCloneResult{}, fmt.Errorf("master rels part not found: %s", masterRelsPath)
	}
	masterRels, err := parseRelationshipsXML(masterRelsData)
	if err != nil {
		return common.SlideMasterCloneResult{}, fmt.Errorf("parse master rels: %w", err)
	}

	themePart, newThemePart, err := e.cloneMasterTheme(masterRels, sourceMaster)
	if err != nil {
		return common.SlideMasterCloneResult{}, err
	}

	// clone layouts first
	for oldLayout, clonedLayout := range layoutMap {
		layoutXML, ok := e.parts.Get(oldLayout)
		if !ok {
			return common.SlideMasterCloneResult{}, fmt.Errorf("layout part not found: %s", oldLayout)
		}
		e.parts.Set(clonedLayout, append([]byte(nil), layoutXML...))

		layoutRelsPath := common.RelsPathFor(oldLayout)
		layoutRelsData, ok := e.parts.Get(layoutRelsPath)
		if !ok {
			return common.SlideMasterCloneResult{}, fmt.Errorf("layout rels missing: %s", layoutRelsPath)
		}
		layoutRels, err := parseRelationshipsXML(layoutRelsData)
		if err != nil {
			return common.SlideMasterCloneResult{}, fmt.Errorf("parse layout rels: %w", err)
		}
		for i := range layoutRels {
			if layoutRels[i].Type == common.RelTypeSlideMaster {
				layoutRels[i].Target = common.MakeRelativePath(clonedLayout, newMaster)
			}
		}
		rendered, err := renderRelationshipsXML(layoutRels)
		if err != nil {
			return common.SlideMasterCloneResult{}, err
		}
		e.parts.Set(common.RelsPathFor(clonedLayout), []byte(rendered))
	}

	// clone master xml and rels
	e.parts.Set(newMaster, append([]byte(nil), masterXML...))
	for i := range masterRels {
		if masterRels[i].Type == common.RelTypeSlideLayout {
			oldLayout := common.CanonicalPartPath(path.Join(path.Dir(sourceMaster), masterRels[i].Target))
			if newLayout, ok := layoutMap[oldLayout]; ok {
				masterRels[i].Target = common.MakeRelativePath(newMaster, newLayout)
			}
		}
		if masterRels[i].Type == common.RelTypeTheme && newThemePart != "" {
			masterRels[i].Target = common.MakeRelativePath(newMaster, newThemePart)
		}
	}
	renderedMasterRels, err := renderRelationshipsXML(masterRels)
	if err != nil {
		return common.SlideMasterCloneResult{}, err
	}
	e.parts.Set(common.RelsPathFor(newMaster), []byte(renderedMasterRels))

	// register master in presentation rels + XML master list
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
		return common.SlideMasterCloneResult{}, err
	}
	e.presentationXML = updatedPresentationXML

	outTheme := themePart
	if newThemePart != "" {
		outTheme = newThemePart
	}
	return common.SlideMasterCloneResult{
		MasterPart: newMaster,
		ThemePart:  outTheme,
		LayoutMap:  layoutMap,
	}, nil
}

func (e *PresentationEditor) layoutMasterPart(layoutPart string) (string, error) {
	relsPath := common.RelsPathFor(layoutPart)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return "", fmt.Errorf("layout rels part not found: %s", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return "", fmt.Errorf("parse layout rels: %w", err)
	}
	for _, rel := range rels {
		if rel.Type == common.RelTypeSlideMaster {
			return common.CanonicalPartPath(path.Join(path.Dir(layoutPart), rel.Target)), nil
		}
	}
	return "", fmt.Errorf("layout %s has no slideMaster relationship", layoutPart)
}

func (e *PresentationEditor) layoutsForMaster(masterPart string) ([]string, error) {
	relsParts := e.parts.KeysWithPrefix("ppt/slideLayouts/_rels/slideLayout")
	out := make([]string, 0, len(relsParts))
	for _, relsPath := range relsParts {
		if !strings.HasSuffix(relsPath, ".xml.rels") {
			continue
		}
		layoutName := strings.TrimSuffix(path.Base(relsPath), ".rels")
		layoutPart := path.Join("ppt/slideLayouts", layoutName)
		relsData, ok := e.parts.Get(relsPath)
		if !ok {
			continue
		}
		rels, err := parseRelationshipsXML(relsData)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", relsPath, err)
		}
		for _, rel := range rels {
			if rel.Type != common.RelTypeSlideMaster {
				continue
			}
			target := common.CanonicalPartPath(path.Join(path.Dir(layoutPart), rel.Target))
			if target == masterPart {
				out = append(out, layoutPart)
				break
			}
		}
	}
	sort.Strings(out)
	return out, nil
}

func (e *PresentationEditor) cloneMasterTheme(masterRels []common.EditorRelationship, sourceMaster string) (oldTheme string, newTheme string, err error) {
	for _, rel := range masterRels {
		if rel.Type != common.RelTypeTheme {
			continue
		}
		oldTheme = common.CanonicalPartPath(path.Join(path.Dir(sourceMaster), rel.Target))
		themeXML, ok := e.parts.Get(oldTheme)
		if !ok {
			return "", "", fmt.Errorf("theme part not found: %s", oldTheme)
		}
		newTheme = fmt.Sprintf("ppt/theme/theme%d.xml", e.nextPartNumber(themeNumPattern, "ppt/theme"))
		e.parts.Set(newTheme, append([]byte(nil), themeXML...))
		return oldTheme, newTheme, nil
	}
	return "", "", nil
}

func (e *PresentationEditor) nextPartNumber(pattern *regexp.Regexp, dir string) int {
	maxNum := 0
	for _, part := range e.parts.KeysWithPrefix(dir + "/") {
		base := path.Base(part)
		m := pattern.FindStringSubmatch(base)
		if len(m) != 2 {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if n > maxNum {
			maxNum = n
		}
	}
	return maxNum + 1
}

func parseLayoutName(layoutXML []byte) string {
	s := string(layoutXML)
	const marker = `name="`
	pos := strings.Index(s, marker)
	if pos < 0 {
		return ""
	}
	start := pos + len(marker)
	end := strings.Index(s[start:], `"`)
	if end < 0 {
		return ""
	}
	return s[start : start+end]
}
