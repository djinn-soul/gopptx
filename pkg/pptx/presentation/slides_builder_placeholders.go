package presentation

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

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
