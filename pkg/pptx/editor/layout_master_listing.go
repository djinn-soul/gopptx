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

// Master/layout XML scanning patterns. Fixed literals, compiled once at init.
var (
	shapeNamePattern          = regexp.MustCompile(`(?s)<p:nvSpPr>.*?<p:cNvPr[^>]*name="([^"]+)"`)
	placeholderElementPattern = regexp.MustCompile(`(?s)<p:ph[^>]*>`)
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
			Shapes:       shapes,
			Placeholders: placeholders,
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
		// Get shapes and placeholders from the layout
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

// GetMasterShapes returns the shapes in a slide master.
func (e *PresentationEditor) GetMasterShapes(masterPart string) []string {
	xmlData, ok := e.parts.Get(masterPart)
	if !ok {
		return nil
	}
	return parseShapesFromMasterLayoutXML(xmlData)
}

// GetMasterPlaceholders returns the placeholders in a slide master.
func (e *PresentationEditor) GetMasterPlaceholders(masterPart string) []common.PlaceholderInfo {
	xmlData, ok := e.parts.Get(masterPart)
	if !ok {
		return nil
	}
	return parsePlaceholdersFromMasterLayoutXML(xmlData)
}

// GetLayoutShapes returns the shapes in a slide layout.
func (e *PresentationEditor) GetLayoutShapes(layoutPart string) []string {
	xmlData, ok := e.parts.Get(layoutPart)
	if !ok {
		return nil
	}
	return parseShapesFromMasterLayoutXML(xmlData)
}

// GetLayoutPlaceholders returns the placeholders in a slide layout.
func (e *PresentationEditor) GetLayoutPlaceholders(layoutPart string) []common.PlaceholderInfo {
	xmlData, ok := e.parts.Get(layoutPart)
	if !ok {
		return nil
	}
	return parsePlaceholdersFromMasterLayoutXML(xmlData)
}

// parseShapesFromMasterLayoutXML extracts shape names from master/layout XML.
func parseShapesFromMasterLayoutXML(content []byte) []string {
	// Look for non-placeholder shapes (spTree elements that are not ph elements)
	var shapes []string

	// Simple pattern: find <p:sp> elements and extract name from <p:nvSpPr><p:cNvPr>
	// This is a simplified implementation
	matches := shapeNamePattern.FindAllStringSubmatch(string(content), -1)
	for _, m := range matches {
		if len(m) > 1 {
			shapes = append(shapes, m[1])
		}
	}
	return shapes
}

// parsePlaceholdersFromMasterLayoutXML extracts placeholder info from master/layout XML.
func parsePlaceholdersFromMasterLayoutXML(content []byte) []common.PlaceholderInfo {
	var placeholders []common.PlaceholderInfo

	// Find all placeholder elements <p:ph>
	matches := placeholderElementPattern.FindAllString(string(content), -1)
	for _, m := range matches {
		placeholders = append(placeholders, common.PlaceholderInfo{
			Type:  parsePlaceholderAttrString(m, phTypeAttrPattern),
			Index: parsePlaceholderAttrIndex(m),
			Name:  parsePlaceholderAttrString(m, phNameAttrPattern),
		})
	}

	return placeholders
}
