package presentation

import (
	"archive/zip"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/notes"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// PresentationMetadata defines non-content properties of a PPTX.
type PresentationMetadata = common.PresentationMetadata

// SlideSize defines presentation dimensions in EMUs.
type SlideSize = common.SlideSize

// Default slide sizes.
var (
	SlideSize4x3  = common.SlideSize4x3
	SlideSize16x9 = common.SlideSize16x9
)

type chartPart struct {
	slideIndex int
	partNumber int
	spec       pptxxml.ChartSpec
}

func BuildChartParts(slides []elements.SlideContent) []chartPart {
	out := make([]chartPart, 0)
	for i, slide := range slides {
		spec, ok := slideChartSpec(slide)
		if ok {
			out = append(out, chartPart{
				slideIndex: i,
				partNumber: len(out) + 1,
				spec:       *spec,
			})
		}

		for _, override := range slide.PlaceholderOverrides {
			if override.Chart != nil {
				out = append(out, chartPart{
					slideIndex: i,
					partNumber: len(out) + 1,
					spec:       *override.Chart.ToChartSpec(),
				})
			}
		}
	}
	return out
}

func chartPartBySlide(parts []chartPart) map[int][]chartPart {
	bySlide := make(map[int][]chartPart, len(parts))
	for _, part := range parts {
		bySlide[part.slideIndex] = append(bySlide[part.slideIndex], part)
	}
	return bySlide
}

func writeChartFiles(zw *zip.Writer, parts []chartPart) error {
	for _, part := range parts {
		path := fmt.Sprintf("ppt/charts/chart%d.xml", part.partNumber)
		content := pptxxml.ChartPartXML(&part.spec)
		if err := common.WriteFile(zw, path, content); err != nil {
			return err
		}
	}
	return nil
}

func slideChartSpec(slide elements.SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.Chart != nil {
		return slide.Chart.ToChartSpec(), true
	}
	if slide.BarHorizontal != nil {
		return slide.BarHorizontal.ToChartSpec(), true
	}
	if slide.BarStacked != nil {
		return slide.BarStacked.ToChartSpec(), true
	}
	if slide.BarStacked100 != nil {
		return slide.BarStacked100.ToChartSpec(), true
	}
	if slide.Line != nil {
		return slide.Line.ToChartSpec(), true
	}
	if slide.LineMarkers != nil {
		return slide.LineMarkers.ToChartSpec(), true
	}
	if slide.LineStacked != nil {
		return slide.LineStacked.ToChartSpec(), true
	}
	if slide.Scatter != nil {
		return slide.Scatter.ToChartSpec(), true
	}
	if slide.Area != nil {
		return slide.Area.ToChartSpec(), true
	}
	if slide.AreaStacked != nil {
		return slide.AreaStacked.ToChartSpec(), true
	}
	if slide.AreaStacked100 != nil {
		return slide.AreaStacked100.ToChartSpec(), true
	}
	if slide.Pie != nil {
		return slide.Pie.ToChartSpec(), true
	}
	if slide.Doughnut != nil {
		return slide.Doughnut.ToChartSpec(), true
	}
	if slide.Bubble != nil {
		return slide.Bubble.ToChartSpec(), true
	}
	if slide.Radar != nil {
		return slide.Radar.ToChartSpec(), true
	}
	if slide.RadarFilled != nil {
		return slide.RadarFilled.ToChartSpec(), true
	}
	if slide.StockHLC != nil {
		return slide.StockHLC.ToChartSpec(), true
	}
	if slide.StockOHLC != nil {
		return slide.StockOHLC.ToChartSpec(), true
	}
	if slide.Combo != nil {
		return slide.Combo.ToChartSpec(), true
	}
	return nil, false
}

func slideChartKindDefined(slide elements.SlideContent) bool {
	_, ok := slideChartSpec(slide)
	return ok
}

func WritePackageFiles(zw *zip.Writer, meta PresentationMetadata, slides []elements.SlideContent, slideCount int) error {
	mediaCatalog, err := media.BuildMediaCatalog(slides)
	if err != nil {
		return err
	}
	chartParts := BuildChartParts(slides)
	chartBySlide := chartPartBySlide(chartParts)
	notesParts := notes.BuildRenderedNotesParts(slides)
	notesTargets := notes.NotesTargetBySlide(notesParts)
	hasNotes := len(notesParts) > 0

	files := []struct {
		name    string
		content string
	}{
		{"[Content_Types].xml", pptxxml.ContentTypes(slideCount, mediaCatalog.ImageExtensions(), len(chartParts), notes.NotesSlideNumbers(notesParts), hasNotes)},
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
		if err := common.WriteFile(zw, item.name, item.content); err != nil {
			return err
		}
	}

	if err := writeMediaFiles(zw, mediaCatalog); err != nil {
		return err
	}
	if err := writeChartFiles(zw, chartParts); err != nil {
		return err
	}
	if err := notes.WriteNotesFiles(zw, notesParts); err != nil {
		return err
	}

	for i, slide := range slides {
		slideNumber := i + 1

		var tableSpec *pptxxml.TableSpec
		if slide.Table != nil {
			spec, err := slide.Table.ToTableSpec(slideNumber)
			if err != nil {
				return err
			}
			tableSpec = spec
		}

		imageRefs := make([]pptxxml.ImageRef, 0, len(slide.Images))
		imageTargets := make([]string, 0, len(slide.Images))
		for imageIndex, image := range slide.Images {
			mediaName, ok := mediaCatalog.MediaNameForImage(image)
			if !ok {
				return fmt.Errorf("slide %d image %d was not registered", slideNumber, imageIndex+1)
			}
			relID := fmt.Sprintf("rId%d", imageIndex+2)

			var crop *pptxxml.ImageCropRef
			if image.Crop != (shapes.ImageCrop{}) {
				crop = &pptxxml.ImageCropRef{
					Left:   int64(image.Crop.Left * 100000),
					Right:  int64(image.Crop.Right * 100000),
					Top:    int64(image.Crop.Top * 100000),
					Bottom: int64(image.Crop.Bottom * 100000),
				}
			}

			imageRefs = append(imageRefs, pptxxml.ImageRef{
				RelID:        relID,
				Name:         fmt.Sprintf("Picture %d", imageIndex+1),
				X:            image.X,
				Y:            image.Y,
				CX:           image.CX,
				CY:           image.CY,
				Rotation:     int64(image.Rotation * 60000),
				FlipH:        image.FlipH,
				FlipV:        image.FlipV,
				Crop:         crop,
				Shadow:       image.Shadow,
				Reflection:   image.Reflection,
				AltText:      image.AltText,
				IsDecorative: image.IsDecorative,
			})
			imageTargets = append(imageTargets, fmt.Sprintf("../media/%s", mediaName))
		}

		// Process placeholder images to generate RIDs
		placeholderImageRefs := make(map[int]*pptxxml.ImageRef)
		placeholderTableSpecs := make(map[int]*pptxxml.TableSpec)
		currentRID := len(imageTargets) + 2
		for _, override := range slide.PlaceholderOverrides {
			if override.Image != nil {
				mediaName, ok := mediaCatalog.MediaNameForImage(*override.Image)
				if !ok {
					continue
				}
				rid := fmt.Sprintf("rId%d", currentRID)
				currentRID++
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
				spec, err := override.Table.ToTableSpec(slideNumber)
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
			if slideChartKindDefined(slide) {
				part := parts[partIdx]
				partIdx++
				rid := fmt.Sprintf("rId%d", currentRID)
				currentRID++
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

			for _, override := range slide.PlaceholderOverrides {
				if override.Chart != nil {
					if partIdx >= len(parts) {
						return fmt.Errorf("slide %d: missing chart part for placeholder index %d", slideNumber, override.Index)
					}
					part := parts[partIdx]
					partIdx++
					rid := fmt.Sprintf("rId%d", currentRID)
					currentRID++
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

		hyperlinkRIDs, hyperlinks, _ := elements.BuildSlideHyperlinkRels(slide, currentRID)

		bulletStyles := elements.ToXMLBulletParagraphStyles(slide.BulletStyles)
		bulletRuns := elements.ToXMLTextRunRows(slide.BulletRuns, hyperlinkRIDs)
		shapeSpecs := shapes.ToXMLShapeSpecs(slide.Shapes, hyperlinkRIDs)
		connectorSpecs := shapes.ToXMLConnectorSpecs(slide.Connectors, slide.Shapes)

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

		layoutMode := elements.SlideLayoutXMLMode(slide.Layout)
		shapeIDs := elements.CalculateShapeIDs(slide)
		animationsXML := elements.SlideAnimationsXML(slide, shapeIDs)

		titleSpec := pptxxml.TitleSpec{
			Text:      slide.Title,
			SizePt:    slide.TitleSize,
			Color:     slide.TitleColor,
			Bold:      slide.TitleBold,
			Italic:    slide.TitleItalic,
			Underline: slide.TitleUnderline,
			Align:     slide.TitleAlign,
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
			bulletStyles,
			bulletRuns,
			contentStyle,
			tableSpec,
			chartFrame,
			imageRefs,
			shapeSpecs,
			connectorSpecs,
			placeholderSpecs,
			slide.BackgroundColor,
			elements.SlideTransitionXML(slide),
			animationsXML,
			meta.SlideSize.Width,
			meta.SlideSize.Height,
		)
		slidePath := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
		if err := common.WriteFile(zw, slidePath, slideXML); err != nil {
			return err
		}

		relsPath := fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", slideNumber)
		if err := common.WriteFile(
			zw,
			relsPath,
			pptxxml.SlideRelationshipsWithMultiCharts(
				elements.SlideLayoutTarget(slide.Layout),
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

func writeMediaFiles(zw *zip.Writer, catalog *media.MediaCatalog) error {
	for _, asset := range catalog.Assets() {
		path := fmt.Sprintf("ppt/media/%s", asset.MediaName())
		if err := common.WriteBinaryFile(zw, path, asset.Data()); err != nil {
			return err
		}
	}
	return nil
}
