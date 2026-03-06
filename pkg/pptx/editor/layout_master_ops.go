package editor

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

// AddSlideMaster adds a new slide master to the presentation.
// Returns the path to the newly created master.
func (e *PresentationEditor) AddSlideMaster() (string, error) {
	// Find the next available master number
	masterNum := editorslide.NextPartNumber(
		e.parts.KeysWithPrefix("ppt/slideMasters/"),
		masterNumPattern,
		nextPartPatternSubmatchSize,
	)
	masterPart := fmt.Sprintf("ppt/slideMasters/slideMaster%d.xml", masterNum)

	// Create a basic master XML
	masterXML := []byte(editorslide.DefaultSlideMaster())
	e.parts.Set(masterPart, masterXML)

	// Create master relationships
	masterRelsPath := fmt.Sprintf("ppt/slideMasters/_rels/slideMaster%d.xml.rels", masterNum)
	masterRels := editorslide.DefaultSlideMasterRelationships()
	e.parts.Set(masterRelsPath, []byte(masterRels))

	// Register the master in Content_Types.xml
	e.addContentTypeOverride(masterPart, "application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml")

	// Register the master in presentation.xml
	if err := e.registerNewMaster(masterPart); err != nil {
		return "", err
	}

	return masterPart, nil
}

// RemoveSlideMaster removes a slide master and its associated layouts.
func (e *PresentationEditor) RemoveSlideMaster(masterPart string) error {
	masterPart = common.CanonicalPartPath(masterPart)
	if !e.parts.Has(masterPart) {
		return fmt.Errorf("master part not found: %s", masterPart)
	}

	// Prevent removing the last slide master
	masters, err := e.ListSlideMasters()
	if err != nil {
		return err
	}
	if len(masters) <= 1 {
		return errors.New("cannot remove the last slide master")
	}

	// Get all layouts for this master
	layouts, err := e.layoutsForMaster(masterPart)
	if err != nil {
		return err
	}

	// Remove all layouts first
	for _, layoutPart := range layouts {
		if err := e.RemoveSlideLayout(layoutPart); err != nil {
			return err
		}
	}

	// Remove master relationships
	masterRelsPath := common.RelsPathFor(masterPart)
	e.parts.Delete(masterRelsPath)

	// Remove master XML
	e.parts.Delete(masterPart)

	// Remove from presentation.xml relationships
	if err := e.removeMasterFromPresentation(masterPart); err != nil {
		return err
	}

	return nil
}

// AddSlideLayout adds a new slide layout to an existing master.
// Returns the path to the newly created layout.
func (e *PresentationEditor) AddSlideLayout(masterPart, layoutName string) (string, error) {
	masterPart = common.CanonicalPartPath(masterPart)
	if !e.parts.Has(masterPart) {
		return "", fmt.Errorf("master part not found: %s", masterPart)
	}

	// Find the next available layout number
	layoutNum := editorslide.NextPartNumber(
		e.parts.KeysWithPrefix("ppt/slideLayouts/"),
		layoutNumPattern,
		nextPartPatternSubmatchSize,
	)
	layoutPart := fmt.Sprintf("ppt/slideLayouts/slideLayout%d.xml", layoutNum)

	// Determine master number from part path
	masterNum := extractMasterNumber(masterPart)

	// Create a basic layout XML
	layoutXML := editorslide.DefaultSlideLayout(layoutName, layoutNum, masterNum)
	e.parts.Set(layoutPart, []byte(layoutXML))

	// Create layout relationships
	layoutRelsPath := fmt.Sprintf("ppt/slideLayouts/_rels/slideLayout%d.xml.rels", layoutNum)
	layoutRels := editorslide.DefaultSlideLayoutRelationships(masterNum)
	e.parts.Set(layoutRelsPath, []byte(layoutRels))

	// Register the layout in Content_Types.xml
	e.addContentTypeOverride(layoutPart, "application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml")

	// Register the layout in master relationships
	if err := e.registerNewLayout(masterPart, layoutPart); err != nil {
		return "", err
	}

	return layoutPart, nil
}

