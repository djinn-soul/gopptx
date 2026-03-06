package charts

import "github.com/djinn-soul/gopptx/internal/pptxxml"

const (
	DataLabelPositionCenter     = "ctr"
	DataLabelPositionInsideEnd  = "inEnd"
	DataLabelPositionInsideBase = "inBase"
	DataLabelPositionOutsideEnd = "outEnd"
	DataLabelPositionBestFit    = "bestFit"
	DataLabelPositionLeft       = "l"
	DataLabelPositionRight      = "r"
	DataLabelPositionTop        = "t"
	DataLabelPositionBottom     = "b"
)

// DataLabelSettings controls chart data-label formatting.
type DataLabelSettings struct {
	UseCustom      bool
	Position       string
	ShowLegendKey  bool
	ShowValue      bool
	ShowCategory   bool
	ShowSeriesName bool
	ShowPercent    bool
	ShowBubbleSize bool
}

func applyDataLabelSettings(spec *pptxxml.ChartSpec, settings DataLabelSettings) {
	spec.DataLabelPosition = settings.Position
	if !settings.UseCustom {
		return
	}
	spec.DataLabelShowLegendKey = boolPtr(settings.ShowLegendKey)
	spec.DataLabelShowValue = boolPtr(settings.ShowValue)
	spec.DataLabelShowCategoryName = boolPtr(settings.ShowCategory)
	spec.DataLabelShowSeriesName = boolPtr(settings.ShowSeriesName)
	spec.DataLabelShowPercent = boolPtr(settings.ShowPercent)
	spec.DataLabelShowBubbleSize = boolPtr(settings.ShowBubbleSize)
}

func boolPtr(v bool) *bool {
	value := v
	return &value
}
