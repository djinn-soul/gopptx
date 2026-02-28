package editor

import (
	"bytes"
	"errors"
	"fmt"
	"path"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// AddImage adds a new image to the slide from a local file path.
func (e *PresentationEditor) AddImage(slideIndex int, imagePath string, x, y, w, h float64) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, errors.New("slide index out of range")
	}

	// 1. Register image in media inventory and part store
	relID, err := e.getOrCreateImageRelID(slideIndex, imagePath)
	if err != nil {
		return 0, fmt.Errorf("register image: %w", err)
	}

	// 2. Generate image XML
	slideRef := e.slides[slideIndex]
	content, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return 0, errors.New("read slide part: not found")
	}

	maxID := maxObjectID(content)
	newID := maxID + 1

	imageRef := pptxxml.ImageRef{Name: fmt.Sprintf("Picture %d", newID)}

	// Internal helper from slide_image_xml.go (we need to make it accessible or replicate)
	// For now, let's use a simplified version here
	imageXML := fmt.Sprintf(`
<p:pic>
  <p:nvPicPr>
    <p:cNvPr id="%d" name="%s"/>
    <p:cNvPicPr><a:picLocks noChangeAspect="1"/></p:cNvPicPr>
    <p:nvPr/>
  </p:nvPicPr>
  <p:blipFill>
    <a:blip r:embed="%s"/>
    <a:stretch><a:fillRect/></a:stretch>
  </p:blipFill>
  <p:spPr>
    <a:xfrm>
      <a:off x="%d" y="%d"/>
      <a:ext cx="%d" cy="%d"/>
    </a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
  </p:spPr>
</p:pic>`, newID, imageRef.Name, relID, int64(x), int64(y), int64(w), int64(h))

	// 3. Update slide XML
	endTree := []byte("</p:spTree>")
	idx := bytes.LastIndex(content, endTree)
	if idx == -1 {
		return 0, errors.New("invalid slide xml: missing spTree end")
	}

	var buf bytes.Buffer
	buf.Write(content[:idx])
	buf.WriteString(imageXML)
	buf.Write(content[idx:])

	e.parts.Set(slideRef.Part, buf.Bytes())

	return newID, nil
}

// getOrCreateImageRelID registers an image in the media inventory and creates a slide relationship.
func (e *PresentationEditor) getOrCreateImageRelID(slideIndex int, imagePath string) (string, error) {
	// 1. Register image in media inventory and part store
	partPath, err := e.registerEditorImage(imagePath, nil, "")
	if err != nil {
		return "", err
	}

	// 2. Add relationship to slide if not already present
	slideRef := e.slides[slideIndex]
	rels, err := e.slideRelationships(slideRef.Part)
	if err != nil {
		return "", err
	}

	target := "../media/" + path.Base(partPath)
	for _, r := range rels {
		if r.Type == common.RelTypeImage && r.Target == target {
			return r.ID, nil
		}
	}

	nextNum := 1
	for _, r := range rels {
		if n, ok := common.ParseRelationshipNumber(r.ID); ok && n >= nextNum {
			nextNum = n + 1
		}
	}
	relID := fmt.Sprintf("rId%d", nextNum)

	rels = append(rels, common.EditorRelationship{
		ID:     relID,
		Type:   common.RelTypeImage,
		Target: target,
	})

	e.parts.Set(common.SlideRelsPartName(slideRef.Part), []byte(renderRelationshipsXML(rels)))
	return relID, nil
}
