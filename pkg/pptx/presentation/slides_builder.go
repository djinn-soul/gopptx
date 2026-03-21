package presentation

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

const firstSlideRelIDNumber, rotationEmuFactor = 2, 60000

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
		ridNext: firstSlideRelIDNumber, // rId1 is always slideLayout
		targets: make([]string, 0),
	}
}

func (b *slidePartBuilder) build(
	idx int,
	slide elements.SlideContent,
	chartBySlide map[int][]ChartPart,
	smartArtBySlide map[int][]SmartArtPart,
) (*slideParts, error) {
	p, err := b.buildBaseSlideParts(slide)
	if err != nil {
		return nil, err
	}

	b.addSmartArtParts(idx, smartArtBySlide, p)
	if err := b.addChartParts(idx, slide, chartBySlide, p); err != nil {
		return nil, err
	}
	if err := validatePlaceholderChartTargets(p.placeholderChartRels, b.num); err != nil {
		return nil, err
	}

	return p, nil
}

func (b *slidePartBuilder) buildBaseSlideParts(slide elements.SlideContent) (*slideParts, error) {
	p := &slideParts{
		title:        b.buildTitleSpec(slide),
		contentStyle: b.buildContentStyleSpec(slide),
		// Placeholder chart RIDs must be reserved before placeholder media is mapped
		// so media RIDs remain contiguous from rId2 in relationship output.
		placeholderChartRels: b.allocatePlaceholderChartRels(slide.PlaceholderOverrides),
		backgroundRID:        b.mapBackground(slide.Background),
		smartArtFrames:       make([]pptxxml.SmartArtFrame, 0),
		smartArtRels:         make([]pptxxml.SmartArtRel, 0),
		transitionXML:        b.mapTransition(slide),
	}

	if err := b.addSlideTable(slide, p); err != nil {
		return nil, err
	}
	if err := b.addSlideImages(slide, p); err != nil {
		return nil, err
	}
	placeholders, err := b.mapPlaceholders(slide.PlaceholderOverrides, p.placeholderChartRels)
	if err != nil {
		return nil, err
	}
	p.placeholders = placeholders
	return p, nil
}

func (b *slidePartBuilder) addSlideTable(slide elements.SlideContent, p *slideParts) error {
	if slide.Table == nil {
		return nil
	}
	spec, err := slide.Table.ToTableSpec(b.num)
	if err != nil {
		return err
	}
	p.table = spec
	return nil
}

func (b *slidePartBuilder) addSlideImages(slide elements.SlideContent, p *slideParts) error {
	imageRefs, err := b.mapImages(slide.Images)
	if err != nil {
		return err
	}
	p.imageRefs = imageRefs
	return nil
}

func (b *slidePartBuilder) addSmartArtParts(
	idx int,
	smartArtBySlide map[int][]SmartArtPart,
	p *slideParts,
) {
	saList, ok := smartArtBySlide[idx]
	if !ok {
		return
	}
	for _, saPart := range saList {
		b.appendSmartArtPart(p, saPart)
	}
}

