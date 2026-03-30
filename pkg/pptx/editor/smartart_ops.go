package editor

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

var reSmartArtRelIDs = regexp.MustCompile(`r:(dm|lo|qs|cs)=["']([^"']+)["']`)
var reSmartArtRelTag = regexp.MustCompile(`<dgm:relIds\b[^>]*>`)

// smartArtPartRefs holds all resolved part paths and relationship IDs for a SmartArt diagram.
type smartArtPartRefs struct {
	DataPath    string
	LayoutPath  string
	StylePath   string
	ColorPath   string
	DrawingPath string

	DataRelID    string
	LayoutRelID  string
	StyleRelID   string
	ColorRelID   string
	DrawingRelID string
}

// resolveSmartArtParts finds all 5 diagram parts for the SmartArt graphic frame at shapeID.
func (e *PresentationEditor) resolveSmartArtParts(
	slideRef common.EditorSlideRef, shapeID int,
) (*smartArtPartRefs, error) {
	slideXML, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return nil, fmt.Errorf("slide part %q not found", slideRef.Part)
	}

	dmRelID, loRelID, qsRelID, csRelID := extractAllSmartArtRelIDs(string(slideXML), shapeID)
	if dmRelID == "" {
		return nil, fmt.Errorf("shapeID %d not found or is not a SmartArt graphic frame", shapeID)
	}

	relsPath := common.RelsPathFor(slideRef.Part)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return nil, fmt.Errorf("slide rels %q not found", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return nil, fmt.Errorf("parse rels: %w", err)
	}

	refs := &smartArtPartRefs{
		DataRelID:   dmRelID,
		LayoutRelID: loRelID,
		StyleRelID:  qsRelID,
		ColorRelID:  csRelID,
	}

	for _, rel := range rels {
		target := common.ResolveRelationshipTarget(slideRef.Part, rel.Target)
		switch rel.ID {
		case dmRelID:
			refs.DataPath = target
		case loRelID:
			refs.LayoutPath = target
		case qsRelID:
			refs.StylePath = target
		case csRelID:
			refs.ColorPath = target
		}
		if rel.Type == relTypeDiagramDrawing && refs.DrawingPath == "" {
			refs.DrawingPath = target
			refs.DrawingRelID = rel.ID
		}
	}

	if refs.DataPath == "" {
		return nil, fmt.Errorf("could not resolve SmartArt data part for shapeID %d", shapeID)
	}
	return refs, nil
}

// extractAllSmartArtRelIDs extracts r:dm, r:lo, r:qs, r:cs from dgm:relIds for the given shapeID.
func extractAllSmartArtRelIDs(slideXML string, shapeID int) (string, string, string, string) {
	idPattern := regexp.MustCompile(fmt.Sprintf(`<p:cNvPr\b[^>]*\bid=["']%d["']`, shapeID))
	searchFrom := 0
	for {
		frameStartRel := strings.Index(slideXML[searchFrom:], "<p:graphicFrame")
		if frameStartRel < 0 {
			return "", "", "", ""
		}
		frameStart := searchFrom + frameStartRel

		frameEndRel := strings.Index(slideXML[frameStart:], "</p:graphicFrame>")
		if frameEndRel < 0 {
			return "", "", "", ""
		}
		frameEnd := frameStart + frameEndRel + len("</p:graphicFrame>")
		frameFragment := slideXML[frameStart:frameEnd]
		searchFrom = frameEnd

		if !idPattern.MatchString(frameFragment) {
			continue
		}

		relTag := reSmartArtRelTag.FindString(frameFragment)
		if relTag == "" {
			return "", "", "", ""
		}

		var dm, lo, qs, cs string
		for _, m := range reSmartArtRelIDs.FindAllStringSubmatch(relTag, -1) {
			switch m[1] {
			case "dm":
				dm = m[2]
			case "lo":
				lo = m[2]
			case "qs":
				qs = m[2]
			case "cs":
				cs = m[2]
			}
		}
		return dm, lo, qs, cs
	}
}

// removeSlideRelationships removes relationships by ID from a slide's .rels file.
func (e *PresentationEditor) removeSlideRelationships(slidePart string, relIDs []string) error {
	relsPath := common.RelsPathFor(slidePart)
	data, ok := e.parts.Get(relsPath)
	if !ok {
		return nil // nothing to remove
	}
	rels, err := parseRelationshipsXML(data)
	if err != nil {
		return fmt.Errorf("parse %s: %w", relsPath, err)
	}
	remove := make(map[string]bool, len(relIDs))
	for _, id := range relIDs {
		remove[id] = true
	}
	filtered := rels[:0]
	for _, rel := range rels {
		if !remove[rel.ID] {
			filtered = append(filtered, rel)
		}
	}
	return e.writeRelationships(relsPath, filtered)
}

// removeContentTypeOverride removes a content-type override entry from [Content_Types].xml.
func (e *PresentationEditor) removeContentTypeOverride(partPath string) {
	ctPath := "[Content_Types].xml"
	data, ok := e.parts.Get(ctPath)
	if !ok {
		return
	}
	partNameRooted := "/" + partPath
	// Find and remove the Override element for this part.
	re := regexp.MustCompile(`<Override\s+PartName="` + regexp.QuoteMeta(partNameRooted) + `"[^>]*/>`)
	updated := re.ReplaceAll(data, nil)
	if len(updated) != len(data) {
		e.parts.Set(ctPath, updated)
	}
}

