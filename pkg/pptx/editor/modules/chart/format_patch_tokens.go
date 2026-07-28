package chart

const (
	chartElementClosePrefix      = "</c:"
	axisTickLabelTagPrefix       = "<c:tickLblPos"
	chartExtensionListTagPrefix  = "<c:extLst"
	dataLabelPositionCenterInput = "center"
	dataLabelPositionCenter      = "ctr"
	dataLabelPositionInsideEnd   = "inEnd"
	dataLabelPositionInsideBase  = "inBase"
	dataLabelPositionOutsideEnd  = "outEnd"
	dataLabelPositionBestFit     = "bestFit"

	// Series-child tags, in the CT_Ser order the patch helpers anchor against.
	seriesDataLabelsTag = "<c:dLbls"
	seriesTrendlineTag  = "<c:trendline"
	seriesErrorBarsTag  = "<c:errBars"
	seriesCategoryTag   = "<c:cat>"
	seriesXValuesTag    = "<c:xVal>"
	seriesValuesTag     = "<c:val>"
	seriesSmoothTag     = "<c:smooth"
	seriesCloseTag      = "</c:ser>"
)
