package editor

import (
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	layoutmaster "github.com/djinn-soul/gopptx/pkg/pptx/editor/layoutmaster"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

// AddSlideMaster adds a new slide master to the presentation.
// Returns the path to the newly created master.
func (e *PresentationEditor) AddSlideMaster() (string, error) {
	masterNum := editorslide.NextPartNumber(
		e.parts.KeysWithPrefix("ppt/slideMasters/"),
		masterNumPattern,
		nextPartPatternSubmatchSize,
	)
	masterPart := fmt.Sprintf("ppt/slideMasters/slideMaster%d.xml", masterNum)

	masterXML := []byte(editorslide.DefaultSlideMaster())
	e.parts.Set(masterPart, masterXML)

	masterRelsPath := fmt.Sprintf("ppt/slideMasters/_rels/slideMaster%d.xml.rels", masterNum)
	masterRels := editorslide.DefaultSlideMasterRelationships()
	e.parts.Set(masterRelsPath, []byte(masterRels))

	e.addContentTypeOverride(masterPart, "application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml")
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

	masters, err := e.ListSlideMasters()
	if err != nil {
		return err
	}
	if len(masters) <= 1 {
		return errors.New("cannot remove the last slide master")
	}

	layouts, err := e.layoutsForMaster(masterPart)
	if err != nil {
		return err
	}
	for _, layoutPart := range layouts {
		if err := e.RemoveSlideLayout(layoutPart); err != nil {
			return err
		}
	}

	e.parts.Delete(common.RelsPathFor(masterPart))
	e.parts.Delete(masterPart)
	return e.removeMasterFromPresentation(masterPart)
}

// AddSlideLayout adds a new slide layout to an existing master.
// Returns the path to the newly created layout.
func (e *PresentationEditor) AddSlideLayout(masterPart, layoutName string) (string, error) {
	masterPart = common.CanonicalPartPath(masterPart)
	if !e.parts.Has(masterPart) {
		return "", fmt.Errorf("master part not found: %s", masterPart)
	}

	layoutNum := editorslide.NextPartNumber(
		e.parts.KeysWithPrefix("ppt/slideLayouts/"),
		layoutNumPattern,
		nextPartPatternSubmatchSize,
	)
	layoutPart := fmt.Sprintf("ppt/slideLayouts/slideLayout%d.xml", layoutNum)
	masterNum := layoutmaster.ExtractMasterNumber(masterPart)

	layoutXML := editorslide.DefaultSlideLayout(layoutName)
	e.parts.Set(layoutPart, []byte(layoutXML))

	layoutRelsPath := fmt.Sprintf("ppt/slideLayouts/_rels/slideLayout%d.xml.rels", layoutNum)
	layoutRels := editorslide.DefaultSlideLayoutRelationships(masterNum)
	e.parts.Set(layoutRelsPath, []byte(layoutRels))

	e.addContentTypeOverride(layoutPart, "application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml")
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

	masterPart, err := editorslide.ResolveLayoutMasterPart(layoutPart, e.parts.Get, parseRelationshipsXML)
	if err != nil {
		return err
	}

	e.parts.Delete(common.RelsPathFor(layoutPart))
	e.parts.Delete(layoutPart)
	return e.removeLayoutFromMaster(masterPart, layoutPart)
}
