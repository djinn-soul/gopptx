package editor

import (
	"fmt"
	"path"
	"sort"
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
