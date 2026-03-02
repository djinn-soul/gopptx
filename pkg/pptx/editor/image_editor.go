package editor

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// AddImage adds a new image to the slide from a local file path with optional parameters.
func (e *PresentationEditor) AddImage(slideIndex int, imagePath string, x, y, w, h float64, opts *common.ShapeUpdate) (int, error) {
	return e.addImageGeneric(slideIndex, imagePath, nil, "", x, y, w, h, opts)
}

// AddImageFromBytes adds a new image to the slide from raw data.
func (e *PresentationEditor) AddImageFromBytes(slideIndex int, data []byte, format string, x, y, w, h float64, opts *common.ShapeUpdate) (int, error) {
	return e.addImageGeneric(slideIndex, "", data, format, x, y, w, h, opts)
}

func (e *PresentationEditor) addImageGeneric(slideIndex int, imagePath string, data []byte, format string, x, y, w, h float64, opts *common.ShapeUpdate) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, errors.New("slide index out of range")
	}
	if len(data) > 0 && strings.TrimSpace(format) == "" {
		return 0, errors.New("image format is required when adding image bytes")
	}

	// 1. Register image in media inventory and part store
	var relID string
	var err error
	if len(data) > 0 {
		partPath, err2 := e.RegisterImage(data, format)
		if err2 != nil {
			return 0, err2
		}
		relID, err = e.getOrCreateSlideRel(slideIndex, partPath)
	} else {
		relID, err = e.getOrCreateImageRelID(slideIndex, imagePath)
	}
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

	name := fmt.Sprintf("Picture %d", newID)

	// Build blip XML
	blipXML := fmt.Sprintf(`<a:blip r:embed="%s"/>`, relID)
	srcRectXML := ""
	if opts != nil && opts.Crop != nil {
		c := opts.Crop
		srcRectXML = fmt.Sprintf(`<a:srcRect l="%d" r="%d" t="%d" b="%d"/>`,
			int(c.Left*100000), int(c.Right*100000), int(c.Top*100000), int(c.Bottom*100000))
	}

	// Build xfrm XML with rotation/flip
	xfrmAttr := ""
	if opts != nil {
		if opts.Rotation != nil {
			xfrmAttr += fmt.Sprintf(` rot="%d"`, int(*opts.Rotation*60000))
		}
		if opts.FlipH != nil && *opts.FlipH {
			xfrmAttr += ` flipH="1"`
		}
		if opts.FlipV != nil && *opts.FlipV {
			xfrmAttr += ` flipV="1"`
		}
	}

	imageXML := fmt.Sprintf(`
<p:pic>
  <p:nvPicPr>
    <p:cNvPr id="%d" name="%s"/>
    <p:cNvPicPr><a:picLocks noChangeAspect="1"/></p:cNvPicPr>
    <p:nvPr/>
  </p:nvPicPr>
  <p:blipFill>
    %s
    %s
    <a:stretch><a:fillRect/></a:stretch>
  </p:blipFill>
  <p:spPr>
    <a:xfrm%s>
      <a:off x="%d" y="%d"/>
      <a:ext cx="%d" cy="%d"/>
    </a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
  </p:spPr>
</p:pic>`, newID, name, blipXML, srcRectXML, xfrmAttr, int64(x), int64(y), int64(w), int64(h))

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

func (e *PresentationEditor) getOrCreateSlideRel(slideIndex int, partPath string) (string, error) {
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

	nextNum := common.NextRelationshipNumber(rels)
	relID := fmt.Sprintf("rId%d", nextNum)

	rels = append(rels, common.EditorRelationship{
		ID:     relID,
		Type:   common.RelTypeImage,
		Target: target,
	})

	e.parts.Set(common.SlideRelsPartName(slideRef.Part), []byte(renderRelationshipsXML(rels)))
	return relID, nil
}

// getOrCreateImageRelID registers an image in the media inventory and creates a slide relationship.
func (e *PresentationEditor) getOrCreateImageRelID(slideIndex int, imagePath string) (string, error) {
	// 1. Register image in media inventory and part store
	partPath, err := e.registerEditorImage(imagePath, nil, "")
	if err != nil {
		return "", err
	}
	// 2. Add relationship to slide if not already present
	return e.getOrCreateSlideRel(slideIndex, partPath)
}

// GetImageMetadata returns dimensions and format for a specific image shape.
func (e *PresentationEditor) GetImageMetadata(slideIndex, shapeID int) (*common.ImageMetadata, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}

	slideRef := e.slides[slideIndex]
	content, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return nil, errors.New("read slide part: not found")
	}

	// Find the shape and its blip relID
	shapes, err := scanShapesWithOffsets(content, false)
	if err != nil {
		return nil, err
	}

	var relID string
	for _, s := range shapes {
		if s.ID == shapeID && s.Type == shapeTypePicture {
			// Extract relID from XML (simplified: scan for r:embed)
			shapeXML := content[s.Start:s.End]
			re := regexp.MustCompile(`r:embed="([^"]+)"`)
			match := re.FindSubmatch(shapeXML)
			if len(match) > 1 {
				relID = string(match[1])
			}
			break
		}
	}

	if relID == "" {
		return nil, fmt.Errorf("image shape %d not found or has no embed rel", shapeID)
	}

	// Resolve relID to part path
	rels, err := e.slideRelationships(slideRef.Part)
	if err != nil {
		return nil, err
	}

	var partPath string
	for _, r := range rels {
		if r.ID == relID {
			partPath = common.CanonicalPartPath(path.Join("ppt/slides", r.Target))
			break
		}
	}

	if partPath == "" {
		return nil, fmt.Errorf("could not resolve relationship %s", relID)
	}

	data, ok := e.parts.Get(partPath)
	if !ok {
		return nil, fmt.Errorf("media part %s not found", partPath)
	}

	// Decode image config
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image %s: %w", partPath, err)
	}

	return &common.ImageMetadata{
		Width:  config.Width,
		Height: config.Height,
		Format: format,
	}, nil
}
