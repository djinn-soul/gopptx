package presentation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

func renderSlides(
	pw *pptxxml.PackageWriter,
	meta PresentationMetadata,
	slides []elements.SlideContent,
	mediaCatalog *media.MediaCatalog,
	chartBySlide map[int][]chartPart,
	notesTargets map[int]string,
	masterCount int,
) error {
	for i, slide := range slides {
		num := i + 1
		builder := &slidePartBuilder{
			num:     num,
			catalog: mediaCatalog,
			targets: make([]string, 0),
			ridNext: 2,
		}

		parts, err := builder.build(i, slide, chartBySlide)
		if err != nil {
			return err
		}

		hyperlinkRIDs, hyperlinks, _ := elements.BuildSlideHyperlinkRels(slide, builder.ridNext)

		slideXML := pptxxml.SlideWithLayout(
			elements.SlideLayoutXMLMode(slide.Layout),
			parts.title,
			slide.Bullets,
			elements.ToXMLBulletParagraphStyles(slide.BulletStyles),
			elements.ToXMLTextRunRows(slide.BulletRuns, hyperlinkRIDs),
			parts.contentStyle,
			parts.table,
			parts.chartFrame,
			parts.imageRefs,
			shapes.ToXMLShapeSpecs(slide.Shapes, hyperlinkRIDs),
			shapes.ToXMLConnectorSpecs(slide.Connectors, slide.Shapes),
			parts.placeholders,
			elements.ToXMLBackgroundSpec(slide.Background, parts.backgroundRID),
			parts.transitionXML,
			elements.SlideAnimationsXML(slide, elements.CalculateShapeIDs(slide)),
			slide.ShowSlideNumber,
			func() string {
				if slide.FooterText != "" {
					return slide.FooterText
				}
				return meta.FooterText
			}(),
			meta.ShowDateTime,
			meta.SlideSize.Width,
			meta.SlideSize.Height,
		)

		layoutTarget := elements.SlideLayoutTarget(slide.Layout)
		if masterCount > 1 {
			masterNum := (i % masterCount) + 1
			layoutTarget = layoutTargetForMaster(layoutTarget, masterNum)
		}

		pw.AddPart(fmt.Sprintf("ppt/slides/slide%d.xml", num), slideXML)
		pw.AddPart(fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", num), pptxxml.SlideRelationshipsWithMultiCharts(
			layoutTarget,
			builder.targets,
			parts.chartRel,
			parts.placeholderChartRels,
			notesTargets[num],
			hyperlinks,
		))
	}
	return nil
}

func layoutTargetForMaster(baseTarget string, masterNum int) string {
	if masterNum <= 1 {
		return baseTarget
	}
	const prefix = "../slideLayouts/slideLayout"
	const suffix = ".xml"
	if !strings.HasPrefix(baseTarget, prefix) || !strings.HasSuffix(baseTarget, suffix) {
		return baseTarget
	}
	n, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(baseTarget, prefix), suffix))
	if err != nil || n < 1 {
		return baseTarget
	}
	globalLayout := (masterNum-1)*6 + n
	return fmt.Sprintf("%s%d%s", prefix, globalLayout, suffix)
}

type slidePartBuilder struct {
	num     int
	catalog *media.MediaCatalog
	targets []string
	ridNext int
}

type slideParts struct {
	title                pptxxml.TitleSpec
	contentStyle         pptxxml.ContentStyleSpec
	table                *pptxxml.TableSpec
	imageRefs            []pptxxml.ImageRef
	backgroundRID        string
	transitionXML        string
	placeholders         []pptxxml.PlaceholderOverrideSpec
	chartFrame           *pptxxml.ChartFrame
	chartRel             *pptxxml.ChartRel
	placeholderChartRels []pptxxml.ChartRel
}

