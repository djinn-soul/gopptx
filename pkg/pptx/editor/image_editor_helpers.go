package editor

import (
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	imageCropScale     = 100000
	imageRotationScale = 60000
)

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
}
