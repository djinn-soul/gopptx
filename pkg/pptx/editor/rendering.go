package editor

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func renderEditorSlideParts(slide elements.SlideContent, slideNumber int, notesTarget string, width, height int64) (string, string, error) {
	tableSpec, err := renderEditorTableSpec(slide, slideNumber)
	if err != nil {
		return "", "", err
	}

	layoutMode := elements.SlideLayoutXMLMode(slide.Layout)
	hyperlinkRIDs, hyperlinks, _ := elements.BuildSlideHyperlinkRels(slide, 2)
	shapeIDs := elements.CalculateShapeIDs(slide)

	titleSpec := pptxxml.TitleSpec{
		Text:      slide.Title,
		SizePt:    slide.TitleSize,
		Color:     slide.TitleColor,
		Bold:      slide.TitleBold,
		Italic:    slide.TitleItalic,
		Underline: slide.TitleUnderline,
	}
	contentStyle := pptxxml.ContentStyleSpec{
		SizePt:    slide.ContentSize,
		Color:     slide.ContentColor,
		Bold:      slide.ContentBold,
		Italic:    slide.ContentItalic,
		Underline: slide.ContentUnderline,
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
		nil,
		shapes.ToXMLShapeSpecs(slide.Shapes, hyperlinkRIDs),
		shapes.ToXMLConnectorSpecs(slide.Connectors, slide.Shapes),
		nil,
		elements.SlideTransitionXML(slide),
		elements.SlideAnimationsXML(slide, shapeIDs),
		width,
		height,
	)
	relsXML := pptxxml.SlideRelationshipsWithHyperlinks(
		elements.SlideLayoutTarget(slide.Layout),
		nil,
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
