package pptx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

// Create builds a valid PPTX with generated slide titles.
func Create(title string, slideCount int) ([]byte, error) {
	if title == "" {
		return nil, fmt.Errorf("presentation title cannot be empty")
	}
	if slideCount < 1 {
		return nil, fmt.Errorf("slide count must be at least 1")
	}

	slides := make([]SlideContent, 0, slideCount)
	for i := 1; i <= slideCount; i++ {
		slideTitle := title
		if i > 1 {
			slideTitle = fmt.Sprintf("Slide %d", i)
		}
		slides = append(slides, NewSlide(slideTitle))
	}

	return CreateWithMetadata(PresentationMetadata{Title: title}, slides)
}

// CreateWithSlides builds a PPTX from caller-provided slide content.
func CreateWithSlides(title string, slides []SlideContent) ([]byte, error) {
	return CreateWithMetadata(PresentationMetadata{Title: title}, slides)
}

// CreateWithMetadata builds a PPTX from metadata and caller-provided slide content.
func CreateWithMetadata(meta PresentationMetadata, slides []SlideContent) ([]byte, error) {
	if meta.Title == "" {
		return nil, fmt.Errorf("presentation title cannot be empty")
	}
	if len(slides) == 0 {
		return nil, fmt.Errorf("at least one slide is required")
	}
	for i, slide := range slides {
		if err := slide.Validate(i + 1); err != nil {
			return nil, err
		}
	}

	if meta.SlideSize.Width == 0 || meta.SlideSize.Height == 0 {
		meta.SlideSize = SlideSize4x3
	}

	buf := bytes.NewBuffer(nil)
	zw := zip.NewWriter(buf)
	count := len(slides)

	if err := writePackageFiles(zw, meta, slides, count); err != nil {
		_ = zw.Close()
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WriteFile is a convenience helper that writes the generated PPTX to disk.
func WriteFile(path string, title string, slides []SlideContent) error {
	data, err := CreateWithMetadata(PresentationMetadata{Title: title, SlideSize: SlideSize4x3}, slides)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func writePackageFiles(zw *zip.Writer, meta PresentationMetadata, slides []SlideContent, slideCount int) error {
	mediaCatalog, err := buildMediaCatalog(slides)
	if err != nil {
		return err
	}
	chartParts := buildChartParts(slides)
	chartBySlide := chartPartBySlide(chartParts)
	notesParts := buildRenderedNotesParts(slides)
	notesTargets := notesTargetBySlide(notesParts)
	hasNotes := len(notesParts) > 0

	files := []struct {
		name    string
		content string
	}{
		{"[Content_Types].xml", pptxxml.ContentTypes(slideCount, mediaCatalog.imageExtensions(), len(chartParts), notesSlideNumbers(notesParts), hasNotes)},
		{"_rels/.rels", pptxxml.RootRelationships()},
		{"ppt/_rels/presentation.xml.rels", pptxxml.PresentationRelationships(slideCount, hasNotes)},
		{"ppt/presentation.xml", pptxxml.Presentation(meta.Title, slideCount, hasNotes, meta.SlideSize.Width, meta.SlideSize.Height)},
		{"ppt/slideLayouts/slideLayout1.xml", pptxxml.SlideLayoutTitleAndContent()},
		{"ppt/slideLayouts/_rels/slideLayout1.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideLayouts/slideLayout2.xml", pptxxml.SlideLayoutTitleOnly()},
		{"ppt/slideLayouts/_rels/slideLayout2.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideLayouts/slideLayout3.xml", pptxxml.SlideLayoutBlank()},
		{"ppt/slideLayouts/_rels/slideLayout3.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideLayouts/slideLayout4.xml", pptxxml.SlideLayoutCenteredTitle()},
		{"ppt/slideLayouts/_rels/slideLayout4.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideLayouts/slideLayout5.xml", pptxxml.SlideLayoutTitleAndBigContent()},
		{"ppt/slideLayouts/_rels/slideLayout5.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideLayouts/slideLayout6.xml", pptxxml.SlideLayoutTwoColumn()},
		{"ppt/slideLayouts/_rels/slideLayout6.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideMasters/slideMaster1.xml", pptxxml.SlideMaster()},
		{"ppt/slideMasters/_rels/slideMaster1.xml.rels", pptxxml.SlideMasterRelationships()},
		{"ppt/theme/theme1.xml", pptxxml.Theme()},
		{"docProps/core.xml", pptxxml.CoreProperties(pptxxml.CorePropertiesInfo{
			Title:       meta.Title,
			Subject:     meta.Subject,
			Creator:     meta.Creator,
			Description: meta.Description,
		})},
		{"docProps/app.xml", pptxxml.AppProperties(slideCount, len(notesParts), meta.SlideSize.Width, meta.SlideSize.Height)},
	}
	if hasNotes {
		files = append(files,
			struct {
				name    string
				content string
			}{"ppt/notesMasters/notesMaster1.xml", pptxxml.NotesMaster()},
			struct {
				name    string
				content string
			}{"ppt/notesMasters/_rels/notesMaster1.xml.rels", pptxxml.NotesMasterRelationships()},
		)
	}

	for _, item := range files {
		if err := writeFile(zw, item.name, item.content); err != nil {
			return err
		}
	}

	if err := writeMediaFiles(zw, mediaCatalog); err != nil {
		return err
	}
	if err := writeChartFiles(zw, chartParts); err != nil {
		return err
	}
	if err := writeNotesFiles(zw, notesParts); err != nil {
		return err
	}

	for i, slide := range slides {
		slideNumber := i + 1

		var tableSpec *pptxxml.TableSpec
		if slide.Table != nil {
			spec, err := buildTableSpec(*slide.Table, slideNumber)
			if err != nil {
				return err
			}
			tableSpec = spec
		}

		imageRefs := make([]pptxxml.ImageRef, 0, len(slide.Images))
		imageTargets := make([]string, 0, len(slide.Images))
		for imageIndex, image := range slide.Images {
			mediaName, ok := mediaCatalog.mediaNameForImage(image)
			if !ok {
				return fmt.Errorf("slide %d image %d was not registered", slideNumber, imageIndex+1)
			}
			relID := fmt.Sprintf("rId%d", imageIndex+2)

			var crop *pptxxml.ImageCropRef
			if image.Crop != (ImageCrop{}) {
				crop = &pptxxml.ImageCropRef{
					Left:   int64(image.Crop.Left * 100000),
					Right:  int64(image.Crop.Right * 100000),
					Top:    int64(image.Crop.Top * 100000),
					Bottom: int64(image.Crop.Bottom * 100000),
				}
			}

			imageRefs = append(imageRefs, pptxxml.ImageRef{
				RelID:      relID,
				Name:       fmt.Sprintf("Picture %d", imageIndex+1),
				X:          image.X,
				Y:          image.Y,
				CX:         image.CX,
				CY:         image.CY,
				Rotation:   int64(image.Rotation * 60000),
				FlipH:      image.FlipH,
				FlipV:      image.FlipV,
				Crop:       crop,
				Shadow:     image.Shadow,
				Reflection: image.Reflection,
			})
			imageTargets = append(imageTargets, fmt.Sprintf("../media/%s", mediaName))
		}

		// Process placeholder images to generate RIDs
		placeholderImageRefs := make(map[int]*pptxxml.ImageRef)
		placeholderTableSpecs := make(map[int]*pptxxml.TableSpec)
		nextRID := len(imageTargets) + 2
		for _, override := range slide.PlaceholderOverrides {
			if override.Image != nil {
				mediaName, ok := mediaCatalog.mediaNameForImage(*override.Image)
				if !ok {
					// Fallback or error? buildMediaCatalog should have caught it.
					continue
				}
				rid := fmt.Sprintf("rId%d", nextRID)
				nextRID++
				imageTargets = append(imageTargets, fmt.Sprintf("../media/%s", mediaName))

				ref := &pptxxml.ImageRef{
					RelID:      rid,
					Name:       "Placeholder Picture",
					FlipH:      override.Image.FlipH,
					FlipV:      override.Image.FlipV,
					Shadow:     override.Image.Shadow,
					Reflection: override.Image.Reflection,
				}
				placeholderImageRefs[override.Index] = ref
			}
			if override.Table != nil {
				spec, err := buildTableSpec(*override.Table, slideNumber)
				if err != nil {
					return err
				}
				placeholderTableSpecs[override.Index] = spec
			}
		}

		var chartFrame *pptxxml.ChartFrame
		var chartRel *pptxxml.ChartRel
		placeholderChartFrames := make(map[int]*pptxxml.ChartFrame)
		placeholderChartRels := make([]pptxxml.ChartRel, 0)

		if parts, ok := chartBySlide[i]; ok {
			partIdx := 0
			// Primary chart (if slide.Chart etc is set)
			if slideChartKindDefined(slide) {
				part := parts[partIdx]
				partIdx++
				rid := fmt.Sprintf("rId%d", nextRID)
				nextRID++
				chartFrame = &pptxxml.ChartFrame{
					RelID: rid,
					X:     part.spec.X,
					Y:     part.spec.Y,
					CX:    part.spec.CX,
					CY:    part.spec.CY,
				}
				chartRel = &pptxxml.ChartRel{
					RID:    rid,
					Target: fmt.Sprintf("../charts/chart%d.xml", part.partNumber),
				}
			}

			// Placeholder charts
			for _, override := range slide.PlaceholderOverrides {
				if override.Chart != nil {
					part := parts[partIdx]
					partIdx++
					rid := fmt.Sprintf("rId%d", nextRID)
					nextRID++
					frame := &pptxxml.ChartFrame{
						RelID: rid,
						X:     part.spec.X,
						Y:     part.spec.Y,
						CX:    part.spec.CX,
						CY:    part.spec.CY,
					}
					placeholderChartFrames[override.Index] = frame
					placeholderChartRels = append(placeholderChartRels, pptxxml.ChartRel{
						RID:    rid,
						Target: fmt.Sprintf("../charts/chart%d.xml", part.partNumber),
					})
				}
			}
		}

		hyperlinkRIDs, hyperlinks, _ := buildSlideHyperlinkRels(slide, nextRID)

		bulletStyles := toXMLBulletParagraphStyles(slide.BulletStyles)
		bulletRuns := toXMLTextRunRows(slide.BulletRuns, hyperlinkRIDs)
		shapeSpecs := toXMLShapeSpecs(slide.Shapes, hyperlinkRIDs)
		connectorSpecs := toXMLConnectorSpecs(slide.Connectors, slide.Shapes)

		placeholderSpecs := make([]pptxxml.PlaceholderOverrideSpec, 0, len(slide.PlaceholderOverrides))
		for _, override := range slide.PlaceholderOverrides {
			var imageRef *pptxxml.ImageRef
			if ref, ok := placeholderImageRefs[override.Index]; ok {
				imageRef = ref
			}
			var tableRef *pptxxml.TableSpec
			if spec, ok := placeholderTableSpecs[override.Index]; ok {
				tableRef = spec
			}
			var chartRef *pptxxml.ChartFrame
			if frame, ok := placeholderChartFrames[override.Index]; ok {
				chartRef = frame
			}

			placeholderSpecs = append(placeholderSpecs, pptxxml.PlaceholderOverrideSpec{
				Index: override.Index,
				Type:  override.Type,
				Text:  override.Text,
				Image: imageRef,
				Table: tableRef,
				Chart: chartRef,
			})
		}

		layoutMode := slideLayoutXMLMode(slide.Layout)
		shapeIDs := calculateShapeIDs(slide)
		animationsXML := slideAnimationsXML(slide, shapeIDs)

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
			bulletStyles,
			bulletRuns,
			contentStyle,
			tableSpec,
			chartFrame,
			imageRefs,
			shapeSpecs,
			connectorSpecs,
			placeholderSpecs,
			slideTransitionXML(slide),
			animationsXML,
		)
		slidePath := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
		if err := writeFile(zw, slidePath, slideXML); err != nil {
			return err
		}

		relsPath := fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", slideNumber)
		if err := writeFile(
			zw,
			relsPath,
			pptxxml.SlideRelationshipsWithMultiCharts(
				slideLayoutTarget(slide.Layout),
				imageTargets,
				chartRel,
				placeholderChartRels,
				notesTargets[slideNumber],
				hyperlinks,
			),
		); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(zw *zip.Writer, path string, content string) error {
	w, err := zw.Create(path)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(content))
	return err
}

func writeMediaFiles(zw *zip.Writer, catalog *mediaCatalog) error {
	for _, asset := range catalog.ordered {
		path := fmt.Sprintf("ppt/media/%s", asset.mediaName)
		w, err := zw.Create(path)
		if err != nil {
			return err
		}
		if _, err := w.Write(asset.data); err != nil {
			return err
		}
	}
	return nil
}
