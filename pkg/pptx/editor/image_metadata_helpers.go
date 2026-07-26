package editor

import (
	"bytes"
	"fmt"
	"image"
	"path"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func imageRelIDForShape(content []byte, shapeID int) (string, error) {
	shapes, err := scanShapesWithOffsets(content, false)
	if err != nil {
		return "", err
	}

	for _, shape := range shapes {
		if shape.ID != shapeID || shape.Type != shapeTypePicture {
			continue
		}
		match := embeddedImageRelPattern.FindSubmatch(content[shape.Start:shape.End])
		if len(match) > 1 {
			return string(match[1]), nil
		}
		break
	}

	return "", fmt.Errorf("image shape %d not found or has no embed rel", shapeID)
}

func (e *PresentationEditor) imagePartData(
	slidePart string,
	relID string,
) ([]byte, string, error) {
	rels, err := e.slideRelationships(slidePart)
	if err != nil {
		return nil, "", err
	}

	for _, rel := range rels {
		if rel.ID != relID {
			continue
		}
		partPath := common.CanonicalPartPath(path.Join("ppt/slides", rel.Target))
		data, ok := e.parts.Get(partPath)
		if !ok {
			return nil, "", fmt.Errorf("media part %s not found", partPath)
		}
		return data, partPath, nil
	}

	return nil, "", fmt.Errorf("could not resolve relationship %s", relID)
}

func decodeImageConfig(data []byte, partPath string) (image.Config, string) {
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err == nil {
		return config, format
	}

	format = strings.TrimPrefix(strings.ToLower(path.Ext(partPath)), ".")
	if format == "" {
		format = formatPNG
	}
	return image.Config{}, format
}

func imageBlipXML(relID string, opts *common.ShapeUpdate) string {
	if opts == nil || (!strings.EqualFold(opts.Format, "svg") && !opts.IsSVG) {
		return fmt.Sprintf(`<a:blip r:embed="%s"/>`, relID)
	}

	return fmt.Sprintf(
		`<a:blip r:embed="%s"><a:extLst><a:ext uri="{96DAC542-722E-43E9-B0A7-37CDD9819D23}">`+
			`<asvg:svgBlip xmlns:asvg="http://schemas.microsoft.com/office/drawing/2016/SVG/main" `+
			`r:embed="%s"/></a:ext></a:extLst></a:blip>`,
		relID,
		relID,
	)
}
