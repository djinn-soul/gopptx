package presentation

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

type slidePartBuilder struct {
	num     int
	catalog *media.Catalog
	ridNext int
	targets []string
}

func newSlidePartBuilder(num int, catalog *media.Catalog) *slidePartBuilder {
	return &slidePartBuilder{
		num:     num,
		catalog: catalog,
		ridNext: 2, // rId1 is always slideLayout
		targets: make([]string, 0),
	}
}

func (b *slidePartBuilder) build(
	idx int,
	slide elements.SlideContent,
	chartBySlide map[int][]ChartPart,
	smartArtBySlide map[int][]SmartArtPart,
) (*slideParts, error) {
	p := &slideParts{
		title:        b.buildTitleSpec(slide),
		contentStyle: b.buildContentStyleSpec(slide),
		// Placeholder chart RIDs must be reserved before placeholder media is mapped
		// so media RIDs remain contiguous from rId2 in relationship output.
		placeholderChartRels: b.allocatePlaceholderChartRels(slide.PlaceholderOverrides),
	}

	if slide.Table != nil {
		spec, err := slide.Table.ToTableSpec(b.num)
		if err != nil {
			return nil, err
		}
		p.table = spec
	}

	imageRefs, mapErr := b.mapImages(slide.Images)
	if mapErr != nil {
		return nil, mapErr
	}
	p.imageRefs = imageRefs

	p.backgroundRID = b.mapBackground(slide.Background)
	p.smartArtFrames = make([]pptxxml.SmartArtFrame, 0)
	p.smartArtRels = make([]pptxxml.SmartArtRel, 0)
	p.transitionXML = b.mapTransition(slide)
	placeholders, err := b.mapPlaceholders(slide.PlaceholderOverrides, p.placeholderChartRels)
	if err != nil {
		return nil, err
	}
	p.placeholders = placeholders

	if saList, ok := smartArtBySlide[idx]; ok {
		for _, saPart := range saList {
			dmRID := b.nextRID()
			loRID := b.nextRID()
			qsRID := b.nextRID()
			csRID := b.nextRID()
			drRID := b.nextRID()

			p.smartArtFrames = append(p.smartArtFrames, pptxxml.SmartArtFrame{
				DataRelID:    dmRID,
				LayoutRelID:  loRID,
				StyleRelID:   qsRID,
				ColorRelID:   csRID,
				X:            saPart.spec.X,
				Y:            saPart.spec.Y,
				CX:           saPart.spec.CX,
				CY:           saPart.spec.CY,
				AltText:      saPart.spec.AltText,
				IsDecorative: saPart.spec.IsDecorative,
			})

			p.smartArtRels = append(p.smartArtRels,
				pptxxml.SmartArtRel{
					RID:    dmRID,
					Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/diagramData",
					Target: fmt.Sprintf("../diagrams/data%d.xml", saPart.partNumber),
				},
				pptxxml.SmartArtRel{
					RID:    loRID,
					Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/diagramLayout",
					Target: fmt.Sprintf("../diagrams/layout%d.xml", saPart.partNumber),
				},
				pptxxml.SmartArtRel{
					RID:    qsRID,
					Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/diagramQuickStyle",
					Target: fmt.Sprintf("../diagrams/quickStyle%d.xml", saPart.partNumber),
				},
				pptxxml.SmartArtRel{
					RID:    csRID,
					Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/diagramColors",
					Target: fmt.Sprintf("../diagrams/colors%d.xml", saPart.partNumber),
				},
				pptxxml.SmartArtRel{
					RID:    drRID,
					Type:   "http://schemas.microsoft.com/office/2007/relationships/diagramDrawing",
					Target: fmt.Sprintf("../diagrams/drawing%d.xml", saPart.partNumber),
				},
			)
		}
	}

	if chartList, ok := chartBySlide[idx]; ok && len(chartList) > 0 {
		listIdx := 0
		if slideChartKindDefined(slide) {
			rid := b.nextRID()
			p.chartRel = &pptxxml.ChartRel{
				RID:    rid,
				Target: fmt.Sprintf("../charts/chart%d.xml", chartList[0].partNumber),
			}
			p.chartFrame = &pptxxml.ChartFrame{
				RelID:        rid,
				X:            chartList[0].spec.X,
				Y:            chartList[0].spec.Y,
				CX:           chartList[0].spec.CX,
				CY:           chartList[0].spec.CY,
				AltText:      chartList[0].spec.AltText,
				IsDecorative: chartList[0].spec.IsDecorative,
			}
			listIdx = 1
		}

		// Remaining charts are from placeholders
		for ; listIdx < len(chartList); listIdx++ {
			placeholderChartIdx := listIdx
			if slideChartKindDefined(slide) {
				placeholderChartIdx--
			}
			if placeholderChartIdx >= len(p.placeholderChartRels) {
				return nil, fmt.Errorf("slide %d: missing placeholder chart relationship slot", b.num)
			}
			p.placeholderChartRels[placeholderChartIdx].Target = fmt.Sprintf(
				"../charts/chart%d.xml",
				chartList[listIdx].partNumber,
			)
		}
	}
	for _, rel := range p.placeholderChartRels {
		if strings.TrimSpace(rel.Target) == "" {
			return nil, fmt.Errorf("slide %d: missing chart part for placeholder chart", b.num)
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
			Shadow:       img.Shadow,
			Reflection:   img.Reflection,
			AltText:      img.AltText,
			IsDecorative: img.IsDecorative,
			Crop:         mapToXMLCrop(img.Crop),
		})
		b.targets = append(b.targets, fmt.Sprintf("../media/%s", mediaName))
	}
	return refs, nil
}

