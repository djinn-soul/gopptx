package editor

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"  // register GIF decoder for DecodeConfig metadata reads
	_ "image/jpeg" // register JPEG decoder for DecodeConfig metadata reads
	_ "image/png"  // register PNG decoder for DecodeConfig metadata reads
	"path"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodmedia "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/media"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

const (
	imageCropScale     = 100000
	imageRotationScale = 60000
)

// AddImage adds a new image to the slide from a local file path with optional parameters.
func (e *PresentationEditor) AddImage(
	slideIndex int,
	imagePath string,
	x, y, w, h float64,
	opts *common.ShapeUpdate,
) (int, error) {
	return e.addImageGeneric(slideIndex, imagePath, nil, "", x, y, w, h, opts)
}

// AddImageFromBytes adds a new image to the slide from raw data.
func (e *PresentationEditor) AddImageFromBytes(
	slideIndex int,
	data []byte,
	format string,
	x, y, w, h float64,
	opts *common.ShapeUpdate,
) (int, error) {
	return e.addImageGeneric(slideIndex, "", data, format, x, y, w, h, opts)
}

func (e *PresentationEditor) addImageGeneric(
	slideIndex int,
	imagePath string,
	data []byte,
	format string,
	x, y, w, h float64,
	opts *common.ShapeUpdate,
) (int, error) {
	if err := editormodmedia.ValidateMediaSlideIndex(slideIndex, len(e.slides)); err != nil {
		return 0, err
	}
	if len(data) > 0 && strings.TrimSpace(format) == "" {
		return 0, errors.New("image format is required when adding image bytes")
	}

	relID, err := resolveAddImageRelID(e, slideIndex, imagePath, data, format)
	if err != nil {
		return 0, fmt.Errorf("register image: %w", err)
	}

	slideRef := e.slides[slideIndex]
	content, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return 0, errors.New("read slide part: not found")
	}

	maxID := editorshape.MaxObjectID(content, cNvPrIDPattern, cNvPrSubmatchSize)
	newID := maxID + 1

	imageXML := buildImageShapeXML(newID, relID, x, y, w, h, opts)
	updatedContent, err := editormodmedia.AppendShapeXMLToSlide(content, imageXML)
	if err != nil {
		return 0, err
	}
	e.parts.Set(slideRef.Part, updatedContent)

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

func resolveAddImageRelID(
	e *PresentationEditor,
	slideIndex int,
	imagePath string,
	data []byte,
	format string,
) (string, error) {
	if len(data) > 0 {
		partPath, err := e.RegisterImage(data, format)
		if err != nil {
			return "", err
		}
		return e.getOrCreateSlideRel(slideIndex, partPath)
	}
	return e.getOrCreateImageRelID(slideIndex, imagePath)
}

func buildImageCropXML(opts *common.ShapeUpdate) string {
	if opts == nil || opts.Crop == nil {
		return ""
	}
	c := opts.Crop
	return fmt.Sprintf(
		`<a:srcRect l="%d" r="%d" t="%d" b="%d"/>`,
		int(c.Left*imageCropScale),
		int(c.Right*imageCropScale),
		int(c.Top*imageCropScale),
		int(c.Bottom*imageCropScale),
	)
}

func buildImageTransformAttrs(opts *common.ShapeUpdate) string {
	if opts == nil {
		return ""
	}
	var attrs strings.Builder
	if opts.Rotation != nil {
		attrs.WriteString(fmt.Sprintf(` rot="%d"`, int(*opts.Rotation*imageRotationScale)))
	}
	if opts.FlipH != nil && *opts.FlipH {
		attrs.WriteString(` flipH="1"`)
	}
	if opts.FlipV != nil && *opts.FlipV {
		attrs.WriteString(` flipV="1"`)
	}
	return attrs.String()
}

func buildImageShapeXML(
	newID int,
	relID string,
	x, y, w, h float64,
	opts *common.ShapeUpdate,
) string {
	name := fmt.Sprintf("Picture %d", newID)
	blipXML := fmt.Sprintf(`<a:blip r:embed="%s"/>`, relID)
	srcRectXML := buildImageCropXML(opts)
	xfrmAttr := buildImageTransformAttrs(opts)

	return fmt.Sprintf(`
<p:pic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
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
}
