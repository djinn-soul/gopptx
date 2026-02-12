package editor

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func renderEditorSlideParts(e *PresentationEditor, slide elements.SlideContent, slideNumber int, notesTarget string, width, height int64) (string, string, error) {
	tableSpec, err := renderEditorTableSpec(slide, slideNumber)
	if err != nil {
		return "", "", err
	}

	imageRefs := make([]pptxxml.ImageRef, 0, len(slide.Images))
	imageTargets := make([]string, 0, len(slide.Images))

	for i, img := range slide.Images {
		data := img.Data
		format := img.Format
		if len(data) == 0 && img.Path != "" {
			d, err := os.ReadFile(img.Path)
			if err != nil {
				return "", "", fmt.Errorf("read image %d: %w", i+1, err)
			}
			data = d
			if format == "" {
				if idx := strings.LastIndex(img.Path, "."); idx >= 0 {
					format = img.Path[idx+1:]
				}
			}
		}

		if len(data) == 0 {
			return "", "", fmt.Errorf("slide %d image %d has no data or path", slideNumber, i+1)
		}

		partPath, err := e.RegisterImage(data, format)
		if err != nil {
			return "", "", err
		}

		relID := fmt.Sprintf("rId%d", i+2)
		imageTargets = append(imageTargets, "../media/"+path.Base(partPath))

		ref := pptxxml.ImageRef{
			RelID:        relID,
			Name:         fmt.Sprintf("Picture %d", i+1),
			X:            img.X.Emu(),
			Y:            img.Y.Emu(),
			CX:           img.CX.Emu(),
			CY:           img.CY.Emu(),
			Rotation:     int64(img.Rotation * 60000),
			FlipH:        img.FlipH,
			FlipV:        img.FlipV,
			Shadow:       img.Shadow,
			Reflection:   img.Reflection,
			AltText:      img.AltText,
			IsDecorative: img.IsDecorative,
		}
		if img.Crop != (shapes.ImageCrop{}) {
			ref.Crop = &pptxxml.ImageCropRef{
				Left:   int64(img.Crop.Left * 100000),
				Right:  int64(img.Crop.Right * 100000),
				Top:    int64(img.Crop.Top * 100000),
				Bottom: int64(img.Crop.Bottom * 100000),
			}
		}
		imageRefs = append(imageRefs, ref)
	}

	backgroundRID := ""
	if slide.Background != nil && slide.Background.Type == elements.SlideBackgroundPicture && slide.Background.PictureFill != nil {
		img := *slide.Background.PictureFill
		data := img.Data
		format := img.Format
		if len(data) == 0 && img.Path != "" {
			d, _ := os.ReadFile(img.Path)
			data = d
		}
		if len(data) > 0 {
			partPath, _ := e.RegisterImage(data, format)
			backgroundRID = fmt.Sprintf("rId%d", len(imageTargets)+2)
			imageTargets = append(imageTargets, "../media/"+path.Base(partPath))
		}
	}

	layoutMode := elements.SlideLayoutXMLMode(slide.Layout)
	hyperlinkRIDs, hyperlinks, _ := elements.BuildSlideHyperlinkRels(slide, len(imageTargets)+2)
	shapeIDs := elements.CalculateShapeIDs(slide)

	titleSpec := pptxxml.TitleSpec{
		Text:      slide.Title,
		SizePt:    slide.TitleSize,
		Color:     slide.TitleColor,
		Bold:      slide.TitleBold,
		Italic:    slide.TitleItalic,
		Underline: slide.TitleUnderline,
		Align:     slide.TitleAlign,
		Font:      slide.TitleFont,
	}
	contentStyle := pptxxml.ContentStyleSpec{
		SizePt:    slide.ContentSize,
		Color:     slide.ContentColor,
		Bold:      slide.ContentBold,
		Italic:    slide.ContentItalic,
		Underline: slide.ContentUnderline,
		VAlign:    slide.ContentVAlign,
	}

	slideXML := pptxxml.SlideWithLayout(
		layoutMode,
		titleSpec,
		slide.Bullets,
		elements.ToXMLBulletParagraphStyles(slide.BulletStyles),
		elements.ToXMLTextRunRows(slide.BulletRuns, hyperlinkRIDs),
		contentStyle,
		tableSpec,
		nil,
		imageRefs,
		shapes.ToXMLShapeSpecs(slide.Shapes, hyperlinkRIDs),
		shapes.ToXMLConnectorSpecs(slide.Connectors, slide.Shapes),
		nil,
		elements.ToXMLBackgroundSpec(slide.Background, backgroundRID),
		elements.SlideTransitionXML(slide),
		elements.SlideAnimationsXML(slide, shapeIDs),
		slide.ShowSlideNumber,
		"",
		false,
		width,
		height,
	)
	relsXML := pptxxml.SlideRelationshipsWithHyperlinks(
		elements.SlideLayoutTarget(slide.Layout),
		imageTargets,
		nil,
		notesTarget,
		hyperlinks,
	)
	return slideXML, relsXML, nil
}

func renderEditorTableSpec(slide elements.SlideContent, slideNumber int) (*pptxxml.TableSpec, error) {
	if slide.Table == nil {
		return nil, nil
	}
	return slide.Table.ToTableSpec(slideNumber)
}

func editorEnsureSlideRelsExist(parts map[string][]byte, slidePart string) error {
	relsPath := common.SlideRelsPartName(slidePart)
	if _, ok := parts[relsPath]; ok {
		return nil
	}
	return fmt.Errorf("missing slide relationships part %q", relsPath)
}
