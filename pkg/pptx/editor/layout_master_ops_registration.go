package editor

import (
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	layoutmaster "github.com/djinn-soul/gopptx/pkg/pptx/editor/layoutmaster"
)

func (e *PresentationEditor) registerNewMaster(masterPart string) error {
	e.recalculateNextRelIDNum()
	updatedRels, updatedPresentationXML, nextRelIDNum, err := layoutmaster.AddMasterRelationship(
		e.nonSlideRels,
		e.presentationXML,
		e.nextRelIDNum,
		masterPart,
	)
	if err != nil {
		return err
	}
	e.nonSlideRels = updatedRels
	e.presentationXML = updatedPresentationXML
	e.nextRelIDNum = nextRelIDNum
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

	masterRels = layoutmaster.AppendLayoutRelationship(masterRels, masterPart, layoutPart)

	rendered := renderRelationshipsXML(masterRels)
	e.parts.Set(masterRelsPath, []byte(rendered))
	return nil
}

func (e *PresentationEditor) removeMasterFromPresentation(masterPart string) error {
	masterTarget := common.MakeRelativePath(common.PresentationXMLPath, masterPart)

	newRels := make([]common.EditorRelationship, 0, len(e.nonSlideRels))
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

	newPresRels := layoutmaster.FilterOutRelationshipTarget(presRels, masterTarget)
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
	newMasterRels := layoutmaster.FilterOutRelationshipTarget(masterRels, layoutTarget)

	rendered := renderRelationshipsXML(newMasterRels)
	e.parts.Set(masterRelsPath, []byte(rendered))
	return nil
}