func (b *slidePartBuilder) mapPlaceholders(overrides []shapes.PlaceholderContent, placeholderChartRels []pptxxml.ChartRel) ([]pptxxml.PlaceholderOverrideSpec, error) {
	specs := make([]pptxxml.PlaceholderOverrideSpec, 0, len(overrides))
	chartIdx := 0
	for _, o := range overrides {
		// Validation for create-path (where builder is used)
		if o.Target != nil {
			if o.Target.Type == "" && o.Target.Index == 0 && o.Target.Name != "" {
				return nil, fmt.Errorf("name-only target %q is not supported for create-path overrides", o.Target.Name)
			}
		}

		spec := pptxxml.PlaceholderOverrideSpec{
			Index: o.Index,
			Type:  o.Type,
			Text:  o.Text,
		}
		if o.Override != nil {
			spec.X = mapOptionalLength(o.Override.X)
			spec.Y = mapOptionalLength(o.Override.Y)
			spec.CX = mapOptionalLength(o.Override.CX)
			spec.CY = mapOptionalLength(o.Override.CY)
			spec.TextStyle = mapPlaceholderTextStyle(o.Override.TextStyle)
		}
		if o.Image != nil {
			mediaName, ok := b.catalog.MediaNameForImage(*o.Image)
			if ok {
				rid := b.nextRID()
				spec.Image = &pptxxml.ImageRef{
					RelID:      rid,
					Name:       "Placeholder Picture",
					X:          o.Image.X.Emu(),
					Y:          o.Image.Y.Emu(),
					CX:         o.Image.CX.Emu(),
					CY:         o.Image.CY.Emu(),
					Rotation:   int64(o.Image.Rotation * 60000),
					FlipH:      o.Image.FlipH,
					FlipV:      o.Image.FlipV,
					Shadow:     o.Image.Shadow,
					Reflection: o.Image.Reflection,
					Crop:       mapToXMLCrop(o.Image.Crop),
				}
				b.targets = append(b.targets, fmt.Sprintf("../media/%s", mediaName))
			}
		}
		if o.Table != nil {
			tableSpec, err := o.Table.ToTableSpec(b.num)
			if err != nil {
				return nil, err
			}
			spec.Table = tableSpec
		}
		if o.Chart != nil && chartIdx < len(placeholderChartRels) {
			chartSpec := o.Chart.ToChartSpec()
			spec.Chart = &pptxxml.ChartFrame{
				RelID:        placeholderChartRels[chartIdx].RID,
				X:            chartSpec.X,
				Y:            chartSpec.Y,
				CX:           chartSpec.CX,
				CY:           chartSpec.CY,
				AltText:      chartSpec.AltText,
				IsDecorative: chartSpec.IsDecorative,
			}
			chartIdx++
		}
		specs = append(specs, spec)
	}
	return specs, nil
}

func mapToXMLCrop(crop shapes.ImageCrop) *pptxxml.ImageCropRef {
	const cropScaleFactor = 100000
	if crop == (shapes.ImageCrop{}) {
		return nil
	}
	return &pptxxml.ImageCropRef{
		Left:   int64(crop.Left * cropScaleFactor),
		Right:  int64(crop.Right * cropScaleFactor),
		Top:    int64(crop.Top * cropScaleFactor),
		Bottom: int64(crop.Bottom * cropScaleFactor),
	}
}
