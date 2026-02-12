package presentation

import (
	"archive/zip"
	"encoding/xml"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/notes"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// PresentationMetadata defines non-content properties of a PPTX.
type PresentationMetadata struct {
	common.PresentationMetadata
	Theme  *styling.Theme
	Master *elements.SlideMaster
}

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

	// Register master images in the media catalog and build targets/specs.
	masterImageTargets, masterImageRefs := buildMasterImageInfo(meta.Master, mediaCatalog)
	masterSpec := mapMasterToSpec(meta.Master, masterImageRefs)

	chartParts := BuildChartParts(slides)
	chartBySlide := chartPartBySlide(chartParts)
	notesParts := notes.BuildRenderedNotesParts(slides)
	notesTargets := notes.NotesTargetBySlide(notesParts)
	hasNotes := len(notesParts) > 0

	files := []struct {
		name    string
		content string
	}{
		{"[Content_Types].xml", pptxxml.ContentTypes(slideCount, mediaCatalog.ImageExtensions(), len(chartParts), notes.NotesSlideNumbers(notesParts), hasNotes, len(meta.CustomXML))},
		{"_rels/.rels", pptxxml.RootRelationships()},
		{"ppt/_rels/presentation.xml.rels", pptxxml.PresentationRelationships(slideCount, hasNotes, len(meta.CustomXML))},
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
		{"ppt/slideMasters/slideMaster1.xml", pptxxml.SlideMaster(masterSpec)},
		{"ppt/slideMasters/_rels/slideMaster1.xml.rels", pptxxml.SlideMasterRelationships(masterImageTargets)},
		{"ppt/theme/theme1.xml", pptxxml.Theme(mapThemeToSpec(meta.Theme))},
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
				RelID: relID,
				Name:  fmt.Sprintf("Picture %d", imageIndex+1),
				X:     image.X.Emu(),
				Y:     image.Y.Emu(),
				CX:    image.CX.Emu(),
				CY:    image.CY.Emu(),

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

		backgroundRID := ""
		if slide.Background != nil && slide.Background.Type == elements.SlideBackgroundPicture && slide.Background.PictureFill != nil {
			mediaName, ok := mediaCatalog.MediaNameForImage(*slide.Background.PictureFill)
			if ok {
				backgroundRID = fmt.Sprintf("rId%d", len(imageTargets)+2)
				imageTargets = append(imageTargets, fmt.Sprintf("../media/%s", mediaName))
			}
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

		backgroundSpec := elements.ToXMLBackgroundSpec(slide.Background, backgroundRID)

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
			backgroundSpec,
			elements.SlideTransitionXML(slide),
			animationsXML,
			slide.ShowSlideNumber,
			meta.FooterText,
			meta.ShowDateTime,
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

	// Write custom XML parts.
	for i, part := range meta.CustomXML {
		// Ensure part content is well-formed XML.
		if err := xml.Unmarshal([]byte(part.Content), new(interface{})); err != nil {
			return fmt.Errorf("custom XML part %d contains invalid XML: %w", i+1, err)
		}
		path := fmt.Sprintf("customXml/item%d.xml", i+1)
		if err := common.WriteFile(zw, path, part.Content); err != nil {
			return err
		}

		// Generate itemProps.
		itemID, err := common.NewGUID()
		if err != nil {
			return fmt.Errorf("generate custom XML itemID for part %d: %w", i+1, err)
		}
		propsContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<ds:datastoreItem ds:itemID="%s" xmlns:ds="http://schemas.openxmlformats.org/officeDocument/2006/customXml">
<ds:schemaRefs/>
</ds:datastoreItem>`, itemID)
		propsPath := fmt.Sprintf("customXml/itemProps%d.xml", i+1)
		if err := common.WriteFile(zw, propsPath, propsContent); err != nil {
			return err
		}
	}

	return nil
}

func mapThemeToSpec(theme *styling.Theme) *pptxxml.ThemeSpec {
	if theme == nil {
		return nil
	}
	spec := &pptxxml.ThemeSpec{
		Name: theme.Name,
		Colors: pptxxml.ColorSchemeSpec{
			Name:     theme.Colors.Name,
			Dk1:      theme.Colors.Dk1,
			Lt1:      theme.Colors.Lt1,
			Dk2:      theme.Colors.Dk2,
			Lt2:      theme.Colors.Lt2,
			Accent1:  theme.Colors.Accent1,
			Accent2:  theme.Colors.Accent2,
			Accent3:  theme.Colors.Accent3,
			Accent4:  theme.Colors.Accent4,
			Accent5:  theme.Colors.Accent5,
			Accent6:  theme.Colors.Accent6,
			Hlink:    theme.Colors.Hlink,
			FolHlink: theme.Colors.FolHlink,
		},
		Fonts: pptxxml.FontSchemeSpec{
			Name:      theme.Fonts.Name,
			MajorFont: theme.Fonts.MajorFont,
			MinorFont: theme.Fonts.MinorFont,
		},
	}
	return spec
}

func mapMasterToSpec(master *elements.SlideMaster, imageRefs []pptxxml.ImageRef) *pptxxml.SlideMasterSpec {
	if master == nil {
		return nil
	}
	spec := &pptxxml.SlideMasterSpec{
		FooterText: master.FooterText,
		Images:     imageRefs,
	}
	if master.ColorMapping != nil {
		spec.ColorMapping = &pptxxml.ColorMappingSpec{
			BG1: master.ColorMapping.BG1,
			TX1: master.ColorMapping.TX1,
		}
	}
	if master.Background != nil {
		spec.Background = elements.ToXMLBackgroundSpec(master.Background, "")
	}
	// Map shapes (no hyperlinks on master shapes).
	masterShapes := make([]shapes.Shape, 0, len(master.Shapes))
	for _, sd := range master.Shapes {
		masterShapes = append(masterShapes, sd.ToShape())
	}
	spec.Shapes = shapes.ToXMLShapeSpecs(masterShapes, nil)
	return spec
}

// buildMasterImageInfo registers master images and returns relationship targets and ImageRef specs.
func buildMasterImageInfo(master *elements.SlideMaster, catalog *media.MediaCatalog) ([]string, []pptxxml.ImageRef) {
	if master == nil || len(master.Images) == 0 {
		return nil, nil
	}
	targets := make([]string, 0, len(master.Images))
	refs := make([]pptxxml.ImageRef, 0, len(master.Images))
	for i, img := range master.Images {
		mediaName, err := catalog.RegisterImage(img)
		if err != nil {
			continue // skip unresolved master images
		}
		// Master image RIDs start at rId8 (rId1-6 are layouts, rId7 is theme).
		relID := fmt.Sprintf("rId%d", 8+i)
		targets = append(targets, fmt.Sprintf("../media/%s", mediaName))
		refs = append(refs, pptxxml.ImageRef{
			RelID: relID,
			Name:  fmt.Sprintf("Master Picture %d", i+1),
			X:     img.X.Emu(),
			Y:     img.Y.Emu(),
			CX:    img.CX.Emu(),
			CY:    img.CY.Emu(),

			AltText:      img.AltText,
			IsDecorative: img.IsDecorative,
		})
	}
	return targets, refs
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
