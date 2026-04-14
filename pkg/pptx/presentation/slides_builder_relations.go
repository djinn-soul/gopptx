package presentation

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

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