// RemoveSlideLayout removes a slide layout.
func (e *PresentationEditor) RemoveSlideLayout(layoutPart string) error {
	layoutPart = common.CanonicalPartPath(layoutPart)
	if !e.parts.Has(layoutPart) {
		return fmt.Errorf("layout part not found: %s", layoutPart)
	}

	// Find the master for this layout
	masterPart, err := editorslide.ResolveLayoutMasterPart(layoutPart, e.parts.Get, parseRelationshipsXML)
	if err != nil {
		return err
	}

	// Remove layout relationships
	layoutRelsPath := common.RelsPathFor(layoutPart)
	e.parts.Delete(layoutRelsPath)

	// Remove layout XML
	e.parts.Delete(layoutPart)

	// Remove from master relationships
	if err := e.removeLayoutFromMaster(masterPart, layoutPart); err != nil {
		return err
	}

	return nil
}

func (e *PresentationEditor) registerNewMaster(masterPart string) error {
	e.recalculateNextRelIDNum()
	newMasterRelID := fmt.Sprintf("rId%d", e.nextRelIDNum)
	e.nextRelIDNum++

	e.nonSlideRels = append(e.nonSlideRels, common.EditorRelationship{
		ID:     newMasterRelID,
		Type:   common.RelTypeSlideMaster,
		Target: common.MakeRelativePath(common.PresentationXMLPath, masterPart),
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

func (e *PresentationEditor) registerNewLayout(masterPart, layoutPart string) error {
	masterRelsPath := common.RelsPathFor(masterPart)
	masterRelsData, ok := e.parts.Get(masterRelsPath)
	if !ok {
		return fmt.Errorf("master rels part not found: %s", masterRelsPath)
	}

	masterRels, err := parseRelationshipsXML(masterRelsData)
	if err != nil {
		return fmt.Errorf("parse master rels: %w", err)
	}

	// Find next available rId
	maxID := 0
	for _, rel := range masterRels {
		id, ok := common.ParseRelationshipNumber(rel.ID)
		if ok && id > maxID {
			maxID = id
		}
	}

	newRelID := fmt.Sprintf("rId%d", maxID+1)
	masterRels = append(masterRels, common.EditorRelationship{
		ID:     newRelID,
		Type:   common.RelTypeSlideLayout,
		Target: common.MakeRelativePath(masterPart, layoutPart),
	})

	rendered := renderRelationshipsXML(masterRels)
	e.parts.Set(masterRelsPath, []byte(rendered))
	return nil
}

func (e *PresentationEditor) removeMasterFromPresentation(masterPart string) error {
	masterTarget := common.MakeRelativePath(common.PresentationXMLPath, masterPart)

	newRels := make([]common.EditorRelationship, 0)
	for _, rel := range e.nonSlideRels {
		if rel.Type == common.RelTypeSlideMaster && rel.Target == masterTarget {
			continue
		}
		newRels = append(newRels, rel)
	}
	e.nonSlideRels = newRels

	// Also update presentation.xml
	presentationRelPath := common.PresentationRelPath
	presRelsData, ok := e.parts.Get(presentationRelPath)
	if !ok {
		return nil
	}

	presRels, err := parseRelationshipsXML(presRelsData)
	if err != nil {
		return fmt.Errorf("parse presentation rels: %w", err)
	}

	newPresRels := make([]common.EditorRelationship, 0)
	for _, rel := range presRels {
		if rel.Target == masterTarget {
			continue
		}
		newPresRels = append(newPresRels, rel)
	}

	rendered := renderRelationshipsXML(newPresRels)
	e.parts.Set(presentationRelPath, []byte(rendered))
	return nil
}

func (e *PresentationEditor) removeLayoutFromMaster(masterPart, layoutPart string) error {
	masterRelsPath := common.RelsPathFor(masterPart)
	masterRelsData, ok := e.parts.Get(masterRelsPath)
	if !ok {
		return fmt.Errorf("master rels not found: %s", masterRelsPath)
	}

	masterRels, err := parseRelationshipsXML(masterRelsData)
	if err != nil {
		return fmt.Errorf("parse master rels: %w", err)
	}

	layoutTarget := common.MakeRelativePath(masterPart, layoutPart)
	newMasterRels := make([]common.EditorRelationship, 0)
	for _, rel := range masterRels {
		if rel.Target == layoutTarget {
			continue
		}
		newMasterRels = append(newMasterRels, rel)
	}

	rendered := renderRelationshipsXML(newMasterRels)
	e.parts.Set(masterRelsPath, []byte(rendered))
	return nil
}

func extractMasterNumber(masterPart string) int {
	re := regexp.MustCompile(`slideMaster(\d+)\.xml`)
	matches := re.FindStringSubmatch(masterPart)
	if len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return num
		}
	}
	return 1
}
