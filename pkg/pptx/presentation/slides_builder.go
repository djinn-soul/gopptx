package presentation

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
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

func (b *slidePartBuilder) nextRID() string {
	rid := fmt.Sprintf("rId%d", b.ridNext)
	b.ridNext++
	return rid
}
