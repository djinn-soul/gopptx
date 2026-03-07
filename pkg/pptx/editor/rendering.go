package editor

import (
	"errors"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

const (
	firstImageRelationshipID = 2
	rotationDegreeToOOXML    = 60000
	cropFractionToOOXML      = 100000
)

var errNoSlideTable = errors.New("slide has no table")

func renderEditorSlideParts(
	e *PresentationEditor,
	slide elements.SlideContent,
	slidePart string,
	slideNumber int,
	existingNotesTarget string,
	width, height int64,
) (string, string, error) {
	tableSpec, err := renderEditorTableSpec(slide, slideNumber)
	if err != nil && !errors.Is(err, errNoSlideTable) {
		return "", "", err
	}
	if errors.Is(err, errNoSlideTable) {
		tableSpec = nil
	}

	imageRefs, imageTargets, err := e.renderSlideImages(slide.Images, slideNumber)
	if err != nil {
		return "", "", err
	}

	backgroundRID, backgroundTarget, err := e.renderBackgroundImageTarget(slide.Background, len(imageTargets))
	if err != nil {
		return "", "", err
	}
	if backgroundTarget != "" {
		imageTargets = append(imageTargets, backgroundTarget)
	}

	layoutMode := elements.SlideLayoutXMLMode(slide.Layout)
	hyperlinkRIDs, hyperlinks, nextRID := elements.BuildSlideHyperlinkRels(
		slide,
		len(imageTargets)+firstImageRelationshipID,
	)

	// Process placeholder overrides
	placeholderSpecs, phImageTargets, phChartRels, err := renderEditorPlaceholderSpecs(
		e,
		slide,
		slidePart,
		slideNumber,
		nextRID,
	)
	if err != nil {
		return "", "", err
	}
	imageTargets = append(imageTargets, phImageTargets...)

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
		nil, // chartFrame
		imageRefs,
		shapes.ToXMLShapeSpecs(slide.Shapes, hyperlinkRIDs),
		shapes.ToXMLConnectorSpecs(slide.Connectors, slide.Shapes),
		placeholderSpecs,
		nil, // smartArtFrames
		elements.ToXMLBackgroundSpec(slide.Background, backgroundRID),
		elements.SlideTransitionXML(slide),
		elements.SlideAnimationsXML(slide, elements.CalculateShapeIDs(slide)),
		slide.ShowSlideNumber,
		"",    // footerText
		false, // showDateTime
		width,
		height,
	)

	notesTarget := e.renderSlideNotesTarget(slide, slideNumber, existingNotesTarget)

	relsXML := pptxxml.SlideRelationshipsWithMultiCharts(
		elements.SlideLayoutTarget(slide.Layout),
		imageTargets,
		nil,
		phChartRels,
		nil, // smartArtRels
		notesTarget,
		hyperlinks,
		"", // commentsTarget
	)
	return slideXML, relsXML, nil
}

func (e *PresentationEditor) renderSlideImages(
	images []shapes.Image,
	slideNumber int,
) ([]pptxxml.ImageRef, []string, error) {
	imageRefs := make([]pptxxml.ImageRef, 0, len(images))
	imageTargets := make([]string, 0, len(images))
	for i, img := range images {
		ref, target, err := e.renderSlideImageRef(img, i, slideNumber)
		if err != nil {
			return nil, nil, err
		}
		imageRefs = append(imageRefs, ref)
		imageTargets = append(imageTargets, target)
	}
	return imageRefs, imageTargets, nil
}

func (e *PresentationEditor) renderSlideImageRef(
	img shapes.Image,
	index int,
	slideNumber int,
) (pptxxml.ImageRef, string, error) {
	return editorslide.RenderSlideImageRef(
		img,
		index,
		slideNumber,
		firstImageRelationshipID,
		e.registerEditorImage,
	)
}

func renderEditorTableSpec(slide elements.SlideContent, slideNumber int) (*pptxxml.TableSpec, error) {
	if slide.Table == nil {
		return nil, errNoSlideTable
	}
	return slide.Table.ToTableSpec(slideNumber)
}

func (e *PresentationEditor) registerEditorImage(pathValue string, data []byte, format string) (string, error) {
	return editorslide.RegisterEditorImage(pathValue, data, format, e.registerImageFromPath, e.RegisterImage)
}

func (e *PresentationEditor) renderBackgroundImageTarget(
	background *elements.SlideBackground,
	currentImageCount int,
) (string, string, error) {
	return editorslide.RenderBackgroundImageTarget(
		background,
		currentImageCount,
		firstImageRelationshipID,
		e.registerEditorImage,
	)
}
