package editor

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

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

		// Get shapes and placeholders from the master
		shapes := e.GetMasterShapes(masterPart)
		placeholders := e.GetMasterPlaceholders(masterPart)

		infos = append(infos, common.SlideMasterInfo{
			Part:         masterPart,
			Shapes:       shapes,
			Placeholders: placeholders,
		})
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
	layoutIDs := e.layoutIDsForMaster(masterPart)
	infos := make([]common.SlideLayoutInfo, 0, len(layouts))
	for _, part := range layouts {
		xmlData, ok := e.parts.Get(part)
		if !ok {
			return nil, fmt.Errorf("layout part not found: %s", part)
		}
		// Get shapes and placeholders from the layout
		shapes := e.GetLayoutShapes(part)
		placeholders := e.GetLayoutPlaceholders(part)

		infos = append(infos, common.SlideLayoutInfo{
			Part:         part,
			Name:         editorslide.ParseLayoutName(xmlData),
			MasterPart:   masterPart,
			LayoutID:     layoutIDs[part],
			Shapes:       shapes,
			Placeholders: placeholders,
		})
	}
	return infos, nil
}

func (e *PresentationEditor) ListSlideLayouts() ([]common.SlideLayoutInfo, error) {
	masters, err := e.ListSlideMasters()
	if err != nil || len(masters) == 0 {
		layoutParts := e.parts.KeysWithPrefix("ppt/slideLayouts/slideLayout")
		infos := make([]common.SlideLayoutInfo, 0, len(layoutParts))
		for _, part := range layoutParts {
			if !strings.HasSuffix(part, ".xml") {
				continue
			}
			masterPart, _ := editorslide.ResolveLayoutMasterPart(part, e.parts.Get, parseRelationshipsXML)
			xmlData, ok := e.parts.Get(part)
			if !ok {
				continue
			}
			shapes := e.GetLayoutShapes(part)
			placeholders := e.GetLayoutPlaceholders(part)
			infos = append(infos, common.SlideLayoutInfo{
				Part:         part,
				Name:         editorslide.ParseLayoutName(xmlData),
				MasterPart:   masterPart,
				Shapes:       shapes,
				Placeholders: placeholders,
			})
		}
		sort.Slice(infos, func(i, j int) bool { return infos[i].Part < infos[j].Part })
		return infos, nil
	}

	var allInfos []common.SlideLayoutInfo
	for _, master := range masters {
		mLayouts, err := e.ListMasterLayouts(master.Part)
		if err == nil {
			allInfos = append(allInfos, mLayouts...)
		}
	}
	return allInfos, nil
}

func (e *PresentationEditor) GetSlideLayoutRef(slideIndex int) (string, string, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return "", "", fmt.Errorf("slide index %d out of range", slideIndex)
	}

	slidePart := e.slides[slideIndex].Part
	relsPath := common.RelsPathFor(slidePart)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return "", "", fmt.Errorf("slide rels part not found: %s", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return "", "", fmt.Errorf("parse slide rels: %w", err)
	}

	var layoutPart string
	for _, rel := range rels {
		if rel.Type != common.RelTypeSlideLayout {
			continue
		}
		layoutPart = common.ResolveRelationshipTarget(slidePart, rel.Target)
		break
	}
	if layoutPart == "" {
		return "", "", fmt.Errorf("slide %d has no layout relationship", slideIndex)
	}

	masterPart, err := editorslide.ResolveLayoutMasterPart(layoutPart, e.parts.Get, parseRelationshipsXML)
	if err != nil {
		return "", "", err
	}
	return layoutPart, masterPart, nil
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
	ridToPart := make(map[string]string)
	for _, rel := range rels {
		if rel.Type == common.RelTypeSlideLayout {
			layoutPart := common.CanonicalPartPath(path.Join(path.Dir(masterPart), rel.Target))
			ridToPart[rel.ID] = layoutPart
		}
	}

	masterXML, _ := e.parts.Get(masterPart)
	content := string(masterXML)
	matches := slideLayoutIDPattern.FindAllStringSubmatch(content, -1)

	out := make([]string, 0, len(matches))
	seen := make(map[string]bool)
	for _, m := range matches {
		if len(m) > 2 {
			rID := m[2]
			if part, ok := ridToPart[rID]; ok && e.parts.Has(part) {
				out = append(out, part)
				seen[part] = true
			}
		}
	}
	for _, part := range ridToPart {
		if !seen[part] && e.parts.Has(part) {
			out = append(out, part)
		}
	}
	return out, nil
}

// slideLayoutIDPattern captures the sldLayoutId id and its relationship id. The
// id is what SlideMaster.get_layout(slide_layout_id) looks a layout up by.
var slideLayoutIDPattern = regexp.MustCompile(
	`<p:sldLayoutId[^>]*\bid="(\d+)"[^>]*r:id="([^"]+)"`,
)

// slideLayoutIDGroups is the full match plus the id and r:id capture groups.
const slideLayoutIDGroups = 3

// layoutIDsForMaster maps each layout part to the sldLayoutId id the master
// lists it under. Parts not referenced from sldLayoutIdLst are absent.
func (e *PresentationEditor) layoutIDsForMaster(masterPart string) map[string]int {
	ids := make(map[string]int)
	masterRelsData, ok := e.parts.Get(common.RelsPathFor(masterPart))
	if !ok {
		return ids
	}
	rels, err := parseRelationshipsXML(masterRelsData)
	if err != nil {
		return ids
	}
	ridToPart := make(map[string]string, len(rels))
	for _, rel := range rels {
		if rel.Type == common.RelTypeSlideLayout {
			ridToPart[rel.ID] = common.CanonicalPartPath(
				path.Join(path.Dir(masterPart), rel.Target),
			)
		}
	}

	masterXML, _ := e.parts.Get(masterPart)
	for _, m := range slideLayoutIDPattern.FindAllStringSubmatch(string(masterXML), -1) {
		if len(m) < slideLayoutIDGroups {
			continue
		}
		layoutID, convErr := strconv.Atoi(m[1])
		if convErr != nil {
			continue
		}
		if part, found := ridToPart[m[2]]; found {
			ids[part] = layoutID
		}
	}
	return ids
}