func (b *slidePartBuilder) build(
	idx int,
	slide elements.SlideContent,
	chartBySlide map[int][]chartPart,
) (*slideParts, error) {
	p := &slideParts{
		title:        b.buildTitleSpec(slide),
		contentStyle: b.buildContentStyleSpec(slide),
	}

	if slide.Table != nil {
		spec, err := slide.Table.ToTableSpec(b.num)
		if err != nil {
			return nil, err
		}
		p.table = spec
	}

	imageRefs, err := b.mapImages(slide.Images)
	if err != nil {
		return nil, err
	}
	p.imageRefs = imageRefs

	p.backgroundRID = b.mapBackground(slide.Background)
	p.placeholderChartRels = make([]pptxxml.ChartRel, 0)

	b.handleTransitionSound(&slide)
	p.transitionXML = elements.SlideTransitionXML(slide)
	if err := b.mapPlaceholders(&p.placeholders, &p.placeholderChartRels, slide.PlaceholderOverrides); err != nil {
		return nil, err
	}

	if parts, ok := chartBySlide[idx]; ok {
		if err := b.mapCharts(p, parts, slide); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (b *slidePartBuilder) nextRID() string {
	rid := fmt.Sprintf("rId%d", b.ridNext)
	b.ridNext++
	return rid
}

func (b *slidePartBuilder) mapImages(images []shapes.Image) ([]pptxxml.ImageRef, error) {
	refs := make([]pptxxml.ImageRef, 0, len(images))
	for i, img := range images {
		mediaName, ok := b.catalog.MediaNameForImage(img)
		if !ok {
			return nil, fmt.Errorf("slide %d image %d not registered", b.num, i+1)
		}
		rid := b.nextRID()
		refs = append(refs, pptxxml.ImageRef{
			RelID:        rid,
			Name:         fmt.Sprintf("Picture %d", i+1),
			X:            img.X.Emu(),
			Y:            img.Y.Emu(),
			CX:           img.CX.Emu(),
			CY:           img.CY.Emu(),
			Rotation:     int64(img.Rotation * 60000),
			FlipH:        img.FlipH,
			FlipV:        img.FlipV,
			Crop:         mapCrop(img.Crop),
			Shadow:       img.Shadow,
			Reflection:   img.Reflection,
			AltText:      img.AltText,
			IsDecorative: img.IsDecorative,
		})
		b.targets = append(b.targets, fmt.Sprintf("../media/%s", mediaName))
	}
	return refs, nil
}

func (b *slidePartBuilder) mapBackground(bg *elements.SlideBackground) string {
	if bg != nil && bg.Type == elements.SlideBackgroundPicture && bg.PictureFill != nil {
		if mediaName, ok := b.catalog.MediaNameForImage(*bg.PictureFill); ok {
			rid := b.nextRID()
			b.targets = append(b.targets, fmt.Sprintf("../media/%s", mediaName))
			return rid
		}
	}
	return ""
}

func (b *slidePartBuilder) handleTransitionSound(slide *elements.SlideContent) {
	if slide.Transition != nil {
		if opt, ok := slide.Transition.(transitions.TransitionOptions); ok && opt.Sound != nil &&
			strings.HasPrefix(opt.Sound.RelID, "file:") {
			path := strings.TrimPrefix(opt.Sound.RelID, "file:")
			soundMedia := shapes.Image{Path: path}
			if mediaName, ok := b.catalog.MediaNameForImage(soundMedia); ok {
				rid := b.nextRID()
				b.targets = append(b.targets, fmt.Sprintf("../media/%s", mediaName))
				opt.Sound.RelID = rid
				slide.Transition = opt
			}
		}
	}
}

func (b *slidePartBuilder) mapPlaceholders(
	specs *[]pptxxml.PlaceholderOverrideSpec,
	chartRels *[]pptxxml.ChartRel,
	overrides []shapes.PlaceholderContent,
) error {
	imageRefs := make(map[int]*pptxxml.ImageRef)
	tableSpecs := make(map[int]*pptxxml.TableSpec)

	for _, override := range overrides {
		if override.Image != nil {
			if mediaName, ok := b.catalog.MediaNameForImage(*override.Image); ok {
				rid := b.nextRID()
				b.targets = append(b.targets, fmt.Sprintf("../media/%s", mediaName))
				imageRefs[override.Index] = &pptxxml.ImageRef{
					RelID:      rid,
					Name:       "Placeholder Picture",
					FlipH:      override.Image.FlipH,
					FlipV:      override.Image.FlipV,
					Shadow:     override.Image.Shadow,
					Reflection: override.Image.Reflection,
				}
			}
		}
		if override.Table != nil {
			spec, err := override.Table.ToTableSpec(b.num)
			if err != nil {
				return err
			}
			tableSpecs[override.Index] = spec
		}
	}

	for _, override := range overrides {
		*specs = append(*specs, pptxxml.PlaceholderOverrideSpec{
			Index: override.Index,
			Type:  override.Type,
			Text:  override.Text,
			Image: imageRefs[override.Index],
			Table: tableSpecs[override.Index],
		})
	}
	return nil
}

func (b *slidePartBuilder) mapCharts(p *slideParts, parts []chartPart, slide elements.SlideContent) error {
	partIdx := 0
	if slideChartKindDefined(slide) {
		part := parts[partIdx]
		partIdx++
		rid := b.nextRID()
		p.chartFrame = &pptxxml.ChartFrame{
			RelID: rid,
			X:     part.spec.X,
			Y:     part.spec.Y,
			CX:    part.spec.CX,
			CY:    part.spec.CY,
		}
		p.chartRel = &pptxxml.ChartRel{
			RID:    rid,
			Target: fmt.Sprintf("../charts/chart%d.xml", part.partNumber),
		}
	}

	for i, override := range slide.PlaceholderOverrides {
		if override.Chart != nil {
			if partIdx >= len(parts) {
				return fmt.Errorf("slide %d: missing chart part for placeholder index %d", b.num, override.Index)
			}
			part := parts[partIdx]
			partIdx++
			rid := b.nextRID()
			frame := &pptxxml.ChartFrame{
				RelID: rid,
				X:     part.spec.X,
				Y:     part.spec.Y,
				CX:    part.spec.CX,
				CY:    part.spec.CY,
			}
			p.placeholders[i].Chart = frame
			p.placeholderChartRels = append(p.placeholderChartRels, pptxxml.ChartRel{
				RID:    rid,
				Target: fmt.Sprintf("../charts/chart%d.xml", part.partNumber),
			})
		}
	}
	return nil
}

func (b *slidePartBuilder) buildTitleSpec(slide elements.SlideContent) pptxxml.TitleSpec {
	return pptxxml.TitleSpec{
		Text:      slide.Title,
		SizePt:    slide.TitleSize,
		Color:     slide.TitleColor,
		Bold:      slide.TitleBold,
		Italic:    slide.TitleItalic,
		Underline: slide.TitleUnderline,
		Align:     slide.TitleAlign,
		Font:      slide.TitleFont,
	}
}

func (b *slidePartBuilder) buildContentStyleSpec(slide elements.SlideContent) pptxxml.ContentStyleSpec {
	return pptxxml.ContentStyleSpec{
		SizePt:    slide.ContentSize,
		Color:     slide.ContentColor,
		Bold:      slide.ContentBold,
		Italic:    slide.ContentItalic,
		Underline: slide.ContentUnderline,
		VAlign:    slide.ContentVAlign,
	}
}

func mapCrop(crop shapes.ImageCrop) *pptxxml.ImageCropRef {
	if crop == (shapes.ImageCrop{}) {
		return nil
	}
	return &pptxxml.ImageCropRef{
		Left:   int64(crop.Left * 100000),
		Right:  int64(crop.Right * 100000),
		Top:    int64(crop.Top * 100000),
		Bottom: int64(crop.Bottom * 100000),
	}
}
