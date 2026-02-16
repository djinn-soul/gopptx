package editor

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
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
	placeholderSpecs, phImageTargets, phChartRels, err := renderEditorPlaceholderSpecs(e, slide, slideNumber, nextRID)
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
		nil,
		imageRefs,
		shapes.ToXMLShapeSpecs(slide.Shapes, hyperlinkRIDs),
		shapes.ToXMLConnectorSpecs(slide.Connectors, slide.Shapes),
		placeholderSpecs,
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
		notesTarget,
		hyperlinks,
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
	partPath, imgErr := e.registerEditorImage(img.Path, img.Data, img.Format)
	if imgErr != nil {
		if errors.Is(imgErr, errImagePayloadEmpty) {
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
		Rotation:     int64(img.Rotation * rotationDegreeToOOXML),
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
	e.parts.Set(notesPath, []byte(pptxxml.NotesSlide(editorNotesBody(slide))))

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

func editorNotesBody(slide elements.SlideContent) []elements.Paragraph {
	if len(slide.NotesBody) > 0 {
		return slide.NotesBody
	}

	p := elements.NewParagraph()
	p.Runs = append(p.Runs, elements.NewRun(slide.Notes))
	return []elements.Paragraph{p}
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
func renderEditorPlaceholderSpecs(
	e *PresentationEditor,
	slide elements.SlideContent,
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

	for _, override := range slide.PlaceholderOverrides {
		spec := pptxxml.PlaceholderOverrideSpec{
			Index: override.Index,
			Type:  override.Type,
			Text:  override.Text,
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
				return nil, nil, nil, fmt.Errorf("placeholder table %d: %w", override.Index, err)
			}
			spec.Table = tableSpec
		}

		// Handle chart placeholder
		if override.Chart != nil {
			chartSpec := override.Chart.ToChartSpec()
			chartPath := fmt.Sprintf("ppt/charts/chart_ph_%d_%d.xml", slideNumber, override.Index)
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

var errImagePayloadEmpty = errors.New("image has no data or path")

func (e *PresentationEditor) registerEditorImage(pathValue string, data []byte, format string) (string, error) {
	switch {
	case pathValue != "" && len(data) == 0:
		return e.registerImageFromPath(pathValue, format)
	case len(data) > 0:
		return e.RegisterImage(data, format)
	default:
		return "", errImagePayloadEmpty
	}
}

func (e *PresentationEditor) renderBackgroundImageTarget(
	background *elements.SlideBackground,
	currentImageCount int,
) (string, string, error) {
	if background == nil || background.Type != elements.SlideBackgroundPicture || background.PictureFill == nil {
		return "", "", nil
	}

	partPath, err := e.registerEditorImage(
		background.PictureFill.Path,
		background.PictureFill.Data,
		background.PictureFill.Format,
	)
	if err != nil {
		if errors.Is(err, errImagePayloadEmpty) {
			return "", "", nil
		}
		return "", "", fmt.Errorf("read background image: %w", err)
	}

	backgroundRID := fmt.Sprintf("rId%d", currentImageCount+firstImageRelationshipID)
	return backgroundRID, "../media/" + path.Base(partPath), nil
}

func (e *PresentationEditor) renderPlaceholderImageRef(
	override shapes.PlaceholderContent,
	ridIndex int,
) (*pptxxml.ImageRef, string, error) {
	if override.Image == nil {
		return nil, "", nil
	}

	partPath, err := e.registerEditorImage(override.Image.Path, override.Image.Data, override.Image.Format)
	if err != nil {
		if errors.Is(err, errImagePayloadEmpty) {
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

func editorEnsureSlideRelsExistPS(ps *PartStore, slidePart string) error {
	relsPath := common.SlideRelsPartName(slidePart)
	if ps.Has(relsPath) {
		return nil
	}
	return fmt.Errorf("missing slide relationships part %q", relsPath)
}
