package editor

import (
	"fmt"
	"path"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodchart "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/chart"
)

func buildPlaceholderChartFrame(
	e *PresentationEditor,
	slideIndex int,
	payload map[string]any,
) (*pptxxml.ChartFrame, bool, error) {
	rawChart, ok := payload["chart"]
	if !ok {
		return nil, false, nil
	}
	chartMap, ok := rawChart.(map[string]any)
	if !ok {
		return nil, false, NewBridgeError(ErrCodeInvalidPayload, "chart must be an object")
	}

	v := NewPayloadValidator()
	chartType, ok := v.RequireString(chartMap, "chart_type")
	if !ok {
		return nil, false, v.Error()
	}
	title := v.OptionalString(chartMap, "title")
	x, _ := v.OptionalInt64(chartMap, "x")
	y, _ := v.OptionalInt64(chartMap, "y")
	w, _ := v.OptionalInt64(chartMap, "w")
	h, _ := v.OptionalInt64(chartMap, "h")

	chartDef, err := editormodchart.PlaceholderChartDefinition(chartMap, chartType, title, x, y, w, h)
	if err != nil {
		return nil, false, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	frame, err := e.createPlaceholderChartPart(slideIndex, chartDef)
	if err != nil {
		return nil, false, err
	}
	return frame, true, nil
}

func (e *PresentationEditor) createPlaceholderChartPart(
	slideIndex int,
	chartDef charts.ChartDefinition,
) (*pptxxml.ChartFrame, error) {
	slidePart := e.slides[slideIndex].Part
	excelData, err := placeholderChartEmbeddingData(chartDef)
	if err != nil {
		return nil, fmt.Errorf("generate excel: %w", err)
	}

	chartNum := e.nextChartNum
	e.nextChartNum++
	chartPartPath := fmt.Sprintf("ppt/charts/chart%d.xml", chartNum)

	excelPartPath, err := e.registerExcelEmbedding(excelData)
	if err != nil {
		return nil, fmt.Errorf("register excel embedding: %w", err)
	}

	e.parts.Set(chartPartPath, nil)
	e.parts.Set(excelPartPath, excelData)
	e.addContentTypeOverride(chartPartPath, "application/vnd.openxmlformats-officedocument.drawingml.chart+xml")
	e.addContentTypeOverride(excelPartPath, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	slideRelID, err := e.nextSlideRelID(slidePart)
	if err != nil {
		return nil, fmt.Errorf("allocate slide relationship id: %w", err)
	}
	if err := e.addSlideRelationship(
		slidePart,
		slideRelID,
		common.RelTypeChart,
		"../charts/"+path.Base(chartPartPath),
	); err != nil {
		return nil, fmt.Errorf("add slide rel: %w", err)
	}

	chartRelID, err := e.allocChartRelID(chartPartPath)
	if err != nil {
		return nil, fmt.Errorf("allocate chart relationship id: %w", err)
	}
	if err := e.addRelationship(
		chartPartPath,
		chartRelID,
		common.RelTypePackage,
		"../embeddings/"+path.Base(excelPartPath),
	); err != nil {
		return nil, fmt.Errorf("add chart rel: %w", err)
	}

	chartSpec := chartDef.ToChartSpec()
	chartSpec.ExternalDataID = chartRelID
	e.parts.Set(chartPartPath, pptxxml.RenderChart(chartSpec))
	e.chartEmbeddings[chartPartPath] = excelPartPath

	return &pptxxml.ChartFrame{
		RelID: slideRelID,
		X:     chartSpec.X,
		Y:     chartSpec.Y,
		CX:    chartSpec.CX,
		CY:    chartSpec.CY,
	}, nil
}

func placeholderChartEmbeddingData(chartDef charts.ChartDefinition) ([]byte, error) {
	switch c := chartDef.(type) {
	case charts.ScatterChart:
		return editormodchart.GenerateExcelForChartUpdate(
			editormodchart.KindScatter,
			common.ChartDataUpdate{
				Series: []common.ChartSeriesData{
					{
						XValues: c.XValues,
						YValues: c.YValues,
					},
				},
			},
		)
	case *charts.ScatterChart:
		return editormodchart.GenerateExcelForChartUpdate(
			editormodchart.KindScatter,
			common.ChartDataUpdate{
				Series: []common.ChartSeriesData{
					{
						XValues: c.XValues,
						YValues: c.YValues,
					},
				},
			},
		)
	case charts.BubbleChart:
		return editormodchart.GenerateExcelForChartUpdate(
			editormodchart.KindBubble,
			common.ChartDataUpdate{
				Series: []common.ChartSeriesData{
					{
						XValues: c.XValues,
						YValues: c.YValues,
						Sizes:   c.BubbleSizes,
					},
				},
			},
		)
	case *charts.BubbleChart:
		return editormodchart.GenerateExcelForChartUpdate(
			editormodchart.KindBubble,
			common.ChartDataUpdate{
				Series: []common.ChartSeriesData{
					{
						XValues: c.XValues,
						YValues: c.YValues,
						Sizes:   c.BubbleSizes,
					},
				},
			},
		)
	default:
		return editormodchart.GenerateExcelForChart(chartDef.GetCategories(), chartDef.GetValues())
	}
}
