package chart

import (
	"encoding/xml"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type dataLabelValueXML struct {
	Value string `xml:"val,attr"`
}

type dataLabelNumberFormatXML struct {
	FormatCode   string `xml:"formatCode,attr"`
	SourceLinked string `xml:"sourceLinked,attr"`
}

type dataLabelBodyPropertiesXML struct {
	Wrap string `xml:"wrap,attr"`
}

type dataLabelTextPropertiesXML struct {
	BodyProperties dataLabelBodyPropertiesXML `xml:"bodyPr"`
}

type dataLabelsXML struct {
	Position       dataLabelValueXML           `xml:"dLblPos"`
	ShowValue      dataLabelValueXML           `xml:"showVal"`
	ShowCategory   dataLabelValueXML           `xml:"showCatName"`
	ShowSeriesName dataLabelValueXML           `xml:"showSerName"`
	NumberFormat   *dataLabelNumberFormatXML   `xml:"numFmt"`
	TextProperties *dataLabelTextPropertiesXML `xml:"txPr"`
}

type chartPlotXML struct {
	XMLName    xml.Name
	DataLabels *dataLabelsXML `xml:"dLbls"`
}

type chartStateDocumentXML struct {
	Chart struct {
		PlotArea struct {
			Children []chartPlotXML `xml:",any"`
		} `xml:"plotArea"`
	} `xml:"chart"`
}

func parseDataLabelState(chartXML string) common.ChartDataLabelState {
	labels, found := firstPlotDataLabels(chartXML)
	if !found {
		return common.ChartDataLabelState{}
	}
	state := common.ChartDataLabelState{
		Present:        true,
		Position:       strings.TrimSpace(labels.Position.Value),
		ShowValue:      xmlBool(labels.ShowValue),
		ShowCategory:   xmlBool(labels.ShowCategory),
		ShowSeriesName: xmlBool(labels.ShowSeriesName),
	}
	if labels.NumberFormat != nil {
		state.NumberFormat = labels.NumberFormat.FormatCode
		linked := strings.TrimSpace(labels.NumberFormat.SourceLinked) == "1"
		state.FormatLinked = &linked
	}
	if labels.TextProperties != nil {
		switch strings.TrimSpace(labels.TextProperties.BodyProperties.Wrap) {
		case xmlValueNone:
			wordWrap := false
			state.WordWrap = &wordWrap
		case dataLabelWrapSquare:
			wordWrap := true
			state.WordWrap = &wordWrap
		}
	}
	return state
}

func firstPlotDataLabels(chartXML string) (dataLabelsXML, bool) {
	var document chartStateDocumentXML
	//nolint:musttag // Arbitrary plot names require an untagged XMLName.
	if err := xml.Unmarshal([]byte(chartXML), &document); err != nil {
		return dataLabelsXML{}, false
	}
	for _, child := range document.Chart.PlotArea.Children {
		if isChartPlotElement(child.XMLName.Local) {
			if child.DataLabels == nil {
				return dataLabelsXML{}, false
			}
			return *child.DataLabels, true
		}
	}
	return dataLabelsXML{}, false
}

func isChartPlotElement(name string) bool {
	switch name {
	case "barChart", "lineChart", "areaChart", "pieChart", "pie3DChart", "doughnutChart",
		"scatterChart", "bubbleChart", "radarChart", "stockChart":
		return true
	default:
		return false
	}
}

func xmlBool(value dataLabelValueXML) bool {
	return strings.TrimSpace(value.Value) == "1"
}