func (b *slidePartBuilder) appendSmartArtPart(p *slideParts, saPart SmartArtPart) {
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

func (b *slidePartBuilder) addChartParts(
	idx int,
	slide elements.SlideContent,
	chartBySlide map[int][]ChartPart,
	p *slideParts,
) error {
	chartList, ok := chartBySlide[idx]
	if !ok || len(chartList) == 0 {
		return nil
	}

	listIdx := 0
	if slideChartKindDefined(slide) {
		b.assignPrimaryChart(chartList[0], p)
		listIdx = 1
	}
	return b.assignPlaceholderCharts(chartList, slide, listIdx, p)
}

func (b *slidePartBuilder) assignPrimaryChart(chartPart ChartPart, p *slideParts) {
	rid := b.nextRID()
	p.chartRel = &pptxxml.ChartRel{
		RID:    rid,
		Target: fmt.Sprintf("../charts/chart%d.xml", chartPart.partNumber),
	}
	p.chartFrame = &pptxxml.ChartFrame{
		RelID:        rid,
		X:            chartPart.spec.X,
		Y:            chartPart.spec.Y,
		CX:           chartPart.spec.CX,
		CY:           chartPart.spec.CY,
		AltText:      chartPart.spec.AltText,
		IsDecorative: chartPart.spec.IsDecorative,
	}
}

func (b *slidePartBuilder) assignPlaceholderCharts(
	chartList []ChartPart,
	slide elements.SlideContent,
	start int,
	p *slideParts,
) error {
	for listIdx := start; listIdx < len(chartList); listIdx++ {
		placeholderChartIdx := placeholderChartIndex(slide, listIdx)
		if placeholderChartIdx >= len(p.placeholderChartRels) {
			return fmt.Errorf("slide %d: missing placeholder chart relationship slot", b.num)
		}
		p.placeholderChartRels[placeholderChartIdx].Target = fmt.Sprintf(
			"../charts/chart%d.xml",
			chartList[listIdx].partNumber,
		)
	}
	return nil
}

func placeholderChartIndex(slide elements.SlideContent, listIdx int) int {
	if slideChartKindDefined(slide) {
		return listIdx - 1
	}
	return listIdx
}

func validatePlaceholderChartTargets(rels []pptxxml.ChartRel, slideNumber int) error {
	for _, rel := range rels {
		if strings.TrimSpace(rel.Target) == "" {
			return fmt.Errorf("slide %d: missing chart part for placeholder chart", slideNumber)
		}
	}
	return nil
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
			Rotation:     int64(img.Rotation * rotationEmuFactor),
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

func (b *slidePartBuilder) mapPlaceholders(
	overrides []shapes.PlaceholderContent,
	placeholderChartRels []pptxxml.ChartRel,
) ([]pptxxml.PlaceholderOverrideSpec, error) {
	merged := mergePlaceholderOverrides(overrides)
	specs := make([]pptxxml.PlaceholderOverrideSpec, 0, len(merged))
	chartIdx := 0
	for _, o := range merged {
		spec, nextChartIdx, err := b.mapPlaceholderOverrideSpec(o, placeholderChartRels, chartIdx)
		if err != nil {
			return nil, err
		}
		chartIdx = nextChartIdx
		specs = append(specs, spec)
	}
	return specs, nil
}

func (b *slidePartBuilder) mapPlaceholderOverrideSpec(
	o shapes.PlaceholderContent,
	placeholderChartRels []pptxxml.ChartRel,
	chartIdx int,
) (pptxxml.PlaceholderOverrideSpec, int, error) {
	if err := validatePlaceholderTarget(o); err != nil {
		return pptxxml.PlaceholderOverrideSpec{}, chartIdx, err
	}
	spec := buildPlaceholderBaseSpec(o)
	b.applyPlaceholderImage(&spec, o)
	if err := b.applyPlaceholderTable(&spec, o); err != nil {
		return pptxxml.PlaceholderOverrideSpec{}, chartIdx, err
	}
	nextChartIdx := b.applyPlaceholderChart(&spec, o, placeholderChartRels, chartIdx)
	return spec, nextChartIdx, nil
}

func validatePlaceholderTarget(o shapes.PlaceholderContent) error {
	if o.Target == nil {
		return nil
	}
	if o.Target.Type == "" && o.Target.Index == 0 && o.Target.Name != "" {
		return fmt.Errorf("name-only target %q is not supported for create-path overrides", o.Target.Name)
	}
	return nil
}

func buildPlaceholderBaseSpec(o shapes.PlaceholderContent) pptxxml.PlaceholderOverrideSpec {
	spec := pptxxml.PlaceholderOverrideSpec{
		Index: o.Index,
		Type:  o.Type,
		Text:  o.Text,
	}
	if o.Override == nil {
		return spec
	}
	spec.X = mapOptionalLength(o.Override.X)
	spec.Y = mapOptionalLength(o.Override.Y)
	spec.CX = mapOptionalLength(o.Override.CX)
	spec.CY = mapOptionalLength(o.Override.CY)
	spec.TextStyle = mapPlaceholderTextStyle(o.Override.TextStyle)
	spec.ForceRectGeometry = o.Override.ForceRect
	return spec
}

func (b *slidePartBuilder) applyPlaceholderImage(spec *pptxxml.PlaceholderOverrideSpec, o shapes.PlaceholderContent) {
	if o.Image == nil {
		return
	}
	mediaName, ok := b.catalog.MediaNameForImage(*o.Image)
	if !ok {
		return
	}
	rid := b.nextRID()
	spec.Image = &pptxxml.ImageRef{
		RelID:      rid,
		Name:       "Placeholder Picture",
		X:          o.Image.X.Emu(),
		Y:          o.Image.Y.Emu(),
		CX:         o.Image.CX.Emu(),
		CY:         o.Image.CY.Emu(),
		Rotation:   int64(o.Image.Rotation * rotationEmuFactor),
		FlipH:      o.Image.FlipH,
		FlipV:      o.Image.FlipV,
		Shadow:     o.Image.Shadow,
		Reflection: o.Image.Reflection,
		Crop:       mapToXMLCrop(o.Image.Crop),
	}
	b.targets = append(b.targets, fmt.Sprintf("../media/%s", mediaName))
}

func (b *slidePartBuilder) applyPlaceholderTable(
	spec *pptxxml.PlaceholderOverrideSpec,
	o shapes.PlaceholderContent,
) error {
	if o.Table == nil {
		return nil
	}
	tableSpec, err := o.Table.ToTableSpec(b.num)
	if err != nil {
		return err
	}
	spec.Table = tableSpec
	return nil
}

func (b *slidePartBuilder) applyPlaceholderChart(
	spec *pptxxml.PlaceholderOverrideSpec,
	o shapes.PlaceholderContent,
	placeholderChartRels []pptxxml.ChartRel,
	chartIdx int,
) int {
	if o.Chart == nil || chartIdx >= len(placeholderChartRels) {
		return chartIdx
	}
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
	return chartIdx + 1
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
