package editor

import (
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodchart "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/chart"
)

// UpdateChartFormatting applies a partial formatting patch to an existing chart.
func (e *PresentationEditor) UpdateChartFormatting(
	slideIndex int,
	selector common.ChartSelector,
	req common.ChartFormatUpdate,
) error {
	refs, err := e.ListSlideCharts(slideIndex)
	if err != nil {
		return err
	}
	chartRef, err := editormodchart.ResolveChartSelector(refs, selector, slideIndex)
	if err != nil {
		return err
	}
	if err := editormodchart.ValidateChartFormatUpdate(req); err != nil {
		return err
	}

	chartXML, ok := e.parts.Get(chartRef.ChartPart)
	if !ok {
		return fmt.Errorf("chart part %s not found", chartRef.ChartPart)
	}
	patched, err := editormodchart.PatchChartFormatting(chartXML, req)
	if err != nil {
		return err
	}
	e.parts.Set(chartRef.ChartPart, patched)
	return nil
}
