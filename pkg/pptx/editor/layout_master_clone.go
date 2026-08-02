package editor

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strconv"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

var layoutNumPattern = regexp.MustCompile(`^slideLayout(\d+)\.xml$`)
var masterNumPattern = regexp.MustCompile(`^slideMaster(\d+)\.xml$`)
var themeNumPattern = regexp.MustCompile(`^theme(\d+)\.xml$`)

const nextPartPatternSubmatchSize = 2

// CloneLayoutMasterFamily duplicates a layout and its master family inside this
// presentation.
func (e *PresentationEditor) CloneLayoutMasterFamily(layoutPart string) (common.SlideMasterCloneResult, error) {
	return e.cloneLayoutMasterFamilyFrom(e, layoutPart)
}

// ImportLayoutMasterFamily copies a layout and its master family out of another
// presentation into this one, so a slide here can be bound to a layout that
// lives in a different deck. The returned LayoutMap gives the part path each
// source layout now has locally.
func (e *PresentationEditor) ImportLayoutMasterFamily(
	src *PresentationEditor,
	layoutPart string,
) (common.SlideMasterCloneResult, error) {
	if src == nil {
		return common.SlideMasterCloneResult{}, errors.New("source editor cannot be nil")
	}
	return e.cloneLayoutMasterFamilyFrom(src, layoutPart)
}

