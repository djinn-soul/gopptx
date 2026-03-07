package layoutmaster

import (
	"fmt"
	"regexp"
	"strconv"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

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
