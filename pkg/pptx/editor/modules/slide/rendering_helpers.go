package slide

import (
	"errors"
	"fmt"
	"path"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	RotationDegreeToOOXML = 60000
	cropFractionToOOXML   = 100000
)

var ErrImagePayloadEmpty = errors.New("image has no data or path")

type RegisterImageFromPathFn func(pathValue string, format string) (string, error)
type RegisterImageBytesFn func(data []byte, format string) (string, error)

func RegisterEditorImage(
	pathValue string,
	data []byte,
	format string,
	registerFromPath RegisterImageFromPathFn,
	registerFromBytes RegisterImageBytesFn,
) (string, error) {
	switch {
	case pathValue != "" && len(data) == 0:
		return registerFromPath(pathValue, format)
	case len(data) > 0:
		return registerFromBytes(data, format)
	default:
		return "", ErrImagePayloadEmpty
	}
}

func RenderSlideImageRef(
	img shapes.Image,
	index int,
	slideNumber int,
	firstImageRelationshipID int,
	registerImage func(pathValue string, data []byte, format string) (string, error),
) (pptxxml.ImageRef, string, error) {
	partPath, imgErr := registerImage(img.Path, img.Data, img.Format)
	if imgErr != nil {
		if errors.Is(imgErr, ErrImagePayloadEmpty) {
			return pptxxml.ImageRef{}, "", fmt.Errorf("slide %d image %d has no data or path", slideNumber, index+1)
		}
		return pptxxml.ImageRef{}, "", fmt.Errorf("read image %d: %w", index+1, imgErr)
	}

	relID := fmt.Sprintf("rId%d", index+firstImageRelationshipID)
	ref := pptxxml.ImageRef{
		RelID:        relID,
		Name:         fmt.Sprintf("Picture %d", index+1),
		X:            img.X.Emu(),
		Y:            img.Y.Emu(),
		CX:           img.CX.Emu(),
		CY:           img.CY.Emu(),
		Rotation:     int64(img.Rotation * RotationDegreeToOOXML),
		FlipH:        img.FlipH,
		FlipV:        img.FlipV,
		Shadow:       img.Shadow,
		Reflection:   img.Reflection,
		AltText:      img.AltText,
		IsDecorative: img.IsDecorative,
	}
	if img.Crop != (shapes.ImageCrop{}) {
		ref.Crop = &pptxxml.ImageCropRef{
			Left:   int64(img.Crop.Left * cropFractionToOOXML),
			Right:  int64(img.Crop.Right * cropFractionToOOXML),
			Top:    int64(img.Crop.Top * cropFractionToOOXML),
			Bottom: int64(img.Crop.Bottom * cropFractionToOOXML),
		}
	}
	return ref, "../media/" + path.Base(partPath), nil
}

func RenderBackgroundImageTarget(
	background *elements.SlideBackground,
	currentImageCount int,
	firstImageRelationshipID int,
	registerImage func(pathValue string, data []byte, format string) (string, error),
) (string, string, error) {
	if background == nil || background.Type != elements.SlideBackgroundPicture || background.PictureFill == nil {
		return "", "", nil
	}

	partPath, err := registerImage(
		background.PictureFill.Path,
		background.PictureFill.Data,
		background.PictureFill.Format,
	)
	if err != nil {
		if errors.Is(err, ErrImagePayloadEmpty) {
			return "", "", nil
		}
		return "", "", fmt.Errorf("read background image: %w", err)
	}

	backgroundRID := fmt.Sprintf("rId%d", currentImageCount+firstImageRelationshipID)
	return backgroundRID, "../media/" + path.Base(partPath), nil
}

func RenderPlaceholderImageRef(
	override shapes.PlaceholderContent,
	ridIndex int,
	registerImage func(pathValue string, data []byte, format string) (string, error),
) (*pptxxml.ImageRef, string, error) {
	if override.Image == nil {
		return nil, "", nil
	}

	partPath, err := registerImage(override.Image.Path, override.Image.Data, override.Image.Format)
	if err != nil {
		if errors.Is(err, ErrImagePayloadEmpty) {
			return nil, "", nil
		}
		return nil, "", fmt.Errorf("placeholder image %d: %w", override.Index, err)
	}

	imageRef := &pptxxml.ImageRef{
		RelID: fmt.Sprintf("rId%d", ridIndex),
		Name:  "Placeholder Picture",
		X:     override.Image.X.Emu(),
		Y:     override.Image.Y.Emu(),
		CX:    override.Image.CX.Emu(),
		CY:    override.Image.CY.Emu(),
	}
	return imageRef, "../media/" + path.Base(partPath), nil
}

func EditorNotesBody(slide elements.SlideContent) []elements.Paragraph {
	if len(slide.NotesBody) > 0 {
		return slide.NotesBody
	}

	p := elements.NewParagraph()
	p.Runs = append(p.Runs, elements.NewRun(slide.Notes))
	return []elements.Paragraph{p}
}

func MapOptionalLength(l *styling.Length) *int64 {
	if l == nil {
		return nil
	}
	val := l.Emu()
	return &val
}

func MapPlaceholderTextStyle(ts *shapes.PlaceholderTextStyle) *pptxxml.PlaceholderTextStyleSpec {
	if ts == nil {
		return nil
	}
	return &pptxxml.PlaceholderTextStyleSpec{
		SizePt:    ts.SizePt,
		Color:     ts.Color,
		Bold:      ts.Bold,
		Italic:    ts.Italic,
		Underline: ts.Underline,
		Align:     ts.Align,
		Font:      ts.Font,
	}
}