// DeleteSmartArt removes a SmartArt graphic frame and all 5 associated diagram parts.
func (e *PresentationEditor) DeleteSmartArt(slideIndex, shapeID int) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]

	refs, err := e.resolveSmartArtParts(slideRef, shapeID)
	if err != nil {
		return err
	}

	// Delete all diagram parts.
	for _, path := range []string{refs.DataPath, refs.LayoutPath, refs.StylePath, refs.ColorPath, refs.DrawingPath} {
		if path != "" {
			e.parts.Delete(path)
			e.removeContentTypeOverride(path)
		}
	}

	// Remove relationships from slide .rels.
	relIDs := []string{refs.DataRelID, refs.LayoutRelID, refs.StyleRelID, refs.ColorRelID}
	if refs.DrawingRelID != "" {
		relIDs = append(relIDs, refs.DrawingRelID)
	}
	if err := e.removeSlideRelationships(slideRef.Part, relIDs); err != nil {
		return fmt.Errorf("remove SmartArt rels: %w", err)
	}

	// Remove the graphic frame from slide XML.
	return e.RemoveShape(slideIndex, shapeID)
}

// ChangeSmartArtLayout replaces the layout of an existing SmartArt while preserving its nodes.
func (e *PresentationEditor) ChangeSmartArtLayout(slideIndex, shapeID int, newLayout smartart.Layout) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]

	refs, err := e.resolveSmartArtParts(slideRef, shapeID)
	if err != nil {
		return err
	}

	dataXML, ok := e.parts.Get(refs.DataPath)
	if !ok {
		return fmt.Errorf("SmartArt data part %q not found", refs.DataPath)
	}

	nodes, err := smartart.ParseDataModelNodes(dataXML)
	if err != nil {
		return fmt.Errorf("parse SmartArt nodes: %w", err)
	}

	sa := smartart.NewSmartArt(newLayout)
	sa.Nodes = nodes
	if quickStyle := e.readSmartArtUniqueID(refs.StylePath); quickStyle != "" {
		sa = sa.WithQuickStyle(quickStyle)
	}
	if colorStyle := e.readSmartArtUniqueID(refs.ColorPath); colorStyle != "" {
		sa = sa.WithColorStyle(colorStyle)
	}
	if err := sa.Validate(slideIndex); err != nil {
		return err
	}

	spec := sa.ToSpec()

	category := smartArtCategoryFromURI(spec.LayoutURI)
	e.parts.Set(refs.DataPath, []byte(pptxxml.SmartArtDataXML(spec)))
	e.parts.Set(refs.LayoutPath, []byte(pptxxml.SmartArtLayoutXML(spec.LayoutURI, category)))
	e.parts.Set(refs.StylePath, []byte(pptxxml.SmartArtStyleXML(spec.QuickStyleID)))
	e.parts.Set(refs.ColorPath, []byte(pptxxml.SmartArtColorsXML(spec.ColorStyleID)))
	e.parts.Set(refs.DrawingPath, []byte(pptxxml.SmartArtDrawingXML(spec)))
	return nil
}

// SetSmartArtStyle updates the quick style and/or color style of an existing SmartArt.
// Pass empty string to keep the existing value.
func (e *PresentationEditor) SetSmartArtStyle(slideIndex, shapeID int, quickStyle, colorStyle string) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]

	refs, err := e.resolveSmartArtParts(slideRef, shapeID)
	if err != nil {
		return err
	}

	if quickStyle != "" {
		e.parts.Set(refs.StylePath, []byte(pptxxml.SmartArtStyleXML(quickStyle)))
	}
	if colorStyle != "" {
		e.parts.Set(refs.ColorPath, []byte(pptxxml.SmartArtColorsXML(colorStyle)))
	}
	return nil
}

// SetSmartArtNodes replaces the node tree of an existing SmartArt diagram.
func (e *PresentationEditor) SetSmartArtNodes(slideIndex, shapeID int, nodes []smartart.Node) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]

	refs, err := e.resolveSmartArtParts(slideRef, shapeID)
	if err != nil {
		return err
	}

	// Read existing layout URI from layout part.
	layoutXML, ok := e.parts.Get(refs.LayoutPath)
	if !ok {
		return fmt.Errorf("SmartArt layout part %q not found", refs.LayoutPath)
	}
	layoutURI := smartart.ExtractLayoutURI(string(layoutXML))
	if layoutURI == "" {
		return fmt.Errorf("could not determine SmartArt layout URI from %q", refs.LayoutPath)
	}

	// Build a SmartArt with the new nodes and validate.
	sa := smartart.NewSmartArt(smartart.CustomLayout(layoutURI))
	sa.Nodes = nodes
	if err := sa.Validate(slideIndex); err != nil {
		return err
	}

	spec := sa.ToSpec()
	e.parts.Set(refs.DataPath, []byte(pptxxml.SmartArtDataXML(spec)))
	e.parts.Set(refs.DrawingPath, []byte(pptxxml.SmartArtDrawingXML(spec)))
	return nil
}

func (e *PresentationEditor) readSmartArtUniqueID(partPath string) string {
	if partPath == "" {
		return ""
	}
	partXML, ok := e.parts.Get(partPath)
	if !ok {
		return ""
	}
	return smartart.ExtractUniqueID(string(partXML))
}
