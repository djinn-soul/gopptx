package editor

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
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

func (e *PresentationEditor) renderSlideNotesTarget(
	slide elements.SlideContent,
	slideNumber int,
	existingNotesTarget string,
) string {
	notesTarget := strings.TrimSpace(existingNotesTarget)
	if strings.TrimSpace(slide.Notes) == "" {
		return notesTarget
	}

	e.ensureNotesInfrastructure()
	slidePath := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
	notesPath := e.ensureSlideNotesPart(slidePath)
	e.parts.Set(notesPath, []byte(pptxxml.NotesSlide(editorslide.EditorNotesBody(slide))))

	notesRelsPath := common.SlideRelsPartName(notesPath)
	e.parts.Set(notesRelsPath, []byte(pptxxml.NotesSlideRelationships(slideNumber)))
	return "../notesSlides/" + path.Base(notesPath)
}

func (e *PresentationEditor) ensureSlideNotesPart(slidePath string) string {
	notesPath, ok := e.notesInventory[slidePath]
	if ok {
		return notesPath
	}
	notesPath = fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", e.nextNotesNum)
	e.nextNotesNum++
	e.notesInventory[slidePath] = notesPath
	return notesPath
}

func renderEditorTableSpec(slide elements.SlideContent, slideNumber int) (*pptxxml.TableSpec, error) {
	if slide.Table == nil {
		return nil, errNoSlideTable
	}
	return slide.Table.ToTableSpec(slideNumber)
}

// renderEditorPlaceholderSpecs converts SlideContent.PlaceholderOverrides into
// XML specs for the editor rendering path. It returns the specs, any additional
// image relationship targets, chart rels, and an error.
//
//nolint:gocognit // Placeholder rendering must keep explicit branching for supported override kinds.
func renderEditorPlaceholderSpecs(
	e *PresentationEditor,
	slide elements.SlideContent,
	slidePart string,
	slideNumber int,
	startRID int,
) ([]pptxxml.PlaceholderOverrideSpec, []string, []pptxxml.ChartRel, error) {
	if len(slide.PlaceholderOverrides) == 0 {
		return nil, nil, nil, nil
	}

	specs := make([]pptxxml.PlaceholderOverrideSpec, 0, len(slide.PlaceholderOverrides))
	var imageTargets []string
	var chartRels []pptxxml.ChartRel
	currentRID := startRID
	placeholders, err := editorslide.LookupSlidePlaceholders(
		slidePart,
		e.parts.Get,
		func(content []byte) []editorslide.PlaceholderMeta {
			parsed := parsePlaceholdersFromSlideXML(content)
			out := make([]editorslide.PlaceholderMeta, 0, len(parsed))
			for _, ph := range parsed {
				out = append(out, editorslide.PlaceholderMeta{
					Name:  ph.Name,
					Type:  ph.Type,
					Index: ph.Index,
				})
			}
			return out
		},
	)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, override := range slide.PlaceholderOverrides {
		targetType, targetIndex, err := editorslide.ResolvePlaceholderTarget(override, placeholders)
		if err != nil {
			return nil, nil, nil, err
		}
		spec := pptxxml.PlaceholderOverrideSpec{
			Index: targetIndex,
			Type:  targetType,
			Text:  override.Text,
		}

		if override.Override != nil {
			spec.X = editorslide.MapOptionalLength(override.Override.X)
			spec.Y = editorslide.MapOptionalLength(override.Override.Y)
			spec.CX = editorslide.MapOptionalLength(override.Override.CX)
			spec.CY = editorslide.MapOptionalLength(override.Override.CY)
			spec.TextStyle = editorslide.MapPlaceholderTextStyle(override.Override.TextStyle)
		}

		// Handle image placeholder
		if override.Image != nil {
			imageRef, imageTarget, imageErr := e.renderPlaceholderImageRef(override, currentRID)
			if imageErr != nil {
				return nil, nil, nil, imageErr
			}
			if imageRef != nil {
				spec.Image = imageRef
				imageTargets = append(imageTargets, imageTarget)
				currentRID++
			}
		}

		// Handle table placeholder
		if override.Table != nil {
			tableSpec, err := override.Table.ToTableSpec(slideNumber)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("placeholder table %d: %w", targetIndex, err)
			}
			spec.Table = tableSpec
		}

		// Handle chart placeholder
		if override.Chart != nil {
			chartSpec := override.Chart.ToChartSpec()
			chartPath := fmt.Sprintf("ppt/charts/chart_ph_%d_%d.xml", slideNumber, targetIndex)
			e.parts.Set(chartPath, []byte(pptxxml.ChartPartXML(chartSpec)))

			rid := fmt.Sprintf("rId%d", currentRID)
			currentRID++
			spec.Chart = &pptxxml.ChartFrame{
				RelID: rid,
				X:     chartSpec.X,
				Y:     chartSpec.Y,
				CX:    chartSpec.CX,
				CY:    chartSpec.CY,
			}
			chartRels = append(chartRels, pptxxml.ChartRel{
				RID:    rid,
				Target: "../charts/" + path.Base(chartPath),
			})
		}

		specs = append(specs, spec)
	}

	return specs, imageTargets, chartRels, nil
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

func (e *PresentationEditor) renderPlaceholderImageRef(
	override shapes.PlaceholderContent,
	ridIndex int,
) (*pptxxml.ImageRef, string, error) {
	return editorslide.RenderPlaceholderImageRef(override, ridIndex, e.registerEditorImage)
}