// cloneLayoutMasterFamilyFrom reads the layout family from src and writes the
// copy into e. When src is e this is an in-deck clone; otherwise it is a
// cross-deck import and any media the layouts reference is copied too.
func (e *PresentationEditor) cloneLayoutMasterFamilyFrom(
	src *PresentationEditor,
	layoutPart string,
) (common.SlideMasterCloneResult, error) {
	sourceMaster, layoutFamily, err := editorslide.CloneFamilyInputs(
		layoutPart,
		src.parts.Has,
		common.CanonicalPartPath,
		func(part string) (string, error) {
			return editorslide.ResolveLayoutMasterPart(part, src.parts.Get, parseRelationshipsXML)
		},
		src.layoutsForMaster,
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
	masterXML, masterRels, err := e.loadMasterCloneSource(src, sourceMaster)
	if err != nil {
		return common.SlideMasterCloneResult{}, err
	}

	themePart, newThemePart, err := e.cloneMasterTheme(src, masterRels, sourceMaster)
	if err != nil {
		return common.SlideMasterCloneResult{}, err
	}
	if err := e.cloneLayoutParts(src, layoutMap, newMaster); err != nil {
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
	src *PresentationEditor,
	sourceMaster string,
) ([]byte, []common.EditorRelationship, error) {
	masterXML, ok := src.parts.Get(sourceMaster)
	if !ok {
		return nil, nil, fmt.Errorf("master part not found: %s", sourceMaster)
	}
	masterRelsPath := common.RelsPathFor(sourceMaster)
	masterRelsData, ok := src.parts.Get(masterRelsPath)
	if !ok {
		return nil, nil, fmt.Errorf("master rels part not found: %s", masterRelsPath)
	}
	masterRels, err := parseRelationshipsXML(masterRelsData)
	if err != nil {
		return nil, nil, fmt.Errorf("parse master rels: %w", err)
	}
	return masterXML, masterRels, nil
}

func (e *PresentationEditor) cloneLayoutParts(
	src *PresentationEditor,
	layoutMap map[string]string,
	newMaster string,
) error {
	for oldLayout, clonedLayout := range layoutMap {
		layoutXML, layoutOK := src.parts.Get(oldLayout)
		if !layoutOK {
			return fmt.Errorf("layout part not found: %s", oldLayout)
		}
		e.parts.Set(clonedLayout, append([]byte(nil), layoutXML...))

		layoutRelsPath := common.RelsPathFor(oldLayout)
		layoutRelsData, relsOK := src.parts.Get(layoutRelsPath)
		if !relsOK {
			return fmt.Errorf("layout rels missing: %s", layoutRelsPath)
		}
		layoutRels, parseErr := parseRelationshipsXML(layoutRelsData)
		if parseErr != nil {
			return fmt.Errorf("parse layout rels: %w", parseErr)
		}
		for i := range layoutRels {
			switch layoutRels[i].Type {
			case common.RelTypeSlideMaster:
				layoutRels[i].Target = common.MakeRelativePath(clonedLayout, newMaster)
			case common.RelTypeImage:
				// A cross-deck import has to bring the image bytes along; an
				// in-deck clone already shares the media part.
				if src == e {
					continue
				}
				srcImage := common.ResolveRelationshipTarget(oldLayout, layoutRels[i].Target)
				newImage, copyErr := e.copyImageAsset(src, srcImage)
				if copyErr != nil {
					return fmt.Errorf("copy layout image %s: %w", srcImage, copyErr)
				}
				layoutRels[i].Target = common.MakeRelativePath(clonedLayout, newImage)
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
	e.parts.Set(newMaster, e.renumberClonedLayoutIDs(append([]byte(nil), masterXML...)))
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

// sldLayoutIDPattern matches the id attribute of a p:sldLayoutId entry.
var sldLayoutIDPattern = regexp.MustCompile(`(<p:sldLayoutId\b[^>]*\bid=")(\d+)(")`)

// sldMasterIDPattern matches the id attribute of a p:sldMasterId entry.
var sldMasterIDPattern = regexp.MustCompile(`<p:sldMasterId\b[^>]*\bid="(\d+)"`)

// minSlideObjectID is the floor ECMA-376 sets for slide master and slide layout
// ids.
const minSlideObjectID uint32 = 2147483648

// renumberClonedLayoutIDs gives a cloned master's p:sldLayoutId entries ids that
// are free across the whole package. PowerPoint draws slide master and slide
// layout ids from one pool, so a straight copy of a master collides with the
// master it came from and Office rejects the package as unreadable — which our
// own validator, and a round trip through this library, both accept.
func (e *PresentationEditor) renumberClonedLayoutIDs(masterXML []byte) []byte {
	used := e.usedSlideObjectIDs()
	next := minSlideObjectID
	return sldLayoutIDPattern.ReplaceAllFunc(masterXML, func(match []byte) []byte {
		groups := sldLayoutIDPattern.FindSubmatch(match)
		assigned := nextFreeSlideObjectID(used, &next)
		return fmt.Appendf(nil, "%s%d%s", groups[1], assigned, groups[3])
	})
}

// nextFreeSlideObjectID hands out the lowest unused id at or above cursor and
// marks it taken.
func nextFreeSlideObjectID(used map[uint32]bool, cursor *uint32) uint32 {
	for used[*cursor] {
		*cursor++
	}
	assigned := *cursor
	used[assigned] = true
	*cursor++
	return assigned
}

// usedSlideObjectIDs collects every p:sldMasterId and p:sldLayoutId already in
// the package, which share one id pool as far as PowerPoint is concerned.
func (e *PresentationEditor) usedSlideObjectIDs() map[uint32]bool {
	used := make(map[uint32]bool)
	collect := func(data []byte, pattern *regexp.Regexp, group int) {
		for _, match := range pattern.FindAllSubmatch(data, -1) {
			value, err := strconv.ParseUint(string(match[group]), 10, 32)
			if err != nil {
				continue
			}
			used[uint32(value)] = true
		}
	}
	for _, masterPart := range e.parts.KeysWithPrefix("ppt/slideMasters/") {
		if masterXML, ok := e.parts.Get(masterPart); ok {
			collect(masterXML, sldLayoutIDPattern, 2)
		}
	}
	collect([]byte(e.presentationXML), sldMasterIDPattern, 1)
	return used
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

	// The master id has to avoid every layout id too, including the ones the
	// clone just took, or PowerPoint refuses to open the package.
	used := e.usedSlideObjectIDs()
	cursor := minSlideObjectID
	masterID := nextFreeSlideObjectID(used, &cursor)

	updatedPresentationXML, err := editorslide.RewritePresentationSlideMasterListWithID(
		[]byte(e.presentationXML),
		newMasterRelID,
		int64(masterID),
	)
	if err != nil {
		return err
	}
	e.presentationXML = updatedPresentationXML
	return nil
}

func (e *PresentationEditor) cloneMasterTheme(
	src *PresentationEditor,
	masterRels []common.EditorRelationship,
	sourceMaster string,
) (string, string, error) {
	for _, rel := range masterRels {
		if rel.Type != common.RelTypeTheme {
			continue
		}
		oldTheme := common.CanonicalPartPath(path.Join(path.Dir(sourceMaster), rel.Target))
		themeXML, ok := src.parts.Get(oldTheme)
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
